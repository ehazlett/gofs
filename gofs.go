package gofs

import (
	"fmt"
	"io"
)

type Handler interface {
	OpenRW() (io.ReadWriteSeeker, error)
	OpenRO() (io.ReadSeeker, error)
	OpenWO() (io.WriteSeeker, error)

	IsDir() bool
	ListDir() ([]string, error)
	WalkDir(parts ...string) (Handler, error)
	Create(dir bool, parts ...string) (Handler, error)

	// FIXME: Remove()
	// FIXME: Stat()
	// FIXME: WStat()
}

type ChanHandler struct {
	ch chan []byte
}

func (h *ChanHandler) Read(data []byte) (int, error) {
	msg, ok := <-h.ch
	if !ok {
		return 0, fmt.Errorf("IO error")
	}
	n := copy(data, msg)
	return n, nil
}

func (h *ChanHandler) Write(data []byte) (int, error) {
	h.ch <- data
	return len(data), nil
}

func (h ChanHandler) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("can't seek")
}

func (h *ChanHandler) OpenRW() (io.ReadWriteSeeker, error) {
	return h, nil
}

func (h *ChanHandler) OpenRO() (io.ReadSeeker, error) {
	return h, nil
}

func (h *ChanHandler) OpenWO() (io.WriteSeeker, error) {
	return h, nil
}

func (h *ChanHandler) IsDir() bool {
	return false
}

func (h *ChanHandler) ListDir() ([]string, error) {
	return nil, fmt.Errorf("not a directory")
}

func (h *ChanHandler) WalkDir(parts ...string) (Handler, error) {
	return nil, fmt.Errorf("not a directory")
}

func (h *ChanHandler) Create(dir bool, parts ...string) (Handler, error) {
	return nil, fmt.Errorf("not a directory")
}

// FIXME: make thread-safe
type GatewayHandler map[string]*ChanHandler

func (h GatewayHandler) OpenRW() (io.ReadWriteSeeker, error) {
	return nil, fmt.Errorf("is a directory")
}

func (h GatewayHandler) OpenRO() (io.ReadSeeker, error) {
	return nil, fmt.Errorf("is a directory")
}

func (h GatewayHandler) OpenWO() (io.WriteSeeker, error) {
	return nil, fmt.Errorf("is a directory")
}

func (h GatewayHandler) IsDir() bool {
	return true
}

func (h GatewayHandler) ListDir() ([]string, error) {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys, nil
}

func (h GatewayHandler) WalkDir(parts ...string) (Handler, error) {
	if len(parts) == 0 {
		return h, nil
	}
	if len(parts) == 1 {
		if ch, ok := h[parts[0]]; ok {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("no such file or directory")
}

func (h GatewayHandler) Create(dir bool, parts ...string) (Handler, error) {
	if len(parts) != 1 {
		return nil, fmt.Errorf("can only create files at depth=1")
	}
	if dir == true {
		return nil, fmt.Errorf("permission denied: can't create a directory")
	}
	key := parts[0]
	ch, exists := h[key]
	if exists {
		return nil, fmt.Errorf("%s: already exists", parts[0])
	}
	h[key] = &ChanHandler{
		ch: make(chan []byte),
	}
	return ch, nil
}
