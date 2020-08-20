//go:generate protoc -I=. --go_out=. edgeInfo.proto
package protobuf

import "github.com/golang/protobuf/proto"
import "github.com/sirupsen/logrus"

func GetBufEdgeInfo(e *EdgeInfo) []byte {
	b, err := proto.Marshal(e)
	if err != nil {
		logrus.Warning(err)
	}
	return b
}
