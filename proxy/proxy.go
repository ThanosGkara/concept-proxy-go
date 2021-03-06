package proxy

import (
	"bytes"
	"concept-proxy-go/intercept"
	"concept-proxy-go/lb"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/patrickmn/go-cache"
)

// Log the typeform payload and redirect url
func logRequestPayload(request *http.Request) {

	// dec := requestBodyDecoder(request)
	body, err := ioutil.ReadAll(request.Body)
	if err == nil {
		log.Printf("Request payload: \n%s", string(body))
	}

	// In go lang if you read the request body then any susequent calls
	// are unable to read the body again
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

/*
	Reverse Proxy Logic
*/

// Serve a reverse proxy for a given url
func serveReverseProxy(srv map[string]*lb.RoundRobin, pcache *cache.Cache, res http.ResponseWriter, req *http.Request) {

	if req.Body != nil {
		// requestPayload := parseRequestBody(req)
		logRequestPayload(req)
	}

	fmt.Println("URL before: " + req.Host + req.URL.String() + "\n")

	cache_key := req.Host + req.URL.String()

	//Ask the cache to get the cached value
	cachedResponse, found := pcache.Get(cache_key)

	//Cache Hit
	if found {
		fmt.Printf("Cache hit for: %s", cache_key)
		fmt.Fprintf(res, cachedResponse.(string))

	} else { //Cache miss

		// parse the url
		url_, _ := url.Parse(req.Host)

		var service string
		dom := strings.Split(url_.String(), ":")[0]
		fmt.Println(dom)

		// choose backend with roundrobin
		service = (*srv[dom]).Next()
		fmt.Println("Service: " + service)
		url_.Host = service
		url_.Scheme = "http"

		// create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(url_)

		// Update the headers to allow for SSL redirection
		req.Header.Set("X-Forwarded-Host", dom)
		req.URL.Host = url_.Host
		req.URL.Scheme = url_.Scheme

		fmt.Println("URL after: " + req.URL.String() + "\n")

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(intercept.NewCustomWriter(res, pcache, cache_key), req)

	}

}

// Given a request send it to the appropriate url
func ProxyOperation(srv map[string]*lb.RoundRobin, cache *cache.Cache) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		serveReverseProxy(srv, cache, res, req)
	}
}
