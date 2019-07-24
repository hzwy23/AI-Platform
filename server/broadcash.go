package server

import (
	"ai-platform/panda/config"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func GetBroadcast() []string {
	ipNets, err := GetIp()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return Get(ipNets)
}

func Get(ipNets []*net.IPNet) []string {
	broadcasts := make([]string, 0)
	for _, val := range ipNets {
		mask := val.Mask.String()
		ips := strings.Split(val.IP.String(), ".")
		ipStr := ""
		for i := 0; i < 4; i++ {
			m, _ := strconv.ParseUint(mask[2*i:2*(i+1)], 16, 10)
			ip, _ := strconv.ParseUint(ips[i], 10, 10)
			uip := (^uint8(m)) | uint8(ip)
			ipStr += strconv.Itoa(int(uip)) + "."
		}
		if len(ipStr) > 7 {
			broadcasts = append(broadcasts, ipStr[:len(ipStr)-1])
		}
	}
	return broadcasts
}

func GetIp() ([]*net.IPNet, error) {
	ipNets := make([]*net.IPNet, 0)
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return nil, err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ipNets = append(ipNets, ipNet)
					}
				}
			}
		}
	}
	return ipNets, nil
}

// 广播消息
func Broadcast(ip string, port string) {

	laddr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("创建UDP监听失败，失败原因是：", err)
		return
	}
	defer conn.Close()

	logger.Info("开启广播服务, broadcast ip is:", ip, ", port is:",port)
	for {
		platform.UDP(conn)
		time.Sleep(time.Millisecond * 50)
	}
}

func defaultIpAndPort() (string, string) {
	ip := ""
	port := ""
	c, err := config.Load("conf/app.conf", config.INI)
	if err == nil {
		// 从配置文件中读取端口号
		mPort, err := c.Get("ai.platform.broadcast.port")
		if err == nil {
			port = mPort
		}
		mIp, err := c.Get("ai.platform.broadcast.ip")
		if len(strings.TrimSpace(mIp)) > 0 {
			ip = mIp
		}
	}
	return ip, port
}

func init() {
	ip, port := defaultIpAndPort()
	if len(port) <= 0 {
		logger.Error("广播服务无法启动，请设置广播端口")
		return
	}
	if len(ip) > 0 {
		go Broadcast(ip, port)
	} else {
		ips := GetBroadcast()
		for _, val := range ips {
			go Broadcast(val, port)
		}
	}
}
