# Build stage
FROM golang:1.16.7-alpine AS build
WORKDIR /src
COPY . /src
RUN go build

FROM hashicorp/packer:1.7.4 AS packer

# Final stage
FROM alpine

COPY --from=build /src/carpenter /usr/local/bin/carpenter
RUN carpenter help

COPY --from=packer /bin/packer /usr/local/bin/packer
RUN packer version