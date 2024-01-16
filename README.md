# go-pocketbase-htmx-templ-tailwind

go pocketbase htmx templ tailwind template

- [go docs](https://go.dev/doc/)
- [pocketbase docs](https://pocketbase.io/docs/)
- [templ docs](https://github.com/a-h/templ)
- [tailwind docs](https://tailwindcss.com/docs/installation)

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

1. `docker build .`
2. `docker run -p 8080:8080`

save these urls:

- https://git.sunshine.industries/efim/go-ssr-pocketbase-oauth-attempt
- https://www.reddit.com/r/pocketbase/comments/18lhedp/jwt_auth_with_htmx/
