package root

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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
	s := GetEdgeDaemon()
	if len(s)!=0{
		for _, number := range s {
			deviceId := GetEdgeDevice(number.SerialNumber)
			for _, singleDeviceId := range deviceId {
				hub.Bind(number.SerialNumber, singleDeviceId)
			}
		}
	}
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
func GetEdgeDevice(s string) (deviceId []string) {
	var rout = Info.EdgeInfoServer + Info.GetEdgeDevice + "?" + s
	req, err := http.NewRequest("GET", rout, io.ReadCloser(nil))
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDevice"))
		return
	}
	req.Header.Add("token", Info.EdgeInfoServerToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := getEdgeDeviceClient.Do(req)
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDevice: resp"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Printf("getEdgeDevice: httpStatusCode: %d: %v", resp.StatusCode, resp.Body))
		return
	}
	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDevice"))
	}
	message := string(read)
	fmt.Println(message)
	deviceId, err = FindAllKeyString(message, "deviceId")
	if err != nil {
		log.Error(err, "getEdgeDevice")
		return
	}
	return
}

//hold edgeDaemonInfo from func GetEdgeDaemon()
type EdgeDaemonInfo struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []SerialNumberData `json:"data"`
}
type SerialNumberData struct {
	SerialNumber string `json:"serialNumber"`
}

//get repeated serialNumber
func GetEdgeDaemon() (s []SerialNumberData) {
	var rout = Info.EdgeInfoServer + Info.GetActivate
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
	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Printf("getEdgeDevice: httpStatusCode: %d: %v", resp.StatusCode, resp.Body))
		return
	}
	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDaemon"))
	}
	log.Info("getEdgeDaemon: ", string(read))
	edgeDaemonInfo := &EdgeDaemonInfo{}
	err = json.Unmarshal(read, edgeDaemonInfo)
	if err != nil {
		log.Error(errors.Wrap(err, "getEdgeDaemon"))
		return
	}
	s = edgeDaemonInfo.Data
	log.Info(fmt.Printf("getEdgeDaemon: %v", s))
	return s
}

//find All the same key-value and return a alice of them.
func FindAllKeyString(s string, key string) ([]string, error) {
	sliceString := make([]string, 0)
	if !strings.Contains(s, ",") {
		if strings.Contains(s, "\""+key+"\"") {
			t := strings.Split(s, "\"")
			for i, t2 := range t {
				if strings.Contains(t2, ":") {
					if t[i-1] == key {
						sliceString = append(sliceString, t[i+1])
						return sliceString, nil
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
						sliceString = append(sliceString, t[i+1])
						break
					}
				}
			}
		}
	}
	if len(sliceString) != 0 {
		return sliceString, nil
	}
	return nil, fmt.Errorf("FindKeyString: can't find %s from %s", key, s)
}

func FindKey(s, key string) (string, error) {
	rs := []rune(s)
	plusRule := PlusRule()
	var fs []rune
	t := ""
	index := 0
	for _, v2 := range rs {
		if string(v2) == " " || string(v2) == "\n" || string(v2) == "\t" {
			continue
		}
		fs = append(fs, v2)
	}
	for i, v := range fs {
		if v == int32(key[index]) {
			if index == len(key)-1 {
				if fs[i-len(key)] != []rune("\"")[0] || fs[i+2] != []rune(":")[0] {
					index = 0
					continue
				}
				key := i + 3
				t += string(fs[key])
				for !plusRule(fs[key], fs[i+3]) {
					t += string(fs[i+4])
					i++
				}
				return t, nil
			}
			index++
			continue
		}
		index = 0
	}
	return "", fmt.Errorf(fmt.Sprintf("no %s in %s", key, s))
}

func PlusRule() func(rune, rune) bool {
	count := 0
	return func(s, n rune) bool {
		var head rune
		var tail rune
		if s == []rune("[")[0] {
			head, tail = []rune("[")[0], []rune("]")[0]
		} else if s == []rune("{" )[0] {
			head, tail = []rune("{" )[0], []rune("}" )[0]
		} else {
			log.Error(fmt.Errorf("findKey: set key: only { or [ is permmited,what i get is %s", string(s)))
			count = 0
		}
		if n == head {
			count++
		}
		if n == tail {
			count--
		}
		if count != 0 {
			return false
		}
		return true
	}
}
