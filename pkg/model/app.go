package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chengyumeng/khadijah/pkg/config"
	"github.com/chengyumeng/khadijah/pkg/utils/log"
)

func GetAppBody(nsId int64) *AppBody {
	url := fmt.Sprintf("%s/%s/%d/%s?pageSize=%d", config.GlobalOption.System.BaseURL, "api/v1/namespaces", nsId, "apps", PageSize)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+config.GlobalOption.Token)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.AppLogger.Warning(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return nil
	}
	data := new(AppBody)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.AppLogger.Warning(err)
	}
	return data
}
