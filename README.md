# Mattermost Plugin NoEXIF
[![Build Status](https://img.shields.io/circleci/project/github/scottleedavis/mattermost-plugin-noexif/master.svg)](https://circleci.com/gh/scottleedavis/mattermost-plugin-noexif) [![codecov](https://codecov.io/gh/scottleedavis/mattermost-plugin-noexif/branch/master/graph/badge.svg)](https://codecov.io/gh/scottleedavis/mattermost-plugin-noexif)  

This plugin removes EXIF information from images uploaded to [Mattermost](http://mattermost.com).
Currently supports jpg/jpeg, and png files.


## Build
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.example.my-plugin.tar.gz
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

Alternatively, if you are running your `mattermost-server` out of a sibling directory by the same name, use the `deploy` target alone to  unpack the files into the right directory. You will need to restart your server and manually enable your plugin.

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).
