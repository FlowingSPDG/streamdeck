package main

import (
	"context"
	"image"
	"image/color"
	"os"
	"strconv"

	"github.com/FlowingSPDG/streamdeck"
)

type Settings struct {
	Counter int `json:"counter"`
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	setup(client)

	return client.Run(ctx)
}

func setup(client *streamdeck.Client) {
	action := client.Action("dev.samwho.streamdeck.counter")

	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		// Settings will be passed through the event.
		// So we don't need to store them internally.
		// But in case you want to store them, you should store the settings since the settings will appear at this event
		var p streamdeck.WillAppearPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return err
		}

		bg, err := streamdeck.Image(background())
		if err != nil {
			return err
		}

		if err := client.SetImage(ctx, bg, streamdeck.HardwareAndSoftware); err != nil {
			return err
		}

		return client.SetTitle(ctx, strconv.Itoa(p.Settings.Counter), streamdeck.HardwareAndSoftware)
	})

	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		// Settings will be passed through the event.
		// So we don't need to store them internally.
		// But in case you want to store them, you should remove the settings since the settings will disappear at this event
		var p streamdeck.WillDisappearPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return err
		}

		p.Settings.Counter = 0
		return client.SetSettings(ctx, p.Settings)
	})

	action.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.KeyDownPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return err
		}

		p.Settings.Counter++
		if err := client.SetSettings(ctx, p.Settings); err != nil {
			return err
		}

		return client.SetTitle(ctx, strconv.Itoa(p.Settings.Counter), streamdeck.HardwareAndSoftware)
	})
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
