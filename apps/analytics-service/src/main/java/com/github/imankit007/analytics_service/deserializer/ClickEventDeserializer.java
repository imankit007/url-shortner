package com.github.imankit007.analytics_service.deserializer;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;

import com.github.imankit007.analytics_service.model.ClickEvent;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.FieldNamingPolicy;

public class ClickEventDeserializer implements DeserializationSchema<ClickEvent> {

    private static final long serialVersionUID = 1L;
    private transient Gson gson;

    private Gson getGson() {
        if (gson == null) {
            gson = new GsonBuilder()
                    .setFieldNamingPolicy(FieldNamingPolicy.LOWER_CASE_WITH_UNDERSCORES)
                    .create();
        }
        return gson;
    }

    @Override
    public ClickEvent deserialize(byte[] message) throws IOException {
        String json = new String(message, StandardCharsets.UTF_8);
        return getGson().fromJson(json, ClickEvent.class);
    }

    @Override
    public boolean isEndOfStream(ClickEvent nextElement) {
        return false;
    }

    @Override
    public TypeInformation<ClickEvent> getProducedType() {
        return TypeInformation.of(ClickEvent.class);
    }
}
