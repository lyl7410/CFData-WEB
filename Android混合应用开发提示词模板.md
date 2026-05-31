# Android 混合应用开发完整提示词模板

## 项目概述

本项目是一个**Go 后端 + Android WebView 前端**的混合架构应用，实现了跨平台功能复用和原生体验的完美结合。

---

## 一、项目结构设计

```
项目根目录/
├── go_backend/              # Go 后端源码
│   ├── main.go
│   ├── server.go
│   └── ...
└── android_app/             # Android 应用项目
    ├── app/
    │   ├── src/main/
    │   │   ├── java/com/example/app/
    │   │   │   └── MainActivity.java
    │   │   ├── jniLibs/
    │   │   │   ├── arm64-v8a/
    │   │   │   │   └── libbackend.so
    │   │   │   ├── armeabi-v7a/
    │   │   │   │   └── libbackend.so
    │   │   │   └── x86_64/
    │   │   │       └── libbackend.so
    │   │   ├── res/
    │   │   │   ├── mipmap-*/
    │   │   │   └── values/
    │   │   └── AndroidManifest.xml
    │   └── build.gradle
    ├── build.gradle
    └── settings.gradle
```

---

## 二、实现步骤指南

### 阶段 1：Go 后端准备

1. 确保 Go 后端可以编译为独立可执行文件
2. 使用 Go 交叉编译为 Android 支持的架构：
   - ARM64 (arm64-v8a)
   - ARM32 (armeabi-v7a)
   - x86_64 (模拟器)

### 阶段 2：Android 项目初始化

1. 使用 Android Studio 创建新项目或手动创建 Gradle 项目
2. 配置 `settings.gradle` 和根目录 `build.gradle`
3. 配置 `app/build.gradle`

### 阶段 3：核心代码实现

1. 编写 `MainActivity.java`
2. 配置 `AndroidManifest.xml`
3. 添加必要的资源文件

### 阶段 4：集成与测试

1. 将编译好的 Go 二进制文件放入 `jniLibs` 对应架构目录
2. 编译并安装 APK
3. 测试各项功能

---

## 三、核心代码模板

### 1. MainActivity.java 完整模板

```java
package com.your.package;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.ActivityNotFoundException;
import android.content.Intent;
import android.graphics.Color;
import android.graphics.Typeface;
import android.graphics.drawable.GradientDrawable;
import android.net.Uri;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.webkit.JavascriptInterface;
import android.webkit.ValueCallback;
import android.webkit.WebChromeClient;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.FrameLayout;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.TextView;
import android.widget.Toast;

import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class MainActivity extends Activity {
    // ================== 配置区域 ==================
    private static final int PORT = 13335;              // 本地服务端口
    private static final String BACKEND_LIBRARY = "libbackend.so";  // Go 后端库名
    private static final String APP_NAME = "MyApp";     // 应用名称
    private static final int REQUEST_CREATE_DOCUMENT = 1001;
    private static final int REQUEST_FILE_CHOOSER = 1002;

    // ================== 成员变量 ==================
    private final ExecutorService executor = Executors.newSingleThreadExecutor();
    private final Handler mainHandler = new Handler(Looper.getMainLooper());
    private WebView webView;
    private FrameLayout rootView;
    private View loadingView;
    private TextView loadingTitle;
    private TextView loadingMessage;
    private ProgressBar loadingSpinner;
    private Process backendProcess;
    private String pendingExportName;
    private byte[] pendingExportBytes;
    private ValueCallback<Uri[]> filePathCallback;

    // ================== 生命周期 ==================
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        rootView = new FrameLayout(this);
        webView = new WebView(this);
        rootView.addView(webView, new FrameLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
        loadingView = createLoadingView();
        rootView.addView(loadingView, new FrameLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
        setContentView(rootView, new ViewGroup.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT));
        applySystemBarInsets(rootView);
        configureWebView();
        startBackend();
    }

    // ================== UI 创建 ==================
    private View createLoadingView() {
        FrameLayout overlay = new FrameLayout(this);
        overlay.setBackgroundColor(Color.rgb(246, 248, 251));

        LinearLayout card = new LinearLayout(this);
        card.setOrientation(LinearLayout.VERTICAL);
        card.setGravity(Gravity.CENTER_HORIZONTAL);
        card.setPadding(dp(28), dp(30), dp(28), dp(28));
        card.setElevation(dp(8));

        GradientDrawable cardBackground = new GradientDrawable();
        cardBackground.setColor(Color.WHITE);
        cardBackground.setCornerRadius(dp(22));
        cardBackground.setStroke(dp(1), Color.argb(160, 226, 232, 240));
        card.setBackground(cardBackground);

        ImageView logo = new ImageView(this);
        logo.setImageResource(getApplicationInfo().icon);
        logo.setAdjustViewBounds(true);
        logo.setScaleType(ImageView.ScaleType.FIT_CENTER);
        LinearLayout.LayoutParams logoParams = new LinearLayout.LayoutParams(dp(64), dp(64));
        card.addView(logo, logoParams);

        loadingTitle = new TextView(this);
        loadingTitle.setText(APP_NAME);
        loadingTitle.setTextColor(Color.rgb(31, 41, 55));
        loadingTitle.setTextSize(24);
        loadingTitle.setTypeface(Typeface.DEFAULT_BOLD);
        loadingTitle.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams titleParams = new LinearLayout.LayoutParams(ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT);
        titleParams.topMargin = dp(18);
        card.addView(loadingTitle, titleParams);

        loadingMessage = new TextView(this);
        loadingMessage.setText("正在启动本地服务...");
        loadingMessage.setTextColor(Color.rgb(100, 116, 139));
        loadingMessage.setTextSize(14);
        loadingMessage.setGravity(Gravity.CENTER);
        LinearLayout.LayoutParams messageParams = new LinearLayout.LayoutParams(ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT);
        messageParams.topMargin = dp(8);
        card.addView(loadingMessage, messageParams);

        loadingSpinner = new ProgressBar(this);
        loadingSpinner.setIndeterminate(true);
        LinearLayout.LayoutParams spinnerParams = new LinearLayout.LayoutParams(dp(36), dp(36));
        spinnerParams.topMargin = dp(22);
        card.addView(loadingSpinner, spinnerParams);

        FrameLayout.LayoutParams cardParams = new FrameLayout.LayoutParams(Math.min(dp(320), getResources().getDisplayMetrics().widthPixels - dp(48)), ViewGroup.LayoutParams.WRAP_CONTENT, Gravity.CENTER);
        overlay.addView(card, cardParams);
        return overlay;
    }

    private void applySystemBarInsets(View view) {
        view.setOnApplyWindowInsetsListener((v, insets) -> {
            v.setPadding(0, insets.getSystemWindowInsetTop(), 0, insets.getSystemWindowInsetBottom());
            return insets;
        });
        view.requestApplyInsets();
    }

    // ================== WebView 配置 ==================
    @SuppressLint({"SetJavaScriptEnabled", "AddJavascriptInterface"})
    private void configureWebView() {
        WebSettings settings = webView.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setAllowFileAccess(true);
        settings.setAllowContentAccess(true);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_ALWAYS_ALLOW);
        webView.addJavascriptInterface(new AndroidBridge(), "AndroidBridge");
        
        webView.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                Uri uri = request.getUrl();
                if (uri != null && ("http".equals(uri.getScheme()) || "https".equals(uri.getScheme()))) {
                    String host = uri.getHost();
                    if ("127.0.0.1".equals(host) || "localhost".equals(host)) {
                        return false;
                    }
                    startActivity(new Intent(Intent.ACTION_VIEW, uri));
                    return true;
                }
                return false;
            }
        });
        
        webView.setWebChromeClient(new WebChromeClient() {
            @Override
            public boolean onShowFileChooser(WebView view, ValueCallback<Uri[]> callback, FileChooserParams params) {
                if (filePathCallback != null) {
                    filePathCallback.onReceiveValue(null);
                }
                filePathCallback = callback;
                Intent intent = params.createIntent();
                intent.setType("*/*");
                try {
                    startActivityForResult(intent, REQUEST_FILE_CHOOSER);
                } catch (ActivityNotFoundException e) {
                    filePathCallback = null;
                    Toast.makeText(MainActivity.this, "没有可用的文件选择器", Toast.LENGTH_LONG).show();
                    return false;
                }
                return true;
            }
        });
    }

    // ================== 后端启动 ==================
    private void startBackend() {
        executor.execute(() -> {
            try {
                setLoadingMessage("正在准备本地服务...");
                File backend = prepareBackendBinary();
                ProcessBuilder builder = new ProcessBuilder(
                        backend.getAbsolutePath(),
                        "-host", "127.0.0.1",
                        "-port", String.valueOf(PORT)
                );
                builder.directory(getFilesDir());
                builder.redirectErrorStream(true);
                backendProcess = builder.start();
                setLoadingMessage("正在连接本地服务...");
                waitForBackend();
                mainHandler.post(() -> {
                    setLoadingMessage("正在加载界面...");
                    webView.loadUrl("http://127.0.0.1:" + PORT);
                });
            } catch (Exception e) {
                mainHandler.post(() -> showStartupError("启动失败", "本地服务启动失败，请重新打开应用\\n" + e.getMessage()));
            }
        });
    }

    private void setLoadingMessage(String message) {
        mainHandler.post(() -> {
            if (loadingMessage != null) {
                loadingMessage.setText(message);
            }
        });
    }

    private void hideLoadingView() {
        View view = loadingView;
        if (view != null) {
            view.animate().alpha(0f).setDuration(180).withEndAction(() -> view.setVisibility(View.GONE)).start();
        }
    }

    private void showStartupError(String title, String message) {
        if (loadingView != null) {
            loadingView.setVisibility(View.VISIBLE);
            loadingView.setAlpha(1f);
        }
        if (loadingTitle != null) {
            loadingTitle.setText(title);
            loadingTitle.setTextColor(Color.rgb(185, 28, 28));
        }
        if (loadingMessage != null) {
            loadingMessage.setText(message);
        }
        if (loadingSpinner != null) {
            loadingSpinner.setVisibility(View.GONE);
        }
        Toast.makeText(this, title + ": " + message, Toast.LENGTH_LONG).show();
    }

    private File prepareBackendBinary() throws IOException {
        File nativeBackend = new File(getApplicationInfo().nativeLibraryDir, BACKEND_LIBRARY);
        if (nativeBackend.isFile()) {
            return nativeBackend;
        }
        throw new IOException("未找到内置后端: " + nativeBackend.getAbsolutePath());
    }

    private void waitForBackend() throws InterruptedException {
        long deadline = System.currentTimeMillis() + 15000;
        while (System.currentTimeMillis() < deadline) {
            if (isBackendReady()) {
                injectExportBridge();
                return;
            }
            Thread.sleep(300);
        }
        throw new IllegalStateException("本地服务启动超时");
    }

    private boolean isBackendReady() {
        HttpURLConnection connection = null;
        try {
            URL url = new URL("http://127.0.0.1:" + PORT + "/");  // 修改为你的健康检查端点
            connection = (HttpURLConnection) url.openConnection();
            connection.setConnectTimeout(500);
            connection.setReadTimeout(500);
            return connection.getResponseCode() < 500;
        } catch (IOException ignored) {
            return false;
        } finally {
            if (connection != null) {
                connection.disconnect();
            }
        }
    }

    // ================== JavaScript 桥接注入 ==================
    private void injectExportBridge() {
        mainHandler.post(() -> webView.setWebViewClient(new WebViewClient() {
            @Override
            public void onPageFinished(WebView view, String url) {
                super.onPageFinished(view, url);
                // 这里注入自定义的 JS 桥接代码
                String script = "(function(){"
                        + "if(!window.AndroidBridge||window.__androidBridgeInjected)return;"
                        + "window.__androidBridgeInjected=true;"
                        + "// 在这里覆盖或添加你的 JS 方法"
                        + "})();";
                view.evaluateJavascript(script, null);
                hideLoadingView();
            }

            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                Uri uri = request.getUrl();
                if (uri != null && ("http".equals(uri.getScheme()) || "https".equals(uri.getScheme()))) {
                    String host = uri.getHost();
                    if ("127.0.0.1".equals(host) || "localhost".equals(host)) {
                        return false;
                    }
                    startActivity(new Intent(Intent.ACTION_VIEW, uri));
                    return true;
                }
                return false;
            }
        }));
    }

    // ================== JavaScript 桥接类 ==================
    public class AndroidBridge {
        @JavascriptInterface
        public void saveTextFile(String fileName, String content) {
            mainHandler.post(() -> {
                pendingExportName = sanitizeFileName(fileName);
                byte[] bom = new byte[]{(byte) 0xEF, (byte) 0xBB, (byte) 0xBF};
                byte[] data = content == null ? new byte[0] : content.getBytes(StandardCharsets.UTF_8);
                pendingExportBytes = new byte[bom.length + data.length];
                System.arraycopy(bom, 0, pendingExportBytes, 0, bom.length);
                System.arraycopy(data, 0, pendingExportBytes, bom.length, data.length);
                Intent intent = new Intent(Intent.ACTION_CREATE_DOCUMENT);
                intent.addCategory(Intent.CATEGORY_OPENABLE);
                intent.setType("text/plain");
                intent.putExtra(Intent.EXTRA_TITLE, pendingExportName);
                startActivityForResult(intent, REQUEST_CREATE_DOCUMENT);
            });
        }

        // 在这里添加更多 JS 可调用的方法
        @JavascriptInterface
        public void showToast(String message) {
            mainHandler.post(() -> Toast.makeText(MainActivity.this, message, Toast.LENGTH_SHORT).show());
        }
    }

    // ================== 工具方法 ==================
    private String sanitizeFileName(String fileName) {
        String name = fileName == null ? "results.txt" : fileName.replaceAll("[\\\\/:*?\"<>|]", "_").trim();
        return name.isEmpty() ? "results.txt" : name;
    }

    private int dp(int value) {
        return Math.round(value * getResources().getDisplayMetrics().density);
    }

    // ================== Activity 结果处理 ==================
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if (requestCode == REQUEST_CREATE_DOCUMENT) {
            if (resultCode == RESULT_OK && data != null && data.getData() != null && pendingExportBytes != null) {
                try (OutputStream out = getContentResolver().openOutputStream(data.getData())) {
                    if (out != null) {
                        out.write(pendingExportBytes);
                        Toast.makeText(this, "已保存: " + pendingExportName, Toast.LENGTH_SHORT).show();
                    }
                } catch (IOException e) {
                    Toast.makeText(this, "保存失败: " + e.getMessage(), Toast.LENGTH_LONG).show();
                }
            }
            pendingExportBytes = null;
            pendingExportName = null;
            return;
        }
        if (requestCode == REQUEST_FILE_CHOOSER) {
            if (filePathCallback == null) {
                return;
            }
            Uri[] result = WebChromeClient.FileChooserParams.parseResult(resultCode, data);
            filePathCallback.onReceiveValue(result);
            filePathCallback = null;
        }
    }

    // ================== 返回键处理 ==================
    @Override
    public void onBackPressed() {
        if (webView != null && webView.canGoBack()) {
            webView.goBack();
            return;
        }
        super.onBackPressed();
    }

    // ================== 清理资源 ==================
    @Override
    protected void onDestroy() {
        if (backendProcess != null) {
            backendProcess.destroy();
        }
        executor.shutdownNow();
        if (webView != null) {
            webView.destroy();
        }
        super.onDestroy();
    }
}
```

### 2. AndroidManifest.xml 模板

```xml
<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    <!-- 网络权限 -->
    <uses-permission android:name="android.permission.INTERNET" />
    
    <!-- 按需添加其他权限 -->
    <!-- <uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" /> -->
    <!-- <uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" /> -->

    <application
        android:allowBackup="false"
        android:icon="@mipmap/ic_launcher"
        android:roundIcon="@mipmap/ic_launcher"
        android:usesCleartextTraffic="true"  <!-- 允许明文 HTTP（本地服务需要） -->
        android:theme="@style/AppTheme"
        android:label="你的应用名称"
        android:supportsRtl="true">
        <activity
            android:name=".MainActivity"
            android:configChanges="keyboardHidden|orientation|screenSize"  <!-- 配置变化不重启 -->
            android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
```

### 3. app/build.gradle 模板

```gradle
plugins {
    id "com.android.application"
}

android {
    namespace "com.your.package"
    compileSdk 35

    // 签名配置（可选）
    signingConfigs {
        release {
            def storeFilePath = System.getenv("ANDROID_KEYSTORE_FILE")
            if (storeFilePath) {
                storeFile file(storeFilePath)
                storePassword System.getenv("ANDROID_KEYSTORE_PASSWORD")
                keyAlias System.getenv("ANDROID_KEY_ALIAS")
                keyPassword System.getenv("ANDROID_KEY_PASSWORD")
            }
        }
    }

    defaultConfig {
        applicationId "com.your.package"
        minSdk 26  // Android 8.0
        targetSdk 35
        versionCode 1
        versionName System.getenv("ANDROID_VERSION_NAME") ?: "1.0"
    }

    buildTypes {
        release {
            // 启用签名（如果配置了）
            // signingConfig signingConfigs.release
            // 代码混淆（可选）
            // minifyEnabled true
            // proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }

    packaging {
        jniLibs {
            useLegacyPackaging true  // 使用传统 JNI 打包方式
        }
    }
}
```

### 4. 根目录 build.gradle 模板

```gradle
plugins {
    id "com.android.application" version "8.7.3" apply false
}
```

### 5. settings.gradle 模板

```gradle
pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "你的项目名称"
include ":app"
```

---

## 四、Go 后端编译为 Android 库的指南

### 1. 编译命令示例

将 Go 后端编译为 Android 支持的共享库：

```bash
# ARM64 (大多数现代 Android 设备)
GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o libbackend.so .

# ARM32 (旧设备)
GOOS=android GOARCH=arm CGO_ENABLED=1 go build -buildmode=c-shared -o libbackend.so .

# x86_64 (模拟器)
GOOS=android GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o libbackend.so .
```

注意：可能需要配置 Android NDK 进行交叉编译。

### 2. Go 后端启动参数

确保 Go 后端支持以下参数：

```go
var (
    host = flag.String("host", "127.0.0.1", "监听地址")
    port = flag.Int("port", 13335, "监听端口")
)
```

---

## 五、技术细节与最佳实践

### 1. 架构优势

| 优势 | 说明 |
|------|------|
| **代码复用** | Go 后端同时支持桌面版和安卓版 |
| **功能一致** | Web 界面在所有平台保持一致 |
| **安全性** | 服务仅监听本地回环地址，不对外暴露 |
| **性能良好** | Go 后端性能优异，进程独立运行 |
| **原生体验** | 通过 JS 桥接可调用 Android 原生功能 |

### 2. 安全注意事项

1. **网络隔离**：后端只监听 `127.0.0.1`，不对外暴露
2. **权限最小化**：只申请必要的权限（如 `INTERNET`）
3. **JS 桥接安全**：谨慎处理 `@JavascriptInterface` 方法，验证所有输入
4. **进程管理**：应用销毁时确保后端进程被杀死

### 3. 调试技巧

1. **查看后端日志**：可以修改 `ProcessBuilder` 将输出重定向到文件
2. **WebView 调试**：Chrome DevTools 可以连接调试 WebView（Chrome://inspect）
3. **Health Check**：实现专门的健康检查端点，替代检查 favicon

### 4. 扩展 JavaScript 桥接功能

在 `AndroidBridge` 类中可以添加更多方法：

```java
@JavascriptInterface
public void getDeviceInfo() {
    return android.os.Build.MODEL;
}

@JavascriptInterface
public void shareText(String text) {
    Intent shareIntent = new Intent(Intent.ACTION_SEND);
    shareIntent.setType("text/plain");
    shareIntent.putExtra(Intent.EXTRA_TEXT, text);
    startActivity(Intent.createChooser(shareIntent, "分享"));
}
```

### 5. 自定义主题样式

在 `res/values/styles.xml` 中添加：

```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="AppTheme" parent="android:Theme.Material.Light.NoActionBar">
        <item name="android:statusBarColor">@android:color/transparent</item>
        <item name="android:navigationBarColor">@android:color/black</item>
    </style>
</resources>
```

---

## 六、完整的实施检查清单

- [ ] 创建 Android 项目目录结构
- [ ] 编写 MainActivity.java（使用上面的模板）
- [ ] 配置 AndroidManifest.xml
- [ ] 配置 Gradle 构建文件
- [ ] 添加应用图标资源
- [ ] 编译 Go 后端为各架构的 .so 库
- [ ] 将 .so 库放入 jniLibs 对应目录
- [ ] 测试应用安装和启动
- [ ] 验证所有功能正常工作
- [ ] 配置签名（发布版本）

---

## 七、常见问题排查

### 问题 1：后端启动失败
- 检查 .so 库是否正确放置在 jniLibs 目录
- 验证编译架构与设备架构匹配
- 查看 logcat 输出了解具体错误

### 问题 2：WebView 无法加载页面
- 确认后端已成功启动并监听正确端口
- 检查 `isBackendReady()` 方法的健康检查端点是否正确
- 验证 `usesCleartextTraffic="true"` 已设置

### 问题 3：JS 桥接不工作
- 确保方法添加了 `@JavascriptInterface` 注解
- 验证对象名和方法名在 JS 中正确调用
- 注意线程：UI 操作必须在主线程执行

---

## 八、提示词使用说明

将上述内容保存为 markdown 文档后，下次需要创建类似项目时，直接复制以下提示词给 AI：

---

### AI 使用提示词

```
请为我创建一个 Android 混合应用项目，使用 Go 后端 + WebView 前端的架构。

要求：
1. 采用完整的项目结构
2. 使用提供的代码模板（参考 "Android混合应用开发提示词模板.md"）
3. Go 后端监听本地 127.0.0.1
4. 包含 JavaScript 桥接功能
5. 提供文件保存、选择等原生功能
6. 包含启动加载界面
7. 配置好所有 Gradle 和 Manifest 文件

项目信息：
- 应用包名：com.example.myapp
- 应用名称：MyApp
- 后端端口：13335
- Go 后端库名：libbackend.so

请提供完整的可直接构建的代码。
```

---

## 附录：项目文件快速索引

| 文件 | 说明 |
|------|------|
| [MainActivity.java](file:///workspace/app/src/main/java/com/cfdata/web/MainActivity.java) | 主活动完整实现 |
| [AndroidManifest.xml](file:///workspace/app/src/main/AndroidManifest.xml) | 应用清单配置 |
| [app/build.gradle](file:///workspace/app/build.gradle) | 应用构建配置 |
| [build.gradle](file:///workspace/build.gradle) | 根构建配置 |
| [settings.gradle](file:///workspace/settings.gradle) | Gradle 设置 |

---

**文档版本**：v1.0  
**最后更新**：2024-05-31
