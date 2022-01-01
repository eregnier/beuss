package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func newClient(connectionType string) (net.Conn, error) {
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if connectionType != MESSAGE_GET && connectionType != MESSAGE_PUT {
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

func clientPutMessage(conn net.Conn, queueName string, message []byte) error {
	l := len(queueName)
	if l > QUEUES_NAME_LENGTH {
		return errors.New("queue name too long")
	}
	queueName = queueName + strings.Repeat(" ", QUEUES_NAME_LENGTH-l)
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

func clientGetMessage(conn net.Conn, queueName string) ([]byte, error) {
	conn.Write([]byte(queueName))
	buff := make([]byte, MAX_MESSAGE_SIZE)
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

func clientClose(conn net.Conn) {
	conn.Write([]byte("QUIT"))
	time.Sleep(time.Millisecond * 300)
	conn.Close()
	log.Println("closed connection from client")
}
