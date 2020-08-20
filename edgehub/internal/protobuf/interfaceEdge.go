//go:generate protoc -I=. --go_out=. interfaceEdge.proto
package protobuf

import "github.com/golang/protobuf/proto"
import "github.com/sirupsen/logrus"

func ReadEdge(b []byte) *InterfaceEdge {
	edge := InterfaceEdge{}
	err := proto.Unmarshal(b, &edge)
	if err != nil {
		logrus.Warning(err)
	}
	return &edge
}
