package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/dsoprea/go-exif"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
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
	case "JPG", "TIFF", "PNG", "JP2", "PGF", "MIFF", "HDP", "PSP", "XCF":
		if data, err := ioutil.ReadAll(file); err != nil {
			p.API.LogError(err.Error())
			return nil, ""
		} else {
			if rawExif, err := exif.SearchAndExtractExif(data); err != nil {
				p.API.LogError(err.Error())
				return nil, ""
			} else {
				im := exif.NewIfdMappingWithStandard()
				ti := exif.NewTagIndex()

				entries := make([]IfdEntry, 0)
				visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {
					defer func() {
						if state := recover(); state != nil {
							p.API.LogError(state.(error).Error())
						}
					}()

					ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
					if err != nil {
						p.API.LogError(err.Error())
						return err
					}

					it, err := ti.Get(ifdPath, tagId)
					if err != nil {
						p.API.LogError(err.Error())
						return err
					}

					valueString := ""
					var value interface{}
					if tagType.Type() == exif.TypeUndefined {
						var err2 error
						value, err2 = exif.UndefinedValue(ifdPath, tagId, valueContext, tagType.ByteOrder())
						if err2 != nil {
							p.API.LogError(err2.Error())
							return err2
						} else {
							valueString = fmt.Sprintf("%v", value)
						}
					} else {
						valueString, err = tagType.ResolveAsString(valueContext, true)
						if err != nil {
							p.API.LogError(err.Error())
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

				_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
				if err != nil {
					p.API.LogError(err.Error())
					return nil, ""
				}

				//TODO strip out the EXIF information and save file
				if data, err := json.MarshalIndent(entries, "", "    "); err != nil {
					p.API.LogError(err.Error())
					return nil, ""
				} else {
					p.API.LogInfo(string(data))
				}

			}

		}

	}
	return nil, ""
}
