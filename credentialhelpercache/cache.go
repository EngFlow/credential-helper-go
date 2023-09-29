// Copyright 2023 EngFlow, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package credentialhelpercache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"

	credentialhelper "github.com/EngFlow/credential-helper-go"
)

const (
	// DefaultCacheDuration specifies the default TTL for cache entries.
	DefaultCacheDuration = 30 * time.Minute
)

// Options represents options for the Credential Helper cache.
type Options struct {
	// TTL specifies the time to cache credentials before invoking
	// the Credential Helper again.
	//
	// If not set, TTL defaults to `DefaultCacheDuration`.
	TTL time.Duration
}

// New wraps a `CredentialHelper` with caching.
func New(delegate credentialhelper.CredentialHelper, options Options) (CachingCredentialHelper, error) {
	ttl := options.TTL
	if ttl < 0 {
		return nil, fmt.Errorf("ttl must not be negative, got %v", ttl)
	} else if ttl == 0 {
		ttl = DefaultCacheDuration
	}

	cache := ttlcache.New[credentialhelper.GetCredentialsRequest, credentialhelper.GetCredentialsResponse](
		ttlcache.WithTTL[credentialhelper.GetCredentialsRequest, credentialhelper.GetCredentialsResponse](ttl),
		ttlcache.WithDisableTouchOnHit[credentialhelper.GetCredentialsRequest, credentialhelper.GetCredentialsResponse]())
	go cache.Start()

	c := &cachingCredentialHelper{
		delegate: delegate,

		cache: cache,
	}
	return c, nil
}

// CachingCredentialHelper represents a `CredentialHelper` that
// internally caches credentials.
//
// Use `New()` to create an instance.
type CachingCredentialHelper interface {
	credentialhelper.CredentialHelper

	// Close closes the Credential Helper and releases all associated resources.
	Close() error
}

type cachingCredentialHelper struct {
	credentialhelper.CredentialHelperBase

	delegate credentialhelper.CredentialHelper

	cacheMutex sync.Mutex
	closed     bool
	cache      *ttlcache.Cache[credentialhelper.GetCredentialsRequest, credentialhelper.GetCredentialsResponse]
}

func (c *cachingCredentialHelper) GetCredentials(ctx context.Context, request *credentialhelper.GetCredentialsRequest, extraParameters ...string) (*credentialhelper.GetCredentialsResponse, error) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.closed {
		return nil, errors.New("Cannot get credentials from closed Credential Helper")
	}

	if entry := c.cache.Get(*request); entry != nil {
		response := entry.Value()
		return &response, nil
	}

	response, err := c.delegate.GetCredentials(ctx, request, extraParameters...)
	if err != nil {
		return nil, err
	}

	// TTL of 0 indicates to use the TTL specified when creating the cache.
	c.cache.Set(*request, *response /* ttl= */, 0)

	return response, nil
}

func (c *cachingCredentialHelper) Close() error {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.closed {
		// Already closed.
		return nil
	}
	c.closed = true

	c.cache.Stop()

	return nil
}
