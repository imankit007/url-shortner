package com.github.mankit007.recipe.Dashboard.Service.model;

import java.time.Instant;
import java.util.List;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Document
public class ClickAggregate {

    @Id
    private String id;

    @Field("short_code")
    private String shortCode;

    @Field("tenant_id")
    private String tenantId;

    @Field("original_url")
    private String originalUrl;

    @Field("granularity")
    private String granularity;

    @Field("window_start")
    private Instant windowStart;

    @Field("window_end")
    private Instant windowEnd;

    @Field("click_count")
    private long clickCount;

    @Field("unique_ip_count")
    private long uniqueIpCount;

    @Field("referers")
    private List<RefererEntry> referers;

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class RefererEntry {
        private String referer;
        private long count;
    }
}
