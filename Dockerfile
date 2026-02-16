# syntax=docker/dockerfile:1
FROM gcr.io/distroless/static-debian12:nonroot
COPY --chown=nonroot:nonroot bin/bike-rental-linux-amd64 /usr/bin/bike-rental-api
COPY --chown=nonroot:nonroot docs /app/docs
EXPOSE 8080
ENTRYPOINT ["/usr/bin/bike-rental-api"]
