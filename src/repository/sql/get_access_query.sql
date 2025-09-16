WITH RankedAccess AS (SELECT run,
                             entry_at,
                             exit_at,
                             location,
                             ROW_NUMBER() OVER (PARTITION BY run, location ORDER BY entry_at DESC) AS rn
                      FROM access
                      WHERE exit_at IS NULL
                        AND entry_at + INTERVAL '2.5 hours' >= NOW()
                      GROUP BY run, entry_at, exit_at, location)
SELECT external_id,
       "user".run,
       "user".full_name,
       entry_at,
       exit_at,
       location
FROM RankedAccess
         INNER JOIN "user" on "user".run = RankedAccess.run
WHERE rn = 1
ORDER BY entry_at DESC;