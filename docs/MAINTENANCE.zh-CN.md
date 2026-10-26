# 开发、发布与上游同步 v1.0

## 工作约定

所有变更在新分支，通过 PR 指向 Fork 的 main，禁止直接推 main、自动合并或直接推上游。
仓库 AGENTS.md 禁止修改 `.github/` 及标记“原版源码”的历史文件；本次遵循该限制。CI/CD 是待实施方案，不能将本文档视为已经配置成功的流水线。
用户本次明确授权提交与创建 PR。没有创建发布、签名密钥、Secrets 或自动同步任务。

## 本地复现

安装 Go 1.26.2（最低按 combined_refactor/go.mod）、Windows PowerShell。验证脚本不需要扫描真实公网 IP、不启动服务、不读取签名密钥。

```powershell
./scripts/verify.ps1
# 若 Go 不在 PATH 中：
./scripts/verify.ps1 -Go 'C:/tools/go/bin/go.exe'
# 缓存目录无写权限时，先将 GOCACHE/GOMODCACHE/GOPATH 指向可写目录。
```

脚本运行活跃模块 `go test ./...`、`go vet ./...`，构建 `release_assets/cfdata-test.exe`、Linux 测试版 `release_assets/cfdata-test`、Android arm64 后端；任何一步失败即退出。
这不包括 Gradle APK 构建、Android lint、安装、签名或真机测试。
手动 Windows 冒烟测试建议 `./release_assets/cfdata-test.exe -host 127.0.0.1`，浏览器打开终端地址，完成后退出程序。

Android 复现沿用当前发布方案：Java 17、Gradle 8.10.2、Android SDK 35；仓库没有 Gradle wrapper，后续应补充固定版本及校验值。
在 `combined_refactor/` 下设置 `GOOS=android GOARCH=arm64 CGO_ENABLED=0`，构建至 `../app/src/main/jniLibs/arm64-v8a/libcfdata.so`（先创建目录），再在根目录执行 `gradle assembleDebug --no-daemon`。
只有正式发布才使用专用签名与 `assembleRelease`；不要分发 debug 签名包作为稳定更新通道。

## 拟实施的 CI/CD（后续独立 PR）

| 触发 | 检查 | 权限与输出 |
| --- | --- | --- |
| pull_request 到 main；手动验证 | Go test/vet；Windows runner 运行测试；Linux/Windows/Android 后端构建；网页脚本检查 | contents:read；checkout 禁止持久化凭据；无 Secrets、无发布 |
| 同一 PR 的 Android job | Java17/Gradle8.10.2/SDK35、后端入包、assembleDebug/lint、APK 内含 arm64 后端检查 | 无发布签名；构建产物标记 SHA 和短期保留 |
| 受控增强版 tag | 同一 SHA 的双端构建，APK apksigner 验证、包名/versionCode 校验、SHA256、对应源码与许可 | environment 人工审批；仅发布 job contents:write；先草稿/预发布 |
| 上游检查（后续可选手动/周期） | 比较记录 SHA；有差异才创建或更新同步 PR，验证状态单独显示 | 仅所需 contents/pull-requests:write；不自动合并或发布 |

不要用 pull_request_target 执行 PR 代码或给不可信 PR 注入签名 Secrets。固定或审查第三方 action 版本/SHA，避免把现有工作流中尚未验证的 action 主版本当作可用保证。
Windows 与 Android 产物必须带相同 git SHA、增强版版本及上游基准；记录 Go/JDK/SDK/Gradle 版本、依赖锁定信息、证书输入哈希。现有在线下载 CA 会改变构建输入，后续固定证书刷新 PR，而非发布时不加记录地刷新。
首次发布前配置 Actions、签名与环境门禁；仓库当前工作流可能直接发布 Release，因此不要用打 v* 标签的方式“测试”。

## Android 身份和升级

拟定独立 applicationId：`io.github.lyl7410.cfdata.enhanced`；Java namespace 可保留，Manifest 相对 Activity 名须在构建 APK 中核对。
首次对外分发后固定 applicationId 与证书；每次提升 versionCode。稳定版和 debug 构建如需并存，debug 使用额外 suffix，明确其数据不共用。
四个签名 Secrets 见路线图。密钥生成应由维护者在受控环境执行，离线备份并记录证书指纹，避免机器人替用户生成无法长期保管的发布身份。
升级前备份配置/结果；包名变化会建立独立应用沙箱，原版数据不会自动迁移。先做脱敏导出/导入；数据格式迁移需保留旧副本。

## 上游同步 PR 操作

```text
git remote add upstream https://github.com/PoemMisty/CFData-WEB.git  # 仅首次
git fetch origin
git fetch upstream
git switch -c sync/upstream-YYYYMMDD origin/main
git merge --no-commit --no-ff upstream/main
# 人工审查并解决冲突，运行验证脚本与双端矩阵，更新上游基准文档
git commit
git push origin sync/upstream-YYYYMMDD
# 在 lyl7410/CFData-WEB 创建指向 main 的 PR；人工审核后才合并
```

实际操作时用已核实的 upstream SHA 固定同步目标，PR 写明起止 SHA、冲突路径、测试证据和未通过项。
若无差异，不创建空同步 PR。若仓库限制仍禁止 `.github/`，遇到上游工作流变更必须单独获得修改该目录的明确授权，再处理；不要默默带入禁止路径。
默认合并保留上游历史；不要 rebase/强推 main。功能模块专用提交尽量小，以便必要时移植。
回滚通过新分支 revert 问题提交并建立 PR；数据 schema 升级有独立备份和回退步骤。Android 回退低 versionCode 通常不能直接覆盖安装，应以更高 versionCode 发布修复版。

## 本次执行记录（2026-10-26）

环境：Windows amd64，Go 1.26.2；工具链和缓存置于工作目录外侧的临时工具目录，不提交仓库。

| 检查 | 结果 |
| --- | --- |
| `scripts/verify.ps1` | 全流程退出码 0 |
| `go test ./...` | 通过，cfdata；5 个更新渠道子用例，不访问真实 GitHub API |
| `go vet ./...` | 通过 |
| Windows amd64 构建 | 通过，release_assets/cfdata-test.exe |
| Linux amd64 构建 | 通过，release_assets/cfdata-test；按仓库指南生成测试版 |
| Android arm64 Go 后端构建 | 通过，release_assets/cfdata-test-android-arm64 |
| Windows 测试版 `-h` | 可运行并显示参数帮助；启动逻辑仍查询版本，Fork 暂无 Release 返回 404，未回退上游 |
| 文档相对链接检查、git diff --check | 通过 |
| Android APK/签名/真实设备、完整浏览器操作 | 未执行，仍待 P1；当前环境未配置 Java/Gradle/SDK |
| GitHub PR CI | 未新增 CI；现有发布工作流不监听 PR，不能将本地通过写成远程 CI 通过 |

二进制在忽略目录 release_assets 中，仅作验证；本次不发布 EXE/APK、不运行公网扫描。后端交叉编译不替代完整 APK 与真机验收。
