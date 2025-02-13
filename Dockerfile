FROM gcr.io/distroless/static:latest

# Copy the prebuilt Go binary
COPY wayfarer /

CMD ["/wayfarer"]
