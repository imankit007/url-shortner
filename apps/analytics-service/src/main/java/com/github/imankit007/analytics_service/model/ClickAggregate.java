package com.github.imankit007.analytics_service.model;

import java.io.Serializable;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

public class ClickAggregate implements Serializable {

    private static final long serialVersionUID = 1L;

    private String shortCode;
    private String tenantId;
    private String originalUrl;
    private String granularity; // HOUR, DAY, WEEK
    private long windowStart;
    private long windowEnd;
    private long clickCount;
    private Set<String> uniqueIPs;
    private Map<String, Long> refererCounts;

    public ClickAggregate() {
        this.uniqueIPs = new HashSet<>();
        this.refererCounts = new HashMap<>();
    }

    public String getShortCode() { return shortCode; }
    public void setShortCode(String shortCode) { this.shortCode = shortCode; }

    public String getTenantId() { return tenantId; }
    public void setTenantId(String tenantId) { this.tenantId = tenantId; }

    public String getOriginalUrl() { return originalUrl; }
    public void setOriginalUrl(String originalUrl) { this.originalUrl = originalUrl; }

    public String getGranularity() { return granularity; }
    public void setGranularity(String granularity) { this.granularity = granularity; }

    public long getWindowStart() { return windowStart; }
    public void setWindowStart(long windowStart) { this.windowStart = windowStart; }

    public long getWindowEnd() { return windowEnd; }
    public void setWindowEnd(long windowEnd) { this.windowEnd = windowEnd; }

    public long getClickCount() { return clickCount; }
    public void setClickCount(long clickCount) { this.clickCount = clickCount; }

    public Set<String> getUniqueIPs() { return uniqueIPs; }
    public void setUniqueIPs(Set<String> uniqueIPs) { this.uniqueIPs = uniqueIPs; }

    public Map<String, Long> getRefererCounts() { return refererCounts; }
    public void setRefererCounts(Map<String, Long> refererCounts) { this.refererCounts = refererCounts; }
}
