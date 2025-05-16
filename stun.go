package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/pion/turn/v2"
)

var turnServer *turn.Server

// Define the StunClient struct
type StunClient struct {
	serverAddr string
	// Add other fields as necessary
}

// StunResponse represents a STUN server response
type StunResponse struct {
	// 基础字段
	mappedAddress net.UDPAddr
}

func (r *StunResponse) IsValid() bool {
	// 简单的有效性检查
	return r.mappedAddress.IP != nil && r.mappedAddress.Port != 0
}

func startTurn(publicIP string, listenPort int, realm string, username string, password string) {
	fmt.Printf("Starting TURN server with public IP: %s\n", publicIP)

	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(listenPort))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	usersMap := map[string][]byte{}
	usersMap[username] = turn.GenerateAuthKey(username, realm, password)
	fmt.Printf("Generated auth key for user: %s\n", username)

	turnServer, err = turn.NewServer(turn.ServerConfig{
		Realm: realm,
		// Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			fmt.Printf("Auth request from %s: username=%s, realm=%s\n", srcAddr.String(), username, realm)
			// framework will check auth key
			if key, ok := usersMap[username]; ok {
				fmt.Printf("Auth successful for user: %s\n", username)
				return key, true
			}
			fmt.Printf("Auth failed for user: %s\n", username)
			return nil, false
		},

		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(publicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",             // But actually be listening on every interface
				},
			},
		},
		// 添加事件处理器
		ChannelBindTimeout: time.Hour, // 通道绑定超时时间
	})

	if err != nil {
		log.Panic(err)
	}

	// 打印服务器配置信息
	fmt.Printf("TURN server configuration:\n")
	fmt.Printf("- Public IP: %s\n", publicIP)
	fmt.Printf("- Listen Port: %d\n", listenPort)
	fmt.Printf("- Realm: %s\n", realm)
	fmt.Printf("- Username: %s\n", username)
	fmt.Printf("- Channel Bind Timeout: %s\n", time.Hour)
	fmt.Printf("- Permission Timeout: %s\n", time.Hour)
	fmt.Printf("- Allocation Lifetime: %s\n", time.Hour*8)

	// 启动状态监控
	go func() {
		for {
			time.Sleep(time.Minute)
			if turnServer != nil {
				fmt.Printf("TURN server is running...\n")
			}
		}
	}()
}
