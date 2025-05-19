# nimbus-notify
Service that allows to subscribe to regular emails with weather updates


### Task
https://github.com/mykhailo-hrynko/se-school-5


### Useful commands
migrate -source file://./migrations -database "postgres://user:password@localhost:5432/nimbus-notify?sslmode=disable" up

mockery init nimbus-notify/internal/service
mockery

go test ./...

