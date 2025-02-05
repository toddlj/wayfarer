package com.toddljones.wayfarer;

import jakarta.inject.Singleton;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.stream.Collectors;

@Singleton
public class TelegramNotifier {

    private static final String TELEGRAM_API_BASE = "https://api.telegram.org/bot";
    private static final Logger log = LoggerFactory.getLogger(TelegramNotifier.class);

    private final String botToken;

    TelegramNotifier() {
        this.botToken = System.getenv("TELEGRAM_BOT_TOKEN");
        if (botToken == null) {
            throw new IllegalStateException("TELEGRAM_BOT_TOKEN environment variable must be set");
        }
    }

    public void sendMessage(String userId, String message) {
        String endpoint = TELEGRAM_API_BASE + botToken + "/sendMessage";
        String params = String.format("chat_id=%s&text=%s", userId, message);

        try {
            HttpResponse response = executePost(endpoint, params);
            if (!response.isSuccess()) {
                throw new RuntimeException("Failed to send message. Response Code: " + response.getCode());
            }
            log.info("Message sent successfully to telegram userId={}", userId);
        } catch (IOException e) {
            throw new RuntimeException("Failed to send message to Telegram userId=" + userId, e);
        }
    }

    private HttpResponse executePost(String endpoint, String params) throws IOException {
        HttpURLConnection conn = createConnection(endpoint, "POST");
        conn.setRequestProperty("Content-Type", "application/x-www-form-urlencoded");

        try (OutputStream os = conn.getOutputStream()) {
            os.write(params.getBytes(StandardCharsets.UTF_8));
        }

        return getResponse(conn);
    }

    private HttpURLConnection createConnection(String endpoint, String method) throws IOException {
        URL url = new URL(endpoint);
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod(method);
        conn.setDoOutput(method.equals("POST"));
        return conn;
    }

    private HttpResponse getResponse(HttpURLConnection conn) throws IOException {
        int responseCode = conn.getResponseCode();
        String responseBody = "";

        if (responseCode == HttpURLConnection.HTTP_OK) {
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(conn.getInputStream()))) {
                responseBody = reader.lines().collect(Collectors.joining("\n"));
            }
        }

        return new HttpResponse(responseCode, responseBody);
    }

    // Helper class for HTTP responses
    private static class HttpResponse {
        private final int code;
        private final String body;

        HttpResponse(int code, String body) {
            this.code = code;
            this.body = body;
        }

        boolean isSuccess() {
            return code == HttpURLConnection.HTTP_OK;
        }

        int getCode() {
            return code;
        }

        String getBody() {
            return body;
        }
    }
}
