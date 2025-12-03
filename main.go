package main

import (
    "io"
    "log"
    "net"
    "os"
)

const (
    V2RAY_SERVER_IP = "62.133.63.108"
    TARGET_PORT     = "80"
)

func handleClient(clientConn net.Conn, targetAddr string) {
    defer clientConn.Close()

    remoteConn, err := net.Dial("tcp", targetAddr)
    if err != nil {
        log.Printf("Failed to connect to %s: %v", targetAddr, err)
        return
    }
    defer remoteConn.Close()

    // Client -> Remote (в фоне)
    go func() {
        _, err := io.Copy(remoteConn, clientConn)
        if err != nil {
            log.Printf("Client->Remote copy error: %v", err)
        }
        remoteConn.Close()
    }()

    // Remote -> Client
    _, err = io.Copy(clientConn, remoteConn)
    if err != nil {
        log.Printf("Remote->Client copy error: %v", err)
    }
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    listenAddr := ":" + port
    targetAddr := V2RAY_SERVER_IP + ":" + TARGET_PORT
    
    log.Printf("Starting proxy server: %s -> %s:%s", listenAddr, V2RAY_SERVER_IP, TARGET_PORT)
    
    listener, err := net.Listen("tcp", listenAddr)
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }
    defer listener.Close()

    log.Println("Proxy server started successfully")
    
    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Printf("Accept error: %v", err)
            continue
        }
        log.Printf("New connection from %s", clientConn.RemoteAddr())
        go handleClient(clientConn, targetAddr)
    }
}
