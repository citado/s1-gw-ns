package lora

// nolint: staticcheck
import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/lorawan"
	"github.com/fxamacker/cbor/v2"
	"github.com/golang/protobuf/proto"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// FPort distinguishes the data format.
	FPort uint8 = 5
	// FCnt counts number of frames.
	FCnt = 10
	// Battery level of device.
	Battery = 115
	// Margin is a link margin of device.
	Margin = 7
	// Channel for communication (Gateway).
	Channel = 1
	// CRCStatus of packet (Gateway).
	CRCStatus = 1
	// Bandwidth of communication (Gateway).
	Bandwidth = 125
	// SpreadFactor of communication (Gateway).
	SpreadFactor = 7
	// Frequency of communication (Gateway).
	Frequency = 868300000
	// LoRaSNR indicates Signal to noise ratio (Gateway).
	LoRaSNR = 7
	// RFChain of communication (Gateway).
	RFChain = 1
	// Size of packet (Gateway).
	Size = 23
)

// Gateway simulate data based on LoRaWAN protocol from an ABP device and one gateway.
// It encrypts data and you will need
// a lora server to decode it.
type Gateway struct {
	Config
}

// create new gateway simulation based on given configuration.
func New(cfg Config) Gateway {
	return Gateway{
		Config: cfg,
	}
}

// Topic returns lora gateway mqtt topic.
func (g Gateway) Topic() string {
	return fmt.Sprintf("gateway/%s/event/up", g.MAC)
}

// Generate generates lora message by converting input into cbor and encrypts it.
// nolint: funlen
func (g Gateway) Generate(message interface{}, index int) ([]byte, error) {
	b, err := cbor.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("cannot encode message to cbor: %w", err)
	}

	pterm.Info.Printf("message payload length %d\n", len(b))

	mac, err := hex.DecodeString(g.MAC)
	if err != nil {
		return nil, fmt.Errorf("cannot decode gateway mac: %w", err)
	}

	// converts network and application session keys to AES128
	appSKeySlice, err := hex.DecodeString(g.Keys.ApplicationSKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decode application session key: %w", err)
	}

	var appSKey lorawan.AES128Key

	copy(appSKey[:], appSKeySlice)

	nwkSKeySlice, err := hex.DecodeString(g.Keys.NetworkSKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decode network session key: %w", err)
	}

	var nwkSKey lorawan.AES128Key

	copy(nwkSKey[:], nwkSKeySlice)

	// converts device addr into DevAddr
	devAddrSlice, err := hex.DecodeString(g.Devices[index].Addr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode device adr: %w", err)
	}

	var devAddr lorawan.DevAddr

	copy(devAddr[:], devAddrSlice)

	// https://godoc.org/github.com/brocaar/lorawan#example-PHYPayload--Lorawan10Encode
	fport := FPort

	phy := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.UnconfirmedDataUp,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: &lorawan.MACPayload{
			FHDR: lorawan.FHDR{
				DevAddr: devAddr,
				FCtrl: lorawan.FCtrl{
					ADR:       false,
					ADRACKReq: false,
					ACK:       false,
					FPending:  false,
					ClassB:    false,
				},
				FCnt: FCnt,
				FOpts: []lorawan.Payload{
					&lorawan.MACCommand{
						CID: lorawan.DevStatusAns,
						Payload: &lorawan.DevStatusAnsPayload{
							Battery: Battery,
							Margin:  Margin,
						},
					},
				},
			},
			FPort:      &fport,
			FRMPayload: []lorawan.Payload{&lorawan.DataPayload{Bytes: b}},
		},
		MIC: lorawan.MIC{},
	}

	if err := phy.EncryptFRMPayload(appSKey); err != nil {
		return nil, fmt.Errorf("frame encoding failed: %w", err)
	}

	if err := phy.SetUplinkDataMIC(lorawan.LoRaWAN1_0, 0, 0, 0, nwkSKey, lorawan.AES128Key{}); err != nil {
		return nil, fmt.Errorf("frame mic calculation failed: %w", err)
	}

	phyBytes, err := phy.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal frame to binary: %w", err)
	}

	// lora gateway grpc message based on
	// github.com/brocaar/chirpstack-api/go/v3

	// nolint: exhaustivestruct, gomnd
	raw, err := proto.Marshal(&gw.UplinkFrame{
		TxInfo: &gw.UplinkTXInfo{
			Frequency:  Frequency,
			Modulation: common.Modulation_LORA,
			ModulationInfo: &gw.UplinkTXInfo_LoraModulationInfo{
				LoraModulationInfo: &gw.LoRaModulationInfo{
					Bandwidth:             Bandwidth,
					SpreadingFactor:       SpreadFactor,
					CodeRate:              "4/5",
					PolarizationInversion: false,
				},
			},
		},
		RxInfo: &gw.UplinkRXInfo{
			GatewayId:         mac,
			Time:              timestamppb.New(time.Now()),
			TimeSinceGpsEpoch: nil,
			Rssi:              0,
			LoraSnr:           0.0,
			Channel:           0,
			RfChain:           0,
			Board:             0,
			Antenna:           0,
			Location: &common.Location{
				Latitude:  35.723737,
				Longitude: 50.952981,
				Altitude:  0.0,
				Source:    0,
				Accuracy:  0,
			},
			FineTimestampType: 0,
			FineTimestamp:     nil,
			Context:           []byte{},
			UplinkId:          []byte{},
			CrcStatus:         gw.CRCStatus_CRC_OK,
		},
		PhyPayload: phyBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot marshal protobuf message from gateway: %w", err)
	}

	return raw, nil
}
