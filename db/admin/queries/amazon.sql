-- name: InsertAmazonProductDetail :execrows
INSERT INTO amazon_cache (asin, title, image_medium_url, link) VALUES (?, ?, ?, ?);

-- name: CountAmazonCacheByAsin :one
SELECT count(1)
FROM amazon_cache
WHERE asin = ?;

-- name: GetAmazonImageUrlByAsin :one
SELECT image_medium_url
FROM amazon_cache
WHERE asin = ?;
