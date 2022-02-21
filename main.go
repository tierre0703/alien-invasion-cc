/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"alien-invasion-cc/cmd"
	log "github.com/sirupsen/logrus"
) 
func main() {
	fmt.Println("=========================")
	fmt.Println("Alien Invasion Simulator")
	fmt.Println("=========================")
	cmd.Execute()
}

func init() {
	log.SetLevel(log.WarnLevel)
}