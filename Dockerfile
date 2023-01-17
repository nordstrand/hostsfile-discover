FROM --platform=$BUILDPLATFORM golang:1.19.5-alpine AS build
WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app ./...

FROM alpine
COPY --from=build /out/app /bin
