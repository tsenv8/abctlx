package airbyte

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PlannedStream is one stream add-streams would enable.
type PlannedStream struct {
	Table      string // schema.table
	SyncMode   string
	Cursor     []string
	PrimaryKey [][]string
	Reason     string
}

// PlanResult is what add-streams would do, before anything is sent.
type PlanResult struct {
	Add     []PlannedStream
	Already []string // requested tables that are already enabled
	Missing []string // requested tables the source does not offer
}

// ChooseSyncMode applies the D3 rule to what the source reports for one stream:
// a source-defined primary key and incremental dedupe -> incremental_deduped_history
// on the CDC cursor; otherwise incremental append when offered; otherwise full
// refresh overwrite.
func ChooseSyncMode(a AvailableStream) (mode string, reason string) {
	has := func(m string) bool {
		for _, x := range a.SyncModes {
			if x == m {
				return true
			}
		}
		return false
	}
	switch {
	case len(a.SourceDefinedPrimaryKey) > 0 && has(SYNC_MODE_INCREMENTAL_DEDUPED_HISTORY):
		return SYNC_MODE_INCREMENTAL_DEDUPED_HISTORY, "has primary key"
	case has(SYNC_MODE_INCREMENTAL_APPEND):
		return SYNC_MODE_INCREMENTAL_APPEND, "no primary key, incremental offered"
	default:
		return SYNC_MODE_FULL_REFRESH_OVERWRITE, "source offers full refresh only"
	}
}

// PlanAddStreams works out which of the wanted "schema.table" names to enable.
// Existing streams are left untouched. It does not talk to Airbyte.
func PlanAddStreams(existing []json.RawMessage, avail []AvailableStream, wanted []string) (PlanResult, []json.RawMessage, error) {
	var res PlanResult
	enabled := map[string]bool{}
	for _, raw := range existing {
		var st StreamConfiguration
		if err := json.Unmarshal(raw, &st); err != nil {
			return res, nil, fmt.Errorf("existing stream: %w", err)
		}
		enabled[st.Namespace+"."+st.Name] = true
	}
	offered := map[string]AvailableStream{}
	for _, a := range avail {
		offered[a.StreamNamespace+"."+a.StreamName] = a
	}

	out := append([]json.RawMessage{}, existing...)
	seen := map[string]bool{}
	for _, w := range wanted {
		w = strings.TrimSpace(w)
		if seen[w] {
			continue
		}
		seen[w] = true
		if strings.Count(w, ".") != 1 {
			return res, nil, fmt.Errorf("%q: write tables as schema.table", w)
		}
		if enabled[w] {
			res.Already = append(res.Already, w)
			continue
		}
		a, ok := offered[w]
		if !ok {
			res.Missing = append(res.Missing, w)
			continue
		}
		mode, reason := ChooseSyncMode(a)
		var cursor []string
		var pk [][]string
		if mode != SYNC_MODE_FULL_REFRESH_OVERWRITE {
			cursor = a.DefaultCursorField
			pk = a.SourceDefinedPrimaryKey
		}
		if cursor == nil {
			cursor = []string{}
		}
		if pk == nil {
			pk = [][]string{}
		}
		// Same shape Airbyte returns for existing streams on this instance.
		entry, err := json.Marshal(map[string]any{
			"name": a.StreamName, "namespace": a.StreamNamespace,
			"syncMode": mode, "cursorField": cursor, "primaryKey": pk,
			"includeFiles": false, "selectedFields": []any{}, "mappers": []any{},
		})
		if err != nil {
			return res, nil, err
		}
		out = append(out, entry)
		res.Add = append(res.Add, PlannedStream{Table: w, SyncMode: mode, Cursor: cursor, PrimaryKey: pk, Reason: reason})
	}
	return res, out, nil
}
