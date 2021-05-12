package main

import(
  "bufio"
  "flag"
  "fmt"
  "net"
  "os"
  "regexp"
  "strconv"
  "time"
)


var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")

func main() {

  //getting ip and port from command line
  flag.Parse()
  dest := *host + ":" + strconv.Itoa(*port)


  // socket creation and connection to server
  fmt.Printf("Connecting to %s...\n",dest)
  conn, err := net.Dial("tcp", dest)

  if err != nil {
    if _,t:= err.(*net.OpError);t{
      fmt.Println("Some problem connecting.")
    } else{
      fmt.Println("Unknown error: " + err.Error())
    }
    os.Exit(1)
  }

  
  
  // call function to handle connection asynchronously

  go handleConnection(conn)

  welcome,_ := bufio.NewReader(conn).ReadString('=')
  fmt.Print(string(welcome))

  //getting input from the client and writes it into the socket
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    text, _ := reader.ReadString('\n')

    conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
    _, err := conn.Write([]byte(text))
    if err != nil {
      fmt.Println("Error writing to stream.")
      break
    }
  }

  /*
  for { 
    // what to send?
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    // send to server
    fmt.Fprintf(conn, text + "\n")
    // wait for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message from server: "+message)
  }
  */
}

func handleConnection(conn net.Conn) {
  //socket reader
  scanner := bufio.NewScanner(conn)


  //print welcome message
  for scanner.Scan(){
    line := scanner.Text()
    if line == "=" {
      break
    }
    fmt.Println(line)

  }
  
  for {
      

    ok := scanner.Scan()
    text := scanner.Text()
    fmt.Println(scanner.Text())

    command := handleCommands(text,conn)
    if !command{
      fmt.Printf("\b\b** %s\n> ", text)
    }

    if !ok{
      fmt.Println("Reached EOF on server connection.")
      break

    }
  }
  //read welcome message from the socket
  /*
  message,_ := sock_reader.ReadString('=')
  fmt.Print(string(welcome))
  */

}


func handleCommands(text string,conn net.Conn) bool {
  r, err := regexp.Compile("^%.*%$")

  if err == nil {
    if r.MatchString(text) {

      switch {
      case text == "%quit%":
        fmt.Println("\b\bServer is leaving. Hanging up.")
        os.Exit(0)
      case text == "%getPass%" :
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter the username\t:")
        username, _ := reader.ReadString('\n')
        fmt.Print("Enter the website(starting with www):")
        website, _ := reader.ReadString('\n')
        fmt.Println("Enter the master password\t:")
        secret, _ := reader.ReadString('\n')
        _,err := conn.Write([]byte(username+"\n"+website+"\n"+ secret + "\n"))
        if err != nil {
          fmt.Println("Error writing to stream.")
          break
          }
        message,_:= bufio.NewReader(conn).ReadString('\n')
        fmt.Print(string(message))

      }

     return true 
    }
  }
  return true

}
