
- Use a Log

```go
package main
import (
    "gitlab.warungpintar.co/back-end/libwp/log"
)

func main(){
    
    log.Error("Error",
        log.Field("error", "error for log"),
    )
    log.Debug("Debug",
        log.Field("debug", "debug for log"),
        log.Field("message", "message debug for log"),
    )
    log.Info("Info",
        log.Field("info", "info for log"),
        log.Field("message", "message info for log"),
    )
    log.Warn("Warn",
        log.Field("warning", "warning for log"),
        log.Field("message", "message warning for log"),
    )
}

``` 

- Use a Http Serve

```go
package main
import (
     "fmt"
     "html"
     "net/http"
     httpUtil "gitlab.warungpintar.co/back-end/libwp/http"
)

func main(){
        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                   _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
                   if err != nil {
                           panic(err)
                   }
        })
        s := httpUtil.NewHTTPCmd(8009, 100, 10, handler).GetBaseCmd()
        if err := s.Execute(); err != nil {
                panic(err)
        }
}

``` 


- Use a Breaker

```go
package main
import (
    "net/http"
    "github.com/hashicorp/go-cleanhttp"
    "gitlab.warungpintar.co/back-end/libwp/breaker"
    "gitlab.warungpintar.co/back-end/libwp/log"
)

func main() {
    endpoint := `www.google.com`
    cb := breaker.NewBreaker(
        "",
        100,
        10,
    )
    var res *http.Response
    err := cb.Execute(func() (err error) {
        client := cleanhttp.DefaultClient()
        req, _ := http.NewRequest(http.MethodGet,
            endpoint, nil)
        res, err = client.Do(req)
        return err
    })
    if err != nil {
        log.Error("Error",
            log.Field("error", err.Error()),
        )
        panic(err)
    }
}
```

