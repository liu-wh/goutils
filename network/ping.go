package network

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
)

type PingReport struct {
	MaxLatency time.Duration
	MinLatency time.Duration
	AvgLatency time.Duration
	LostNum    int
	LostPer    int
	Packages   []IcmpReply
}

type IcmpReply struct {
	Duration time.Duration
	Seq      int
}

func Ping(host string, count int) (PingReport, error) {
	var (
		report PingReport
		err    error
	)
	dst, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return report, err
	}
	// 创建 ICMP 连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return report, err
	}
	defer conn.Close()
	for i := 0; i < count; i++ {
		// 构造 ICMP 报文
		msg := &icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte("Hola!"),
			},
		}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			continue
		}
		// 发送 ICMP 报文
		start := time.Now()
		_, err = conn.WriteTo(msgBytes, dst)
		if err != nil {
			return report, err
		}
		// 接收 ICMP 报文
		reply := make([]byte, 1500)
		err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			return report, err
		}
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			report.LostNum += 1
			continue
		}
		duration := time.Since(start)
		// 解析 ICMP 报文
		msg, err = icmp.ParseMessage(protocolICMP, reply[:n])
		if err != nil {
			return report, err
		}
		// 打印结果
		switch msg.Type {
		case ipv4.ICMPTypeEchoReply:
			echoReply, ok := msg.Body.(*icmp.Echo)
			if !ok {
				return report, err
			}
			if peer.String() == dst.String() && echoReply.ID == os.Getpid()&0xffff && echoReply.Seq == 1 {
				// fmt.Printf("reply from %s: seq=%d time=%v\n", dst.String(), echoReply.Seq, duration)
				report.AvgLatency += duration
				report.Packages = append(report.Packages, IcmpReply{Duration: duration, Seq: echoReply.Seq})
				if i == 0 {
					report.MinLatency = duration
					report.MaxLatency = duration
					continue
				}
				if duration > report.MaxLatency {
					report.MaxLatency = duration
				} else if duration < report.MinLatency {
					report.MinLatency = duration
				}
			}
		default:
			return report, fmt.Errorf("unexpected ICMP message type: %v", msg.Type)
		}
	}
	report.AvgLatency = report.AvgLatency / time.Duration(count)
	if report.LostNum > 0 {
		report.LostPer = count / report.LostNum
	}
	return report, nil
}
