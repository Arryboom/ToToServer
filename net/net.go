package net
import (
	"net"
	"log"
	. "ToToServer/config"
	"sync"
	"time"
	"bytes"
	"ToToServer/util"
)
func ListenMode2Tcp(remoteAddr string,remotePort string,localPort string,HTTPHeaderHostReplace string){
	l,err:=net.Listen("tcp",":"+localPort)
	if err!=nil{
		log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s]监听失败，原因：%s",remoteAddr,remotePort,localPort,err.Error())
		return
	}
	log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s]监听成功",remoteAddr,remotePort,localPort)
	defer func(){
		l.Close()
	}()
	for{
		c,err:=l.Accept()
		if err!=nil{
			log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]客户端接受失败，原因：%s",remoteAddr,
				remotePort,localPort,c.RemoteAddr().String(),err.Error())
			continue
		}
		if IsDebug(){
			log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]客户端连接",remoteAddr,
				remotePort,localPort,c.RemoteAddr().String())
		}
		r,err:=net.Dial("tcp",remoteAddr+":"+remotePort)
		if err!=nil{
			log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]远程地址连接失败，原因：%s",remoteAddr,
				remotePort,localPort,c.RemoteAddr().String(),err.Error())
			continue
		}
		if IsDebug(){
			log.Printf("模式2：TCP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]客户端和远程对接成功",remoteAddr,
				remotePort,localPort,c.RemoteAddr().String())
		}
		go tcpBridgeM2(c,r,HTTPHeaderHostReplace!="",HTTPHeaderHostReplace)
		go tcpBridgeM2(r,c,false,"")
	}
}
func tcpBridgeM2(src net.Conn,dst net.Conn,rpl bool,host string){
	defer func(){
		src.Close()
	}()
	buf:=make([]byte,65536)
	for{
		n,err:=src.Read(buf)
		Realbuf:=buf[:n]
		if err!=nil{
			return
		}
		if rpl{
			if pos1:=bytes.Index(Realbuf,[]byte("\r\nHost:"));pos1!=-1{
				if pos2:=bytes.Index(Realbuf[pos1+7:],[]byte("\r\n"));pos2!=-1{
					oldhost:=Realbuf[pos1+7:pos1+7+pos2]
					Realbuf=bytes.Replace(Realbuf,[]byte("\r\nHost:"+util.B2s(oldhost)),[]byte("\r\nHost:"+host),1)
				}
			}
		}
		dst.Write(Realbuf)
	}
}
func ListenMode2Udp(remoteAddr string,remotePort string,localPort string){
	udp_addr, _ := net.ResolveUDPAddr("udp", ":"+localPort)
	l,err:=net.ListenUDP("udp4",udp_addr)
	if err!=nil{
		log.Printf("模式2：UDP[远程地址:%s,远程端口:%s,本地端口:%s]监听失败，原因：%s",remoteAddr,remotePort,localPort,err.Error())
		return
	}
	log.Printf("模式2：UDP[远程地址:%s,远程端口:%s,本地端口:%s]监听成功",remoteAddr,remotePort,localPort)
	defer func(){
		l.Close()
	}()
	buf:=make([]byte,65536)
	for{
		n,addr,err:=l.ReadFromUDP(buf)
		if err!=nil{
			continue
		}
		r:=getUDPRemoteConnM2(addr)
		if r!=nil{
			(*r).Write(buf[:n])
			continue
		}
		remoteConn,err:=net.Dial("udp", remoteAddr+":"+remotePort)
		if err!=nil{
			log.Printf("模式2：UDP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]远程地址连接失败，原因：%s",remoteAddr,
				remotePort,localPort,addr.String(),err.Error())
			continue
		}
		if IsDebug(){
			log.Printf("模式2：UDP[远程地址:%s,远程端口:%s,本地端口:%s,客户端地址:%s]客户端和远程对接成功",remoteAddr,
				remotePort,localPort,addr.String())
		}
		remoteConn.Write(buf[:n])
		setUDPRemoteConnM2(addr,&remoteConn)
		go checkTimeout(addr)
		go udpBridgeM2(l,addr,&remoteConn)
	}
}
type udpRemoteConn struct {
	Conn *net.Conn
	Timeout time.Time
}
var mapUDPRemote map[string]*udpRemoteConn
var mapUDPRemoteLock sync.Mutex
func getUDPRemoteConnM2(clientAddr *net.UDPAddr)(*net.Conn){
	addr:=clientAddr.String()
	mapUDPRemoteLock.Lock()
	r,e:=mapUDPRemote[addr]
	mapUDPRemoteLock.Unlock()
	if e{
		return r.Conn
	}
	return nil
}
func setUDPRemoteConnM2(clientAddr *net.UDPAddr,remoteConn *net.Conn){
	addr:=clientAddr.String()
	mapUDPRemoteLock.Lock()
	mapUDPRemote[addr]=&udpRemoteConn{Conn:remoteConn,Timeout:time.Now()}
	mapUDPRemoteLock.Unlock()
}
func updateTimeout(clientAddr *net.UDPAddr){
	addr:=clientAddr.String()
	mapUDPRemoteLock.Lock()
	r,e:=mapUDPRemote[addr]
	if e{
		r.Timeout=time.Now()
	}
	mapUDPRemoteLock.Unlock()
}
func checkTimeout(clientAddr *net.UDPAddr){
	sec:=float64(ServerConfig.Mode2UdpAliveTime)
	for{
		mapUDPRemoteLock.Lock()
		c:=mapUDPRemote[clientAddr.String()]
		if time.Now().Sub(c.Timeout).Seconds()>sec{
			if IsDebug(){
				log.Printf("模式2：UDP[远程地址:%s,客户端地址:%s]心跳(%g秒)超时，已断开远程连接",(*c.Conn).RemoteAddr(),clientAddr.String(),sec)
			}
			(*c.Conn).Close()
			delete(mapUDPRemote,clientAddr.String())
			mapUDPRemoteLock.Unlock()
			return
		}
		mapUDPRemoteLock.Unlock()
		time.Sleep(time.Second*1)
	}
}
func udpBridgeM2(udpConn *net.UDPConn,clientAddr *net.UDPAddr,remoteConn *net.Conn){
	buf:=make([]byte,65536)
	for{
		n,err:=(*remoteConn).Read(buf)
		if err!=nil{
			return
		}
		updateTimeout(clientAddr)
		(*udpConn).WriteToUDP(buf[:n],clientAddr)
	}
}
func init(){
	mapUDPRemote=make(map[string]*udpRemoteConn)
}