import { MusicService } from '../bindings/ifly/internal/music'
import type { MusicPlan } from '../bindings/ifly/internal/music/models'
import type { SemanticAnalysis } from '../bindings/ifly/internal/semantic/models'
import { prepareAnalysisImage } from './structuralAnalysis'

// composePreview takes an imported image URL and validated model description and returns Go's repeatable score.
export async function composePreview(url: string, semantic: SemanticAnalysis): Promise<MusicPlan> {
  return MusicService.Compose(await prepareAnalysisImage(url), semantic)
}
