package tool

import (
	"ForensicsTool/utils"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
	"github.com/tidwall/gjson"
)

func AnalyzeMetaMask(filePath string) (*orderedmap.OrderedMap, error) {
	/*
		@param	filePath:	persist-root文件
		@return:			解析结果、错误
	*/
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		logger.Errorf("%s is not a file", filePath)
		return nil, fmt.Errorf("%s is not a file", filePath)
	}
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var walletInfoList []*orderedmap.OrderedMap
	var contactInfoList []*orderedmap.OrderedMap
	var transactionInfoList []*orderedmap.OrderedMap
	var browserHistoryList []*orderedmap.OrderedMap
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	json := gjson.ParseBytes(data)
	if json.Get("engine").Exists() {
		engine := json.Get("engine")
		browser := json.Get("browser")
		for _, arr := range engine.Get("backgroundState.AccountTrackerController.accounts").Array() {
			walletInfo := orderedmap.New()
			walletInfo.SetEscapeHTML(false)
			s := engine.Get(fmt.Sprintf("backgroundState.PreferencesController.identities.%s.importTime", arr.String())).Int()
			//t := time.UnixMicro(s)
			//importTime := t.Format("2006-01-02 15:04:05")
			importTime := utils.DefaultTimestampToDatetime(s, "", "")
			balance, err := strconv.ParseInt(engine.Get(fmt.Sprintf("backgroundState.AccountTrackerController.accounts.%s.balance", arr.String())).String(), 16, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			walletInfo.Set("添加时间", importTime)
			walletInfo.Set("钱包名称", engine.Get(fmt.Sprintf("backgroundState.PreferencesController.identities.%s.name", arr)).String())
			walletInfo.Set("钱包地址", arr.String())
			walletInfo.Set("余额（ETH）", strconv.FormatInt(balance/(10^18), 10))
			walletInfoList = append(walletInfoList, walletInfo)
		}
		for _, contactId := range engine.Get("backgroundState.AccountTrackerController.addressBook").Array() {
			for _, contact := range contactId.Get(contactId.String()).Array() {
				contactInfo := orderedmap.New()
				contactInfo.SetEscapeHTML(false)
				contactInfo.Set("账户名", contact.Get("name").String())
				contactInfo.Set("钱包地址", contact.Get("address").String())
				contactInfoList = append(contactInfoList, contactInfo)
			}
		}
		for _, transaction := range engine.Get("backgroundState.TransactionController.transactions").Array() {
			transactionInfo := orderedmap.New()
			transactionInfo.SetEscapeHTML(false)
			s := transaction.Get("time").Int()
			//t := time.UnixMicro(s)
			transactionTime := utils.DefaultTimestampToDatetime(s, "", "")
			value, err := strconv.ParseInt(transaction.Get("txParams.value").String(), 16, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			transactionValue := value / (10 ^ 18)
			transactionInfo.Set("交易时间", transactionTime)
			transactionInfo.Set("发送方", transaction.Get("txParams.from").String())
			transactionInfo.Set("接收方", transaction.Get("txParams.to").String())
			transactionInfo.Set("交易金额（ETH）", strconv.FormatInt(transactionValue, 10))
			transactionInfo.Set("交易哈希", transaction.Get("hash").String())
			transactionInfo.Set("错误", transaction.Get("error").String())
			transactionInfoList = append(transactionInfoList, transactionInfo)
		}
		for _, history := range browser.Get("history").Array() {
			historyInfo := orderedmap.New()
			historyInfo.SetEscapeHTML(false)
			historyInfo.Set("名称", history.Get("name").String())
			historyInfo.Set("URL", history.Get("url").String())
			browserHistoryList = append(browserHistoryList, historyInfo)
		}
		result.Set("钱包信息", walletInfoList)
		result.Set("账户信息", contactInfoList)
		result.Set("交易信息", transactionInfoList)
		result.Set("浏览历史", browserHistoryList)
	}
	return result, nil
}
