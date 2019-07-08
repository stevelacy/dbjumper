package cli

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	dbjumper "github.com/stevelacy/dbjumper/pkg"
	"github.com/stevelacy/dbjumper/proxy"
)

var config dbjumper.Config
var configPath = "./dbjumper.yaml"

// Init creates the cli
func Init() error {

	envPath := os.Getenv("config_file")
	if envPath != "" {
		configPath = envPath
	}
	cfg, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(string(cfg)), &config)
	if err != nil {
		return err
	}

	log.Printf("starting on %s", config.Address)
	err = proxy.Start(&config)
	if err != nil {
		log.Fatal(err)
	}
	go forever()
	select {}
}

func forever() {
	for {
		time.Sleep(1)
	}
}
