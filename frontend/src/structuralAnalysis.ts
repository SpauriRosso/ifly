import { AnalysisService } from '../bindings/changeme'
import type { StructuralAnalysis } from '../bindings/changeme/models'

// prepareAnalysisImage downsamples a validated preview into a bounded PNG for local Go analysis.
export async function prepareAnalysisImage(url: string): Promise<string> {
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
    return canvas.toDataURL('image/png')
  } finally {
    bitmap.close()
  }
}

// analyzePreview asks Go for deterministic structural measurements of a bounded local PNG.
export async function analyzePreview(url: string): Promise<StructuralAnalysis> {
  return AnalysisService.Analyze(await prepareAnalysisImage(url))
}
