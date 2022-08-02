package urlconfig_test

import (
	"context"
	"errors"
	"io"
	"net/url"
	"testing"

	"github.com/davidsbond/urlconfig"
)

func TestRegistry_Register(t *testing.T) {
	registry := urlconfig.NewRegistry[string]()
	factory := func(ctx context.Context, u *url.URL) (string, error) {
		return u.String(), nil
	}

	t.Run("It should not return an error for a new scheme", func(t *testing.T) {
		if err := registry.Register("example", factory); err != nil {
			t.Fatal("expected no error, but got one")
		}
	})

	t.Run("It should return an error for a duplicate scheme", func(t *testing.T) {
		if !errors.Is(registry.Register("example", factory), urlconfig.ErrSchemeExists) {
			t.Fatal("expected ErrSchemeExists, but didn't get it")
		}
	})
}

func TestRegistry_Configure(t *testing.T) {
	ctx := context.Background()
	registry := urlconfig.NewRegistry[string]()
	factory := func(ctx context.Context, u *url.URL) (string, error) {
		if u.Host == "error" {
			return "", io.EOF
		}

		return u.String(), nil
	}

	if err := registry.Register("example", factory); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	t.Run("It should return the configured value with a valid URL", func(t *testing.T) {
		expected := "example://test"
		actual, err := registry.Configure(ctx, expected)
		if err != nil {
			t.Fatal("expected no error, but got one")
		}

		if expected != actual {
			t.Fatalf("expcted %s, got %s", expected, actual)
		}
	})

	t.Run("It should return an error for a scheme that is not registered", func(t *testing.T) {
		_, err := registry.Configure(ctx, "nope://nope")
		if !errors.Is(err, urlconfig.ErrUnknownScheme) {
			t.Fatal("expected ErrUnknownScheme, but didn't get it")
		}
	})

	t.Run("It should propagate errors from the factory function", func(t *testing.T) {
		_, err := registry.Configure(ctx, "example://error")
		if err == nil {
			t.Fatal("expected an error, but didn't get one")
		}
	})
}
