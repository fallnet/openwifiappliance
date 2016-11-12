package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/golang/geo/s2"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: parser <scandir>")
		os.Exit(1)
	}

	d := os.Args[1]
	ifs := ifaces(d)
	if len(ifs) == 0 {
		fmt.Printf("no ifaces found in %q\n", d)
		os.Exit(1)
	}

	gpsc := NewGPSChan(path.Join(d, "gps", "nmea.log"))

	pc := NewCache()

	var wg sync.WaitGroup
	wg.Add(len(ifs) + 1)
	for _, c := range ifs {
		go func(pktc <-chan packet) {
			defer wg.Done()
			for pkt := range pktc {
				pc.Add(pkt)
			}
		}(c)
	}

	gpsm := make(map[time.Time]s2.LatLng)
	go func() {
		defer wg.Done()
		for g := range gpsc {
			gpsm[g.t.Local()] = g.pt
		}

	}()
	wg.Wait()

	found := 0
	gps2macs := make(map[s2.LatLng][]string)
	for k, v := range pc.cache {
		t := v.first.Truncate(time.Second)
		if pt, ok := gpsm[t]; ok {
			gps2macs[pt] = append(gps2macs[pt], k)
			found++
		}
	}

	for pt, macs := range gps2macs {
		fmt.Println(pt.Lat, pt.Lng, strings.Join(macs, ","))
	}

	//	fmt.Println("missed:", len(pc.cache) - found, "of", len(pc.cache))
}

func ifaces(d string) (ret []<-chan packet) {
	fi, err := ioutil.ReadDir(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, info := range fi {
		if !strings.HasPrefix(info.Name(), "wlx") {
			continue
		}
		ret = append(ret, NewIFace(path.Join(d, info.Name())))
	}
	return ret
}
