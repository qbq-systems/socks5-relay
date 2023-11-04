package main

import (
    "context"
    "github.com/things-go/go-socks5"
    "golang.org/x/net/proxy"
    "log"
    "net"
)

type Socks5Relay struct {
    RemoteHost     string
    RemotePort     string
    RemoteUser     string
    RemotePassword string
    LocalPort      int
    listener       net.Listener
    remote         string
    service        *socks5.Server
}

func NewSocks5Relay(remoteHost, remotePort, remoteUser, remotePassword string) *Socks5Relay {
    return &Socks5Relay{
        RemoteHost:     remoteHost,
        RemotePort:     remotePort,
        RemoteUser:     remoteUser,
        RemotePassword: remotePassword,
        LocalPort:      0,
    }
}

func (s *Socks5Relay) Connect() error {

    s.remote = s.RemoteHost + ":" + s.RemotePort
    log.Println(s.remote, "creating relay ...")

    var relay proxy.Dialer
    var err error
    relay, err = proxy.SOCKS5("tcp", s.remote, &proxy.Auth{
        User:     s.RemoteUser,
        Password: s.RemotePassword,
    }, proxy.Direct)
    if err != nil {
        log.Println(s.remote, "SOCKS5 error:", err)
        return err
    }

    worker := func(ctx context.Context, network, addr string) (net.Conn, error) {
        return relay.Dial(network, addr)
    }
    s.service = socks5.NewServer(socks5.WithDial(worker))

    result := make(chan interface{})

    go func() {
        var listener net.Listener
        for {
            listener, err = net.Listen("tcp", ":0")
            if err == nil {
                break
            }
        }
        s.listener = listener
        port := s.listener.Addr().(*net.TCPAddr).Port
        result <- port
        //
        // todo Perhaps a health check is needed every n seconds
        //
        s.service.Serve(s.listener)
    }()

    received := <-result
    data, isInt := received.(int)
    if isInt {
        log.Println(s.remote, "local port is", data)
        s.LocalPort = data
    } else {
        log.Println(s.remote, "listen error:", data)
        return err
    }

    return nil
}

func (s *Socks5Relay) Stop() {
    log.Println(s.remote, "stopping process with local port", s.LocalPort, "...")
    s.listener.Close()
}
