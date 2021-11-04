package main

import (
	"bytes"
	// "encoding/json"
	// "os"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	// "github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

/*
	Structs
*/

type proxyConfig struct {
	Proxy struct {
		Listen struct {
			Address string `yaml:"address"`
			Port    string `yaml:"port"`
		} `yaml:"listen"`
		Services []struct {
			Name   string `yaml:"name"`
			Domain string `yaml:"domain"`
			Hosts  []struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"hosts"`
		} `yaml:"services"`
	} `yaml:"proxy"`
}

/*
	Utilities
*/

// Get env var or default
// func getEnv(key, fallback string) string {
// 	if value, ok := os.LookupEnv(key); ok {
// 		return value
// 	}
// 	return fallback
// }

/*
	Getters
*/

/*
	Logging
*/

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
func serveReverseProxy(c *proxyConfig, res http.ResponseWriter, req *http.Request) {

	if req.Body != nil {
		// requestPayload := parseRequestBody(req)
		logRequestPayload(req)
	}

	// parse the url
	url_, _ := url.Parse(req.Host)

	var service string
	rand.Seed(time.Now().UnixNano())
	dom := strings.Split(url_.String(), ":")[0]
	fmt.Println(dom[0])
	for _, serv := range c.Proxy.Services {
		if dom == serv.Domain {
			backend := rand.Intn(len(serv.Hosts))
			fmt.Println("Forward to service " + serv.Name)
			service = serv.Hosts[backend].Address + ":" + serv.Hosts[backend].Port
			break
		}
	}
	fmt.Println("Service: " + service)
	url_.Host = service
	url_.Scheme = "http"
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url_)

	// Update the headers to allow for SSL redirection
	req.Header.Set("X-Forwarded-Host", dom)
	req.URL.Host = url_.Host
	req.URL.Scheme = url_.Scheme
	req.Host = url_.Host

	fmt.Println("URL after: " + req.URL.String() + "\n")
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func proxyOperation(c *proxyConfig) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		serveReverseProxy(c, res, req)
	}
}

// Generate config struct
func generateConfig(c *proxyConfig) {
	filename, _ := filepath.Abs("./config.yml")
	yamlConfigFile, err := ioutil.ReadFile(filename)

	fmt.Printf("Server config: \n%s", yamlConfigFile)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlConfigFile, c)
	if err != nil {
		panic(err)
	}
}

/*
	Entry
*/

func main() {
	// Read config file
	var config proxyConfig
	generateConfig(&config)

	// pretty.Println(config)
	proxylisten := config.Proxy.Listen.Address + ":" + config.Proxy.Listen.Port

	// start server
	http.HandleFunc("/", proxyOperation(&config))
	if err := http.ListenAndServe(proxylisten, nil); err != nil {
		panic(err)
	}
}
