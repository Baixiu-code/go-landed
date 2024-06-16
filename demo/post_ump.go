package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	//post url
	domain := "http://open.ump.jd.com"
	uri := "/operateMonitorData"
	url := domain + uri
	fmt.Println("请求 url:", url)
	key := "cartsoa.function.beta.frameWorkCartResult.resultCode.error.CART_UNCHECK_ALL.-2.app"
	//请求 data, operateType:method .操作方法监控点.
	dataMap := getOperateMonitorData(key, "update")
	//map to json
	jsonDataMap, err := json.Marshal(dataMap)
	// json to string
	strJsonDataMapStr := string(jsonDataMap)
	strJsonDataMapStr = unescapeJSONString(strJsonDataMapStr)
	//解决类型不匹配问题
	strJsonDataMapStr = strings.Replace(strJsonDataMapStr, "\"frequency\":\"1\"", "\"frequency\":1", -1)
	//解决 json 格式问题
	strJsonDataMapStr = strings.Replace(strJsonDataMapStr, "\"{\"critical\"", "{\"critical\"", -1)
	strJsonDataMapStr = strings.Replace(strJsonDataMapStr, "}}}\"", "}}}", -1)
	fmt.Println("请求 req body:", strJsonDataMapStr)
	// string to byte
	data := []byte(strJsonDataMapStr)
	//创建一个 post 请求
	var req = getRequest(data, uri)
	client := &http.Client{}
	//发送 post 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	//close read 流
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭 read 失败:", err)
		}
	}(resp.Body)
	//读取 响应 body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取返回失败:", err)
		return
	}
	fmt.Println("请求响应:", string(body))
}

/*
*
get req body by umpKey.
通过 ump key 来修改报警配置.
*/
func getOperateMonitorData(umpKey string, typeUmp string) map[string]string {
	//根据 ump key 读取 alarm key 下的配置
	alarmMap := getOperateMonitorDataFromAlarm(umpKey)
	//map to json
	alarmMapJson, _ := json.Marshal(alarmMap)
	//json to str
	alarmMapStr := string(alarmMapJson)
	alarmMapStr = unescapeJSONString(alarmMapStr)
	fmt.Println("getOperateMonitorData.alarmMapStr:", alarmMapStr)
	operateMonitorData := map[string]string{
		"operateType": "method",
		"type":        typeUmp,
		"key":         umpKey,
		"appName":     "appcartsoa",
		"platform":    "jdos",
		"frequency":   "1",
		"alarm":       alarmMapStr,
	}

	return operateMonitorData

}

/*
get operate monitor data from alarm .构造 alarm 节点
*/
func getOperateMonitorDataFromAlarm(umpKey string) map[string]map[string]map[string]int {
	//get critical map
	criticalMap := getCriticalMap(umpKey)
	//map to json
	operateMonitorDataFromAlarmData := map[string]map[string]map[string]int{
		"critical": criticalMap,
	}
	return operateMonitorDataFromAlarmData
}

/*
*
// (可选) 报警级别 warning、critical
// - 如果是update操作，alarm对象及其子对象warning、critical等可以置为空对象{}，意为删除原有该对象。
// - warning、critical单独设置其中任意一个，另外一个保持不变。
*/
func getCriticalMap(umpKey string) map[string]map[string]int {

	invokeCountMap := getInvokeCount(umpKey)
	alarmConfigMap := getAlarmConfigMap(umpKey)
	criticalMap := map[string]map[string]int{
		"invokeCount": invokeCountMap,
		"alarmConfig": alarmConfigMap,
	}
	return criticalMap
}

func getAlarmConfigMap(umpKey string) map[string]int {
	alarmConfigMap := map[string]int{
		"alarmSwitch": 1,
		"sign":        1,
		"alarmWay":    1,
	}
	return alarmConfigMap
}

/*
*
// (可选) 方法性能报警配置
*/
func getPerformanceMap(umpKey string) map[string]string {
	performance := map[string]string{
		"tp50":     "300", // (可选) tp50阈值，但至少配置一项性能指标，毫秒值
		"tp90":     "400", // (可选) tp90阈值，但至少配置一项性能指标，毫秒值
		"tp99":     "500", // (可选) tp99阈值，但至少配置一项性能指标，毫秒值
		"tp999":    "600", // (可选) tp999阈值，但至少配置一项性能指标，毫秒值
		"max":      "700", // (可选) max最大值阈值，但至少配置一项性能指标，毫秒值
		"avg":      "500", // (可选) avg平均值阈值，但至少配置一项性能指标，毫秒值
		"count":    "5",   // (必选) 连续出现次数，以上任意一个性能指标连续超过此配置次数时报警
		"interval": "5",   // (必选) 报警收敛，单位分钟，在配置的分钟之内只报一次警
	}
	return performance
}

/*
*
(可选) 方法调用次数报警
*/
func getInvokeCount(umpKey string) map[string]int {
	invokeCountMap := map[string]int{
		"count":     5,   // (必选) 连续出现次数
		"range":     10,  // (必选) 调用次数取值时间范围，5分钟的倍数，
		"operator":  1,   // (必选) 比较方式：(1:大于等于 2:小于等于 3:等于 4:大于 5:小于)
		"threshold": 100, // (必选) 调用次数阈值
		"interval":  5,   // (必选) 报警收敛，单位分钟，在配置的分钟之内只报一次警
	}
	return invokeCountMap
}

/*
get ump key 生效失效状态
*/
func getOperateKeyStatusReqBodyData(umpKey string) map[string]string {
	operateKeyStatusMap := map[string]string{
		"operateType":   "OperateKeyStatus",
		"appName":       "jdos_appcartsoa",
		"scope":         "Key",
		"scopeValues":   umpKey,
		"platform":      "jdos",
		"operateStatus": "enable",
	}
	return operateKeyStatusMap
}

/*
*
get request by body and uri
*/
func getRequest(body []byte, uri string) *http.Request {
	//post url
	domain := "http://open.ump.jd.com"
	url := domain + uri
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	token := "jwt::eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjaGVuZmFuZ2xpbjEiLCJpc3MiOiJ1bXAiLCJzdWIiOiJ7XCJlcnBcIjpcImNoZW5mYW5nbGluMVwiLFwicXBtXCI6MjAwLFwic3RhcnRUaW1lc3RhbXBcIjoxNzE4MjY4MDU2OTc1LFwiZW5kVGltZXN0YW1wXCI6MTcxODg3Mjg1Njk3NSxcImJvdW5kQXBwQ29vcmRzXCI6W1wiYXBwY2FydHNvYShqZG9zKVwiXX0ifQ.yuIY7r0fzljdgHTzsVSMByFC5T-ybFNBW_0p7IbHQ4M"
	req.Header.Set("token", token)
	fmt.Println("请求 header:", req.Header)
	return req
}

// 去除 JSON 字符串中的转义字符
func unescapeJSONString(s string) string {
	escaped := []byte(s)
	unescaped := make([]byte, 0, len(escaped))

	for i := 0; i < len(escaped); i++ {
		if escaped[i] == '\\' && i+1 < len(escaped) {
			switch escaped[i+1] {
			case '"', '\\', '/', '\'':
				unescaped = append(unescaped, escaped[i+1])
				i++
			case 'b':
				unescaped = append(unescaped, '\b')
				i++
			case 'f':
				unescaped = append(unescaped, '\f')
				i++
			case 'n':
				unescaped = append(unescaped, '\n')
				i++
			case 'r':
				unescaped = append(unescaped, '\r')
				i++
			case 't':
				unescaped = append(unescaped, '\t')
				i++
			case 'u':
				// 处理 Unicode 转义字符
				if i+5 < len(escaped) {
					hex := escaped[i+2 : i+6]
					unicode := rune((hexDigitToInt(hex[0]) << 12) | (hexDigitToInt(hex[1]) << 8) | (hexDigitToInt(hex[2]) << 4) | hexDigitToInt(hex[3]))
					unescaped = append(unescaped, []byte(string(unicode))...)
					i += 5
				} else {
					unescaped = append(unescaped, escaped[i])
				}
			default:
				unescaped = append(unescaped, escaped[i])
			}
		} else {
			unescaped = append(unescaped, escaped[i])
		}
	}

	return string(unescaped)
}

// 将 16 进制数字转换为整数
func hexDigitToInt(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10
	case 'A' <= c && c <= 'F':
		return int(c-'A') + 10
	default:
		return 0
	}
}
