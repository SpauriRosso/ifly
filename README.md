# ifly
i'm bored so i'm doing stupid shit, don't ask why ! I'm merely following Ovide's teaching

## Discord Rich Presence

IFLY uses Discord application `707372989921886250` by default. Launch Discord Desktop, then launch IFLY. To use another Discord application ID, set an override before starting IFLY in PowerShell:

```powershell
$env:IFLY_DISCORD_APP_ID = "ANOTHER_APPLICATION_ID"
wails3 dev
```

IFLY shows **Playing IFLY** with **Drawing a piece of art** before an image is open and **Admiring a piece of art** while an image is visible. Its bundled artwork is [build/appicon.png](build/appicon.png). To replace Discord's question-mark icon, upload this PNG in the [Discord Developer Portal](https://discord.com/developers/applications/707372989921886250): set it as the application icon under **General Information**, then add it under **Rich Presence > Art Assets** with the key `ifly`. IFLY sends that asset key in every activity update. Discord cannot read the PNG from a local executable, so bundling it alone cannot make it appear on the profile card.
