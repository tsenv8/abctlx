package airbyte

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPlanAddStreams(t *testing.T) {
	existing := []json.RawMessage{json.RawMessage(`{"name":"keep","namespace":"md","syncMode":"incremental_deduped_history","mappers":[{"x":1}]}`)}
	all := []string{"full_refresh_append", "full_refresh_overwrite", "incremental_append", "incremental_deduped_history"}
	avail := []AvailableStream{
		{StreamName: "keep", StreamNamespace: "md", SyncModes: all},
		{StreamName: "room", StreamNamespace: "md", SyncModes: all, DefaultCursorField: []string{"_ab_cdc_lsn"}, SourceDefinedPrimaryKey: [][]string{{"id"}, {"branch_id"}}},
		{StreamName: "hist", StreamNamespace: "public", SyncModes: all[:2]},
	}
	res, out, err := PlanAddStreams(existing, avail, []string{"md.keep", "md.room", "public.hist", "md.nope", "md.room"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Already) != 1 || len(res.Missing) != 1 || len(res.Add) != 2 || len(out) != 3 {
		t.Fatalf("%+v out=%d", res, len(out))
	}
	if res.Add[0].SyncMode != SYNC_MODE_INCREMENTAL_DEDUPED_HISTORY || res.Add[1].SyncMode != SYNC_MODE_FULL_REFRESH_OVERWRITE {
		t.Fatalf("modes: %+v", res.Add)
	}
	if string(out[0]) != string(existing[0]) || !strings.Contains(string(out[0]), `"mappers":[{"x":1}]`) {
		t.Fatal("existing stream changed")
	}
	if !strings.Contains(string(out[1]), `"primaryKey":[["id"],["branch_id"]]`) || !strings.Contains(string(out[2]), `"primaryKey":[]`) {
		t.Fatalf("entries: %s %s", out[1], out[2])
	}
	if _, _, err := PlanAddStreams(nil, avail, []string{"room"}); err == nil {
		t.Fatal("bare table name accepted")
	}
}
