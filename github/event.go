package github

import (
	_ "embed"

	"github.com/starudream/go-lib/core/v2/codec/json"
)

type Event struct {
	Key     string        `json:"key"`
	Name    string        `json:"name"`
	Desc    string        `json:"desc"`
	Actions []EventAction `json:"actions"`
}

type EventAction struct {
	// Category        string   `json:"category"`
	Action string `json:"action"`
	// Availability    []string `json:"availability"`
	// DescriptionHtml string   `json:"descriptionHtml"`
	SummaryHtml string `json:"summaryHtml"`
}

// var _ tea.SelectItem = (*Event)(nil)

var (
	//go:embed events.json
	eventsRaw []byte

	Events = json.MustUnmarshalTo[[]Event](eventsRaw)
)
