//go:generate protoc -I=. --go_out=. system.proto
package protobuf

import (
	"log"
	"net"
	"os"
	"runtime"
)

const pwd = "internal/protobuf/system.go"

func GetSystemInfo() SystemInfo {
	var mynet = []*NetInterfaces{new(NetInterfaces)}
	var s SystemInfo
	netInterface, err := net.Interfaces()
	if err != nil {
		log.Printf("%s : %s \n", pwd, err)
	}
	for i1,_ := range netInterface {
		v1 := netInterface[i1]
		addr := v1.HardwareAddr
		if addr.String() == "" {
			log.Printf("%s : cant read MAC addr", pwd)
			break
		}
		ip,err := v1.Addrs()
		if err != nil {
			log.Printf("%s : %s ",pwd,err)
			break
		}
		v1_ := mynet[i1]
		for i2,_ :=range ip{
			v1_.IPAddr=append(v1_.IPAddr,ip[i2].String())
		}
		v1_.HardwareAddr=addr.String()
		v1_.Name=v1.Name
		v1_.Index=int32(v1.Index)
		v1_.MTU=int32(v1.MTU)
		mynet=append(mynet,v1_)

	}
	s.CPUNum = int32(runtime.GOMAXPROCS(0))
	s.ARC = runtime.GOARCH
	s.OS = runtime.GOOS
	name, err := os.Hostname()
	if err != nil {
		log.Printf("internal/protobuf/system.go :%s \n", err)
	}
	s.HostName = name
	s.NetInterfaces = mynet
	cpuInfo := GetCPUInfo()
	s.CPUInfo=&cpuInfo
	memInfo:=GetMEMInfo()
	s.MEMInfo=&memInfo
	return s
}
