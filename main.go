package main

import (
	"go-proxy/proxy"
	"net/http"
)

/*
	Utilities
*/

/*
	Entry
*/

func main() {
	// Read config file
	var config proxy.ProxyConfig
	proxy.GenerateConfig(&config)

	// pretty.Println(config)
	proxylisten := config.Proxy.Listen.Address + ":" + config.Proxy.Listen.Port

	// start server
	http.HandleFunc("/", proxy.ProxyOperation(&config))
	if err := http.ListenAndServe(proxylisten, nil); err != nil {
		panic(err)
	}
}
