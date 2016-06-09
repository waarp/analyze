package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (r RestClient) url(uri string) string {
	return fmt.Sprintf("%s%s", r.BaseUrl, uri)
}

// Runs the request req and unmarshal the response in the target object.
// It returns an error if the server returns a sttus code in the 4xx or 5xx
// ranges.
func (r *RestClient) exec(req *Request, target interface{}) error {
	req.SetUser(r.User)
	if r.NeedsTimestamp {
		req.SetTimestamp(time.Now().Format(time.RFC3339))
	}

	if len(r.signingKey) != 0 {
		req.Sign(r.Password, r.signingKey)
	}

	resp, err := r.http.Do(req.Request)
	if err2 := processErrors(resp, err); err2 != nil {
		return err2
	}

	rv := NewResponse(target)
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(rv); err != nil {
		return fmt.Errorf("Cannot unmarshall response: %s", err.Error())
	}
	return nil
}

func processErrors(resp *http.Response, err error) error {
	if err != nil {
		switch {

		case strings.Contains(err.Error(), "x509"):
			return fmt.Errorf("Certificate validation error: %s", err.Error())

		case strings.Contains(err.Error(), "net/http: timeout"):
			return fmt.Errorf("A connection has been made, but the server did not respond")

		case strings.Contains(err.Error(), "connection refused"):
			return fmt.Errorf("Cannot open a connection to the server (%s)",
				err.Error())

		default:
			return fmt.Errorf(`An error occured while executing the request: %s`, err.Error())
		}
	}

	switch resp.StatusCode {

	case http.StatusBadRequest:
		return fmt.Errorf("Server returned 'Bad Request'")

	case http.StatusUnauthorized:
		return fmt.Errorf("Bad authentication (username, timestamp, signature)")

	case http.StatusForbidden:
		return fmt.Errorf("Request '%s %s' forbidden",
			resp.Request.Method, resp.Request.URL)

	case http.StatusNotFound:
		return fmt.Errorf("url %s not found", resp.Request.URL)

	case http.StatusNotImplemented:
		return fmt.Errorf("url %s not found or method %s invalid for this url",
			resp.Request.URL, resp.Request.Method)

	case http.StatusInternalServerError:
		return fmt.Errorf("Internal server error")

	case http.StatusOK:
		// noop

	default:
		return fmt.Errorf("Unhandled status code : %d\n", resp.StatusCode)
	}
	return nil
}
