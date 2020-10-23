# stopless
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/i-pva/stopless/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/i-pva/stopless)](https://goreportcard.com/report/github.com/i-pva/stopless)

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
