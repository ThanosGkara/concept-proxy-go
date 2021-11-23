package main

import (
	"concept-proxy-go/lb"
	"concept-proxy-go/proxy"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v2"
)

/*
	Structs
*/

type ProxyConfig struct {
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

// Generate config struct
func generateConfig(c *ProxyConfig) {
	filename, _ := filepath.Abs("config.yml")
	yamlConfigFile, err := ioutil.ReadFile(filename)

	fmt.Printf("Server config!! \n%s", yamlConfigFile)

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
	var config ProxyConfig
	generateConfig(&config)

	// pretty.Println(config)
	proxylisten := config.Proxy.Listen.Address + ":" + config.Proxy.Listen.Port

	srv := make(map[string]*lb.RoundRobin, len(config.Proxy.Services))

	// populate Services with each backends to be used
	// with the LB
	for _, s := range config.Proxy.Services {
		hosts := make([]string, len(s.Hosts))
		for i, h := range s.Hosts {
			hosts[i] = h.Address + ":" + h.Port
		}
		srv_tmp, err := lb.New(hosts...)
		if srv_tmp != nil {
			srv[s.Domain] = &srv_tmp
		} else {
			fmt.Println(err)
		}
	}

	// Create an intitia memory cache
	var pageCache = cache.New(5*time.Minute, 10*time.Minute)

	// start server
	http.HandleFunc("/", proxy.ProxyOperation(srv, pageCache))
	if err := http.ListenAndServe(proxylisten, nil); err != nil {
		panic(err)
	}
}
