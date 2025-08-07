from flask import Flask, request, jsonify
import os
import platform

app = Flask(__name__)

# 检测操作系统
is_windows = platform.system() == 'Windows'

# DLL 路径设置为变量
# DLL_DIR = r"D:\camera\HiTiRollTypesSDK_Windows_1.10.14\dllx64"
# DLL 路径设置为当前目录
DLL_DIR = os.getcwd()  # 使用当前工作目录
DLL_NAME = "HTRTApiW.dll"
dll_path = os.path.join(DLL_DIR, DLL_NAME)

# 检查文件是否存在
print(f"DLL 路径: {dll_path}")
print(f"DLL 文件存在: {os.path.exists(dll_path)}")

# 如果是 Windows 系统，尝试加载 DLL
if is_windows:
    import ctypes
    # Add c_wchar and other required types to the import
    from ctypes import (
        c_char_p, c_uint, c_ulong, c_wchar, 
        Structure, POINTER, byref, WINFUNCTYPE
    )
    
    try:
        # Load the UNICODE version DLL explicitly
        hiti_dll = ctypes.WinDLL(dll_path)
        print("DLL 加载成功")
        
        # Define function prototypes according to SDK documentation
        # HITI_CheckPrinterStatus (section 3.2.2)
        hiti_dll.HITI_CheckPrinterStatusW.argtypes = [
            ctypes.c_wchar_p,    # LPTSTR szPrinter
            POINTER(c_uint)      # DWORD *lpdwStatus
        ]
        hiti_dll.HITI_CheckPrinterStatusW.restype = c_uint
        
        # HITI_DoCommand (section 3.2.3)
        hiti_dll.HITI_DoCommandW.argtypes = [
            ctypes.c_wchar_p,    # LPTSTR szPrinter
            c_uint              # DWORD dwCommand
        ]
        hiti_dll.HITI_DoCommandW.restype = c_uint
        
    except Exception as e:
        print(f"DLL 加载失败: {e}")
        is_windows = False
else:
    print("在 macOS 上运行，使用模拟模式")

# 状态码映射
# 状态码映射
status_descriptions = {
    # 通用状态码
    0: "正常",  # 0 (0x0) b'0'  # 添加缺失的正常状态定义
    0x00080000: "打印机忙",          # 524288 (0x80000) b'10000000000000000000'
    0x00000080: "打印机离线或关机",    # 128 (0x80)     b'000000000000000010000000'
    0x00000400: "打印中",            # 1024 (0x400)   b'000000000000010000000000'
    0x00000005: "处理打印数据中",      # 5 (0x5)       b'000000000000000000000101'
    0x00000006: "发送数据到打印机中",   # 6 (0x6)       b'000000000000000000000110'
    0x00050001: "盖子或门打开",        # 327681 (0x50001) b'010100000000000000000001'
    0x00050101: "盖子或门打开",        # 327937 (0x50101) b'010100000001000000000001'
    0x00008000: "缺纸或进纸错误",       # 32768 (0x8000)  b'000000001000000000000000'
    0x00030000: "卡纸",              # 196608 (0x30000) b'001100000000000000000000'
    0x000100FE: "纸张类型不匹配",       # 65790 (0x100FE) b'00010000000011111110'
    0x00008010: "纸盘不匹配",          # 32784 (0x8010)  b'000000001000000000010000'
    0x00008008: "纸盘缺失",            # 32776 (0x8008)  b'000000001000000000001000'
    0x00080004: "色带缺失",            # 524292 (0x80004) b'10000000000000000100'
    0x00080103: "色带用尽",            # 524547 (0x80103) b'10000000000100000011'
    0x00080200: "色带类型不匹配",        # 524800 (0x80200) b'10000000001000000000'
    0x000802FE: "色带错误",            # 525054 (0x802FE) b'10000000001011111110'
    0x00030001: "SRAM错误",           # 196609 (0x30001) b'00110000000000000001'
    0x00030101: "SDRAM错误",          # 196865 (0x30101) b'00110000000100000001'
    0x00030201: "ADC错误",            # 196865 (0x30201) b'00110000001000000001'
    0x00030301: "NVRAM读写错误",       # 197377 (0x30301) b'00110000001100000001'
    0x00030302: "固件校验和错误",        # 197378 (0x30302) b'00110000001100000010'
    0x00030402: "DSP代码校验和错误",     # 197634 (0x30402) b'00110000010000000010'
    0x000304FE: "加热参数表不兼容",       # 197886 (0x304FE) b'00110000010011111110'
    0x00030501: "Cam Platen错误",     # 197889 (0x30501) b'00110000010100000001'
    0x00030601: "ADF错误",            # 198145 (0x30601) b'00110000011000000001'
    0x0000001F: "发送数据失败",         # 31 (0x1F)      b'00000000000000011111'
    0x0000002F: "接收数据失败",         # 47 (0x2F)      b'00000000000000101111'

    # 特定型号状态码 (P720L/P728L/P520L/P750L)
    0x00000100: "盖子打开",  # 256 (0x100)
    0x00000101: "盖子打开失败",  # 257 (0x101)
    0x00000200: "IC芯片缺失",  # 512 (0x200)
    0x00000201: "色带缺失",  # 513 (0x201)
    0x00000202: "色带类型不匹配01",  # 514 (0x202)
    0x00000203: "安全校验失败",  # 515 (0x203)
    0x00000204: "色带类型不匹配02",  # 516 (0x204)
    0x00000205: "色带类型不匹配03",  # 517 (0x205)
    0x00000300: "色带用尽01",  # 768 (0x300)
    0x00000301: "色带用尽02",  # 769 (0x301)
    0x00000302: "打印失败",  # 770 (0x302)
    0x00000400: "缺纸01",  # 1024 (0x400)
    0x00000401: "缺纸02",  # 1025 (0x401)
    0x00000402: "纸张未就绪",  # 1026 (0x402)
    0x00000500: "卡纸01",  # 1280 (0x500)
    0x00000501: "卡纸02",  # 1281 (0x501)
    0x00000502: "卡纸03",  # 1282 (0x502)
    0x00000503: "卡纸04",  # 1283 (0x503)
    0x00000504: "卡纸05",  # 1284 (0x504)
    0x00000600: "纸张不匹配",  # 1536 (0x600)
    0x00000700: "Cam错误01",  # 1792 (0x700)
    0x00000800: "Cam错误02",  # 2048 (0x800)
    0x00000900: "NVRAM错误",  # 2304 (0x900)
    0x00000A00: "IC芯片错误",  # 2560 (0xA00)
    0x00000C00: "ADC错误",  # 3072 (0xC00)
    0x00000D00: "固件校验错误",  # 3328 (0xD00)
    0x00000F00: "切刀错误"  # 3840 (0xF00)

    
}

def get_status_description(status):
    # 1. 精确匹配优先
    if status in status_descriptions:
        return status_descriptions[status]
    
    # 2. 按优先级检查关键错误
    priority_order = [
        # 通用错误优先
        0x00050001, 0x00050101,  # 盖子/门开
        0x00030000,              # 卡纸
        0x00008000,              # 缺纸
        
        # 特定型号关键错误
        0x00000100, 0x00000101,  # 特定型号盖子问题
        0x00000500, 0x00000501,  # 特定型号卡纸
        
        # 其他关键错误...
    ]
    
    for code in priority_order:
        if status & code:
            return status_descriptions.get(code, f"未知组合状态: 0x{status:X}")
    
    # 3. 通用按位检查
    for code, desc in status_descriptions.items():
        if code != 0 and (status & code):
            return desc
            
    return f"未知状态: 0x{status:X}"

@app.route('/status', methods=['GET'])
def check_printer_status():
    printer = request.args.get('printer')
    
    if not printer:
        return jsonify({
            "data": {"message": "缺少打印机名称参数"},
            "code": 400,
            "msg": "缺少打印机名称参数"  # 保持内外message一致
        }), 400
    
    if is_windows:
        status = c_uint(0)
        result = hiti_dll.HITI_CheckPrinterStatusW(
            printer,
            byref(status)
        )
        
        if result == 0:
            print(f"[DEBUG] Printer Status - 十进制: {status.value}, 十六进制: 0x{status.value:X}")
            # 修改data结构
            response_data = {
                "message": get_status_description(status.value),
                "status": status.value
            }
            
            return jsonify({
                "data": response_data,
                "code": 200,
                "msg": response_data["message"]  # 直接引用data.message
            })
        else:
            error_code = ctypes.GetLastError()
            print(f"[ERROR] 检查状态失败 - 错误码: {error_code} ({ctypes.FormatError(error_code)})")
            # 统一错误消息结构
            error_data = {
                "message": ctypes.FormatError(error_code),
                "system_error": error_code
            }
            return jsonify({
                "data": error_data,
                "code": 500,
                "msg": error_data["message"]
            }), 500
    else:
        # 修改模拟模式数据结构
        return jsonify({
            "data": {
                "message": "模拟状态 - 正常 (macOS)",
                "status": 0
            },
            "code": 200,
            "msg": "模拟状态 - 正常 (macOS)"
        })

@app.route('/enumPrinters', methods=['GET'])
def enum_printers():
    if is_windows:
        MAX_PRINTERS = 10
        printers = (HITI_USB_PRINTER * MAX_PRINTERS)()
        cb_needed = c_ulong(0)
        returned = c_ulong(0)
        
        # 使用 Unicode 版本函数调用
        result = hiti_dll.HITI_EnumUsbPrintsW(
            printers,
            ctypes.sizeof(HITI_USB_PRINTER) * MAX_PRINTERS,
            byref(cb_needed),
            byref(returned)
        )
        
        if result == 0:  # 成功
            printers_list = []
            for i in range(returned.value):
                # 直接获取宽字符串
                name = printers[i].PrinterName.split('\x00')[0]
                printers_list.append({
                    "name": name,
                    "modelNo": printers[i].bModelNo,
                    "indexNo": printers[i].bIndexNo
                })
            
            return jsonify({
                "printers": printers_list,
                "count": returned.value
            })
        else:
            return jsonify({"error": f"枚举失败，错误码: {result}"}), 500
    else:
        # macOS模拟数据保持不变
        return jsonify({
            "printers": [
                {"name": "模拟打印机1", "modelNo": 1, "indexNo": 1},
                {"name": "模拟打印机2", "modelNo": 2, "indexNo": 2}
            ],
            "count": 2,
            "note": "这是 macOS 上的模拟数据"
        })

@app.route('/deviceInfo', methods=['GET'])
def get_device_info():
    printer = request.args.get('printer')
    info_type = request.args.get('infoType')
    
    if not printer or not info_type:
        return jsonify({"error": "缺少必要参数：printer或infoType"}), 400

    if is_windows:
        try:
            info_type = int(info_type)
        except ValueError:
            return jsonify({"error": "无效的信息类型，必须为数字"}), 400

        # 定义设备信息类型常量
        HITI_DEVINFO = {
            1: "MFG_SERIAL",
            2: "MODEL_NAME", 
            3: "FIRMWARE_VERSION",
            4: "RIBBON_INFO",
            5: "PRINT_COUNT",
            6: "CUTTER_COUNT"
        }

        # 初始化缓冲区
        # Fix buffer initialization with proper ctypes syntax
        buf_size = c_uint(256)
        buffer = (ctypes.c_wchar * buf_size.value)()  # Use fully qualified c_wchar
        result = hiti_dll.HITI_GetDeviceInfoW(
            printer,
            info_type,
            buffer,
            byref(buf_size)
        )

        if result == 0:
            data = buffer[:buf_size.value]
            
            # 处理特殊数据类型
            if info_type == 4:  # 色带信息
                ribbon_type, ribbon_count = data
                return jsonify({
                    "type": HITI_DEVINFO[info_type],
                    "ribbon_type": ribbon_type,
                    "remaining_count": ribbon_count
                })
            elif info_type in [5, 6]:  # 计数器
                return jsonify({
                    "type": HITI_DEVINFO[info_type],
                    "count": data[0]
                })
            else:  # 字符串信息
                return jsonify({
                    "type": HITI_DEVINFO[info_type],
                    "data": ''.join(data).strip('\x00')
                })
        else:
            error_code = ctypes.GetLastError()
            return jsonify({
                "error": "获取设备信息失败",
                "system_error": error_code,
                "message": ctypes.FormatError(error_code)
            }), 500
    else:
        return jsonify({
            "info": f"模拟信息 - 打印机: {printer}, 类型: {info_type}",
            "note": "macOS 模拟数据"
        })

# 添加执行命令的接口
@app.route('/command', methods=['POST'])
def do_command():
    if not request.is_json:
        return jsonify({"error": "请求必须是JSON格式"}), 400
        
    data = request.get_json()
    printer = data.get('printer')
    command = data.get('command')
    
    if not printer or command is None:
        return jsonify({"error": "缺少必要参数：printer或command"}), 400
    
    if is_windows:
        # 实际的 Windows 实现
        result = hiti_dll.HITI_DoCommandA(printer.encode('utf-8'), int(command))
        
        if result == 0:
            return jsonify({"success": True})
        else:
            return jsonify({"error": f"执行命令失败，错误码: {result}"}), 500
    else:
        # macOS 模拟数据
        return jsonify({
            "success": True,
            "note": f"模拟执行命令 - 打印机: {printer}, 命令: {command} (macOS 模拟模式)"
        })

if __name__ == '__main__':
    try:
        print("="*50)
        print("正在启动 Flask 服务器...")
        print(f"操作系统: {platform.system()}")
        print(f"Python 版本: {platform.python_version()}")
        print(f"工作目录: {os.getcwd()}")
        print(f"Windows 模式: {is_windows}")
        print("="*50)
        
        # 设置 host 为 0.0.0.0 使服务可以从其他设备访问
        app.run(host='0.0.0.0', port=8080, debug=False)
    except Exception as e:
        print(f"启动 Flask 服务器时发生错误: {e}")
        import traceback
        traceback.print_exc()