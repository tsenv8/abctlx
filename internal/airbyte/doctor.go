package airbyte

import (
	"abctlx/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

type CheckStatus string

const (
	CheckOK   CheckStatus = "ok"
	CheckWarn CheckStatus = "warn"
	CheckFail CheckStatus = "FAIL"
	CheckSkip CheckStatus = "skip"
)

// Check is one line of the doctor report. Hint says what to do about a warn or
// a failure.
type Check struct {
	Name   string
	Status CheckStatus
	Detail string
	Hint   string
}

type doctor struct {
	checks []Check
}

func (d *doctor) add(name string, st CheckStatus, detail, hint string) bool {
	d.checks = append(d.checks, Check{name, st, detail, hint})
	return st != CheckFail
}

func (d *doctor) skip(names ...string) {
	for _, n := range names {
		d.add(n, CheckSkip, "needs an earlier check to pass", "")
	}
}

// Doctor checks the chosen profile (empty: the active one) from the file on
// disk to a working API call. It never exits: every problem becomes a failed
// check, and later checks that depend on it are skipped.
func Doctor(ctx context.Context, profile string) []Check {
	d := &doctor{}

	// 1. config file
	path, err := config.Path()
	if err != nil {
		d.add("config file", CheckFail, err.Error(), "")
		d.skip("profile", "profile fields", "client secret", "api reachable", "authentication", "workspace", "connections")
		return d.checks
	}
	info, err := os.Stat(path)
	switch {
	case err != nil:
		d.add("config file", CheckFail, path+": "+err.Error(), "run 'abctlx profile add <name> --url ... --client-id ... --use'")
		d.skip("profile", "profile fields", "client secret", "api reachable", "authentication", "workspace", "connections")
		return d.checks
	case info.Mode().Perm()&0o077 != 0:
		d.add("config file", CheckWarn, fmt.Sprintf("%s is mode %o; it holds an encrypted secret", path, info.Mode().Perm()), "chmod 600 "+path)
	default:
		d.add("config file", CheckOK, path, "")
	}

	// 2. profile
	if err := config.Resolve(profile); err != nil {
		d.add("profile", CheckFail, err.Error(), "")
		d.skip("profile fields", "client secret", "api reachable", "authentication", "workspace", "connections")
		return d.checks
	}
	d.add("profile", CheckOK, config.CurrentName(), "")

	// 3. profile fields
	cfg := config.Data
	var missing []string
	if cfg.URL == "" {
		missing = append(missing, "url")
	}
	if cfg.ClientId == "" {
		missing = append(missing, "clientId")
	}
	if cfg.Port == 0 {
		missing = append(missing, "port")
	}
	if cfg.API_URL == "" {
		missing = append(missing, "apiUrl")
	}
	if len(missing) > 0 {
		d.add("profile fields", CheckFail, "missing: "+strings.Join(missing, ", "), "re-run 'abctlx profile add "+config.CurrentName()+" ...'")
		d.skip("client secret", "api reachable", "authentication", "workspace", "connections")
		return d.checks
	}
	d.add("profile fields", CheckOK, NewAirbyteClient(cfg).GetURL(nil), "")

	// 4. secret
	secret, err := config.ClientSecret()
	secretOK := d.add("client secret", statusOf(err), detailOf(err, "decrypted"), secretHint(err))

	// 5. reachable (needs no secret)
	c := NewAirbyteClient(cfg)
	_, err = c.Request(ctx, http.MethodGet, HEALTH_CHECK_ENDPOINT, nil, nil)
	reachable := d.add("api reachable", statusOf(err), detailOf(err, "health 200"), "is Airbyte up, and are url and port right? ('abctlx profile show')")

	// 6. token
	var token string
	switch {
	case !secretOK || !reachable:
		d.skip("authentication")
	default:
		res, err := c.Request(ctx, http.MethodPost, GENERATE_ACCESS_TOKEN_ENDPOINT, GenerateAccessTokenRequest{ClientId: cfg.ClientId, ClientKey: secret}, nil)
		if err == nil {
			var tr GenerateAccessTokenResponse
			if err = json.Unmarshal(res.Body, &tr); err == nil && tr.AccessToken == "" {
				err = fmt.Errorf("Airbyte returned no access token")
			} else if err == nil {
				token = tr.AccessToken
			}
		}
		d.add("authentication", statusOf(err), detailOf(err, "token issued"), "wrong client id or secret? 'abctlx profile secret "+config.CurrentName()+"'")
	}
	if token == "" {
		d.skip("workspace", "connections")
		return d.checks
	}

	// 7. workspace
	var ws ListWorkspacesResponse
	res, err := c.Request(ctx, http.MethodGet, LIST_WORKSPACES_ENDPOINT, nil, &token)
	if err == nil {
		err = json.Unmarshal(res.Body, &ws)
	}
	switch {
	case err != nil:
		d.add("workspace", CheckFail, err.Error(), "")
	case len(ws.Data) == 0:
		d.add("workspace", CheckFail, "no workspace visible to this application", "give the application access to a workspace in Airbyte")
	default:
		d.add("workspace", CheckOK, fmt.Sprintf("%s (%d visible)", ws.Data[0].Name, len(ws.Data)), "")
	}

	// 8. connections
	var cons ListConnectionResponse
	res, err = c.Request(ctx, http.MethodGet, CONNECTION_ENDPOINT+"?limit=500", nil, &token)
	if err == nil {
		err = json.Unmarshal(res.Body, &cons)
	}
	if err != nil {
		d.add("connections", CheckFail, err.Error(), "")
		return d.checks
	}
	tagged := map[string]int{}
	for _, con := range cons.Data {
		for _, t := range con.Tags {
			tagged[t.Name]++
		}
	}
	detail := fmt.Sprintf("%d connections", len(cons.Data))
	if len(tagged) > 0 {
		var parts []string
		for name, n := range tagged {
			parts = append(parts, fmt.Sprintf("%s=%d", name, n))
		}
		sort.Strings(parts)
		detail += " (tags: " + strings.Join(parts, ", ") + ")"
	}
	if cons.Next != "" {
		d.add("connections", CheckWarn, detail+"; more than 500, only the first 500 counted", "")
	} else {
		d.add("connections", CheckOK, detail, "")
	}

	// 9. optional: kubeconfig, only used by 'res update-connection'
	if cfg.Kubeconfig != "" {
		if _, err := os.Stat(cfg.Kubeconfig); err != nil {
			d.add("kubeconfig", CheckWarn, err.Error(), "only 'res update-connection' uses it")
		} else {
			d.add("kubeconfig", CheckOK, cfg.Kubeconfig, "")
		}
	}
	return d.checks
}

func statusOf(err error) CheckStatus {
	if err != nil {
		return CheckFail
	}
	return CheckOK
}

func detailOf(err error, ok string) string {
	if err != nil {
		return err.Error()
	}
	return ok
}

func secretHint(err error) string {
	if err == nil {
		return ""
	}
	return "set $ABCTLX_PASSPHRASE, or run in a terminal to be asked; 'abctlx profile secret " + config.CurrentName() + "' replaces the secret"
}
