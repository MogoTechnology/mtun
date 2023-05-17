package main

import (
	"github.com/icechen128/mtun/client"
	"github.com/icechen128/mtun/server"
	"go.uber.org/zap"
	"os"
)

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	mode := "server"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "server":
		startServer()
	case "client":
		startClient()
	}
}

func startServer() {
	err := server.New().Start()
	if err != nil {
		zap.L().Fatal("Server start failed", zap.Error(err))
	}
}

func startClient() {
	err := client.StartClient()
	if err != nil {
		zap.L().Fatal("Client start failed", zap.Error(err))
	}
}
