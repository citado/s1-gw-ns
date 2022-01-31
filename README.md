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
- create _fake-profile_ service profile for _docker-compsoe_ with recommended defaults.
- create _fake-gateway_ gateway with `b827ebffff70c80a` and sane defaults.
- create _fake_dp_ device profile for _docker_compose_ with recommended defaults.

## What we do?

Simulate few gateways that send packets with predefined pattern (normal process or etc.)
and the evaluate one way delay and packet delivery ratio from gateway to network server.
