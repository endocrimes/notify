package main

import (
  "flag"
  "fmt"
  "net/http"
  "net/http/httputil"
  "net/url"
  "bytes"
  "io/ioutil"
)

var fromHost = flag.String("from", "localhost:80", "The proxy server's host.")
var toHost = flag.String("to", "localhost:8000", "The host that the proxy " +
             " server should forward requests to.")

func main() {
  flag.Parse()
  to, err := url.Parse(*toHost)
  if err != nil {
    fmt.Println(err)
  }
  handler := &RequestHandler{
    Target: to,
  }
  proxy := NewNotificationProxy(handler)
  fmt.Printf("Proxying %s->%s.\n", *fromHost, *toHost)
  http.ListenAndServe(*fromHost, proxy)
}

// @!group Helpers

func NewNotificationProxy(handler *RequestHandler) *httputil.ReverseProxy {
  return &httputil.ReverseProxy{
    Director: handler.handle,
  }
}

func GetValueOrFallbackFromValues(v url.Values, key string, fallback string) string {
  value := v.Get(key)
  if value == "" {
    value = fallback
  }

  return value
}

// @!group Request Handler

type RequestHandler struct {
  Target *url.URL
}

func (r *RequestHandler)handle(req *http.Request) {
    query := req.URL.Query()
    message := GetValueOrFallbackFromValues(query, "message", "...")
    service := GetValueOrFallbackFromValues(query, "service", "GENERIC")
    text := fmt.Sprintf("[%s] %s", service, message)
    json := fmt.Sprintf("{\"text\":\"%s\"}", text)
    body := ioutil.NopCloser(bytes.NewBufferString(json))

    req.URL = r.Target
    req.Host = r.Target.Host
    req.Body = body
    req.Method = "POST"
    req.Header = http.Header{}
}
