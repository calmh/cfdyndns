package main

import (
	"log"
	"net"
	"strings"

	"github.com/calmh/cfdns"
)

type dnsUpdater struct {
	changes <-chan net.IP
	client  *cfdns.Client
	name    string
	stop    chan struct{}
}

func newDNSUpdater(name string, client *cfdns.Client, changes <-chan net.IP) *dnsUpdater {
	return &dnsUpdater{
		changes: changes,
		client:  client,
		name:    name,
		stop:    make(chan struct{}),
	}
}

func (u *dnsUpdater) Serve() {
	var zoneID string
	zones, err := u.client.ListZones()
	if err != nil {
		log.Println("Listing zones:", err)
		return
	}

	for _, zone := range zones {
		if strings.HasSuffix(name, "."+zone.Name) {
			zoneID = zone.ID
			break
		}
	}

	if zoneID == "" {
		log.Println("No zone found for name", name, "- set manually using -zone")
		return
	}

	recs, err := u.client.ListDNSRecords(zoneID)
	if err != nil {
		log.Println("Listing records:", err)
		return
	}

	var curRec cfdns.DNSRecord
	for _, rec := range recs {
		if rec.Name == name {
			curRec = rec
			break
		}
	}

	if curRec.ID != "" {
		log.Println("Current DNS IP is", curRec.Content)
	}

	for {
		select {
		case newIP := <-u.changes:
			ip := newIP.String()
			if curRec.ID == "" {
				err := u.client.CreateDNSRecord(zoneID, name, "A", ip)
				if err != nil {
					log.Println("Creating record:", err)
					return
				}
				log.Println("Created record", name, "->", ip)
			} else {
				curRec.Content = ip
				err := u.client.UpdateDNSRecord(curRec)
				if err != nil {
					log.Println("Updating record:", err)
					return
				}
				log.Println("Updated record", name, "->", ip)
			}

		case <-u.stop:
			return
		}
	}
}

func (u *dnsUpdater) Stop() {
	close(u.stop)
}
