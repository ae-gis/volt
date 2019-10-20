
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
