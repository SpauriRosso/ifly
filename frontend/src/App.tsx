import { useEffect, useRef, useState } from 'react'
import { ImageData, loadImage } from './imageImport'
import './app.css'

// App manages the import session and releases each preview URL when replaced or unmounted.
function App() {
  const [image, setImage] = useState<ImageData | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const imageRef = useRef<ImageData | null>(null)
  const requestRef = useRef(0)

  // importFile validates and decodes a file, discarding obsolete asynchronous results.
  async function importFile(file: File) {
    const request = ++requestRef.current
    setLoading(true)
    setError('')
    try {
      const next = await loadImage(file)
      if (request !== requestRef.current) {
        URL.revokeObjectURL(next.url)
        return
      }
      if (imageRef.current) URL.revokeObjectURL(imageRef.current.url)
      imageRef.current = next
      setImage(next)
    } catch (cause) {
      if (request === requestRef.current) setError(cause instanceof Error ? cause.message : 'Impossible de lire cette image.')
    } finally {
      if (request === requestRef.current) setLoading(false)
    }
  }

  // handleSelection imports the chosen local file and clears the input for repeated selections.
  function handleSelection(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0]
    event.currentTarget.value = ''
    if (file) void importFile(file)
  }

  // handleDrop accepts a local file without navigating the webview away from the app.
  function handleDrop(event: React.DragEvent<HTMLElement>) {
    event.preventDefault()
    const file = event.dataTransfer.files[0]
    if (file) void importFile(file)
  }

  useEffect(() => () => {
    requestRef.current++
    if (imageRef.current) URL.revokeObjectURL(imageRef.current.url)
  }, [])

  return <main className="workspace" onDragOver={event => event.preventDefault()} onDrop={handleDrop}>
    <aside className="sidebar">
      <div className="wordmark">IFLY<span className="wordmark-dot">.</span></div>
      <div className="sidebar-label">LABORATOIRE · PROTOTYPE B</div>
      <button className="import-button" onClick={() => inputRef.current?.click()} disabled={loading}>
        {loading ? 'Ouverture…' : image ? 'Changer d’image' : 'Importer une image'}
      </button>
      <input ref={inputRef} type="file" accept=".png,.jpg,.jpeg,.webp,image/png,image/jpeg,image/webp" onChange={handleSelection} className="visually-hidden" aria-label="Choisir une image" />
      <p className="sidebar-hint">PNG, JPEG ou WebP · 20 Mo maximum</p>
      <div className="sidebar-spacer" />
      <div className="stage-index">01 / IMPORT</div>
      <p className="sidebar-foot">Une image pour commencer. Son interprétation et son évolution viendront ensuite.</p>
    </aside>
    <section className="main-panel" aria-label="Aperçu de l’œuvre">
      <header className="topbar">
        <div><span className="eyebrow">ŒUVRE EN PRÉPARATION</span><h1>Votre image, avant le temps.</h1></div>
        <span className="status-pill">Import de laboratoire</span>
      </header>
      <div className={`preview ${image ? 'has-image' : ''}`}>
        {image ? <img src={image.url} alt={`Aperçu de ${image.fileName}`} /> :
          <div className="empty-state"><div className="empty-symbol">✦</div><strong>Un espace pour votre image</strong><span>Importez une image ou glissez-la ici.</span></div>}
      </div>
      <section className="diagnostics" aria-label="Diagnostic d’import">
        <div className="diagnostic-title">DIAGNOSTIC <span>01 — FICHIER</span></div>
        {error && <p className="error" role="alert">{error}</p>}
        {image ? <dl className="facts">
          <div><dt>Fichier</dt><dd title={image.fileName}>{image.fileName}</dd></div>
          <div><dt>Format</dt><dd>{image.format}</dd></div>
          <div><dt>Dimensions</dt><dd>{image.width} × {image.height} px</dd></div>
          <div><dt>Taille</dt><dd>{(image.byteSize / 1024 / 1024).toFixed(2)} Mo</dd></div>
        </dl> : <p className="diagnostic-empty">Aucune image chargée. Les caractéristiques apparaîtront ici.</p>}
      </section>
    </section>
  </main>
}

export default App
