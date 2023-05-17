module github.com/mattermost/mattermost-plugin-starter-template

go 1.16

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattermost/mattermost-server/server/public v0.0.4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	golang.org/x/crypto v0.9.0 // indirect
)

replace github.com/mattermost/mattermost-server/server/public => ../mattermost-server/server/public
