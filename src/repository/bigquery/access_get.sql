WITH RankedAccess AS (
  SELECT 
    external_id,
    run,
    entry_at,
    location,
    ROW_NUMBER() OVER (PARTITION BY external_id, location ORDER BY entry_at DESC) AS rn
  FROM `access.today_access`
  WHERE exit_at IS NULL
    AND TIMESTAMP_ADD(entry_at, INTERVAL 150 MINUTE) >= CURRENT_TIMESTAMP()
  GROUP BY external_id, run, entry_at, exit_at, location
)
SELECT
  r.external_id,
  u.run,
  u.full_name,
  r.entry_at,
  r.location
FROM RankedAccess r
INNER JOIN `access.user` u ON u.external_id = r.external_id
WHERE r.rn = 1
ORDER BY r.entry_at DESC