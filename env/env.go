package env

import (
	"log"
	"os"
	"strconv"
)

const QUEUES_NAME_LENGTH = 32
const MAX_MESSAGE_SIZE = 1024 * 100
const CONSUME_DELAY = 1
const (
	ConnHost    = "localhost"
	ConnPort    = "6552"
	ConnType    = "tcp"
	MESSAGE_GET = "GET"
	MESSAGE_PUT = "PUT"
)

func GetHost() string {
	host := os.Getenv("BHOST")
	if host != "" {
		return host
	} else {
		return ConnHost
	}
}

func GetPort() string {
	port := os.Getenv("BPORT")
	if port != "" {
		return port
	} else {
		return ConnPort
	}
}

func GetMaxMessageSize() int {
	size := os.Getenv("BMAXMESSAGESIZE")
	if size == "" {
		return MAX_MESSAGE_SIZE
	} else {
		value, err := strconv.Atoi(size)
		if err != nil {
			log.Fatalf("wrong message size env value given")
		}
		return value
	}
}
