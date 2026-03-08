package com.github.mankit007.recipe.Dashboard.Service.controller;

import java.util.List;

import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.ClickSummary;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.PaginatedClicks;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.RefererStat;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TimeSeriesPoint;
import com.github.mankit007.recipe.Dashboard.Service.dto.AnalyticsDTO.TopLink;
import com.github.mankit007.recipe.Dashboard.Service.service.AnalyticsService;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v1/analytics/{tenantId}")
@CrossOrigin
@RequiredArgsConstructor
public class AnalyticsController {

    private final AnalyticsService analyticsService;

    @GetMapping("/summary")
    public Mono<ClickSummary> getSummary(@PathVariable String tenantId) {
        return analyticsService.getSummary(tenantId);
    }

    @GetMapping("/clicks")
    public Mono<PaginatedClicks> getRecentClicks(
            @PathVariable String tenantId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        return analyticsService.getRecentClicks(tenantId, page, size);
    }

    @GetMapping("/top-links")
    public Mono<List<TopLink>> getTopLinks(
            @PathVariable String tenantId,
            @RequestParam(defaultValue = "10") int limit,
            @RequestParam(defaultValue = "HOUR") String granularity) {
        return analyticsService.getTopLinks(tenantId, limit, granularity);
    }

    @GetMapping("/timeseries")
    public Mono<List<TimeSeriesPoint>> getTimeSeries(
            @PathVariable String tenantId,
            @RequestParam(required = false) String shortCode,
            @RequestParam(defaultValue = "HOUR") String granularity) {
        return analyticsService.getTimeSeries(tenantId, shortCode, granularity);
    }

    @GetMapping("/referers")
    public Mono<List<RefererStat>> getRefererBreakdown(
            @PathVariable String tenantId,
            @RequestParam(required = false) String shortCode) {
        return analyticsService.getRefererBreakdown(tenantId, shortCode);
    }
}
