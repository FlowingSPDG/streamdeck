package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/FlowingSPDG/streamdeck"
	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/shirou/gopsutil/cpu"
)

const (
	imgX = 72
	imgY = 72
)

type PropertyInspectorSettings struct {
	ShowText bool `json:"showText,omitempty"`
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	f, err := os.Create(path.Join(exeDir, "streamdeck-cpu.log"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// log.SetOutput(f)
	log.SetOutput(os.Stdout)

	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run(ctx context.Context) error {
	fmt.Println("args:", strings.Join(os.Args, " "))
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}
	fmt.Println("params:", params)

	client := streamdeck.NewClient(ctx, params)
	setup(client)

	return client.Run()
}

func setup(client *streamdeck.Client) {
	action := client.Action("dev.samwho.streamdeck.cpu")

	pi := &PropertyInspectorSettings{}
	contexts := sync.Map{}

	action.RegisterHandler(streamdeck.SendToPlugin, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		b, _ := json.MarshalIndent(event, "", "	")
		fmt.Printf("event:%s\n", b)
		return json.Unmarshal(event.Payload, pi)
	})

	action.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		b, _ := json.MarshalIndent(event, "", "	")
		fmt.Printf("event:%s\n", b)
		return json.Unmarshal(event.Payload, pi)
	})

	action.RegisterHandler(streamdeck.KeyUp, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		b, _ := json.MarshalIndent(event, "", "	")
		fmt.Printf("event:%s\n", b)
		return json.Unmarshal(event.Payload, pi)
	})

	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		b, _ := json.MarshalIndent(event, "", "	")
		fmt.Printf("event:%s\n", b)
		contexts.Store(event.Context, struct{}{})
		return nil
	})

	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		b, _ := json.MarshalIndent(event, "", "	")
		fmt.Printf("event:%s\n", b)
		contexts.Delete(event.Context)
		return nil
	})

	readings := make([]float64, imgX, imgX)

	go func() {
		for range time.Tick(time.Second / 4) {
			for i := 0; i < imgX-1; i++ {
				readings[i] = readings[i+1]
			}

			r, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Printf("error getting CPU reading: %v\n", err)
			}
			readings[imgX-1] = r[0]

			contexts.Range(func(key, value any) bool {
				ctxStr := value.(string)
				ctx := context.Background()
				ctx = sdcontext.WithContext(ctx, ctxStr)

				img, err := streamdeck.Image(graph(readings))
				if err != nil {
					fmt.Printf("error creating image: %v\n", err)
					return true
				}

				if err := client.SetImage(ctx, img, streamdeck.HardwareAndSoftware); err != nil {
					fmt.Printf("error setting image: %v\n", err)
					return true
				}

				title := ""
				if pi.ShowText {
					title = fmt.Sprintf("CPU\n%d%%", int(r[0]))
				}

				if err := client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware); err != nil {
					fmt.Printf("error setting title: %v\n", err)
					return true
				}
				return true
			})
		}
	}()
}

func graph(readings []float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, imgX, imgY))
	for x := 0; x < imgX; x++ {
		reading := readings[x] / 100
		upto := int(float64(imgY) * reading)
		for y := 0; y < upto; y++ {
			img.Set(x, imgY-y, color.RGBA{R: 255, A: 255})
		}
		for y := upto; y < imgY; y++ {
			img.Set(x, imgY-y, color.Black)
		}
	}
	return img
}
