package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Global window instance
var window *Window

// setupWebServer starts the web server for serving static content and returns the port number
// Returns 0 and error if failed to start server
func setupWebServer() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to listen: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	
	go func() {
		e := echo.New()
		e.Logger.SetLevel(log.OFF)
		e.Logger.SetOutput(io.Discard)
		e.Static("/", "app")
		
		if err := e.Server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
		}
	}()
	
	return port, nil
}

func main() {
	// Start web server and get random port
	port, err := setupWebServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start web server: %v\n", err)
		os.Exit(1)
	}
	
	// Check for debug mode
	debug := false
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		debug = true
	} else if len(os.Args) > 1 {
		// Exit if unknown arguments
		os.Exit(0)
	}
	
	// Create and initialize window
	window = NewWindow(debug)
	defer window.webview.Destroy()
	
	// Initialize webview
	window.initWebView()
	
	// Bind JavaScript functions
	window.bindJavaScriptFunctions()
	
	// Navigate to local server
	window.webview.Navigate(fmt.Sprintf("http://127.0.0.1:%d/", port))
	
	// Center window on screen
	window.center()
	
	// Run the application
	window.webview.Run()
}