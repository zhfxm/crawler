package pool

import (
	"encoding/json"
	"fmt"

	"github.com/zhfxm/simple-crawler/collect"
	"github.com/zhfxm/simple-crawler/collector"
	"go.uber.org/zap"
)

type AntpoolCurrent struct {
	Code string              `json:"code"`
	Msg  string              `json:"msg"`
	Data *AntpoolCurrentData `json:"data"`
}

type AntpoolCurrentData struct {
	UserId       string `json:"userId"`
	HsNow        string `json:"hsNow"`
	HsNowUnit    string `json:"hsNowUnit"`
	HsLast1D     string `json:"hsLast1D"`
	HsLast1DUnit string `json:"hsLast1DUnit"`
}

func AntpoolCurrentParseFunc(contents []byte, req *collect.Request) collect.ParseResult {
	res := collect.ParseResult{}
	res.Items = []interface{}{}
	hs := &AntpoolCurrent{}
	json.Unmarshal(contents, hs)
	if hs.Code != "000000" {
		req.Logger.Error("antpool current hash rate error", zap.String("code", hs.Code), zap.String("msg", hs.Msg), zap.String("url", req.Url))
		return res
	}
	// hsInfo := &pojo.CurrentHashrateInfo{}
	// hsInfo.HsNow = hs.Data.HsNow
	// hsInfo.HsNowUnit = hs.Data.HsNowUnit
	// hsInfo.Hs24 = hs.Data.HsLast1D
	// hsInfo.Hs24Unit = hs.Data.HsLast1DUnit

	cellMap := make(map[string]interface{})
	dataCell := &collector.DataCell{}
	cellMap["hsNow"] = hs.Data.HsNow
	cellMap["hsNowUnit"] = hs.Data.HsNowUnit
	cellMap["Hs24"] = hs.Data.HsLast1D
	cellMap["Hs24Unit"] = hs.Data.HsLast1DUnit
	dataCell.Data = cellMap
	fmt.Println("iterm result....")
	return collect.ParseResult{
		Items: []interface{}{dataCell},
	}
}
