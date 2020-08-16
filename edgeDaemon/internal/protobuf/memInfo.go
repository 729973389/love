//go:generate protoc -I=. --go_out=. memInfo.proto
package protobuf

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)


func GetMEMInfo () MEMInfo{


	cmd := exec.Command("vmstat","-s")
	var out bytes.Buffer
	cmd.Stdout=&out
	err:=cmd.Run()
	if err != nil {
		log.Println(err)
		return MEMInfo{}
	}

	process :=make([][]string,0)

	for i:=0;i<6;i++ {
		lines,err := out.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		token := strings.Split(lines," ")
		ft:=make([]string ,0)
		for _,t := range token {
			if t!=""&&t!="\t"{
				ft=append(ft,t)
			}
		}
		process=append(process,ft)
	}
	memInfo:=MEMInfo{
		Total:process[0][0],
		Used:process[1][0],
		Free:process[4][0],
		Buffer:process[5][0],
	}
	return memInfo
}
