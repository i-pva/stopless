package stopless

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var notifySignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGKILL,
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGHUP,
}

type Server struct {
	http.Server
	context.Context
	signalHooks         map[os.Signal]func()
	customNotifySignals []os.Signal
	terminationPeriod   time.Duration
}

func (srv *Server) ListenAndServe() error {
	go srv.handleSignals()
	err := srv.Server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	go srv.handleSignals()
	err := srv.Server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func (srv *Server) restart() error {
	return nil
}

func (srv *Server) Notify(sig ...os.Signal) {
	srv.customNotifySignals = append(srv.customNotifySignals, sig...)
}

func (srv *Server) handleSignals() {
	var sig os.Signal
	listener := make(chan os.Signal, 1)

	if len(srv.customNotifySignals) != 0 {
		notifySignals = srv.customNotifySignals
	}

	signal.Notify(
		listener,
		notifySignals...,
	)

	for {
		sig = <-listener

		hook, ok := srv.signalHooks[sig]
		if ok {
			hook()
		}

		switch sig {
		case syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM:
			log.Printf("received %v.", sig)
			err := srv.shutdown()
			if err != nil {
				log.Fatalf("shutdown err: %v\n", err)
			}
			break
		default:
			log.Printf("received %v: nothing i care about...\n", sig)
		}
	}
}

func (srv *Server) shutdown() error {
	if srv.terminationPeriod == 0 {
		srv.terminationPeriod = 30 * time.Second
	}

	if srv.Context == nil {
		srv.Context = context.Background()
	}

	ctx, cancel := context.WithTimeout(srv.Context, srv.terminationPeriod)
	defer cancel()
	err := srv.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %v\n", err)
	}
	return nil
}
