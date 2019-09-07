# Mattermost Plugin NoEXIF
[![Build Status](https://img.shields.io/circleci/project/github/scottleedavis/mattermost-plugin-noexif/master.svg)](https://circleci.com/gh/scottleedavis/mattermost-plugin-noexif) [![codecov](https://codecov.io/gh/scottleedavis/mattermost-plugin-noexif/branch/master/graph/badge.svg)](https://codecov.io/gh/scottleedavis/mattermost-plugin-noexif)  [![Releases](https://img.shields.io/github/release/scottleedavis/mattermost-plugin-noexif.svg)](https://github.com/scottleedavis/mattermost-plugin-noexif/releases/latest)
 

This plugin removes EXIF information from images uploaded to [Mattermost](http://mattermost.com).
Currently supports jpg/jpeg, and png files.

_Requires Mattermost 5.9 or higher_

## Build
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.github.scottleedavis.mattermost-plugin-noexif.tar.gz
```
