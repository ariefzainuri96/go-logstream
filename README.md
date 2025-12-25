## air-verse for hot reload

1. <https://github.com/air-verse/air>
2. run air with `air -c .air.toml`

## swagger openapi

1. Install the swag CLI tool if not installed: Install the executable globally by running `go install github.com/swaggo/swag/cmd/swag@latest`
2. to update the documentation, head to cmd/api
3. move to `base-entity.go` comment `DeletedAt` and uncomment the below implementation
4. run this command `swag init`
5. after success, revert back the change of commenting `DeletedAt`

## Migration

1. Head to migrate github `https://github.com/golang-migrate/migrate/tree/master/cmd/migrate` -> this link contains cli installation
2. Create migration file `migrate create -seq -ext sql -dir ././migrations create_users`
3. Perform manual migration `export $(grep -v '^#' .env | xargs) && migrate -path ./migrations -database="${DB_ADDR}" up`

## .air.toml

1. current working is for linux because we are using docker for running this apps
   bin = "./bin/api"
   cmd = "go build -o ./bin/api ./cmd/api/"

2. if you running locally, change .air.toml line 7-8 to:
   bin = "./bin/api.exe"
   cmd = "go build -o ./bin/ ./cmd/api/"
