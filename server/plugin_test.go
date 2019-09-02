package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
	pngstructure "github.com/dsoprea/go-png-image-structure"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestFileWillBeUpload(t *testing.T) {

	t.Run("JPG EXIF removal", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		data, err := ioutil.ReadFile("../assets/exif.jpg")
		assert.Nil(t, err)

		fi := &model.FileInfo{
			Extension: "JPG",
		}

		r := bytes.NewReader(data)

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		_, reason := p.FileWillBeUploaded(nil, fi, r, w)
		assert.Equal(t, reason, "")

		jmp := jpegstructure.NewJpegMediaParser()
		sl, err := jmp.ParseBytes(buf.Bytes())
		_, _, err = sl.Exif()
		assert.NotNil(t, err)

	})

	t.Run("PNG EXIF removal", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		data, err := ioutil.ReadFile("../assets/exif.png")
		assert.Nil(t, err)

		fi := &model.FileInfo{
			Extension: "JPG",
		}

		r := bytes.NewReader(data)

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		_, reason := p.FileWillBeUploaded(nil, fi, r, w)
		assert.Equal(t, reason, "")

		pmp := pngstructure.NewPngMediaParser()
		cs, err := pmp.ParseBytes(buf.Bytes())
		_, _, err = cs.Exif()
		assert.NotNil(t, err)

	})
}
