# Build stage
FROM golang:alpine AS build
WORKDIR /src
ADD . /src
RUN go build

# Final stage
FROM alpine
COPY --from=build /src/carpenter /usr/local/bin/carpenter
RUN carpenter help