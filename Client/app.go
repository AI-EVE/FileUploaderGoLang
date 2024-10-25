package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)
var fileName string = "upload-me.txt"
func main() {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("fialed to open a file", err)
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	fileSize := stat.Size()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("failed to connect", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 16000)
	reader := bufio.NewReaderSize(file, 16000)
	writer := bufio.NewWriterSize(conn, 16000)
	totalRead := 0
	totalWritten := 0
	conn.Write([]byte(fileName))
	for {
		n, err := reader.Read(buffer)
		if n == 0 {
			fmt.Println("File has been read completely")
			writer.Flush()
			fmt.Printf("\rwriting progress 100 / %.4f\n", float32(totalRead) / float32(fileSize) * 100.0)
			return
		}
		totalRead += n
		n, err = writer.Write(buffer[:n])
		if err != nil {
			fmt.Println("failed to write to the server", err)
			return
		}
		totalWritten += n

		fmt.Printf("\rwriting progress 100 / %.4f\n", float32(totalWritten) / float32(fileSize) * 100.0)

		time.Sleep(1 * time.Second)
	}

	// for i := 0; i < 1000000; i++ {
	// 	file.WriteString(fmt.Sprintf("%d ", i))
	// }


}
