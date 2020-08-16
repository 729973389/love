package http

import (
	"log"
	http "net/http"
)




func init() {
	http.Handle("/lxd/",http.StripPrefix("/lxd/",http.FileServer(http.Dir("/home/ubuntu/lovefile"))))
	http.Handle("/jbzhu/",http.StripPrefix("/jbzhu/",http.FileServer(http.Dir("/home/ubuntu/jbzhufile"))))
	log.Println("serving ")
     log.Fatal(http.ListenAndServe(":80",nil))
}

