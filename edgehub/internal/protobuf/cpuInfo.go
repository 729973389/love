//go:generate protoc --proto_path=. --go_out=. cpuInfo.proto
package protobuf

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	CPU float32
}

func GetCPUInfo() CPUInfo {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	process := make([]*Process, 0)
	for {
		lines, err := out.ReadString('\n')
		if err != nil {
			break
		}
		token := strings.Split(lines, " ")
		ft := make([]string, 0)
		for _, t := range token {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		cpu, err := strconv.ParseFloat(ft[2], 32)
		if err != nil {
			cpu = 0
		}
		process = append(process, &Process{CPU: float32(cpu)})
	}
	var cpuTotal float32 = 0
	for _, p := range process {
		cpuTotal += p.CPU
	}
	cpuInfo := CPUInfo{
		CPUUsed: cpuTotal,
	}
	return cpuInfo
}
