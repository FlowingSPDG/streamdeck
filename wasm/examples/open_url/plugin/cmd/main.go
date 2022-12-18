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
	c.LogMessage("Start Backend")
	done := make(chan struct{})
	<-done
}
