'''
s1_gw_ns, generate lorwan packets, publishes them on mqtt
to create load on chirpstack network server.
'''
import base64
import json
import random
import time

from paho.mqtt import client


def connect_mqtt(client_id: str, broker: str, port: int) -> client.Client:
    '''
    connect_mqtt creates a connection to mqtt broker.
    it blocks until the successful connection.
    '''

    def on_connect(_client, _userdata, _flags, return_code):
        if return_code == 0:
            print("connected to mqtt broker!")
        else:
            print(f"failed to connect, return code {return_code}")

    # Set Connecting Client ID
    cli = client.Client(client_id)
    cli.on_connect = on_connect
    cli.connect(broker, port)
    cli.loop_start()
    return cli


if __name__ == "__main__":
    print("Welcome to our simulator")

    # lorawan gateway mac address
    # from the isrc project.
    GATEWAY_MAC = 'b827ebffff70c80a'
    # emqx broker address runs with docker-compose.
    BROKER = '127.0.0.1'
    # emqx broker port runs with docker-compose.
    PORT = 1883
    # chirpstack mqtt topic
    TOPIC = f'gateway/{GATEWAY_MAC}/event/up'
    # mqtt client id
    CLIENT_ID = f's1-gw-ns-{random.randint(0, 1000)}'
    # counts number of frames
    FRAME_COUNT = 10
    DEV_ADDR = "26011CF6"

    cli = connect_mqtt(CLIENT_ID, BROKER, PORT)

    while True:
        cli.publish(TOPIC, json.dumps({
            "rxinfo": {
                "board": 0,
                "antena": 0,
                "channel": 1,
                "code_rate": "4/5",
                "crc_status": 1,
            },
            "phyPayload": {
                "mhdr": {
                    "mType": 2,  # UnconfirmedDataUp
                    "major": 0,  # LoRaWANR1
                },
                "macPayload": {
                    "fhdr": {
                        "devAddr": DEV_ADDR,
                        "fCtrl": {
                            "adr": False,
                            "adrAckReq": False,
                            "ack": False,
                            "classB": False,
                        },
                        "fCnt": FRAME_COUNT,
                        "fOpts": [],
                    },
                    "fPort": 5,
                    "frmPayload": [
                        {
                            "bytes": base64.b64encode(b'Hello World').decode(),
                        }
                    ],
                }
            }
        }))
        # sleeps 10 seconds
        time.sleep(10)
