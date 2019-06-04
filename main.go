package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	cfg     config
	version = "master"
	build   = "0000-00-00 00:00:00"
)

func init() {
	// set default values in config
	cfg.BindIP = "0.0.0.0"
	cfg.BindPort = "8080"
	cfg.APIURI = "/api"
	cfg.Projects = make(map[string]project)

}

func main() {

	// logs debug info
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// parse flags
	cfgFilePath := flag.String("config", "./config.yml", "config file")
	v := flag.Bool("version", false, "print version")
	flag.Parse()

	// handle flags
	if *v == true {
		// version request
		fmt.Printf("version: %s, build: %s\n", version, build)
	} else {
		// read config
		err := readConfig(cfgFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		// debug config print
		fmt.Printf("config: #%v\n", cfg)

		// run rest API
		http.HandleFunc(cfg.APIURI+"/", postHandler)
		err = http.ListenAndServe(cfg.BindIP+":"+cfg.BindPort, nil)
		if err != nil {
			log.Println(err)
		}
	}
}
