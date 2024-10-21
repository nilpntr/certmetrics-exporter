ARG VERSION
FROM sammobach/go:latest as build
LABEL maintainer="Sam Mobach <hello@sammobach.com>"
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN mkdir -p /app/bin && go build -buildvcs=false -ldflags "-s -w -X 'github.com/nilpntr/certmetrics-exporter/cmd.version=$VERSION'" -o bin/certmetrics-exporter github.com/nilpntr/certmetrics-exporter
CMD ["/app/bin/certmetrics-exporter"]

FROM alpine:3.20
LABEL maintainer="Sam Mobach <hello@sammobach.com>"
RUN mkdir /app
COPY --from=build /app/bin /app
WORKDIR /app
CMD ["/app/certmetrics-exporter"]