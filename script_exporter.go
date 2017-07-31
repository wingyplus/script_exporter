// TODO: load script metrics into prometheus
// TODO: run script
// TODO: add timeout when running script (with context)
// TODO: run script in parallel

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ScriptsConfig []Script `yaml:"scripts_config"`
}

type Script struct {
	Name    string        `yaml:"name"`
	Cmd     string        `yaml:"cmd"`
	Args    []string      `yaml:"args,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

var (
	scriptsConfig = flag.String("scripts_config", "./scripts_config.yml", "scripts configuration path")
)

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(*scriptsConfig)
	if err != nil {
		log.Fatal(err)
	}

	var conf Config
	yaml.Unmarshal(b, &conf)

	fmt.Println(conf)
}
