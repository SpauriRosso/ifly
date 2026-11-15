import { AnalysisService } from '../bindings/changeme'
import type { StructuralAnalysis } from '../bindings/changeme/models'

// analyzePreview downsamples a validated local image and asks Go for deterministic structural metrics.
export async function analyzePreview(url: string): Promise<StructuralAnalysis> {
  const response = await fetch(url)
  const bitmap = await createImageBitmap(await response.blob(), { imageOrientation: 'from-image' })
  try {
    const ratio = Math.min(1, 256 / Math.max(bitmap.width, bitmap.height))
    const canvas = document.createElement('canvas')
    canvas.width = Math.max(1, Math.round(bitmap.width * ratio))
    canvas.height = Math.max(1, Math.round(bitmap.height * ratio))
    const context = canvas.getContext('2d')
    if (!context) throw new Error('Le rendu de l’image est indisponible.')
    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, canvas.width, canvas.height)
    context.drawImage(bitmap, 0, 0, canvas.width, canvas.height)
    return await AnalysisService.Analyze(canvas.toDataURL('image/png'))
  } finally {
    bitmap.close()
  }
}
