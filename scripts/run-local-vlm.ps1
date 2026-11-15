$ErrorActionPreference = 'Stop'
$root = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\local-vlm'))
$server = Join-Path $root 'runtime\llama-server.exe'
$model = Join-Path $root 'SmolVLM2-2.2B-Instruct-Q4_K_M.gguf'
$projector = Join-Path $root 'mmproj-SmolVLM2-2.2B-Instruct-Q8_0.gguf'

if (-not (Test-Path -LiteralPath $server) -or -not (Test-Path -LiteralPath $model) -or -not (Test-Path -LiteralPath $projector)) {
    throw 'Local VLM files are missing. Run scripts/setup-local-vlm.ps1 first.'
}

& $server --model $model --mmproj $projector --no-mmproj-offload --n-gpu-layers 0 --host 127.0.0.1 --port 11480 --ctx-size 2048 --cors-origins localhost
