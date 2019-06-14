module github.com/mattermost/mattermost-plugin-sample

go 1.12

require (
	github.com/mattermost/mattermost-server v5.12.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)

// Workaround for https://github.com/golang/go/issues/30831 and fallout.
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
