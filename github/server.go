package github

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/google/go-github/v56/github"

	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/signalutil"
)

func Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", eventHandler)

	server := &http.Server{Handler: mux}

	ln, err := net.Listen("tcp", "")
	if err != nil {
		return err
	}

	slog.Info("listening on %s", ln.Addr())

	go func() {
		if e := server.Serve(ln); e != nil && !errors.Is(e, http.ErrServerClosed) {
			slog.Error("serve error: %v", e)
		}
	}()

	stop := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)

		slog.Info("server stopped")
	}

	slog.Info("server started")

	<-signalutil.Defer(stop).Done()

	return nil
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, nil) // todo
	if err != nil {
		return
	}

	tpe, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return
	}

	switch event := tpe.(type) {
	case *github.Event:
		_ = event.Type
	}
}
