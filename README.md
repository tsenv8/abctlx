# abctlx

Airbyte controls, extended. A CLI over Airbyte's public API.

## Profiles

One profile is one Airbyte instance. Nothing is read from `.env`.

```bash
abctlx profile add local --url http://192.168.68.133 --client-id <id> --use   # asks for the secret and a passphrase
abctlx profile test
abctlx profile list | show | use <name> | remove <name> | secret <name>
abctlx -P other connection        # --profile picks one for a single command
```

The client secret is stored AES-256-GCM encrypted in `~/.config/abctlx/profiles.json` (0600). The passphrase is never stored: set `$ABCTLX_PASSPHRASE` or type it once per run. The secret can also be given with `--secret` or `$ABCTLX_CLIENT_SECRET` instead of the prompt (a flag ends up in shell history).

## Connections

```bash
abctlx connection                                   # list
abctlx connection streams --id <conn>               # enabled streams
abctlx connection streams --id <conn> --available   # what the source offers
abctlx connection streams --id <conn> --json        # raw API response
```
