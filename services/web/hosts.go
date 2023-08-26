package web

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Host struct {
	Address string `json:"address"`
	Port    string `json:"runningPort"`
}

type Hosts struct {
	Hosts            []Host `json:"hosts"`
	curIndex, length uint
}

func GetHostsFromFile(fileName string) (*Hosts, error) {
	var hosts Hosts
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &hosts)
	hosts.curIndex = 0
	hosts.length = uint(len(hosts.Hosts))
	return &hosts, nil
}

func (hst *Hosts) GetNext() (*Host, error) {
	if hst.length == 0 {
		return nil, errors.New("Empty list of hosts")
	}
	hst.curIndex = (hst.curIndex+1) % hst.length
	return &hst.Hosts[hst.curIndex], nil
}