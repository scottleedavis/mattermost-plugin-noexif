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
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"golang.org/x/image/tiff"
)

type FileTrack struct {
	fileId    string
	timestamp time.Time
}

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	filesTracked []FileTrack
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {

	switch strings.ToUpper(info.Extension) {
	case "JPG", "JPEG", "PNG", "TIFF":
		if data, err := ioutil.ReadAll(file); err != nil {
			p.API.LogError(err.Error())
			return nil, ""
		} else {
			img, _, _ := image.Decode(bytes.NewReader(data))
			switch strings.ToUpper(info.Extension) {
			case "JPG", "JPEG":
				if err := jpeg.Encode(output, img, nil); err != nil {
					p.API.LogError(err.Error())
				} else {
					p.TrackFileId(info)
				}
			case "PNG":
				if err := png.Encode(output, img); err != nil {
					p.API.LogError(err.Error())
				} else {
					p.TrackFileId(info)
				}
			case "TIFF":
				if err := tiff.Encode(output, img, nil); err != nil {
					p.API.LogError(err.Error())
				}
			}
		}
	}

	return nil, ""
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	for _, id := range post.FileIds {
		if contains(p.filesTracked, id) {
			go p.NotifyStatus(post.UserId, post.ChannelId)
			p.UnTrackFileId(id)
		}
	}
}

func (p *Plugin) UnTrackFileId(id string) {
	newTracked := []FileTrack{}
	for _, ft := range p.filesTracked {
		if ft.fileId != id && time.Now().Sub(ft.timestamp).Minutes() < 5 {
			newTracked = append(newTracked, ft)
		}
	}
	p.filesTracked = newTracked
}

func (p *Plugin) TrackFileId(info *model.FileInfo) {
	fileTrack := FileTrack{
		fileId:    info.Id,
		timestamp: time.Now(),
	}
	p.filesTracked = append(p.filesTracked, fileTrack)
}

func (p *Plugin) NotifyStatus(userId string, channelId string) {
	time.Sleep(time.Millisecond * 500)
	p.API.SendEphemeralPost(userId, &model.Post{
		ChannelId: channelId,
		Message:   "Your uploaded photos have had their location and EXIF data removed before being uploaded to the server.",
	})
}

func contains(s []FileTrack, id string) bool {
	for _, a := range s {
		if a.fileId == id {
			return true
		}
	}
	return false
}
