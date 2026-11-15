import { MusicService } from '../bindings/changeme'
import type { MusicPlan, SemanticAnalysis } from '../bindings/changeme/models'
import { prepareAnalysisImage } from './structuralAnalysis'

// composePreview takes an imported image URL and validated model description and returns Go's repeatable score.
export async function composePreview(url: string, semantic: SemanticAnalysis): Promise<MusicPlan> {
  return MusicService.Compose(await prepareAnalysisImage(url), semantic)
}
