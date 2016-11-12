package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s2"
)

type gpsUpdate struct {
	t  time.Time
	pt s2.LatLng
}

func NewGPSChan(gpsPath string) <-chan gpsUpdate {
	f, err := os.Open(gpsPath)
	if err != nil {
		panic(err)
	}
	gpsc := make(chan gpsUpdate)
	go func() {
		defer close(gpsc)
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			l := readType(r, "$GPRMC")
			if len(l) == 0 {
				break
			}
			if u := gprmc2update(l); u != nil {
				gpsc <- *u
			}
		}
	}()
	return gpsc
}

func readType(r *bufio.Reader, ty string) string {
	for {
		l, lerr := r.ReadString('\n')
		if lerr != nil {
			return ""
		}
		if strings.HasPrefix(l, ty) {
			return l
		}
	}
	return ""
}

func gprmc2update(s string) *gpsUpdate {
	// $GPRMC,050322.00,A,1234.56789,N,12345.67890,W,6.310,302.06,051116,,,D*74
	fields := strings.Split(s, ",")
	// 050322.00 => 05:03:22
	t := fields[1]
	hr, min, sec := t[0:2], t[2:4], t[4:6]

	hrv, _ := strconv.ParseInt(hr, 10, 32)
	minv, _ := strconv.ParseInt(min, 10, 32)
	secv, _ := strconv.ParseInt(sec, 10, 32)

	// 051116 => November 5th, 2016
	d := fields[9]
	day, mon, year := d[0:2], d[2:4], d[4:6]
	dayv, _ := strconv.ParseInt(day, 10, 32)
	monv, _ := strconv.ParseInt(mon, 10, 32)
	yearv, _ := strconv.ParseInt(year, 10, 32)

	tv := time.Date(
		int(yearv)+2000,
		time.Month(monv),
		int(dayv),
		int(hrv),
		int(minv),
		int(secv),
		0,
		time.UTC)

	// DDMM.MMMMM,N = latitude
	lat, latdir := fields[3], fields[4]
	if len(lat) == 0 {
		return nil
	}

	latdeg, _ := strconv.ParseInt(lat[0:2], 10, 32)
	latmin, _ := strconv.ParseFloat(lat[2:], 64)
	latv := float64(latdeg) + (latmin / 60.0)
	if latdir == "S" {
		latv = -latv
	}

	// DDDMM.MMMMM,W = lnggitude
	lng, lngdir := fields[5], fields[6]
	if len(lng) == 0 {
		return nil
	}
	lngdeg, _ := strconv.ParseInt(lng[0:3], 10, 32)
	lngmin, _ := strconv.ParseFloat(lng[3:], 64)
	lngv := float64(lngdeg) + (lngmin / 60.0)
	if lngdir == "W" {
		lngv = -lngv
	}

	return &gpsUpdate{tv, s2.LatLngFromDegrees(latv, lngv)}
}
