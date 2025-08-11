SELECT a.run, u.full_name, a.location, a.entry_at, a.exit_at 
FROM access.today_access a
INNER JOIN access.user u ON u.external_id = a.external_id