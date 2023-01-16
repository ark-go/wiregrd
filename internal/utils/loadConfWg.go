package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	//"reflect"
	"strings"
)

// для интерфейсов
// type peer struct {
// 	PublicKey          string
// 	Endpoint           string
// 	AllowedIps         string
// 	LatestHandshake    string
// 	Transfer           string
// 	PersistenKeepAlive string
// }

type peer map[string]string

// для пиров
type ifaceinfo struct {
	//Name string
	// PublicKey     string
	// PrivateKey    string
	// ListeningPort string
	Info      map[string]string
	Peers     []peer
	currpeers peer
}

// рабочая структура
type confwg struct {
	curiface *ifaceinfo   // текщий интерфейс
	ifices   []*ifaceinfo //все интерфейсы
	curstep  int
}

const (
	constIface int = iota
	constPeer
)

// Добавляем/создаем интерфейс
func (c *confwg) addIface(name string) {
	ii := &ifaceinfo{
		//Name: name,
		Info: map[string]string{},
	}

	c.ifices = append(c.ifices, ii)
	c.curiface = ii // добавили и сделалитекущим
	return
}
func (c *confwg) addPeer(pubKey string) {
	// pe := &peer{
	// 	PublicKey: pubKey,
	// }
	pe := map[string]string{}
	c.curiface.Peers = append(c.curiface.Peers, pe)
	//pe := map[string]string{"public key": pubKey}
	//c.curiface.peers = append(c.curiface.peers, pe)
	c.curiface.currpeers = pe // добавили и сделалитекущим
	return
}

func LoadConf(name string) (*confwg, error) {
	cmd := exec.Command("sudo", "/bin/sh", "-c", name, "show")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		log.Printf("error: %s: %v\n", name, err)
		return nil, err
	}
	return LoadLinesBuffer(buf)
}

func LoadLinesBuffer(buf bytes.Buffer) (*confwg, error) {
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	cfg := &confwg{}
	for scanner.Scan() {
		CreateStruct(cfg, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return cfg, nil
}

func CreateStruct(cfg *confwg, str string) {
	words := strings.SplitN(str, ":", 2)
	if len(words) < 2 {
		log.Println("+++++++++++>", words)
		return
	}
	str1 := strings.ToLower(strings.Trim(words[0], " "))
	str2 := strings.Trim(words[1], " ")

	switch str1 {
	case "interface":
		log.Println("iface:", str2)
		cfg.addIface(str2)
		cfg.curstep = constIface

	//-------------------------------------------
	case "peer":
		cfg.addPeer(str2)
		cfg.curstep = constPeer
	}
	if cfg.curstep == constPeer {
		cfg.curiface.currpeers[str1] = str2
	}
	if cfg.curstep == constIface {
		cfg.curiface.Info[str1] = str2
	}
	log.Println(words[0], " >> ", words[1])
	//PrintCfg(cfg)

}

func PrintCfg(cfg *confwg) {
	for _, v := range cfg.ifices {
		//log.Println("interface:", v.Name)
		log.Println("- - -- - - -")
		for k, v := range v.Info {

			log.Printf("%s:\t\t%s\n", k, v)
		}
		for _, p := range v.Peers {
			log.Println("peer: _________")
			for k, v := range p {

				log.Printf("%s:\t\t%s\n", k, v)
			}
		}

		// for _, vv := range v.peers {
		// 	//log.Println(vv.AllowedIps)
		// 	// values := reflect.ValueOf(*vv)
		// 	// types := values.Type()
		// 	// for i := 0; i < values.NumField(); i++ {
		// 	// 	log.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
		// 	// }
		// }
	}
	jsonString, err := json.Marshal(cfg.ifices)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(jsonString))
	}
}
