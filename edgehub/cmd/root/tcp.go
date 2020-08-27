package root

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var addr = flag.String("port", ":43211", "http service address")

func Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	hub := NewHub()
	go hub.Run(ctx)
	go Serve(hub)
	flag.Parse()
	//Mux holds the map that server looks up from pattern to handler
	router := http.NewServeMux()
	//hook the handler,do something before or after it.
	router.Handle("/mid", Middleware(http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		Servews(ctx, hub, writer, request)
	}))))
	router.HandleFunc("/hub", func(w http.ResponseWriter, r *http.Request) { Servews(ctx, hub, w, r) })
	server := http.Server{Addr: fmt.Sprintf(":%s", GetConfig().Socket), Handler: router}
	fmt.Println("listening", *addr)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Run", err)
		} else {
			log.Info("main", "http.server exit normally")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Warning("closing run")
			return
		}
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Listening ")
		next.ServeHTTP(w, r)
		log.Println("Disconnect ")

	})

}
