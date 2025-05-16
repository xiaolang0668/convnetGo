package main

import (
	"net"
	"strings"

	"github.com/labstack/gommon/log"
)

func ConnectServer(server string, port string) error {

	log.Info("ConnectTo:", server, ":", port)
	var err error
	//5秒超时

	service := server + ":" + port
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	client.g_conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

	var clientMessage clientMessage
	clientMessage.CMDType = 1
	//生成一个随机的guid
	UUID := client.UUID
	//测试UUID
	//驱动参数有-t
	// if client.ClientMode == MAINCLIENT {
	// 	UUID = "test"
	// } else {
	// 	UUID = "TEST"
	// }

	// read or write on conn
	go HandleConn()
	client.IsConnected = true
	client.RetryConnect = true

	sendToConn(client.g_conn, WS_REGISTE, []interface{}{"1.0", UUID, client.ClientID,
		client.ClientMode, client.Mac, client.MyCvnIP, client.AllowTcpPortRange, client.AllowUdpPortRange})

	return nil
}

func HandleConn() {
	defer func() {
		client.g_conn.Close()
		client.logout()
		client.IsConnected = false
		log.Info("client exit")
	}()
	handleConnection(client.g_conn)
}

func Split_string(s string) []string {
	a := strings.Split(s, ",")
	return a
}

func ExecComand(cmdField []string) {
	switch StrToProtocol(cmdField[0]) {

	default:
		log.Info("尚未实现", cmdField)
	}
}

var udpserver *net.UDPConn

func mymacstr() string {
	str := client.Mac
	return strings.ToUpper(strings.Replace(str, ":", "", -1))
}

func Getmymac(etherName string) string {

	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		//mac := inter.HardwareAddr //获取本机MAC地址
		if etherName == inter.Name {
			//fmt.Println("MAC = ", mac)
			return inter.HardwareAddr.String()
		}
	}

	return ""
}

func cmdCalltoUserRespDecode(cmdField []string) {

}
