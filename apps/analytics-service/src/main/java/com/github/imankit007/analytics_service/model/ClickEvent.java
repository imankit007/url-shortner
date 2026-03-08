package com.github.imankit007.analytics_service.model;

import java.io.Serializable;

public class ClickEvent implements Serializable {

    private static final long serialVersionUID = 1L;

    private String shortCode;
    private String originalUrl;
    private String tenantId;
    private String timestamp;
    private String userAgent;
    private String referer;
    private String ipAddress;

    public ClickEvent() {
    }

    public ClickEvent(String shortCode, String originalUrl, String tenantId,
                      String timestamp, String userAgent, String referer, String ipAddress) {
        this.shortCode = shortCode;
        this.originalUrl = originalUrl;
        this.tenantId = tenantId;
        this.timestamp = timestamp;
        this.userAgent = userAgent;
        this.referer = referer;
        this.ipAddress = ipAddress;
    }

    public String getShortCode() { return shortCode; }
    public void setShortCode(String shortCode) { this.shortCode = shortCode; }

    public String getOriginalUrl() { return originalUrl; }
    public void setOriginalUrl(String originalUrl) { this.originalUrl = originalUrl; }

    public String getTenantId() { return tenantId; }
    public void setTenantId(String tenantId) { this.tenantId = tenantId; }

    public String getTimestamp() { return timestamp; }
    public void setTimestamp(String timestamp) { this.timestamp = timestamp; }

    public String getUserAgent() { return userAgent; }
    public void setUserAgent(String userAgent) { this.userAgent = userAgent; }

    public String getReferer() { return referer; }
    public void setReferer(String referer) { this.referer = referer; }

    public String getIpAddress() { return ipAddress; }
    public void setIpAddress(String ipAddress) { this.ipAddress = ipAddress; }

    @Override
    public String toString() {
        return "ClickEvent{" +
                "shortCode='" + shortCode + '\'' +
                ", originalUrl='" + originalUrl + '\'' +
                ", tenantId='" + tenantId + '\'' +
                ", timestamp='" + timestamp + '\'' +
                ", userAgent='" + userAgent + '\'' +
                ", referer='" + referer + '\'' +
                ", ipAddress='" + ipAddress + '\'' +
                '}';
    }
}
