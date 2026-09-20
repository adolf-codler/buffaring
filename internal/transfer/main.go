package transfer

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"adolf-codler/buffaring/internal/discovery"
)

const TCPPort = "8902"

// status prints a \r-overwriting status line to stderr so stdout stays clean for piping.
func status(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "\r\033[K"+format, a...)
}

// done moves to a new line on stderr after the last status message.
func done() {
	fmt.Fprintln(os.Stderr)
}

// SendBuff discovers a receiver via UDP, waits for it to be ready, then sends the buffer over TCP.
//
// Protocol:
//  1. Sender broadcasts discovery.Phrase on SenderPort (so the receiver can find us).
//  2. Sender listens on ReceiverPort for the receiver's ACK broadcast.
//  3. Once ACK is received, the receiver's TCP listener is up — dial and send.
func SendBuff(buff string) error {
	// Step 1: broadcast so the receiver knows we exist
	status("Broadcasting discovery ...")
	go func() {
		_ = discovery.Broadcast(discovery.SenderPort, discovery.Phrase)
	}()

	// Step 2: wait for the receiver to ACK (meaning its TCP listener is ready)
	status("Waiting for receiver ...")
	receiverAddr, err := discovery.Listen(discovery.ReceiverPort, discovery.AckPhrase)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	// Step 3: small grace period then dial TCP
	status("Connecting to %s ...", receiverAddr.IP)
	time.Sleep(200 * time.Millisecond)

	conn, err := net.Dial("tcp", net.JoinHostPort(receiverAddr.IP.String(), TCPPort))
	if err != nil {
		return fmt.Errorf("TCP dial error: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(buff)); err != nil {
		return fmt.Errorf("TCP write error: %w", err)
	}

	status("Buffer sent ✓")
	done()
	return nil
}

// RecvBuff listens for a sender's UDP broadcast, opens a TCP listener,
// ACKs back, and reads the incoming buffer.
//
// Protocol:
//  1. Receiver listens on SenderPort for the sender's broadcast.
//  2. Receiver opens TCP listener.
//  3. Receiver broadcasts ACK on ReceiverPort so the sender knows TCP is ready.
//  4. Receiver accepts TCP connection and reads the buffer.
func RecvBuff() (string, error) {
	// Step 1: discover the sender
	status("Listening for sender ...")
	_, err := discovery.Listen(discovery.SenderPort, discovery.Phrase)
	if err != nil {
		return "", fmt.Errorf("discovery failed: %w", err)
	}

	// Step 2: open TCP listener first, before we ACK
	status("Starting TCP listener ...")
	listener, err := net.Listen("tcp", ":"+TCPPort)
	if err != nil {
		return "", fmt.Errorf("TCP listen error: %w", err)
	}
	defer listener.Close()

	// Step 3: broadcast ACK so the sender knows we're ready
	status("Waiting for connection ...")
	go func() {
		_ = discovery.Broadcast(discovery.ReceiverPort, discovery.AckPhrase)
	}()

	// Step 4: accept TCP connection and read
	conn, err := listener.Accept()
	if err != nil {
		return "", fmt.Errorf("TCP accept error: %w", err)
	}
	defer conn.Close()

	status("Receiving buffer ...")
	data, err := io.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("TCP read error: %w", err)
	}

	status("Buffer received ✓")
	done()
	return string(data), nil
}
