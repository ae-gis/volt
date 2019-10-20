package http

import (
        "fmt"
        "html"
        "net/http"
)

func main() {
        s := NewHTTPCmd(
                8009,
                100,
                10,
                http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
                        if err != nil {
                                panic(err)
                        }
                }),
        ).GetBaseCmd()
        if err := s.Execute(); err != nil {
                panic(err)
        }
}
