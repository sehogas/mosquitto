package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Recibido [%s] en tópico [%s]\n", msg.Payload(), msg.Topic())
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Conectado!")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Conexión perdida: %s\n", err.Error())
}

/*
func newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("certificado.crt")
	if err != nil {
		panic(err.Error())
	}

	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{RootCAs: certpool}
}
*/

func main() {
	var broker = "ws://localhost:9001"
	//var broker = "tcp://localhost:1883"
	//var broker = "ssl://test.mosquitto.org:8883"

	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("ejemplo_go_mqtt")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler
	//options.SetTLSConfig(newTLSConfig())

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "/eventos"

	token = client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Suscripto a tópico [%s]\n", topic)

	num := 1000
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("mensaje %d", i)
		fmt.Printf("Publicando [%s] en tópico [%s]\n", text, topic)
		token = client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)

	}

	client.Disconnect(100)
}
