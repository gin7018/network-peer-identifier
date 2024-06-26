package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
)

func make_packet() []byte {
	message := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getegid() & 0xffff,
			Seq:  0,
			Data: []byte("hi from flip-phone!"),
		},
	}

	encoded_msg, err := message.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	return encoded_msg
}

func receive_messages(conn *icmp.PacketConn) {

	buffer := make([]byte, 1024)
	bSize, peer_ip_addr, err := conn.ReadFrom(buffer)
	fmt.Printf("buffer size %v\n", bSize)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response from ip: " + peer_ip_addr.String())

	peer_message, err := icmp.ParseMessage(1, buffer[:bSize])
	if err != nil {
		log.Fatal(err)
	}

	switch peer_message.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("received message from %v\n", peer_ip_addr)
	default:
		fmt.Printf("Failed: %+v\n", peer_message)
	}

}

func main() {
	fmt.Println("begin network discovery...")

	connection, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		log.Fatal("error setting up connection: " + err.Error())
	}

	defer connection.Close()

	packet := make_packet()
	// send the broadcast packet to 255.255.255.255
	broadcast_addr, err := net.ResolveIPAddr("ip4", "255.255.255.255")
	if err != nil {
		log.Fatal("error resolving ip addr: " + err.Error())
	}
	fmt.Printf("ip addr resolved. ")

	if _, err := connection.WriteTo(packet, broadcast_addr); err != nil {
		log.Fatal("error writing to addr: " + err.Error())
	}
	fmt.Println("broadcast message sent. waiting for response")

	receive_messages(connection)

}
