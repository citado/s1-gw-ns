package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/citado/s1-gw-ns/internal/chirpstack"
	"github.com/citado/s1-gw-ns/internal/lora"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fxamacker/cbor/v2"
	"github.com/pterm/pterm"
)

const (
	DefaultMessageTimeout = 10 * time.Second
	IDLen                 = 16
)

func (a *Application) onMessage(client mqtt.Client, msg mqtt.Message) {
	go func(msg mqtt.Message) {
		var payload chirpstack.RxMessage

		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			pterm.Error.Printf("cannot unmarshal incomming message %s\n", err)

			return
		}

		data := make(map[string]interface{})

		if err := cbor.Unmarshal(payload.Data, &data); err != nil {
			pterm.Error.Printf("cannot unmarshal incomming cbor message %s\n", err)

			return
		}

		id, ok := data["id"].(uint64)
		if !ok {
			pterm.Error.Printf("cannot convert id to int\n")

			return
		}

		pterm.Info.Printf("id %d\n", id)

		d := time.Since(payload.RxInfo[0].Time)
		pterm.Info.Printf("latency %s\n", d)

		a.signal <- Message{
			Delay: d,
			ID:    id,
		}

		pterm.Info.Println("packet process done")
	}(msg)
}

func (a *Application) onConnect(client mqtt.Client) {
	pterm.Info.Println("connected")

	client.Subscribe("application/+/device/+/event/up", 1, a.onMessage)
}

func (a *Application) onDisconnect(client mqtt.Client, err error) {
	pterm.Error.Printf("connection lost due to %s", err)
}

type Config struct {
	Port int    `koanf:"port"`
	Addr string `koanf:"addr"`

	Total int           `koanf:"total"`
	Delay time.Duration `koanf:"delay"`
}

type Message struct {
	ID    uint64
	Delay time.Duration
}

type Application struct {
	Client    mqtt.Client
	Durations []time.Duration
	Gateways  []lora.Gateway
	signal    chan Message

	Total int
	Delay time.Duration
}

func New(cfg Config) *Application {
	app := new(Application)

	id := make([]byte, IDLen)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Addr, cfg.Port))
	opts.SetClientID(fmt.Sprintf("fake_gateway_%s", hex.EncodeToString(id)))
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(app.onConnect)
	opts.SetConnectionLostHandler(app.onDisconnect)
	client := mqtt.NewClient(opts)

	app.Client = client
	app.Delay = cfg.Delay
	app.Total = cfg.Total

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

func (a *Application) PublishSubscribe() {
	a.Durations = nil
	a.signal = make(chan Message)

	for _, gateway := range a.Gateways {
		go func(gateway lora.Gateway) {
			for i := 0; i < a.Total; i++ {
				// generate empty packet
				packet, err := gateway.Generate(map[string]interface{}{
					"id": i,
				})
				if err != nil {
					pterm.Fatal.Println(err.Error())
				}

				token := a.Client.Publish(gateway.Topic(), 1, false, packet)
				if token.Wait() && token.Error() != nil {
					pterm.Fatal.Println(token.Error())
				}

				pterm.Info.Printf("message [%d] is sent over mqtt\n", i)
				time.Sleep(a.Delay)
			}
		}(gateway)
	}

	// wait for all messages to be receieved from application server.
	for i := 0; i < a.Total; i++ {
		select {
		case m := <-a.signal:
			a.Durations = append(a.Durations, m.Delay)
		case <-time.After(DefaultMessageTimeout):
			pterm.Error.Printf("missed event\n")

			a.Durations = append(a.Durations, -1)
		}
	}
}

func (a *Application) Publish() {
	a.Durations = nil
	a.signal = make(chan Message)

	var wg sync.WaitGroup

	wg.Add(len(a.Gateways))

	for _, gateway := range a.Gateways {
		go func(gateway lora.Gateway) {
			for i := 0; i < a.Total; i++ {
				// generate empty packet
				packet, err := gateway.Generate(map[string]interface{}{
					"id": i,
				})
				if err != nil {
					pterm.Fatal.Println(err.Error())
				}

				token := a.Client.Publish(gateway.Topic(), 1, false, packet)
				if token.Wait() && token.Error() != nil {
					pterm.Fatal.Println(token.Error())
				}

				pterm.Info.Printf("message [%d] is sent over mqtt\n", i)
				time.Sleep(a.Delay)
			}
			wg.Done()
		}(gateway)
	}

	wg.Wait()
}

func (a *Application) Subscribe() {
	a.Durations = nil
	a.signal = make(chan Message)

	// wait for all messages to be receieved from application server.
	for i := 0; i < a.Total; i++ {
		select {
		case m := <-a.signal:
			a.Durations = append(a.Durations, m.Delay)
		case <-time.After(DefaultMessageTimeout):
			pterm.Error.Printf("missed event\n")

			a.Durations = append(a.Durations, -1)
		}
	}
}
