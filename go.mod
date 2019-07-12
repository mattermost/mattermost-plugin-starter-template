module github.com/mattermost/mattermost-plugin-starter-template

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible // indirect
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/mattermost/go-i18n v1.11.0 // indirect
	github.com/mattermost/mattermost-server v5.12.0+incompatible
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)

// Workaround for https://github.com/golang/go/issues/30831 and fallout.
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
