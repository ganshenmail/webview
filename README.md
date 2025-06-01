# Go WebView 桌面应用框架

## 概述
本项目是一个使用 Go 语言和 webview_go 库构建的桌面应用框架。它采用混合开发模式，结合了原生窗口管理和Web技术实现的用户界面。

## 功能特性
- **窗口管理**: 最小化、最大化、关闭、居中、调整窗口大小以及动态设置标题
- **内置Web服务器**: 从 `app` 目录提供静态内容服务
- **系统集成**: 提供剪贴板、进程和文件系统操作的API
- **可配置**: 通过JSON配置文件自定义窗口和服务器设置
- **动态控制**: 支持运行时修改窗口标题和尺寸

## 安装
1. 确保已安装Go语言环境(推荐1.16+版本)
2. 克隆本仓库
3. 安装依赖:
   ```bash
   go mod download
   ```

## 构建与运行
构建并运行应用程序:
```bash
go build -o main.exe
./main.exe
```

## 默认配置
应用程序使用以下默认配置(代码中定义):

### 窗口配置
- 标题: "MATAVIEW"
- 初始尺寸: 1280x720 像素
- 最小尺寸: 800x600 像素
- 自动居中: 是

### 服务器配置
- 地址: 127.0.0.1
- 端口: 随机

### 调试模式
通过命令行参数 `--debug` 启用调试模式

> 提示: 如需修改配置，需要直接编辑main.go源代码并重新编译

## 项目结构
```
webview_go/
├── app/                # 前端资源目录
│   ├── index.html      # 主页面
│   ├── prism.css       # 代码高亮样式
│   └── prism.js        # 代码高亮脚本
├── main.go             # 主程序入口
├── winapi.go           # Windows API封装
├── window.go           # 窗口管理逻辑
├── rsrc.syso           # Windows资源文件
├── go.mod              # Go模块定义
├── go.sum              # 依赖校验文件
├── main.exe            # 编译后的可执行文件
└── README.md           # 项目文档
```

## 注意事项
1. Windows API定义应放在winapi.go中
2. 避免在window.go中重复定义Windows常量和函数
3. 如需扩展Windows API功能，请修改winapi.go文件

## JavaScript API

通过 `mata` 对象提供的API:

### 窗口控制
```javascript
// 设置窗口标题
mata.win.title("新标题");

// 调整窗口尺寸
mata.win.resize(宽度, 高度);
```

## 扩展应用
1. 在 `app` 目录中添加新的HTML/JS/CSS文件
2. 扩展Go后端功能
3. 通过 `mata` JavaScript对象暴露新的API
