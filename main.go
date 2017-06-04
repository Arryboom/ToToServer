package main

import (
	"ToToServer/config"
	"ToToServer/net"
	"log"
	"strings"
	"time"
)
func main(){
	if config.ServerConfig.EnableMode2{
		log.Println("模式2开启：是")
		for _,s:=range config.ServerConfig.Mode2Config{
			if strings.ToLower(s.Protocol)=="tcp" {
				go net.ListenMode2Tcp(s.RemoteAddr,s.RemotePort,s.LocalPort,s.HTTPHeaderHostReplace)
			}else{
				go net.ListenMode2Udp(s.RemoteAddr,s.RemotePort,s.LocalPort)
			}
		}
	}
	for{
		time.Sleep(time.Second)
	}
}