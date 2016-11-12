package main

import (
	"io/ioutil"
	"path"
	"strings"
)

func NewIFace(ifaceDir string) <-chan packet {
	fi, err := ioutil.ReadDir(ifaceDir)
	if err != nil {
		panic(err)
	}
	var paths []string
	for _, info := range fi {
		if !strings.HasPrefix(info.Name(), "pcap") {
			continue
		}
		paths = append(paths, path.Join(ifaceDir, info.Name()))
	}
	// XXX paths not in order
	pktc := make(chan packet, 16*len(paths))
	go func() {
		defer close(pktc)
		for _, p := range paths {
			curpktc, perr := NewMonChan(p)
			if perr != nil {
				return
			}
			for pkt := range curpktc {
				pktc <- pkt
			}
		}
	}()
	return pktc
}
