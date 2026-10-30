# Android P1 验证清单

状态：本地实现，未发布、未签名、未真机验收。

## 构建环境

Java 17、Gradle 8.10.2、Android SDK platform 35/build-tools 35.0.0，Go 1.26.2。当前仓库没有 wrapper；使用 scripts/build-android-debug.ps1 显式指定工具路径，不运行现有发布工作流。

## 设备验收（记录 SHA、APK SHA256、设备型号、Android/WebView 版本）

- 安装 debug APK，包名应为 io.github.lyl7410.cfdata.enhanced.debug；Manifest 的 Activity 应为 com.cfdata.web.MainActivity。
- 与原版同时启动，确认独立数据、动态端口及页面；占用 13335 不应影响增强版。
- 退出重启、旋转、前后台恢复；启动失败必须显示错误，不加载其他进程页面。
- 使用授权小规模输入验证扫描/取消、配置保存、导出中文文件名。
- 提升 versionCode 后以同一测试签名覆盖安装，验证数据保留；正式签名覆盖与错误签名拒绝留到发布前。

## 发布门槛

首次公开前固定 applicationId、长期证书和版本递增规则；受控保管密钥。正式签名、Secrets、Release、自动同步本次均不执行。工作流修改仍受 AGENTS.md 限制。

## 本地验证记录（2026-10-29）

- Windows amd64 / Go 1.26.2：scripts/verify.ps1 全流程退出 0。
- go test ./...、go vet ./... 通过；含更新渠道与私有端口握手测试。
- Windows amd64、Linux amd64、Android arm64 后端构建通过。
- git diff --check、APK 构建脚本 PowerShell 语法检查通过。
- 当前环境未检测到 Java/Gradle/Android SDK/adb；未执行 Gradle 配置、APK 构建、Android lint 或真机验收，不能宣称 Android 可交付。
- 尚未提交、推送或创建 PR；未修改 .github，未生成正式密钥或配置 Secrets。
