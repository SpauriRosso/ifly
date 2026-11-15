// createEdgePreview derives a disposable grayscale edge-map URL from a local image preview.
export async function createEdgePreview(url: string): Promise<string> {
  const response = await fetch(url)
  const bitmap = await createImageBitmap(await response.blob(), { imageOrientation: 'from-image' })
  try {
    const scale = Math.min(1, 1024 / Math.max(bitmap.width, bitmap.height))
    const canvas = document.createElement('canvas')
    canvas.width = Math.max(1, Math.round(bitmap.width * scale))
    canvas.height = Math.max(1, Math.round(bitmap.height * scale))
    const context = canvas.getContext('2d', { willReadFrequently: true })
    if (!context) throw new Error('La carte des contours est indisponible.')
    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, canvas.width, canvas.height)
    context.drawImage(bitmap, 0, 0, canvas.width, canvas.height)
    const source = context.getImageData(0, 0, canvas.width, canvas.height)
    const output = context.createImageData(canvas.width, canvas.height)
    const gray = new Float32Array(canvas.width * canvas.height)
    for (let i = 0; i < gray.length; i++) {
      const pixel = i * 4
      gray[i] = .2126 * source.data[pixel] + .7152 * source.data[pixel + 1] + .0722 * source.data[pixel + 2]
    }
    for (let y = 0; y < canvas.height; y++) {
      for (let x = 0; x < canvas.width; x++) {
        const i = y * canvas.width + x
        const horizontal = x + 1 < canvas.width ? Math.abs(gray[i] - gray[i + 1]) : 0
        const vertical = y + 1 < canvas.height ? Math.abs(gray[i] - gray[i + canvas.width]) : 0
        const strength = Math.min(255, (horizontal + vertical) * 2)
        const pixel = i * 4
        output.data[pixel] = strength
        output.data[pixel + 1] = strength
        output.data[pixel + 2] = strength
        output.data[pixel + 3] = 255
      }
    }
    context.putImageData(output, 0, 0)
    const blob = await new Promise<Blob>((resolve, reject) => canvas.toBlob(value => value ? resolve(value) : reject(new Error('La carte des contours est indisponible.')), 'image/png'))
    return URL.createObjectURL(blob)
  } finally {
    bitmap.close()
  }
}
