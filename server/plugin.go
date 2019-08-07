package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
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
			return nil, ""
		} else {
			img, _, _ := image.Decode(bytes.NewReader(data))
			switch strings.ToUpper(info.Extension) {
			case "JPG", "JPEG":
				if err := jpeg.Encode(output, img, nil); err != nil {
					p.API.LogError(err.Error())
				}
			case "PNG":
				if err := png.Encode(output, img); err != nil {
					p.API.LogError(err.Error())
				}
			}
		}
	}

	return nil, ""
}
