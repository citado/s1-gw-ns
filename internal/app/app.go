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
	go func() {
		var payload chirpstack.RxMessage

		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			pterm.Error.Printf("cannot unmarshal incomming message %s", err)
		}

		d := time.Since(payload.RxInfo[0].Time)
		pterm.Info.Printf("latency %s\n", d)

		a.Durations = append(a.Durations, d)
		a.signal <- 1

		pterm.Info.Println("packet process done")
	}()
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

	Total int
	Delay time.Duration
}

type Application struct {
	Client    mqtt.Client
	Durations []time.Duration
	Gateways  []lora.Gateway
	signal    chan int

	Total int
	Delay time.Duration
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
	app.Delay = cfg.Delay
	app.Total = cfg.Total
	app.signal = make(chan int)

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

func (a *Application) Run() {
	for _, gateway := range a.Gateways {
		go func(gateway lora.Gateway) {
			for i := 0; i < a.Total; i++ {
				packet, err := gateway.Generate(map[string]interface{}{})
				if err != nil {
					pterm.Fatal.Println(err.Error())
				}

				token := a.Client.Publish(gateway.Topic(), 0, false, packet)
				token.Wait()

				time.Sleep(a.Delay)
			}
		}(gateway)
	}

	for i := 0; i < a.Total*len(a.Gateways); i++ {
		<-a.signal
	}
}
