package main

import (
	"log"
	"net"

	"github.com/pion/stun"
)

func main() {
	addr := "0.0.0.0:3479"
	conn, err := net.ListenPacket("udp4", addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer conn.Close()
	log.Printf("Listening on %s", addr)

	buf := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Failed to read packet: %v", err)
			continue
		}

		message := new(stun.Message)
		message.Raw = buf[:n]
		if err := message.Decode(); err != nil {
			log.Printf("Failed to decode STUN message: %v", err)
			continue
		}

		if message.Type.Class == stun.ClassRequest {
			response := stun.New()
			response.TransactionID = message.TransactionID
			response.Type = stun.NewType(stun.MethodBinding, stun.ClassSuccessResponse)
			response.WriteHeader()

			// Add XOR-MAPPED-ADDRESS attribute
			xorAddr := &stun.XORMappedAddress{
				IP:   addr.(*net.UDPAddr).IP,
				Port: addr.(*net.UDPAddr).Port,
			}
			if err := xorAddr.AddTo(response); err != nil {
				log.Printf("Failed to add XOR-MAPPED-ADDRESS attribute: %v", err)
				continue
			}

			if _, err := conn.WriteTo(response.Raw, addr); err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		}
	}
}
