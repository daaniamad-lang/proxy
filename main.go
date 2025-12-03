package main

import (
    "io"
    "net"
    "os"
)

const (
    V2RAY_SERVER_IP = "62.133.63.108"  // Замените на ваш IP
    TARGET_PORT     = "80"
)

func handleClient(clientConn net.Conn, targetAddr string) {
    defer clientConn.Close()

    remoteConn, err := net.Dial("tcp", targetAddr)
    if err != nil {
        return
    }
    defer remoteConn.Close()

    go func() {
        io.Copy(remoteConn, clientConn)
    }()

    io.Copy(clientConn, remoteConn)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    listenAddr := ":" + port
    targetAddr := V2RAY_SERVER_IP + ":" + TARGET_PORT
    
    listener, err := net.Listen("tcp", listenAddr)
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    for {
        clientConn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(clientConn, targetAddr)
    }
}
