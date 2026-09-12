package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FlowingSPDG/streamdeck/v2"
	sdcontext "github.com/FlowingSPDG/streamdeck/v2/context"
)

type Settings struct {
	Counter       int    `json:"counter"`
	ButtonText    string `json:"buttonText"`
	Color         string `json:"color"`
	AutoIncrement bool   `json:"autoIncrement"`
}

type PropertyInspectorMessage struct {
	Action string `json:"action"`
}

type AllStatesResponse struct {
	Action string                 `json:"action"`
	States map[string]ButtonState `json:"states"`
}

type ResetCompleteResponse struct {
	Action string `json:"action"`
}

type ButtonState struct {
	Settings   Settings
	LastUpdate time.Time
}

type SettingsManager struct {
	mu           sync.Mutex
	buttonStates map[string]ButtonState
	tickers      map[string]*time.Ticker
}

func NewSettingsManager() *SettingsManager {
	return &SettingsManager{
		buttonStates: make(map[string]ButtonState),
		tickers:      make(map[string]*time.Ticker),
	}
}

func (sm *SettingsManager) StoreButtonState(contextID string, state ButtonState) {
	sm.mu.Lock()
	sm.buttonStates[contextID] = state
	sm.mu.Unlock()
}

func (sm *SettingsManager) LoadButtonState(contextID string) (ButtonState, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	state, ok := sm.buttonStates[contextID]
	return state, ok
}

func (sm *SettingsManager) DeleteButtonState(contextID string) {
	sm.mu.Lock()
	delete(sm.buttonStates, contextID)
	sm.mu.Unlock()
	sm.stopAutoIncrement(contextID)
}

func (sm *SettingsManager) UpdateButtonState(contextID string, settings Settings) {
	state, exists := sm.LoadButtonState(contextID)
	if !exists {
		state = ButtonState{Settings: settings, LastUpdate: time.Now()}
	} else {
		state.Settings = settings
		state.LastUpdate = time.Now()
	}
	sm.StoreButtonState(contextID, state)
	if settings.AutoIncrement {
		sm.startAutoIncrement(contextID)
	} else {
		sm.stopAutoIncrement(contextID)
	}
}

func (sm *SettingsManager) startAutoIncrement(contextID string) {
	sm.stopAutoIncrement(contextID)
	ticker := time.NewTicker(2 * time.Second)
	sm.mu.Lock()
	sm.tickers[contextID] = ticker
	sm.mu.Unlock()

	go func() {
		for range ticker.C {
			state, exists := sm.LoadButtonState(contextID)
			if !exists {
				ticker.Stop()
				sm.mu.Lock()
				delete(sm.tickers, contextID)
				sm.mu.Unlock()
				return
			}
			state.Settings.Counter++
			state.LastUpdate = time.Now()
			sm.StoreButtonState(contextID, state)
		}
	}()
}

func (sm *SettingsManager) stopAutoIncrement(contextID string) {
	sm.mu.Lock()
	ticker, ok := sm.tickers[contextID]
	if ok {
		delete(sm.tickers, contextID)
	}
	sm.mu.Unlock()
	if ok {
		ticker.Stop()
	}
}

func (sm *SettingsManager) GetAllButtonStates() map[string]ButtonState {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	result := make(map[string]ButtonState, len(sm.buttonStates))
	for k, v := range sm.buttonStates {
		result[k] = v
	}
	return result
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
	sm := NewSettingsManager()
	action := streamdeck.NewAction[Settings](client, "dev.samwho.streamdeck.settings_manager")

	action.OnWillAppear(func(ctx context.Context, e streamdeck.WillAppearEvent[Settings]) error {
		if e.Payload.Settings.ButtonText == "" {
			e.Payload.Settings.ButtonText = "Click Me"
		}
		if e.Payload.Settings.Color == "" {
			e.Payload.Settings.Color = "blue"
		}
		sm.UpdateButtonState(sdcontext.Context(ctx), e.Payload.Settings)
		return applyVisuals(ctx, action, e.Payload.Settings)
	})

	action.OnWillDisappear(func(ctx context.Context, e streamdeck.WillDisappearEvent[Settings]) error {
		sm.DeleteButtonState(sdcontext.Context(ctx))
		e.Payload.Settings.Counter = 0
		return action.SetSettings(ctx, e.Payload.Settings)
	})

	action.OnKeyDown(func(ctx context.Context, e streamdeck.KeyDownEvent[Settings]) error {
		e.Payload.Settings.Counter++
		sm.UpdateButtonState(sdcontext.Context(ctx), e.Payload.Settings)
		if err := action.SetSettings(ctx, e.Payload.Settings); err != nil {
			return err
		}
		return action.SetTitle(ctx, fmt.Sprintf("%s\n%d", e.Payload.Settings.ButtonText, e.Payload.Settings.Counter), streamdeck.HardwareAndSoftware)
	})

	action.OnDidReceiveSettings(func(ctx context.Context, e streamdeck.DidReceiveSettingsEvent[Settings]) error {
		sm.UpdateButtonState(sdcontext.Context(ctx), e.Payload.Settings)
		return applyVisuals(ctx, action, e.Payload.Settings)
	})

	action.OnPropertyInspectorMessage(func(ctx context.Context, payload PropertyInspectorMessage) error {
		switch payload.Action {
		case "getAllStates":
			return action.SendToPropertyInspector(ctx, AllStatesResponse{
				Action: "allStates",
				States: sm.GetAllButtonStates(),
			})
		case "resetAll":
			sm.mu.Lock()
			for key, value := range sm.buttonStates {
				value.Settings.Counter = 0
				sm.buttonStates[key] = value
			}
			sm.mu.Unlock()
			return action.SendToPropertyInspector(ctx, ResetCompleteResponse{Action: "resetComplete"})
		}
		return nil
	})

	return client.Run(ctx)
}

func applyVisuals(ctx context.Context, action *streamdeck.Action[Settings], settings Settings) error {
	bg, err := streamdeck.Image(createBackground(settings.Color))
	if err != nil {
		return err
	}
	if err := action.SetImage(ctx, bg, streamdeck.HardwareAndSoftware); err != nil {
		return err
	}
	title := fmt.Sprintf("%s\n%d", settings.ButtonText, settings.Counter)
	return action.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
}

func createBackground(colorName string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))
	var bgColor color.Color
	switch colorName {
	case "red":
		bgColor = color.RGBA{R: 255, A: 255}
	case "green":
		bgColor = color.RGBA{G: 255, A: 255}
	case "blue":
		bgColor = color.RGBA{B: 255, A: 255}
	case "yellow":
		bgColor = color.RGBA{R: 255, G: 255, A: 255}
	case "purple":
		bgColor = color.RGBA{R: 128, B: 128, A: 255}
	default:
		bgColor = color.RGBA{R: 64, G: 64, B: 64, A: 255}
	}
	for x := 0; x < 72; x++ {
		for y := 0; y < 72; y++ {
			img.Set(x, y, bgColor)
		}
	}
	return img
}
