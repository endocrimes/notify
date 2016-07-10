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

func GetValueOrFallbackFromValues(v url.Values, key string, fallback string) string {
  value := v.Get(key)
  if value == "" {
    value = fallback
  }

  return value
}

func RewriteRequest(req *http.Request, to *url.URL) {
    query := req.URL.Query()
    message := GetValueOrFallbackFromValues(query, "message", "...")
    service := GetValueOrFallbackFromValues(query, "service", "GENERIC")
    text := fmt.Sprintf("[%s] %s", service, message)
    json := fmt.Sprintf("{\"text\":\"%s\"}", text)
    body := ioutil.NopCloser(bytes.NewBufferString(json))

    req.URL = to
    req.Host = to.Host
    req.Body = body
    req.Method = "POST"
    req.Header = http.Header{}
}

func NewNotificationProxy(to *url.URL) *httputil.ReverseProxy {
  director := func(req *http.Request) {
    RewriteRequest(req, to)
  }

  return &httputil.ReverseProxy{
    Director: director,
  }
}

func main() {
  flag.Parse()
  to, err := url.Parse(*toHost)
  if err != nil {
    fmt.Println(err)
  }
  proxy := NewNotificationProxy(to)
  fmt.Printf("Proxying %s->%s.\n", *fromHost, *toHost)
  http.ListenAndServe(*fromHost, proxy)
}
