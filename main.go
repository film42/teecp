package main

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	Bind  string   `json:"bind"`
	Proxy string   `json:"proxy"`
	Tees  []string `json:"tees"`
}

func main() {
	InitLogger()

	configPathPtr := flag.String("config", "config.json", "Path to the teecp config.")
	debugPtr := flag.Bool("debug", false, "Enable debug logging.")
	flag.Parse()

	if !(*debugPtr) {
		DisableDebugLogging()
	}

	configFile, err := os.Open(*configPathPtr)
	if err != nil {
		logger.Fatal.Println("Error opening config file:", *configPathPtr)
		return
	}

	config := new(Config)
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Fatal.Println("Error parsing config file:", err)
		return
	}

	logger.Info.Printf("Adding Proxy:\tclient -> %s\t%s -> client\n", config.Proxy, config.Proxy)
	for _, server := range config.Tees {
		logger.Info.Printf("Adding Tee:\tclient -> %s\t%s -> sink\n", server, server)
	}

	teecp := NewTeecp(config)
	err = teecp.ListenAndServe(config)
	if err != nil {
		logger.Fatal.Println("Error starting server:", err)
	}
}
