package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
)

const QUEUES_NAME_LENGTH = 32
const MAX_MESSAGE_SIZE = 1024 * 100
const (
	connHost    = "localhost"
	connPort    = "6552"
	connType    = "tcp"
	MESSAGE_GET = "GET"
	MESSAGE_PUT = "PUT"
)

type Queue struct {
	Messages [][]byte
	Name     string
}

type InputMessage struct {
	Queue   string
	content []byte
}

var queues []Queue

func main() {
	queues = make([]Queue, 0)
	listenNetwork()
}

func onMessage(message []byte) error {
	if len(message) < QUEUES_NAME_LENGTH {
		return errors.New("message too short : rejected")
	}
	queueName := strings.TrimSpace(string(message[:QUEUES_NAME_LENGTH]))
	content := message[QUEUES_NAME_LENGTH:]
	inputMessage := InputMessage{queueName, content}
	found := false
	for x, q := range queues {
		if q.Name == queueName {
			queues[x].Messages = append(q.Messages, inputMessage.content)
			found = true
		}
	}
	if !found {
		m := [][]byte{inputMessage.content}
		queues = append(queues, Queue{m, queueName})
	}
	return nil
}

func readMessage(queueName string) ([]byte, error) {
	err := errors.New("no message in given queue")
	for x, q := range queues {
		if q.Name == queueName {
			l := len(q.Messages)
			if l == 0 {
				return nil, err
			}
			message := q.Messages[0]
			queues[x].Messages = queues[x].Messages[1:]
			return message, nil
		}
	}
	return nil, err
}

func handleIncomingConnection(c net.Conn) {
	messageType := ""
	for {
		log.Println("Client " + c.RemoteAddr().String() + " connected.")
		recvData := make([]byte, MAX_MESSAGE_SIZE)
		n, err := c.Read(recvData)
		if err != nil {
			log.Println("network error on data receive", err)
		}
		if n == 4 && string(recvData[:n]) == "QUIT" {
			c.Close()
			log.Println("Closed connection from server")
			return
		}
		if messageType == MESSAGE_PUT && n > 0 {
			log.Println("received message PUT")
			err := onMessage(recvData[:n])
			if err != nil {
				log.Println(err)
				c.Write([]byte("KO"))
			} else {
				c.Write([]byte("OK"))
			}
		}
		if messageType == MESSAGE_GET && n > 0 {
			log.Println("received message GET")
			message, err := readMessage(string(recvData[:n]))
			if err != nil {
				if err.Error() == "no message in given queue" {
					c.Write([]byte("EMPTY"))
				} else {
					log.Println(err)
					c.Write([]byte("KO"))
				}
			} else {
				c.Write(message)
			}
		}

		if messageType == "" && n > 0 {
			if string(recvData[:n]) == MESSAGE_GET {
				messageType = MESSAGE_GET
				c.Write([]byte("OK"))
			}
			if string(recvData[:n]) == MESSAGE_PUT {
				messageType = MESSAGE_PUT
				c.Write([]byte("OK"))
			}
			log.Println("connection configuration to " + messageType + " OK")
		}
		if messageType == "" && n > 0 {
			log.Println("Invalid state, expected connection is already configured as GET|PUT")
			c.Close()
		}

	}
}

func listenNetwork() {
	log.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("Error connecting:", err.Error())
			return
		}
		log.Println("Client connected.")
		go handleIncomingConnection(c)
	}
}
