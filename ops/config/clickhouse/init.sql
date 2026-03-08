CREATE DATABASE IF NOT EXISTS url_shortener;

-- Raw click events table
CREATE TABLE IF NOT EXISTS url_shortener.click_events
(
    short_code   String,
    original_url String,
    tenant_id    String,
    timestamp    DateTime64(3, 'UTC'),
    user_agent   String,
    referer      String,
    ip_address   String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (tenant_id, short_code, timestamp)
TTL toDateTime(timestamp) + INTERVAL 1 YEAR;

-- Hourly materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS url_shortener.click_aggregates_hourly
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(window_start)
ORDER BY (tenant_id, short_code, window_start)
AS
SELECT
    tenant_id,
    short_code,
    any(original_url) AS original_url,
    toStartOfHour(timestamp) AS window_start,
    count() AS click_count,
    uniqExact(ip_address) AS unique_ip_count
FROM url_shortener.click_events
GROUP BY tenant_id, short_code, toStartOfHour(timestamp);

-- Daily materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS url_shortener.click_aggregates_daily
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(window_start)
ORDER BY (tenant_id, short_code, window_start)
AS
SELECT
    tenant_id,
    short_code,
    any(original_url) AS original_url,
    toStartOfDay(timestamp) AS window_start,
    count() AS click_count,
    uniqExact(ip_address) AS unique_ip_count
FROM url_shortener.click_events
GROUP BY tenant_id, short_code, toStartOfDay(timestamp);

-- Weekly materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS url_shortener.click_aggregates_weekly
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(window_start)
ORDER BY (tenant_id, short_code, window_start)
AS
SELECT
    tenant_id,
    short_code,
    any(original_url) AS original_url,
    toStartOfWeek(timestamp) AS window_start,
    count() AS click_count,
    uniqExact(ip_address) AS unique_ip_count
FROM url_shortener.click_events
GROUP BY tenant_id, short_code, toStartOfWeek(timestamp);
