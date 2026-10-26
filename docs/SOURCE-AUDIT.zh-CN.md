# 源码与构建审查 v1.0

审查日期：2026-10-26。

## 基准

Fork main 与上游 main 均为 `49d5571f4e44d1fe265c44f9fa289bc296b865f3`，审查时左右差异提交数为 0/0。
上游：https://github.com/PoemMisty/CFData-WEB 。本记录是该提交的快照，不代表上游以后不变化。

## 实际结构

| 路径 | 实际用途及维护约定 |
| --- | --- |
| `combined_refactor/go.mod` | 活跃 Go 模块 cfdata，要求 Go 1.26.0；本次使用 1.26.2 验证 |
| `combined_refactor/main.go`, `server.go` | CLI/Web 服务、WebSocket、版本查询 |
| `combined_refactor/types.go` | embed 同目录 index.html、login.html、favicon.png |
| `combined_refactor/ca.go` | embed CA 证书；现有发布流程会联网刷新证书 |
| `combined_refactor/tasks_*.go` | 官方/非标扫描、测速及复验逻辑，桌面与 APK 共用 |
| 根目录 `cfdata.go`, `index.html`, `go.mod` | 历史代码；根模块 Go 1.25.4，不能当作当前产品入口 |
| `app/src/main/java/com/cfdata/web/MainActivity.java` | WebView 外壳，以 ProcessBuilder 启动 nativeLibraryDir/libcfdata.so；本地 127.0.0.1:13335 |
| `app/build.gradle` | AGP 8.7.3、SDK 35、minSdk 26；当前 applicationId/namespace 为 com.cfdata.web，versionCode 固定 1 |
| `.github/workflows/build.yml` | tag v* / 手动触发，编译 Windows/Linux/macOS amd64/arm64 与 Android arm64，SHA256、签名 APK、直接创建非草稿 Release |

Android 的 libcfdata.so 是按既有方案打包的 Go 可执行后端，并非 JNI 库。当前仅支持 arm64 APK，不能据此承诺其他 ABI。
Windows 运行 Go 程序后用浏览器访问网页；目前不是 Windows 原生桌面安装包。
Android 有文件选择、保存文件、深浅色等 Java 桥接，系统能力仍须单独验收。

## 发现及处理

1. Fork 版本 API 和下载链接仍指向上游：首阶段统一改到 Fork；Fork 无 Release 时维持错误提示，不回退推荐原版。
2. 既有流水线没有 PR 测试门禁，发布需要四个 Android 签名 Secrets；Fork 不会继承上游 Secrets。本次提供方案与本地脚本，未启用发布。
3. APK 与原版共用包名，无法作为独立增强版并存；versionCode 固定也不适合作为持续升级策略。列入下一阶段，首次分发前必须解决。
4. 当前没有 Go 测试文件；首阶段添加更新渠道离线测试。更完整的扫描/导出/配置测试待补。
5. Go 服务默认监听全部地址，Android 显式限制回环。Windows 本地测试也应显式使用 `-host 127.0.0.1`。
6. Android 固定端口可能与并行运行的原版冲突；独立包名并不能消除端口冲突，后续须设计可分配端口与生命周期处理。
7. 本审查不是完整依赖漏洞或 WebView 安全审计；后续补充桥接来源限制、会话隔离和依赖清单检查。

## 许可证

仓库 LICENSE 是 GNU GPL version 3 完整文本。保留 LICENSE、版权、来源与无担保说明；修改版本注明修改内容与日期。
对外提供 EXE/APK 时，同一发布页提供对应提交的完整可构建源码和构建/安装说明，包含自定义模块、构建脚本及所需许可说明，不能只发布二进制。
后续核对 Go 依赖、证书及其他资源的许可并生成第三方清单。不得推断商标授权，也不将增强版冒充上游官方。

## 验证范围

首阶段验证命令和结果见 [维护文档](MAINTENANCE.zh-CN.md)。Go 交叉编译不等于 Android APK 构建或手机实测，未执行的项目明确保留为待办。
