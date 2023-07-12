# Go Utils

## Network Util
### Install
`go get github.com/liu-wh/goutils/network`
### Ping - Get network latency from local host to remote host with icmp protocol, not exec.run
```
func main() {
	host := "www.baidu.com"
	num := 10
	report, err := network.Ping(host, num)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Lost=%d(%d%% loss)\nMinimum = %v, Maximum = %v, Average = %v", report.LostNum, report.LostPer, report.MinLatency, report.MaxLatency, report.AvgLatency)
}
```
![image](https://github.com/liu-wh/goutils/assets/52809998/db86c42b-53e0-447e-adc5-f4c87d929824)

## Convert Util
### Install
`go get github.com/liu-wh/goutils/convert`
### Str2bytes, Bytes2str - Use pointers to do string and byte conversion, which is more efficient
### Bytes2Human - Convert the byte int to k, m, g string