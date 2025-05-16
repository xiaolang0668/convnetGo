package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/songgao/water"
)

func setupIfce(ipNet net.IPNet, dev string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		sargs := fmt.Sprintf("interface ip set address name='REPLACE_ME' source=static addr=REPLACE_ME mask=REPLACE_ME gateway=none")
		args := strings.Split(sargs, " ")
		args[4] = fmt.Sprintf("name=%s", dev)
		args[6] = fmt.Sprintf("addr=%s", ipNet.IP)
		args[7] = fmt.Sprintf("mask=%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])
		cmd = exec.Command("netsh", args...)
	} else if runtime.GOOS == "linux" {
		ipAddr := ipNet.IP.String()
		ipMask := net.IP(ipNet.Mask).String()
		cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", ipAddr, ipMask), "dev", dev)
		log.Info("cmdexec: ", cmd.String())
		if err := cmd.Run(); err != nil {
			log.Errorf("Failed to assign IP address: %v", err)
			return
		}

		//启动网卡
		time.Sleep(time.Second * 1)
		cmd = exec.Command("ip", "link", "set", dev, "up")
		log.Info("cmdexec: ", cmd.String())
		if err := cmd.Run(); err != nil {
			log.Errorf("Failed to bring up the interface: %v", err)
			return
		}
	} else {
		log.Info("Unsupported OS:", runtime.GOOS)
		return
	}

	log.Info(cmd.String())
	if err := cmd.Run(); err != nil {
		log.Info(err)
	}
}

func setupArpinfo(mac, ip string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("arp", "-d", "*")
		if err := cmd.Run(); err != nil {
			log.Info(err)
		}
		//将:替换为-
		mac = strings.Replace(mac, ":", "-", -1)
		cmd = exec.Command("arp", "-s", ip, mac)
	}
	if runtime.GOOS == "linux" {
		cmd := exec.Command("arp", "-s", ip, mac)
		if err := cmd.Run(); err != nil {
			log.Info(err)
		}
	}
}

func teardownIfce(ifce *water.Interface) {
	client.g_ifce = nil
	if err := ifce.Close(); err != nil {
		log.Info(err)
	}
}
