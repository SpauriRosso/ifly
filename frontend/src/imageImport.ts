export type ImageFormat = 'PNG' | 'JPEG' | 'WebP'

export interface ImageData {
  fileName: string
  format: ImageFormat
  width: number
  height: number
  byteSize: number
  url: string
}

const MAX_BYTES = 20 * 1024 * 1024
const MAX_SIDE = 8192
const MAX_PIXELS = 24_000_000

// readJpegSize scans JPEG markers and returns the encoded dimensions, or null for malformed data.
function readJpegSize(bytes: Uint8Array): [number, number] | null {
  let offset = 2
  while (offset + 4 <= bytes.length) {
    if (bytes[offset++] !== 0xff) return null
    let marker = bytes[offset++]
    while (marker === 0xff && offset < bytes.length) marker = bytes[offset++]
    if (marker === 0xd9 || marker === 0xda) return null
    if (marker === 0x01 || (marker >= 0xd0 && marker <= 0xd7)) continue
    if (offset + 2 > bytes.length) return null
    const length = (bytes[offset] << 8) | bytes[offset + 1]
    if (length < 2 || offset + length > bytes.length) return null
    if ([0xc0, 0xc1, 0xc2, 0xc3, 0xc5, 0xc6, 0xc7, 0xc9, 0xca, 0xcb, 0xcd, 0xce, 0xcf].includes(marker)) {
      if (length < 7) return null
      return [(bytes[offset + 5] << 8) | bytes[offset + 6], (bytes[offset + 3] << 8) | bytes[offset + 4]]
    }
    offset += length
  }
  return null
}

// inspectImage reads a bounded file header and returns its actual format and encoded dimensions.
function inspectImage(bytes: Uint8Array): { format: ImageFormat; size: [number, number] } | null {
  const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength)
  const ascii = (at: number, length: number) => String.fromCharCode(...bytes.subarray(at, at + length))
  if (bytes.length >= 24 && ascii(0, 8) === '\x89PNG\r\n\x1a\n' && ascii(12, 4) === 'IHDR') {
    return { format: 'PNG', size: [view.getUint32(16), view.getUint32(20)] }
  }
  if (bytes.length >= 4 && bytes[0] === 0xff && bytes[1] === 0xd8 && bytes[2] === 0xff) {
    const size = readJpegSize(bytes)
    return size ? { format: 'JPEG', size } : null
  }
  if (bytes.length >= 30 && ascii(0, 4) === 'RIFF' && ascii(8, 4) === 'WEBP') {
    const chunk = ascii(12, 4)
    if (chunk === 'VP8X') return { format: 'WebP', size: [1 + bytes[24] + (bytes[25] << 8) + (bytes[26] << 16), 1 + bytes[27] + (bytes[28] << 8) + (bytes[29] << 16)] }
    if (chunk === 'VP8 ' && bytes[23] === 0x9d && bytes[24] === 0x01 && bytes[25] === 0x2a) {
      return { format: 'WebP', size: [view.getUint16(26, true) & 0x3fff, view.getUint16(28, true) & 0x3fff] }
    }
    if (chunk === 'VP8L' && bytes[20] === 0x2f) {
      return { format: 'WebP', size: [1 + (((bytes[22] & 0x3f) << 8) | bytes[21]), 1 + (((bytes[24] & 0x0f) << 10) | (bytes[23] << 2) | (bytes[22] >> 6))] }
    }
  }
  return null
}

// loadImage validates a local file, decodes it with browser orientation handling, and returns a disposable preview URL.
export async function loadImage(file: File): Promise<ImageData> {
  if (file.size === 0 || file.size > MAX_BYTES) throw new Error('Image vide ou supérieure à 20 Mo.')
  const extension = file.name.split('.').pop()?.toLowerCase()
  const allowed: Record<ImageFormat, string[]> = { PNG: ['png'], JPEG: ['jpg', 'jpeg'], WebP: ['webp'] }
  const bytes = new Uint8Array(await file.slice(0, Math.min(file.size, 256 * 1024)).arrayBuffer())
  const details = inspectImage(bytes)
  if (!details || !extension || !allowed[details.format].includes(extension)) {
    throw new Error('Format non reconnu ou extension incompatible. Utilisez PNG, JPEG ou WebP.')
  }
  const [width, height] = details.size
  if (!width || !height || width > MAX_SIDE || height > MAX_SIDE || width * height > MAX_PIXELS) {
    throw new Error('Dimensions trop grandes : maximum 8 192 pixels par côté et 24 mégapixels.')
  }
  const bitmap = await createImageBitmap(file, { imageOrientation: 'from-image' })
  try {
    const normal = bitmap.width === width && bitmap.height === height
    const rotated = bitmap.width === height && bitmap.height === width
    if (bitmap.width * bitmap.height > MAX_PIXELS || (!normal && !rotated)) {
      throw new Error('Les dimensions décodées ne correspondent pas à une image prise en charge.')
    }
  } finally {
    bitmap.close()
  }
  return { fileName: file.name, format: details.format, width: bitmap.width, height: bitmap.height, byteSize: file.size, url: URL.createObjectURL(file) }
}
