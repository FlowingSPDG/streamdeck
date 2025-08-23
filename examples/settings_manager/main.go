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

	"github.com/FlowingSPDG/streamdeck"
	"github.com/puzpuzpuz/xsync/v3"
)

// Settings ボタンの設定を表す構造体
type Settings struct {
	Counter       int    `json:"counter"`
	ButtonText    string `json:"buttonText"`
	Color         string `json:"color"`
	AutoIncrement bool   `json:"autoIncrement"`
}

// ButtonState ボタンの状態を表す構造体
type ButtonState struct {
	Settings   Settings
	LastUpdate time.Time
	IsActive   bool
}

// SettingsManager 設定管理を行う構造体
type SettingsManager struct {
	// contextをキーとしたボタン状態の保存
	buttonStates *xsync.MapOf[string, ButtonState]
	// 自動インクリメント用のticker管理
	tickers *xsync.MapOf[string, *time.Ticker]
	// ロック用のmutex
	mu sync.RWMutex
}

// NewSettingsManager 新しい設定マネージャーを作成
func NewSettingsManager() *SettingsManager {
	return &SettingsManager{
		buttonStates: xsync.NewMapOf[string, ButtonState](),
		tickers:      xsync.NewMapOf[string, *time.Ticker](),
	}
}

// StoreButtonState ボタンの状態を保存
func (sm *SettingsManager) StoreButtonState(contextID string, state ButtonState) {
	sm.buttonStates.Store(contextID, state)
	log.Printf("Stored button state for context %s: %+v", contextID, state)
}

// LoadButtonState ボタンの状態を読み込み
func (sm *SettingsManager) LoadButtonState(contextID string) (ButtonState, bool) {
	return sm.buttonStates.Load(contextID)
}

// DeleteButtonState ボタンの状態を削除
func (sm *SettingsManager) DeleteButtonState(contextID string) {
	sm.buttonStates.Delete(contextID)
	sm.stopAutoIncrement(contextID)
	log.Printf("Deleted button state for context %s", contextID)
}

// UpdateButtonState ボタンの状態を更新
func (sm *SettingsManager) UpdateButtonState(contextID string, settings Settings) {
	state, exists := sm.LoadButtonState(contextID)
	if !exists {
		state = ButtonState{
			Settings:   settings,
			LastUpdate: time.Now(),
			IsActive:   true,
		}
	} else {
		state.Settings = settings
		state.LastUpdate = time.Now()
	}

	sm.StoreButtonState(contextID, state)

	// 自動インクリメントの設定を更新
	if settings.AutoIncrement {
		sm.startAutoIncrement(contextID)
	} else {
		sm.stopAutoIncrement(contextID)
	}
}

// startAutoIncrement 自動インクリメントを開始
func (sm *SettingsManager) startAutoIncrement(contextID string) {
	// 既存のtickerがあれば停止
	sm.stopAutoIncrement(contextID)

	ticker := time.NewTicker(2 * time.Second)
	sm.tickers.Store(contextID, ticker)

	go func() {
		for range ticker.C {
			// This runs outside of StreamDeck event.
			state, exists := sm.LoadButtonState(contextID)
			if !exists || !state.IsActive {
				ticker.Stop()
				sm.tickers.Delete(contextID)
				return
			}

			// カウンターをインクリメント
			state.Settings.Counter++
			state.LastUpdate = time.Now()
			sm.StoreButtonState(contextID, state)

			log.Printf("Auto-incremented counter for context %s: %d", contextID, state.Settings.Counter)
		}
	}()
}

// stopAutoIncrement 自動インクリメントを停止
func (sm *SettingsManager) stopAutoIncrement(contextID string) {
	if ticker, ok := sm.tickers.Load(contextID); ok {
		ticker.Stop()
		sm.tickers.Delete(contextID)
		log.Printf("Stopped auto-increment for context %s", contextID)
	}
}

// GetAllButtonStates 全てのボタン状態を取得（デバッグ用）
func (sm *SettingsManager) GetAllButtonStates() map[string]ButtonState {
	result := make(map[string]ButtonState)
	sm.buttonStates.Range(func(key string, value ButtonState) bool {
		result[key] = value
		return true
	})
	return result
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := streamdeck.NewClient(ctx, params)
	settingsManager := NewSettingsManager()

	setup(client, settingsManager)

	// デバッグ用：定期的に全てのボタン状態をログ出力
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			states := settingsManager.GetAllButtonStates()
			log.Printf("Current button states: %+v", states)
		}
	}()

	return client.Run(ctx)
}

func setup(client *streamdeck.Client, sm *SettingsManager) {
	action := client.Action("dev.samwho.streamdeck.settings_manager")

	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.WillAppearPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return fmt.Errorf("failed to unmarshal WillAppear payload: %w", err)
		}

		// デフォルト設定
		if p.Settings.ButtonText == "" {
			p.Settings.ButtonText = "Click Me"
		}
		if p.Settings.Color == "" {
			p.Settings.Color = "blue"
		}

		// ボタン状態を保存
		sm.UpdateButtonState(event.Context, p.Settings)

		// 背景画像を設定
		bg, err := streamdeck.Image(createBackground(p.Settings.Color))
		if err != nil {
			return fmt.Errorf("failed to create background image: %w", err)
		}

		if err := client.SetImage(ctx, bg, streamdeck.HardwareAndSoftware); err != nil {
			return fmt.Errorf("failed to set image: %w", err)
		}

		// タイトルを設定
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.WillDisappearPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return fmt.Errorf("failed to unmarshal WillDisappear payload: %w", err)
		}

		// ボタン状態を削除
		sm.DeleteButtonState(event.Context)

		// カウンターをリセット
		p.Settings.Counter = 0
		return client.SetSettings(ctx, p.Settings)
	})

	action.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.KeyDownPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return fmt.Errorf("failed to unmarshal KeyDown payload: %w", err)
		}

		// カウンターをインクリメント
		p.Settings.Counter++

		// ボタン状態を更新
		sm.UpdateButtonState(event.Context, p.Settings)

		// 設定を保存
		if err := client.SetSettings(ctx, p.Settings); err != nil {
			return fmt.Errorf("failed to set settings: %w", err)
		}

		// タイトルを更新
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	action.RegisterHandler(streamdeck.DidReceiveSettings, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var p streamdeck.DidReceiveSettingsPayload[Settings]
		if err := event.UnmarshalPayload(&p); err != nil {
			return fmt.Errorf("failed to unmarshal DidReceiveSettings payload: %w", err)
		}

		// ボタン状態を更新
		sm.UpdateButtonState(event.Context, p.Settings)

		// 背景画像を更新
		bg, err := streamdeck.Image(createBackground(p.Settings.Color))
		if err != nil {
			return fmt.Errorf("failed to create background image: %w", err)
		}

		if err := client.SetImage(ctx, bg, streamdeck.HardwareAndSoftware); err != nil {
			return fmt.Errorf("failed to set image: %w", err)
		}

		// タイトルを更新
		title := fmt.Sprintf("%s\n%d", p.Settings.ButtonText, p.Settings.Counter)
		return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
	})

	action.RegisterHandler(streamdeck.SendToPlugin, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		var payload map[string]interface{}
		if err := event.UnmarshalPayload(&payload); err != nil {
			return fmt.Errorf("failed to unmarshal SendToPlugin payload: %w", err)
		}

		// Property Inspectorからのメッセージを処理
		if action, ok := payload["action"].(string); ok {
			switch action {
			case "getAllStates":
				// 全てのボタン状態をProperty Inspectorに送信
				states := sm.GetAllButtonStates()
				return client.SendToPropertyInspector(ctx, map[string]interface{}{
					"action": "allStates",
					"states": states,
				})
			case "resetAll":
				// 全てのボタンのカウンターをリセット
				sm.buttonStates.Range(func(key string, value ButtonState) bool {
					value.Settings.Counter = 0
					sm.StoreButtonState(key, value)
					return true
				})
				return client.SendToPropertyInspector(ctx, map[string]interface{}{
					"action": "resetComplete",
				})
			}
		}

		return nil
	})
}

// createBackground 色に応じた背景画像を作成
func createBackground(colorName string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))

	var bgColor color.Color
	switch colorName {
	case "red":
		bgColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "green":
		bgColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	case "blue":
		bgColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	case "yellow":
		bgColor = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case "purple":
		bgColor = color.RGBA{R: 128, G: 0, B: 128, A: 255}
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
