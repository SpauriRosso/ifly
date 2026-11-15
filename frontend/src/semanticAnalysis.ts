import { SemanticService } from '../bindings/ifly/internal/semantic'
import type { SemanticAnalysis } from '../bindings/ifly/internal/semantic/models'
import { prepareAnalysisImage } from './structuralAnalysis'

// interpretPreview sends a bounded local PNG to Go and returns the validated, unverified VLM description.
export async function interpretPreview(url: string): Promise<SemanticAnalysis> {
  return SemanticService.Analyze(await prepareAnalysisImage(url))
}

// cancelInterpretation asks Go to stop an in-flight local model request when its image is discarded.
export async function cancelInterpretation(): Promise<void> {
  await SemanticService.Cancel()
}
