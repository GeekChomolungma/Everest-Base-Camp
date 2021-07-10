package handlerfactory

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/huobirdcenter/huobi_golang/logging/perflogger"
)

type HuoBiFactory struct {
}

func (HBFac *HuoBiFactory) Create() BaseHandler {
	return &HuoBiHandler{}
}

type HuoBiHandler struct {
}

func (HBF *HuoBiHandler) Get(url string) (string, error) {
	logger := perflogger.GetInstance()
	logger.Start()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	logger.StopAndLog("GET", url)

	return string(result), err
}

func (HBF *HuoBiHandler) Post(url string, body string) (string, error) {
	logger := perflogger.GetInstance()
	logger.Start()

	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	logger.StopAndLog("POST", url)

	return string(result), err
}
