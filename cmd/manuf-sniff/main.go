package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/picatz/manuf/pkg/index"
)

func handleError(msg string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s: %v\n", msg, err))
		os.Exit(1)
	}
}

func findMatch(records index.Records, macAddr string) (*index.Record, bool) {
	macAddr = strings.ToUpper(strings.ReplaceAll(macAddr, ":", ""))

	for _, record := range records {
		if strings.HasPrefix(macAddr, string(record.Assignment)) {
			return record, true
		}
	}

	return nil, false
}

func main() {
	dir, err := os.UserCacheDir()
	handleError("failed to get current user's cache directory: %w", err)
	path := filepath.Join(dir, "manuf.csv")

	records, err := index.RecordsFromFile(path)
	handleError("failed to read records from cache", err)

	// Use the first, non-loopback interface if no interface is given.
	// https://github.com/picatz/iface/blob/master/iface.go#L90
	var iface string
	ifaces, err := net.Interfaces()
	handleError("failed list interfaces", err)
	for _, ifc := range ifaces {
		if ifc.HardwareAddr != nil && ifc.Flags&net.FlagUp != 0 && ifc.Flags&net.FlagBroadcast != 0 {
			iface = ifc.Name
			break
		}
	}

	flag.StringVar(&iface, "interface", iface, "network interface to listen on")
	flag.Parse()

	pcapHandle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	handleError(fmt.Sprintf("failed to open %q for listening", iface), err)

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())
	for packet := range packetSource.Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeEthernet:
				eth, _ := layer.(*layers.Ethernet)
				var (
					srcManufName string = "?"
					dstManufName string = "?"
				)
				if srcRecord, ok := findMatch(records, eth.SrcMAC.String()); ok {
					srcManufName = srcRecord.OrganizationName
				}
				if dstRecord, ok := findMatch(records, eth.DstMAC.String()); ok {
					dstManufName = dstRecord.OrganizationName
				}

				log.Printf("%v (%v) -> %v (%v)", eth.SrcMAC, srcManufName, eth.DstMAC, dstManufName)
			}
		}
	}
}
