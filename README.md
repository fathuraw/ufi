# ufi

**UniFi** from the terminal. Manage your UDM, UDR, or UCG — devices, clients, networks, firewall, WiFi, DNS, ACLs — without touching the web UI.

Built on the [UniFi Integration API](https://developers.ui.com/) (`/integration/v1`).

## Install

```bash
go install github.com/fathuraw/ufi@latest
```

Or from source:

```bash
git clone git@github.com:fathuraw/ufi.git && cd ufi
go build -o ufi .
```

## Quick start

```bash
# 1. Generate an API key in your UniFi console:
#    Settings > System > Advanced > API Keys

# 2. Login (stores key in your system keychain)
ufi login

# 3. Go
ufi site list
ufi device list
ufi client list
```

Or skip login entirely:

```bash
export UFI_API_KEY=your-key
ufi --host https://192.168.1.1 device list
```

## Commands

### Sites & devices

```bash
ufi site list                                      # list all sites

ufi device list                                     # list devices
ufi device get <id>                                 # device details
ufi device stats <id>                               # CPU, memory, traffic
ufi device restart <id>                             # reboot a device
ufi device adopt <mac>                              # adopt by MAC
ufi device remove <id>                              # forget a device
ufi device port power-cycle <deviceId> <portIdx>    # PoE power cycle
```

### Clients

```bash
ufi client list                # all connected clients
ufi client get <id>            # client details
ufi client block <id>          # block a client
ufi client unblock <id>        # unblock
```

### Networks

```bash
ufi network list
ufi network get <id>
ufi network create --name "IoT" --vlan-id 30 --purpose corporate
ufi network update <id> --name "IoT Devices"
ufi network delete <id>
ufi network refs <id>          # check dependencies before deleting
```

### Firewall

```bash
# Policies
ufi firewall policy list
ufi firewall policy create --name "Block IoT" --action DENY --source-zone <zoneId>
ufi firewall policy update <id> --enabled=false
ufi firewall policy delete <id>
ufi firewall policy reorder <id1>,<id2>,<id3>

# Zones
ufi firewall zone list
ufi firewall zone create --name "IoT Zone" --network-ids <netId1>,<netId2>
ufi firewall zone update <id> --name "Smart Home"
ufi firewall zone delete <id>
```

### WiFi

```bash
ufi wifi list
ufi wifi get <id>
ufi wifi create --name "Guest" --security wpa2 --password "hunter2"
ufi wifi update <id> --enabled=false
ufi wifi delete <id>
```

### ACL rules

```bash
ufi acl list
ufi acl create --name "Allow printer" --action ALLOW --source-mac AA:BB:CC:DD:EE:FF
ufi acl update <id> --enabled=false
ufi acl delete <id>
ufi acl reorder <id1>,<id2>
```

### DNS policies

```bash
ufi dns list
ufi dns create --name "Cloudflare" --dns-servers 1.1.1.1,1.0.0.1
ufi dns update <id> --dns-servers 9.9.9.9
ufi dns delete <id>
```

## Output

Table by default, `--json` for piping:

```bash
ufi device list                    # human-readable table
ufi device list --json             # JSON (pipe to jq, scripts, etc.)
ufi device list --json | jq '.[].name'
```

## Pagination & filtering

```bash
ufi client list --limit 20 --offset 40
ufi device list --filter "state eq CONNECTED"
```

## Global flags

| Flag | Description |
|------|-------------|
| `--host` | Controller URL (e.g. `https://192.168.1.1`) |
| `--site` | Site ID (default from config) |
| `--api-key` | API key (overrides keychain) |
| `--json` | JSON output |
| `--insecure` | Skip TLS verification (self-signed certs) |
| `--config` | Config file path (default `~/.ufi.yaml`) |

## Authentication

First, generate an API key in your UniFi console: **Settings > System > Advanced > API Keys**.

Then run `ufi login` — it prompts for your controller URL and API key, and stores them:

```bash
ufi login          # store credentials
ufi logout         # remove credentials
```

### Where tokens are stored

| OS | Primary storage | Location |
|----|----------------|----------|
| macOS | Keychain | `Keychain Access > login > ufi-cli` |
| Linux (desktop) | Secret Service (GNOME Keyring / KDE Wallet) | Via D-Bus `org.freedesktop.secrets` |
| Windows | Credential Manager | `Control Panel > Credential Manager > ufi-cli` |
| Linux (headless) | Encrypted file | `~/.ufi/credentials.enc` |

The encrypted file fallback uses **AES-256-GCM** with a key derived from your passphrase via **Argon2id**. If `ufi login` can't reach the system keyring, it prompts you for a passphrase and creates this file instead.

API keys are **never** written to plaintext config files.

### Resolution order

When ufi needs the API key, it checks these in order (first one wins):

1. `UFI_API_KEY` environment variable — for CI, scripts, containers
2. `--api-key` flag — inline override
3. System keyring — set by `ufi login`
4. Encrypted file (`~/.ufi/credentials.enc`) — fallback when no keyring

### Config file

`~/.ufi.yaml` stores host and site only (no secrets):

```yaml
host: https://192.168.1.1
site: default
```
