package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Dumper     Dumper     `yaml:"dumper"`
	Influx     Influx     `yaml:"influx"`
	Prometheus Prometheus `yaml:"prometheus"`
}

type Prometheus struct {
	Endpoints      string            `yaml:"endpoints"`
	ScrapeInterval int               `yaml:"scrape_interval"`
	NodeIP         map[string]string `yaml:"node_ip"`
	Metrics        []string          `yaml:"metrics"`
}

type Dumper struct {
	Interval int `yaml:"interval"`
}

type Influx struct {
	Endpoints string `yaml:"endpoints"`
	User      string `yaml:"user"`
	Pwd       string `yaml:"pwd"`
	Database  string `yaml:"database"`
}

func NewConfig(configPath string) *Config {
	var config Config
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("read config failed")
	}
	err = yaml.Unmarshal(yamlFile, &config)
	return &config
}
