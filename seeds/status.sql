SELECT
    pid,
    state,
    query_start,
    query
FROM
    pg_start_activity
WHERE
    query LIKE '%pg_sleep%' ANS state != 'idle';