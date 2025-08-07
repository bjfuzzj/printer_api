package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"syscall"
	"unsafe"
)

// 打印机状态常量定义
const (
	// 状态组
	GROUP_USUALLY   = 0x00010000
	GROUP_SETTING   = 0x00020000
	GROUP_HARDWARE  = 0x00040000
	GROUP_SYSTEM    = 0x00080000
	GROUP_FLSHPROG  = 0x00100000

	// 错误状态 - 使用负数表示法
	STATUS_ERROR = -2147483648

	// 通常状态
	STATUS_USUALLY_IDLE        = GROUP_USUALLY | 0x0001
	STATUS_USUALLY_PRINTING    = GROUP_USUALLY | 0x0002
	STATUS_USUALLY_STANDSTILL  = GROUP_USUALLY | 0x0004
	STATUS_USUALLY_PAPER_END   = GROUP_USUALLY | 0x0008
	STATUS_USUALLY_RIBBON_END  = GROUP_USUALLY | 0x0010
	STATUS_USUALLY_COOLING     = GROUP_USUALLY | 0x0020
	STATUS_USUALLY_MOTCOOLING  = GROUP_USUALLY | 0x0040

	// 设置状态
	STATUS_SETTING_COVER_OPEN   = GROUP_SETTING | 0x0001
	STATUS_SETTING_PAPER_JAM    = GROUP_SETTING | 0x0002
	STATUS_SETTING_RIBBON_ERR   = GROUP_SETTING | 0x0004
	STATUS_SETTING_PAPER_ERR    = GROUP_SETTING | 0x0008
	STATUS_SETTING_DATA_ERR     = GROUP_SETTING | 0x0010
	STATUS_SETTING_SCRAPBOX_ERR = GROUP_SETTING | 0x0020

	// 硬件错误状态
	STATUS_HARDWARE_ERR01 = GROUP_HARDWARE | 0x0001
	STATUS_HARDWARE_ERR02 = GROUP_HARDWARE | 0x0002
	STATUS_HARDWARE_ERR03 = GROUP_HARDWARE | 0x0004
	STATUS_HARDWARE_ERR04 = GROUP_HARDWARE | 0x0008
	STATUS_HARDWARE_ERR05 = GROUP_HARDWARE | 0x0010
	STATUS_HARDWARE_ERR06 = GROUP_HARDWARE | 0x0020
	STATUS_HARDWARE_ERR07 = GROUP_HARDWARE | 0x0040
	STATUS_HARDWARE_ERR08 = GROUP_HARDWARE | 0x0080
	STATUS_HARDWARE_ERR09 = GROUP_HARDWARE | 0x0100
	STATUS_HARDWARE_ERR10 = GROUP_HARDWARE | 0x0200

	// 系统错误状态
	STATUS_SYSTEM_ERR01 = GROUP_SYSTEM | 0x0001
)

// DLL函数声明
var (
	cyStat64DLL        *syscall.DLL
	portInitialize     *syscall.Proc
	getStatus          *syscall.Proc
	getFirmwVersion    *syscall.Proc
	getSerialNo        *syscall.Proc
	getSensorInfo      *syscall.Proc
	getCounterL        *syscall.Proc
	getCounterA        *syscall.Proc
	getCounterB        *syscall.Proc
	getMediaCounter    *syscall.Proc
	getMediaLotNo      *syscall.Proc
	getResolutionH     *syscall.Proc
	getResolutionV     *syscall.Proc
	getFreeBuffer      *syscall.Proc
	getPQTY            *syscall.Proc
)

// 打印机信息结构体
type PrinterInfo struct {
	Status          string `json:"status"`
	StatusCode      int32  `json:"status_code"` // 改回 int32
	FirmwareVersion string `json:"firmware_version"`
	SerialNumber    string `json:"serial_number"`
	SensorInfo      string `json:"sensor_info"`
	CounterL        int32  `json:"counter_l"`
	CounterA        int32  `json:"counter_a"`
	CounterB        int32  `json:"counter_b"`
	MediaCounter    int32  `json:"media_counter"`
	MediaLotNo      string `json:"media_lot_no"`
	ResolutionH     int32  `json:"resolution_h"`
	ResolutionV     int32  `json:"resolution_v"`
	FreeBuffer      int32  `json:"free_buffer"`
	PQTY            int32  `json:"pqty"`
	Error           string `json:"error,omitempty"`
}

// 全局变量存储端口号
var printerPort int32 = -1

// 初始化DLL
func initDLL() error {
	var err error
	cyStat64DLL, err = syscall.LoadDLL("CyStat64.dll")
	if err != nil {
		return fmt.Errorf("加载DLL失败: %v", err)
	}

	// 获取函数地址
	portInitialize, err = cyStat64DLL.FindProc("PortInitialize")
	if err != nil {
		return fmt.Errorf("找不到PortInitialize函数: %v", err)
	}

	getStatus, err = cyStat64DLL.FindProc("GetStatus")
	if err != nil {
		return fmt.Errorf("找不到GetStatus函数: %v", err)
	}

	getFirmwVersion, err = cyStat64DLL.FindProc("GetFirmwVersion")
	if err != nil {
		return fmt.Errorf("找不到GetFirmwVersion函数: %v", err)
	}

	getSerialNo, err = cyStat64DLL.FindProc("GetSerialNo")
	if err != nil {
		return fmt.Errorf("找不到GetSerialNo函数: %v", err)
	}

	getSensorInfo, err = cyStat64DLL.FindProc("GetSensorInfo")
	if err != nil {
		return fmt.Errorf("找不到GetSensorInfo函数: %v", err)
	}

	getCounterL, err = cyStat64DLL.FindProc("GetCounterL")
	if err != nil {
		return fmt.Errorf("找不到GetCounterL函数: %v", err)
	}

	getCounterA, err = cyStat64DLL.FindProc("GetCounterA")
	if err != nil {
		return fmt.Errorf("找不到GetCounterA函数: %v", err)
	}

	getCounterB, err = cyStat64DLL.FindProc("GetCounterB")
	if err != nil {
		return fmt.Errorf("找不到GetCounterB函数: %v", err)
	}

	getMediaCounter, err = cyStat64DLL.FindProc("GetMediaCounter")
	if err != nil {
		return fmt.Errorf("找不到GetMediaCounter函数: %v", err)
	}

	getMediaLotNo, err = cyStat64DLL.FindProc("GetMediaLotNo")
	if err != nil {
		return fmt.Errorf("找不到GetMediaLotNo函数: %v", err)
	}

	getResolutionH, err = cyStat64DLL.FindProc("GetResolutionH")
	if err != nil {
		return fmt.Errorf("找不到GetResolutionH函数: %v", err)
	}

	getResolutionV, err = cyStat64DLL.FindProc("GetResolutionV")
	if err != nil {
		return fmt.Errorf("找不到GetResolutionV函数: %v", err)
	}

	getFreeBuffer, err = cyStat64DLL.FindProc("GetFreeBuffer")
	if err != nil {
		return fmt.Errorf("找不到GetFreeBuffer函数: %v", err)
	}

	getPQTY, err = cyStat64DLL.FindProc("GetPQTY")
	if err != nil {
		return fmt.Errorf("找不到GetPQTY函数: %v", err)
	}

	return nil
}

// 初始化打印机端口
func initPrinter() error {
	// 将"USB001"转换为UTF-16
	portName, err := syscall.UTF16PtrFromString("USB001")
	if err != nil {
		return fmt.Errorf("转换端口名失败: %v", err)
	}

	// 调用PortInitialize
	ret, _, _ := portInitialize.Call(uintptr(unsafe.Pointer(portName)))
	printerPort = int32(ret)

	if printerPort < 0 {
		return fmt.Errorf("初始化打印机端口失败，返回值: %d", printerPort)
	}

	log.Printf("打印机端口初始化成功，端口号: %d", printerPort)
	return nil
}

// 获取字符串信息的辅助函数
func getStringInfo(proc *syscall.Proc, portNum int32) (string, error) {
	buffer := make([]byte, 256)
	ret, _, _ := proc.Call(uintptr(portNum), uintptr(unsafe.Pointer(&buffer[0])))

	if int32(ret) < 0 {
		return "", fmt.Errorf("调用失败，返回值: %d", int32(ret))
	}

	// 找到字符串结束位置
	length := int32(ret)
	if length > 255 {
		length = 255
	}

	return string(buffer[:length]), nil
}

// 获取整数信息的辅助函数
func getIntInfo(proc *syscall.Proc, portNum int32) (int32, error) {
	ret, _, _ := proc.Call(uintptr(portNum))
	result := int32(ret)

	if result < 0 {
		return 0, fmt.Errorf("调用失败，返回值: %d", result)
	}

	return result, nil
}

// 解析状态码为可读字符串
func parseStatus(statusCode int32) string {
	if statusCode == STATUS_ERROR {
		return "ERROR"
	}

	if statusCode&GROUP_USUALLY != 0 {
		switch statusCode {
		case STATUS_USUALLY_IDLE:
			return "空闲"
		case STATUS_USUALLY_PRINTING:
			return "打印中"
		case STATUS_USUALLY_STANDSTILL:
			return "待机"
		case STATUS_USUALLY_PAPER_END:
			return "纸张用完"
		case STATUS_USUALLY_RIBBON_END:
			return "色带用完"
		case STATUS_USUALLY_COOLING:
			return "打印头冷却中"
		case STATUS_USUALLY_MOTCOOLING:
			return "电机冷却中"
		}
	} else if statusCode&GROUP_SETTING != 0 {
		switch statusCode {
		case STATUS_SETTING_COVER_OPEN:
			return "盖子打开"
		case STATUS_SETTING_PAPER_JAM:
			return "卡纸"
		case STATUS_SETTING_RIBBON_ERR:
			return "色带错误"
		case STATUS_SETTING_PAPER_ERR:
			return "纸张定义错误"
		case STATUS_SETTING_DATA_ERR:
			return "数据错误"
		case STATUS_SETTING_SCRAPBOX_ERR:
			return "废料盒错误"
		}
	} else if statusCode&GROUP_HARDWARE != 0 {
		switch statusCode {
		case STATUS_HARDWARE_ERR01:
			return "打印头电压错误"
		case STATUS_HARDWARE_ERR02:
			return "打印头位置错误"
		case STATUS_HARDWARE_ERR03:
			return "风扇停止错误"
		case STATUS_HARDWARE_ERR04:
			return "切刀错误"
		case STATUS_HARDWARE_ERR05:
			return "压纸轮错误"
		case STATUS_HARDWARE_ERR06:
			return "打印头温度异常"
		case STATUS_HARDWARE_ERR07:
			return "介质温度异常"
		case STATUS_HARDWARE_ERR08:
			return "色带张力错误"
		case STATUS_HARDWARE_ERR09:
			return "RFID模块错误"
		case STATUS_HARDWARE_ERR10:
			return "电机温度异常"
		}
	} else if statusCode&GROUP_SYSTEM != 0 {
		return "系统错误"
	} else if statusCode&GROUP_FLSHPROG != 0 {
		return "固件更新模式"
	}

	return fmt.Sprintf("未知状态: 0x%X", statusCode)
}

// 获取完整的打印机信息
func getPrinterInfo() (*PrinterInfo, error) {
	if printerPort < 0 {
		return nil, fmt.Errorf("打印机未初始化")
	}

	info := &PrinterInfo{}

	// 获取状态
	statusCode, err := getIntInfo(getStatus, printerPort)
	if err != nil {
		info.Error = fmt.Sprintf("获取状态失败: %v", err)
	} else {
		info.StatusCode = statusCode
		info.Status = parseStatus(statusCode)
	}

	// 获取固件版本
	if firmwareVersion, err := getStringInfo(getFirmwVersion, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取固件版本失败: %v; ", err)
	} else {
		info.FirmwareVersion = firmwareVersion
	}

	// 获取序列号
	if serialNumber, err := getStringInfo(getSerialNo, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取序列号失败: %v; ", err)
	} else {
		info.SerialNumber = serialNumber
	}

	// 获取传感器信息
	if sensorInfo, err := getStringInfo(getSensorInfo, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取传感器信息失败: %v; ", err)
	} else {
		info.SensorInfo = sensorInfo
	}

	// 获取介质批号
	if mediaLotNo, err := getStringInfo(getMediaLotNo, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取介质批号失败: %v; ", err)
	} else {
		info.MediaLotNo = mediaLotNo
	}

	// 获取各种计数器
	if counterL, err := getIntInfo(getCounterL, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取生命周期计数器失败: %v; ", err)
	} else {
		info.CounterL = counterL
	}

	if counterA, err := getIntInfo(getCounterA, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取计数器A失败: %v; ", err)
	} else {
		info.CounterA = counterA
	}

	if counterB, err := getIntInfo(getCounterB, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取计数器B失败: %v; ", err)
	} else {
		info.CounterB = counterB
	}

	if mediaCounter, err := getIntInfo(getMediaCounter, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取介质计数器失败: %v; ", err)
	} else {
		info.MediaCounter = mediaCounter
	}

	// 获取分辨率
	if resolutionH, err := getIntInfo(getResolutionH, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取水平分辨率失败: %v; ", err)
	} else {
		info.ResolutionH = resolutionH
	}

	if resolutionV, err := getIntInfo(getResolutionV, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取垂直分辨率失败: %v; ", err)
	} else {
		info.ResolutionV = resolutionV
	}

	// 获取缓冲区信息
	if freeBuffer, err := getIntInfo(getFreeBuffer, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取空闲缓冲区失败: %v; ", err)
	} else {
		info.FreeBuffer = freeBuffer
	}

	if pqty, err := getIntInfo(getPQTY, printerPort); err != nil {
		info.Error += fmt.Sprintf("获取PQTY失败: %v; ", err)
	} else {
		info.PQTY = pqty
	}

	return info, nil
}

// 标准API响应结构体
type APIResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

// 打印机状态响应数据结构体
type PrinterStatusData struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

// HTTP处理函数
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	info, err := getPrinterInfo()
	if err != nil {
		// 错误响应
		errorResponse := APIResponse{
			Data: nil,
			Code: 500,
			Msg:  err.Error(),
		}
		jsonData, _ := json.MarshalIndent(errorResponse, "", "  ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonData)
		return
	}

	// 构建响应数据 - 根据不同状态返回对应的status值
	var statusValue int32
	if info.StatusCode == STATUS_USUALLY_IDLE {
		statusValue = 0 // 空闲时返回0
	} else if info.StatusCode == STATUS_USUALLY_PRINTING {
		statusValue = 2 // 打印中时返回2
	} else {
		statusValue = info.StatusCode // 其他状态返回原始状态码
	}

	responseData := PrinterStatusData{
		Message: info.Status,
		Status:  statusValue,
	}

	// 构建标准API响应
	apiResponse := APIResponse{
		Data: responseData,
		Code: 200,
		Msg:  responseData.Message, // 直接引用data.message
	}

	jsonData, err := json.MarshalIndent(apiResponse, "", "  ")
	if err != nil {
		errorResponse := APIResponse{
			Data: nil,
			Code: 500,
			Msg:  "JSON编码失败: " + err.Error(),
		}
		errorData, _ := json.MarshalIndent(errorResponse, "", "  ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorData)
		return
	}

	w.Write(jsonData)
}

// 健康检查处理函数
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "printer-status-api"}`))
}

// 主页处理函数
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>打印机状态API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px; margin: 5px; }
        .button:hover { background-color: #45a049; }
        pre { background-color: #f4f4f4; padding: 15px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1>打印机状态API服务</h1>
        <p>欢迎使用打印机状态监控API服务</p>
        
        <h2>可用接口:</h2>
        <ul>
            <li><a href="/status" class="button">获取打印机状态</a> - GET /status</li>
            <li><a href="/health" class="button">健康检查</a> - GET /health</li>
        </ul>
        
        <h2>使用示例:</h2>
        <pre>
# 获取打印机状态
curl http://localhost:8080/status

# 健康检查
curl http://localhost:8080/health
        </pre>
        
        <h2>返回数据格式:</h2>
        <pre>
{
  "status": "空闲",
  "status_code": 65537,
  "firmware_version": "1.0.0",
  "serial_number": "ABC123456",
  "sensor_info": "传感器信息",
  "counter_l": 1000,
  "counter_a": 500,
  "counter_b": 300,
  "media_counter": 200,
  "media_lot_no": "LOT001",
  "resolution_h": 300,
  "resolution_v": 300,
  "free_buffer": 1024,
  "pqty": 100
}
        </pre>
    </div>
</body>
</html>
`
	w.Write([]byte(html))
}

func main() {
	log.Println("正在启动打印机状态API服务...")

	// 初始化DLL
	if err := initDLL(); err != nil {
		log.Fatalf("初始化DLL失败: %v", err)
	}
	log.Println("DLL初始化成功")

	// 初始化打印机
	if err := initPrinter(); err != nil {
		log.Printf("警告: 初始化打印机失败: %v", err)
		log.Println("服务将继续运行，但可能无法获取打印机状态")
	}

	// 设置路由
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/health", healthHandler)

	// 启动服务器
	log.Println("服务器正在监听端口 8080...")
	log.Println("访问 http://localhost:8080 查看API文档")
	log.Println("访问 http://localhost:8080/status 获取打印机状态")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}