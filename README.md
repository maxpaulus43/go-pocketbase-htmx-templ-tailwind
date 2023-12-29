# go-pocketbase-htmx-templ-tailwind
go pocketbase htmx templ tailwind template

## Prereqs

1. install golang 1.21
2. install tailwind cli: `https://tailwindcss.com/blog/standalone-cli`
3. install templ: `go install github.com/a-h/templ/cmd/templ@latest`
4. install air: `go install github.com/cosmtrek/air@latest`
5. run `go get`

## To Develop

1. `air` // watches for file changes and generates template files and builds css changes

## To Build for Production

1. `templ generate && ./tailwindcss -i tailwind.css -o ./public/css/tailwind.css && go build -o ./main .`
2. copy `public` folder and `main` binary to a server 
3. expose port 8080 via systemd or reverse proxy (see [pocketbase setup guide](https://pocketbase.io/docs/going-to-production/#minimal-setup))
4. run `main serve --http=0.0.0.0:8080`

## To Build for Production (Docker)

```dockerfile
FROM alpine:latest

ARG PB_VERSION=0.20.1

RUN apk add --no-cache \
    unzip \
    ca-certificates

# download and unzip PocketBase
ADD https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip /tmp/pb.zip
RUN unzip /tmp/pb.zip -d /pb/

# uncomment to copy the local pb_migrations dir into the image
# COPY ./pb_migrations /pb/pb_migrations

# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pb/pb_hooks

COPY ./public /public

EXPOSE 8080

# start PocketBase
CMD ["/pb/pocketbase", "serve", "--http=0.0.0.0:8080"]
```
