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
	Settings   Settings  `json:"settings"`
	LastUpdate time.Time `json:"lastUpdate"`
}

type SettingsManager struct {
	mu           sync.Mutex
	buttonStates map[string]ButtonState
	tickers      map[string]*time.Ticker
	onTick       func(contextID string, settings Settings)
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
		sm.ensureAutoIncrement(contextID)
	} else {
		sm.stopAutoIncrement(contextID)
	}
}

func (sm *SettingsManager) ensureAutoIncrement(contextID string) {
	sm.mu.Lock()
	_, running := sm.tickers[contextID]
	sm.mu.Unlock()
	if running {
		return
	}
	sm.startAutoIncrement(contextID)
}

func (sm *SettingsManager) startAutoIncrement(contextID string) {
	ticker := time.NewTicker(2 * time.Second)
	sm.mu.Lock()
	if _, running := sm.tickers[contextID]; running {
		sm.mu.Unlock()
		ticker.Stop()
		return
	}
	sm.tickers[contextID] = ticker
	sm.mu.Unlock()

	go func() {
		for range ticker.C {
			state, exists := sm.LoadButtonState(contextID)
			if !exists {
				sm.stopAutoIncrement(contextID)
				return
			}
			state.Settings.Counter++
			state.LastUpdate = time.Now()
			sm.StoreButtonState(contextID, state)
			if sm.onTick != nil {
				sm.onTick(contextID, state.Settings)
			}
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
	action := streamdeck.NewAction[Settings](client, "dev.samwho.streamdeck.settings-manager")
	sm.onTick = func(contextID string, settings Settings) {
		for _, inst := range action.Contexts() {
			if sdcontext.Context(inst) == contextID {
				_ = action.SetSettings(inst, settings)
				_ = applyVisuals(inst, action, settings)
				return
			}
		}
	}

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
			for _, inst := range action.Contexts() {
				id := sdcontext.Context(inst)
				state, ok := sm.LoadButtonState(id)
				if !ok {
					state.Settings.ButtonText = "Click Me"
					state.Settings.Color = "blue"
				}
				state.Settings.Counter = 0
				sm.UpdateButtonState(id, state.Settings)
				if err := action.SetSettings(inst, state.Settings); err != nil {
					return err
				}
				if err := applyVisuals(inst, action, state.Settings); err != nil {
					return err
				}
			}
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

func accent(colorName string) color.RGBA {
	switch colorName {
	case "red":
		return color.RGBA{R: 226, G: 75, B: 74, A: 255}
	case "green":
		return color.RGBA{R: 61, G: 204, B: 122, A: 255}
	case "yellow":
		return color.RGBA{R: 230, G: 192, B: 74, A: 255}
	case "purple":
		return color.RGBA{R: 155, G: 107, B: 255, A: 255}
	default:
		return color.RGBA{R: 61, G: 126, B: 255, A: 255}
	}
}

func fillRect(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, c)
		}
	}
}

func createBackground(colorName string) image.Image {
	const size = 72
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	fillRect(img, 0, 0, size, size, color.RGBA{R: 44, G: 44, B: 44, A: 255})
	fillRect(img, 0, size-6, size, size, accent(colorName))

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	type slider struct{ x, knobY int }
	for _, s := range []slider{{20, 22}, {34, 30}, {48, 18}} {
		fillRect(img, s.x, 14, s.x+4, 50, white)
		fillRect(img, s.x-3, s.knobY, s.x+7, s.knobY+6, white)
	}
	return img
}
