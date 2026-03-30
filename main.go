package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coderj001/lazypost/config"
	"github.com/coderj001/lazypost/ui"
	"github.com/jroimartin/gocui"
)

// Views
const (
	UrlEndpointView  = "UrlEndpoint"
	ParamsView       = "Params"
	ResponseBodyView = "ResponseBody"
	HeadersView      = "Headers"
	MethodView       = "Method"
	FloatingView     = "floating"
)

var (
	AllViews = []string{
		UrlEndpointView,
		ParamsView,
		HeadersView,
		ResponseBodyView,
	}
	active        = 0
	defaultURL    = "https://httpbin.org/get"
	httpMethod    = "GET"
	httpClient    = &http.Client{Timeout: 10 * time.Second}
	httpMethods   = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	currentMethod = 0
)

// moved: all create*View logic relocated to ui/layout.go Setup()
// createFloatingEditorView remains here for floating edit (not regular layout)

func createFloatingEditorView(g *gocui.Gui, maxX, maxY int) (*gocui.View, error) {
	v, err := g.SetView(FloatingView, maxX/2-20, maxY/2-3, maxX/2+20, maxY/2+3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Title = " Edit Value "
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.SelFgColor = gocui.ColorWhite
		v.BgColor = gocui.ColorCyan
		v.Frame = true
	}
	return v, nil
}

// Layout logic moved to ui/layout.go
// func layout(g *gocui.Gui) error {}

func initKeybindings(g *gocui.Gui) error {
	cfg, err := config.LoadKeybindings()
	if err != nil {
		return fmt.Errorf("failed to load keybindings: %w", err)
	}

	warnings := cfg.Validate()
	for _, w := range warnings {
		log.Printf("keybinding warning: %s", w)
	}

	bindings := map[string]func(*gocui.Gui, *gocui.View) error{
		"quit":          quit,
		"quit-alt":      quit,
		"next-view":     nextView,
		"prev-view":     prevView,
		"send-request":  sendRequest,
		"start-editor":  startEditor,
		"switch-method": switchMethod,
	}

	for action, keyStr := range cfg.Keybindings {
		handler, ok := bindings[action]
		if !ok {
			continue
		}
		key, mod, err := config.ParseKey(keyStr)
		if err != nil {
			log.Printf("skipping invalid keybinding %s: %v", action, err)
			continue
		}
		if err := g.SetKeybinding("", key, mod, handler); err != nil {
			return fmt.Errorf("failed to set keybinding for %s: %w", action, err)
		}
	}

	return nil
}

// Parse parameters from the ParamsView
func parseParams(view *gocui.View) map[string]string {
	params := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(view.Buffer()), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			params[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return params
}

// Parse headers from the HeadersView
func parseHeaders(view *gocui.View) map[string]string {
	headers := make(map[string]string)
	if view == nil {
		return headers
	}

	lines := strings.Split(strings.TrimSpace(view.Buffer()), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return headers
}

// Helper function to format response for display
func formatResponse(resp *http.Response, body []byte) string {
	var result strings.Builder

	// Format response status
	result.WriteString(fmt.Sprintf("Status: %s\n", resp.Status))
	result.WriteString("Headers:\n")

	// Format response headers
	for key, values := range resp.Header {
		result.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}

	result.WriteString("\nBody:\n")

	// Try to format JSON for readability
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
		result.Write(prettyJSON.Bytes())
	} else {
		// If not JSON, just print the raw body
		result.Write(body)
	}

	return result.String()
}

// Switch HTTP method
func switchMethod(g *gocui.Gui, v *gocui.View) error {
	currentMethod = (currentMethod + 1) % len(httpMethods)
	httpMethod = httpMethods[currentMethod]

	urlView, err := g.View(UrlEndpointView)
	if err != nil {
		return err
	}

	urlView.Title = fmt.Sprintf("(1) %s", httpMethod)

	return nil
}

// Send HTTP request
func sendRequest(g *gocui.Gui, v *gocui.View) error {
	urlView, err := g.View(UrlEndpointView)
	if err != nil {
		return err
	}

	paramsView, err := g.View(ParamsView)
	if err != nil {
		return err
	}

	headersView, err := g.View(HeadersView)
	if err != nil {
		return err
	}

	responseView, err := g.View(ResponseBodyView)
	if err != nil {
		return err
	}

	// Clear the response view before making a new request
	responseView.Clear()
	fmt.Fprintln(responseView, "Sending request...")
	g.Update(func(*gocui.Gui) error { return nil })

	// Get URL from the URL view
	targetURL := strings.TrimSpace(urlView.Buffer())
	params := parseParams(paramsView)
	headers := parseHeaders(headersView)

	// Create request based on HTTP method
	var req *http.Request
	var reqErr error

	switch httpMethod {
	case "GET":
		// For GET requests, parameters go in the URL
		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			fmt.Fprintf(responseView, "Error parsing URL: %v", err)
			return nil
		}

		// Add query parameters
		q := parsedURL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		parsedURL.RawQuery = q.Encode()

		req, reqErr = http.NewRequest("GET", parsedURL.String(), nil)
	case "POST", "PUT", "PATCH":
		// For POST/PUT/PATCH, parameters go in the body
		jsonParams := make(map[string]interface{})
		for key, value := range params {
			jsonParams[key] = value
		}

		jsonBody, err := json.Marshal(jsonParams)
		if err != nil {
			fmt.Fprintf(responseView, "Error creating JSON body: %v", err)
			return nil
		}

		req, reqErr = http.NewRequest(httpMethod, targetURL, bytes.NewBuffer(jsonBody))

		// Default content type to JSON if not specified
		if _, hasContentType := headers["Content-Type"]; !hasContentType {
			headers["Content-Type"] = "application/json"
		}
	case "DELETE":
		req, reqErr = http.NewRequest("DELETE", targetURL, nil)
	}

	if reqErr != nil {
		fmt.Fprintf(responseView, "Error creating request: %v", reqErr)
		return nil
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(responseView, "Error sending request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(responseView, "Error reading response: %v", err)
		return nil
	}

	// Format and display the response
	responseView.Clear()
	fmt.Fprintln(responseView, formatResponse(resp, body))

	return nil
}

func startEditor(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	editorView, err := createFloatingEditorView(g, maxX, maxY)
	if err != nil {
		return err
	}

	currentView := g.CurrentView().Name()

	v.Editable = true
	v.Editor = gocui.DefaultEditor
	v.SelFgColor = gocui.ColorWhite
	v.BgColor = gocui.ColorCyan

	// Set the editor view's content to the current value of the item being edited
	switch currentView {
	case UrlEndpointView:
		editorView.Clear()
		fmt.Fprint(editorView, defaultURL)
	}

	g.SetCurrentView(FloatingView) //set focus

	// Set up keybindings for the floating editor view
	if err := g.SetKeybinding(FloatingView, gocui.KeyEnter, gocui.ModNone, saveEditedValue); err != nil {
		return err
	}
	if err := g.SetKeybinding(FloatingView, gocui.KeyEsc, gocui.ModNone, cancelEdit); err != nil {
		return err
	}

	return nil
}

func saveEditedValue(g *gocui.Gui, v *gocui.View) error {
	editedText := strings.TrimSpace(v.Buffer())
	editorViewName := v.Name()

	switch editorViewName {
	case FloatingView:
		defaultURL, _ = url.JoinPath("https://", editedText)
	}

	urlView, err := g.View(UrlEndpointView)
	if err != nil {
		return err
	}
	urlView.Clear()
	fmt.Fprintln(urlView, defaultURL)

	if err := g.DeleteKeybinding(FloatingView, gocui.KeyEnter, gocui.ModNone); err != nil {
		return err
	}
	if err := g.DeleteKeybinding(FloatingView, gocui.KeyEsc, gocui.ModNone); err != nil {
		return err
	}

	if err := g.DeleteView(FloatingView); err != nil {
		return err
	}

	g.SetCurrentView(UrlEndpointView)
	return nil
}

func cancelEdit(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteKeybinding(FloatingView, gocui.KeyEnter, gocui.ModNone); err != nil {
		return err
	}
	if err := g.DeleteKeybinding(FloatingView, gocui.KeyEsc, gocui.ModNone); err != nil {
		return err
	}

	if err := g.DeleteView(FloatingView); err != nil {
		return err
	}

	g.SetCurrentView(UrlEndpointView)
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView().Name()

	currentIdx := -1
	for i, name := range AllViews {
		if name == currentView {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(AllViews)
	g.SetCurrentView(AllViews[nextIdx])
	return nil
}

func prevView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView().Name()

	currentIdx := -1
	for i, name := range AllViews {
		if name == currentView {
			currentIdx = i
			break
		}
	}

	prevIdx := (currentIdx - 1 + len(AllViews)) % len(AllViews)
	g.SetCurrentView(AllViews[prevIdx])
	return nil
}

// quit handles Ctrl+C to exit the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	// Call deep UI Setup instead of SetManagerFunc/layout
	layout := ui.NewUILayout()
	if err := layout.Setup(g); err != nil {
		log.Panicln("UI setup failed:", err)
	}

	//KeyBinding
	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
