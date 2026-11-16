import { useEffect, useRef, useState } from 'react'
import type { MusicPlan } from '../bindings/ifly/internal/music/models'
import type { SemanticAnalysis } from '../bindings/ifly/internal/semantic/models'
import type { StructuralAnalysis } from '../bindings/ifly/internal/analysis/models'
import { Service as PresenceService } from '../bindings/ifly/internal/presence'
import { ImageData, loadImage } from './imageImport'
import { analyzePreview } from './structuralAnalysis'
import { createEdgePreview } from './edgePreview'
import { cancelInterpretation, interpretPreview } from './semanticAnalysis'
import { composePreview } from './musicPreview'
import { MusicPlayer } from './musicPlayer'
import { INSTRUMENT_LABELS } from './musicInstruments'
import { TIMBRES, type MusicVoice } from './musicTimbres'
import { createTempoMap } from './musicTiming'
import './app.css'

const MOOD_LABELS: Record<string, string> = { calm: 'Calme', tense: 'Tendue', energetic: 'Énergique', neutral: 'Neutre' }
const MODE_LABELS: Record<string, string> = { pentatonic: 'pentatonique', minor: 'mineur', major: 'majeur', dorian: 'dorien' }
const TIMBRE_LABELS: Record<string, string> = { warm: 'chaleureux', dark: 'sombre', bright: 'clair', soft: 'doux' }

// percent formats a normalized metric for the diagnostic panel.
function percent(value: number): string {
  return `${Math.round(value * 100)} %`
}

// describeMusicForm takes a generated plan and returns its phrase labels and approximate loop duration in French.
function describeMusicForm(plan: MusicPlan): string {
  const seconds = Math.round(createTempoMap(plan.sections, plan.beatsPerLoop).loopSeconds)
  const duration = seconds >= 60 ? `${Math.floor(seconds / 60)} min ${String(seconds % 60).padStart(2, '0')} s` : `${seconds} s`
  return `${plan.sections?.map(section => section.label).join(' → ') ?? 'Forme indisponible'} · ${duration}`
}

// describeTempoRange takes a score and returns the slowest and fastest phrase tempos for its inspector summary.
function describeTempoRange(plan: MusicPlan): string {
  const tempi = plan.sections?.map(section => section.tempoBpm) ?? [plan.tempoBpm]
  return `${Math.min(...tempi)}–${Math.max(...tempi)} BPM`
}

// describeInstrumentPalette takes a score and returns the names of only the instruments actually present in its notes.
function describeInstrumentPalette(plan: MusicPlan): string {
  const profile = TIMBRES[plan.timbre]
  if (!profile) return 'Instruments indisponibles'
  const names = new Set<string>()
  for (const note of plan.notes ?? []) {
    const voice = note.voice as MusicVoice
    const instrument = profile.voices[voice]?.instrument
    if (instrument) names.add(INSTRUMENT_LABELS[instrument])
    else if (voice === 'pulse') names.add('Percussion légère')
  }
  return [...names].join(' · ')
}

// App owns the temporary artwork, its asynchronous analysis, and disposable preview resources.
function App() {
  const [image, setImage] = useState<ImageData | null>(null)
  const [analysis, setAnalysis] = useState<StructuralAnalysis | null>(null)
  const [semantic, setSemantic] = useState<SemanticAnalysis | null>(null)
  const [semanticError, setSemanticError] = useState('')
  const [musicPlan, setMusicPlan] = useState<MusicPlan | null>(null)
  const [musicError, setMusicError] = useState('')
  const [playback, setPlayback] = useState<'stopped' | 'playing' | 'paused'>('stopped')
  const [audioLoading, setAudioLoading] = useState(false)
  const [edgeUrl, setEdgeUrl] = useState<string | null>(null)
  const [showEdges, setShowEdges] = useState(false)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState<'import' | 'analysis' | 'semantic' | 'music' | 'edges' | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const imageRef = useRef<ImageData | null>(null)
  const edgeRef = useRef<string | null>(null)
  const playerRef = useRef<MusicPlayer | null>(null)
  const playerAction = useRef(false)
  const generation = useRef(0)

  useEffect(() => {
    void PresenceService.SetArtworkOpen(image !== null).catch(() => undefined)
  }, [image])

  // discardMusic stops the current preview and clears the score when its source image or interpretation changes.
  function discardMusic(): void {
    playerRef.current?.stop()
    setPlayback('stopped')
    setAudioLoading(false)
    setMusicPlan(null)
    setMusicError('')
  }

  // importFile validates a new image and discards results from any earlier import or analysis.
  async function importFile(file: File) {
    const current = ++generation.current
    void cancelInterpretation().catch(() => undefined)
    discardMusic()
    setBusy('import')
    setError('')
    setAnalysis(null)
    setSemantic(null)
    setSemanticError('')
    setShowEdges(false)
    setEdgeUrl(null)
    if (edgeRef.current) URL.revokeObjectURL(edgeRef.current)
    edgeRef.current = null
    try {
      const next = await loadImage(file)
      if (current !== generation.current) {
        URL.revokeObjectURL(next.url)
        return
      }
      if (imageRef.current) URL.revokeObjectURL(imageRef.current.url)
      imageRef.current = next
      setImage(next)
    } catch (cause) {
      if (current === generation.current) setError(cause instanceof Error ? cause.message : 'Impossible d’ouvrir cette image.')
    } finally {
      if (current === generation.current) setBusy(null)
    }
  }

  // handleSelection imports the chosen file and resets the input so that it can be selected again.
  function handleSelection(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0]
    event.currentTarget.value = ''
    if (file) void importFile(file)
  }

  // handleDrop accepts one local image without navigating the application webview.
  function handleDrop(event: React.DragEvent<HTMLElement>) {
    event.preventDefault()
    const file = event.dataTransfer.files[0]
    if (file) void importFile(file)
  }

  // runAnalysis measures the current preview in Go and ignores results for an image that was replaced.
  async function runAnalysis() {
    if (!imageRef.current || busy) return
    const current = generation.current
    discardMusic()
    setBusy('analysis')
    setError('')
    setSemantic(null)
    setSemanticError('')
    try {
      const result = await analyzePreview(imageRef.current.url)
      if (current !== generation.current) return
      setAnalysis(result)
      await runSemanticAnalysis()
    } catch (cause) {
      if (current === generation.current) setError(cause instanceof Error ? cause.message : 'Analyse locale impossible.')
    } finally {
      if (current === generation.current) setBusy(null)
    }
  }

  // runSemanticAnalysis requests a local VLM description and keeps failures separate from measured data.
  async function runSemanticAnalysis() {
    if (!imageRef.current) return
    const current = generation.current
    discardMusic()
    setBusy('semantic')
    setSemanticError('')
    try {
      const result = await interpretPreview(imageRef.current.url)
      if (current !== generation.current) return
      setSemantic(result)
      setBusy('music')
      try {
        const plan = await composePreview(imageRef.current.url, result)
        if (current === generation.current) setMusicPlan(plan)
      } catch (cause) {
        if (current === generation.current) setMusicError(cause instanceof Error ? cause.message : 'Plan musical indisponible.')
      }
    } catch (cause) {
      if (current === generation.current) {
        const message = cause instanceof Error ? cause.message : ''
        setSemanticError(message.includes('unavailable') || message.includes('connect')
          ? 'Modèle local indisponible. Exécutez scripts/setup-local-vlm.ps1 une fois, puis relancez IFLY. Les mesures visuelles restent utilisables.'
          : 'Interprétation locale impossible. Vérifiez le modèle et réessayez.')
      }
    } finally {
      if (current === generation.current) setBusy(null)
    }
  }

  // togglePlayback takes the prepared score and starts, suspends, or resumes local audio from the user's click.
  async function togglePlayback() {
    if (!musicPlan || busy || playerAction.current) return
    playerAction.current = true
    const current = generation.current
    const player = playerRef.current ?? new MusicPlayer()
    playerRef.current = player
    try {
      if (playback === 'playing') {
        await player.pause()
        if (current === generation.current) setPlayback('paused')
      } else if (playback === 'paused') {
        await player.resume()
        if (current === generation.current) setPlayback('playing')
      } else {
        setAudioLoading(true)
        await player.start(musicPlan)
        if (current === generation.current) setPlayback('playing')
      }
      if (current === generation.current) setMusicError('')
    } catch (cause) {
      player.stop()
      if (current === generation.current) {
        setPlayback('stopped')
        setMusicError(cause instanceof Error ? cause.message : 'Lecture audio impossible.')
      }
    } finally {
      if (current === generation.current) setAudioLoading(false)
      playerAction.current = false
    }
  }

  // toggleEdges alternates between the source image and a locally derived contour map.
  async function toggleEdges() {
    if (showEdges) {
      setShowEdges(false)
      return
    }
    if (edgeUrl) {
      setShowEdges(true)
      return
    }
    if (!imageRef.current || !analysis || busy) return
    const current = generation.current
    setBusy('edges')
    setError('')
    try {
      const next = await createEdgePreview(imageRef.current.url)
      if (current !== generation.current) {
        URL.revokeObjectURL(next)
        return
      }
      edgeRef.current = next
      setEdgeUrl(next)
      setShowEdges(true)
    } catch (cause) {
      if (current === generation.current) setError(cause instanceof Error ? cause.message : 'Impossible de préparer les contours.')
    } finally {
      if (current === generation.current) setBusy(null)
    }
  }

  useEffect(() => () => {
    generation.current++
    void cancelInterpretation().catch(() => undefined)
    playerRef.current?.stop()
    if (imageRef.current) URL.revokeObjectURL(imageRef.current.url)
    if (edgeRef.current) URL.revokeObjectURL(edgeRef.current)
  }, [])

  return <main className="app-shell" onDragOver={event => event.preventDefault()} onDrop={handleDrop}>
    <header className="app-header">
      <div className="brand">IFLY<span>.</span></div>
      <span className="header-divider" />
      <span className="header-title">Atelier d’image</span>
      <span className="header-file">{image?.fileName ?? 'Aucune œuvre ouverte'}</span>
      <span className="header-mode">LABORATOIRE B</span>
    </header>

    <aside className="tool-rail" aria-label="Outils">
      <button className="rail-tool selected" title="Importer une image" aria-label="Importer une image" onClick={() => inputRef.current?.click()}>▣<small>Import</small></button>
      <button className="rail-tool" title="Analyser l’image" aria-label="Analyser l’image" onClick={() => void runAnalysis()} disabled={!image || busy !== null}>◈<small>Analyse</small></button>
      <div className="rail-divider" />
      <span className="rail-hint" title="Les outils de dessin arriveront dans le prototype A">Outils de dessin à venir</span>
    </aside>

    <section className="stage" aria-label="Œuvre importée">
      <div className="stage-toolbar">
        <div className="toolbar-actions">
          <button className="button primary" onClick={() => inputRef.current?.click()} disabled={busy === 'import'}>{image ? 'Changer l’image' : 'Importer une image'}</button>
          <button className="button" onClick={() => void runAnalysis()} disabled={!image || busy !== null}>{busy === 'analysis' ? 'Analyse en cours…' : 'Analyser'}</button>
          <input ref={inputRef} type="file" accept=".png,.jpg,.jpeg,.webp,image/png,image/jpeg,image/webp" onChange={handleSelection} className="visually-hidden" aria-label="Choisir une image" />
        </div>
        <button className="play-state" onClick={() => void togglePlayback()} disabled={!musicPlan || busy !== null || audioLoading} title={musicPlan ? 'Écouter ou mettre en pause la composition' : 'Analysez une image pour préparer la musique'} aria-label={audioLoading ? 'Chargement des instruments' : playback === 'playing' ? 'Mettre la musique en pause' : playback === 'paused' ? 'Reprendre la musique' : 'Écouter la composition'}>{audioLoading ? '…' : playback === 'playing' ? '❚❚' : '▶'} <span>{audioLoading ? 'Chargement…' : playback === 'playing' ? 'Pause' : playback === 'paused' ? 'Reprendre' : 'Écouter'}</span></button>
      </div>
      <div className="canvas-area">
        {image ? <img src={showEdges && edgeUrl ? edgeUrl : image.url} alt={showEdges ? `Carte des contours de ${image.fileName}` : `Œuvre importée : ${image.fileName}`} /> : <div className="canvas-empty"><div className="canvas-mark">✳</div><strong>Déposez une image sur le plan de travail</strong><span>PNG, JPEG ou WebP · 20 Mo maximum</span></div>}
      </div>
      <div className="status-bar"><span>{image ? `${image.width} × ${image.height} px · ${image.format}` : 'Aucune image'}</span><span>{audioLoading ? 'Chargement des instruments' : playback === 'playing' ? 'Composition instrumentale en cours' : playback === 'paused' ? 'Musique en pause' : showEdges ? 'Vue des contours' : busy === 'analysis' ? 'Mesures en cours' : busy === 'semantic' ? 'Interprétation locale en cours' : busy === 'music' ? 'Composition en cours' : musicPlan ? 'Composition instrumentale prête' : semantic ? 'Interprétation expérimentale disponible' : analysis ? 'Analyse structurelle terminée' : 'Prêt'}</span></div>
    </section>

    <aside className="inspector" aria-label="Analyse et propriétés">
      <div className="inspector-tabs"><strong>PROPRIÉTÉS</strong><span>ANALYSE</span></div>
      <section className="inspector-section">
        <h2>Document</h2>
        {image ? <dl className="document-facts"><div><dt>Nom</dt><dd title={image.fileName}>{image.fileName}</dd></div><div><dt>Dimensions</dt><dd>{image.width} × {image.height}</dd></div><div><dt>Taille</dt><dd>{(image.byteSize / 1024 / 1024).toFixed(2)} Mo</dd></div></dl> : <p className="hint">Importez une image pour commencer.</p>}
      </section>
      <section className="inspector-section">
        <h2>Structure visuelle</h2>
        {error && <p className="error" role="alert">{error}</p>}
        {!analysis && <p className="hint">{image ? 'Cliquez sur Analyser pour mesurer cette image.' : 'Les mesures apparaîtront ici.'}</p>}
        {analysis && <><div className="metric-grid">
          <div><span>Luminosité</span><strong>{percent(analysis.brightness)}</strong></div>
          <div><span>Contraste</span><strong>{percent(analysis.contrast)}</strong></div>
          <div><span>Contours</span><strong>{percent(analysis.edgeDensity)}</strong></div>
          <div><span>Couleurs</span><strong>{percent(analysis.colorDiversity)}</strong></div>
          <div><span>Complexité</span><strong>{percent(analysis.complexity)}</strong></div>
        </div><p className="analysis-note">Mesures locales sans interprétation du sujet ni de l’émotion.</p><button className="edge-button" onClick={() => void toggleEdges()} disabled={busy === 'edges'}>{busy === 'edges' ? 'Préparation…' : showEdges ? 'Revoir l’image' : 'Voir les contours'}</button></>}
      </section>
      <section className="inspector-section next-section">
        <h2>Interprétation locale</h2>
        {semantic && <><p className="semantic-description">{semantic.description}</p><p className="hint">Sujets proposés : {semantic.subjects?.length ? semantic.subjects.join(', ') : 'non identifiés'}</p><p className="analysis-note">{semantic.uncertain ? 'Le modèle signale une incertitude.' : 'Proposition du modèle à vérifier : les croquis ambigus peuvent être mal reconnus.'}</p></>}
        {!semantic && busy === 'semantic' && <p className="hint" role="status">Le modèle local examine l’image…</p>}
        {semanticError && <p className="error" role="alert">{semanticError}</p>}
        {!semantic && !semanticError && busy !== 'semantic' && <p className="hint">{analysis ? 'La description locale sera proposée ici.' : 'Analysez l’image pour commencer.'}</p>}
        {analysis && !semantic && busy === null && <button className="edge-button" onClick={() => void runSemanticAnalysis()}>Réessayer l’interprétation</button>}
        <div className="pipeline"><span className={analysis ? 'done' : ''}>01 · Structure</span><span className={semantic ? 'done' : ''}>02 · Sens expérimental</span><span className={musicPlan ? 'done' : ''}>03 · Musique</span><span>04 · Évolution</span></div>
      </section>
      <section className="inspector-section music-section">
        <h2>Composition instrumentale</h2>
        {musicPlan ? <>
          <p className="semantic-description">{MOOD_LABELS[musicPlan.mood] ?? 'Neutre'} · {describeTempoRange(musicPlan)} · mode {MODE_LABELS[musicPlan.mode] ?? musicPlan.mode}</p>
          <p className="hint">Forme : {describeMusicForm(musicPlan)} · timbre {TIMBRE_LABELS[musicPlan.timbre] ?? musicPlan.timbre}</p>
          <p className="hint">Instruments : {describeInstrumentPalette(musicPlan)}</p>
          <p className="analysis-note">Instruments choisis automatiquement pour cette image. La pièce recommence après sa conclusion.</p>
        </> : <p className="hint">La musique est préparée après l’interprétation.</p>}
        {musicError && <p className="error" role="alert">{musicError}</p>}
      </section>
    </aside>
  </main>
}

export default App
