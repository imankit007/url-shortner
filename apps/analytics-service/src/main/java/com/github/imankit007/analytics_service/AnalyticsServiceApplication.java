package com.github.imankit007.analytics_service;

import java.time.Duration;
import java.time.Instant;

import org.apache.flink.api.common.eventtime.SerializableTimestampAssigner;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows;

import com.github.imankit007.analytics_service.aggregation.ClickAggregator;
import com.github.imankit007.analytics_service.deserializer.ClickEventDeserializer;
import com.github.imankit007.analytics_service.model.ClickEvent;
import com.github.imankit007.analytics_service.sink.MongoAggregateSink;
import com.github.imankit007.analytics_service.sink.MongoRawClickSink;

public class AnalyticsServiceApplication {

    private static final String KAFKA_BOOTSTRAP_SERVERS = "localhost:9092";
    private static final String KAFKA_TOPIC = "click-events";
    private static final String KAFKA_GROUP_ID = "analytics-flink-pipeline";
    private static final String MONGO_URI = "mongodb://localhost:27017";
    private static final String MONGO_DATABASE = "url-shortner-db";

    public static void main(String[] args) throws Exception {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        // Kafka source
        KafkaSource<ClickEvent> kafkaSource = KafkaSource.<ClickEvent>builder()
                .setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS)
                .setTopics(KAFKA_TOPIC)
                .setGroupId(KAFKA_GROUP_ID)
                .setStartingOffsets(OffsetsInitializer.earliest())
                .setValueOnlyDeserializer(new ClickEventDeserializer())
                .build();

        // Watermark strategy using the event timestamp
        WatermarkStrategy<ClickEvent> watermarkStrategy = WatermarkStrategy
                .<ClickEvent>forBoundedOutOfOrderness(Duration.ofSeconds(30))
                .withTimestampAssigner(
                        (SerializableTimestampAssigner<ClickEvent>) (event, recordTimestamp) -> {
                            try {
                                return Instant.parse(event.getTimestamp()).toEpochMilli();
                            } catch (Exception e) {
                                return recordTimestamp;
                            }
                        });

        DataStream<ClickEvent> clickStream = env
                .fromSource(kafkaSource, watermarkStrategy, "Kafka Click Events");

        // Sink 1: Raw click events to MongoDB
        clickStream.addSink(new MongoRawClickSink(MONGO_URI, MONGO_DATABASE))
                .name("MongoDB Raw Click Sink");

        // Sink 2: Hourly tumbling window aggregation
        clickStream
                .keyBy(ClickEvent::getShortCode)
                .window(TumblingEventTimeWindows.of(Duration.ofHours(1)))
                .process(new ClickAggregator("HOUR"))
                .addSink(new MongoAggregateSink(MONGO_URI, MONGO_DATABASE))
                .name("MongoDB Hourly Aggregate Sink");

        // Sink 3: Daily tumbling window aggregation
        clickStream
                .keyBy(ClickEvent::getShortCode)
                .window(TumblingEventTimeWindows.of(Duration.ofDays(1)))
                .process(new ClickAggregator("DAY"))
                .addSink(new MongoAggregateSink(MONGO_URI, MONGO_DATABASE))
                .name("MongoDB Daily Aggregate Sink");

        // Sink 4: Weekly tumbling window aggregation
        clickStream
                .keyBy(ClickEvent::getShortCode)
                .window(TumblingEventTimeWindows.of(Duration.ofDays(7)))
                .process(new ClickAggregator("WEEK"))
                .addSink(new MongoAggregateSink(MONGO_URI, MONGO_DATABASE))
                .name("MongoDB Weekly Aggregate Sink");

        env.execute("URL Shortener Analytics Pipeline");
    }
}
