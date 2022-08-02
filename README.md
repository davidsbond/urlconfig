# urlconfig

This package aims to provide an abstraction for package developers that want to configure things using URLs. It provides
a generic `Registry` type that will invoke factory functions based on a URL scheme.

As this package uses parameterized types you must use go 1.18 or later. 

## Example usage

```go
package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/davidsbond/urlconfig"
)

type (
	Blob interface {
		io.ReadWriter
	}
)

func main() {
	ctx := context.Background()
	registry := urlconfig.NewRegistry[Blob]()

	// If you want an in-memory blob, just return a bytes.Buffer.
	err := registry.Register("memory", func(ctx context.Context, u *url.URL) (Blob, error) {
		return bytes.NewBuffer([]byte{}), nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	// If you want to persist to disk, use a file.
	err = registry.Register("file", func(ctx context.Context, u *url.URL) (Blob, error) {
		return os.Create(u.Path)
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Get ourselves an in-memory blob writer.
	_, err = registry.Configure(ctx, "memory://blob")
	if err != nil {
		log.Fatalln(err)
	}

	// Get ourselves a file-backed blob writer.
	_, err = registry.Configure(ctx, "file:///path/to/file")
	if err != nil {
		log.Fatalln(err)
	}
}
```
