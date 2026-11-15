import { useEffect, useRef, useState } from 'react'
import type { StructuralAnalysis } from '../bindings/changeme/models'
import { ImageData, loadImage } from './imageImport'
import { analyzePreview } from './structuralAnalysis'
import { createEdgePreview } from './edgePreview'
import './app.css'

// percent formats a normalized metric for the diagnostic panel.
function percent(value: number): string {
  return `${Math.round(value * 100)} %`
}

// App owns the temporary artwork, its asynchronous analysis, and disposable preview resources.
function App() {
  const [image, setImage] = useState<ImageData | null>(null)
  const [analysis, setAnalysis] = useState<StructuralAnalysis | null>(null)
  const [edgeUrl, setEdgeUrl] = useState<string | null>(null)
  const [showEdges, setShowEdges] = useState(false)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState<'import' | 'analysis' | 'edges' | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const imageRef = useRef<ImageData | null>(null)
  const edgeRef = useRef<string | null>(null)
  const generation = useRef(0)

  // importFile validates a new image and discards results from any earlier import or analysis.
  async function importFile(file: File) {
    const current = ++generation.current
    setBusy('import')
    setError('')
    setAnalysis(null)
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
    setBusy('analysis')
    setError('')
    try {
      const result = await analyzePreview(imageRef.current.url)
      if (current === generation.current) setAnalysis(result)
    } catch (cause) {
      if (current === generation.current) setError(cause instanceof Error ? cause.message : 'Analyse locale impossible.')
    } finally {
      if (current === generation.current) setBusy(null)
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
        <div className="play-state" title="La composition musicale sera disponible après l’analyse sémantique locale">▶ <span>Lecture à venir</span></div>
      </div>
      <div className="canvas-area">
        {image ? <img src={showEdges && edgeUrl ? edgeUrl : image.url} alt={showEdges ? `Carte des contours de ${image.fileName}` : `Œuvre importée : ${image.fileName}`} /> : <div className="canvas-empty"><div className="canvas-mark">✳</div><strong>Déposez une image sur le plan de travail</strong><span>PNG, JPEG ou WebP · 20 Mo maximum</span></div>}
      </div>
      <div className="status-bar"><span>{image ? `${image.width} × ${image.height} px · ${image.format}` : 'Aucune image'}</span><span>{showEdges ? 'Vue des contours' : busy === 'analysis' ? 'Analyse locale en cours' : analysis ? 'Analyse structurelle terminée' : 'Prêt'}</span></div>
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
        <h2>Suite de l’expérience</h2>
        <p className="hint">{analysis ? 'Prochaine étape : intégrer et évaluer un modèle visuel local pour comprendre le sujet et l’atmosphère. La musique viendra après ce test.' : 'Analysez l’image pour préparer son interprétation.'}</p>
        <div className="pipeline"><span className={analysis ? 'done' : ''}>01 · Structure</span><span>02 · Sens</span><span>03 · Musique</span><span>04 · Évolution</span></div>
      </section>
    </aside>
  </main>
}

export default App
