package main

import (
	"log"

	incus "github.com/lxc/incus/client"
)

func main() {
	log.Println("Minecraft Server Panel v1.0 starting . . .")
	client, err := incus.ConnectIncusUnix("", nil)
	if err != nil {
		log.Fatalf("connect: %s\n", err)
	}
	serverInfo, _, err := client.GetServer()
	if err != nil {
		log.Fatalf("getServer: %s\n", err)
	}
	if serverInfo.APIVersion != "1.0" {
		log.Fatalln("Wrong Incus server API version, refusing to connect for now!")
	}

}
