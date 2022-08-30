package main

import (
	"embed"
	"log"
	"os"
	"os/signal"
	"score-calculate/service"
	"syscall"
)

var content embed.FS

func main() {
	service.StartService()

	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	service.ShutDown()
	log.Println("Server exiting")
}
