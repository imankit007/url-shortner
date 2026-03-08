package com.github.imankit007.analytics_service.aggregation;

import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

import org.apache.flink.streaming.api.functions.windowing.ProcessWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.util.Collector;

import com.github.imankit007.analytics_service.model.ClickAggregate;
import com.github.imankit007.analytics_service.model.ClickEvent;

public class ClickAggregator
        extends ProcessWindowFunction<ClickEvent, ClickAggregate, String, TimeWindow> {

    private static final long serialVersionUID = 1L;

    private final String granularity;

    public ClickAggregator(String granularity) {
        this.granularity = granularity;
    }

    @Override
    public void process(String key,
                        ProcessWindowFunction<ClickEvent, ClickAggregate, String, TimeWindow>.Context context,
                        Iterable<ClickEvent> events,
                        Collector<ClickAggregate> out) {

        long clickCount = 0;
        Set<String> uniqueIPs = new HashSet<>();
        Map<String, Long> refererCounts = new HashMap<>();
        String tenantId = null;
        String originalUrl = null;
        String shortCode = null;

        for (ClickEvent event : events) {
            clickCount++;

            if (shortCode == null) {
                shortCode = event.getShortCode();
                tenantId = event.getTenantId();
                originalUrl = event.getOriginalUrl();
            }

            if (event.getIpAddress() != null && !event.getIpAddress().isEmpty()) {
                uniqueIPs.add(event.getIpAddress());
            }

            String referer = event.getReferer();
            if (referer == null || referer.isEmpty()) {
                referer = "Direct";
            }
            refererCounts.merge(referer, 1L, Long::sum);
        }

        TimeWindow window = context.window();

        ClickAggregate aggregate = new ClickAggregate();
        aggregate.setShortCode(shortCode);
        aggregate.setTenantId(tenantId);
        aggregate.setOriginalUrl(originalUrl);
        aggregate.setGranularity(granularity);
        aggregate.setWindowStart(window.getStart());
        aggregate.setWindowEnd(window.getEnd());
        aggregate.setClickCount(clickCount);
        aggregate.setUniqueIPs(uniqueIPs);
        aggregate.setRefererCounts(refererCounts);

        out.collect(aggregate);
    }
}
