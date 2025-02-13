FROM gcr.io/distroless/static:latest

WORKDIR /app

ARG TARGETARCH

# Copy the prebuilt Go binary
COPY bin/wayfarer-${TARGETARCH} /app/wayfarer

CMD ["/app/wayfarer"]
