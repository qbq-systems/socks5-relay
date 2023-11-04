package main

import (
    "encoding/json"
    "fmt"
    "github.com/savsgio/atreugo/v11"
    "github.com/valyala/fasthttp"
    "log"
    "os"
    "strconv"
    "time"
)

type RestRequest struct {
    Host     string `json:"host"`
    Port     string `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
}

type RestResponse struct {
    Port    int    `json:"port"`
    Message string `json:"message"`
}

var relays = make(map[string]*Socks5Relay)

func main() {

    portEnv := os.Getenv("LISTEN_PORT")
    if portEnv == "" {
        portEnv = "18000"
    }
    port, err := strconv.Atoi(portEnv)
    if err != nil {
        log.Fatal(err)
    }

    server := atreugo.New(atreugo.Config{
        Addr: ":" + strconv.Itoa(port),
    })
    server.POST("/", func(ctx *atreugo.RequestCtx) error {
        return handler("POST", ctx)
    })
    server.DELETE("/", func(ctx *atreugo.RequestCtx) error {
        return handler("DELETE", ctx)
    })

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}

func handler(method string, ctx *atreugo.RequestCtx) error {

    sendResponse := func(ctx *atreugo.RequestCtx, port int, message string, err bool) error {
        var status int
        if err {
            status = fasthttp.StatusBadRequest
        } else {
            status = fasthttp.StatusOK
        }

        response := RestResponse{
            Port:    port,
            Message: message,
        }
        responseJSON, _ := json.Marshal(response)

        ctx.SetContentType("application/json")
        ctx.SetStatusCode(status)
        ctx.SetBody(responseJSON)

        return nil
    }

    body := ctx.PostBody()
    var data RestRequest
    if err := json.Unmarshal(body, &data); err != nil {
        return sendResponse(ctx, 0, fmt.Sprintf("Error unmarshaling JSON: %v", err), true)
    }

    remote := data.Host + ":" + data.Port
    val, exists := relays[remote]

    var port = 0
    var message = ""
    var responseError = false

    if method == "POST" {
        if exists {
            message = "Exists"
            port = val.LocalPort
        } else {
            result := make(chan bool)
            go func() {
                step := 0
                for {
                    if step == 20 {
                        message = "Exhausted connection attempts"
                        responseError = true
                    }
                    relay := NewSocks5Relay(data.Host, data.Port, data.User, data.Password)
                    err := relay.Connect()
                    if err == nil {
                        relays[remote] = relay
                        port = relays[remote].LocalPort
                        message = "Added"
                        result <- true
                        break
                    }
                    step += 1
                    time.Sleep(1 * time.Second)
                }
            }()
            _ = <-result
        }
    } else {
        if exists {
            relays[remote].Stop()
            delete(relays, remote)
            message = "Deleted"
        } else {
            message = "Not exists"
        }
    }

    return sendResponse(ctx, port, message, responseError)
}
