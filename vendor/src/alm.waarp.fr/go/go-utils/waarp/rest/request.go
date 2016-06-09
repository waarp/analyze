package rest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

// A request that enbeds http.Request but provides authentication mechanisms
// for Waarp R66 REST API. it is mainly designed to be used by RestClient.
type Request struct {
	*http.Request
	username  string
	timestamp string
	data      interface{}
}

// Creates a new request
func NewRequest(method, url string, data interface{}) (*Request, error) {
	bodyContent, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Cannot marshal request to JSON: %s", err.Error())
	}

	hReq, err := http.NewRequest(method, url, bytes.NewBuffer(bodyContent))
	if err != nil {
		return nil, fmt.Errorf("Cannot create request: %s", err.Error())
	}
	return &Request{Request: hReq, data: &data}, nil
}

func (r *Request) SetUser(username string) {
	if username != "" {
		r.username = username
		r.Header.Add("X-Auth-User", username)
	}
}

func (r *Request) SetTimestamp(timestamp string) {
	r.timestamp = timestamp
	r.Header.Add("X-Auth-Timestamp", timestamp)
}

func (r *Request) Sign(password string, key []byte) {
	r.Header.Add("X-Auth-Key", r.computeSignature(password, key))
}

func (r *Request) computeSignature(password string, key []byte) string {
	sigData := map[string]string{}
	keys := []string{}

	sigData["x-auth-user"] = r.username
	sigData["x-auth-timestamp"] = r.timestamp
	keys = append(keys, "x-auth-user", "x-auth-timestamp")

	sort.Strings(keys)

	finalString := ""
	for i, k := range keys {
		sep := "&"
		if i == 0 {
			sep = "?"
		}
		finalString = fmt.Sprintf("%s%s%s=%s", finalString, sep, k, sigData[k])
	}

	finalString = r.URL.Path + finalString + "&X-Auth-InternalKey=" + password

	h := hmac.New(sha256.New, key)
	h.Write([]byte(finalString))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// func extractData(data interface{}) (map[string]string, []string) {
//  sigData := map[string]string{}
//  keys := []string{}

//  v := reflect.ValueOf(data)
//  t := reflect.TypeOf(data)
//  for i := 0; i < t.NumField(); i++ {
//      f := v.Field(i)
//      if f.Interface() != reflect.Zero(f.Type()).Interface() {
//          fieldName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
//          fieldName = strings.ToLower(fieldName)

//          var val string
//          switch f.Type().String() {
//          case "string":
//              val = f.String()
//          case "int":
//              val = fmt.Sprintf("%d", f.Int())
//          case "time.Time":
//              val = f.Interface().(time.Time).Format(time.RFC3339)
//          case "*time.Time":
//              val = f.Interface().(*time.Time).Format(time.RFC3339)
//          case "bool":
//              val = fmt.Sprintf("%t", f.Bool())
//          default:
//              val = f.String()
//          }
//          sigData[fieldName] = strings.ToLower(val)
//          keys = append(keys, fieldName)
//      }
//  }
//  return sigData, keys
// }

// A generic response as sent by Waarp R66. Use NewResponse to get a
// Response object with a deserializer for the specific request
type Response struct {
	Method string            `json:"X-method"` //"X-method":"GET",
	Path   string            //"path":"/transfers",
	Base   string            //"base":"transfers",
	Qs     map[string]string `json:"uri"` //"uri":{},
	Answer struct {
		Filters map[string]interface{} `json:"filter"`
		Results interface{}
		Count   int
		Limit   int
	}
	Command    string            //"command":"MULTIGET",
	Message    string            //"message":"OK",
	StatusCode int               `json:"code"` //"code":200,
	Cookie     map[string]string //"cookie":{}
}

func NewResponse(results interface{}) *Response {
	r := &Response{}
	r.Answer.Results = results
	return r
}
