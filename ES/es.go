package ES

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"storageApi/conf"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

func getMetaData(name string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/objects_metadata/objects/%s_%d?_source", conf.GetConfig().Env.Es, name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hit struct {
		Total int
		Hits  []hit
	}
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	esUrl := fmt.Sprintf("http://%s/objects_metadata/_search?q=name:%s&size=1&sort=version:desc", conf.GetConfig().Env.Es, url.PathEscape(name))
	r, e := http.Get(esUrl)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hit.Hits) != 0 {
		meta = sr.Hit.Hits[0].Source
	}
	return
}

func GetMetaData(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetaData(name, version)
}

func PutMetaData(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/objects_metadata/objects/%s_%d?op_type=create", conf.GetConfig().Env.Es, name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetaData(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name string, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetaData(name, version.Version+1, size, hash)
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/objects_metadata/_search?sort=name,version&from=%d&dsize=%d", conf.GetConfig().Env.Es, from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hit.Hits {
		metas = append(metas, sr.Hit.Hits[i].Source)
	}
	return metas, nil
}
