# NekoBox for Android

[![API](https://img.shields.io/badge/API-21%2B-brightgreen.svg?style=flat)](https://android-arsenal.com/api?level=21)
[![Releases](https://img.shields.io/github/v/release/MatsuriDayo/NekoBoxForAndroid)](https://github.com/MatsuriDayo/NekoBoxForAndroid/releases)
[![License: GPL-3.0](https://img.shields.io/badge/license-GPL--3.0-orange.svg)](https://www.gnu.org/licenses/gpl-3.0)

sing-box / universal proxy toolchain for Android.

一款使用 sing-box 的 Android 通用代理软件.

> **关于本 Fork**:本项目 Fork 自 [MatsuriDayo/NekoBoxForAndroid](https://github.com/MatsuriDayo/NekoBoxForAndroid),
> 因上游作者已长期未维护,无法适配 sing-box 1.11+ 的配置变更和 Android 平台的限制。
> 本 Fork 在 1.4.2 基础上,针对 **sing-box 1.14+** 做了一系列兼容性与构建修复,详见下方[修复内容](#修复内容)。

## 下载 / Downloads

[![GitHub All Releases](https://img.shields.io/github/downloads/Matsuridayo/NekoBoxForAndroid/total?label=downloads-total&logo=github&style=flat-square)](https://github.com/MatsuriDayo/NekoBoxForAndroid/releases)

[GitHub Releases 下载](https://github.com/MatsuriDayo/NekoBoxForAndroid/releases)

**Google Play 版本自 2024 年 5 月起已被第三方控制，为非开源版本，请不要下载。**

**The Google Play version has been controlled by a third party since May 2024 and is a non-open
source version. Please do not download it.**

## 修复内容

### 1. Legacy Inbound Fields 迁移
**错误**:`decode config: inbounds[1]: legacy inbound fields are deprecated in sing-box 1.11.0 and removed in sing-box 1.13.0`

**原因**:`sniff` / `sniff_override_destination` / `domain_strategy` 等旧字段被移除。

**修复**:迁移为 route rule actions。在 `ConfigBuilder.kt` 中为需要 sniff/resolve 的 inbound 收集 tag,
并在 route 规则最前面插入 `action = "sniff"` / `action = "resolve"` 规则。
`SingBoxOptions.java` 的 `Rule_DefaultOptions` 新增 `sniffer` / `timeout` / `strategy` 字段。

### 2. Android netlink 禁用
**错误**:`create network monitor: netlink socket in Android is banned by Google`

**原因**:sing-box 1.11+ 会查找 `adapter.PlatformInterface` 决定是否使用 app 自带的接口监视器;
未注册时回退到原始 netlink 套接字(Android 禁止非 root 应用使用)。

**修复**:在 `libcore/box.go` 中同时注册 `platform.Interface` 和 `adapter.PlatformInterface`。
新增 `libcore/platform_adapter.go` 实现 `boxAdapterPlatformInterface`,使 `UsePlatformDefaultInterfaceMonitor()` 返回 `true`。

### 3. Legacy Tun Address Fields
**错误**:`initialize inbound[0]: legacy tun address fields are deprecated in sing-box 1.10.0 and removed in sing-box 1.12.0`

**原因**:`inet4_address` / `inet6_address` 等旧字段被合并。

**修复**:`SingBoxOptions.java` 的 `Inbound_TunOptions` 新增 `address` / `route_address` / `route_exclude_address`;
`ConfigBuilder.kt` 改用合并后的 `address` 字段。

### 4. DNS Detour 空 direct 检查
**错误**:`start dns/https[dns-direct]: detour to an empty direct outbound makes no sense`

**原因**:sing-box 1.14 新增检查,拒绝 detour 到空 direct outbound;但 DNS 通过 direct 解析是有意为之。

**修复**:`sing-box/dns/transport_dialer.go` 的 `NewLocalDialer` / `NewRemoteDialer` 设置 `DisableEmptyDirectCheck: true`。

### 5. geosite / geoip 数据库缺失
**错误**:`parse rule-set[0]: open /storage/emulated/0/Android/data/moe.nb4a.debug/files/geosite.db: no such file or directory`

**原因**:APK 未打包 geosite/geoip 数据库;`generateRuleSet` 对 `geosite:`/`geoip:` 设了 `format = "binary"`,
导致 sing-box 把前缀当字面文件路径解析。

**修复**:
- 运行 `buildScript/lib/assets.sh` 下载 geosite.db / geoip.db 到 `app/src/main/assets/sing-box/`;
- `SingBoxOptionsUtil.kt` 的 `generateRuleSet` 去掉 `format = "binary"`,让 sing-box 走注入的 headless 规则路径。

### 6. RuleSet missing format
**错误**:`decode config: route.rule_set[0]: missing format`

**原因**:`geosite:`/`geoip:` 伪路径无文件扩展名,`ruleSetDefaultFormat` 返回空。

**修复**:`sing-box/option/rule_set.go` 的 `ruleSetDefaultFormat` 识别 `geoip:`/`geosite:` 前缀,返回默认 `source` format。

## 构建说明

### 环境要求

- **JDK 17+**
- **Go 1.25+**
- **Android SDK**:Platform 35,Build-tools 35.0.0
- **Android NDK**:25.0.8775105(或更高)
- **Gradle**:8.9+(可直接用项目自带的 `./gradlew`)
- **gomobile-matsuri / gobind-matsuri**:用于把 Go 核心编译为 Android AAR

### 构建步骤

**1. 准备 Go 工具链**

从 [MatsuriDayo/gomobile](https://github.com/MatsuriDayo/gomobile) 获取 `gomobile-matsuri` 与 `gobind-matsuri`,安装到 `$GOPATH/bin`:

```bash
go install github.com/MatsuriDayo/gomobile/cmd/gomobile-matsuri@latest
go install github.com/MatsuriDayo/gomobile/cmd/gobind-matsuri@latest
```

**2. 构建 Go 核心库(libcore.aar)**

```bash
cd libcore
export PATH=$PATH:$(go env GOPATH)/bin
bash build.sh
```

生成的 `libcore.aar` 会自动复制到 `app/libs/`。

**3. 构建 APK**

```bash
cd ..
./gradlew assembleOssDebug
```

产物:`app/build/outputs/apk/oss/debug/NekoBox-1.4.2-arm64-v8a-debug.apk`

### 安装
拷贝 APK 到设备后手动安装。

## 更新日志 & Telegram 发布频道 / Changelog & Telegram Channel

https://t.me/Matsuridayo

## 项目主页 & 文档 / Homepage & Documents

https://matsuridayo.github.io

## 支持的代理协议 / Supported Proxy Protocols

* SOCKS (4/4a/5)
* HTTP(S)
* SSH
* Shadowsocks
* VMess
* Trojan
* VLESS
* AnyTLS
* ShadowTLS
* TUIC
* Hysteria 1/2
* WireGuard
* Trojan-Go (trojan-go-plugin)
* NaïveProxy (naive-plugin)
* Mieru (mieru-plugin)

请到[这里](https://matsuridayo.github.io/nb4a-plugin/)下载插件以获得完整的代理支持.

Please visit [here](https://matsuridayo.github.io/nb4a-plugin/) to download plugins for full proxy
supports.

## 支持的订阅格式 / Supported Subscription Format

* 一些广泛使用的格式 (如 Shadowsocks, ClashMeta 和 v2rayN)
* sing-box 出站

仅支持解析出站，即节点。分流规则等信息会被忽略。

* Some widely used formats (like Shadowsocks, ClashMeta and v2rayN)
* sing-box outbound

Only resolving outbound, i.e. nodes, is supported. Information such as diversion rules are ignored.

## 免责声明
本 Fork 仅用于本地适配与学习,不包含恶意代码。上游许可证为 GPL-3.0。
