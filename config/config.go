package config

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type Mode2ConfigSturct struct{
	RemoteAddr string
	RemotePort string
	Protocol string
	LocalPort string
	HTTPHeaderHostReplace string
}
type ServerConfigStruct struct{
	ServerPassword string
	ServerPort string
	EnableIPV6Support bool
	EnableMode1 bool
	EnableMode2 bool
	Mode2UdpAliveTime int
	Debug bool
	Mode2Config []Mode2ConfigSturct
}
var ServerConfig ServerConfigStruct
func IsDebug() bool{
	return ServerConfig.Debug
}
func init(){
	configJson,err:=ioutil.ReadFile("config.json")
	if err!=nil{
		fmt.Println("配置文件(config.json)加载错误！原因：",err.Error())
		os.Exit(-1)
	}
	jsonErr:=json.Unmarshal(configJson,&ServerConfig)
	if jsonErr!=nil{
		fmt.Println("配置文件(config.json)解析错误！原因：",err.Error())
		os.Exit(-1)
	}
}