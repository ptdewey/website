package atproto

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/bluesky-social/indigo/atproto/auth/oauth"

	"github.com/ptdewey/cedar/internal/config"
)

var oauthScopes = []string{
	"atproto",
	"repo:site.standard.publication",
	"repo:site.standard.document",
}

type callbackResult struct {
	Code  string
	ISS   string
	State string
	Err   error
}

// RunOAuthFlow performs the ATProto OAuth loopback flow and saves credentials to disk.
func RunOAuthFlow(cfg *config.Config) error {
	ctx := context.Background()

	// Find a free port for the loopback callback server.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("starting callback listener: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	store := newFileAuthStore()
	oauthCfg := oauth.NewLocalhostConfig(redirectURI, oauthScopes)
	app := oauth.NewClientApp(&oauthCfg, store)

	fmt.Printf("Starting OAuth flow for %s...\n", cfg.ATProto.Handle)
	authURL, err := app.StartAuthFlow(ctx, cfg.ATProto.Handle)
	if err != nil {
		return fmt.Errorf("starting auth flow: %w", err)
	}

	// Start loopback callback server.
	codeCh := make(chan callbackResult, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		errParam := r.URL.Query().Get("error")
		if errParam != "" {
			errDesc := r.URL.Query().Get("error_description")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Authorization failed: %s - %s", errParam, errDesc)
			codeCh <- callbackResult{Err: fmt.Errorf("authorization failed: %s - %s", errParam, errDesc)}
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body><h1>Authorization successful!</h1><p>You can close this tab.</p></body></html>")
		codeCh <- callbackResult{
			Code:  r.URL.Query().Get("code"),
			ISS:   r.URL.Query().Get("iss"),
			State: r.URL.Query().Get("state"),
		}
	})

	srv := &http.Server{Handler: mux}
	go srv.Serve(listener) //nolint:errcheck

	fmt.Printf("\nOpening browser for authorization...\n")
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser. Please visit:\n%s\n", authURL)
	}

	fmt.Println("Waiting for authorization...")
	result := <-codeCh

	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx) //nolint:errcheck

	if result.Err != nil {
		return result.Err
	}

	// Hand the callback parameters to the SDK to complete the flow.
	params := url.Values{}
	params.Set("code", result.Code)
	params.Set("iss", result.ISS)
	params.Set("state", result.State)

	sessData, err := app.ProcessCallback(ctx, params)
	if err != nil {
		return fmt.Errorf("processing callback: %w", err)
	}

	fmt.Printf("\nAuthenticated as %s\n", sessData.AccountDID)
	fmt.Println("Credentials saved to " + authStatePath)
	return nil
}

// openBrowser opens url in the system default browser.
func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", rawURL)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
