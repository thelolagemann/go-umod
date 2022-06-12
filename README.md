# go-umod

> An unofficial umod API client, written in golang. 

## Getting started

### Installation
```shell
go get -u github.com/thelolagemann/go-umod
```

### Example Usage
```go
package main

import (
	"fmt"
	"github.com/thelolagemann/go-umod"
)

func main() {
	results, err := umod.Search("helicopter")
	if err != nil {
		// handle err
	}
	printResults(results)
	
	if results.NextPageURL != "" {
		next, err := results.NextPage()
		if err != nil {
			// handle err
    }
  	printResults(next)
  }
}

func printResults(results *umod.SearchResponse) {
	for _, r := range results.Data {
		fmt.Println(r)
  }
}
```