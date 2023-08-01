package test

import (
	"context"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"io/ioutil"
	"net/http"
	"testing"
)

// 测试用户注册接口
func TestUserRegister(t *testing.T) {
	url := serverAddr + "/tiktok/user/register?username=test&password=123456"
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

// 测试用户登录接口
func TestUserLogin(t *testing.T) {
	url := serverAddr + "/tiktok/user/login?username=test&password=123456"
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

// 测试用户信息接口
func TestUserInfo(t *testing.T) {
	token := userToken
	url := serverAddr + "/tiktok/user?user_id=1014&token=" + token
	method := "GET"

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

func TestUserRegisterService(t *testing.T) {
	c, err := userservice.NewClient("user",
		client.WithHostPorts("127.0.0.1:8889"))
	if err != nil {
		t.Errorf("New client error: %#v", err)
		return
	}

	resp, err := c.UserRegister(context.Background(), &user.UserRegisterRequest{
		Username: "test",
		Password: "123456",
	})
	if err != nil {
		t.Errorf("user register error: %#v\n", err)
		return
	}
	fmt.Printf("%#v\n", resp)
}

func TestUserLoginService(t *testing.T) {
	c, err := userservice.NewClient("user",
		client.WithHostPorts("127.0.0.1:8889"))
	if err != nil {
		t.Errorf("New Client error: %v", err)
		return
	}

	resp, err := c.UserRegister(context.Background(), &user.UserRegisterRequest{
		Username: "test",
		Password: "123456",
	})
	if err != nil {
		t.Errorf("user register error: %#v\n", err)
		return
	}
	fmt.Printf("%#v\n", resp)
}
