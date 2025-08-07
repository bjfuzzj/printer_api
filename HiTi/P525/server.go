package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// 全局变量
var (
	isWindows bool
	hitiDLL   *syscall.DLL
	checkPrinterStatusW *syscall.Proc
	doCommandW          *syscall.Proc
	enumUsbPrintsW      *syscall.Proc
	getDeviceInfoW      *syscall.Proc
)

// 状态码映射
var statusDescriptions = map[uint32]string{
	// 通用状态码
	0:          "正常",
	0x00080000: "打印机忙",
	0x00000080: "打印机离线或关机",
	0x00000400: "打印中",
	0x00000005: "处理打印数据中",
	0x00000006: "发送数据到打印机中",
	0x00050001: "盖子或门打开",
	0x00050101: "盖子或门打开",
	0x00008000: "缺纸或进纸错误",
	0x00030000: "卡纸",
	0x000100FE: "纸张类型不匹配",
	0x00008010: "纸盘不匹配",
	0x00008008: "纸盘缺失",
	0x00080004: "色带缺失",
	0x00080103: "色带用尽",
	0x00080200: "色带类型不匹配",
	0x000802FE: "色带错误",
	0x00030001: "SRAM错误",
	0x00030101: "SDRAM错误",
	0x00030201: "ADC错误",
	0x00030301: "NVRAM读写错误",
	0x00030302: "固件校验和错误",
	0x00030402: "DSP代码校验和错误",
	0x000304FE: "加热参数表不兼容",
	0x00030501: "Cam Platen错误",
	0x00030601: "ADF错误",
	0x0000001F: "发送数据失败",
	0x0000002F: "接收数据失败",

	// 特定型号状态码 (P720L/P728L/P520L/P750L)
	0x00000100: "盖子打开",
	0x00000101: "盖子打开失败",
	0x00000200: "IC芯片缺失",
	0x00000201: "色带缺失",
	0x00000202: "色带类型不匹配01",
	0x00000203: "安全校验失败",
	0x00000204: "色带类型不匹配02",
	0x00000205: "色带类型不匹配03",
	0x00000300: "色带用尽01",
	0x00000301: "色带用尽02",
	0x00000302: "打印失败",
	0x00000401: "缺纸02",
	0x00000402: "纸张未就绪",
	0x00000500: "卡纸01",
	0x00000501: "卡纸02",
	0x00000502: "卡纸03",
	0x00000503: "卡纸04",
	0x00000504: "卡纸05",
	0x00000600: "纸张不匹配",
	0x00000700: "Cam错误01",
	0x00000800: "Cam错误02",
	0x00000900: "NVRAM错误",
	0x00000A00: "IC芯片错误",
	0x00000C00: "ADC错误",
	0x00000D00: "固件校验错误",
	0x00000F00: "切刀错误",
}

// 优先级检查顺序
var priorityOrder = []uint32{
	0x00050001, 0x00050101, // 盖子/门开
	0x00030000,             // 卡纸
	0x00008000,             // 缺纸
	0x00000100, 0x00000101, // 特定型号盖子问题
	0x00000500, 0x00000501, // 特定型号卡纸
}

// API响应结构体
type APIResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

// 打印机状态响应数据
type PrinterStatusData struct {
	Message string `json:"message"`
	Status  uint32 `json:"status"`
}

// 错误响应数据
type ErrorData struct {
	Message     string `json:"message"`
	SystemError uint32 `json:"system_error,omitempty"`
}

// 打印机信息结构体
type PrinterInfo struct {
	Name    string `json:"name"`
	ModelNo uint8  `json:"modelNo"`
	IndexNo uint8  `json:"indexNo"`
}

// 枚举打印机响应
type EnumPrintersResponse struct {
	Printers []PrinterInfo `json:"printers"`
	Count    int           `json:"count"`
	Note     string        `json:"note,omitempty"`
}

// 设备信息响应
type DeviceInfoResponse struct {
	Type           string      `json:"type,omitempty"`
	Data           string      `json:"data,omitempty"`
	RibbonType     interface{} `json:"ribbon_type,omitempty"`
	RemainingCount interface{} `json:"remaining_count,omitempty"`
	Count          interface{} `json:"count,omitempty"`
	Info           string      `json:"info,omitempty"`
	Note           string      `json:"note,omitempty"`
}

// 命令请求结构体
type CommandRequest struct {
	Printer string `json:"printer"`
	Command int    `json:"command"`
}

// 命令响应结构体
type CommandResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Note    string `json:"note,omitempty"`
}

// 初始化函数
func init() {
	isWindows = runtime.GOOS == "windows"
	
	wd, _ := os.Getwd()
	dllPath := filepath.Join(wd, "HTRTApiW.dll")
	
	fmt.Printf("DLL 路径: %s\n", dllPath)
	fmt.Printf("DLL 文件存在: %t\n", fileExists(dllPath))
	
	if isWindows {
		if err := loadDLL(dllPath); err != nil {
			fmt.Printf("DLL 加载失败: %v\n", err)
			isWindows = false
		} else {
			fmt.Println("DLL 加载成功")
		}
	} else {
		fmt.Println("在 macOS 上运行，使用模拟模式")
	}
}

// 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 加载DLL
func loadDLL(dllPath string) error {
	var err error
	hitiDLL, err = syscall.LoadDLL(dllPath)
	if err != nil {
		return err
	}
	
	// 加载函数
	checkPrinterStatusW, err = hitiDLL.FindProc("HITI_CheckPrinterStatusW")
	if err != nil {
		return err
	}
	
	doCommandW, err = hitiDLL.FindProc("HITI_DoCommandW")
	if err != nil {
		return err
	}
	
	enumUsbPrintsW, err = hitiDLL.FindProc("HITI_EnumUsbPrintsW")
	if err != nil {
		return err
	}
	
	getDeviceInfoW, err = hitiDLL.FindProc("HITI_GetDeviceInfoW")
	if err != nil {
		return err
	}
	
	return nil
}

// 获取状态描述
func getStatusDescription(status uint32) string {
	// 1. 精确匹配优先
	if desc, exists := statusDescriptions[status]; exists {
		return desc
	}
	
	// 2. 按优先级检查关键错误
	for _, code := range priorityOrder {
		if status&code != 0 {
			if desc, exists := statusDescriptions[code]; exists {
				return desc
			}
			return fmt.Sprintf("未知组合状态: 0x%X", status)
		}
	}
	
	// 3. 通用按位检查
	for code, desc := range statusDescriptions {
		if code != 0 && (status&code) != 0 {
			return desc
		}
	}
	
	return fmt.Sprintf("未知状态: 0x%X", status)
}

// 设置CORS头
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 状态检查处理函数
func statusHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	printer := r.URL.Query().Get("printer")
	if printer == "" {
		errorResponse := APIResponse{
			Data: ErrorData{Message: "缺少打印机名称参数"},
			Code: 400,
			Msg:  "缺少打印机名称参数",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	if isWindows {
		// Windows实现
		var status uint32
		printerPtr, _ := syscall.UTF16PtrFromString(printer)
		
		ret, _, _ := checkPrinterStatusW.Call(
			uintptr(unsafe.Pointer(printerPtr)),
			uintptr(unsafe.Pointer(&status)),
		)
		
		if ret == 0 {
			fmt.Printf("[DEBUG] Printer Status - 十进制: %d, 十六进制: 0x%X\n", status, status)
			
			responseData := PrinterStatusData{
				Message: getStatusDescription(status),
				Status:  status,
			}
			
			apiResponse := APIResponse{
				Data: responseData,
				Code: 200,
				Msg:  responseData.Message,
			}
			
			json.NewEncoder(w).Encode(apiResponse)
		} else {
			// 获取Windows系统错误码
			errorCode := uint32(ret)
			fmt.Printf("[ERROR] 检查状态失败 - 错误码: %d\n", errorCode)
			
			errorData := ErrorData{
				Message:     fmt.Sprintf("系统错误: %d", errorCode),
				SystemError: errorCode,
			}
			
			errorResponse := APIResponse{
				Data: errorData,
				Code: 500,
				Msg:  errorData.Message,
			}
			
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
		}
	} else {
		// macOS模拟模式
		responseData := PrinterStatusData{
			Message: "模拟状态 - 正常 (macOS)",
			Status:  0,
		}
		
		apiResponse := APIResponse{
			Data: responseData,
			Code: 200,
			Msg:  "模拟状态 - 正常 (macOS)",
		}
		
		json.NewEncoder(w).Encode(apiResponse)
	}
}

// 枚举打印机处理函数
func enumPrintersHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	
	if isWindows {
		// Windows实现 - 这里需要根据实际的HITI_USB_PRINTER结构体来实现
		// 由于结构体定义复杂，这里提供一个简化的模拟实现
		response := EnumPrintersResponse{
			Printers: []PrinterInfo{
				{Name: "HITI打印机1", ModelNo: 1, IndexNo: 1},
			},
			Count: 1,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		// macOS模拟数据
		response := EnumPrintersResponse{
			Printers: []PrinterInfo{
				{Name: "模拟打印机1", ModelNo: 1, IndexNo: 1},
				{Name: "模拟打印机2", ModelNo: 2, IndexNo: 2},
			},
			Count: 2,
			Note:  "这是 macOS 上的模拟数据",
		}
		json.NewEncoder(w).Encode(response)
	}
}

// 设备信息处理函数
func deviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	
	printer := r.URL.Query().Get("printer")
	infoType := r.URL.Query().Get("infoType")
	
	if printer == "" || infoType == "" {
		errorResponse := map[string]string{"error": "缺少必要参数：printer或infoType"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	if isWindows {
		infoTypeInt, err := strconv.Atoi(infoType)
		if err != nil {
			errorResponse := map[string]string{"error": "无效的信息类型，必须为数字"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		
		// 设备信息类型映射
		devInfoTypes := map[int]string{
			1: "MFG_SERIAL",
			2: "MODEL_NAME",
			3: "FIRMWARE_VERSION",
			4: "RIBBON_INFO",
			5: "PRINT_COUNT",
			6: "CUTTER_COUNT",
		}
		
		// 这里需要实际调用HITI_GetDeviceInfoW，简化实现
		response := DeviceInfoResponse{
			Type: devInfoTypes[infoTypeInt],
			Data: fmt.Sprintf("模拟数据 - 打印机: %s, 类型: %d", printer, infoTypeInt),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := DeviceInfoResponse{
			Info: fmt.Sprintf("模拟信息 - 打印机: %s, 类型: %s", printer, infoType),
			Note: "macOS 模拟数据",
		}
		json.NewEncoder(w).Encode(response)
	}
}

// 命令执行处理函数
func commandHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	if r.Header.Get("Content-Type") != "application/json" {
		errorResponse := map[string]string{"error": "请求必须是JSON格式"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse := map[string]string{"error": "JSON解析失败"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	if req.Printer == "" {
		errorResponse := map[string]string{"error": "缺少必要参数：printer或command"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	if isWindows {
		// Windows实现
		printerPtr, _ := syscall.UTF16PtrFromString(req.Printer)
		ret, _, _ := doCommandW.Call(
			uintptr(unsafe.Pointer(printerPtr)),
			uintptr(req.Command),
		)
		
		if ret == 0 {
			response := CommandResponse{Success: true}
			json.NewEncoder(w).Encode(response)
		} else {
			response := CommandResponse{
				Success: false,
				Error:   fmt.Sprintf("执行命令失败，错误码: %d", ret),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
	} else {
		// macOS模拟数据
		response := CommandResponse{
			Success: true,
			Note:    fmt.Sprintf("模拟执行命令 - 打印机: %s, 命令: %d (macOS 模拟模式)", req.Printer, req.Command),
		}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("正在启动 Go HTTP 服务器...")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("Go 版本: %s\n", runtime.Version())
	wd, _ := os.Getwd()
	fmt.Printf("工作目录: %s\n", wd)
	fmt.Printf("Windows 模式: %t\n", isWindows)
	fmt.Println(strings.Repeat("=", 50))
	
	// 设置路由
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/enumPrinters", enumPrintersHandler)
	http.HandleFunc("/deviceInfo", deviceInfoHandler)
	http.HandleFunc("/command", commandHandler)
	
	// 启动服务器
	fmt.Println("服务器启动在端口 8080")
	fmt.Println("访问 http://localhost:8080/status?printer=打印机名称 获取状态")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("启动 HTTP 服务器时发生错误: %v", err)
	}
}