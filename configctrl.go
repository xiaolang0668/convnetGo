package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/labstack/gommon/log"
)

func DeleteLineFromFile(filename, publicID string) {
	// 打开文件（读取）
	file, err := os.Open(filename)
	if err != nil {
		log.Error("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 读取所有行并过滤掉要删除的行
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, publicID) {
			lines = append(lines, line)
		}
	}

	// 重新打开文件用于写入（清空原内容）
	newFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Error("无法重新打开文件用于写入:", err)
		return
	}
	defer newFile.Close()

	// 将过滤后的内容写回文件
	writer := bufio.NewWriter(newFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Error("写入文件失败:", err)
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		log.Error("刷新缓冲区失败:", err)
	}
}

func AppendLineToFile(filename, publicID string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("无法打开autoConnectPeer.txt文件:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(publicID + "\n")
	if err != nil {
		log.Error("无法写入autoConnectPeer.txt文件:", err)
	}
}

func IsPublicIDInAutoConnectPeer(publicID string) bool {
	file, err := os.Open("autoConnectPeer.txt")
	if err != nil {
		log.Error("无法打开autoConnectPeer.txt文件:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, publicID) {
			return true
		}
	}
	return false
}

func saveClientConfig(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	//排除一部分属性进行保存
	// 创建一个新的 ClientConfig 实例，并填充数据
	clientConfig := ClientConfig{
		Server:              client.Server,
		UUID:                client.UUID,
		ClientID:            client.ClientID,
		ServerPort:          client.ServerPort,
		ServerTurnPort:      client.ServerTurnPort,
		ServerTurnUser:      client.ServerTurnUser,
		ServerTurnPass:      client.ServerTurnPass,
		ServerTurnRealm:     client.ServerTurnRealm,
		AutoConnectPassword: client.AutoConnectPassword,
		AllowTcpPortRange:   client.AllowTcpPortRange,
		AllowUdpPortRange:   client.AllowUdpPortRange,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(clientConfig)
}
