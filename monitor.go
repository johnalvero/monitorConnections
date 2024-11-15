package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Connection struct {
	SrcIP    string
	SrcPort  uint16
	DstIP    string
	DstPort  uint16
	Protocol string
}

// State table to track TCP connections and avoid duplicate outputs for UDP
var connectionStateTable = make(map[string]*Connection)
var excludePublic bool
var excludeUDP bool
var generateCSV bool
var device string

func main() {
	// Parse command-line flags
	flag.BoolVar(&excludePublic, "exclude-public", false, "Exclude connections to/from public IP addresses")
	flag.BoolVar(&excludeUDP, "exclude-udp", false, "Exclude UDP connections")
	flag.BoolVar(&generateCSV,"generate-csv", false, "Generate CSV file for plotting")
	flag.StringVar(&device, "device", "eth0", "Network device to monitor (e.g., eth0)")
	flag.Parse()

	// Open the device for capturing
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening device %s: %v", device, err)
	}
	defer handle.Close()

	// Set a BPF filter based on the excludeUDP flag
	bpfFilter := "tcp"
	if !excludeUDP {
		bpfFilter += " or udp"
	}
	err = handle.SetBPFFilter(bpfFilter)
	if err != nil {
		log.Fatal(err)
	}

	// Use a packet source to process each packet
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		processPacket(packet)
	}
}

func processPacket(packet gopacket.Packet) {
	// Extract the IP layer
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}

	ip, _ := ipLayer.(*layers.IPv4)

	// Check for private IP addresses based on the flag
	if excludePublic && (!isPrivateIP(ip.SrcIP) || !isPrivateIP(ip.DstIP)) {
		return
	}

	// Check for TCP and UDP layers
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		processTCPConnection(ip, tcp)
	} else if !excludeUDP {
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			processUDPConnection(ip, udp)
		}
	}
}

func processTCPConnection(ip *layers.IPv4, tcp *layers.TCP) {
	// Build a unique connection key
	connKey := fmt.Sprintf("%s:%d -> %s:%d", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)

	// Check if this is a SYN packet (TCP connection initiation)
	if tcp.SYN && !tcp.ACK {
		// Track connection initiation in the state table
		connectionStateTable[connKey] = &Connection{
			SrcIP:    ip.SrcIP.String(),
			SrcPort:  uint16(tcp.SrcPort),
			DstIP:    ip.DstIP.String(),
			DstPort:  uint16(tcp.DstPort),
			Protocol: "TCP",
		}
	} else if tcp.ACK && !tcp.SYN && !tcp.FIN && !tcp.RST {
		// This is the final ACK in a three-way handshake
		if conn, exists := connectionStateTable[connKey]; exists {
			if generateCSV {
				printCSVConnection(conn)
			} else {
				printConnection(conn)
			}
			delete(connectionStateTable, connKey) // Remove connection after printing
		}
	}
}

func processUDPConnection(ip *layers.IPv4, udp *layers.UDP) {
	// Build a unique connection key
	connKey := fmt.Sprintf("%s:%d -> %s:%d", ip.SrcIP, udp.SrcPort, ip.DstIP, udp.DstPort)

	// Check if the connection already exists in the state table
	if _, exists := connectionStateTable[connKey]; !exists {
		// Create a new connection entry for UDP and print it immediately
		conn := &Connection{
			SrcIP:    ip.SrcIP.String(),
			SrcPort:  uint16(udp.SrcPort),
			DstIP:    ip.DstIP.String(),
			DstPort:  uint16(udp.DstPort),
			Protocol: "UDP",
		}
		connectionStateTable[connKey] = conn
		if generateCSV {
			printCSVConnection(conn)
		} else {
			printConnection(conn)
		}
	}
}

func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func printConnection(conn *Connection) {
	fmt.Printf("Initiated Connection - Src: %s:%d -> Dst: %s:%d, Protocol: %s\n",
		conn.SrcIP, conn.SrcPort, conn.DstIP, conn.DstPort, conn.Protocol)
	time.Sleep(time.Millisecond * 100) // To ensure output synchronization in console
}

func printCSVConnection(conn *Connection) {
	//fmt.Printf("(\"%s\",\"%s\",\"%s:%d\"),\n",conn.SrcIP, conn.DstIP, conn.Protocol, conn.DstPort)
	fmt.Printf("%s,%s,%s:%d\n",conn.SrcIP, conn.DstIP, conn.Protocol, conn.DstPort)
	time.Sleep(time.Millisecond * 100) // To ensure output synchronization in console
}
