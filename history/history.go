package history

import (
	"encoding/json"
	"io"
	"os"
)

type History interface {
	io.Closer
	Add(id, description string)
	Exists(id string) bool
}

type fileHistory struct {
	path string
	ids  map[string]string
}

func Load(path string) (History, error) {
	h := &fileHistory{path: path, ids: make(map[string]string)}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return h, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&h.ids)
	return h, err
}

func (h *fileHistory) Close() error {
	f, err := os.Create(h.path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(h.ids)
}

func (h *fileHistory) Add(id, description string) {
	h.ids[id] = description
}

func (h *fileHistory) Exists(id string) bool {
	_, set := h.ids[id]
	return set
}
