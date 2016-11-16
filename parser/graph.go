package main

import "fmt"

type macGraph map[string]map[string]struct{}

type macSubgraph map[string]struct{}

// subgraph is the reachable subgraph for a set of verices
func (mg macGraph) subgraph(s string, m map[string]struct{}) macSubgraph {
	if m == nil {
		m = make(map[string]struct{})
	}
	m[s] = struct{}{}
	// visit all endpoints x where s -> x
	for dst := range mg[s] {
		if _, ok := m[dst]; !ok {
			// not yet visited
			mg.subgraph(dst, m)
		}
	}
	// visit all endpoints x where x -> s
	for src, dsts := range mg {
		if _, ok := dsts[s]; !ok {
			// no path src -> s
			continue
		}
		if _, ok := m[src]; !ok {
			// not yet visited
			mg.subgraph(src, m)
		}
	}
	return m
}

func (mg macGraph) add(pkt packet) {
	if pkt.dstAddr == "" || pkt.srcAddr == "" {
		return
	}
	if pkt.dstAddr == "ff:ff:ff:ff:ff:ff" {
		return
	}
	m := mg[pkt.srcAddr]
	if m == nil {
		m = make(map[string]struct{})
		mg[pkt.srcAddr] = m
	}
	m[pkt.dstAddr] = struct{}{}
}

// Dot generates a dot file for a subgraph.
func (mg macGraph) Dot(subgraph macSubgraph) string {
	s := "digraph {\n"
	for mac := range subgraph {
		for dst := range mg[mac] {
			s += fmt.Sprintf("%q -> %q;\n", mac, dst)
		}
	}
	return s + "}\n"
}
