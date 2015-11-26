package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

type ipChecker struct {
	cmd      string
	changes  chan net.IP
	interval time.Duration
	stop     chan struct{}
}

func newIPChecker(cmd string, interval time.Duration) *ipChecker {
	return &ipChecker{
		cmd:      cmd,
		changes:  make(chan net.IP),
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (c *ipChecker) get() (net.IP, error) {
	cmd := exec.Command("/bin/sh", "-c", c.cmd)
	bs, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	bs = bytes.TrimSpace(bs)

	if !ipExp.Match(bs) {
		return nil, fmt.Errorf("%q does not look like an IP\n", bs)
	}

	return net.ParseIP(string(bs)), nil
}

func (c *ipChecker) Serve() {
	t := time.NewTimer(0)
	var curIP net.IP
	for {
		select {
		case <-t.C:
			ip, err := c.get()
			if err != nil {
				log.Println("Get IP:", err)
			}
			if !curIP.Equal(ip) {
				log.Println("Current external IP is", ip)
				curIP = ip
				c.changes <- ip
			}
			t.Reset(c.interval)

		case <-c.stop:
			t.Stop()
			return
		}
	}
}

func (c *ipChecker) Stop() {
	close(c.stop)
}
