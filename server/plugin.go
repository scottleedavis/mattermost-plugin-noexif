package main

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	exifremove "github.com/scottleedavis/go-exif-remove"
)

// Plugin comment
type Plugin struct {
	plugin.MattermostPlugin
}

//FileWillBeUploaded hook
func (p *Plugin) FileWillBeUploaded(_ *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	switch strings.ToUpper(info.Extension) {
	case "JPG", "JPEG", "PNG":
		data, err := ioutil.ReadAll(file)
		if err != nil {
			errMsg := "Failed to read image: " + err.Error()
			p.API.LogWarn(errMsg)
			return nil, errMsg
		}

		_, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			errMsg := "Original image is corrupt: " + err.Error()
			p.API.LogWarn(errMsg)
			return nil, errMsg
		}

		filtered, err := exifremove.Remove(data)
		if err != nil {
			errMsg := "Failed to remove exif information: " + err.Error()
			p.API.LogWarn(errMsg)
			return nil, errMsg
		}

		if _, _, err = image.Decode(bytes.NewReader(filtered)); err != nil {
			errMsg := "Failed to decode filtered image: " + err.Error()
			p.API.LogWarn(errMsg)
			return nil, errMsg
		}

		if _, err := output.Write(filtered); err != nil {
			errMsg := "Failed to write new image: " + err.Error()
			p.API.LogWarn(errMsg)
			return nil, errMsg
		}
	}

	return nil, ""
}
