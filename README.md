# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes with configurable rules.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config portwatch.yaml
```

Example `portwatch.yaml`:

```yaml
interval: 30s
alert:
  - type: log
    path: /var/log/portwatch.log
rules:
  allow:
    - 22
    - 80
    - 443
  deny: all
```

portwatch will poll open ports at the specified interval and emit an alert whenever a port outside the allowed list is detected or a previously open port disappears.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `portwatch.yaml` | Path to config file |
| `--interval` | `30s` | Poll interval |
| `--once` | `false` | Run a single scan and exit |

```bash
# Run a one-time scan and print results
portwatch --once
```

---

## License

MIT © yourname