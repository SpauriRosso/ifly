package main

import (
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"ifly/internal/analysis"
	"ifly/internal/music"
	"ifly/internal/presence"
	"ifly/internal/semantic"
)

// assets contains the built React application served by the Wails webview.
//
//go:embed all:frontend/dist
var assets embed.FS

// main starts the IFLY desktop window, binds its services, and closes Discord presence on exit.
func main() {
	presenceService := presence.NewService(os.Getenv("IFLY_DISCORD_APP_ID"))
	defer presenceService.Close()
	app := application.New(application.Options{
		Name:        "IFLY",
		Description: "A local audiovisual artwork laboratory",
		Services: []application.Service{
			application.NewService(&analysis.AnalysisService{}),
			application.NewService(semantic.NewSemanticService()),
			application.NewService(&music.MusicService{}),
			application.NewService(presenceService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "IFLY — Laboratoire",
		Width:            1200,
		Height:           760,
		BackgroundColour: application.NewRGB(233, 231, 226),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
