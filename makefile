build-www:
	go build -o bin/ryan_dashboard.app ./cmd/dashboard

production-build-www:
	GOOS=linux GOARCH=amd64  go build -o bin/ryan_dashboard.app.amd64 ./cmd/dashboard

start-dev: tmp 
	make -j 3  templ serve sync_assets

format-html:
	templ fmt ./www 

templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false  -v

serve:
	air \
  --build.cmd "go build -o tmp/bin/ryan_dashboard ./cmd/dashboard " --build.bin "tmp/bin/ryan_dashboard" --build.delay "100" \
  --build.exclude_dir "node_modules" \
  --build.include_ext "go,env,html" \
  --build.stop_on_error "false" \
  --misc.clean_on_exit true

sync_assets: 
	air \
  --build.cmd "templ generate --notify-proxy" \
  --build.bin "true" \
  --build.delay "100" \
  --build.exclude_dir "" \
  --build.include_dir "assets,html" \
  --build.include_ext "js,css,html"

tmp:
	@echo "Creating tmp directory..." 
	mkdir -p ./tmp/bin