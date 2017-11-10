package main

import (
	"io"
	"net"
	"sync"
)

type Teecp struct {
	config *Config
}

func NewTeecp(config *Config) *Teecp {
	return &Teecp{config: config}
}

func newStream(src io.Reader, dest io.Writer, wg *sync.WaitGroup) {
	io.Copy(dest, src)
	wg.Done()
}

func (t *Teecp) handleConnection(srcConn net.Conn) {
	logger.Info.Println("Connection was opened")

	defer srcConn.Close()

	// Connect to proxy server
	proxyConn, err := net.Dial("tcp", t.config.Proxy)
	if err != nil {
		logger.Fatal.Println("Error connecting to proxy", t.config.Proxy, err)
		return
	}
	logger.Debug.Println("Connected to proxy", t.config.Proxy)
	defer proxyConn.Close()

	// Connect to tee servers
	teeConns := []io.Writer{}
	for _, server := range t.config.Tees {
		teeConn, err := net.Dial("tcp", server)
		if err != nil {
			logger.Fatal.Println("Error connecting to tee", server, err)
			return
		}
		logger.Debug.Println("Connected to tee", server)
		defer teeConn.Close()
		teeConns = append(teeConns, teeConn)
	}

	// Fan-out writes from the src to all destinations.
	conns := append(teeConns, proxyConn)
	multiWriter := io.MultiWriter(conns...)

	// Build streams
	wg := sync.WaitGroup{}

	// Build stream from proxy->client
	wg.Add(1)
	go newStream(proxyConn, srcConn, &wg)

	// Build stream from tee->sink to drain the tee
	sink := new(sinkReadWriter)
	for i, _ := range teeConns {
		teeConn := teeConns[i].(net.Conn)
		wg.Add(1)
		go newStream(teeConn, sink, &wg)
	}

	// Write any incoming bytes to all destinations.
	_, err = io.Copy(multiWriter, srcConn)
	if err != nil {
		logger.Info.Println(err)
	}

	// Force all conns to close.
	logger.Debug.Println("Client side was closed")
	srcConn.Close()
	for i, _ := range conns {
		conn := conns[i].(net.Conn)
		conn.Close()
	}

	// Wait for all streams to finish.
	wg.Wait()

	logger.Info.Println("Connection was closed")
}

func (t *Teecp) ListenAndServe(config *Config) error {
	listener, err := net.Listen("tcp", config.Bind)
	if err != nil {
		return err
	}

	logger.Info.Println("Teecp listening on", config.Bind)

	acceptedConnChannel := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Warn.Println("Could not accept socket:", err)
				continue
			}

			acceptedConnChannel <- conn
		}
	}()

	for {
		go t.handleConnection(<-acceptedConnChannel)
	}

	return nil
}
