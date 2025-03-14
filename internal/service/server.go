package server

import (
	"context"
	"errors"
	"fmt"

	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/api/middleware"
	v1 "github.com/DillonEnge/jolt/internal/api/v1"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/internal/sessions"
	"github.com/DillonEnge/jolt/templates"
	"github.com/DillonEnge/seaweedfs-go-client"
	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

func Start(address string, dbPool *pgxpool.Pool, nc *nats.Conn, config *api.Config) func(context.Context) error {
	sm := sessions.NewSessionManager()

	authClient := auth.NewClient(config)

	fsClient := seaweedfs.NewClient(seaweedfs.Config{
		MasterURL:  config.SeaweedFS.MasterURL,
		VolumesURL: config.SeaweedFS.VolumesURL,
	})

	db := database.New(dbPool)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /navbar", makeH(v1.HandleNavbar(sm, authClient)))

	mux.Handle("GET /search", templ.Handler(templates.Search()))

	mux.HandleFunc("GET /listings/popular", makeH(v1.HandlePopularListings(db, authClient, sm)))
	mux.HandleFunc("GET /listings", makeH(v1.HandleListings(db, authClient, sm)))
	mux.HandleFunc("POST /listings", makeH(v1.HandlePostListings(db, fsClient, config)))
	mux.HandleFunc("DELETE /listings", makeH(v1.HandleDeleteListings(dbPool)))
	mux.HandleFunc("PATCH /listings", makeH(v1.HandlePatchListing(db)))

	mux.HandleFunc("GET /create-listing", makeH(v1.HandleCreateListing(sm, authClient)))

	mux.Handle("GET /negotiations", makeH(v1.HandleNegotiations(dbPool, authClient, sm)))
	mux.Handle("POST /negotiations", makeH(v1.HandlePostNegotiation(db, authClient, sm)))

	mux.Handle("GET /chat", makeH(v1.HandleChat(dbPool, sm, authClient, config)))

	mux.Handle("GET /ws/messages", makeH(v1.HandleMessageWS(db, authClient, nc, sm)))
	mux.Handle("GET /messages", makeH(v1.HandleMessages(dbPool)))
	mux.Handle("POST /messages", makeH(v1.HandlePostMessage(dbPool, sm, authClient)))

	mux.Handle(
		"GET /static/",
		middleware.NoCache(
			http.StripPrefix("/static/",
				http.FileServer(http.Dir("./templates/static")),
			),
		),
	)

	mux.HandleFunc("GET /my-listings", makeH(v1.HandleMyListings(dbPool, authClient, sm)))

	mux.HandleFunc("GET /loader", makeH(v1.HandleLoader()))

	mux.HandleFunc("GET /signin", makeH(v1.HandleSignin(sm, authClient)))

	mux.HandleFunc("GET /signout", makeH(v1.HandleLogout(sm)))

	mux.HandleFunc("GET /", makeH(v1.HandleBase(sm, authClient, config)))

	h := middleware.NewHandlerWithMiddleware(
		mux,
		middleware.Logger,
		sm.LoadAndSave,
	)

	s := &http.Server{
		Addr:              address,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("Listening...", "address", address)
		err := s.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return s.Shutdown
}

func Service(ctx context.Context, dbPool *pgxpool.Pool, nc *nats.Conn, config *api.Config) (func(), error) {
	shutdown := Start(fmt.Sprintf(":%d", config.Port), dbPool, nc, config)

	stopService := func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}

	return stopService, nil
}

func makeH(h api.HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("encountered error while serving request", "err", err)
			w.WriteHeader(err.Status)
			errJSON := fmt.Sprintf(
				`{"error": "%s"}`,
				strings.ReplaceAll(err.Error(), `"`, `'`),
			)
			w.Write([]byte(errJSON))
		}
	}
}
