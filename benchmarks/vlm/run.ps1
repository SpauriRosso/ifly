param([string]$Label = 'smolvlm-256m')

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$results = @()
$prompt = 'Describe this simple drawing in one short sentence. Mention only visible objects and actions. If the scene or emotion is unclear, say uncertain. Do not invent details.'

foreach ($name in @('landscape', 'abstract', 'person', 'happy_house', 'threat')) {
    $imagePath = Join-Path $root "$name.png"
    $imageURL = 'data:image/png;base64,' + [Convert]::ToBase64String([IO.File]::ReadAllBytes($imagePath))
    $request = @{
        model = 'local-vlm'
        temperature = 0
        max_tokens = 128
        stream = $false
        messages = @(@{
            role = 'user'
            content = @(
                @{ type = 'text'; text = $prompt },
                @{ type = 'image_url'; image_url = @{ url = $imageURL } }
            )
        })
    } | ConvertTo-Json -Depth 10 -Compress
    $timer = [Diagnostics.Stopwatch]::StartNew()
    try {
        $reply = Invoke-RestMethod -Uri 'http://127.0.0.1:11480/v1/chat/completions' -Method Post -Body $request -ContentType 'application/json' -TimeoutSec 120
        $timer.Stop()
        $results += [pscustomobject]@{
            image = $name
            seconds = [math]::Round($timer.Elapsed.TotalSeconds, 2)
            text = $reply.choices[0].message.content
            tokens = $reply.usage.completion_tokens
            error = $null
        }
    } catch {
        $timer.Stop()
        $results += [pscustomobject]@{
            image = $name
            seconds = [math]::Round($timer.Elapsed.TotalSeconds, 2)
            text = $null
            tokens = $null
            error = $_.Exception.Message
        }
    }
}

$results | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath (Join-Path $root "$Label-results.json") -Encoding utf8
$results | Format-Table -AutoSize
