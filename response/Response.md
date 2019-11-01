
- Response 

```go
package main

import (
    "strconv"
    "http"
    "github.com/ae-gis/volt/response"
)

func main() {
    handler := func(w http.ResponseWriter, r *http.Request) {
       
        body := response.New() // initial new response struct
        if r.Method != http.MethodPost {
            body.SetErrors(response.ErrorContext{
                Code:    strconv.Itoa(http.StatusMethodNotAllowed),
                Message: http.StatusText(http.StatusMethodNotAllowed),
            })
            log.Error("CourierDasboard",
                log.Field("Error", errors.New("request Method Not Allowed")))
            response.Status(r, http.StatusMethodNotAllowed) 
            response.WriteJSON(w, r, body) // response to JSON
            return
        }
        body.SetData(map[string]interface{}{
            "Code": strconv.Itoa(http.StatusOK),
            "Message": "Successful dude!",
        }) // set data to body response
        response.Status(r, constant.StatusSuccess) // set status code
        response.WriteJSON(w, r, body) // response to JSON
    }
    
    http.ListenAndServe(":8000", handler)
}
```
