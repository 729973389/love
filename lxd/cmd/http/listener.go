package http

import (
	"html/template"
	"log"
	http "net/http"
)


var SendMSG  = func(r http.ResponseWriter,w *http.Request) {
	t,_:=template.ParseFiles("./internal/html/index.html")
	b := map[string]string{
		"name":"lxd",
		"someStr":"hello",
	}
		t.Execute(r,b)

}
func Root(r http.ResponseWriter,w *http.Request) {
	t,_:=template.ParseFiles("./internal/html/index.html")
	b := map[string]string{
		"name":"lxd",
		"someStr":"hello",
	}
	t.Execute(r,b)

}

func init() {
	////http.Handle("/lxd",http.StripPrefix("/lxd",http.FileServer(http.Dir("/home/ubuntu/lovefile"))))
	//http.Handle("/jbzhu",http.StripPrefix("/jbzhu",http.FileServer(http.Dir("/home/ubuntu/jbzhufile"))))
	log.Println("serving ")
     log.Fatal(http.ListenAndServe(":54321",http.FileServer(http.Dir("/home/ubuntu/lovefile"))))
}

