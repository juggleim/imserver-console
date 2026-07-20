# Push multi-package migration notes

## Verified source schema

The authoritative initialization DDL in the sibling `im-server` and
`im-server-cluster` repositories defines:

- `androidpushconfs.uniq_channel` on `(app_key, push_channel)`.
- `ioscertificates.uniq_package` on `(app_key)`.

Those constraints are stricter than the new business identities, so a database
matching the checked-in DDL cannot already contain duplicate
`(app_key, push_channel, package)` or `(app_key, package)` groups. The console's
local MySQL credentials were not accepted during implementation, so deployment
must still run the following read-only preflight against the target environment:

```sql
SHOW CREATE TABLE androidpushconfs;
SHOW INDEX FROM androidpushconfs;
SELECT app_key, push_channel, package, COUNT(*) AS row_count
FROM androidpushconfs
GROUP BY app_key, push_channel, package
HAVING COUNT(*) > 1;

SHOW CREATE TABLE ioscertificates;
SHOW INDEX FROM ioscertificates;
SELECT app_key, package, COUNT(*) AS row_count
FROM ioscertificates
GROUP BY app_key, package
HAVING COUNT(*) > 1;
```

Do not run `20260720.sql` until both duplicate queries return no rows.

## Non-destructive rollback

Application rollback does not require deleting either new index. Before
restoring the old, stricter indexes, stop writes and inspect whether multiple
packages now exist per old identity:

```sql
SELECT app_key, push_channel, COUNT(*) AS row_count
FROM androidpushconfs
GROUP BY app_key, push_channel
HAVING COUNT(*) > 1;

SELECT app_key, COUNT(*) AS row_count
FROM ioscertificates
GROUP BY app_key
HAVING COUNT(*) > 1;
```

If either query returns rows, export and reconcile them with the application
owner; never delete package configurations automatically. Only after every old
identity has at most one row may an operator create `uniq_channel` on
`androidpushconfs(app_key, push_channel)` and `uniq_package` on
`ioscertificates(app_key)`, then remove the new indexes.
