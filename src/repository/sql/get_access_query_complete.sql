WITH RankedAccess AS (SELECT external_id,
                             entry_at,
                             exit_at,
                             location,
                             ROW_NUMBER() OVER (PARTITION BY external_id, location ORDER BY entry_at DESC) AS rn
                      FROM access
                      GROUP BY external_id, entry_at, exit_at, location)
SELECT RankedAccess.external_id,
       "user".run,
       "user".full_name,
       entry_at,
       exit_at,
       location
FROM RankedAccess
         INNER JOIN "user" on "user".external_id = RankedAccess.external_id
WHERE rn = 1
ORDER BY entry_at DESC;