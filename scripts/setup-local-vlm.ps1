param([string]$SourceCache = '')

$ErrorActionPreference = 'Stop'
$root = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\local-vlm'))
New-Item -ItemType Directory -Force -Path $root | Out-Null

$assets = @(
    @{
        Name = 'llama-cpu.zip'
        URL = 'https://github.com/ggml-org/llama.cpp/releases/download/b12415/llama-b12415-bin-win-cpu-x64.zip'
        SHA256 = '83d65aca06314b8fa38378a49c4fb16e08d5547460f7bc916b737204ba6a39c1'
    },
    @{
        Name = 'SmolVLM2-2.2B-Instruct-Q4_K_M.gguf'
        URL = 'https://huggingface.co/ggml-org/SmolVLM2-2.2B-Instruct-GGUF/resolve/main/SmolVLM2-2.2B-Instruct-Q4_K_M.gguf?download=true'
        SHA256 = '0cf76814555b8665149075b74ab6b5c1d428ea1d3d01c1918c12012e8d7c9f58'
    },
    @{
        Name = 'mmproj-SmolVLM2-2.2B-Instruct-Q8_0.gguf'
        URL = 'https://huggingface.co/ggml-org/SmolVLM2-2.2B-Instruct-GGUF/resolve/main/mmproj-SmolVLM2-2.2B-Instruct-Q8_0.gguf?download=true'
        SHA256 = 'ae07ea1facd07dd3230c4483b63e8cda96c6944ad2481f33d531f79e892dd024'
    }
)

# InstallAsset copies a verified local download or fetches the pinned official file and verifies its SHA-256.
function InstallAsset($asset) {
    $destination = Join-Path $root $asset.Name
    if (-not (Test-Path -LiteralPath $destination)) {
        $cached = if ($SourceCache) { Join-Path $SourceCache $asset.Name } else { '' }
        if ($cached -and (Test-Path -LiteralPath $cached)) {
            Copy-Item -LiteralPath $cached -Destination $destination
        } else {
            & curl.exe --fail --location --retry 2 --output $destination $asset.URL
            if ($LASTEXITCODE -ne 0) { throw "Download failed: $($asset.Name)" }
        }
    }
    $actual = (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $asset.SHA256) { throw "SHA-256 mismatch: $($asset.Name)" }
    Write-Host "Verified $($asset.Name)"
}

foreach ($asset in $assets) { InstallAsset $asset }
$runtime = Join-Path $root 'runtime'
if (-not (Test-Path -LiteralPath (Join-Path $runtime 'llama-server.exe'))) {
    Expand-Archive -LiteralPath (Join-Path $root 'llama-cpu.zip') -DestinationPath $runtime -Force
}
Write-Host "Local VLM ready in $root"
