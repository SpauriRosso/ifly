import { SemanticService } from '../bindings/changeme'
import type { SemanticAnalysis } from '../bindings/changeme/models'
import { prepareAnalysisImage } from './structuralAnalysis'

// interpretPreview sends a bounded local PNG to Go and returns the validated, unverified VLM description.
export async function interpretPreview(url: string): Promise<SemanticAnalysis> {
  return SemanticService.Analyze(await prepareAnalysisImage(url))
}
