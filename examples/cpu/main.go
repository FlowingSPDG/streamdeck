package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/FlowingSPDG/streamdeck/v2"
	sdcontext "github.com/FlowingSPDG/streamdeck/v2/context"
	"github.com/shirou/gopsutil/v4/cpu"
)

const (
	imgX = 72
	imgY = 72
)

type Settings struct{}

type PropertyInspectorSettings struct {
	ShowText bool `json:"showText,omitempty"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	action := streamdeck.NewAction[Settings](client, "dev.samwho.streamdeck.cpu")

	var (
		piMu    sync.RWMutex
		pi      PropertyInspectorSettings
		ctxMu   sync.RWMutex
		visible = map[string]struct{}{}
	)

	action.OnPropertyInspectorMessage(func(ctx context.Context, msg PropertyInspectorSettings) error {
		piMu.Lock()
		pi = msg
		piMu.Unlock()
		return nil
	})

	action.OnWillAppear(func(ctx context.Context, e streamdeck.WillAppearEvent[Settings]) error {
		ctxMu.Lock()
		visible[e.Context] = struct{}{}
		ctxMu.Unlock()
		return nil
	})

	action.OnWillDisappear(func(ctx context.Context, e streamdeck.WillDisappearEvent[Settings]) error {
		ctxMu.Lock()
		delete(visible, e.Context)
		ctxMu.Unlock()
		return nil
	})

	readings := make([]float64, imgX)

	go func() {
		ticker := time.NewTicker(time.Second / 4)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				copy(readings, readings[1:])
				r, err := cpu.Percent(0, false)
				if err != nil || len(r) == 0 {
					log.Printf("cpu.Percent: %v", err)
					continue
				}
				readings[imgX-1] = r[0]

				img, err := streamdeck.Image(graph(readings))
				if err != nil {
					log.Printf("image: %v", err)
					continue
				}

				piMu.RLock()
				showText := pi.ShowText
				piMu.RUnlock()

				ctxMu.RLock()
				ids := make([]string, 0, len(visible))
				for id := range visible {
					ids = append(ids, id)
				}
				ctxMu.RUnlock()

				for _, id := range ids {
					sdctx := sdcontext.WithContext(ctx, id)
					if err := action.SetImage(sdctx, img, streamdeck.HardwareAndSoftware); err != nil {
						log.Printf("set image: %v", err)
						continue
					}
					title := ""
					if showText {
						title = fmt.Sprintf("CPU\n%d%%", int(r[0]))
					}
					if err := action.SetTitle(sdctx, title, streamdeck.HardwareAndSoftware); err != nil {
						log.Printf("set title: %v", err)
					}
				}
			}
		}
	}()

	return client.Run(ctx)
}

func graph(readings []float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	for x := 0; x < imgX; x++ {
		reading := readings[x] / 100
		upto := int(float64(imgY) * reading)
		if upto < 0 {
			upto = 0
		}
		if upto > imgY {
			upto = imgY
		}
		for y := 0; y < upto; y++ {
			img.Set(x, imgY-1-y, color.RGBA{R: 255, A: 255})
		}
		for y := upto; y < imgY; y++ {
			img.Set(x, imgY-1-y, color.Black)
		}
	}
	return img
}
