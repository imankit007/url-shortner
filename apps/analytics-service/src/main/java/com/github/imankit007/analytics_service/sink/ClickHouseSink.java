package com.github.imankit007.analytics_service.sink;

import java.io.IOException;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.Timestamp;
import java.time.Instant;

import org.apache.flink.api.connector.sink2.Sink;
import org.apache.flink.api.connector.sink2.SinkWriter;

import com.github.imankit007.analytics_service.model.ClickEvent;

public class ClickHouseSink implements Sink<ClickEvent> {


    private final String jdbcUrl;

    public ClickHouseSink(String jdbcUrl) {
        this.jdbcUrl = jdbcUrl;
    }

    @SuppressWarnings("deprecation")
    @Override
    public SinkWriter<ClickEvent> createWriter(InitContext context) {
        return new ClickHouseWriter(jdbcUrl);
    }

    private static class ClickHouseWriter implements SinkWriter<ClickEvent> {

        private final String jdbcUrl;
        private transient Connection connection;

        ClickHouseWriter(String jdbcUrl) {
            this.jdbcUrl = jdbcUrl;
        }

        private void ensureConnection() throws Exception {
            if (connection == null || connection.isClosed()) {
                connection = DriverManager.getConnection(jdbcUrl);
            }
        }

        @Override
        public void write(ClickEvent event, Context context) throws IOException, InterruptedException {
            try {
                ensureConnection();

                String sql = "INSERT INTO url_shortener.click_events " +
                        "(short_code, original_url, tenant_id, timestamp, user_agent, referer, ip_address) " +
                        "VALUES (?, ?, ?, ?, ?, ?, ?)";

                try (PreparedStatement stmt = connection.prepareStatement(sql)) {
                    stmt.setString(1, event.getShortCode());
                    stmt.setString(2, event.getOriginalUrl());
                    stmt.setString(3, event.getTenantId());
                    stmt.setTimestamp(4, Timestamp.from(Instant.parse(event.getTimestamp())));
                    stmt.setString(5, event.getUserAgent());
                    stmt.setString(6, event.getReferer());
                    stmt.setString(7, event.getIpAddress());
                    stmt.executeUpdate();
                }
            } catch (Exception e) {
                throw new IOException("Failed to write to ClickHouse", e);
            }
        }

        @Override
        public void flush(boolean endOfInput) {
            // ClickHouse JDBC handles flushing per statement
        }

        @Override
        public void close() throws Exception {
            if (connection != null && !connection.isClosed()) {
                connection.close();
            }
        }
    }
}
