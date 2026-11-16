# ifly
i'm bored so i'm doing stupid shit, don't ask why ! I'm merely following Ovide's teaching

## Discord Rich Presence

Create an application named **IFLY** in the [Discord Developer Portal](https://discord.com/developers/applications), copy its **Application ID**, and set it before starting IFLY in PowerShell:

```powershell
$env:IFLY_DISCORD_APP_ID = "YOUR_APPLICATION_ID"
wails3 dev
```

To keep the ID for future terminals, run `setx IFLY_DISCORD_APP_ID "YOUR_APPLICATION_ID"` once, then open a new terminal. Discord Desktop must be running. IFLY shows **Playing IFLY** with **Drawing a piece of art** before an image is open and **Admiring a piece of art** while an image is visible. Presence is optional; without the environment variable, IFLY starts normally.
