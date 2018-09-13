// myserver.go
// Rich Robinson
// Sept 2018

package main
import (
    "strings"
    "os/exec"
//    "sync"
    "context"
    "time"
    "fmt"
    "log"
    "bufio"
    "net/http"
    "golang.org/x/net/trace"
    "golang.org/x/net/websocket"
)

const port = 4050

type Event struct {
    // Fields of this struct are exported for the echo module
    // to write to. Field tags specify names for JSON.
    H string `json:"h"`
    W string `json:"w"`
    X int `json: x`
    Y int `json: y`
    C string `json: c`
    A []string `json: a`
    R string `json: r`
}

func main() {
    go http.Handle("/wstime", websocket.Handler(wsTime))
    go http.Handle("/wscall", websocket.Handler(wsCall))
    go http.Handle("/", http.FileServer(http.Dir("static/html")))

    log.Printf("Server listening on port %d", port)
    go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func cpuTemp()(res string) {
    var x string
    if strings.Contains( sysType(), "ARM" ) {
        x = exeCmd( "/opt/vc/bin/vcgencmd", "measure_temp")
    } else {
        x = "cpuTemp: SysType Not ARM"
    }
    return x
}

func sysType()(res string) {
    var x string
    scanner := bufio.NewScanner(strings.NewReader(exeCmd( "cat", "/proc/cpuinfo")))
    for scanner.Scan() {
        if strings.Contains( scanner.Text(), "model name") {
            x = scanner.Text() [ 13:len( scanner.Text() ) ]
            return x
//            fmt.Println( x )
            break
        }
    }
     return "SysType Not Found"
}

func exeCmd(command string, args... string) (res string) {
// e.g.   log.Println( exeCmd("ls", "..", "-l",) )
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, command, args... )
    out, err := cmd.Output()

    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("exeCmd: Command timed out")
        return
    }
    if err != nil {
        fmt.Println("Non-zero exit code:", err)
        return "exeCmd: Error in external command"
        }
    return string(out)
}

func handleWsCall(ws *websocket.Conn, e Event) error {
    tr := trace.New("websocket.Receive", "receive")
    defer tr.Finish()
//    log.Printf("RC Event %v\n", e)
    err := websocket.JSON.Send(ws, e)
    if err != nil {
        return fmt.Errorf("Can't send: %s", err.Error())
    }
    return nil
}

func wsCall(ws *websocket.Conn) {
    log.Printf("Client %s connected", ws.RemoteAddr())
    for {
        var event Event
        err := websocket.JSON.Receive(ws, &event)
        if err != nil {
            log.Printf("Receive: %s; closing connection...", err.Error())
            if err = ws.Close(); err != nil {
                log.Println("Error closing:", err.Error())
            }
            break
        } else 
        {
// allow for non-call event (i.e. echo)
            if event.C != "" {
                event.H = "Hello"
                event.W = "World"
                switch event.C {
                    case "cpuTemp":
                        event.R = cpuTemp()
                    case "sysType":
                        event.R = sysType()
                    default:
                        event.R = none()
                }
            }
            if err := handleWsCall(ws, event); err != nil {
                log.Println(err.Error())
                break
            }
        }
    }
}

func none () (res string) {
    return ("Function Not Found")
}

func wsTime(ws *websocket.Conn) {
    for range time.Tick(1 * time.Second) {
        websocket.Message.Send(ws, time.Now().Format("Mon, 02 Jan 2006 15:04:05 PST"))
    }
}
