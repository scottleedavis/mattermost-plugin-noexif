package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
	pngstructure "github.com/dsoprea/go-png-image-structure"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFileWillBeUpload(t *testing.T) {
	t.Run("original image is corrupt", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogInfo", mock.Anything).Maybe()
			return api
		}

		api := setupAPI()
		api.On("LogWarn", mock.AnythingOfType("string"))
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		fi := &model.FileInfo{
			Extension: "JPG",
		}

		r := bytes.NewReader([]byte{})

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		fi, reason := p.FileWillBeUploaded(nil, fi, r, w)
		assert.Equal(t, reason, "Original image is corrupt: image: unknown format")
		assert.Nil(t, fi)
	})

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

		fi, reason := p.FileWillBeUploaded(nil, fi, r, w)
		assert.Equal(t, reason, "")
		assert.Nil(t, fi)

		jmp := jpegstructure.NewJpegMediaParser()
		sl, err := jmp.ParseBytes(buf.Bytes())
		assert.NoError(t, err)
		require.NotNil(t, sl)

		_, _, err = sl.Exif()
		assert.Error(t, err)
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

		fi, reason := p.FileWillBeUploaded(nil, fi, r, w)
		assert.Equal(t, reason, "")
		assert.Nil(t, fi)

		pmp := pngstructure.NewPngMediaParser()
		cs, err := pmp.ParseBytes(buf.Bytes())
		assert.NoError(t, err)
		require.NotNil(t, cs)

		_, _, err = cs.Exif()
		assert.Error(t, err)
	})
}
