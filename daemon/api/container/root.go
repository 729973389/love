package container

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"os/exec"
	"strings"
)

type Docker struct {
	*protobuf.Docker
	Resp chan string
}

//pull images from cloud
func (c *Docker) Pull() error {
	cmd := "docker pull test"
	resp, err := Exec(cmd)
	if err != nil {
		return errors.Wrap(err, "docker pull")
	}
	resp = "docker pull: " + resp
	c.Resp <- resp
	return nil
}

//run the docker image with specific parameter
func (c *Docker) Run() error {
	run := "docker run "
	var restart, portExpose, portPublish, limitCPU, limitMEM, name string
	if c.Container.Restart != "" {
		restart = fmt.Sprintf("--restart %s ", c.Container.Restart)
	}
	if c.Container.PortExpose != "" {
		portExpose = fmt.Sprintf("--expose %s ", c.Container.PortExpose)
	}
	if c.Container.PortPublish != "" {
		portPublish = fmt.Sprintf("--publish %s ", c.Container.PortPublish)
	}
	if c.Container.LimitCPU != "" {
		limitCPU = fmt.Sprintf("--cpus %s ", c.Container.LimitCPU)
	}
	if c.Container.LimitMEM != "" {
		limitMEM = fmt.Sprintf("--kernel-memory %s ", c.Container.LimitMEM)
	}
	if c.Container.Name != "" {
		name = fmt.Sprintf("--name %s ", c.Container.Name)
	}
	cmd := run + restart + portExpose + portPublish + limitCPU + limitMEM + name
	resp, err := Exec(cmd)
	if err != nil {
		return errors.Wrap(err, "docker run")
	}
	resp = "docker run: " + resp
	c.Resp <- resp
	return nil
}

//remove exits Docker
func (c *Docker) Remove() error {
	cmd := fmt.Sprintf("docker stop %s", c.Container.Name)
	resp, err := Exec(cmd)
	if err != nil {
		return errors.Wrap(err, "docker stop")
	}
	resp = "docker stop " + resp
	c.Resp <- resp
	cmd = fmt.Sprintf("docker rm %s", c.Container.Name)
	resp, err = Exec(cmd)
	if err != nil {
		return errors.Wrap(err, "docker rm")
	}
	resp = "docker rm: " + resp
	c.Resp <- resp
	return nil
}
//Update delete the old container, then pull new image and run to create a new container
func (c *Docker) Update() (err error) {
	defer func() {
		errors.Wrap(err, "update")
	}()
	if err = c.Remove(); err != nil {
		return
	}
	if err = c.Pull(); err != nil {
		return
	}
	if err = c.Run(); err != nil {
		return
	}
	return nil
}

//Exec execute a set of string cmd and return the response
func Exec(s string) (r string, err error) {
	defer func() {
		//if err is nil,Wrap returns nil
		errors.Wrap(err, "exec")
	}()
	var parameter []string
	ss := strings.Split(s, " ")
	for _, v := range ss {
		parameter = append(parameter, v)
	}
	if len(parameter) < 2 {
		err = fmt.Errorf("command have no enough parameter")
		return
	}
	cmd := exec.Command(parameter[0], parameter[1:]...)
	var read bytes.Buffer
	cmd.Stdout = &read
	if err = cmd.Run(); err != nil {
		return
	}
	return read.String(), nil
}
