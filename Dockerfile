FROM golang:1.17-alpine3.14 as build-stage
WORKDIR /app/
ENV GO111MODULE=on

# for caching modules
COPY go.mod .
COPY go.sum .
RUN go mod download

ARG VERSION
ARG RELEASE
ARG BUILD_DATE
ARG GIT_REVISION

COPY ./ /app/

RUN export GO_VERSION=$(go version | awk '{print $3" "$4}') && \
	go build -o app \
	-ldflags="-X main.VERSION=$VERSION-$RELEASE -X 'main.BUILD_DATE=$BUILD_DATE' \
    	-X 'main.GIT_REVISION=$GIT_REVISION' -X 'main.GO_VERSION=$GO_VERSION'"

# production stage
FROM alpine:3.14 as production-stage
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

# copy backend binary
COPY --from=build-stage /app/app /app/app

EXPOSE 9800
ENTRYPOINT ["/app/app"]
CMD ["-h"]
