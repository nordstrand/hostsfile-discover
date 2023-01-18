# docker build -t hostdisc --platform=linux/arm .

FROM --platform=$BUILDPLATFORM golang:1.19.5-alpine AS build
WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=7 go build -o /out/app ./...

FROM alpine
ENTRYPOINT [ "/bin/app" ]
COPY --from=build /out/app /bin
