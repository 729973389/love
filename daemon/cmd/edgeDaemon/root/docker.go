package root

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

//demo
func RunDocker() {
	c, err := http.Get("http://192.168.40.129:4333/test.tar")
	if err != nil {
		log.Error(err)
		return
	}
	file, err := os.OpenFile("test.tar", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return
	}
	tar, err := ioutil.ReadAll(c.Body)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = file.Write(tar)
	if err != nil {
		log.Error(err)
		return
	}
	cmd := exec.Command("docker", "load", "--input", "test.tar")
	var read bytes.Buffer
	read.Reset()
	cmd.Stdout = &read
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(read.String())
	cmd = exec.Command("docker", "run", "test:v1")
	read.Reset()
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(read.String())
}
