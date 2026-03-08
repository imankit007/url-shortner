package com.github.imankit007.analytics_service;

import java.time.Duration;
import java.time.Instant;

import org.apache.flink.api.common.eventtime.SerializableTimestampAssigner;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

import com.github.imankit007.analytics_service.deserializer.ClickEventDeserializer;
import com.github.imankit007.analytics_service.model.ClickEvent;
import com.github.imankit007.analytics_service.sink.ClickHouseSink;

public class AnalyticsServiceApplication {

    private static final String KAFKA_BOOTSTRAP_SERVERS = "localhost:9092";
    private static final String KAFKA_TOPIC = "click-events";
    private static final String KAFKA_GROUP_ID = "analytics-flink-pipeline";
    private static final String CLICKHOUSE_JDBC_URL = "jdbc:clickhouse://localhost:8123/url_shortener";

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

        // Sink: Raw click events to ClickHouse
        // Materialized views in ClickHouse automatically aggregate hourly/daily/weekly
        clickStream
                .sinkTo(new ClickHouseSink(CLICKHOUSE_JDBC_URL))
                .name("ClickHouse Sink");

        env.execute("URL Shortener Analytics Pipeline");
    }
}
