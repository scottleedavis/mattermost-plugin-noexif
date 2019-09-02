package main

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/scottleedavis/go-exif-remove"
)

type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {

	switch strings.ToUpper(info.Extension) {
	case "JPG", "JPEG", "PNG":
		if data, err := ioutil.ReadAll(file); err != nil {
			p.API.LogError(err.Error())
			return nil, err.Error()
		} else {
			if _, _, err = image.Decode(bytes.NewReader(data)); err != nil {
				errMsg := "ERROR: original image is corrupt " + err.Error()
				p.API.LogInfo(errMsg)
				return nil, errMsg
			}
			if filtered, err := exifremove.Remove(data); err != nil {
				p.API.LogError(err.Error())
				return nil, err.Error()
			} else {
				if _, err := output.Write(filtered); err != nil {
					p.API.LogError(err.Error())
				}
			}
		}
	}
	return nil, ""
}
