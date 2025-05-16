package main

import "net"

const (
	ALL_DATA = iota //0       //发送数据
	WS_REGISTE
	WS_REGISTE_RESP
	WS_REGISTE_FAIL
	C_GETWS_SERVER_INFO
	C_GETWS_SERVER_INFO_RESP
	C_CONNTOWS_SERVER
	C_CONNTOWS_PEERCALL
	C_CONNTOWS_PEERCALL_RESP
	C_CONNTOWS_PEERCALL_RESP_NOTONLINE
	C_CONNTOWS_PEERCALL_FIN
	C_CONNTOWS_PEERDISCONNECT

	UNKNKOWN
)

const (
	CLIENTMODE = iota
	MAINCLIENT
	CLIENT
	SERVER
)

func sendToConn(conn net.Conn, msgtype int, message []interface{}) {
	var clientMessage clientMessage
	clientMessage.CMDType = msgtype
	clientMessage.Message = message
	clientMessage.Version = "1.0"
	jsonStr := ToJson(clientMessage) + "\r\n"
	//封包，发送两位大端序的长度
	len := (int32)(len(jsonStr))

	send := append(IntToBytes(len), []byte(jsonStr)...)
	if conn != nil {
		_, err := conn.Write(send)
		if err != nil {
			conn.Close()
		}
	}
}
