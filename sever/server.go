package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"runtime"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func main(){
	
	flag.Parse()
	
	fmt.Println("Start server...")

	//listen on port 8000

	src := *addr + ":" + strconv.Itoa(*port)

	// does bind() and listen()
	listener, _ := net.Listen("tcp",src)
	fmt.Printf("Listening on %s.\n",src)

	//add socketfile.close() to the bottom of 
	//the stack for resource clean-up after everything is over

	defer listener.Close()

	//listen forever for a connection
	for{
		conn, err := listener.Accept()
		if err != nil{
			fmt.Printf("Some connection error :%s\n",err)
		}

		//call the socket communication function asynchronously by running it as a goroutine
		go handleConnection(conn)

	}

	/*
	for{
		message,_:= bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Recvd : ",string(message))
	}
	*/
}


func handleConnection(conn net.Conn){
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client Connected from " + remoteAddr)

	// write first message into client socket

	fmt.Fprintf(conn,"\nWelcome to our simple stateless password manager\n")
	fmt.Fprintf(conn,"Here are the commands that you can use:-\n")
	fmt.Fprintf(conn,"/getPass\t Enter username , website(starting with www) and master secret to get the password\n")
	fmt.Fprintf(conn,"/quit\t\t Quit server connection\n")
	fmt.Fprintf(conn,"/help\t\t display help options\n")
	fmt.Fprintf(conn,"/getSpec\t Get program specifications\n=")

	scanner := bufio.NewScanner(conn)
	for{
		ok := scanner.Scan()
		if !ok{
			break
		}
		handleMessage(scanner.Text(),conn)

	}
	fmt.Println("Client at " + remoteAddr + "disconnected.")
}


func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)
	
	if len(message) >0 && message[0] == '/' {
		switch{
		case message == "/help":
			resp := "Time : " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			fmt.Fprintf(conn,"Here are the commands that you can use:-\n")
			fmt.Fprintf(conn,"/getPass\t Enter username , website(starting with www) and master secret to get the password\n")
			fmt.Fprintf(conn,"/quit\t\t Quit server connection\n")
			fmt.Fprintf(conn,"/help\t\t display help options\n")
			fmt.Fprintf(conn,"/getSpec\t Get program specifications\n=")
			conn.Write([]byte(resp))

		case message == "/getPass" :
			
			reader := bufio.NewReader(conn)
			username,_ := reader.ReadString('\n')
			website,_ := reader.ReadString('\n')
			secret,_ := reader.ReadString('\n')

			conn.Write([]byte("Generated password\t:"+getPass(username,website,secret) + "\n"))

		case message == "/quit":
			fmt.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			fmt.Println("< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)
		case message == "/getSpec" :
			fmt.Fprintf(conn,"Go version:\t %s\n", runtime.Version())
			fmt.Fprintf(conn,"Connection:\t TCP socket\n")
			fmt.Fprintf(conn,"Platform: \t %s\n",runtime.GOOS)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		
		}
	}
}

func getPass(username string, website string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(username+website))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}