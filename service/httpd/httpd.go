// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package httpd // import "miniflux.app/service/httpd"

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"miniflux.app/api"
	"miniflux.app/config"
	"miniflux.app/fever"
	"miniflux.app/logger"
	"miniflux.app/reader/feed"
	"miniflux.app/storage"
	"miniflux.app/ui"
	"miniflux.app/worker"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

// Serve starts a new HTTP server.
func Serve(store *storage.Storage, pool *worker.Pool, feedHandler *feed.Handler) *http.Server {
	certFile := config.Opts.CertFile()
	keyFile := config.Opts.CertKeyFile()
	certDomain := config.Opts.CertDomain()
	certCache := config.Opts.CertCache()
	listenAddr := config.Opts.ListenAddr()
	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      setupHandler(store, feedHandler, pool),
	}

	switch {
	case os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()):
		startSystemdSocketServer(server)
	case strings.HasPrefix(listenAddr, "/"):
		startUnixSocketServer(server, listenAddr)
	case certDomain != "" && certCache != "":
		config.Opts.HTTPS = true
		startAutoCertTLSServer(server, certDomain, certCache)
	case certFile != "" && keyFile != "":
		config.Opts.HTTPS = true
		server.Addr = listenAddr
		startTLSServer(server, certFile, keyFile)
	default:
		server.Addr = listenAddr
		startHTTPServer(server)
	}

	return server
}

func startSystemdSocketServer(server *http.Server) {
	go func() {
		f := os.NewFile(3, "systemd socket")
		listener, err := net.FileListener(f)
		if err != nil {
			logger.Fatal(`Unable to create listener from systemd socket: %v`, err)
		}

		logger.Info(`Listening on systemd socket`)
		if err := server.Serve(listener); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}()
}

func startUnixSocketServer(server *http.Server, socketFile string) {
	os.Remove(socketFile)

	go func(sock string) {
		listener, err := net.Listen("unix", sock)
		if err != nil {
			logger.Fatal(`Server failed to start: %v`, err)
		}
		defer listener.Close()

		if err := os.Chmod(sock, 0666); err != nil {
			logger.Fatal(`Unable to change socket permission: %v`, err)
		}

		logger.Info(`Listening on Unix socket %q`, sock)
		if err := server.Serve(listener); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}(socketFile)
}

func startAutoCertTLSServer(server *http.Server, certDomain, certCache string) {
	server.Addr = ":https"
	certManager := autocert.Manager{
		Cache:      autocert.DirCache(certCache),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(certDomain),
	}

	// Handle http-01 challenge.
	s := &http.Server{
		Handler: certManager.HTTPHandler(nil),
		Addr:    ":http",
	}
	go s.ListenAndServe()

	go func() {
		logger.Info(`Listening on %q by using auto-configured certificate for %q`, server.Addr, certDomain)
		if err := server.Serve(certManager.Listener()); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}()
}

func startTLSServer(server *http.Server, certFile, keyFile string) {
	// See https://blog.cloudflare.com/exposing-go-on-the-internet/
	// And https://wiki.mozilla.org/Security/Server_Side_TLS
	server.TLSConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	go func() {
		logger.Info(`Listening on %q by using certificate %q and key %q`, server.Addr, certFile, keyFile)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}()
}

func startHTTPServer(server *http.Server) {
	go func() {
		logger.Info(`Listening on %q without TLS`, server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(`Server failed to start: %v`, err)
		}
	}()
}

func setupHandler(store *storage.Storage, feedHandler *feed.Handler, pool *worker.Pool) *mux.Router {
	router := mux.NewRouter()

	if config.Opts.BasePath() != "" {
		router = router.PathPrefix(config.Opts.BasePath()).Subrouter()
	}

	router.Use(middleware)

	fever.Serve(router, store)
	api.Serve(router, store, feedHandler)
	ui.Serve(router, store, pool, feedHandler)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Name("healthcheck")

	return router
}
