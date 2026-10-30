package main

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReadyPort(t *testing.T) {
	a, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	if a.Addr().String() == b.Addr().String() {
		t.Fatal("ports must differ")
	}
	path := filepath.Join(t.TempDir(), "ready")
	port := a.Addr().(*net.TCPAddr).Port
	if err := writeReadyPort(path, port); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != strconv.Itoa(port) {
		t.Fatalf("handshake: %q %v", data, err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temporary file remains")
	}
	for _, invalid := range []int{0, -1, 65536} {
		if writeReadyPort(path, invalid) == nil {
			t.Fatal("invalid port accepted")
		}
	}
	if writeReadyPort(filepath.Join(path, "missing"), port) == nil {
		t.Fatal("write failure ignored")
	}
}
