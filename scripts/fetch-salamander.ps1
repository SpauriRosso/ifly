# Downloads the compact Salamander Grand Piano selection bundled for offline playback.
$ErrorActionPreference = 'Stop'
$root = Join-Path $PSScriptRoot '..\frontend\public\samples\salamander'
$baseUrl = 'https://raw.githubusercontent.com/Tonejs/audio/master/salamander'
$notes = @('C2', 'Ds2', 'Fs2', 'A2', 'C3', 'Ds3', 'Fs3', 'A3', 'C4', 'Ds4', 'Fs4', 'A4', 'C5', 'Ds5', 'Fs5', 'A5', 'C6', 'Ds6', 'Fs6')

New-Item -ItemType Directory -Force -Path $root | Out-Null
foreach ($note in $notes) {
    $destination = Join-Path $root "$note.mp3"
    if (-not (Test-Path -LiteralPath $destination)) {
        & curl.exe --fail --location --silent --show-error --output $destination "$baseUrl/$note.mp3"
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to download Salamander sample $note"
        }
    }
    $length = (Get-Item -LiteralPath $destination).Length
    if ($length -lt 20000 -or $length -gt 1000000) {
        throw "Unexpected Salamander sample size for $destination`: $length bytes"
    }
}
