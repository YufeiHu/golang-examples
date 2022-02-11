package main

import (
	"context"
    "encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ztb-tester/internal/http_server"
)

const (
	basePath = "/v1/demo"
)

type responseStruct struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func writeAndLog(w http.ResponseWriter, req *http.Request, status int, err error) {
	if err != nil {
		fmt.Printf("Http handler: %s %s (%s) -> %d. error: %v\n", req.Method, req.RequestURI, req.RemoteAddr, status, err)
	} else {
		fmt.Printf("Http handler: %s %s (%s) -> %d\n", req.Method, req.RequestURI, req.RemoteAddr, status)
	}
	w.WriteHeader(status)
}

func NewHttpHandler(ctx context.Context) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(basePath + "/", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			_ = req.Body.Close()
		}()

		if req.Method != "GET" {
			writeAndLog(w, req, http.StatusMethodNotAllowed, nil)
			return
		}

		response := responseStruct{
			"1",
			"yufhu",
		}

		body, err := json.Marshal(response)
		if err != nil {
			writeAndLog(w, req, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			writeAndLog(w, req, http.StatusInternalServerError, err)
			return
		}
	})

	return mux
}

func getInterruptableCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	// detect and cancel on SIGINT or SIGTERM
	go func() {
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-c:
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func interruptableServe(parentCtx context.Context, listener net.Listener, server *http.Server) error {
	go func() {
		<-parentCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Failed to shutdown the server: %v\n", err)
		}
	}()

	return server.Serve(listener)
}

func main() {
	ctx, cancel := getInterruptableCtx(context.Background())
	defer cancel()

	handler := http_server.NewHttpHandler(ctx)
	server := http.Server{
		Handler: handler,
	}

	addr := "localhost:4443"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to create a listener on %s\n", addr)
		return
	}

	err = interruptableServe(ctx, listener, &server)
	if err != nil {
		fmt.Printf("Failed to serve on %s\n", addr)
		return
	}
}
