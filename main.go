package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	// "github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

/*
	Structs
*/

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

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
func logRequestPayload(requestionPayload requestPayloadStruct) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition)
}

// Log the env variables required for a reverse proxy
func logSetup(serverAddress string) {
	a_condtion_url := os.Getenv("A_CONDITION_URL")
	b_condtion_url := os.Getenv("B_CONDITION_URL")
	default_condtion_url := os.Getenv("DEFAULT_CONDITION_URL")

	log.Printf("Server will run on: %s\n", serverAddress)
	log.Printf("Redirecting to A url: %s\n", a_condtion_url)
	log.Printf("Redirecting to B url: %s\n", b_condtion_url)
	log.Printf("Redirecting to Default url: %s\n", default_condtion_url)
}

/*
	Reverse Proxy Logic
*/

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) *json.Decoder {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

// Parse the requests body
func parseRequestBody(request *http.Request) requestPayloadStruct {
	decoder := requestBodyDecoder(request)

	var requestPayload requestPayloadStruct
	err := decoder.Decode(&requestPayload)

	if err != nil {
		panic(err)
	}

	return requestPayload
}

// Serve a reverse proxy for a given url
func serveReverseProxy(c *proxyConfig, res http.ResponseWriter, req *http.Request) {

	if req.Body != nil {
		requestPayload := parseRequestBody(req)

		logRequestPayload(requestPayload)
	}

	// parse the url
	url, _ := url.Parse(req.Host)

	var service string
	rand.Seed(time.Now().UnixNano())
	for _, serv := range c.Proxy.Services {
		if url.Host == serv.Domain {
			backend := rand.Intn(len(serv.Hosts))
			addr := serv.Hosts[backend]
			service = addr.Address + addr.Port
		}
	}
	url.Host = service

	// create the reverse proxy
	// proxy := httputil.NewSingleHostReverseProxy(service)
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

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

	// Log setup values
	logSetup(proxylisten)

	// start server
	http.HandleFunc("/", proxyOperation(&config))
	if err := http.ListenAndServe(proxylisten, nil); err != nil {
		panic(err)
	}
}
