package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"syscall"
)

const (
	APP_SOCK_PATH = "/var/run/go.sock"
	APP_UID       = 1000
	APP_GID       = 1001
)

func main() {
	fmt.Println("Here we GO!")
	// HTTP Server
	sockPath := APP_SOCK_PATH
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		log.Fatalf("[FATAL] [main] net.Listen(). err=%s", err)
	}
	if err := os.Chown(sockPath, APP_UID, APP_GID); err != nil {
		log.Fatalf("[FATAL] [main] os.Chown(). err=%s", err)
	}
	if err := os.Chmod(sockPath, 0664); err != nil {
		log.Fatalf("[FATAL] [main] os.Chmod(). err=%s", err)
	}
	go func() {
		log.Printf("[INFO] [main] Start server. sock=%s", sockPath)
		mux := http.NewServeMux()
		mux.HandleFunc("/", hello)
		if err := fcgi.Serve(listener, mux); err != nil {
			log.Fatalf("[FATAL] [main] fcgi.Serve(). err=%s", err)
		}
	}()

	shutdown(listener)
}

func shutdown(listener net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	listener.Close()
	log.Fatalf("[FATAL] [main] Caught a signal. signal=%s", s)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

/*
import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
)

type FastCGIServer struct{}

func (s FastCGIServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Callback called!")
	w.Write([]byte("This is a FastCGI example server.\n"))
}

func main() {
	fmt.Println("Starting server...")
	l, _ := net.Listen("tcp", "127.0.0.1:9001")
	b := new(FastCGIServer)
	fcgi.Serve(l, b)
}
*/
