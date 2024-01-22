package main

import (
	"analysis_ini/analysis"
	"fmt"
)

func main() {
	iniFilePath := "../config/config.ini"
	config, err := analysis.UnMarshalWithIniPath(iniFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", config)
}
