# go-cors

Package **cors** provides tools for working with **Cross-Origin Resource Sharing** (**CORS**), for the Go programming language.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-cors

[![GoDoc](https://godoc.org/github.com/reiver/go-cors?status.svg)](https://godoc.org/github.com/reiver/go-cors)

## Examples

Here is an example using `cors.ProxyHandler`:

```golang
import (
	"os"

	"github.com/reiver/go-cors"
)

// ...

var corsProxyHandler cors.ProxyHandler = cors.ProxyHandler{
	LogWriter: os.Stdout, // write logs to os.Stdout
}

// ...

err := http.ListenAndServe(addr, &corsProxyHandler)
```

## Import

To import package **cors** use `import` code like the following:
```
import "github.com/reiver/go-cors"
```

## Installation

To install package **cors** do the following:
```
GOPROXY=direct go get github.com/reiver/go-cors
```

## Author

Package **cors** was written by [Charles Iliya Krempeaux](http://reiver.link)
