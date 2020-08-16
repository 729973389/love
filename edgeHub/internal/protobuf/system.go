//go:generate protoc -I=. --go_out=. system.proto
package protobuf

import (
	"log"
	"os"
	"runtime"
)

const pwd = "internal/protobuf/system.go"

func GetSystemInfo() SystemInfo {
	var s SystemInfo
	s.CPUNum = int32(runtime.GOMAXPROCS(0))
	s.ARC = runtime.GOARCH
	s.OS = runtime.GOOS
	name, err := os.Hostname()
	if err != nil {
		log.Printf("%s:%s \n", pwd,err)
	}
	s.HostName = name
	cpuInfo := GetCPUInfo()
	s.CPUInfo=&cpuInfo
	memInfo:=GetMEMInfo()
	s.MEMInfo=&memInfo
	diskInfo := GetDiskInfo()
	s.DiskInfo=&diskInfo
	return s
}
