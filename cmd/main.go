package main

import (
	"fmt"
	"repo_manager/internal/admin"
	config2 "repo_manager/internal/config"
	"repo_manager/internal/telegram_bot"
	"sync"
)

func main() {

	config := config2.Config{}
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		admin.Init(&config)
	}()

	go func() {
		defer wg.Done()
		telegram_bot.Init(&config)
	}()

	wg.Wait()
}
