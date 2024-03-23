package cli

import "fmt"

type History struct {
	pages      []*Page
	currentIdx int
}

func (h *History) Back() error {
	if h.currentIdx == 0 {
		return fmt.Errorf("you cannot go back from your root page")
	}

	h.currentIdx = h.currentIdx - 1

	return nil
}

func (h *History) Forward() error {
	if h.currentIdx >= len(h.pages)-1 {
		return fmt.Errorf("you cannot go back from your root page")
	}

	h.currentIdx = h.currentIdx + 1

	return nil
}

func (h *History) Append(name string) {
	h.pages = append(h.pages, &Page{name: name})
	h.currentIdx = len(h.pages) - 1
}

func NewHistory(pages []*Page) *History {
	h := &History{
		pages:      pages,
		currentIdx: 0,
	}

	return h
}
