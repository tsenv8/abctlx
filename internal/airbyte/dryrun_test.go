package airbyte

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"

	"abctlx/internal/config"
)

func dryRunClient(t *testing.T) (AirbyteClient, *[]string, *int) {
	t.Helper()
	var mu sync.Mutex
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.Write([]byte("{}"))
	}))
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	port, _ := strconv.Atoi(u.Port())

	exits := 0
	oldExit, oldDry := dryRunExit, DryRun
	dryRunExit = func() { exits++ }
	DryRun = true
	t.Cleanup(func() { dryRunExit, DryRun = oldExit, oldDry })

	return NewAirbyteClient(config.AirbyteConfig{URL: "http://" + u.Hostname(), Port: port, API_URL: "/api/public"}), &seen, &exits
}

func TestDryRunSkipsMutations(t *testing.T) {
	for _, m := range []string{http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete} {
		c, seen, exits := dryRunClient(t)
		res, err := c.Request(context.Background(), m, CONNECTION_ENDPOINT, map[string]string{"a": "b"}, nil)
		if !errors.Is(err, ErrDryRun) || res != nil {
			t.Fatalf("%s: res=%v err=%v", m, res, err)
		}
		if len(*seen) != 0 {
			t.Fatalf("%s reached the server: %v", m, *seen)
		}
		if *exits != 1 {
			t.Fatalf("%s: exit hook called %d times", m, *exits)
		}
	}
}

func TestDryRunStillReadsAndAuthenticates(t *testing.T) {
	c, seen, exits := dryRunClient(t)
	if _, err := c.Request(context.Background(), http.MethodGet, CONNECTION_ENDPOINT, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Request(context.Background(), http.MethodPost, GENERATE_ACCESS_TOKEN_ENDPOINT, map[string]string{"client_id": "x"}, nil); err != nil {
		t.Fatal(err)
	}
	if len(*seen) != 2 || *exits != 0 {
		t.Fatalf("seen=%v exits=%d", *seen, *exits)
	}
}

func TestMutationsSendWithoutDryRun(t *testing.T) {
	c, seen, _ := dryRunClient(t)
	DryRun = false
	if _, err := c.Request(context.Background(), http.MethodPatch, CONNECTION_ENDPOINT+"/x", map[string]string{"a": "b"}, nil); err != nil {
		t.Fatal(err)
	}
	if len(*seen) != 1 {
		t.Fatalf("seen=%v", *seen)
	}
}
