package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alecthomas/kong"
)


type CLI struct{// {{{
	Send Send `short:"s" help:"flag for sending the buffer"`
	Recv bool `short:"r" help:"flag for receiving the buffer"`
}
type Send struct{
	Buff string `arg:"" help:"buffer to be sent"`
}
func (c *CLI)Run()error{
	if c.Recv{
		//
	} else if c.Send.Buff!=""{
		//
	}
	return nil
}// }}}

const PHRASE = "BUFFINGA"// {{{
const DEFAULT_PORT = ":8901"
func UDPBroadcast(broadPort string)error{
	broadIP, err:=getSubnet()
	if err!=nil{
		return fmt.Errorf("Subnet Error: %w", err)
	} else {
		fmt.Println("Broadcasting at", broadIP)
	}
	broadAddr:= net.JoinHostPort(broadIP, broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		return fmt.Errorf("Resolve Error: %w",err)
	}
	conn, err:=net.DialUDP("udp4", nil, addr)
	if err != nil{
		return fmt.Errorf("Dial Error: %w",err)
	}
	defer conn.Close()
	msg:=[]byte(PHRASE)
	fmt.Println("Waiting for Receiver ...")
	for {
		conn.Write(msg)
		time.Sleep(time.Second*2)
	}
}
func ListenBroadcast(broadPort string)(net.UDPAddr, error){
	broadAddr := net.JoinHostPort("", broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		return net.UDPAddr{}, fmt.Errorf("Resolve Error: %w", err)
	}
	conn, err:=net.ListenUDP("udp4", addr)
	if err != nil{
		return net.UDPAddr{}, fmt.Errorf("Listen Error: %w", err)
	}
	defer conn.Close()
	fmt.Println("Waiting for Sender ...")
	buf:=make([]byte,1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return net.UDPAddr{}, fmt.Errorf("Reading from udp Error: %w", err)
		}
		msg := string(buf[:n])
		fmt.Printf("Received '%s' from %s\n", msg, remoteAddr)
		if msg == PHRASE{
			return *remoteAddr, nil
		} else{
			continue
		}
	}
}
func getSubnet()(string, error){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4()
				mask := ipnet.Mask
				broadcast := net.IP(make([]byte, 4))
				for i := range ip {
					broadcast[i] = ip[i] | ^mask[i]
				}
				return broadcast.String(), nil
			}
		}
	}
	return "255.255.255.255", nil
}// }}}

func SendBuff(conn net.Conn, buff string)error{
	UDPBroadcast(DEFAULT_PORT)
}

func RecvBuff(conn net.Conn)(string, error){

}

func main(){
	var cli CLI
	ctx, err := kong.Parse(&cli)
	if err!=nil{
		log.Fatalf("Kong parsing error: %v", err)
	}
	ctx.Run()
}

