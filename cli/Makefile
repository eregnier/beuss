build:
	go build -o beuss .

sample-put:
	echo "my-message-content" | ./beuss PUT testqueue

sample-get:
	./beuss GET testqueue

sample-on:
	echo "message 1" | ./beuss PUT testqueue
	echo "message 2" | ./beuss PUT testqueue
	echo "message 3" | ./beuss PUT testqueue
	./beuss ON testqueue
