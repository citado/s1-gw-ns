package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
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
	KeepAliveTimeout      = 100 * time.Millisecond
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

		dev, ok := data["device"].(string)
		if !ok {
			pterm.Error.Printf("cannot convert device to string\n")

			return
		}

		pterm.Info.Printf("id %d device %s\n", id, dev)

		d := time.Since(payload.RxInfo[0].Time)
		pterm.Info.Printf("latency %s\n", d)

		a.signal <- Message{
			Delay:  d,
			ID:     id,
			Device: dev,
		}

		pterm.Info.Println("packet process done")
	}(msg)
}

func (a *Application) onConnect(client mqtt.Client) {
	pterm.Info.Println("connected")

	if token := client.Subscribe("application/+/device/+/event/up", 1, a.onMessage); token.Wait() && token.Error() != nil {
		pterm.Fatal.Printf("cannot subscribe %s\n", token.Error())
	}
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
	ID     uint64
	Delay  time.Duration
	Device string
}

type Application struct {
	Port      int
	Addr      string
	Client    mqtt.Client
	Durations map[string][]time.Duration
	Gateways  []lora.Gateway
	signal    chan Message

	Total int
	Delay time.Duration
}

func newClientOptions(addr string, port int) *mqtt.ClientOptions {
	id := make([]byte, IDLen)
	if _, err := rand.Read(id); err != nil { // nolint: gosec
		panic(err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", addr, port))
	opts.SetClientID(fmt.Sprintf("fake_gateway_%s", hex.EncodeToString(id)))
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	opts.SetConnectRetry(true)
	opts.SetOrderMatters(false)
	opts.SetKeepAlive(KeepAliveTimeout)

	return opts
}

func New(cfg Config) *Application {
	app := new(Application)

	opts := newClientOptions(cfg.Addr, cfg.Port)
	opts.SetOnConnectHandler(app.onConnect)
	opts.SetConnectionLostHandler(app.onDisconnect)
	client := mqtt.NewClient(opts)

	app.Client = client
	app.Delay = cfg.Delay
	app.Total = cfg.Total
	app.Port = cfg.Port
	app.Addr = cfg.Addr

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

func (a *Application) publishOnGateway(gateway lora.Gateway) {
	var wg sync.WaitGroup

	wg.Add(len(gateway.Devices))

	// loops over devices and create go-routine for every single of them.
	for j := 0; j < len(gateway.Devices); j++ {
		go func(j int) {
			client := mqtt.NewClient(newClientOptions(a.Addr, a.Port))

			if token := client.Connect(); token.Wait() && token.Error() != nil {
				pterm.Fatal.Println(token.Error())
			}

			for i := 0; i < a.Total; i++ {
				// generates a packet with sequence number and device id.
				packet, err := gateway.Generate(map[string]interface{}{
					"id":     i,
					"device": gateway.Devices[j].Addr,
				}, j)
				if err != nil {
					pterm.Fatal.Println(err.Error())
				}

				token := client.Publish(gateway.Topic(), 1, false, packet)
				<-token.Done()

				if token.Error() != nil {
					pterm.Fatal.Println(token.Error())
				}

				pterm.Info.Printf("message [%d] is sent over mqtt from device [%d]\n", i, j)

				<-time.After(a.Delay)
			}

			pterm.Success.Printf("publishing on device %d is completed\n", j)
			wg.Done()
			client.Disconnect(1)
		}(j)
	}

	wg.Wait()
	pterm.Success.Printf("publishing on gateway %s is completed\n", gateway.MAC)
}

func (a *Application) PublishSubscribe() {
	a.Durations = make(map[string][]time.Duration)
	a.signal = make(chan Message)

	var wg sync.WaitGroup

	wg.Add(len(a.Gateways))

	for _, gateway := range a.Gateways {
		go func(gateway lora.Gateway) {
			a.publishOnGateway(gateway)
			wg.Done()
		}(gateway)
	}

	wg.Wait()

	// wait for all messages to be receieved from application server.
	for i := 0; i < len(a.Gateways); i++ {
		for j := 0; j < len(a.Gateways[i].Devices); j++ {
			for k := 0; k < a.Total; k++ {
				select {
				case m := <-a.signal:
					a.Durations[m.Device] = append(a.Durations[m.Device], m.Delay)
				default:
					pterm.Error.Printf("missed event, we will try on next message\n")
				}
			}
		}
	}
}
