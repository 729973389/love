//go:generate protoc -I=. --go_out=. edgeProperties.proto
package protobuf

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func GetEdgeProperties() EdgeProperties {
	p := &EdgeProperties{
		CPUNum: int32(runtime.GOMAXPROCS(0)),
		ARC:    runtime.GOARCH,
		OS:     runtime.GOOS,
	}
	name, err := os.Hostname()
	if err != nil {
		log.Warning(errors.Wrap(err, "getEdgeProperties"))
		return *p
	}
	p.HostName = name
	return *p
}
