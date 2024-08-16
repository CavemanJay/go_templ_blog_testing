-- name: QueryArticles :many
SELECT * FROM articles;

-- name: QueryArticleBySlug :one
SELECT * FROM articles WHERE slug = ?;