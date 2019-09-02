module github.com/mattermost/mattermost-plugin-sample

go 1.12

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/dsoprea/go-exif v0.0.0-20190901173045-3ce78807c90f
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20190422055009-d6f9ba25cf48
	github.com/dsoprea/go-logging v0.0.0-20190624164917-c4f10aab7696
	github.com/dsoprea/go-png-image-structure v0.0.0-20190624104353-c9b28dcdc5c8
	github.com/dyatlov/go-opengraph v0.0.0-20180429202543-816b6608b3c8 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/hashicorp/go-plugin v1.0.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattermost/mattermost-server v5.9.0+incompatible
	github.com/mattermost/viper v1.0.4 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.3.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/scottleedavis/go-exif-remove v0.0.0-20190902210608-b25a47538306
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1 // indirect
	golang.org/x/image v0.0.0-20190703141733-d6a02ce849c9
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	google.golang.org/grpc v1.20.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

// Workaround for https://github.com/golang/go/issues/30831 and fallout.
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
