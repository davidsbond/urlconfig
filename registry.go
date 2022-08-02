// Package urlconfig provides a generic registry for configuration of values based on URLs. Individual factory functions
// are added to the registry which will be run when checked against registered URL schemes in the registry.
package urlconfig

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

type (
	// The Registry type contains Factory implementations that are used to configure values based on a URL.
	Registry[T any] struct {
		mux       *sync.RWMutex
		factories map[string]Factory[T]
	}

	// The Factory type describes a function that, given a valid URL, returns an instance of the parameterized type
	// T.
	Factory[T any] func(ctx context.Context, u *url.URL) (T, error)
)

var (
	// ErrSchemeExists is the error given when calling Registry.Register with a scheme that is already registered.
	ErrSchemeExists = errors.New("already registered")
	// ErrUnknownScheme is the error given when calling Registry.Configure with a scheme that is not registered.
	ErrUnknownScheme = errors.New("unknown scheme")
)

// NewRegistry returns a new instance of the Registry type that will store Factory implementations that configure values
// of the parameterized type T based on URLs. The scheme of the URL determines which Factory implementation is called.
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		mux:       &sync.RWMutex{},
		factories: make(map[string]Factory[T]),
	}
}

// Register a new Factory implementation for a URL scheme. If the scheme is already registered, returns ErrSchemeExists.
func (r *Registry[T]) Register(scheme string, factory Factory[T]) error {
	r.mux.RLock()
	if _, ok := r.factories[scheme]; ok {
		r.mux.RUnlock()
		return fmt.Errorf("%s: %w", scheme, ErrSchemeExists)
	}
	r.mux.RUnlock()

	r.mux.Lock()
	r.factories[scheme] = factory
	r.mux.Unlock()

	return nil
}

// Configure an instance of T based on the provided URL. If a Factory does not exist for the scheme, returns ErrUnknownScheme.
// is returned.
func (r *Registry[T]) Configure(ctx context.Context, urlStr string) (T, error) {
	var out T

	u, err := url.Parse(urlStr)
	if err != nil {
		return out, err
	}

	r.mux.RLock()
	defer r.mux.RUnlock()
	factory, ok := r.factories[u.Scheme]
	if !ok {
		return out, fmt.Errorf("%w: %s", ErrUnknownScheme, u.Scheme)
	}

	out, err = factory(ctx, u)
	return out, err
}
