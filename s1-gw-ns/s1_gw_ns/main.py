'''
s1_gw_ns, generate lorwan packets, publishes them on mqtt
to create load on chirpstack network server.
'''
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

    cli = connect_mqtt(CLIENT_ID, BROKER, PORT)

    while True:
        cli.publish(TOPIC, json.dumps({
            "mhdr": {
                "mType": "",
                "major": "",
            }
        }))        
        # sleeps 10 seconds
        time.sleep(10)
