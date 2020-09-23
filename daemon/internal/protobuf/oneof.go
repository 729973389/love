//go:generate protoc -I=. --go_out=. oneof.proto
package protobuf

func SetOneOfEdgeInfo(number string) *Message {
	systemInfo := GetSystemInfo()
	m := &Message{Switch: &Message_EdgeInfo{EdgeInfo: &EdgeInfo{SerialNumber: number, Data: &systemInfo}}}
	return m
}

func SetOneOfAuthor(number string, token string, time string, hmac string) *Message {
	author := Author{SerialNumber: number, Token: token, Time: time, Hmac: hmac}
	m := &Message{Switch: &Message_Author{Author: &author}}
	return m
}
