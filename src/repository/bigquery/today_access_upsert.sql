MERGE access.today_access AS target
USING (
    SELECT 
        CAST(external_id AS STRING) as external_id,
        CAST(run AS STRING) as run,
        CAST(location AS INT64) as location,
        CAST(entry_at AS TIMESTAMP) as entry_at,
        CAST(exit_at AS TIMESTAMP) as exit_at
    FROM UNNEST([
        STRUCT<external_id STRING, run STRING, location INT64, entry_at TIMESTAMP, exit_at TIMESTAMP>
        %s
    ])
) AS source
ON target.external_id = source.external_id 
    AND target.location = source.location 
    AND target.entry_at = source.entry_at
WHEN MATCHED THEN
    UPDATE SET exit_at = source.exit_at
WHEN NOT MATCHED THEN
    INSERT (external_id, run, location, entry_at, exit_at)
    VALUES (source.external_id, source.run, source.location, source.entry_at, source.exit_at)