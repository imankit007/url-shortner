package com.github.imankit007.analytics_service.sink;

import java.time.Instant;

import org.apache.flink.streaming.api.functions.sink.SinkFunction;
import org.bson.Document;

import com.github.imankit007.analytics_service.model.ClickEvent;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;

public class MongoRawClickSink implements SinkFunction<ClickEvent> {

    private static final long serialVersionUID = 1L;

    private final String mongoUri;
    private final String databaseName;
    private transient MongoClient mongoClient;
    private transient MongoCollection<Document> collection;

    public MongoRawClickSink(String mongoUri, String databaseName) {
        this.mongoUri = mongoUri;
        this.databaseName = databaseName;
    }

    private void ensureConnection() {
        if (mongoClient == null) {
            mongoClient = MongoClients.create(mongoUri);
            MongoDatabase database = mongoClient.getDatabase(databaseName);
            collection = database.getCollection("click_events");
        }
    }

    @Override
    public void invoke(ClickEvent event, Context context) {
        ensureConnection();

        Document doc = new Document()
                .append("short_code", event.getShortCode())
                .append("original_url", event.getOriginalUrl())
                .append("tenant_id", event.getTenantId())
                .append("timestamp", Instant.parse(event.getTimestamp()))
                .append("user_agent", event.getUserAgent())
                .append("referer", event.getReferer())
                .append("ip_address", event.getIpAddress());

        collection.insertOne(doc);
    }
}
