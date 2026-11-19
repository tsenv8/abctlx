// Package config holds abctlx profiles: one profile is one Airbyte instance.
// Profiles live in ~/.config/abctlx/profiles.json (0600). The client secret is
// stored encrypted (see secret.go); nothing is read from .env or ABCTLX_AB_*.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	DefaultPort      = 2499
	DefaultAPIURL    = "/api/public"
	DefaultNamespace = "airbyte-abctl"
)

// AirbyteConfig is what the Airbyte client needs. The secret is not here: it is
// decrypted on demand with ClientSecret, so commands that never talk to Airbyte
// (profile list, show) never ask for the passphrase.
type AirbyteConfig struct {
	URL        string
	ClientId   string
	Port       int
	API_URL    string
	Kubeconfig string
	Namespace  string
}

type Profile struct {
	URL        string  `json:"url"`
	Port       int     `json:"port"`
	APIURL     string  `json:"apiUrl"`
	ClientID   string  `json:"clientId"`
	Secret     *Sealed `json:"secret,omitempty"`
	Kubeconfig string  `json:"kubeconfig,omitempty"`
	Namespace  string  `json:"namespace,omitempty"`
}

type Store struct {
	Active   string             `json:"active,omitempty"`
	Profiles map[string]Profile `json:"profiles"`
}

// Data is the resolved config of the profile in use. Set by Resolve.
var Data AirbyteConfig

var current Profile
var currentName string

func Path() (string, error) {
	if p := os.Getenv("ABCTLX_CONFIG"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "abctlx", "profiles.json"), nil
}

func Load() (*Store, error) {
	st := &Store{Profiles: map[string]Profile{}}
	p, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return st, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, st); err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	if st.Profiles == nil {
		st.Profiles = map[string]Profile{}
	}
	return st, nil
}

func (s *Store) Save() error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func (s *Store) Names() []string {
	names := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Resolve picks the profile (flag, else the active one) and fills Data.
func Resolve(flag string) error {
	st, err := Load()
	if err != nil {
		return err
	}
	name := flag
	if name == "" {
		name = st.Active
	}
	if name == "" {
		return errors.New("no profile selected: run 'abctlx profile add <name> --url ... --client-id ...' then 'abctlx profile use <name>', or pass --profile")
	}
	p, ok := st.Profiles[name]
	if !ok {
		return fmt.Errorf("no profile named %q (have: %v)", name, st.Names())
	}
	current, currentName = p, name
	Data = AirbyteConfig{
		URL: p.URL, Port: p.Port, API_URL: p.APIURL, ClientId: p.ClientID,
		Kubeconfig: p.Kubeconfig, Namespace: p.Namespace,
	}
	return nil
}

func CurrentName() string { return currentName }

// ClientSecret decrypts the profile's client secret, asking for the passphrase
// ($ABCTLX_PASSPHRASE, else a prompt) the first time.
func ClientSecret() (string, error) {
	if current.Secret == nil {
		return "", fmt.Errorf("profile %s has no client secret: run 'abctlx profile secret %s'", currentName, currentName)
	}
	pass, err := Passphrase(false)
	if err != nil {
		return "", err
	}
	return current.Secret.Open(pass)
}
