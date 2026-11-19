package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Airbyte profiles (one profile = one Airbyte instance)",
	Long: `A profile says how to reach one Airbyte instance: URL, port, API path and the
application client id. The client secret is stored encrypted under a passphrase
($ABCTLX_PASSPHRASE, or asked for once per run); the passphrase is never stored.
Profiles live in ~/.config/abctlx/profiles.json. Nothing is read from .env.`,
}

var (
	pfURL, pfAPIURL, pfClientID, pfSecret, pfKubeconfig, pfNamespace string
	pfPort                                                 int
	pfUse                                                  bool
)

func readSecret() (string, error) {
	if pfSecret != "" {
		return pfSecret, nil
	}
	if v := os.Getenv("ABCTLX_CLIENT_SECRET"); v != "" {
		return v, nil
	}
	return config.ReadHidden("Airbyte client secret: ")
}

// sealSecret asks for the secret and a passphrase and returns it encrypted.
func sealSecret() (*config.Sealed, error) {
	secret, err := readSecret()
	if err != nil {
		return nil, err
	}
	if secret == "" {
		return nil, fmt.Errorf("empty client secret")
	}
	pass, err := config.Passphrase(true)
	if err != nil {
		return nil, err
	}
	return config.Seal(secret, pass)
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add (or replace) a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		if pfURL == "" || pfClientID == "" {
			return fmt.Errorf("--url and --client-id are required")
		}
		sealed, err := sealSecret()
		if err != nil {
			return err
		}
		st.Profiles[args[0]] = config.Profile{
			URL: pfURL, Port: pfPort, APIURL: pfAPIURL, ClientID: pfClientID,
			Secret: sealed, Kubeconfig: pfKubeconfig, Namespace: pfNamespace,
		}
		if st.Active == "" || pfUse {
			st.Active = args[0]
		}
		if err := st.Save(); err != nil {
			return err
		}
		fmt.Printf("Saved profile %s (active: %s).\n", args[0], st.Active)
		return nil
	},
}

var profileSecretCmd = &cobra.Command{
	Use:   "secret <name>",
	Short: "Replace the stored client secret (and passphrase) of a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		p, ok := st.Profiles[args[0]]
		if !ok {
			return fmt.Errorf("no profile named %q", args[0])
		}
		if p.Secret, err = sealSecret(); err != nil {
			return err
		}
		st.Profiles[args[0]] = p
		return st.Save()
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		var body []table.Row
		for _, n := range st.Names() {
			p := st.Profiles[n]
			mark := ""
			if n == st.Active {
				mark = "*"
			}
			body = append(body, table.Row{mark, n, p.URL + ":" + strconv.Itoa(p.Port) + p.APIURL, p.ClientID})
		}
		helpers.Table(table.Row{"", "Name", "Endpoint", "Client ID"}, body, "Profiles")
		return nil
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show a profile (the secret is never printed)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		name := st.Active
		if len(args) == 1 {
			name = args[0]
		}
		p, ok := st.Profiles[name]
		if !ok {
			return fmt.Errorf("no profile named %q", name)
		}
		secret := "not set"
		if p.Secret != nil {
			secret = "stored, encrypted"
		}
		helpers.Table(table.Row{"Field", "Value"}, []table.Row{
			{"Name", name}, {"URL", p.URL}, {"Port", p.Port}, {"API path", p.APIURL},
			{"Client ID", p.ClientID}, {"Client secret", secret},
			{"Kubeconfig", p.Kubeconfig}, {"Namespace", p.Namespace},
		}, "Profile")
		return nil
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Make a profile the active one",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		if _, ok := st.Profiles[args[0]]; !ok {
			return fmt.Errorf("no profile named %q", args[0])
		}
		st.Active = args[0]
		return st.Save()
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := config.Load()
		if err != nil {
			return err
		}
		if _, ok := st.Profiles[args[0]]; !ok {
			return fmt.Errorf("no profile named %q", args[0])
		}
		delete(st.Profiles, args[0])
		if st.Active == args[0] {
			st.Active = ""
		}
		return st.Save()
	},
}

var profileTestCmd = &cobra.Command{
	Use:   "test [name]",
	Short: "Check health and authentication of a profile",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := flagProfile
		if len(args) == 1 {
			name = args[0]
		}
		if err := config.Resolve(name); err != nil {
			return err
		}
		s := airbyte.NewAirbyteService(context.Background())
		if !s.Health().Status {
			return fmt.Errorf("profile %s: Airbyte is not healthy at %s", config.CurrentName(), s.GetClient().GetURL(nil))
		}
		s.GetAccessToken() // exits with the API error if the id/secret are wrong
		fmt.Printf("\nProfile %s: healthy and authenticated.\n", config.CurrentName())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd, profileSecretCmd, profileListCmd, profileShowCmd, profileUseCmd, profileRemoveCmd, profileTestCmd)
	f := profileAddCmd.Flags()
	f.StringVar(&pfURL, "url", "", "Airbyte URL without port, e.g. http://192.168.68.133")
	f.IntVar(&pfPort, "port", config.DefaultPort, "Airbyte port")
	f.StringVar(&pfAPIURL, "api-url", config.DefaultAPIURL, "API path prefix")
	f.StringVar(&pfClientID, "client-id", "", "application client id")
	f.StringVar(&pfSecret, "secret", "", "client secret (default: $ABCTLX_CLIENT_SECRET, else a hidden prompt; a flag lands in shell history)")
	profileSecretCmd.Flags().StringVar(&pfSecret, "secret", "", "client secret (default: $ABCTLX_CLIENT_SECRET, else a hidden prompt)")
	f.StringVar(&pfKubeconfig, "kubeconfig", "", "kubeconfig (only for 'res update-connection')")
	f.StringVar(&pfNamespace, "namespace", config.DefaultNamespace, "Kubernetes namespace (only for 'res update-connection')")
	f.BoolVar(&pfUse, "use", false, "make it the active profile")
}
