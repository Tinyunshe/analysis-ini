package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("../config/123")
	if err != nil {
		fmt.Println(err)
		panic("erro")
	}
	defer file.Close()
	// reader := bufio.NewReader(file)
	// for {
	// 	line, err := reader.ReadString('\n')
	// 	if err == io.EOF {
	// 		return
	// 	}
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	line = strings.TrimRight(line, "\n")
	// 	fmt.Print(line)
	// }
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
