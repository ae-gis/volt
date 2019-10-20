package http

import (
        "fmt"
        "html"
        "net/http"
        "os"
        "sync"
        "testing"

        "github.com/go-chi/chi"
        "github.com/spf13/cobra"
        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
)

const (
        Port         int = 8081
        ReadTimeout  int = 5
        WriteTimeout int = 100
)

func TestNewHttp(t *testing.T) {
        var wg sync.WaitGroup
        stop := make(chan bool)
        wg.Add(1)
        defer wg.Wait()

        httpSignal := NewHTTPCmdSignaled(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
                assert.NoError(t, err)
        }), Port, ReadTimeout, WriteTimeout, stop)

        go func() {
                defer wg.Done()
                err := httpSignal.GetBaseCmd().Execute()
                assert.NoError(t, err)
        }()

        stop <- true
}

func TestHttp(t *testing.T) {
        var wg sync.WaitGroup
        assert := require.New(t)
        stop := make(chan bool)
        wg.Add(1)

        cmd := NewHTTPCmdSignaled(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
                assert.NoError(err)
        }), Port, ReadTimeout, WriteTimeout, stop).GetBaseCmd()

        go func() {
                defer wg.Done()
                _, err := cmd.ExecuteC()
                assert.NoError(err)
        }()

        stop <- true
        wg.Wait()
}

func TestNewHttpCmdWithFilename(t *testing.T) {
        var wg sync.WaitGroup
        wg.Add(1)
        defer wg.Wait()

        stop := make(chan bool)

        cmd := NewHTTPCmdSignaled(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
                assert.NoError(t, err)
        }), Port, ReadTimeout, WriteTimeout, stop).GetBaseCmd()
        go func() {
                defer wg.Done()
                _, err := cmd.ExecuteC()
                assert.NoError(t, err)
        }()

        stop <- true
        os.Args = []string{""}
}

func TestListenAndServe(t *testing.T) {
        var (
                err error
                wg  sync.WaitGroup
        )
        stop := make(chan bool)

        cc := &httpCmd{stop: stop}
        cc.BaseCmd = &cobra.Command{
                Use:   "http",
                Short: "Used to run the http service",
                RunE: func(cmd *cobra.Command, args []string) (err error) {
                        mux := chi.NewMux()
                        return cc.serve(mux)
                },
        }

        wg.Add(1)
        go func() {
                defer wg.Done()
                err = cc.BaseCmd.Execute()
        }()
        assert.NoError(t, err)
        stop <- true
        wg.Wait()
}
