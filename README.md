# Ether-UI
**一个高级Web面板 • 基于Xray Core构建**

![](https://img.shields.io/github/v/release/X-ether/Ether-UI.svg)
![](https://img.shields.io/docker/pulls/X-ether/Ether-UI.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/X-ether/Ether-UI)](https://goreportcard.com/report/github.com/X-ether/Ether-UI)
[![Downloads](https://img.shields.io/github/downloads/X-ether/Ether-UI/total.svg)](https://img.shields.io/github/downloads/X-ether/Ether-UI/total.svg)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)

> **免责声明：** 本项目仅用于个人学习和交流，原作者对于其用途和具体功能毫不知情。请不要将其用于非法用途，绝对禁止将其用于越过任何形式的网络审查。请勿在生产环境中使用。

**如果你觉得这个项目对你有帮助，不妨给一个**:star2:

## 概述
| 功能                                   |      是否启用       |
| -------------------------------------- | :----------------: |
| 多协议支持                             | :heavy_check_mark: |
| 多语言支持                             | :heavy_check_mark: |
| 多客户端/入站                          | :heavy_check_mark: |
| 高级流量路由接口                       | :heavy_check_mark: |
| 客户端 & 流量 & 系统状态监控           | :heavy_check_mark: |
| 基于首次使用的日期 & 流量上限          | :heavy_check_mark: |
| REST API                               | :heavy_check_mark: |
| TG 机器人（数据库备份 + 管理 + 客户端）| :heavy_check_mark: |
| 订阅服务（链接 + 信息）                | :heavy_check_mark: |
| 深度搜索                               | :heavy_check_mark: |
| 深色/浅色主题                          | :heavy_check_mark: |

## 安装 & 升级到最新版本

```sh
bash <(curl -Ls https://raw.githubusercontent.com/X-ether/Ether-UI/master/install.sh)
```

## 安装自定义版本

**步骤1：** 要安装你想要的版本，请使用以下命令。例如，安装`1.8.0`版本：

```sh
VERSION=1.8.0 bash <(curl -Ls "https://raw.githubusercontent.com/X-ether/Ether-UI/refs/tags/$VERSION/install.sh") $VERSION
```

## 手动安装 & 升级

### 使用方法

1. 直接将最新版本的压缩包下载到你的服务器，执行以下命令：

```sh
ARCH=$(uname -m)
case "${ARCH}" in
  x86_64 | x64 | amd64) XUI_ARCH="amd64" ;;
  i*86 | x86) XUI_ARCH="386" ;;
  armv8* | armv8 | arm64 | aarch64) XUI_ARCH="arm64" ;;
  armv7* | armv7) XUI_ARCH="armv7" ;;
  *) XUI_ARCH="amd64" ;;
esac

wget https://github.com/X-ether/Ether-UI/releases/latest/download/Ether-UI-linux-${XUI_ARCH}.tar.gz
```

2. 下载完压缩包后，执行以下命令来安装或升级 Ether-UI：

```sh
ARCH=$(uname -m)
case "${ARCH}" in
  x86_64 | x64 | amd64) XUI_ARCH="amd64" ;;
  i*86 | x86) XUI_ARCH="386" ;;
  armv8* | armv8 | arm64 | aarch64) XUI_ARCH="arm64" ;;
  armv7* | armv7) XUI_ARCH="armv7" ;;
  *) XUI_ARCH="amd64" ;;
esac
cd /root/
rm Ether-UI/ /usr/local/Ether-UI/ /usr/bin/Ether-UI -rf
tar zxvf Ether-UI-linux-${XUI_ARCH}.tar.gz
chmod +x Ether-UI/Ether-UI Ether-UI/bin/xray-linux-* Ether-UI/Ether-UI.sh
cp Ether-UI/Ether-UI.sh /usr/bin/Ether-UI
cp -f Ether-UI/Ether-UI.service /etc/systemd/system/
mv Ether-UI/ /usr/local/
systemctl daemon-reload
systemctl enable Ether-UI
systemctl restart Ether-UI
```

## 使用Docker安装

### 使用方法

**步骤1：** 安装 Docker

```shell
curl -fsSL https://get.docker.com | sh
```

**步骤2：** 克隆项目仓库：

```sh
git clone https://github.com/X-ether/Ether-UI.git
cd X-UI
```

**步骤3：** 启动服务

```sh
docker compose up -d
```

或

```shell
mkdir Ether-UI && cd Ether-UI
docker run -itd \
    -p 54321:54321 -p 443:443 -p 80:80 \
    -e XRAY_VMESS_AEAD_FORCED=false \
    -v $PWD/db/:/etc/X-UI/ \
    -v $PWD/cert/:/root/cert/ \
    --name Ether-UI --restart=unless-stopped \
    X-ether/Ether-UI:latest
```

升级到最新版本：

```sh
cd X-UI
docker compose down
docker compose pull X-UI
docker compose up -d
```

从docker中移除 Ether-UI：

```sh
docker stop Ether-UI
docker rm X-UI
cd --
rm -r X-UI
```

构建你自己的镜像：

```shell
docker build -t X-UI .
```

## 语言支持

- 英语
- 中文
- 波斯语
- 俄语
- 越南语

## 功能

- 支持的协议包括VLESS、VMess、Trojan、Shadowsocks、Dokodemo-door、SOCKS、HTTP、Wireguard
- 支持的XTLS协议，包括Vision和REALITY
- 高级流量路由接口，支持PROXY Protocol、反向代理、外部代理和透明代理，以及多域名、SSL证书和端口
- 支持使用Wireguard外部生成Cloudflare WARP
- 提供用于Xray模板配置的交互式JSON接口
- 高级入站和出站配置接口
- 基于首次使用的日期和流量上限的客户端流量限制与到期日期
- 显示在线客户端、流量统计和系统状态监控
- 深度数据库搜索
- 显示到期日期或超出流量上限的已耗尽客户端
- 提供多链路订阅服务
- 数据库导入和导出
- 一键申请SSL证书及自动续期
- 提供安全访问面板和订阅服务的HTTPS支持（自带域名+SSL证书）
- 深色/浅色主题切换

## 推荐操作系统

- Ubuntu 20.04+
- Debian 11+
- CentOS 8+
- OpenEuler 22.03+
- Fedora 36+
- Arch Linux
- Parch Linux
- Manjaro
- Armbian
- AlmaLinux 8.0+
- Rocky Linux 8+
- Oracle Linux 8+
- OpenSUSE Tubleweed
- Amazon Linux 2023

## API 路由

### 使用方法

- `/login` 使用 `PUSH` 用户数据：`{username: '', password: ''}` 进行登录
- `/xui/API/inbounds` 基于以下操作：

| 方法 | 路径                                 | 动作                                     |
| :--: | ------------------------------------ | ---------------------------------------- |
| `GET`  | `"/"`                              | 获取所有入站                            |
| `GET`  | `"/get/:id"`                       | 通过 inbound.id 获取入站                 |
| `GET`  | `"/createbackup"`                  | Telegram 机器人向管理员发送备份          |
| `POST` | `"/add"`                           | 添加入站                                |
| `POST` | `"/del/:id"`                       | 删除入站                                |
| `POST` | `"/update/:id"`                    | 更新入站                                |
| `POST` | `"/addClient/"`                    | 向入站添加客户端                        |
| `POST` | `"/:id/delClient/:clientId"`       | 通过 clientId 删除客户端\*              |
| `POST` | `"/updateClient/:clientId"`        | 通过 clientId 更新客户端\*              |
| `GET`  | `"/getClient/:clientId"`           | 通过 clientId 获取客户端\*              |

