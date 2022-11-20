package util

import (
	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func ThreadPrevURI(thread *bakusai.Thread) string {
	prevThread, err := parser.ParseThreadURI(thread.PrevURI)
	if err != nil {
		return ""
	}

	return prevThread.URI()
}

func ThreadNextURI(thread *bakusai.Thread) string {
	nextThread, err := parser.ParseThreadURI(thread.NextURI)
	if err != nil {
		return ""
	}

	return nextThread.URI()
}
