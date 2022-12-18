package main

import (
	"context"
	"os"

	"github.com/FlowingSPDG/streamdeck"
)

func main() {
	ctx := context.Background()
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		panic(err)
	}
	c := streamdeck.NewClient(ctx, params)
	ac := c.Action("dev.flowingspdg.was.openurl")
	ac.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		client.LogMessage("WillAppear on Backend")
		return nil
	})
	if err := c.Run(); err != nil {
		panic(err)
	}
}
