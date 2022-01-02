package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"

	env "github.com/eregnier/beuss/env"
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
	if len(message) < env.QUEUES_NAME_LENGTH {
		return errors.New("message too short : rejected")
	}
	queueName := strings.TrimSpace(string(message[:env.QUEUES_NAME_LENGTH]))
	content := message[env.QUEUES_NAME_LENGTH:]
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
		recvData := make([]byte, env.MAX_MESSAGE_SIZE)
		n, err := c.Read(recvData)
		if err != nil {
			log.Println("network error on data receive", err)
			c.Close()
			return
		}
		if n == 4 && string(recvData[:n]) == "QUIT" {
			c.Close()
			log.Println("Closed connection from server")
			return
		}
		if messageType == env.MESSAGE_PUT && n > 0 {
			log.Println("received message PUT")
			err := onMessage(recvData[:n])
			if err != nil {
				log.Println(err)
				c.Write([]byte("KO"))
			} else {
				c.Write([]byte("OK"))
			}
		}
		if messageType == env.MESSAGE_GET && n > 0 {
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
			if string(recvData[:n]) == env.MESSAGE_GET {
				messageType = env.MESSAGE_GET
				c.Write([]byte("OK"))
			}
			if string(recvData[:n]) == env.MESSAGE_PUT {
				messageType = env.MESSAGE_PUT
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
	addr := "0.0.0.0:" + env.GetPort()
	log.Println("Starting " + env.ConnType + " server on " + addr)
	l, err := net.Listen(env.ConnType, addr)
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
