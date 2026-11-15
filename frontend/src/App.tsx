import { useEffect, useRef, useState } from 'react'
import type { StructuralAnalysis } from '../bindings/changeme/models'
import { ImageData, loadImage } from './imageImport'
import { analyzePreview } from './structuralAnalysis'
import './app.css'

// percent formats a normalized metric for the diagnostic panel.
function percent(value: number): string {
  return `${Math.round(value * 100)} %`
}

// App owns the temporary artwork, its asynchronous analysis, and disposable preview resources.
function App() {
  const [image, setImage] = useState<ImageData | null>(null)
  const [analysis, setAnalysis] = useState<StructuralAnalysis | null>(null)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState<'import' | 'analysis' | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const imageRef = useRef<ImageData | null>(null)
  const generation = useRef(0)

  // importFile validates a new image and discards results from any earlier import or analysis.
  async function importFile(file: File) {
    const current = ++generation.current
    setBusy('import')
    setError('')
    setAnalysis(null)
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

  useEffect(() => () => {
    generation.current++
    if (imageRef.current) URL.revokeObjectURL(imageRef.current.url)
  }, [])

  return <main className="app-shell" onDragOver={event => event.preventDefault()} onDrop={handleDrop}>
    <aside className="tool-rail" aria-label="Étapes du laboratoire">
      <div className="brand">IFLY<span>.</span></div>
      <div className="rail-caption">LABORATOIRE<br />PROTOTYPE B</div>
      <div className="rail-separator" />
      <div className="step active"><span className="step-number">01</span><div><strong>Image</strong><small>Importer et observer</small></div></div>
      <div className={`step ${analysis ? 'active' : ''}`}><span className="step-number">02</span><div><strong>Structure</strong><small>Mesures locales</small></div></div>
      <div className="step muted"><span className="step-number">03</span><div><strong>Sens & musique</strong><small>À construire</small></div></div>
      <div className="rail-bottom">L’œuvre reste modifiable dans ce laboratoire. Aucune validation définitive n’est encore proposée.</div>
    </aside>

    <section className="stage" aria-label="Œuvre importée">
      <header className="stage-toolbar">
        <div className="toolbar-actions">
          <button className="button primary" onClick={() => inputRef.current?.click()} disabled={busy === 'import'}>{image ? 'Changer l’image' : 'Importer une image'}</button>
          <button className="button" onClick={() => void runAnalysis()} disabled={!image || busy !== null}>{busy === 'analysis' ? 'Analyse en cours…' : 'Analyser'}</button>
          <input ref={inputRef} type="file" accept=".png,.jpg,.jpeg,.webp,image/png,image/jpeg,image/webp" onChange={handleSelection} className="visually-hidden" aria-label="Choisir une image" />
        </div>
        <div className="play-state" title="La lecture sera disponible après l’intégration de la composition musicale">▶ &nbsp; Lecture indisponible</div>
      </header>
      <div className="canvas-area">
        {image ? <img src={image.url} alt={`Œuvre importée : ${image.fileName}`} /> : <div className="canvas-empty"><div className="canvas-mark">✳</div><strong>Votre image commence ici</strong><span>Importez un PNG, JPEG ou WebP, ou glissez-le dans cette fenêtre.</span></div>}
      </div>
      <section className="inspection" aria-label="Diagnostic">
        <div className="inspection-heading"><strong>DIAGNOSTIC</strong><span>{image ? image.fileName : 'Aucune image'}</span></div>
        {error && <p className="error" role="alert">{error}</p>}
        {!image && <p className="hint">L’image occupera cet espace. L’analyse structurelle sera disponible après import.</p>}
        {image && !analysis && !error && <p className="hint">{image.format} · {image.width} × {image.height} px · {(image.byteSize / 1024 / 1024).toFixed(2)} Mo — Cliquez sur « Analyser » pour mesurer la structure.</p>}
        {analysis && <><div className="metric-grid">
          <div><span>Luminosité</span><strong>{percent(analysis.brightness)}</strong></div>
          <div><span>Contraste</span><strong>{percent(analysis.contrast)}</strong></div>
          <div><span>Contours</span><strong>{percent(analysis.edgeDensity)}</strong></div>
          <div><span>Couleurs</span><strong>{percent(analysis.colorDiversity)}</strong></div>
          <div><span>Complexité</span><strong>{percent(analysis.complexity)}</strong></div>
        </div><p className="analysis-note">Mesures visuelles locales uniquement. La description du sujet et de l’atmosphère nécessite encore un modèle installé localement.</p></>}
      </section>
    </section>
  </main>
}

export default App
