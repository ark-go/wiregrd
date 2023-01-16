package main

import (
	"log"

	"github.com/ark-go/wiregrd/internal/utils"
	"github.com/ark-go/wiregrd/internal/wgkey"
)

func init() {
}

func main() {
	keys, err := wgkey.GetKeys()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(keys.PrivateKey)
		log.Println(keys.PublicKey)
	}
	keys1, err := wgkey.GetKeysFromPrivate("0ABlbm7qQkIlKjlhm3d041kCthFElwgsQBnu2W9x9kQ=") //keys.PrivateKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("- - - - ")
	log.Println(keys1.PrivateKey)
	log.Println(keys1.PublicKey)
	cfg, err := utils.LoadConf("wg")
	if err != nil {
		log.Println(err.Error())
		return
	}
	utils.PrintCfg(cfg)
	//log.Println(conf)
}
