package main

import (
	"context"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/FlowingSPDG/streamdeck/v2"
)

type Settings struct {
	Counter int `json:"counter"`
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	action := streamdeck.NewAction[Settings](client, "dev.samwho.streamdeck.counter")

	action.OnWillAppear(func(ctx context.Context, e streamdeck.WillAppearEvent[Settings]) error {
		bg, err := streamdeck.Image(background())
		if err != nil {
			return err
		}
		if err := action.SetImage(ctx, bg, streamdeck.HardwareAndSoftware); err != nil {
			return err
		}
		return action.SetTitle(ctx, strconv.Itoa(e.Payload.Settings.Counter), streamdeck.HardwareAndSoftware)
	})

	action.OnWillDisappear(func(ctx context.Context, e streamdeck.WillDisappearEvent[Settings]) error {
		e.Payload.Settings.Counter = 0
		return action.SetSettings(ctx, e.Payload.Settings)
	})

	action.OnKeyDown(func(ctx context.Context, e streamdeck.KeyDownEvent[Settings]) error {
		e.Payload.Settings.Counter++
		if err := action.SetSettings(ctx, e.Payload.Settings); err != nil {
			return err
		}
		return action.SetTitle(ctx, strconv.Itoa(e.Payload.Settings.Counter), streamdeck.HardwareAndSoftware)
	})

	return client.Run(ctx)
}

func background() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))
	for x := 0; x < 72; x++ {
		for y := 0; y < 72; y++ {
			img.Set(x, y, color.Black)
		}
	}
	return img
}
