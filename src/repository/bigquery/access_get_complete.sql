WITH RankedAccess AS (
  SELECT
    external_id,
    run,
    entry_at,
    exit_at,
    location,
    ROW_NUMBER() OVER (PARTITION BY external_id, location ORDER BY entry_at DESC) AS rn
  FROM `access.today_access`
  GROUP BY external_id, run, entry_at, exit_at, location
)
SELECT
  r.external_id,
  u.run,
  u.full_name,
  r.entry_at,
  r.exit_at,
  r.location
FROM RankedAccess r
INNER JOIN `access.user` u ON u.external_id = r.external_id
WHERE r.rn = 1
ORDER BY r.entry_at DESC