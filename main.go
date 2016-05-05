package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/kobeld/goutils"
)

const (
	DEFAULT_HOST = "127.0.0.1:5555"
)

var (
	localAddr  = flag.String("local", "", "The host of Forwarder")
	remoteAddr = flag.String("remote", "", "The remote host that tcp is forwarded to")
)

func forward(srcConn net.Conn, remoteAddr string) {

	dstConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("ERROR: Dial failed: %v", err)
	}

	defer func() {
		srcConn.Close()
		dstConn.Close()

		log.Printf("INFO: Connection closed for: %s\n\n", srcConn.RemoteAddr().String())
	}()

	goutils.CoveredGo(func() {
		// Inbound forward
		io.Copy(srcConn, dstConn)
	})

	// Outbound forward
	io.Copy(dstConn, srcConn)

}

func main() {

	flag.Parse()

	if *localAddr == "" {
		log.Printf("INFO: No \"-local={host:port}\" arg provided, using the default host: %s\n", DEFAULT_HOST)
		*localAddr = DEFAULT_HOST
	}

	if *remoteAddr == "" {
		log.Fatalln("ERROR: Should provide the remote address: \"-remote={host:port}\"")
		return
	} else {
		log.Printf("INFO: The remote address is: %s\n", *remoteAddr)
	}

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatalf("ERROR: Failed to setup listener: %v", err)
	}
	log.Printf("INFO: Listening and serving incoming TCP on %s\n", *localAddr)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("INFO: Accepted connection from: %s\n", conn.RemoteAddr().String())

		goutils.CoveredGo(func() {
			forward(conn, *remoteAddr)
		})
	}
}
