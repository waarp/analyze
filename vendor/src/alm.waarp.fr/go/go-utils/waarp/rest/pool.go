package rest

import (
	"fmt"
	"net/url"
	"time"
)

// Represents a pool of clients, one for each scheme and host.
// It acts like a cache
var Pool = pool{clients: map[string]*RestClient{}}

// internal type for the cache
type pool struct {
	clients  map[string]*RestClient
	defaults DefaultConfig
}

// Gets a client given its uri connectionString. connectionString has the same syntax
// as the one given to NewRestClient.
//
// If the client already exists for the provided scheme and host, the existing
// client is returned. It is not updated with the parameters given in connectionString.
//
// If no client exist, a new one is initialized with
// initialized with the given connectionString. If a default config has been set,
// parameters missing from the connectionString are filled with default values.
// the connectionString has the highest precedence.
func (p *pool) Get(connectionString string) (*RestClient, error) {
	key, err := makeKey(connectionString)
	if err != nil {
		return nil, err
	}

	if c, ok := p.clients[key]; ok {
		return c, nil
	}

	c, err := NewRestClient(connectionString)
	if err != nil {
		return nil, err
	}

	if err := p.defaults.apply(c); err != nil {
		return nil, err
	}

	p.clients[key] = c
	return c, nil
}

// Deletes a client from the pool and destroys it. Only the scheme and
// host parts of the connectionString are used to identify the cliet
func (p *pool) Del(connectionString string) {
	key, err := makeKey(connectionString)
	if err != nil {
		return
	}
	delete(p.clients, key)
}

// Reset the default configuration for new clients initialized from the pool.
// It does not update the existing clients and only affects new ones.
func (p *pool) SetDefaultConfig(dc DefaultConfig) {
	p.defaults = dc
}

func makeKey(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("Invalid connection string '%s'", uri)
	}
	return fmt.Sprintf("%s://%s", url.Scheme, url.Host), nil
}

// The default configuration used to initialize new clients from the pool
type DefaultConfig struct {
	User       string
	Password   string
	Timeout    time.Duration
	SigningKey string
	RootCA     string
}

// Apply the default configuration to a given rest client.
func (d *DefaultConfig) apply(c *RestClient) error {
	if c.User == "" && d.User != "" {
		c.User = d.User
	}

	if c.Password == "" && d.Password != "" {
		c.Password = d.Password
	}

	if c.User == "" && d.User != "" {
		c.User = d.User
	}

	if c.transport.ResponseHeaderTimeout == 10*time.Second && d.Timeout != 0 {
		c.transport.ResponseHeaderTimeout = d.Timeout
		c.http.Timeout = c.transport.ResponseHeaderTimeout
	}

	if len(c.signingKey) == 0 && d.SigningKey != "" {
		if err := c.SetSigningKey(d.SigningKey); err != nil {
			return err
		}
	}

	if c.transport.TLSClientConfig.RootCAs == nil && d.RootCA != "" {
		if err := c.SetRootCA(d.RootCA); err != nil {
			return err
		}
	}
	return nil
}
