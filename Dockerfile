# Build stage
FROM golang:alpine AS build
WORKDIR /src
ADD . /src
RUN go build

FROM hashicorp/packer AS packer

# Final stage
FROM alpine
COPY --from=build /src/carpenter /usr/local/bin/carpenter
RUN carpenter help

COPY --from=packer /bin/packer /usr/local/bin/packer
RUN packer version