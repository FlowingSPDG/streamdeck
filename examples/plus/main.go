package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/FlowingSPDG/streamdeck/v2"
)

type Settings struct {
	Value int `json:"value"`
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
	action := streamdeck.NewAction[Settings](client, "dev.samwho.streamdeck.plus")

	action.OnWillAppear(func(ctx context.Context, e streamdeck.WillAppearEvent[Settings]) error {
		return update(ctx, action, e.Payload.Settings)
	})

	action.OnDialRotate(func(ctx context.Context, e streamdeck.DialRotateEvent[Settings]) error {
		e.Payload.Settings.Value += e.Payload.Ticks
		if err := action.SetSettings(ctx, e.Payload.Settings); err != nil {
			return err
		}
		return update(ctx, action, e.Payload.Settings)
	})

	action.OnTouchTap(func(ctx context.Context, e streamdeck.TouchTapEvent[Settings]) error {
		e.Payload.Settings.Value = 0
		if err := action.SetSettings(ctx, e.Payload.Settings); err != nil {
			return err
		}
		return update(ctx, action, e.Payload.Settings)
	})

	action.OnDialDown(func(ctx context.Context, e streamdeck.DialDownEvent[Settings]) error {
		return action.ShowOk(ctx)
	})

	return client.Run(ctx)
}

func update(ctx context.Context, action *streamdeck.Action[Settings], settings Settings) error {
	if err := action.SetFeedbackLayout(ctx, "$B1"); err != nil {
		return err
	}
	if err := action.SetFeedback(ctx, streamdeck.Feedback{
		"title": streamdeck.TextFeedback{Value: "Level"},
		"value": streamdeck.TextFeedback{Value: fmt.Sprintf("%d", settings.Value)},
	}); err != nil {
		return err
	}
	return action.SetTitle(ctx, fmt.Sprintf("%d", settings.Value), streamdeck.HardwareAndSoftware)
}
