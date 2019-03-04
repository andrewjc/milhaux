package testing

import "net"

func GetTestConnection() (client, server net.Conn) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	var serverConn net.Conn
	go func() {
		defer ln.Close()
		server, _ = ln.Accept()
	}()

	clientConn, _ := net.Dial("tcp", ln.Addr().String())

	return clientConn, serverConn
}

