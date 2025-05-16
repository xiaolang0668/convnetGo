# ConvnetGo

ConvnetGo 是一个基于 Go 语言开发的 P2P 网络连接工具，支持 Windows 和 Linux 系统。它允许用户通过 TURN/STUN 服务建立点对点连接，实现安全的网络通信。

## 功能特点

- 支持 Windows 和 Linux 系统
- 基于 WebRTC 的 P2P 连接
- TAP 虚拟网卡支持
- 自动重连机制
- Web 界面控制
- 支持 TCP/UDP 端口范围配置
- 支持服务端和客户端模式
- 自动保存连接配置

## 系统要求

### Windows
- 需要安装 TAP 驱动
  - 下载地址：https://openvpn.net/community-downloads/
  - 安装时仅选择 TAP 驱动组件

### Linux
- 需要 root 权限来创建和配置 TAP 设备

## 配置说明

配置文件 `convnet.json` 包含以下主要参数：

```json
{
  "Server": "服务器地址",
  "UUID": "客户端唯一标识",
  "ClientID": "客户端昵称",
  "ServerPort": "服务器端口",
  "AllowConnect": "是否允许被连接",
  "AllowTcpPortRange": [
    {
      "Start": 10000,
      "End": 20000
    }
  ],
  "AllowUdpPortRange": [
    {
      "Start": 10000,
      "End": 20000
    }
  ]
}
```

## 运行模式

### 服务端模式
```bash
convnetgo -s
```
- 默认监听端口：13903 (TCP)
- TURN 服务端口：13902 (UDP)

### 客户端模式
```bash
convnetgo
```
- 启动后会自动打开 Web 控制界面：http://127.0.0.1:8094
- 自动生成 UUID 和随机昵称（首次运行）
- 支持自动重连服务器

## 连接管理

### 自动连接
- 通过 `autoConnectPeer.txt` 文件管理自动连接的节点
- 格式示例：`CVNID://saiboot.com:13903@[PublicID]`

### 连接状态
- 通过 Web 界面查看连接状态
- 支持手动连接/断开操作
- 实时显示在线用户列表

## API 接口

主要 HTTP API 端点：
- `/api/user/list` - 获取用户列表
- `/api/info` - 获取客户端信息
- `/api/info/update` - 更新用户信息
- `/api/peer/connect` - 连接到指定节点
- `/api/client/connect` - 连接到服务器
- `/api/client/disconnect` - 断开服务器连接

## 技术实现

### 网络通信
- 使用 WebRTC 进行 P2P 通信
- 基于 TAP 设备实现虚拟网络
- 支持 NAT 穿透
- 使用 TCP 长连接保持会话

### 安全机制

1. **数据加密**
   - 使用 AES-CBC 模式加密敏感数据
   - 随机生成 IV (初始化向量)
   - 采用 PKCS7 填充标准
   - Base64 编码传输加密数据

2. **身份认证**
   - 支持双重身份认证机制
   - 私有身份（UUID）用于本地认证
   - 公开身份（PublicID）用于 P2P 连接

3. **TURN 服务器安全**
   - 使用 realm 域隔离
   - 动态生成的认证密钥
   - 基于用户名和密码的访问控制
   - TURN 凭证通过 AES-CBC 加密传输

## 注意事项

1. 首次运行会自动生成配置文件
2. Windows 系统必须预先安装 TAP 驱动
3. 确保配置的端口范围在防火墙中已开放
4. 服务端需要同时开放 TCP 和 UDP 端口
5. 建议妥善保存生成的 UUID，它是客户端的唯一标识

## TODO 功能

- 黑名单功能
- 端口屏蔽
- 密码访问控制