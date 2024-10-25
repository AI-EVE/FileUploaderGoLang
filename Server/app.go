package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listner ,err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("fialed to create a listner", err)
		return
	}
	defer listner.Close()

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)
	go func ()  {
		<- closeChan

		listner.Close()
		fmt.Println("Shutting down the server...")
		os.Exit(0)
	}()


	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("failed to accept a connection", err)
			continue
		}



		go func ()  {
			defer conn.Close()

			fileName := make([]byte, 50)
			n, _ := conn.Read(fileName)
	
			file, err := os.OpenFile(string(fileName[:n]), os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				fmt.Println("Failed to open or create the file:", err)
				return
			}
			
			defer file.Close()
		
			buffer := make([]byte, 16000)
			writer := bufio.NewWriterSize(file, 16000)
			
			for {
				n, err := conn.Read(buffer)
				if err != nil &&  err != io.EOF{
					fmt.Println("Failed to read from connection:", err)
					return
				}

				if n == 0 {
					fmt.Println("file was read completly")
					writer.Flush()
					file.Close()
					return
				}

				_, err = writer.Write(buffer[:n])
				if err != nil {
					fmt.Println("Failed to write to file:", err)
					return
				}
			}
		}()
	}
}