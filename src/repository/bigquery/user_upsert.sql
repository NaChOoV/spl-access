MERGE access.user AS target
USING (
    SELECT 
        CAST(external_id AS STRING) as external_id,
        CAST(run AS STRING) as run,
        CAST(full_name AS STRING) as full_name,
    FROM UNNEST([
        STRUCT<external_id STRING, run STRING, full_name STRING>
        %s
    ])
) AS source
ON target.external_id = source.external_id
WHEN MATCHED THEN
    UPDATE SET full_name = source.full_name, run = source.run
WHEN NOT MATCHED THEN
    INSERT (external_id, run, full_name)
    VALUES (source.external_id, source.run, source.full_name)