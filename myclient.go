// myclient.go
// Rich Robinson
// Sept 2018

package main

import (
    "fmt"
    "log"
    "time"
    "os"
    "os/signal"
    "syscall"
    "golang.org/x/net/websocket"
)

type Event struct {
    H string `json:"h"`
    W string `json:"w"`
    X int `json: x`
    Y int `json: y`
    C string `json: c`
    A []string `json: a`
    R string `json: r`
}

// server
var c = "ws://c.local"
var rock = "ws://rock.local"
// const serverport = 4050

func main() {
// initialise getout
    signalChannel := make(chan os.Signal, 2)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
    go func() {
        sig := <-signalChannel
        switch sig {
        case os.Interrupt:
            log.Println("Stopping on Interrupt")
            time.Sleep(2 * time.Second)
            os.Exit(1)
        case syscall.SIGTERM:
            log.Println("Stopping on Terminate")
            time.Sleep(2 * time.Second)
            os.Exit(0)
        }
    }()
// main code follows
    for{
// first one no longer allowed
//        log.Println(WsCall("ls", "..", "-la"))
// do you need a different port per server? seems ok
        log.Println(WsCall(rock, 4050, "sysType", ""))
        log.Println(WsCall(c, 4050, "cpuTemp", ""))

        time.Sleep(2 * time.Second)
    }
}

func WsCall(server string, port int, command string, args... string) (res string) {
    addr := fmt.Sprintf("%s:%d/wscall", server, port)
    
    origin, err := os.Hostname()
    if err != nil {
        log.Fatal("Hostname not found", err)
    }
//  origin is own host address (note the .local)
    conn, err := websocket.Dial(addr, "", "ws://" + origin + ".local")
    if err != nil {
        log.Fatal("websocket.Dial error", err)
    }

    e := Event{} ; e.C = command ; e.A = args

    err = websocket.JSON.Send(conn, e)
    if err != nil {
        log.Fatal("websocket.JSON.Send error", err)
    }

    var reply Event
    err = websocket.JSON.Receive(conn, &reply)
    if err != nil {
        log.Fatal("websocket.JSON.Receive error", err)
    }
    log.Printf("reply: %s", reply.R )

    if err = conn.Close(); err != nil {
        log.Fatal("conn.Close error", err)
    }
    return e.W
}
