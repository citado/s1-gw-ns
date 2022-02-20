# Simulation 1: LoRaWAN Gateway -- Network Server

## Introduction

In this simulation we generate ad-hoc packets with [chirpstack](https://www.chirpstack.io/) gateway bridge strcuture
and then forward them into network server to see dalay and packet delivery ratio (PDR) in the stress conditions.
`docker-compose` is written based on [chirpstack-docker](https://github.com/brocaar/chirpstack-docker).

## Step by Step

- `docker-compose up`
- login to application server on `:8080` with `admin:admin`.
- create _docker-compsoe_ network server with `chirpstack-network-server:8000`.
- create _fake_ gateway profile with for _docker-compose_ with recommended defaults.
- create _fake-profile_ service profile for _docker-compsoe_ with recommended defaults (remember to have gateway metadata).
- create _fake-gateway_ gateway with `b827ebffff70c80a` and sane defaults.
- create _fake_dp_ device profile for _docker_compose_ with recommended defaults.
- create _citado_ application with _fake-profile_.
- create _fake-device_ in _citado_ application with random DevEUI.
- fill _fake-device_ activation fields.

## What we do?

Simulate few gateways that send packets with predefined pattern (normal process or etc.)
and the evaluate one way delay and packet delivery ratio from gateway to network server.

## Results

If we use 1s between each message then we don't loss any message but with 100ms or less we will loss messages
and we don't receive them at application layer.
