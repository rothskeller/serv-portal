//go:build windows

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sys/windows/svc"
	"sunnyvaleserv.org/portal/server"
)

type service struct{}

func (s service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptShutdown | svc.AcceptStop
	var (
		listener net.Listener
		ws       http.Server
		err      error
	)
	status <- svc.Status{State: svc.StartPending}
	log.Print("SERVER START")
	if listener, err = net.Listen("tcp", "localhost:7190"); err != nil {
		log.Fatalf("net.Listen: %s", err)
	}
	ws.Handler = server.Server
	go ws.Serve(listener)
	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	for c := range r {
		switch c.Cmd {
		case svc.Interrogate:
			status <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			goto STOP
		default:
			log.Printf("Unexpected service control request #%d", c)
		}
	}
STOP:
	status <- svc.Status{State: svc.StopPending}
	if err = ws.Shutdown(context.Background()); err != nil {
		log.Printf("server.Shutdown: %s", err)
	}
	log.Print("SERVER STOP")
	return false, 1
}

func main() {
	os.Chdir(`C:\SERV`)
	if err := svc.Run("serv-portal", service{}); err != nil {
		log.Fatalf("svc.Run: %s", err)
	}
}
