# Splatoon (Wii U) NEX server

Pretendo's own Splatoon NEX server implementation (combined auth + secure),
patched to authenticate players against a real, already-verified account
instead of requiring a local account service.

Upstream: [PretendoNetwork/splatoon](https://github.com/PretendoNetwork/splatoon).
`game_server_id` `10162b00`, access key `6f599f81` (from
[kinnay.github.io](https://kinnay.github.io/)'s public Wii U NEX game
database), title IDs `0005000010162b00` / `...2c00` / `...176900` / `...176a00`.

## What's different from upstream

Upstream's `globals.PasswordFromPID` calls out to an account gRPC service's
`GetNEXPassword` - meaning every PID has to already be registered as a local
account in whatever that gRPC service backs. This patch
(`globals/password_from_pid.go`) instead derives the password directly from
the PID with a shared HMAC secret (`PN_NEX_PASSWORD_SECRET`), so any PID
already verified elsewhere (by whatever issues this title's NEX token) can
log in - no local account registration, no account gRPC service required.

Whatever issues the NEX token needs to compute the password field the exact
same way: `base64url(HMAC-SHA256(secret, little-endian PID bytes))`.

## Configuration

| Variable | Purpose |
|---|---|
| `PN_SPLATOON_AUTHENTICATION_SERVER_PORT` | UDP port for the authentication server |
| `PN_SPLATOON_SECURE_SERVER_HOST` | Hostname/IP advertised to clients for the secure server |
| `PN_SPLATOON_SECURE_SERVER_PORT` | UDP port for the secure server |
| `PN_SPLATOON_ACCOUNT_GRPC_HOST` / `_PORT` / `_API_KEY` | Account gRPC service (friends list, etc. - not used for password lookup anymore, see above) |
| `PN_SPLATOON_FRIENDS_GRPC_HOST` / `_PORT` / `_API_KEY` | Friends gRPC service |
| `PN_SPLATOON_LOCAL_AUTH` | Set to `1` to allow `settings.json`-based local test accounts as a fallback |
| `PN_SPLATOON_POSTGRES_URI` | Postgres connection string (matchmaking/BOSS state) |
| `PN_NEX_PASSWORD_SECRET` | Hex-encoded, at least 32 bytes. Must match whatever issues the NEX token for this `game_server_id` |

## Running

See upstream's own build/run instructions - this is a drop-in patch, not a
restructure.
