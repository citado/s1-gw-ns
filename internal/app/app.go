package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/citado/s1-gw-ns/internal/chirpstack"
	"github.com/citado/s1-gw-ns/internal/lora"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pterm/pterm"
)

func (a *Application) onMessage(client mqtt.Client, msg mqtt.Message) {
	var payload chirpstack.RxMessage

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		pterm.Error.Printf("cannot unmarshal incomming message %s", err)
	}

	d := time.Since(payload.RxInfo[0].Time)
	pterm.Info.Printf("latency %s\n", d)

	a.Durations = append(a.Durations, d)
}

func (a *Application) onConnect(client mqtt.Client) {
	pterm.Info.Println("connected")

	client.Subscribe("application/+/device/+/event/up", 0, a.onMessage)
}

func (a *Application) onDisconnect(client mqtt.Client, err error) {
	pterm.Error.Printf("connection lost due to %s", err)
}

type Config struct {
	Port int
	Addr string
}

type Application struct {
	Client    mqtt.Client
	Durations []time.Duration
	Gateways  []lora.Gateway
}

func New(cfg Config) *Application {
	app := new(Application)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Addr, cfg.Port))
	opts.SetClientID("fake_gateway")
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(app.onConnect)
	opts.SetConnectionLostHandler(app.onDisconnect)
	client := mqtt.NewClient(opts)

	app.Client = client

	return app
}

func (a *Application) Connect() {
	if token := a.Client.Connect(); token.Wait() && token.Error() != nil {
		pterm.Fatal.Println(token.Error())
	}
}

func (a *Application) Gateway(cfg lora.Config) {
	a.Gateways = append(a.Gateways, lora.New(cfg))
}
