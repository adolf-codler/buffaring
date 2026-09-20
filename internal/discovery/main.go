package discovery

import (
	"fmt"
	"net"
	"os"
	"time"
)

const Phrase = "BUFFINGA"
const AckPhrase = "BUFFINGA-ACK"

// SenderPort is the port the sender broadcasts on (receiver listens here).
const SenderPort = "8901"

// ReceiverPort is the port the receiver broadcasts ACK on (sender listens here).
const ReceiverPort = "8903"

// status prints a \r-overwriting status line to stderr so stdout stays clean for piping.
func status(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "\r\033[K"+format, a...)
}

// Broadcast continuously sends a message on the given broadcast port.
func Broadcast(broadPort string, message string) error {
	broadIP, err := getSubnet()
	if err != nil {
		return fmt.Errorf("subnet error: %w", err)
	}

	broadAddr := net.JoinHostPort(broadIP, broadPort)
	addr, err := net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil {
		return fmt.Errorf("resolve error: %w", err)
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	msg := []byte(message)
	for {
		conn.Write(msg)
		time.Sleep(time.Second * 1)
	}
}

// Listen waits for a UDP message matching expectedPhrase on the given port
// and returns the remote address.
func Listen(port string, expectedPhrase string) (net.UDPAddr, error) {
	broadAddr := net.JoinHostPort("", port)
	addr, err := net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil {
		return net.UDPAddr{}, fmt.Errorf("resolve error: %w", err)
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return net.UDPAddr{}, fmt.Errorf("listen error: %w", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return net.UDPAddr{}, fmt.Errorf("reading from UDP: %w", err)
		}
		msg := string(buf[:n])
		if msg == expectedPhrase {
			status("Discovered peer at %s", remoteAddr.IP)
			return *remoteAddr, nil
		}
	}
}

func getSubnet() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4()
				mask := ipnet.Mask
				broadcast := net.IP(make([]byte, 4))
				for i := range ip {
					broadcast[i] = ip[i] | ^mask[i]
				}
				return broadcast.String(), nil
			}
		}
	}
	return "255.255.255.255", nil
}
