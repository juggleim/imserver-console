# APNs P8 Credential Contract

This standard-library-only package is shared with the IM server through the
`github.com/juggleim/imserver-console/commons/apnscredentials` import path.

```go
func Validate(privateKey []byte, keyID, teamID string) error
func ValidTopic(topic string) bool
const MaxPrivateKeySize = 16 * 1024
```

`Validate` accepts exactly one unencrypted PKCS#8 PEM ECDSA P-256 key, at most
16 KiB, and 10-character uppercase alphanumeric Apple Key ID and Team ID.
It does not verify Apple account/topic/environment/VoIP authorization.

`ValidTopic` accepts a configured bundle of 1-100 ASCII bytes from
`[A-Za-z0-9.-]` only, without trimming. Validate the base bundle before
appending `.voip`; the derived VoIP topic can therefore reach 105 bytes.
Console normalizes surrounding whitespace as before and rejects interior
whitespace, slashes, non-ASCII characters and oversized bundles.

## Raw Storage

- `IosCertificateDao.P8PrivateKey []byte` maps to `p8_private_key BLOB NULL`
  and is excluded from JSON with `json:"-"`.
- Save the exact uploaded PEM bytes after validation, including surrounding
  whitespace and line endings. Validation does not mutate caller-owned bytes.
- P8 uses only raw storage. The shared package and console have no master-key
  dependency, encryption compatibility layer or conversion API.
- Database readers and backups can access the original private key. Protect
  database/backup access accordingly; write-only APIs are not at-rest encryption.
- SQL credential writes use a silent GORM session even if debug logging is
  enabled. Errors never include input values. Do not log bodies, private keys,
  certificate passwords or JWTs.

## API

Existing `/admingateway/apps/iospushcer/{get,list,set,upload}` paths remain.
Add `auth_type`, `p8_key_id`, `p8_team_id` and multipart `p8_file`; retain
`ioscer`, `voip_ioscer`, `cert_pwd`, `voip_cert_pwd`, `is_product`, `package`
and `original_package`. New rows default to P12; omitted auth_type on edits
keeps the stored mode. Empty passwords/metadata and absent uploads preserve
the stored values. `set` edits an existing row; creation uses `upload`.

Queries return safe metadata including `p8_key_name`, `has_p8_key` and
`config_version`, never any password, certificate bytes or raw key bytes.
`has_p8_key` means the raw column is nonempty. A filename alone does not make
a key available. Presence does not certify Apple authorization or replace
validation during save/runtime.
This intentionally changes old iOS password-prefill behavior. Android is
unchanged. No key-download endpoint exists. Config versions start at 1 and
increment on every successful edit. Inactive-mode credentials are retained.
All edits read, merge, validate and write under the same row lock.
Metadata-only P8 edits validate the retained raw key. Missing multipart files
preserve raw bytes. P8 edits without a saved or uploaded key fail validation.
Explicit P12 rollback and legacy P12 metadata
edits remain possible with the required version and retained P12 credentials.

JSON and multipart edits from the UI send the `config_version` read when the
draft opened. Under the lock, all supplied versions must match before merging
any fields. Stored P8 rows, switches to P8 and P8 credential edits require a
version; only legacy P12 edits without P8 inputs allow it to be omitted.
Stale or missing required versions return HTTP 409 with body code 409 and no
credentials. The UI preserves the draft and offers an explicit discard/reload,
never automatically retries stale IDs against a newly rotated private key.

Legacy P12 metadata-only edits do not require ordinary certificate bytes,
filenames or nonempty passwords: VoIP-only and passwordless rows remain
editable. Empty password inputs still preserve the stored value. New P12
creation retains the prior upload requirements. Replacing an existing file
requires nonempty bytes and a filename, not a nonempty password. Console does
not cryptographically validate P12; that remains a runtime responsibility.

## Deployment

1. Back up the database and protect access to databases and backups that will
   contain raw private keys. Coordinate console and IM upgrades before enabling P8.
2. Apply `commons/dbcommons/sqls/20260908.sql`, which independently adds six
   columns with information_schema guards: `auth_type`, `p8_key_id`,
   `p8_team_id`, `p8_private_key BLOB NULL`, `p8_key_name`, `config_version`.
   Existing rows default to P12 and version 1; existing P12 columns are retained.
   Console and IM migrations have independent version markers and must run
   sequentially on pinned connections. Verify the six columns, retained P12
   columns and `(app_key, package)` unique index explicitly;
   mocks and migration markers alone do not prove MySQL DDL acceptance.
   `dbcommons.Upgrade() error` stops at any version-read, SQL-file or
   version-write failure; standalone startup checks this error before serving
   HTTP. Only a missing version row or MySQL 1146 on the initial version-table
   lookup allows version-zero bootstrap. Invalid stored versions fail closed.
3. Ensure console and IM database configurations point to the same
   `ioscertificates` table. Upgrade all sending nodes and the console before
   uploading a P8 key. No deployment master key is required.
4. Deliver the rebuilt embedded `webconsole/web/dist`. Publish/pin the console
   module for IM builds that do not use the existing local module replace.
5. Enable one bundle. Check authorized sandbox/production and ordinary/VoIP
   delivery in staging. Configuration updates propagate on the next send
   after runtime refresh, within approximately five minutes when DB is healthy.

Rollback: explicitly select retained P12 while upgraded code is still running,
wait for refresh or restart sending nodes, verify P12, then downgrade. Never
drop the new columns as part of rollback. Invalid P8 does not fall back to P12.
Older binaries without P8 support must use P12 after rollback.
The P8-to-P12 editor accepts retained or newly uploaded VoIP-only certificates
without requiring an ordinary certificate or a nonempty password. New P12
configuration requirements are unchanged.
Apple key rotation is a new P8 upload; revocation still requires Apple-side
action and coordinated refresh/restart for emergency invalidation.

## Historical Deployment

The previous 2026-09-08 iteration used AES-256-GCM, a versioned envelope and a
deployment Secret. Historical verification records are not rewritten to imply
raw storage was deployed. The main repository's same-name
OpenSpec records deployments to hk-ali and ali, and user-confirmed sandbox
ordinary delivery on hk-ali. Those results apply to the encrypted iteration.
The user confirmed that P8 is unreleased, so the current schema is simplified
directly rather than retaining compatibility with that experimental iteration.
This work performs no deployment or remote Secret operations.

## Acceptance Boundaries

Offline tests generate temporary keys and use SQL mocks. They cover the
exact raw bytes, invalid or missing input, locked merging, version
updates, quiet SQL, safe DTOs, access bindings, missing files, multipart bounds
and migration statement ordering. Use `RUN_DB_TEST=0` for the offline console
suite; do not run the main pushmanager external-provider tests recursively.
They do not prove real MySQL DDL idempotency/concurrency or Apple authorization.
The raw revision still requires disposable MySQL shared-table migration and live
sandbox/production/VoIP acceptance; earlier encrypted deployment evidence does
not certify it. No actual key files or servers are needed by these tests.
