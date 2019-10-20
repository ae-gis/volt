package http

import (
        "errors"
        "fmt"
        "net/http"
        "net/url"
        "os"
        "os/signal"
        "sync"
        "syscall"
        "time"

        "github.com/ae-gis/volt/log"
        "github.com/spf13/cobra"
)

func WelkomText() string {
        return `
========================================================================================

           __      ___   ___ ___ ___ _  _ 
  _  _ ___ \ \    / /_\ | _ \ _ \_ _| \| |
 | || (_-<  \ \/\/ / _ \|   /  _/| || . |
 | .,_/__/   \_/\_/_/ \_\_|_\_| |___|_|\_|
 |_|
========================================================================================
- port    : %d
-----------------------------------------------------------------------------------------`
}

type HTTPCmd interface {
        serve(router http.Handler) error
        server(cmd *cobra.Command, args []string) (err error)
        GetBaseCmd() *cobra.Command
}

type httpCmd struct {
        stop <-chan bool

        Port         int
        ReadTimeout  int
        WriteTimeout int
        Filename     string
        BaseCmd      *cobra.Command
        handler      http.Handler
        srv          *Server
}

// Init Object HTTP Command
func NewHTTPCmd(
        port,
        readTimeout,
        writeTimeout int,
        handler http.Handler,
) HTTPCmd {
        return NewHTTPCmdSignaled(
                handler,
                port,
                readTimeout,
                writeTimeout, nil)
}

// Init Object HTTP Command
func NewHTTPCmdSignaled(
        handler http.Handler,
        port,
        readTimeout,
        writeTimeout int,
        stop <-chan bool,
) HTTPCmd {
        cc := &httpCmd{stop: stop}
        cc.handler = handler
        cc.Port = port
        cc.ReadTimeout = readTimeout
        cc.WriteTimeout = writeTimeout
        cc.BaseCmd = &cobra.Command{
                Use:   "http",
                Short: "Used to run the http service",
                RunE:  cc.server,
        }
        return cc
}

func (h *httpCmd) server(cmd *cobra.Command, args []string) (err error) {
        if h.handler == nil {
                panic(errors.New("handler function is nil"))
        }

        // Description µ micro service
        fmt.Println(
                fmt.Sprintf(
                        WelkomText(),
                        h.Port,
                ))
        return h.serve(h.handler)
}

func (h *httpCmd) serve(router http.Handler) error {
        addrURL := url.URL{Scheme: "http", Host: fmt.Sprintf(":%d", h.Port)}
        log.Info(fmt.Sprintf("started server %s", addrURL.String()))
        h.srv = StartWebServer(
                addrURL,
                h.ReadTimeout,
                h.WriteTimeout,
                router,
        )
        defer h.srv.Stop()
        sc := make(chan os.Signal, 10)
        signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
        select {
        case s := <-sc:
                log.Info(fmt.Sprintf("shutting down server with signal %q", s.String()))
        case <-h.stop:
                log.Info("shutting down server with stop channel")
        case <-h.srv.StopNotify():
                log.Info("shutting down server with stop signal")

        }
        return nil
}

func (h *httpCmd) GetBaseCmd() *cobra.Command {
        return h.BaseCmd
}

// Server warps http.Server.
type Server struct {
        mu         sync.RWMutex
        addrURL    url.URL
        httpServer *http.Server

        stopc chan struct{}
        donec chan struct{}
}

// StartWebServer starts a web server
func StartWebServer(addr url.URL, readTimeout, writeTimeout int, handler http.Handler) *Server {
        stopc := make(chan struct{})
        srv := &Server{
                addrURL: addr,
                httpServer: &http.Server{
                        Addr:         addr.Host,
                        Handler:      handler,
                        ReadTimeout:  time.Duration(readTimeout) * time.Second,
                        WriteTimeout: time.Duration(writeTimeout) * time.Second,
                },
                stopc: stopc,
                donec: make(chan struct{}),
        }
        go func() {
                defer func() {
                        if err := recover(); err != nil {
                                log.Warn(
                                        "shutting down server with err ",
                                        log.Field("error", fmt.Sprintf(`(%v)`, err)),
                                )
                                os.Exit(0)
                        }
                        close(srv.donec)
                }()
                if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        log.Fatal(
                                "shutting down server with err ",
                                log.Field("error", err),
                        )
                }
        }()
        return srv
}

// StopNotify returns receive-only stop channel to notify the server has stopped.
func (srv *Server) StopNotify() <-chan struct{} {
        return srv.stopc
}

// Stop stops the server. Useful for testing.
func (srv *Server) Stop() {
        log.Warn(fmt.Sprintf("stopping server %s", srv.addrURL.String()))
        srv.mu.Lock()
        if srv.httpServer == nil {
                srv.mu.Unlock()
                return
        }
        close(srv.stopc)
        _ = srv.httpServer.Close()
        <-srv.donec
        srv.mu.Unlock()
        log.Warn(fmt.Sprintf("stopped server %s", srv.addrURL.String()))
}
