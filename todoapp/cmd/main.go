package main

import (
	"flag"
	"kuberneteslab/todoapp/pkg/server"
)

func main() {
	configFilePath := flag.String("config", "config/local.json", "path to config file")
	flag.Parse()
	server.Start(*configFilePath)
}
