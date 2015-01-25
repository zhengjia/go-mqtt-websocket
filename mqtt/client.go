package mqtt

import (
	"code.google.com/p/go.net/websocket"
	"crypto/rand"
	"encoding/json"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"io"
	"log"
	"strings"
	"time"
)

type Proxy struct {
	Conn   *websocket.Conn
	Client *MQTT.MqttClient
	Done   chan bool
}

// TODO add `auth_code`
type Request struct {
	Action  string
	Topic   string
	Message string
}

type Response struct {
	Status        int
	StatusMessage string
	Topic         string
	Message       string
}

func GetClient() (c *MQTT.MqttClient, err error) {
	// tcp://test.mosquitto.org:1883
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientId(getRandStr())
	opts.SetCleanSession(true)
	opts.SetOnConnectionLost(onConnectionLost)
	c = MQTT.NewClient(opts)
	_, err = c.Start()
	return
}

func (proxy *Proxy) Start() {
	for {
		buf := make([]byte, 1024)
		l, err := proxy.Conn.Read(buf)
		time.Sleep(2 * time.Second)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
			}
		}
		if l > 0 {
			request, _ := proxy.parseRequest(string(buf[:l]))
			proxy.processRequest(request)
		}
	}
	log.Println("Closing connection")
	proxy.Close()
}

func (proxy *Proxy) Close() {
	proxy.Client.ForceDisconnect()
	proxy.Done <- true
	return
}

func (proxy *Proxy) Subscribe(topic string) {
	var topicFilter, _ = MQTT.NewTopicFilter(topic, 0)
	if receipt, err := proxy.Client.StartSubscription(proxy.onMessageReceived, topicFilter); err != nil {
		log.Println(err)
		// TODO Notify client the failure
	} else {
		<-receipt
		log.Println("Subscribed Topic: ", topic)
	}
}

func (proxy *Proxy) Publish(topic string, message string) {
	receipt := proxy.Client.Publish(MQTT.QOS_ONE, topic, message)
	<-receipt
	log.Println("Published Topic: ", topic)
}

func (proxy *Proxy) EndSubscription(topic string) error {
	receipt, err := proxy.Client.EndSubscription(topic)
	if err != nil {
		log.Println(err)
	} else {
		<-receipt
	}
	return err
}

func (proxy *Proxy) parseRequest(raw string) (request *Request, err error) {
	request = new(Request)
	err = json.NewDecoder(strings.NewReader(raw)).Decode(request)
	// TODO handle returned error
	if err != nil {
		log.Println("error:", err)
	}
	// TODO validate request
	return
}

func (proxy *Proxy) processRequest(request *Request) {
	switch request.Action {
	case "subscribe":
		proxy.Subscribe(request.Topic)
	case "publish":
		proxy.Publish(request.Topic, request.Message)
	}
}

func (proxy *Proxy) onMessageReceived(client *MQTT.MqttClient, message MQTT.Message) {
	var response = &Response{Topic: message.Topic(), Message: string(message.Payload())}
	jsonResponse, _ := json.Marshal(response)
	proxy.Conn.Write(jsonResponse)
}

func onConnectionLost(client *MQTT.MqttClient, reason error) {
	log.Println("Lost connection with server", reason)
	// TODO Pass in conn and close it
}

func getRandStr() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length := 6
	randStr := make([]byte, length)
	rand.Read(randStr)
	for i, b := range randStr {
		randStr[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(randStr)
}
