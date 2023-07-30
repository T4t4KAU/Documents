package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAddNewFavorite(t *testing.T) {
	url := serverAddr + fmt.Sprintf("/douyin/favorite/action?token=%s&video_id=1014&action_type=1", userToken)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
