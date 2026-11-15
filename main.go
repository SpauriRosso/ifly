package main

import (
	"embed"
	"log"

	"changeme/internal/analysis"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// assets contains the built React application served by the Wails webview.
//
//go:embed all:frontend/dist
var assets embed.FS

// main starts the IFLY desktop window and binds structural and semantic analysis services.
func main() {
	app := application.New(application.Options{
		Name:        "IFLY",
		Description: "A local audiovisual artwork laboratory",
		Services: []application.Service{
			application.NewService(&analysis.AnalysisService{}),
			application.NewService(NewSemanticService()),
			application.NewService(&MusicService{}),
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
