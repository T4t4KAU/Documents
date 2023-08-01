package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRelationAction(t *testing.T) {
	url := serverAddr + "/tiktok/relation/action?token=" + userToken + "&to_user_id=1003&action_type=1"
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
