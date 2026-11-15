package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"ifly/internal/analysis"
	"ifly/internal/music"
	"ifly/internal/semantic"
)

// assets contains the built React application served by the Wails webview.
//
//go:embed all:frontend/dist
var assets embed.FS

// main starts the IFLY desktop window and binds the analysis, semantic, and music services.
func main() {
	app := application.New(application.Options{
		Name:        "IFLY",
		Description: "A local audiovisual artwork laboratory",
		Services: []application.Service{
			application.NewService(&analysis.AnalysisService{}),
			application.NewService(semantic.NewSemanticService()),
			application.NewService(&music.MusicService{}),
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
