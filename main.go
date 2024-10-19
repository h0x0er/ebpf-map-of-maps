package main

import (
	"log/slog"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 Mapper  ./bpf/maps.c -- -I./bpf/includes
func main() {

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		slog.Error("removeLocks", "error", err)
		os.Exit(1)
	}

	slog.Info("loading bpf objects")

	mapperObjs := new(MapperObjects)
	err := LoadMapperObjects(mapperObjs, nil)
	if err != nil {
		slog.Error("error loading mapper objects", "error", err)
		return
	}

	defer mapperObjs.Close()

	// Store section: for storing innner_map in outer_rmap

	slog.Info("creating innnerMap specs")
	// creating inner_map spec with upated max_entries
	innerMapSpecs := &ebpf.MapSpec{
		Type:       ebpf.LRUHash,
		KeySize:    4,
		ValueSize:  1,
		MaxEntries: 2,
	}

	slog.Info("creating inner_map")
	// creating inner map
	inner, err := ebpf.NewMap(innerMapSpecs)
	if err != nil {
		slog.Error("error creating inner map", "error", err)
		return
	}
	defer inner.Close()

	slog.Info("storing in inner_map", "key", 1, "val", 250)
	err = inner.Put(uint32(1), uint8(250))
	if err != nil {
		slog.Error("error putting into inner map", "error", err)
		return
	}

	slog.Info("putting inner_map in outer_map", "key", 1)
	// putting inner_map into outer_map
	err = mapperObjs.OuterMap.Put(uint32(1), inner)

	if err != nil {
		slog.Error("error putting into outer map", "error", err)
		return
	}

	// Lookup Section: for fetching inner_map in outer_map

	slog.Info("feching outer map for stored map", "key", 1)
	var val uint32
	err = mapperObjs.OuterMap.Lookup(uint32(1), &val)
	if err != nil {
		slog.Error("unable to lookup outer map", "error", err)
		return
	}

	slog.Info("storedMap", "id", val)

	storedMap, err := ebpf.NewMapFromID(ebpf.MapID(val))

	if err != nil {
		slog.Error("unable to create inner map", "error", err)
		return
	}

	var val2 uint8
	err = storedMap.Lookup(uint32(1), &val2)
	if err != nil {
		slog.Error("unable to lookup inner map", "error", err)
		return
	}

	slog.Info("storedMap", "val", val2)
}
