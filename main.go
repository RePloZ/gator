package main

import (
	"fmt"

	"github.com/RePloZ/gator/internal/config"
)

func main() {
	cfg := config.Config{}
	_ = cfg.Read()
	cfg.SetUser("lane")
	_ = cfg.Read()
	fmt.Print(cfg)
}
