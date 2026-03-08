package com.github.mankit007.recipe.Dashboard.Service.model;

import java.time.Instant;

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
@Document(collection = "click_events")
public class ClickEvent {

    @Id
    private String id;

    @Field("short_code")
    private String shortCode;

    @Field("original_url")
    private String originalUrl;

    @Field("tenant_id")
    private String tenantId;

    @Field("timestamp")
    private Instant timestamp;

    @Field("user_agent")
    private String userAgent;

    @Field("referer")
    private String referer;

    @Field("ip_address")
    private String ipAddress;
}
