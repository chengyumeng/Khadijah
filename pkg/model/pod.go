package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chengyumeng/khadijah/pkg/config"
)

func GetPodBody(appId int64, podType string) *PodBody {
	url := fmt.Sprintf("%s/%s/%d/%ss?pageSize=%d", config.GlobalOption.System.BaseURL, "api/v1/apps", appId, podType, PageSize)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+config.GlobalOption.Token)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warning(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return nil
	}
	data := new(PodBody)
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Warning(err)
	}
	return data
}
