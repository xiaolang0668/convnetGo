package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"net"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"
)

const (
	ConstSaveDataLength = 4
)

var MainClientInfo *mainClientInfo

func loadClientConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &client)
	if err != nil {
		return err
	}

	return nil
}

func reader(readerChannel chan []byte) {

	for {
		select {
		case data := <-readerChannel:
			log.Debug((string)(data))
			//解析为clientmessage对象
			var clientMessage clientMessage
			//Log(string(data))
			err := json.Unmarshal([]byte(data), &clientMessage)

			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}
			switch clientMessage.CMDType {
			case WS_REGISTE_RESP:
				{
					client.PublicID = clientMessage.Message[0].(string)
					log.Print("我的PUBLICID:", "CVNID://"+client.Server+":"+client.ServerPort+"@"+client.PublicID)
					//client.ClientID = strings.Split(client.PublicID, ":")[1]
					go TapLoopData()
					go LoadAutoConnectPeerList()
					//sendToConn(client.g_conn, C_GETWS_SERVER_INFO, []interface{}{client.PublicID})
					//log.Info("获取对方服务信息：", client.PublicID)
				}
			case WS_REGISTE_FAIL:
				{
					log.Fatal("客户端UUID已经在线，程序退出")
				}
			case C_CONNTOWS_PEERCALL_RESP:
				{
					log.Print("对方请求连接，准备应答", (clientMessage.Message[0].(map[string]interface{})["PublicID"].(string)))

					PublicID := clientMessage.Message[0].(map[string]interface{})["PublicID"].(string)

					SDP := clientMessage.Message[1].([]interface{})[0].(string)
					step := clientMessage.Message[2].(interface{}).(string)
					accesspass := ""
					if len(clientMessage.Message) > 3 {
						accesspass = clientMessage.Message[3].(interface{}).(string)
					}

					log.Debug(PublicID, "《《《《《《《《《《《《》》》》》》》》》》》》》》", client.PublicID)

					if PublicID == client.PublicID {
						log.Info("不可以呼叫自己")
					}

					user := GetUserByPublicID(PublicID)
					if user == nil {
						user = new(User)
					}

					if step == "E" {
						user.AccessPass = "ERROR"
						log.Info("密码错误")
						continue
					}

					if !user.AllowConnect && client.AutoConnectPassword != "" {

						if accesspass != client.AutoConnectPassword {
							log.Print(accesspass, client.AutoConnectPassword)
							log.Print("对方请求连接，密码错误，拒绝连接", PublicID)
							sendToConn(client.g_conn, C_CONNTOWS_PEERCALL, []interface{}{PublicID, "", "E"})
							continue
						}
					}

					clientinfo := formatjsontoMainClientInfo(clientMessage)
					user.PublicID = clientinfo.PublicID
					user.MacAddress = clientinfo.Mac
					user.UserNickName = clientinfo.Name
					user.IsOnline = true //默认在线
					id := Strtoint(strings.Split(clientinfo.PublicID, ":")[1])
					RenewUserList(user)
					user.CvnIP = string(GetCvnIPstring(id))

					user.SDP = SDP
					log.Debug("更新SDP", step)
					go peerConnectionUpdate(user, step) //更新SDP
					//对方向我发起连接
					log.Debug(clientMessage.Message)

				}
			case C_CONNTOWS_PEERCALL_RESP_NOTONLINE:
				{
					log.Info(clientMessage.Message[0].(string), "对方不在线")
					//DelUserFromUserList(clientMessage.Message[0].(string))
					PublicID := clientMessage.Message[0].(string)
					user := GetUserByPublicID(PublicID)
					user.IsOnline = false
					user.Pc = nil
					user.Dc = nil
				}
			case C_GETWS_SERVER_INFO_RESP: //获取我在服务器的身份
				{
					MainClientInfo = formatjsontoMainClientInfo(clientMessage)
					log.Info("服务信息：", ToJson(MainClientInfo))
				}
			}
		}
	}
}

func formatjsontoMainClientInfo(clientMessage clientMessage) *mainClientInfo {
	var MainClientInfo *mainClientInfo
	MainClientInfo = new(mainClientInfo)
	MainClientInfo.PublicID = clientMessage.Message[0].(map[string]interface{})["PublicID"].(string)
	MainClientInfo.Version = clientMessage.Message[0].(map[string]interface{})["Version"].(string)
	MainClientInfo.Name = clientMessage.Message[0].(map[string]interface{})["Name"].(string)
	MainClientInfo.ClientMode = clientMessage.Message[0].(map[string]interface{})["ClientMode"].(float64)
	MainClientInfo.Mac = clientMessage.Message[0].(map[string]interface{})["Mac"].(string)
	MainClientInfo.IP = clientMessage.Message[0].(map[string]interface{})["IP"].(string)
	//clientMessage.Message[4]转换为[]PortRange
	MainClientInfo.AllowTcpPortRange = ToPortRange(clientMessage.Message[0].(map[string]interface{})["AllowTcpPortRange"].([]interface{}))
	MainClientInfo.AllowUdpPortRange = ToPortRange(clientMessage.Message[0].(map[string]interface{})["AllowUdpPortRange"].([]interface{}))
	return MainClientInfo
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// 解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstSaveDataLength {
			break
		}

		messageLength := BytesToInt(buffer[i : i+ConstSaveDataLength])
		if length < i+ConstSaveDataLength+messageLength {
			break
		}
		data := buffer[i+ConstSaveDataLength : i+ConstSaveDataLength+messageLength]
		readerChannel <- data[0:messageLength]

		i += ConstSaveDataLength + messageLength - 1
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	for {
		buffer := make([]byte, BUFFERSIZE)
		n, err := conn.Read(buffer)

		if err != nil {
			log.Info(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		// if buffer[n-1] == 10 && buffer[n-2] == 13 {
		// 	readerChannel <- buffer[:n]
		// }
		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}
func Strtoint(intstr string) int {
	i, _ := strconv.ParseInt(intstr, 10, 0)
	return int(i)
}

func Strtoint64(intstr string) int64 {
	i, _ := strconv.ParseInt(intstr, 10, 0)
	return i
}

func Inttostr(intnum int) string {
	return strconv.Itoa(intnum)
}

func ProtocolToStr(protostr int) string {
	return strconv.Itoa(protostr)
}

func StrToProtocol(str string) int {
	i, _ := strconv.ParseInt(str, 10, 0)
	return int(i)
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func IntToBytes(i int32) []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, i)
	return byteBuffer.Bytes()
}
func sendCmd(str string) {
	sendCmdBuff([]byte(str))
}
func sendCmdBuff(context []byte) {
	client.g_conn.Write(append(append([]byte(""), IntToBytes(int32(len(context)))...), context...))
	// var buffer bytes.Buffer
	// asd := IntToBytes(int32(len(context)))
	// buffer.Write(asd)
	// buffer.Write(context)
	// c.Write(buffer.Bytes())
}
