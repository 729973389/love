package http

import (
	"bufio"
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/config"
	"io/ioutil"
	"os"
)
var urlFile = "http.json"
var socketFile="socket.json"
func GetConfig() *config.Url {
	for i := 0; i < 2; i++ {
		b, err := ioutil.ReadFile(urlFile)
		if err != nil {
			log.WithError(err).WithField("http.json", "Error")
			SetConfig()
			if i == 0 {
				continue
			}
			return nil
		}
		configUrl := &config.Url{}
		err = json.Unmarshal(b, configUrl)
		if err != nil {
			log.WithError(err).WithField("http.json", "Error")
			return nil
		}
		log.Println(configUrl)
		return configUrl
	}
	return nil

}

func SetConfig() {
	configUrl := &config.Url{}
	configUrl.SendData = "/api/v2/edge/data/create"
	b, err := json.MarshalIndent(configUrl, "", " ")
	if err != nil {
		log.WithError(err).WithField("marshal", "Error")
		return
	}
	file, err := os.OpenFile(urlFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(errors.Wrap(err, urlFile))
		return
	}
	writer := bufio.NewWriter(file)
	writer.Write(b)
	writer.Flush()

}

//import (
//	"bufio"
//	"encoding/json"
//	log "github.com/sirupsen/logrus"
//	"github.com/wuff1996/edgeHub/config"
//	"io/ioutil"
//	"os"
//
//	//"os"
//	//"text/template"
//	//"github.com/spf13/cobra"
//	"github.com/pkg/errors"
//	//"github.com/spf13/viper"
//)

//
//func init(){
//	viper.SetConfigType("json")
//	viper.SetConfigName("httpConfig")
//	viper.AddConfigPath(".")
//	viper.AddConfigPath("/etc/easyfetch/edgehub")
//	if err := viper.ReadInConfig();err != nil{
//		switch err.(type) {
//		case viper.ConfigFileNotFoundError:
//			log.Warning("No httpConfig file found")
//		default:
//			log.WithError(err).Fatal("Read httpConfig file error")
//
//		}
//	}
//
//}
//const configTemplate=`
//#edgehub http configfile
//[[http]]
//#interface=source addr interface
//#method=it's method
//interface={{.Http.Interface}}
//method={{.Http.Method}}

//`

//var configCmd=&cobra.Command{
//	Use: "configfile",
//	Short: "Print the edgehub configfile",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		t := template.Must(template.New("config").Parse(configTemplate))
//		err :=t.Execute(os.Stdout,&config.C)
//		if err != nil {
//			return errors.Wrap(err,"exec config template error")
//		}
//		return nil
//	},
//}

//var rootCmd = &cobra.Command{
//	Use:   "wuff",
//	Short: "I ",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		log.Println("running")
//		return nil
//	},
//}

//func init(){
//	cobra.OnInitialize(initConfig)
//	rootCmd.AddCommand(configCmd)
//
//}
//var cfgFile string
//func initConfig(){
//	if cfgFile!=""{
//		_,err := ioutil.ReadFile(cfgFile)
//		if err != nil {
//			log.WithError(err).WithField("config",cfgFile)
//		}
//	}
//
//}
//var run =func() {
//	log.Println("runing ")
//}
