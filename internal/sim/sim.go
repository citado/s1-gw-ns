package sim

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/citado/s1-gw-ns/internal/lora"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Gateways []lora.Config `yaml:"gateways"`
}

func Read() Config {
	var instance Config

	f, err := os.Open("sim.yaml")
	if err != nil {
		log.Fatalf("cannot open the simulation configuration file %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("cannot read the simulation configuration file %s", err)
	}

	if err := yaml.Unmarshal(b, &instance); err != nil {
		log.Fatalf("cannot parse the simulation configuration file %s", err)
	}

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling config: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Simulation ================
	%s
	======================================================
	`
	log.Printf(tmpl, string(indent))

	return instance
}

func Write(cfg Config) {
	f, err := os.Create("sim.yaml")
	if err != nil {
		log.Fatalf("cannot open the simulation configuration file %s", err)
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("cannot parse the simulation configuration file %s", err)
	}

	if _, err := f.Write(b); err != nil {
		log.Fatalf("cannot write into simulation configuration file %s", err)
	}

	f.Close()
}
