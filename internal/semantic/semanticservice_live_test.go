package semantic

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"testing"
	"time"
)

// TestSemanticServiceLive checks the full Go service against a local VLM when a real sketch path is supplied.
func TestSemanticServiceLive(t *testing.T) {
	path := os.Getenv("IFLY_VLM_TEST_IMAGE")
	if path == "" {
		t.Skip("set IFLY_VLM_TEST_IMAGE to run the local model integration test")
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	source, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	width, height := source.Bounds().Dx(), source.Bounds().Dy()
	if width < 1 || height < 1 {
		t.Fatal("empty image")
	}
	thumb := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			thumb.Set(x, y, source.At(source.Bounds().Min.X+x*width/256, source.Bounds().Min.Y+y*height/256))
		}
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, thumb); err != nil {
		t.Fatal(err)
	}
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(encoded.Bytes())
	started := time.Now()
	result, err := NewSemanticService().Analyze(dataURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("elapsed=%s description=%q subjects=%v uncertain=%t", time.Since(started), result.Description, result.Subjects, result.Uncertain)
}
