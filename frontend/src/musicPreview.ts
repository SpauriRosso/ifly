import { MusicService } from '../bindings/ifly/internal/music'
import type { MusicPlan } from '../bindings/ifly/internal/music/models'
import type { SemanticAnalysis } from '../bindings/ifly/internal/semantic/models'
import { prepareAnalysisImage } from './structuralAnalysis'

// composePreview takes an imported image URL and validated model description and returns Go's repeatable score.
export async function composePreview(url: string, semantic: SemanticAnalysis): Promise<MusicPlan> {
  return MusicService.Compose(await prepareAnalysisImage(url), semantic)
}

// composeVariant takes an image, scene description, visible seed, and theme, then asks Go for a new deterministic score.
export async function composeVariant(url: string, semantic: SemanticAnalysis, seed: string, theme: string): Promise<MusicPlan> {
  return MusicService.ComposeVariant(await prepareAnalysisImage(url), semantic, seed, theme)
}

// createMusicSeed takes no input and returns eight cryptographically random bytes as 16 hexadecimal digits.
export function createMusicSeed(): string {
  return Array.from(crypto.getRandomValues(new Uint8Array(8)), byte => byte.toString(16).padStart(2, '0')).join('')
}
