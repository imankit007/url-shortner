package com.github.mankit007.recipe.Dashboard.Service.service;

import java.sql.Timestamp;
import java.util.List;

import org.springframework.data.domain.PageRequest;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;

import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.ClickRecord;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.ClickSummary;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.PaginatedClicks;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.RefererStat;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TimeSeriesPoint;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TopLink;
import com.github.mankit007.recipe.Dashboard.Service.repository.ClickEventRepository;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Schedulers;

@Service
@RequiredArgsConstructor
public class AnalyticsService {

    private final JdbcTemplate jdbcTemplate;
    private final ClickEventRepository clickEventRepository;

    public Mono<ClickSummary> getSummary(String tenantId) {
        return Mono.fromCallable(() -> {
            String sql = """
                    SELECT
                        count() AS total_clicks,
                        uniqExact(ip_address) AS unique_ips,
                        uniqExact(short_code) AS unique_short_codes
                    FROM url_shortener.click_events
                    WHERE tenant_id = ?
                    """;

            return jdbcTemplate.queryForObject(sql, (rs, rowNum) ->
                    ClickSummary.builder()
                            .totalClicks(rs.getLong("total_clicks"))
                            .uniqueIPs(rs.getLong("unique_ips"))
                            .uniqueShortCodes(rs.getLong("unique_short_codes"))
                            .build(),
                    tenantId);
        }).subscribeOn(Schedulers.boundedElastic());
    }

    public Mono<PaginatedClicks> getRecentClicks(String tenantId, int page, int size) {
        return clickEventRepository.countByTenantId(tenantId)
                .flatMap(totalCount ->
                        clickEventRepository
                                .findByTenantIdOrderByTimestampDesc(tenantId, PageRequest.of(page, size))
                                .map(event -> ClickRecord.builder()
                                        .shortCode(event.getShortCode())
                                        .originalUrl(event.getOriginalUrl())
                                        .timestamp(event.getTimestamp())
                                        .userAgent(event.getUserAgent())
                                        .referer(event.getReferer())
                                        .ipAddress(event.getIpAddress())
                                        .build())
                                .collectList()
                                .map(records -> PaginatedClicks.builder()
                                        .clicks(records)
                                        .totalCount(totalCount)
                                        .page(page)
                                        .size(size)
                                        .build()));
    }

    public Mono<List<TopLink>> getTopLinks(String tenantId, int limit, String granularity) {
        return Mono.fromCallable(() -> {
            String table = resolveAggregateTable(granularity);

            String sql = String.format("""
                    SELECT
                        short_code,
                        any(original_url) AS original_url,
                        sum(click_count) AS total_clicks
                    FROM %s
                    WHERE tenant_id = ?
                    GROUP BY short_code
                    ORDER BY total_clicks DESC
                    LIMIT ?
                    """, table);

            return jdbcTemplate.query(sql, (rs, rowNum) ->
                    TopLink.builder()
                            .shortCode(rs.getString("short_code"))
                            .originalUrl(rs.getString("original_url"))
                            .totalClicks(rs.getLong("total_clicks"))
                            .build(),
                    tenantId, limit);
        }).subscribeOn(Schedulers.boundedElastic());
    }

    public Mono<List<TimeSeriesPoint>> getTimeSeries(String tenantId, String shortCode, String granularity) {
        return Mono.fromCallable(() -> {
            String table = resolveAggregateTable(granularity);

            String sql;
            Object[] params;

            if (shortCode != null && !shortCode.isBlank()) {
                sql = String.format("""
                        SELECT
                            window_start,
                            short_code,
                            sum(click_count) AS click_count,
                            sum(unique_ip_count) AS unique_ip_count
                        FROM %s
                        WHERE tenant_id = ? AND short_code = ?
                        GROUP BY window_start, short_code
                        ORDER BY window_start ASC
                        """, table);
                params = new Object[]{tenantId, shortCode};
            } else {
                sql = String.format("""
                        SELECT
                            window_start,
                            '' AS short_code,
                            sum(click_count) AS click_count,
                            sum(unique_ip_count) AS unique_ip_count
                        FROM %s
                        WHERE tenant_id = ?
                        GROUP BY window_start
                        ORDER BY window_start ASC
                        """, table);
                params = new Object[]{tenantId};
            }

            return jdbcTemplate.query(sql, (rs, rowNum) -> {
                Timestamp windowStart = rs.getTimestamp("window_start");
                return TimeSeriesPoint.builder()
                        .windowStart(windowStart != null ? windowStart.toInstant() : null)
                        .windowEnd(null)
                        .shortCode(rs.getString("short_code"))
                        .clickCount(rs.getLong("click_count"))
                        .uniqueIpCount(rs.getLong("unique_ip_count"))
                        .build();
            }, params);
        }).subscribeOn(Schedulers.boundedElastic());
    }

    public Mono<List<RefererStat>> getRefererBreakdown(String tenantId, String shortCode) {
        return Mono.fromCallable(() -> {
            String sql;
            Object[] params;

            if (shortCode != null && !shortCode.isBlank()) {
                sql = """
                        SELECT
                            if(referer = '', 'Direct', referer) AS referer,
                            count() AS cnt
                        FROM url_shortener.click_events
                        WHERE tenant_id = ? AND short_code = ?
                        GROUP BY referer
                        ORDER BY cnt DESC
                        """;
                params = new Object[]{tenantId, shortCode};
            } else {
                sql = """
                        SELECT
                            if(referer = '', 'Direct', referer) AS referer,
                            count() AS cnt
                        FROM url_shortener.click_events
                        WHERE tenant_id = ?
                        GROUP BY referer
                        ORDER BY cnt DESC
                        """;
                params = new Object[]{tenantId};
            }

            return jdbcTemplate.query(sql, (rs, rowNum) ->
                    RefererStat.builder()
                            .referer(rs.getString("referer"))
                            .count(rs.getLong("cnt"))
                            .build(),
                    params);
        }).subscribeOn(Schedulers.boundedElastic());
    }

    private String resolveAggregateTable(String granularity) {
        return switch (granularity.toUpperCase()) {
            case "DAY" -> "url_shortener.click_aggregates_daily";
            case "WEEK" -> "url_shortener.click_aggregates_weekly";
            default -> "url_shortener.click_aggregates_hourly";
        };
    }
}
