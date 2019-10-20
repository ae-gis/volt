
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
