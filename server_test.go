package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	b "github.com/eregnier/beuss/client"
	env "github.com/eregnier/beuss/env"
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
	connPUT, err := b.NewClient(env.MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := b.NewClient(env.MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer b.ClientClose(connGET)
	defer b.ClientClose(connPUT)

	err = b.ClientPutMessage(connPUT, "soupai", []byte("{\"hello\": \"world\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	err = b.ClientPutMessage(connPUT, "soupai", []byte("{\"g\": \"g\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	var message []byte
	message, err = b.ClientGetMessage(connGET, "soupai")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"hello\": \"world\"}" {
		t.Errorf("got wrong message")
	}
	message, err = b.ClientGetMessage(connGET, "soupai")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"g\": \"g\"}" {
		t.Errorf("got wrong message")
	}
	_, err = b.ClientGetMessage(connGET, "soupai")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}

}

func TestPutAndGetMessageTwoQueues(t *testing.T) {
	connPUT, err := b.NewClient(env.MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := b.NewClient(env.MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer b.ClientClose(connGET)
	defer b.ClientClose(connPUT)

	err = b.ClientPutMessage(connPUT, "q1", []byte("{\"hello\": \"world\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	err = b.ClientPutMessage(connPUT, "q2", []byte("{\"g\": \"g\"}"))
	if err != nil {
		log.Println("Error sending message :/", err)
	}

	var message []byte
	_, err = b.ClientGetMessage(connGET, "q3")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error %s", err)
	}
	message, err = b.ClientGetMessage(connGET, "q2")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"g\": \"g\"}" {
		t.Errorf("got wrong message")
	}
	message, err = b.ClientGetMessage(connGET, "q1")
	if err != nil {
		t.Errorf("fail to get message from client")
	} else if string(message) != "{\"hello\": \"world\"}" {
		t.Errorf("got wrong message")
	}
	_, err = b.ClientGetMessage(connGET, "q2")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	_, err = b.ClientGetMessage(connGET, "q1")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
}

func TestBenchSeq(t *testing.T) {
	connPUT, err := b.NewClient(env.MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := b.NewClient(env.MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer b.ClientClose(connGET)
	defer b.ClientClose(connPUT)
	n := time.Now().UnixMilli()
	for i := 0; i < 10000; i++ {
		err = b.ClientPutMessage(connPUT, "bench", []byte("{\"hello\": \"world\"}"))
		if err != nil {
			log.Println("Error sending message :/", err)
		}
	}
	for i := 0; i < 10000; i++ {
		message, err := b.ClientGetMessage(connGET, "bench")
		if err != nil {
			t.Errorf("fail to get message from client")
		} else if string(message) != "{\"hello\": \"world\"}" {
			t.Errorf("got wrong message")
		}
	}
	_, err = b.ClientGetMessage(connGET, "bench")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	if time.Now().UnixMilli()-n > 3000 {
		t.Errorf("processing time too long")
	}
}

func TestBenchParallel(t *testing.T) {
	connPUT, err := b.NewClient(env.MESSAGE_PUT)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	connGET, err := b.NewClient(env.MESSAGE_GET)
	if err != nil {
		log.Println("error on client GET creation", err)
		os.Exit(1)
	}
	defer b.ClientClose(connGET)
	defer b.ClientClose(connPUT)
	n := time.Now().UnixMilli()
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		func() {

			err = b.ClientPutMessage(connPUT, "bench", []byte("{\"hello\": \"world\"}"))
			if err != nil {
				log.Println("Error sending message :/", err)
			}
			message, err := b.ClientGetMessage(connGET, "bench")
			if err != nil {
				t.Errorf("fail to get message from client")
			} else if string(message) != "{\"hello\": \"world\"}" {
				t.Errorf("got wrong message")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	_, err = b.ClientGetMessage(connGET, "bench")
	if err.Error() != "empty queue" {
		t.Errorf("expected empty queue error")
	}
	if time.Now().UnixMilli()-n > 3000 {
		t.Errorf("processing time too long")
	}
}
