package main

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	webview "github.com/webview/webview_go"
)

const (
	// Window state constants
	SW_MAXIMIZE = 3
	SW_MINIMIZE = 6
	SW_RESTORE  = 9

	// Default window dimensions
	defaultWidth  = 1280
	defaultHeight = 720
	defaultDPI    = 96.0
)

// Window represents the application window
type Window struct {
	webview webview.WebView
	handle  unsafe.Pointer
	width   int
	height  int
	scale   float64
}

// NewWindow creates and initializes a new window
func NewWindow(debug bool) *Window {
	w := &Window{
		webview: webview.New(debug),
		width:   defaultWidth,
		height:  defaultHeight,
		scale:   1.0,
	}
	
	w.handle = w.webview.Window()
	
	// Get system DPI for proper scaling
	if dpi, err := GetDpiForSystem(); err == nil && dpi > 0 {
		w.scale = float64(dpi) / defaultDPI
	}
	
	// Set default window title and size
	w.webview.SetTitle("My Application")
	w.resize(defaultWidth, defaultHeight)
	
	return w
}

// initWebView initializes the webview with necessary JavaScript
func (w *Window) initWebView() {
	// Disable refresh and context menu in production mode
	w.webview.Init("document.addEventListener('keydown', function(event) {" +
		"if ((event.ctrlKey && event.key === 'r') || (event.ctrlKey && event.key === 'R') || event.key === 'F5') {" +
		"event.preventDefault();" +
		"}});")
	w.webview.Init("document.addEventListener('contextmenu', function(event) {event.preventDefault();});")
	
	// Initialize window API
	w.webview.Init("window.mata={}")
	w.webview.Init("window.mata.win={}")
}

// bindJavaScriptFunctions binds Go functions to JavaScript
func (w *Window) bindJavaScriptFunctions() {
	// Bind Go functions to JavaScript
	bindings := map[string]interface{}{
		"mata_win_minimize": w.minimize,
		"mata_win_maximize": w.maximize,
		"mata_win_restore":  w.restore,
		"mata_win_close":    w.close,
		"mata_win_center":   w.center,
		"mata_win_title":    w.setTitle,
		"mata_win_resize":   w.resize,
		"mata_alert":        w.showAlert,
	}

	for name, fn := range bindings {
		w.webview.Bind(name, fn)
	}

	// Expose JavaScript API
	w.webview.Init(`
		window.mata.win = {
			close: mata_win_close,
			minimize: mata_win_minimize,
			maximize: mata_win_maximize,
			restore: mata_win_restore,
			center: mata_win_center,
			title: mata_win_title,
			resize: mata_win_resize
		};
		window.mata.alert = mata_alert;
	`)
}

// showAlert displays an alert dialog with the given message
func (w *Window) showAlert(msg string) {
	w.webview.Eval(fmt.Sprintf("alert('%s')", strings.ReplaceAll(msg, "'", "\\'")))
}

// Window management methods

// resize sets the window dimensions
func (w *Window) resize(width, height int) {
	if width > 0 && height > 0 {
		w.width = width
		w.height = height
		w.webview.SetSize(width, height, webview.HintNone)
	}
}

// setTitle updates the window title
func (w *Window) setTitle(title string) {
	w.webview.SetTitle(title)
}

// minimize minimizes the window
func (w *Window) minimize() {
	if err := ShowWindow(w.handle, SW_MINIMIZE); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to minimize window: %v\n", err)
	}
}

// maximize maximizes the window
func (w *Window) maximize() {
	if err := ShowWindow(w.handle, SW_MAXIMIZE); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to maximize window: %v\n", err)
	}
}

// restore restores the window from minimized/maximized state
func (w *Window) restore() {
	if err := ShowWindow(w.handle, SW_RESTORE); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to restore window: %v\n", err)
	}
}

// close terminates the application
func (w *Window) close() {
	w.webview.Terminate()
}

// center positions the window in the center of the screen
func (w *Window) center() {
	screenWidth, err1 := GetSystemMetrics(SM_CXSCREEN)
	screenHeight, err2 := GetSystemMetrics(SM_CYSCREEN)
	
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Failed to get screen metrics: %v, %v\n", err1, err2)
		return
	}
	
	realWidth := float64(w.width) * w.scale
	realHeight := float64(w.height) * w.scale
	x := (screenWidth - int(realWidth)) / 2
	y := (screenHeight - int(realHeight)) / 2
	
	if err := SetWindowPos(w.handle, 0, x, y, int(realWidth), int(realHeight), 0); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to center window: %v\n", err)
	}
}

// Windows API functions are defined in winapi.go