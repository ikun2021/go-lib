package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 消息回调函数：当收到消息时执行
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("收到消息:Topic=[%s] Payload=[%s]\n", msg.Topic(), string(msg.Payload()))
}

// 连接成功回调
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("连接成功！")
}

// 连接丢失回调
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("连接丢失: %v\n", err)
}

func main() {
	// =================配置部分=================
	var broker = "tcp://192.168.2.159:1883" // 如果是远程服务器，请填服务器IP

	// ClientID 必须唯一，通常加个随机数或时间戳
	// 如果两个客户端使用相同的 ID，旧的那个会被踢下线
	var clientID = "go_client_" + fmt.Sprintf("%d", time.Now().Unix())

	// 如果你在 EMQX 开启了认证，请填这里；默认匿名访问可留空
	var username = "zhangSan"
	var password = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InpoYW5nU2FuIn0.gNPxkIhlhZ1mjAr543iCRr89nptL6D-Aid35UlEYZn4"
	// =========================================
	// 创建客户端配置
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	// 设置回调函数
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// 开启自动重连 (非常重要)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Second)

	// 创建客户端实例
	client := mqtt.NewClient(opts)

	// 建立连接
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// === 订阅主题 ===
	topic := "sensor/data"
	token := client.Subscribe(topic, 1, nil) // QoS 1
	token.Wait()
	fmt.Printf("已订阅主题: %s\n", topic)

	// === 启动一个协程每秒发布一条消息 ===
	go func() {
		for {
			text := fmt.Sprintf("当前时间: %s", time.Now().Format(time.RFC3339))
			token := client.Publish(topic, 0, false, text)
			token.Wait()
			fmt.Println("已发送消息 ->", text)
			time.Sleep(2 * time.Second)
		}
	}()

	// 阻塞主进程，防止程序退出
	// 监听 Ctrl+C 信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("正在断开连接...")
	client.Disconnect(250)
}
