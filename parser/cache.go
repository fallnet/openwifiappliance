package main

import (
	"sync"
	"time"
)

type pktCache struct {
	mu    sync.Mutex
	cache map[string]*cacheMac
}

type cacheMac struct {
	first time.Time
	last  time.Time
	seen  int
}

func NewCache() *pktCache { return &pktCache{cache: make(map[string]*cacheMac)} }

func (pc *pktCache) Add(p packet) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	cm := pc.cache[p.srcAddr]
	if cm == nil {
		cm = &cacheMac{first: p.t, last: p.t}
		pc.cache[p.srcAddr] = cm
	}
	if p.t.After(cm.last) {
		cm.last = p.t
	}
	if cm.first.After(p.t) {
		cm.first = p.t
	}
	cm.seen++
}
