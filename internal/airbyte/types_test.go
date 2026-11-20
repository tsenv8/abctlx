package airbyte

import (
	"encoding/json"
	"strings"
	"testing"
)

// A streams-only PATCH must carry every existing stream unchanged (including
// fields abctlx does not model) and must not send schedule or other fields.
func TestUpdateStreamsBodyKeepsRawStreamsAndOmitsRest(t *testing.T) {
	in := `{"configurations":{"streams":[{"name":"a","namespace":"md","syncMode":"incremental_deduped_history","primaryKey":[["id"]],"mappers":[{"x":1}]}]}}`
	var con ConnectionData
	if err := json.Unmarshal([]byte(in), &con); err != nil {
		t.Fatal(err)
	}
	if len(con.Configurations.Streams) != 1 {
		t.Fatalf("streams not read: %d", len(con.Configurations.Streams))
	}
	out, _ := json.Marshal(UpdateConnectionRequest{Configurations: &StreamConfigurations{Streams: con.Configurations.Streams}})
	if !strings.Contains(string(out), `"mappers":[{"x":1}]`) {
		t.Fatalf("unmodelled field dropped: %s", out)
	}
	for _, bad := range []string{"schedule", "status", "name\":\"\""} {
		if strings.Contains(string(out), bad) {
			t.Fatalf("PATCH carries %q: %s", bad, out)
		}
	}
}
