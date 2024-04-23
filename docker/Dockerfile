FROM node:18.17.1 AS nextjs-build
WORKDIR /app
COPY view/package.json view/package-lock.json view/panda.config.ts ./
RUN npm install
COPY view/ .
RUN npm run build

FROM golang:1.21.4 AS go-build
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/ 
COPY config/ config/ 
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -v -o lgtmeme ./cmd/lgtmeme

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-build /go/src/app/lgtmeme .
COPY --from=nextjs-build /app/out ./view/out
EXPOSE 8080
CMD ["./lgtmeme"]