package com.github.imankit007.analytics_service.sink;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.apache.flink.streaming.api.functions.sink.SinkFunction;
import org.bson.Document;

import com.github.imankit007.analytics_service.model.ClickAggregate;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import com.mongodb.client.model.Filters;
import com.mongodb.client.model.ReplaceOptions;

public class MongoAggregateSink implements SinkFunction<ClickAggregate> {

    private static final long serialVersionUID = 1L;

    private final String mongoUri;
    private final String databaseName;
    private transient MongoClient mongoClient;
    private transient MongoDatabase database;

    public MongoAggregateSink(String mongoUri, String databaseName) {
        this.mongoUri = mongoUri;
        this.databaseName = databaseName;
    }

    private void ensureConnection() {
        if (mongoClient == null) {
            mongoClient = MongoClients.create(mongoUri);
            database = mongoClient.getDatabase(databaseName);
        }
    }

    @Override
    public void invoke(ClickAggregate aggregate, Context context) {
        ensureConnection();

        String collectionName = switch (aggregate.getGranularity()) {
            case "HOUR" -> "click_aggregates_hourly";
            case "DAY" -> "click_aggregates_daily";
            case "WEEK" -> "click_aggregates_weekly";
            default -> "click_aggregates_hourly";
        };

        MongoCollection<Document> collection = database.getCollection(collectionName);

        List<Document> refererDocs = new ArrayList<>();
        for (Map.Entry<String, Long> entry : aggregate.getRefererCounts().entrySet()) {
            refererDocs.add(new Document("referer", entry.getKey()).append("count", entry.getValue()));
        }

        Document doc = new Document()
                .append("short_code", aggregate.getShortCode())
                .append("tenant_id", aggregate.getTenantId())
                .append("original_url", aggregate.getOriginalUrl())
                .append("granularity", aggregate.getGranularity())
                .append("window_start", Instant.ofEpochMilli(aggregate.getWindowStart()))
                .append("window_end", Instant.ofEpochMilli(aggregate.getWindowEnd()))
                .append("click_count", aggregate.getClickCount())
                .append("unique_ip_count", aggregate.getUniqueIPs().size())
                .append("referers", refererDocs);

        // Upsert by (short_code, granularity, window_start)
        collection.replaceOne(
                Filters.and(
                        Filters.eq("short_code", aggregate.getShortCode()),
                        Filters.eq("granularity", aggregate.getGranularity()),
                        Filters.eq("window_start", Instant.ofEpochMilli(aggregate.getWindowStart()))
                ),
                doc,
                new ReplaceOptions().upsert(true)
        );
    }
}
