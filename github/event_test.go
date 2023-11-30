package github

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/utils/testutil"
	"github.com/starudream/go-lib/resty/v2"
)

func TestEvents(t *testing.T) {
	for i := 0; i < len(Events); i++ {
		e := Events[i]
		t.Logf("%02d: %-40s %-50s %s", i+1, e.Key, e.Name, e.Desc)
		for j := 0; j < len(e.Actions); j++ {
			a := e.Actions[j]
			t.Logf("\t\t%02d: %s", j+1, a.Action)
		}
	}
}

// save `https://github.com/owner/repo/settings/hooks/new` page source to current folder
func TestGenEvents(t *testing.T) {
	file, err := os.Open("Add webhook.html")
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("Add webhook.html not found")
		} else {
			t.Fatal(err)
		}
	}
	defer gh.Close(file)

	resp, err := resty.R().Get("https://raw.githubusercontent.com/github/docs/main/src/webhooks/data/fpt/schema.json")
	testutil.Nil(t, err)
	testutil.Equal(t, 200, resp.StatusCode())

	mm, err := json.UnmarshalTo[map[string]map[string]EventAction](resp.Body())
	testutil.LogNoErr(t, err, mm)

	root, err := html.Parse(file)
	testutil.Nil(t, err)

	nodes := NodeSearch(root)
	events := make([]Event, len(nodes))

	for i, node := range nodes {
		// NodePrint(node)
		events[i] = Event{
			Key:  NodeText(node, "input", "value"),
			Name: NodeText(node, "input"),
			Desc: NodeText(node, "p"),
		}
		var actions []EventAction
		for _, action := range mm[events[i].Key] {
			actions = append(actions, action)
		}
		sort.Slice(actions, func(i, j int) bool { return actions[i].Action < actions[j].Action })
		events[i].Actions = actions
	}

	err = os.WriteFile("events.json", json.MustMarshalIndent(events), 0644)
	testutil.Nil(t, err)
}

func NodeSearch(node *html.Node) (nodes []*html.Node) {
	if node == nil {
		return nil
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		for _, attr := range child.Attr {
			if attr.Key == "class" {
				for _, v := range strings.Split(attr.Val, " ") {
					if v == "hook-event" {
						nodes = append(nodes, child)
					}
				}
			}
		}
		if ns := NodeSearch(child); len(ns) > 0 {
			nodes = append(nodes, ns...)
		}
	}
	return
}

func NodeText(node *html.Node, tag string, attr ...string) string {
	if node == nil {
		return ""
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			switch tag {
			case "input":
				if len(attr) > 0 {
					for _, a := range child.Attr {
						if a.Key == attr[0] {
							return a.Val
						}
					}
				}
				return strings.TrimSpace(child.NextSibling.Data)
			case "p":
				return strings.TrimSpace(child.FirstChild.Data)
			}
		}
		if s := NodeText(child, tag, attr...); s != "" {
			return s
		}
	}
	return ""
}

func NodePrint(node *html.Node, ii ...int) {
	if node == nil {
		return
	}
	i := 0
	if len(ii) > 0 {
		i = ii[0]
	}
	if i == 0 {
		fmt.Println("----- >>> -----")
		fmt.Printf("%d   %d %q %v\n", i, node.Type, node.Data, node.Attr)
	}
	j := 0
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		j++
		fmt.Printf("%d:%d %d %q %v\n", i, j, child.Type, child.Data, child.Attr)
		NodePrint(child, i+1)
	}
}
