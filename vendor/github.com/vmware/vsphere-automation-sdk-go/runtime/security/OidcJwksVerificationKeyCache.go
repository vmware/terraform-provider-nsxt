/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// OidcJwksVerificationKeyCache cache for JWKS keys implementation of VerificationKeyCache
type OidcJwksVerificationKeyCache struct {
	useHttps                   bool
	oidcPath                   string
	host                       string
	minTimeBetweenRefreshCalls time.Duration
	client                     *http.Client

	jwksUri string
	issuer  string

	lastRefreshTime time.Time
	lock            sync.Mutex
	// holds jwt verification keys of type []interface{}
	atomicKeyAccessor atomic.Value
	// represents the logical time when the last refresh call started. Increment before each
	// call. Starting times range between [1, maxGeneration]. When maxGeneration is reached,
	// the generationCounter is reset to 1
	generationCounter int
}

const maxGeneration = 1_000_000

// NewOidcJwksVerificationKeyCache creates a VerificationKeyCache which uses the OIDC discovery and JWKS endpoints of a
// given vCenter. The created cache trusts only the issuer of the given vCenter. Other VCs in ELM/HLM are not trusted.
// The cache by default is created with the following config:
// - disabled https
// - oidc path - "/openidconnect/.well-known/openid-configuration"
// - host - "localhost:1080"
// - minimum time between refresh calls set to 1 second
// - pointer to the default http.Client
// Any of these default properties can be configured by providing appropriate CacheConfigOptions
func NewOidcJwksVerificationKeyCache(configOptions ...CacheConfigOption) (*OidcJwksVerificationKeyCache, error) {
	cache := OidcJwksVerificationKeyCache{
		useHttps:                   false,
		oidcPath:                   "/openidconnect/.well-known/openid-configuration",
		host:                       "localhost:1080",
		minTimeBetweenRefreshCalls: 1 * time.Second,
		client:                     &http.Client{},
		generationCounter:          1,
	}

	for _, fn := range configOptions {
		err := fn(&cache)
		if err != nil {
			return nil, err
		}
	}

	return &cache, nil
}

func (cache *OidcJwksVerificationKeyCache) init() error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.atomicKeyAccessor.Load() != nil {
		return nil
	}

	oidcDoc, err := cache.fetchOidcDoc()
	if err != nil {
		return fmt.Errorf("failed to instantiate OidcJwksVerificationKeyCache - %w", err)
	}
	cache.jwksUri = oidcDoc.JwksUri
	cache.issuer = oidcDoc.Issuer
	keys, err := cache.fetchJwks()
	if err != nil {
		return fmt.Errorf("failed to instantiate OidcJwksVerificationKeyCache - %w", err)
	}
	cache.atomicKeyAccessor.Store(keys)
	return nil
}

func (cache *OidcJwksVerificationKeyCache) checkInit() error {
	if cache.atomicKeyAccessor.Load() == nil {
		return cache.init()
	}
	return nil
}

func (cache *OidcJwksVerificationKeyCache) Get(issuer string) ([]interface{}, error) {
	err := cache.checkInit()
	if err != nil {
		return nil, err
	}
	if cache.issuer != issuer {
		return nil, fmt.Errorf("unknown issuer - %v", issuer)
	}
	return cache.atomicKeyAccessor.Load().([]interface{}), nil
}

func (cache *OidcJwksVerificationKeyCache) Refresh(issuer string) error {
	// loaded prior init in order to avoid an immediate rpc after lazy init
	keys := cache.atomicKeyAccessor.Load()
	err := cache.checkInit()
	if err != nil {
		return err
	}
	if cache.issuer != issuer {
		return fmt.Errorf("unknown issuer - %v", issuer)
	}

	generation := cache.generationCounter
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if generation != cache.generationCounter {
		// RPC started after current Refresh() invocation; skip refresh call
		return nil
	} else if !(reflect.DeepEqual(keys, cache.atomicKeyAccessor.Load())) {
		// RPC started before we entered refresh(); reuse result
		return nil
	}

	cache.delayIfNeeded()
	// record logical time before RPC start
	cache.generationCounter = (cache.generationCounter % maxGeneration) + 1

	newKeys, err := cache.fetchJwks()
	if err != nil {
		return err
	}
	cache.atomicKeyAccessor.Store(newKeys)
	return nil
}

func (cache *OidcJwksVerificationKeyCache) fetchOidcDoc() (*OidcDiscoveryDoc, error) {
	oidcUrl := url.URL{Scheme: cache.getScheme(),
		Host: cache.host,
		Path: cache.oidcPath}
	resp, err := cache.httpGet(oidcUrl.String())
	if err != nil {
		return nil, err
	}
	return ParseOidcConfig(resp)
}

func (cache *OidcJwksVerificationKeyCache) fetchJwks() ([]interface{}, error) {
	resp, err := cache.httpGet(cache.jwksUri)
	cache.lastRefreshTime = time.Now()
	if err != nil {
		return nil, err
	}
	return ParseJwks(resp)
}

func (cache *OidcJwksVerificationKeyCache) httpGet(uri string) ([]byte, error) {
	getReq, err := cache.prepareRequest(uri)
	if err != nil {
		return nil, err
	}
	resp, err := cache.client.Do(getReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// TODO utilize io.ReadAll when go is upgraded to 1.16
	parsedResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET to %v failed with %v - %s", uri, resp.StatusCode, string(parsedResp))
	}
	return parsedResp, nil
}

func (cache *OidcJwksVerificationKeyCache) prepareRequest(url string) (*http.Request, error) {
	getReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	getReq.URL.Scheme = cache.getScheme()
	getReq.URL.Host = cache.host
	return getReq, err
}

func (cache *OidcJwksVerificationKeyCache) getScheme() string {
	if cache.useHttps {
		return "https"
	}
	return "http"
}

func (cache *OidcJwksVerificationKeyCache) delayIfNeeded() {
	refreshIssuedAt := time.Now()
	delay := refreshIssuedAt.Sub(cache.lastRefreshTime)
	if delay <= 0 {
		time.Sleep(cache.minTimeBetweenRefreshCalls)
	} else if delay < cache.minTimeBetweenRefreshCalls {
		time.Sleep(cache.minTimeBetweenRefreshCalls - delay)
	}
}

type CacheConfigOption func(*OidcJwksVerificationKeyCache) error

// WithHost specifies the host (and possibly port) where the OIDC/JWKS endpoints reside.
// Value should be host or host:port. Must not be empty
func WithHost(host string) CacheConfigOption {
	return func(cache *OidcJwksVerificationKeyCache) error {
		if host == "" {
			return fmt.Errorf("host value must not be empty")
		}
		cache.host = host
		return nil
	}
}

// UseHttps specifies https or http scheme usage
func UseHttps(useHttps bool) CacheConfigOption {
	return func(cache *OidcJwksVerificationKeyCache) error {
		cache.useHttps = useHttps
		return nil
	}
}

// WithMinTimeBetweenRefreshCalls specifies minimum time between requests made by the OidcJwksVerificationKeyCache
// to the JWKS endpoint. Requests issued in shorter timespans than the given minimum time are delayed until the
// time constraint is satisfied. Must not be negative
func WithMinTimeBetweenRefreshCalls(minTimeBetweenRefreshCalls time.Duration) CacheConfigOption {
	return func(cache *OidcJwksVerificationKeyCache) error {
		if minTimeBetweenRefreshCalls < 0 {
			return fmt.Errorf("minimum time between refresh calls must not be a negative duration")
		}
		cache.minTimeBetweenRefreshCalls = minTimeBetweenRefreshCalls
		return nil
	}
}

// WithClient specifies the client used for communication to the OIDC and JWKS endpoints
func WithClient(client *http.Client) CacheConfigOption {
	return func(cache *OidcJwksVerificationKeyCache) error {
		if client == nil {
			return fmt.Errorf("cache client must not be nil")
		}
		cache.client = client
		return nil
	}
}

// WithOidcPath specifies the OIDC discovery path (appended to the host address)
func WithOidcPath(oidcPath string) CacheConfigOption {
	return func(cache *OidcJwksVerificationKeyCache) error {
		cache.oidcPath = oidcPath
		return nil
	}
}
