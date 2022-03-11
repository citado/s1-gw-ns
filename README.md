# Simulation 1: LoRaWAN Gateway -- Chirpstack Evaluation

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/citado/s1-gw-ns/lint?label=lint&logo=github&style=flat-square)

## Introduction

In this simulation we generate ad-hoc packets with [chirpstack](https://www.chirpstack.io/) gateway bridge strcuture
and then forward them into network server to see dalay and packet delivery ratio (PDR) in the stress conditions.
`docker-compose` is written based on [chirpstack-docker](https://github.com/brocaar/chirpstack-docker).

## What we do?

Simulate few gateways that send packets with predefined pattern (normal process or etc.)
and the evaluate one way delay and packet delivery ratio from gateway to network server
to application server and then into our simulator again.

## Main Configuration

There is a yaml configuration file which is used in **all** and **gen** subcommands.

```yaml
---
tries: 10
app:
  total: 1000
  delay: 1s
```

this configuration contains everything except simlated devices. you need to
generate them with **gen** command.

## Generate Devices

We always need more than one device for our simulations so we can generate them with:

```sh
# generate 10 devices with 10 gateway. randomly assigns these devices to gateway.
# each device will connect to one, two or three gateways.
./s1-gw-ns gen --dev-count 10 --gw-count 10
```

the result is generated into `sim.yaml` as follow:

```yaml
gateways:
  - mac: a80fb59d09cfe066
    keys:
      networkskey: DB56B6C3002A4763A79E64573C629D97
      applicationskey: 94B49CD7BC621BC46571D019640804AA
    devices:
      - addr: a80fb59d
        deveui: a80fb59d09cfe066
```

## Ignite the Rocket

Everything is ready? are you ready for having fun with chirpstack and put it on fire?
then ignite the rocket.

```sh
./s1-gw-ns all
```

## Jetstream

for not losing messages between in mqtt broker we are using nats jetstream with mqtt interface
which is way more better than emqx.

```sh
# create chirpstack stream (pay attention to the subject)
nats stream add --subjects gateway.*.event.up
```

please note that by using jetstream's streams you may have more than 1 result for a message.

## Results

Results are stored in csv files, each row started with device id and its receiving delays.
For each try we have separeted file.
