package root

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var addr = flag.String("port", ":43211", "http service address,e.g. :43211")

//the entrance
func Run(ctx context.Context) {
	defer log.Warning("EXIT: RUN")
	hub := NewHub()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		hub.Run(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		Serve(ctx, hub)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunDWS(ctx, hub)
	}()
	flag.Parse()
	//Mux holds the map that server looks up from pattern to handler
	router := http.NewServeMux()
	//hook the handler,do something before or after it.
	router.Handle("/edgeHub", Middleware(http.Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		Servews(ctx, hub, writer, request)
	}))))
	server := http.Server{Addr: fmt.Sprintf(":%s", Info.Socket), Handler: router}
	defer server.Close()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("server: ", err)
		} else {
			log.Info("main", "http.server exit normally")
		}
	}()
	wg.Wait()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Listening: %s", r.RemoteAddr))
		next.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("Disconnect: %s", r.RemoteAddr))
	})
}

//FindKeyString finds string value when given a specified parameter-list(json,key string)
func FindKeyString(s string, key string) (string, error) {
	if !strings.Contains(s, ",") {
		if strings.Contains(s, "\""+key+"\"") {
			t := strings.Split(s, "\"")
			for i, t2 := range t {
				if strings.Contains(t2, ":") {
					if t[i-1] == key {
						return t[i+1], nil
					}
				}
			}
		}
	}
	line := strings.Split(s, ",")
	for _, v := range line {
		if strings.Contains(v, "\""+key+"\":") {
			t := strings.Split(v, "\"")
			for i, t2 := range t {
				if strings.Contains(t2, ":") {
					if t[i-1] == key {
						return t[i+1], nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("FindKeyString: can't find %s from %s", key, s)
}

var getEdgeDeviceClient *http.Client

func init() {
	getEdgeDeviceClient = &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
				Timeout:   60 * time.Second}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		},
	}
}

//get edge device relation when given a serialNumber.
func GetEdgeDevice(s string) {
	var rout = Info.GetEdgeDevice + "?" + s
	req, err := http.NewRequest("GET", rout, io.ReadCloser(nil))
	if err != nil {
		log.Error(errors.Wrap(err, "GetEdgeDevice"))
		return
	}
	req.Header.Add("token", Info.GetEdgeDeviceToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := getEdgeDeviceClient.Do(req)
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDevice"))
		return
	}
	defer resp.Body.Close()
	read := make([]byte, 512)
	_, _ = resp.Body.Read(read)
	fmt.Println(string(read))
}
