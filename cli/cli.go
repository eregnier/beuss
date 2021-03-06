package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	b "github.com/eregnier/beuss/client"
	env "github.com/eregnier/beuss/env"
)

func main() {
	if len(os.Args) < 3 || (os.Args[1] != "GET" && os.Args[1] != "PUT" && os.Args[1] != "ON") {
		fmt.Println(`
Usage :
  PUT : echo "content" | ./beuss PUT <queueName>
  GET : ./beuss GET <queueName> <output.ext>
  ON : ./beuss ON <queueName>
		`)
		os.Exit(0)
	}

	if os.Args[1] == "PUT" {

		connPUT, err := b.NewClient(env.MESSAGE_PUT)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln("could not read from stdin")
		}
		b.ClientPutMessage(connPUT, os.Args[2], bytes)
		b.ClientClose(connPUT)
	}
	if os.Args[1] == "GET" {
		connGET, err := b.NewClient(env.MESSAGE_GET)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		message, err := b.ClientGetMessage(connGET, os.Args[2])
		if err != nil {
			log.Println("could not get message", err)
		} else {
			if len(os.Args) == 4 {
				if err := os.WriteFile(os.Args[3], message, 0666); err != nil {
					fmt.Println("error writing file", err)
				}
			} else {
				fmt.Print(string(message))
			}
		}
		b.ClientClose(connGET)
	}
	if os.Args[1] == "ON" {
		connGET, err := b.NewClient(env.MESSAGE_GET)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		b.ClientOnMessage(connGET, os.Args[2], func(message []byte) {
			fmt.Print(string(message))
		})
		defer b.ClientClose(connGET)
	}
}
