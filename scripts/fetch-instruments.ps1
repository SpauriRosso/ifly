# Downloads the selected FluidR3 General MIDI samples used by the offline instrument player.
$ErrorActionPreference = 'Stop'
$root = Join-Path $PSScriptRoot '..\frontend\public\samples\fluidr3'
$baseUrl = 'https://raw.githubusercontent.com/gleitz/midi-js-soundfonts/gh-pages/FluidR3_GM'
$instruments = [ordered]@{
    acoustic_grand_piano = @('C3', 'G3', 'D4', 'A4', 'E5', 'B5', 'G6')
    acoustic_guitar_nylon = @('C3', 'G3', 'D4', 'A4', 'E5', 'B5')
    acoustic_bass = @('C2', 'G2', 'D3', 'A3', 'E4')
    violin = @('C4', 'G4', 'D5', 'A5', 'E6')
    cello = @('C3', 'G3', 'D4', 'A4', 'E5')
    flute = @('C4', 'G4', 'D5', 'A5', 'E6')
    voice_oohs = @('C4', 'G4', 'D5', 'A5', 'E6')
    woodblock = @('C4')
    timpani = @('C3')
}

foreach ($instrument in $instruments.Keys) {
    $directory = Join-Path $root $instrument
    New-Item -ItemType Directory -Force -Path $directory | Out-Null
    foreach ($note in $instruments[$instrument]) {
        $destination = Join-Path $directory "$note.mp3"
        if (-not (Test-Path -LiteralPath $destination)) {
            $url = "${baseUrl}/${instrument}-mp3/${note}.mp3"
            Invoke-WebRequest -Uri $url -OutFile $destination
        }
        $length = (Get-Item -LiteralPath $destination).Length
        if ($length -lt 1024 -or $length -gt 500000) {
            throw "Unexpected sample size for $destination`: $length bytes"
        }
    }
}
