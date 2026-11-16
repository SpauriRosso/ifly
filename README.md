# ifly
i'm bored so i'm doing stupid shit, don't ask why ! I'm merely following Ovide's teaching

## Discord Rich Presence

IFLY uses Discord application `707372989921886250` by default. Launch Discord Desktop, then launch IFLY. To use another Discord application ID, set an override before starting IFLY in PowerShell:

```powershell
$env:IFLY_DISCORD_APP_ID = "ANOTHER_APPLICATION_ID"
wails3 dev
```

IFLY shows **Playing IFLY** with **Drawing a piece of art** before an image is open and **Admiring a piece of art** while an image is visible. Its bundled artwork is [build/appicon.png](build/appicon.png). For Discord to display that image on other users' profiles, upload it as the application icon or a Rich Presence asset in the [Discord Developer Portal](https://discord.com/developers/applications/707372989921886250); Discord cannot fetch image bytes from a local executable.
