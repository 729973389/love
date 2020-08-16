//go:generate protoc -I=. --go_out=. diskInfo.proto
package protobuf
import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func GetDiskInfo()DiskInfo{
	cmd := exec.Command("df")
	var out bytes.Buffer
	cmd.Stdout=&out
	err:=cmd.Run()
	if err != nil {
		log.Println(err)
		return DiskInfo{}
	}
	var rs =make([]string,0)
	for{
		line,err :=out.ReadString('\n')
		if  err != nil {
			break
		}
		token := strings.Split(line," ")
		ft := make([]string,0)
		for _,v:=range token {
			strings.Replace(v,"\t","",-1)
			if v!=""{
				ft=append(ft,v)
			}
			if len(v)==2 {
				for _,j := range []byte(v) {
					if j=='/'{
						rs=append(ft,"")
					}
				}
			}

		}

	}
	d:=DiskInfo{
		Rate:rs[4],
		Used:rs[2],
		Free:rs[3],
	}
	return d
}

