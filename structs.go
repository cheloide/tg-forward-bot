package main

import "github.com/hellowearemito/go-telegram-structs"

type GetUpdatesResponse struct {
	Ok     bool              `json:"ok"`
	Result []telegram.Update `json:"result"`
}

type Settings struct {
	Token       string `json:"token"`
	ForwardFrom int64  `json:"forwardFrom"`
	ForwardTo   int64  `json:"forwardTo"`
}
