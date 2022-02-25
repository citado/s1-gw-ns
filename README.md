# Simulation 1: LoRaWAN Gateway -- Chirpstack Evaluation

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/citado/s1-gw-ns/lint?label=lint&logo=github&style=flat-square)

## Introduction

In this simulation we generate ad-hoc packets with [chirpstack](https://www.chirpstack.io/) gateway bridge strcuture
and then forward them into network server to see dalay and packet delivery ratio (PDR) in the stress conditions.
`docker-compose` is written based on [chirpstack-docker](https://github.com/brocaar/chirpstack-docker).

## Step by Step in Chirpstack

- `docker-compose up`
- login to application server on `:8080` with `admin:admin`.
- create _docker-compsoe_ network server with `chirpstack-network-server:8000`.
- create _fake_ **gateway profile** with for _docker-compose_ with recommended defaults.
- create _fake-profile_ **service profile** for _docker-compsoe_ with recommended defaults (remember to have gateway metadata).
- create _fake-gateway_ **gateway** with `b827ebffff70c80a` and sane defaults.
- create _fake_dp_ device profile for _docker_compose_ with recommended defaults.
- create _citado_ application with _fake-profile_.

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
gateways:
  - mac: "b827ebffff70c80a"
    keys:
      network_skey: "DB56B6C3002A4763A79E64573C629D97"
      application_skey: "94B49CD7BC621BC46571D019640804AA"
    devices:
      # generated by gen command.
      - addr: a80fb59d
        dev_eui: a80fb59d09cfe066
      - addr: 09cfe066
        dev_eui: 556414c55c87df12
      - addr: 556414c5
        dev_eui: 4c9814eab0894b9c
```

## Generate Devices

We need more than one device for our simulations so we can generate them with:

```sh
# generate 100 devices
./s1-gw-ns gen --count 100
```

then you can use the result devices into the main configuration.

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
nats stream add --subjects gateway.+.event.up
```

## Results

Results are stored in csv files, each row started with device id and its receiving delays.
For each try we have separeted file.
