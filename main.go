package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/calmh/cfdns"
	"github.com/thejerf/suture"
)

var (
	getIPCommand = "curl -s4 http://icanhazip.com/"
	zoneID       = ""
	name         = ""
	authEmail    = ""
	authKey      = ""
	ttl          = 300
	interval     = 300 * time.Second
	ipExp        = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
)

func main() {
	flag.StringVar(&getIPCommand, "cmd", getIPCommand, "Command to get external IP")
	flag.StringVar(&zoneID, "zone", zoneID, "Cloudflare Zone ID")
	flag.StringVar(&name, "name", name, "DNS record to update")
	flag.StringVar(&authEmail, "email", authEmail, "Cloudflare Auth Email")
	flag.StringVar(&authKey, "key", authKey, "Cloudflare API Key")
	flag.IntVar(&ttl, "ttl", ttl, "DNS record TTL (seconds)")
	flag.DurationVar(&interval, "intv", interval, "External IP check interval")
	flag.Parse()

	if name == "" {
		fmt.Println("Option -name is mandatory")
		os.Exit(1)
	}
	if authEmail == "" {
		fmt.Println("Option -email is mandatory")
		os.Exit(1)
	}
	if authKey == "" {
		fmt.Println("Option -key is mandatory")
		os.Exit(1)
	}

	client := cfdns.NewClient(authEmail, authKey)
	main := suture.NewSimple("main")
	checker := newIPChecker(getIPCommand, interval)
	main.Add(checker)
	updater := newDNSUpdater(name, client, checker.changes)
	main.Add(updater)
	main.Serve()
}
