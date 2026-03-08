package com.github.mankit007.recipe.Dashboard.Service.dto;

import java.time.Instant;
import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

public final class AnalyticsDTO {

    private AnalyticsDTO() {
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ClickSummary {
        private long totalClicks;
        private long uniqueIPs;
        private long uniqueShortCodes;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class TimeSeriesPoint {
        private Instant windowStart;
        private Instant windowEnd;
        private String shortCode;
        private long clickCount;
        private long uniqueIpCount;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class TopLink {
        private String shortCode;
        private String originalUrl;
        private long totalClicks;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class RefererStat {
        private String referer;
        private long count;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ClickRecord {
        private String shortCode;
        private String originalUrl;
        private Instant timestamp;
        private String userAgent;
        private String referer;
        private String ipAddress;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class PaginatedClicks {
        private List<ClickRecord> clicks;
        private long totalCount;
        private int page;
        private int size;
    }
}
