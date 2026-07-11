# FALTASI POS — Control Centre

A CLI for the sales team to activate and renew Balce POS customer licenses,
and for the system admin to manage packages — talking directly to the
wapangaji Django backend.

## Running it

Download the binary for your OS from the latest release and run it —
no install step needed.

- **Windows**: `faltasi-windows-amd64.exe`
- **macOS (Apple Silicon)**: `faltasi-darwin-arm64`
- **macOS (Intel)**: `faltasi-darwin-amd64`
- **Linux**: `faltasi-linux-amd64`

On first run you'll pick a language (Kiswahili/English), then log in with
your existing Balce partner phone number and password. Your session is
cached locally for up to 12 hours so you don't need to log in every time.

## What it does

- **Activate/Renew License** — look up a customer by their POS Hardware ID
  (or phone number if new), pick a package, record how they paid, confirm,
  done. The license activates instantly.
- **Search Customer** — quick read-only lookup, no payment involved.
- **Payment History** — recent manually-recorded payments.
- **Manage Packages** (system admin only) — create or edit pricing
  packages. Editing a package never changes what existing customers
  already paid for — only new activations use the new price/days.

## Development

```
go build ./cmd/faltasi
```

Requires Go 1.22+. Uses [Charm](https://charm.sh)'s `bubbletea`/`lipgloss`/`huh`
for the terminal UI.

## Releasing

Push a tag matching `v*` (e.g. `v1.0.0`) — GitHub Actions cross-compiles
for Windows, macOS (Intel + Apple Silicon), and Linux, and publishes them
as release assets automatically.
