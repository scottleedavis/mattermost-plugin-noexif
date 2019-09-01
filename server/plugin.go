package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-png-image-structure"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

const (
	JpegMediaType  = "jpeg"
	PngMediaType   = "png"
	OtherMediaType = "other"
)

type MediaContext struct {
	MediaType string
	RootIfd   *exif.Ifd
	RawExif   []byte
	Media     interface{}
}

type IfdEntry struct {
	IfdPath     string      `json:"ifd_path"`
	FqIfdPath   string      `json:"fq_ifd_path"`
	IfdIndex    int         `json:"ifd_index"`
	TagId       uint16      `json:"tag_id"`
	TagName     string      `json:"tag_name"`
	TagTypeId   uint16      `json:"tag_type_id"`
	TagTypeName string      `json:"tag_type_name"`
	UnitCount   uint32      `json:"unit_count"`
	Value       interface{} `json:"value"`
	ValueString string      `json:"value_string"`
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {

	switch strings.ToUpper(info.Extension) {
	case "JPG", "JPEG", "PNG":
		if data, err := ioutil.ReadAll(file); err != nil {
			p.API.LogError(err.Error())
			return nil, ""
		} else {
			jmp := jpegstructure.NewJpegMediaParser()
			pmp := pngstructure.NewPngMediaParser()
			mc := &MediaContext{
				MediaType: OtherMediaType,
				RootIfd:   nil,
				RawExif:   nil,
				Media:     nil,
			}

			if jmp.LooksLikeFormat(data) {
				mc.MediaType = JpegMediaType
			} else if pmp.LooksLikeFormat(data) {
				mc.MediaType = PngMediaType
			}

			switch mc.MediaType {
			case JpegMediaType:
				sl, _ := jmp.ParseBytes(data)
				if err != nil {
					return nil, ""
				}

				mc.Media = sl

				rootIfd, rawExif, err := sl.Exif()
				if err != nil {
					return nil, ""
				}

				mc.RootIfd = rootIfd
				mc.RawExif = rawExif

			case PngMediaType:
				cs, err := pmp.ParseBytes(data)
				if err != nil {
					return nil, ""
				}

				mc.Media = cs

				rootIfd, rawExif, err := cs.Exif()
				if err != nil {
					return nil, ""
				}

				mc.RootIfd = rootIfd
				mc.RawExif = rawExif
			default:
				return nil, ""
			}

			entries := p.extractEXIF(mc)

			if data, err := json.MarshalIndent(entries, "", "    "); err != nil {
				p.API.LogError(err.Error())
				return nil, ""
			} else {
				p.API.LogInfo(string(data))
			}

			//TODO remove these entries
		}

	}
	return nil, ""
}

func (p *Plugin) extractEXIF(mc *MediaContext) (entries []IfdEntry) {
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	entries = make([]IfdEntry, 0)
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {
		defer func() {
			if state := recover(); state != nil {
				p.API.LogError(state.(error).Error())
			}
		}()

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
		if err != nil {
			return err
		}

		it, err := ti.Get(ifdPath, tagId)
		if err != nil {
			return err
		}

		valueString := ""
		var value interface{}
		if tagType.Type() == exif.TypeUndefined {
			value, err = exif.UndefinedValue(ifdPath, tagId, valueContext, tagType.ByteOrder())
			if err != nil {
				return err
			} else {
				valueString = fmt.Sprintf("%v", value)
			}
		} else {
			valueString, err = tagType.ResolveAsString(valueContext, true)
			if err != nil {
				return err
			}
			value = valueString
		}

		entry := IfdEntry{
			IfdPath:     ifdPath,
			FqIfdPath:   fqIfdPath,
			IfdIndex:    ifdIndex,
			TagId:       tagId,
			TagName:     it.Name,
			TagTypeId:   tagType.Type(),
			TagTypeName: tagType.Name(),
			UnitCount:   valueContext.UnitCount,
			Value:       value,
			ValueString: valueString,
		}

		entries = append(entries, entry)

		return nil
	}

	exif.Visit(exif.IfdStandard, im, ti, mc.RawExif, visitor)

	return entries
}
