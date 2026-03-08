package com.github.mankit007.recipe.Dashboard.Service.service;

import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.springframework.data.domain.PageRequest;
import org.springframework.data.mongodb.core.ReactiveMongoTemplate;
import org.springframework.data.mongodb.core.aggregation.Aggregation;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.stereotype.Service;

import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.ClickRecord;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.ClickSummary;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.PaginatedClicks;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.RefererStat;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TimeSeriesPoint;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TopLink;
import com.github.mankit007.recipe.Dashboard.Service.model.ClickAggregate;
import com.github.mankit007.recipe.Dashboard.Service.repository.ClickAggregateRepository;
import com.github.mankit007.recipe.Dashboard.Service.repository.ClickEventRepository;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
public class AnalyticsService {

    private final ClickEventRepository clickEventRepository;
    private final ClickAggregateRepository clickAggregateRepository;
    private final ReactiveMongoTemplate mongoTemplate;

    public Mono<ClickSummary> getSummary(String tenantId) {
        Aggregation aggregation = Aggregation.newAggregation(
                Aggregation.match(Criteria.where("tenant_id").is(tenantId)),
                Aggregation.group()
                        .count().as("totalClicks")
                        .addToSet("ip_address").as("uniqueIPs")
                        .addToSet("short_code").as("uniqueShortCodes"));

        return mongoTemplate.aggregate(aggregation, "click_events", Map.class)
                .next()
                .map(result -> ClickSummary.builder()
                        .totalClicks(((Number) result.get("totalClicks")).longValue())
                        .uniqueIPs(((List<?>) result.get("uniqueIPs")).size())
                        .uniqueShortCodes(((List<?>) result.get("uniqueShortCodes")).size())
                        .build())
                .defaultIfEmpty(ClickSummary.builder()
                        .totalClicks(0)
                        .uniqueIPs(0)
                        .uniqueShortCodes(0)
                        .build());
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
        String collectionName = resolveCollectionName(granularity);

        Aggregation aggregation = Aggregation.newAggregation(
                Aggregation.match(Criteria.where("tenant_id").is(tenantId)),
                Aggregation.group("short_code")
                        .sum("click_count").as("totalClicks")
                        .first("original_url").as("originalUrl"),
                Aggregation.sort(org.springframework.data.domain.Sort.by(
                        org.springframework.data.domain.Sort.Direction.DESC, "totalClicks")),
                Aggregation.limit(limit));

        return mongoTemplate.aggregate(aggregation, collectionName, Map.class)
                .map(doc -> TopLink.builder()
                        .shortCode((String) doc.get("_id"))
                        .originalUrl((String) doc.get("originalUrl"))
                        .totalClicks(((Number) doc.get("totalClicks")).longValue())
                        .build())
                .collectList();
    }

    public Mono<List<TimeSeriesPoint>> getTimeSeries(String tenantId, String shortCode, String granularity) {
        if (shortCode != null && !shortCode.isBlank()) {
            return clickAggregateRepository
                    .findByTenantIdAndShortCodeAndGranularityOrderByWindowStartAsc(
                            tenantId, shortCode, granularity.toUpperCase())
                    .map(this::toTimeSeriesPoint)
                    .collectList();
        }

        return clickAggregateRepository
                .findByTenantIdAndGranularityOrderByWindowStartDesc(
                        tenantId, granularity.toUpperCase())
                .map(this::toTimeSeriesPoint)
                .collectSortedList(Comparator.comparing(TimeSeriesPoint::getWindowStart));
    }

    public Mono<List<RefererStat>> getRefererBreakdown(String tenantId, String shortCode) {
        Criteria criteria = Criteria.where("tenant_id").is(tenantId);
        if (shortCode != null && !shortCode.isBlank()) {
            criteria = criteria.and("short_code").is(shortCode);
        }

        Aggregation aggregation = Aggregation.newAggregation(
                Aggregation.match(criteria),
                Aggregation.unwind("referers"),
                Aggregation.group("referers.referer")
                        .sum("referers.count").as("count"),
                Aggregation.sort(org.springframework.data.domain.Sort.by(
                        org.springframework.data.domain.Sort.Direction.DESC, "count")));

        String collectionName = "click_aggregates_hourly";

        return mongoTemplate.aggregate(aggregation, collectionName, Map.class)
                .map(doc -> RefererStat.builder()
                        .referer(normalizeReferer((String) doc.get("_id")))
                        .count(((Number) doc.get("count")).longValue())
                        .build())
                .collectList();
    }

    private TimeSeriesPoint toTimeSeriesPoint(ClickAggregate aggregate) {
        return TimeSeriesPoint.builder()
                .windowStart(aggregate.getWindowStart())
                .windowEnd(aggregate.getWindowEnd())
                .shortCode(aggregate.getShortCode())
                .clickCount(aggregate.getClickCount())
                .uniqueIpCount(aggregate.getUniqueIpCount())
                .build();
    }

    private String resolveCollectionName(String granularity) {
        return switch (granularity.toUpperCase()) {
            case "DAY" -> "click_aggregates_daily";
            case "WEEK" -> "click_aggregates_weekly";
            default -> "click_aggregates_hourly";
        };
    }

    private String normalizeReferer(String referer) {
        if (referer == null || referer.isBlank()) {
            return "Direct";
        }
        return referer;
    }
}
