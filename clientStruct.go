package main

// 定义一个用于保存配置的结构体
type ClientConfig struct {
	Server              string      `json:"Server"`
	ServerPort          string      `json:"ServerPort"`
	ServerTurnPort      string      `json:"ServerTurnPort"`
	ServerTurnUser      string      `json:"ServerTurnUser"`
	ServerTurnPass      string      `json:"ServerTurnPass"`
	ServerTurnRealm     string      `json:"ServerTurnRealm"`
	UUID                string      `json:"UUID"`
	ClientID            string      `json:"ClientID"`
	AutoConnectPassword string      `json:"AutoConnectPassword"`
	AllowTcpPortRange   []PortRange `json:"AllowTcpPortRange"`
	AllowUdpPortRange   []PortRange `json:"AllowUdpPortRange"`
}
