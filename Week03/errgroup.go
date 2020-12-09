package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	group, egctx := errgroup.WithContext(context.Background())
	var servErr = make(chan error)
	group.Go(func() error {
		handler := http.NewServeMux()
		handler.HandleFunc("/err", func(writer http.ResponseWriter, request *http.Request) {
			servErr <- errors.New("http-req err")
		})
		httpServ := &http.Server{
			Addr:    ":8000",
			Handler: handler,
		}
		go func(s *http.Server) {
			servErr <- httpServ.ListenAndServe()
		}(httpServ)
		select {
		case <-servErr:
			return errors.New("http server error")
		case <-egctx.Done():
			fmt.Println("close the http-server")
			return httpServ.Close()
		}
	})
	group.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quit:
			return errors.New("syscall signaling")
		case <-egctx.Done():
			return errors.New("")
		}
	})
	err := group.Wait()
	fmt.Printf("app is over by reason[%v]\n", err)
}
