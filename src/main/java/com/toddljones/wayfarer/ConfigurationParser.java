package com.toddljones.wayfarer;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;

public class ConfigurationParser {

    Configuration parse(File file) {
        if (!file.exists()) {
            throw new IllegalArgumentException("File does not exist: " + file);
        }
        String configFileContents;
        try {
            configFileContents = Files.readString(file.toPath());
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        try {
            ObjectMapper objectMapper = new ObjectMapper(new YAMLFactory());
            return objectMapper.readValue(configFileContents, Configuration.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }
}
