package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if os.Getenv("DEBUG") != "1" {
		log.SetOutput(ioutil.Discard)
	}
	queues = make([]Queue, 0)
	go listenNetwork()
	time.Sleep(time.Millisecond * 100)
}

func TestPutAndGetMessage(t *testing.T) {
	connPUT, err := newClient(MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := newClient(MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer clientClose(connGET)
	defer clientClose(connPUT)

	err = clientPutMessage(connPUT, "soupai", []byte("{\"hello\": \"world\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	err = clientPutMessage(connPUT, "soupai", []byte("{\"g\": \"g\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	var message []byte
	message, err = clientGetMessage(connGET, "soupai")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"hello\": \"world\"}" {
		t.Errorf("got wrong message")
	}
	message, err = clientGetMessage(connGET, "soupai")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"g\": \"g\"}" {
		t.Errorf("got wrong message")
	}
	_, err = clientGetMessage(connGET, "soupai")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}

}

func TestPutAndGetMessageTwoQueues(t *testing.T) {
	connPUT, err := newClient(MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := newClient(MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer clientClose(connGET)
	defer clientClose(connPUT)

	err = clientPutMessage(connPUT, "q1", []byte("{\"hello\": \"world\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	err = clientPutMessage(connPUT, "q2", []byte("{\"g\": \"g\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	var message []byte
	_, err = clientGetMessage(connGET, "q3")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error %s", err)
	}
	message, err = clientGetMessage(connGET, "q2")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"g\": \"g\"}" {
		t.Errorf("got wrong message")
	}
	message, err = clientGetMessage(connGET, "q1")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"hello\": \"world\"}" {
		t.Errorf("got wrong message")
	}
	_, err = clientGetMessage(connGET, "q2")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	_, err = clientGetMessage(connGET, "q1")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
}

func TestBenchSeq(t *testing.T) {
	connPUT, err := newClient(MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := newClient(MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer clientClose(connGET)
	defer clientClose(connPUT)
	n := time.Now().UnixMilli()
	for i := 0; i < 10000; i++ {
		err = clientPutMessage(connPUT, "bench", []byte("{\"hello\": \"world\"}"))
		if err != nil {
			log.Println("Error sending message :/", err)
		}
	}
	for i := 0; i < 10000; i++ {
		message, err := clientGetMessage(connGET, "bench")
		if err != nil {
			t.Errorf("fail to get message from client")
		} else if string(message) != "{\"hello\": \"world\"}" {
			t.Errorf("got wrong message")
		}
	}
	_, err = clientGetMessage(connGET, "bench")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	if time.Now().UnixMilli()-n > 3000 {
		t.Errorf("processing time too long")
	}
}

func TestBenchParallel(t *testing.T) {
	connPUT, err := newClient(MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := newClient(MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer clientClose(connGET)
	defer clientClose(connPUT)
	n := time.Now().UnixMilli()
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		func() {

			err = clientPutMessage(connPUT, "bench", []byte("{\"hello\": \"world\"}"))
			if err != nil {
				log.Println("Error sending message :/", err)
			}
			message, err := clientGetMessage(connGET, "bench")
			if err != nil {
				t.Errorf("fail to get message from client")
			} else if string(message) != "{\"hello\": \"world\"}" {
				t.Errorf("got wrong message")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	_, err = clientGetMessage(connGET, "bench")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	if time.Now().UnixMilli()-n > 3000 {
		t.Errorf("processing time too long")
	}
}
