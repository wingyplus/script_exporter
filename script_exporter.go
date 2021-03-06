// TODO: add timeout when running script (with context)
// TODO: run script in parallel

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ScriptsConfig []Script `yaml:"scripts_config"`
}

type Script struct {
	Name    string        `yaml:"name"`
	Cmd     string        `yaml:"cmd"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

var (
	scriptsConfig = flag.String("scripts_config", "./scripts_config.yml", "scripts configuration path")
)

func probeHandler(w http.ResponseWriter, r *http.Request, scripts []Script) {
	probeSuccessGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "script_success",
		Help: "Script success",
	}, []string{"script"})

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)

	for _, script := range scripts {
		shCmd := exec.Command("sh")
		stdin, err := shCmd.StdinPipe()
		if err != nil {
			fmt.Println(err.Error())
		}

		stdin.Write([]byte(script.Cmd))
		stdin.Close()

		gauge := probeSuccessGauge.With(prometheus.Labels{"script": script.Name})

		if err := shCmd.Start(); err != nil {
			fmt.Println("sh start error", err.Error())
		}
		if err := shCmd.Wait(); err != nil {
			fmt.Println("sh wait error", err.Error())
			gauge.Set(0)
		} else {
			gauge.Set(1.0)
		}
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(*scriptsConfig)
	if err != nil {
		log.Fatal(err)
	}

	var conf Config
	yaml.Unmarshal(b, &conf)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, conf.ScriptsConfig)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
