package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/FlowingSPDG/streamdeck"

	"github.com/FlowingSPDG/streamdeck/wasm/examples/open_url/models"
)

func main() {
	ctx := context.Background()
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		panic(err)
	}

	c := streamdeck.NewClient(ctx, params)
	ac := c.Action("dev.flowingspdg.wasm.openurl")
	ac.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		client.LogMessage("WillAppear on Backend")
		payload := streamdeck.WillAppearPayload[models.Settings]{}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return err
		}

		if payload.Settings.IsDefault() {
			payload.Settings.Initialize()
			client.SetSettings(ctx, payload.Settings)
		}

		go func() {
			for {
				// 無限に回り続けるので、何度もWillAppearが開かれるとその分goroutineが生成されてしまう
				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{
					URL: "https://www.elgato.com/",
				})

				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{
					URL: "https://go.dev/",
				})

				time.Sleep(time.Second)
				c.SendToPropertyInspector(ctx, &models.Settings{
					URL: "https://github.com/FlowingSPDG/streamdeck",
				})
			}
		}()
		return nil
	})
	if err := c.Run(); err != nil {
		panic(err)
	}
}
