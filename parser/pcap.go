package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type packet struct {
	t       time.Time
	srcAddr string
}

func NewMonChan(pcapFile string) (<-chan packet, error) {
	if _, err := os.Stat(pcapFile); os.IsNotExist(err) {
		return nil, err
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(pcapFile, ".gz") {
		str := "zcat " + pcapFile + " | /usr/sbin/tcpdump -tt -e -n -r -"
		cmd = exec.Command("/bin/bash", "-c", str)
	} else {
		cmd = exec.Command("/usr/sbin/tcpdump", "-tt", "-e", "-n", "-r", pcapFile)
	}
	if cmd == nil {
		return nil, fmt.Errorf("no command for %q", pcapFile)
	}
	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	pktc := make(chan packet, 16)
	go func() {
		defer close(pktc)
		r := bufio.NewReader(rc)
		for {
			l, lerr := r.ReadString('\n')
			if lerr != nil {
				break
			}
			pkt := line2pkt(l)
			if pkt == nil {
				continue
			}
			pktc <- *pkt
		}
		cmd.Wait()
		rc.Close()
	}()
	return pktc, nil
}

func line2pkt(s string) *packet {
	splt := "SA:"
	switch {
	case strings.Contains(s, "SA:"):
	case strings.Contains(s, "RA:"):
		splt = "RA:"
	default:
		return nil
	}
	sa := strings.SplitN(strings.SplitN(s, splt, 2)[1], " ", 2)[0]
	t := strings.SplitN(s, " ", 2)[0]
	f, _ := strconv.ParseFloat(t, 64)
	return &packet{
		t:       time.Unix(int64(f), int64((f-float64(int64(f)))*1e9)),
		srcAddr: sa,
	}
}
