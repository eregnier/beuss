package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	env "github.com/eregnier/beuss/env"
)

func NewClient(connectionType string) (net.Conn, error) {
	conn, err := net.Dial(env.ConnType, env.GetHost()+":"+env.GetPort())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if connectionType != env.MESSAGE_GET && connectionType != env.MESSAGE_PUT {
		return nil, errors.New("wrong connection type given, expected on of PUT|GET")
	}
	conn.Write([]byte(connectionType))
	buff := make([]byte, 2)
	n, err := conn.Read(buff)
	if err != nil {
		return nil, errors.New("error on server connection configuration")
	}
	if string(buff[:n]) == "OK" {
		log.Println("connection to server OK")
		return conn, nil
	} else {
		conn.Close()
		return nil, fmt.Errorf("could not configure connection type")
	}
}

func ClientPutMessage(conn net.Conn, queueName string, message []byte) error {
	l := len(queueName)
	if l > env.QUEUES_NAME_LENGTH {
		return errors.New("queue name too long")
	}
	queueName = queueName + strings.Repeat(" ", env.QUEUES_NAME_LENGTH-l)
	message = append([]byte(queueName), message...)
	conn.Write(message)
	buff := make([]byte, 2)
	n, _ := conn.Read(buff)
	if string(buff[:n]) == "OK" {
		return nil
	} else {
		return fmt.Errorf("error sending message : received | %s", buff)
	}
}

func ClientGetMessage(conn net.Conn, queueName string) ([]byte, error) {
	conn.Write([]byte(queueName))
	buff := make([]byte, env.MAX_MESSAGE_SIZE)
	n, err := conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("error with read message connexion")
	}
	if string(buff[:5]) == "EMPTY" {
		return nil, fmt.Errorf("empty queue")
	}
	if string(buff[:2]) == "KO" {
		return nil, fmt.Errorf("error with read message response")
	}
	return buff[:n], nil
}

func ClientOnMessage(conn net.Conn, queueName string, callback func(message []byte)) {
	for {
		message, err := ClientGetMessage(conn, queueName)
		if err != nil {
			if err.Error() == "empty queue" {
				time.Sleep(time.Second * env.CONSUME_DELAY)
			} else {
				log.Println("server error |", err)
				return
			}
		} else {
			callback(message)
		}
	}
}

func ClientClose(conn net.Conn) {
	conn.Write([]byte("QUIT"))
	time.Sleep(time.Millisecond * 300)
	conn.Close()
	log.Println("closed connection from client")
}
