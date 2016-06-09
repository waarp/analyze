package rest

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A client for a single instance of Waarp R66 Server
type RestClient struct {
	// The base URL to connect to the rest server (i.e.: https://localhost:8088)
	BaseUrl string

	// The credentials to connect to the server
	User, Password string

	// Whether the server need a timestamp in requests or not
	NeedsTimestamp bool

	// the key used to sign requests
	signingKey []byte

	// An setup http client to use to execute http requests
	http http.Client

	// the transport used by the client
	transport *http.Transport

	// Services accessor
	Transfers *TransferServices
}

// Creates a new R66 REST Client. connectionString has the following syntax :
//
//   http[s]://username:password@host:port[/options]
//
// Following options are accepted:
//
//   timeout
//     timeout delay for http requests in seconds (default: 10)
func NewRestClient(connectionString string) (*RestClient, error) {
	url, err := url.Parse(connectionString)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string '%s'", connectionString)
	}
	if !(url.Scheme == "https" || url.Scheme == "http") {
		return nil, fmt.Errorf("'%s' is not a supported scheme", url.Scheme)
	}

	transport := &http.Transport{
		ResponseHeaderTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			ServerName:         strings.Split(url.Host, ":")[0],
			InsecureSkipVerify: true,
		},
	}

	tmpTimeout := url.Query().Get("timeout")
	if tm, err := strconv.ParseInt(tmpTimeout, 10, 64); err == nil {
		transport.ResponseHeaderTimeout = time.Duration(tm) * time.Second
	}

	rest := &RestClient{
		BaseUrl:   fmt.Sprintf("%s://%s", url.Scheme, url.Host),
		transport: transport,
		http: http.Client{
			Transport: transport,
			Timeout:   transport.ResponseHeaderTimeout,
		},
	}

	if url.User != nil {
		rest.User = url.User.Username()
		if pw, ok := url.User.Password(); ok {
			rest.Password = pw
		}
	}
	rest.Transfers = &TransferServices{client: rest}
	return rest, nil
}

// Reads the signing key from the given file
func (r *RestClient) SetSigningKey(filename string) error {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Cannot read the signing key file %s", filename)
	}
	r.signingKey = key
	r.NeedsTimestamp = true
	return nil
}

func (r *RestClient) SetRootCA(filename string) error {
	caPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Could not load root certificates!")
	}
	caPool.AppendCertsFromPEM(caCert)
	r.transport.TLSClientConfig.RootCAs = caPool
	return nil
}
