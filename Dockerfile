FROM golang:1.22.4 AS build

WORKDIR /src
COPY ./src .
ARG CI_JOB_TOKEN
RUN  echo "building..." && GOPRIVATE=* GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o /service


FROM alpine

RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

WORKDIR /
COPY --from=build /service /service
COPY ./src/config.json /

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

USER nonroot

CMD ["web"]
ENTRYPOINT [ "./service" ]
