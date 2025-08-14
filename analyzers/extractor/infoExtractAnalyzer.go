package extractor

import (
	"ForensicsTool/analyzers/extractor/tool"
	"context"
	"fmt"
	"github.com/iancoleman/orderedmap"
)

type InfoExtractor struct {
	ctx context.Context
}

func NewInfoExtractor() *InfoExtractor {
	return &InfoExtractor{}
}

type ExtractResult struct {
	Data []map[string]string `json:"data"`
	Err  string              `json:"err"`
}

type ExtractMapResult struct {
	Data *orderedmap.OrderedMap `json:"data"`
	Err  string                 `json:"err"`
}

func (p *InfoExtractor) InitCtx(ctx context.Context) {
	p.ctx = ctx
}

func (p *InfoExtractor) ExtractXShell(folder, sid string) *ExtractMapResult {
	/*
		@param folder:	配置文件目录
		@param sid:		用户的用户名+sid
		@return:		解析结果
	*/
	info, err := tool.AnalyzeXshell(folder, sid)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractNavicat(refPath string) *ExtractMapResult {
	/*
		@param	refPath:	NTUSER.DAT注册表文件路径
		@return:			解析结果
	*/
	info, err := tool.AnalyzeNavicat(refPath)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractMobaXterm(filePath, masterPwd string) *ExtractMapResult {
	/*
		@param	filePath:	注册表或配置文件路径
		@param	masterPwd:	主密码
		@return:			解析结果
	*/
	info, err := tool.AnalyzeMobaXterm(filePath, masterPwd)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractDbeaver(folder string) *ExtractMapResult {
	/*
		@param	folder:	包含credentials-config.json和data-sources.json的文件夹
		@return:		解析结果
	*/
	info, err := tool.AnalyzeDbeaver(folder)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractFinalShell(folder string) *ExtractMapResult {
	/*
		@param	folder:	包含连接信息的文件夹
		@return:		解析结果
	*/
	info, err := tool.AnalyzeFinalShell(folder)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractHawk2(filePath, pwd string) *ExtractMapResult {
	/*
		@param	filePath:	Hawk2.xml的文件路径
		@param	pwd:		cipher_key，保存在crypto.KEY_256.xml或crypto.KEY_128.xml中
		@return:			解析结果
	*/
	info, err := tool.AnalyzeHawk2(filePath, pwd)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}

func (p *InfoExtractor) ExtractMetaMask(filePath string) *ExtractMapResult {
	/*
		@param	filePath:	persist-root文件路径
		@return:			解析结果
	*/
	info, err := tool.AnalyzeMetaMask(filePath)
	if err != nil {
		return &ExtractMapResult{nil, fmt.Sprintf("解析失败：%v", err)}
	}
	return &ExtractMapResult{info, ""}
}
