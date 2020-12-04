package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beyondyyh/exconfig"
)

var ecfg *exconfig.Manifest

func init() {
	// use default config
	config := exconfig.DefaultConfig()

	// use self-defined config
	// config := &exconfig.Config{
	// 	ConsulServerAddr: "http://consul-dev.im.weibo.cn:8500/",
	// 	Datacenter:       "kylin_dev",
	// 	KeyPrefix:        "mp_service/release/manifest",
	// }

	var err error
	ecfg, err = exconfig.New(config)
	if err != nil {
		log.Fatal(err)
		ecfg.Close()
	}

	fmt.Printf("env: %+v\n", ecfg.GenerateEnv())
}

func main() {
	var ticker = time.NewTicker(2 * time.Second)
	var timeout = time.After(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			test()
		case <-timeout:
			ecfg.Close()
			return
		}
	}
}

func test() {
	foo, err := exconfig.String(ecfg.Acquire("foo"))
	log.Printf("foo:%+v, err:%+v\n", foo, err)

	reply, err := ecfg.Acquire("foo")
	foos, err := exconfig.Strings(reply, err, "\n")
	log.Printf("foos:%+v, err:%+v\n", foos, err)

	enable, err := exconfig.Bool(ecfg.Acquire("enable"))
	reply, err = ecfg.Acquire("whitelist")
	users, err := exconfig.Sets(reply, err, "\n")
	log.Printf("enable:%v users:%+v, err:%+v isViper:%t\n", enable, users.Elements(), err, users.Contains("3193013134"))

	type jsonSample struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
		Sex       int    `json:"sex"`
	}
	var jguy jsonSample
	reply, err = ecfg.Acquire("json")
	err = exconfig.Json(reply, err, &jguy)
	log.Printf("guy:%+v err:%+v\n", jguy, err)

	type tomlSample struct {
		Guy struct {
			FirstName string `toml:"firstName"`
			LastName  string `toml:"lastName"`
			Age       int    `toml:"age"`
			Sex       int    `toml:"sex"`
		} `toml:"guy"`
	}
	var tguy tomlSample
	reply, err = ecfg.Acquire("toml")
	err = exconfig.Toml(reply, err, &tguy)
	log.Printf("guy:%+v err:%+v\n", tguy, err)

	type yamlSample struct {
		Stages []string `yaml:"stages"`
	}
	var y yamlSample
	reply, err = ecfg.Acquire("yaml")
	err = exconfig.Yaml(reply, err, &y)
	log.Printf("yaml:%+v err:%+v\n", y, err)
}
