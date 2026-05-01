POSTGRESQL_URL=postgres://postgres:aruzhann@localhost:5432/blog_post?sslmode=disable

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	migrate -database $(POSTGRESQL_URL) -path migrations up

.PHONY: migrate-down
migrate-down:
	migrate -database $(POSTGRESQL_URL) -path migrations down

.PHONY: migrate-force
migrate-force:
	migrate -database $(POSTGRESQL_URL) -path migrations force $(version)

.PHONY: migrate-version
migrate-version:
	migrate -database $(POSTGRESQL_URL) -path migrations version