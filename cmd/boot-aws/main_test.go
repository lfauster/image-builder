package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResourceFileCompatibility(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		wantIDs []string
	}{
		{
			name:    "single instance",
			content: `{"region":"us-east-1","security-group":"sg-test","instance":"i-old"}`,
			wantIDs: []string{"i-old"},
		},
		{
			name:    "two instances",
			content: `{"region":"us-east-1","security-group":"sg-test","instances":["i-1","i-2"]}`,
			wantIDs: []string{"i-1", "i-2"},
		},
		{
			name:    "partial launch",
			content: `{"region":"us-east-1","security-group":"sg-test","instances":["i-1"]}`,
			wantIDs: []string{"i-1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var res resources
			if err := json.Unmarshal([]byte(tc.content), &res); err != nil {
				t.Fatal(err)
			}
			if got := res.allInstanceIDs(); !reflect.DeepEqual(got, tc.wantIDs) {
				t.Fatalf("instance IDs: got %v, want %v", got, tc.wantIDs)
			}

			path := filepath.Join(t.TempDir(), "resources.json")
			if err := saveResources(path, &res); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var restored resources
			if err := json.Unmarshal(data, &restored); err != nil {
				t.Fatal(err)
			}
			if got := restored.allInstanceIDs(); !reflect.DeepEqual(got, tc.wantIDs) {
				t.Fatalf("saved instance IDs: got %v, want %v", got, tc.wantIDs)
			}
		})
	}
}
