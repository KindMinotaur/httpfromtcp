package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(linesChannel)

		str := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				linesChannel <- str
				str = ""
			}

			str += string(data)
		}

		if len(str) != 0 {
			linesChannel <- str
		}
	}()

	return linesChannel
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}

	}
}
