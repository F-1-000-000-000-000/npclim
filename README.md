# npclim

**nginx proxy cli manager**

A command-line tool for managing nginx (and nginx-derived) reverse proxy configuration files through the commandline. Tired of manually keeping track of working directories, editing template files, forgetting practically default configs, etc? With npclim creating a new proxy host is as simple as `npclim new`.

## Installation

```bash
git clone https://github.com/F-1-000-000-000-000/npclim
cd npclim
go build -o npclim
```

Optionally move the binary to your PATH:

```bash
sudo mv npclim /usr/local/bin/
```

or add the npclim directory to your $PATH variable.

## Configuration

npclim reads from `~/.config/npclim/config`. All flags can be set here as defaults:

```yaml
proxy-location: /etc/nginx/sites-enabled/
domain: example.com
filename: "{{.Subdomain}}.conf"
```

You can also place a custom template at `~/.config/npclim/template.conf` and it will be used automatically (see [Templates](#templates)).

## Usage

### List proxy hosts

```bash
npclim ls                  # list filenames
npclim ls -l               # list with domain and proxy info
npclim ls /path/to/dir     # list from a specific directory
```

### Create a new proxy host

```bash
npclim new -s hass -d example.com -p 8123
npclim new hass -d example.com -p 8123    # uses hass.conf as filename
npclim new -s hass -p 8123                # domain falls back to config
```

**Flags:**

| Flag | Short | Default | Description |
|---|---|---|---|
| `--subdomain` | `-s` | | Subdomain for the proxy host |
| `--domain` | `-d` | | Base domain |
| `--host` | `-H` | `localhost` | Host to forward traffic to |
| `--port` | `-p` | | Port to forward traffic to |
| `--configuration-template` | `-t` | | Custom template file |
| `--proxy-location` | `-l` | `./` | Output directory |
| `--filename-template` | `-f` | `{{.Subdomain}}.{{.Domain}}.conf` | Filename template |

### Edit a proxy host

```bash
npclim edit hass
```

Opens the configuration file in your `$EDITOR` (falls back to `vim`).

### Remove a proxy host

```bash
npclim rm hass
```

## Templates

npclim uses Go's `text/template` syntax. The following variables are available:

| Variable | Description |
|---|---|
| `{{.Subdomain}}` | The subdomain |
| `{{.Domain}}` | The base domain |
| `{{.Host}}` | The forwarding host |
| `{{.Port}}` | The forwarding port |

The default template generates a basic nginx reverse proxy config:

```nginx
server {
    listen 80;
    server_name {{.Subdomain}}.{{.Domain}};

    location / {
        proxy_pass http://{{.Host}}:{{.Port}};
    }
}
```

Template priority order: `-t` flag → configuration-template set in `~/.config/npclim/config`) → `~/.config/npclim/template.conf` → built-in default.
