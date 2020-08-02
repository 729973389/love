package main

import (
	"context"
	"log"

	// "encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"


	"mergerd/app"
	"mergerd/libmerger"
)

var (
	version = "N/A"
	build   = "N/A"
)

func serveHTTP(wg *sync.WaitGroup, addr string, router http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("main", err)
		} else {
			log.Info("main", "http.Server exit normally:", addr)
		}
	}()

	return srv
}

// func handleAdminDvrTasks(rsp http.ResponseWriter, req *http.Request) {
// 	rsp.Header().Set("Connection", "close")
// 	if b, err := appSvc.GetAllDvrTasks(); err != nil {
// 		log.Error("main", "getAllDvrTasks:", err)
// 		rsp.WriteHeader(500)
// 	} else {
// 		rsp.Header().Set("Content-Type", "application/json")
// 		rsp.Write(b)
// 	}
// }

// fixme: cgo https://golang.org/pkg/os/signal/
func waitForSignals() {
	sigChan := make(chan os.Signal)

	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGPIPE,
	)

	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			log.Info("main", "Received SIGINT.")
			return
		case syscall.SIGTERM:
			log.Debug("main", "Received SIGTERM.")
			return
		case syscall.SIGPIPE:
			log.Info("main", "Received SIGPIPE")
		}
	}
}

func main() {
	log.Infof("main", "Version:%s,Build:%s", version, build)
	//signal.Ignore(os.SIGPIPE)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	libmerger.Start(ctx, &wg)

	// routerAdmin := http.NewServeMux()
	// router.Handle("/metrics", promhttp.Handler())
	// router.HandleFunc("/admin/dvrtasks", handleAdminDvrTasks)
	// wg.Add(1)
	// httpSrv := serveHTTP(&wg, ":9123", router)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/vm/autoPictureMatch", app.HandleVMAutoPictureMatch)
	router.HandleFunc("/api/v1/vm/autoPointMatch", app.HandleVMAutoPointMatch)
	router.HandleFunc("/api/v1/vm/list", app.HandleVM)
	router.HandleFunc("/api/v1/vm/hotarea", app.HandleVMHotArea)
	router.HandleFunc("/api/v1/vm/video/list", app.HandleVMVideoList)
	router.HandleFunc("/api/v1/vm/camera/list", app.HandleVMCameraList)
	router.HandleFunc("/api/v1/vm/camera/preview", app.HandleVMCameraPreview)
	router.HandleFunc("/api/v1/vm/camera/IndePreview", app.HandleIndeCameraPreview)
	router.HandleFunc("/api/v1/vm/camera/config", app.HandleVMCameraConfig)
	router.HandleFunc("/api/v1/vm/cameraOrder", app.HandleVMCameraOrderSet)

	router.HandleFunc("/api/v1/vm/intrusion/ctrl", app.HandleVMIntrusionDetectCtrl)
	router.HandleFunc("/api/v1/vm/intrusion/areaCtrl", app.HandleVMIntrusionDetectAreaCtrl)
	router.HandleFunc("/api/v1/vm/intrusion/config", app.HandleVMIntrusionDetectConfig)
	router.HandleFunc("/api/v1/vm/intrusion/setArea", app.HandleVMSetIntrusionDetectArea)

	router.HandleFunc("/api/v1/vm/count/area", app.HandleVMCountArea)
	router.HandleFunc("/api/v1/vm/count/ctrl", app.HandleVMCountCtrl)

	//router.HandleFunc("/ws/intrusion", app.HandleVMWsIntrusion)
	router.HandleFunc("/ws/count", app.HandleVMWsCount)
	router.HandleFunc("/ws", app.HandleVMWs)

	//router.HandleFunc("/api/v1/vm/tasks", app.HandleVMTasksGet)
	wg.Add(1)
	httpSrv := serveHTTP(&wg, ":8010", router)

	log.Info("main", "started")
	waitForSignals()
	log.Info("main", "shuting down...")
	cancel()

	//graceful shutdown:
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelShutdown()

	httpSrv.Shutdown(ctxShutdown)

	// wait until all modules are done
	wg.Wait()
	log.Info("main", "byebye")
}
