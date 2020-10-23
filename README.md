# stopless

Server has some useful methods for performing graceful 
shutdown.
## Usage
```go
package main

import (
    "net/http"
    "github.com/i-pva/stopless"
)

func main() {

    // Create server
    server := &stopless.Server{
        Server : http.Server{
            Addr:    ":8080",
            Handler: http.NewServeMux(),
        },
    }
	
    // Run server
    err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }

}
```
[License](LICENSE)
-------

This package is released under the MIT license.
