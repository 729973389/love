//go:generate protoc -I=. --go_out=. system.proto
package protobuf

const pwd = "internal/protobuf/system.go"

func GetSystemInfo() SystemInfo {
	var s SystemInfo
	cpuInfo := GetCPUInfo()
	s.CPUInfo = &cpuInfo
	memInfo := GetMEMInfo()
	s.MEMInfo = &memInfo
	diskInfo := GetDiskInfo()
	s.DiskInfo = &diskInfo
	return s
}
