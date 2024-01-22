package main

import (
	"analysis_ini/analysis"
	"fmt"
)

type Config struct {
	ZookeeperClusterAddress []string `ini:"zookeeper_cluster_address"`
	InsecurePort            uint     `ini:"insecure_port"`
	RootDirectory           string   `ini:"root_directory"`
}

func main() {
	config := &Config{}
	iniFilePath := "../config/config.ini"
	err := analysis.UnMarshalWithIniPath(iniFilePath, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", config)
}
