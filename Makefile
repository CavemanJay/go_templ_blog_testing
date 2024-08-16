# Based on: https://github.com/a-h/templ/discussions/596#discussioncomment-9011118

DB := blog.db
MIGRATIONS:=migrations

# live: 
# 	@make -j3 live/templ live/server live/tailwind 

live: 
	@make -j4 live/templ live/server live/tailwind live/sync_assets

live/templ:
	templ generate --watch --proxy="http://localhost:7331" --open-browser=true 

live/tailwind:
	pnpm css:watch

live/sync_assets:
	air --build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

live/server:
	air

sql: 
	-@make sql/down 
	-@make sql/up 
	-@make sql/queries 

sql/up:
	goose -dir $(MIGRATIONS) sqlite3 $(DB) up

sql/down:
	goose -dir $(MIGRATIONS) sqlite3 $(DB) down

sql/queries:
	sqlc generate