FROM golang:1.21
ARG TAILWIND_VERSION=3.4.1
ARG TEMPL_VERSION=0.2.513
WORKDIR /app
COPY go.mod go.sum ./
RUN GOPROXY=direct go mod download -x
COPY . .
ADD https://github.com/tailwindlabs/tailwindcss/releases/download/v$TAILWIND_VERSION/tailwindcss-linux-arm64 ./tailwindcss
ADD https://github.com/a-h/templ/releases/download/v$TEMPL_VERSION/templ_Linux_arm64.tar.gz ./tmp/templ.tar.gz
RUN tar -xvzf ./tmp/templ.tar.gz templ && rm ./tmp/templ.tar.gz && ./templ generate && \
  chmod +x ./tailwindcss && ./tailwindcss -i tailwind.css -o ./public/css/tailwind.css && \
  GOPROXY=direct go build -o ./main .
# uncomment to copy the local pb_migrations dir into the image
# COPY ./pb_migrations /pb/pb_migrations
# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pb/pb_hooks
EXPOSE 8080
CMD ["./main", "serve", "--http=0.0.0.0:8080"]
