FROM golang:1.19.2 AS build

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/fleet-observer

FROM gcr.io/distroless/base

COPY --from=build /build/fleet-observer /bin/fleet-observer
WORKDIR /var/lib/fleet-observer
ENTRYPOINT [ "/bin/fleet-observer" ]
CMD [ "observe" ]
