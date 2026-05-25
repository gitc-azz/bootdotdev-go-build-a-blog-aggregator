package main

import (
	"reflect"
	"testing"
)

func TestUnescapeString(t *testing.T) {
	tcs := map[string]struct {
		inout, expected RSSFeed
	}{
		"simple": {
			inout: RSSFeed{
				Channel: struct {
					Title       string    "xml:\"title\""
					Link        string    "xml:\"link\""
					Description string    "xml:\"description\""
					Item        []RSSItem "xml:\"item\""
				}{
					Title:       "&lt",
					Link:        "",
					Description: "",
					Item: []RSSItem{
						RSSItem{
							Title:       "&lt",
							Link:        "",
							Description: "",
							PubDate:     "",
						},
					},
				},
			},
			expected: RSSFeed{
				Channel: struct {
					Title       string    "xml:\"title\""
					Link        string    "xml:\"link\""
					Description string    "xml:\"description\""
					Item        []RSSItem "xml:\"item\""
				}{
					Title:       "<",
					Link:        "",
					Description: "",
					Item: []RSSItem{
						RSSItem{
							Title:       "<",
							Link:        "",
							Description: "",
							PubDate:     "",
						},
					},
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			unescapeString(&tc.inout)
			if !reflect.DeepEqual(tc.inout, tc.expected) {
				t.Fatalf("actual: %v, expected: %v", tc.inout, tc.expected)
			}
		})
	}
}
