package render

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"le-grimoire/internal/constants"
	"net/url"

	"github.com/chromedp/chromedp"
)

func RenderFramebuffer(ctx context.Context, html string) ([]byte, error) {
	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var pngBuf []byte
	dataURL := "data:text/html," + url.PathEscape(html)

	err := chromedp.Run(allocCtx,
		chromedp.EmulateViewport(int64(constants.DisplayWidth), int64(constants.DisplayHeight)),
		chromedp.Navigate(dataURL),
		chromedp.CaptureScreenshot(&pngBuf),
	)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(pngBuf))
	if err != nil {
		return nil, err
	}

	gray := image.NewPaletted(img.Bounds(), color.Palette{
		color.Gray{Y: 0}, color.Gray{Y: 85}, color.Gray{Y: 170}, color.Gray{Y: 255},
	})
	draw.FloydSteinberg.Draw(gray, img.Bounds(), img, image.Point{})

	return packTo2bpp(gray), nil
}

func packTo2bpp(img *image.Paletted) []byte {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	out := make([]byte, (w*h*2+7)/8)
	bitPos := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			level := img.ColorIndexAt(x, y) // 0-3, 2 bits
			byteIdx := bitPos / 8
			shift := 6 - (bitPos % 8) // MSB-first
			out[byteIdx] |= level << shift
			bitPos += 2
		}
	}
	return out
}
