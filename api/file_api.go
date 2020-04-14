package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

type File struct {
	Id           string `json:"id"`
	Ext          string // filepath.Ext(Name)
	Name         string `json:"name"`
	Size         string `json:"size"`
	MimeType     string `json:"mimeType"`
	ModifiedTime string `json:"modifiedTime"` // RFC3339 TODO: can be Cached

	IsFolder   bool // end with `folder`
	IsAudio    bool // mimeType has prefix `audio`
	IsVideo    bool // mimeType has prefix `video`
	IsPlayable bool // == IsFolder + IsAudio
	IsText     bool // mimeType has prefix `text`
	IsImage    bool //  mimeType has prefix `image`
	IsGApp     bool //  mimeType has prefix `application`

}

func ShowDir(abPath string, auth string, rootDrive string) (bool, []*File) {
	if !strings.HasSuffix(abPath, "/") {
		log.Fatalln("path must end with `/`")
	}

	resp := Do("POST", "ShowDir", abPath, url.Values{"rootId": []string{rootDrive}}.Encode(), auth)

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if string(b) == "null" {
		return false, nil
	}

	a := struct {
		Files []*File `json:"files"`
	}{}

	if err := json.Unmarshal(b, &a); err != nil {
		log.Fatal(err)
	} else {
		for _, f := range a.Files {
			if strings.HasSuffix(f.MimeType, "folder") {
				f.IsFolder = true
				continue
			}

			f.Ext = filepath.Ext(f.Name)
			mime := strings.Split(f.MimeType, "/")[0]
			switch mime {
			case "audio":
				f.IsAudio = true
			case "video":
				f.IsVideo = true
			case "text":
				f.IsText = true
			case "image":
				f.IsImage = true
			case "application":
				f.IsGApp = true

			}
			f.IsPlayable = f.IsVideo || f.IsAudio
		}
		return true, a.Files
	}
	return false, nil
}
