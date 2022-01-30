package main

import (
	"fmt"
	"time"

	"github.com/citado/s1-gw-ns/internal/lora"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pterm/pterm"
)

func connectHandler(client mqtt.Client) {
	pterm.Info.Println("connected")
}

func connectLostHandler(client mqtt.Client, err error) {
	pterm.Error.Printf("connection lost due to %s", err)
}

func main() {
	var (
		broker = "127.0.0.1"
		port   = 1883
	)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("fake_gateway")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		pterm.Fatal.Println(token.Error())
	}

	gw := lora.New(lora.Config{
		MAC: "b827ebffff70c80a",
		Keys: lora.Keys{
			NetworkSKey:     "DB56B6C3002A4763A79E64573C629D97",
			ApplicationSKey: "94B49CD7BC621BC46571D019640804AA",
		},
		Device: lora.Device{
			Addr: "26011CF6",
		},
	})

	p, err := gw.Generate(map[string]interface{}{
		"100": 6750,
		"101": 6606,
	})
	if err != nil {
		pterm.Fatal.Println(err.Error())
	}

	token := client.Publish(gw.Topic(), 0, false, p)
	token.Wait()
	time.Sleep(time.Second)
}
