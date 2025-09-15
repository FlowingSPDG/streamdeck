# Settings Manager Example

This example demonstrates how to use goroutine-safe `xsync.MapOf` to store unmarshaled Payload in a Map with context as the key.

## Features

- **Goroutine-safe state management**: Uses `xsync.MapOf` to safely manage state of multiple button instances
- **Payload Unmarshaling**: Uses `event.UnmarshalPayload()` to safely parse Payload
- **Dynamic settings management**: Dynamically change button text, color, and auto-increment functionality
- **Auto-increment**: Automatically increments counter every 2 seconds when enabled
- **Property Inspector integration**: Display all button states and bulk reset functionality

## Key Implementation Points

### 1. Using xsync.MapOf

```go
type SettingsManager struct {
    // Store button states with context as key
    buttonStates *xsync.MapOf[string, ButtonState]
    // Manage tickers for auto-increment
    tickers *xsync.MapOf[string, *time.Ticker]
    // Mutex for locking
    mu sync.RWMutex
}
```

### 2. Safe Payload Unmarshaling

```go
var p streamdeck.WillAppearPayload[Settings]
if err := event.UnmarshalPayload(&p); err != nil {
    return fmt.Errorf("failed to unmarshal WillAppear payload: %w", err)
}
```

### 3. State Storage and Retrieval

```go
// Store state
sm.StoreButtonState(contextID, state)

// Load state
state, exists := sm.LoadButtonState(contextID)
```

### 4. Auto-increment Functionality

```go
func (sm *SettingsManager) startAutoIncrement(contextID string) {
    ticker := time.NewTicker(2 * time.Second)
    sm.tickers.Store(contextID, ticker)
    
    go func() {
        for range ticker.C {
            state, exists := sm.LoadButtonState(contextID)
            if !exists || !state.IsActive {
                ticker.Stop()
                sm.tickers.Delete(contextID)
                return
            }
            
            // Increment counter
            state.Settings.Counter++
            sm.StoreButtonState(contextID, state)
        }
    }()
}
```

## Usage

1. **Build**: `go build -o settings_manager main.go`
2. **Deploy to Stream Deck**: Place the built file in the Stream Deck plugin directory
3. **Configure via Property Inspector**: Set button text, color, and auto-increment settings
4. **Test**: Click the button to verify counter increments

## Settings

- **Button Text**: Text to display on the button
- **Background Color**: Button background color (blue, red, green, yellow, purple)
- **Auto Increment**: Whether to automatically increment counter every 2 seconds

## Property Inspector Features

- **Save Settings**: Save current settings
- **Get All States**: Display all active button states
- **Reset All Counters**: Reset all button counters to 0

## Technical Benefits

1. **Goroutine Safety**: `xsync.MapOf` allows safe access from multiple goroutines with better performance
2. **Type Safety**: Type-safe maps with generics support, properly handles `any` type Payload
3. **Memory Leak Prevention**: Properly cleans up resources when buttons are removed  
4. **Debug Features**: Periodically logs all button states
5. **Performance**: Better performance than `sync.Map` due to specialized concurrent data structures

## Notes

- `xsync.MapOf` requires Go 1.18+ for generics support
- Auto-increment functionality runs in independent goroutines for each button
- Property Inspector communication uses `SendToPlugin`/`SendToPropertyInspector`
- Uses [puzpuzpuz/xsync](https://github.com/puzpuzpuz/xsync) library for high-performance concurrent data structures
