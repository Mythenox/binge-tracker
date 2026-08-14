-- name: GetLeastViewedSeason :one
WITH ViewtimeTotals AS (
    SELECT show_title, season_number, SUM(viewtime) AS total_viewtime
    FROM episodes
    GROUP BY (show_title, season_number)
)
SELECT s.*
FROM seasons s
JOIN ViewtimeTotals ct ON (s.show_title, s.season_number) = (ct.show_title, season_number)
WHERE 
    ct.total_viewtime = (SELECT MIN(total_viewtime) FROM ViewtimeTotals) AND
    s.show_title = sqlc.arg(show_title)
ORDER BY season_number ASC
LIMIT 1;