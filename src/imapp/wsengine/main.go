package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ws/business"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/config"
	"ws/messenger"
)

type result struct{}

func (r result) ProcessResult(iocProvider containerInterface.IAppIocContainer, result *containerInterface.MessageResult) {
}

func (r result) OnDestroy() {}

func main() {
	//registerBusiness()
	loginBusinessIM()
	//registerMessenger()
	//loginMessenger()
}

func registerBusiness() {
	for true {
		var phoneNumber, code string

		fmt.Println("输入手机号：")
		_, err := fmt.Scanln(&phoneNumber)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		client, _ := business.NewAnonymousClient(phoneNumber)
		err = client.GetSmsCode()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("输入验证码：")
		_, err = fmt.Scanln(&code)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = client.SendReceiveSmsCode(code)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func loginBusinessIM() {
	client := business.NewIMClient("66970241918")

	opts := []config.OptionsFn{
		config.AutoMessageMarkRead(true),
		//config.UseSocks("gate3.rola.info", "2129", "Tt1314_54321-country-HK", "Aa1122334455"),
	}
	err := client.Connect(opts...)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	// 登录是异步的 等待这个时间
	//time.Sleep(time.Second * 10)

	//进行scene操作
	//s := scene.NewSession("85295664081")
	//s.MakeTextMessage("goodbye")
	//client.EnterScene(&s, &result{})

	//s := scene.NewUserSettingsScene()
	//s.QueryNickname()
	//s.QuerySignature()
	//client.EnterScene(&s, &result{})

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	client.Logout()
}

func registerMessenger() {
	client, _ := messenger.NewAnonymousClient("85256048573")
	//client.GetSmsCode()
	client.SendReceiveSmsCode("129334")
}

func loginMessenger() {
	client := messenger.NewIMClient("84327495821")

	opts := []config.OptionsFn{
		config.AutoMessageMarkRead(true),
		//config.UseSocks("gate3.rola.info", "2129", "Tt1314_54321-country-HK", "Aa1122334455"),
	}
	err := client.Connect(opts...)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	// 登录是异步的 等待这个时间
	//time.Sleep(time.Second * 10)

	//进行scene操作
	//s := scene.NewSession("85295664081")
	//s.MakeTextMessage("goodbye")
	//client.EnterScene(&s, &result{})

	//s := scene.NewUserSettingsScene()
	//s.QueryNickname()
	//s.QuerySignature()
	//client.EnterScene(&s, &result{})

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	client.Logout()
}
