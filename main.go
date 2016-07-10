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

func NewNotificationProxy(to *url.URL) *httputil.ReverseProxy {
  director := func(req *http.Request) {
    query := req.URL.Query()
    messages := query["message"]
    services := query["service"]
    var service = "GENERIC"
    if len(services) > 0 {
      service = services[0]
    }

    var message = fmt.Sprintf("[%s]", service)
    if len(messages) > 0 {
      message = fmt.Sprintf("%s %s", message, messages[0])
    }

    req.URL = to
    req.Host = to.Host
    req.Method = "POST"

    bodyString := fmt.Sprintf("{\"text\":\"%s\"}", message)
    req.Body = ioutil.NopCloser(bytes.NewBufferString(bodyString))
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
