package airbyte

import (
	"abctlx/internal/config"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type fakeAirbyte struct {
	health, token, workspaces int // status codes
	conns                     string
}

func (f fakeAirbyte) serve(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	reply := func(code int, body string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write([]byte(body))
		}
	}
	mux.HandleFunc("/api/public/v1/health", reply(f.health, `{"available":true}`))
	mux.HandleFunc("/api/public/v1/applications/token", reply(f.token, `{"access_token":"tok"}`))
	mux.HandleFunc("/api/public/v1/workspaces", reply(f.workspaces, `{"data":[{"WorkspaceId":"w1","Name":"default"}]}`))
	mux.HandleFunc("/api/public/v1/connections", reply(200, f.conns))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// profileFor writes a profiles file for srv and points ABCTLX_CONFIG at it.
func profileFor(t *testing.T, srv *httptest.Server, pass string, mode os.FileMode) {
	t.Helper()
	u, _ := url.Parse(srv.URL)
	port, _ := strconv.Atoi(u.Port())
	sealed, err := config.Seal("secret", "right")
	if err != nil {
		t.Fatal(err)
	}
	st := config.Store{Active: "p", Profiles: map[string]config.Profile{"p": {
		URL: "http://" + u.Hostname(), Port: port, APIURL: "/api/public", ClientID: "cid", Secret: sealed,
	}}}
	path := filepath.Join(t.TempDir(), "profiles.json")
	b, _ := json.Marshal(st)
	if err := os.WriteFile(path, b, mode); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ABCTLX_CONFIG", path)
	t.Setenv("ABCTLX_PASSPHRASE", pass)
	config.ForgetPassphrase()
}

func statuses(cs []Check) map[string]CheckStatus {
	m := map[string]CheckStatus{}
	for _, c := range cs {
		m[c.Name] = c.Status
	}
	return m
}

func TestDoctorAllGood(t *testing.T) {
	srv := fakeAirbyte{200, 200, 200, `{"data":[{"name":"a","tags":[{"name":"client"}]},{"name":"b","tags":[{"name":"client"}]}]}`}.serve(t)
	profileFor(t, srv, "right", 0o600)
	cs := Doctor(context.Background(), "")
	for _, c := range cs {
		if c.Status != CheckOK {
			t.Errorf("%+v", c)
		}
	}
	if got := cs[len(cs)-1]; got.Name != "connections" || got.Detail != "2 connections (tags: client=2)" {
		t.Errorf("%+v", got)
	}
}

func TestDoctorWrongPassphrase(t *testing.T) {
	srv := fakeAirbyte{200, 200, 200, `{"data":[]}`}.serve(t)
	profileFor(t, srv, "wrong", 0o600)
	s := statuses(Doctor(context.Background(), ""))
	if s["client secret"] != CheckFail || s["authentication"] != CheckSkip || s["workspace"] != CheckSkip {
		t.Errorf("%v", s)
	}
	if s["api reachable"] != CheckOK {
		t.Errorf("reachability does not need the secret: %v", s)
	}
}

func TestDoctorBadCredentials(t *testing.T) {
	srv := fakeAirbyte{200, 401, 200, `{"data":[]}`}.serve(t)
	profileFor(t, srv, "right", 0o600)
	s := statuses(Doctor(context.Background(), ""))
	if s["authentication"] != CheckFail || s["workspace"] != CheckSkip || s["connections"] != CheckSkip {
		t.Errorf("%v", s)
	}
}

func TestDoctorAPIDown(t *testing.T) {
	srv := fakeAirbyte{503, 200, 200, `{"data":[]}`}.serve(t)
	profileFor(t, srv, "right", 0o600)
	s := statuses(Doctor(context.Background(), ""))
	if s["api reachable"] != CheckFail || s["authentication"] != CheckSkip {
		t.Errorf("%v", s)
	}
}

func TestDoctorNoConfigFile(t *testing.T) {
	t.Setenv("ABCTLX_CONFIG", filepath.Join(t.TempDir(), "none.json"))
	s := statuses(Doctor(context.Background(), ""))
	if s["config file"] != CheckFail || s["profile"] != CheckSkip {
		t.Errorf("%v", s)
	}
}

func TestDoctorNoActiveProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.json")
	os.WriteFile(path, []byte(`{"profiles":{}}`), 0o600)
	t.Setenv("ABCTLX_CONFIG", path)
	s := statuses(Doctor(context.Background(), ""))
	if s["config file"] != CheckOK || s["profile"] != CheckFail {
		t.Errorf("%v", s)
	}
}

func TestDoctorLooseFileMode(t *testing.T) {
	srv := fakeAirbyte{200, 200, 200, `{"data":[]}`}.serve(t)
	profileFor(t, srv, "right", 0o644)
	if s := statuses(Doctor(context.Background(), "")); s["config file"] != CheckWarn {
		t.Errorf("%v", s)
	}
}
