package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
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
	StartBytes    = 4
	EndBytes      = 4
	OffsetBytes   = 4
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
			_, _, err = image.Decode(bytes.NewReader(data))
			if err != nil {
				p.API.LogInfo("ERROR: original image is corrupt" + err.Error() + "\n")
				return nil, ""
			}
			mc, entries := p.parseEXIF(data)
			if len(entries) == 0 {
				return nil, ""
			}

			filteredBytes := []byte{}

			switch mc.MediaType {
			case JpegMediaType:
				if filteredBytes, err = p.extractJPEGEXIF(mc, data, filteredBytes); err != nil {
					p.API.LogError(err.Error())
					return nil, ""
				}
			}

			if _, err := output.Write(filteredBytes); err != nil {
				p.API.LogError(err.Error())
			}
		}
	}
	return nil, ""
}

func (p *Plugin) extractJPEGEXIF(mc *MediaContext, data []byte, filtered []byte) ([]byte, error) {
	sl := mc.Media.(*jpegstructure.SegmentList)
	_, sExif, err := sl.FindExif()
	if err != nil {
		return nil, errors.New("No EXIF in image")
	}

	p.API.LogInfo(fmt.Sprintf("****(exif) %x %s %x", sExif.Offset, sExif.MarkerName, len(sExif.Data)))

	bytesCount := 0
	startExifBytes := StartBytes
	endExifBytes := EndBytes
	for _, s := range sl.Segments() {

		if s.MarkerName == sExif.MarkerName {
			if startExifBytes == StartBytes {
				startExifBytes = bytesCount
				endExifBytes = startExifBytes + len(s.Data) + OffsetBytes
			} else {
				endExifBytes += len(s.Data) + OffsetBytes
			}
		}
		bytesCount += len(s.Data) + OffsetBytes

		p.API.LogInfo(fmt.Sprintf("%x %s %v (%x)", s.Offset, s.MarkerName, len(s.Data), s.Offset+len(s.Data)))

	}

	filtered = data[:startExifBytes]
	filtered = append(filtered, data[endExifBytes:]...)

	p.API.LogInfo(fmt.Sprintf("********(size) %v %v  (%v)", len(data), len(filtered), len(data)-len(filtered)))

	//TODO enable when > 95% effective EXIF without corruption
	//_, _, err = image.Decode(bytes.NewReader(filtered))
	//if err != nil {
	//	return nil, errors.New("EXIF extraction corrupted image " + err.Error())
	//}

	return filtered, nil
}

func (p *Plugin) parseEXIF(data []byte) (mc *MediaContext, entries []IfdEntry) {

	jmp := jpegstructure.NewJpegMediaParser()
	pmp := pngstructure.NewPngMediaParser()
	mc = &MediaContext{
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
		sl, err := jmp.ParseBytes(data)
		if err != nil {
			return mc, []IfdEntry{}
		}

		mc.Media = sl

		rootIfd, rawExif, err := sl.Exif()
		if err != nil {
			return mc, []IfdEntry{}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif

	case PngMediaType:
		cs, err := pmp.ParseBytes(data)
		if err != nil {
			return mc, []IfdEntry{}
		}

		mc.Media = cs

		rootIfd, rawExif, err := cs.Exif()
		if err != nil {
			return mc, []IfdEntry{}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	default:
		return mc, []IfdEntry{}
	}

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

	return mc, entries
}
