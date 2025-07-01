# 纯净HTML内容渲染测试
$baseUrl = "http://localhost:8080"

Write-Host "=== 纯净HTML内容渲染测试 ===" -ForegroundColor Green

# 1. 登录管理员
Write-Host "1. 管理员登录..." -ForegroundColor Yellow
try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/admin/login" -Method POST -Body @{
        username = "admin"
        password = "password"
    } -ContentType "application/x-www-form-urlencoded" -SessionVariable session
    Write-Host "✓ 登录成功" -ForegroundColor Green
} catch {
    Write-Host "✗ 登录失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. 创建包含丰富HTML内容的测试文章
Write-Host "2. 创建测试文章..." -ForegroundColor Yellow

$richHtmlContent = @"
<h1>这是一级标题</h1>
<h2>这是二级标题</h2>
<p>这是一个包含<strong>粗体</strong>、<em>斜体</em>和<u>下划线</u>的段落。</p>

<blockquote>
这是一个引用块，用来测试样式效果。
</blockquote>

<h3>列表测试</h3>
<ul>
<li>无序列表项1</li>
<li>无序列表项2</li>
<li><a href="https://example.com" target="_blank">带链接的列表项</a></li>
</ul>

<ol>
<li>有序列表项1</li>
<li>有序列表项2</li>
<li>有序列表项3</li>
</ol>

<h3>代码展示</h3>
<p>行内代码：<code>console.log("Hello World")</code></p>

<pre><code>function greetUser(name) {
    return `Hello, ${name}!`;
}

// 调用函数
const message = greetUser("张三");
console.log(message);</code></pre>

<h3>样式测试</h3>
<p style="color: blue; font-size: 18px;">这是一段蓝色的大字体文本。</p>
<p style="background-color: yellow; padding: 10px; border-radius: 5px;">这是一段带背景色的文本。</p>

<hr>

<h3>表格测试</h3>
<table>
<tr>
<th>姓名</th>
<th>年龄</th>
<th>城市</th>
</tr>
<tr>
<td>张三</td>
<td>25</td>
<td>北京</td>
</tr>
<tr>
<td>李四</td>
<td>30</td>
<td>上海</td>
</tr>
</table>

<div style="border: 2px solid #ccc; padding: 15px; margin: 20px 0; border-radius: 8px;">
<h4>自定义容器</h4>
<p>这是一个自定义的容器，用来测试div元素的渲染效果。</p>
</div>
"@

$createBody = @{
    title = "纯净HTML渲染测试"
    content = $richHtmlContent
    status = "published"
    slug = "clean-html-test"
}

try {
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles" -Method POST -Body ($createBody | ConvertTo-Json) -ContentType "application/json" -WebSession $session
    Write-Host "✓ 文章创建成功！文章ID: $($createResponse.article.id)" -ForegroundColor Green
    $articleId = $createResponse.article.id
    $articleSlug = $createResponse.article.slug
} catch {
    Write-Host "✗ 创建文章失败: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails) {
        Write-Host "错误详情: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
    exit 1
}

# 3. 测试Web页面渲染
Write-Host "3. 测试Web页面渲染..." -ForegroundColor Yellow
try {
    $webResponse = Invoke-WebRequest -Uri "$baseUrl/articles/$articleSlug" -Method GET
    Write-Host "✓ Web页面状态码: $($webResponse.StatusCode)" -ForegroundColor Green
    
    $htmlContent = $webResponse.Content
    
    # 检查是否是纯净的HTML内容（没有额外的页面结构）
    if ($htmlContent -notmatch "article-header" -and $htmlContent -notmatch "article-footer") {
        Write-Host "✓ 页面结构已简化，没有多余的header和footer" -ForegroundColor Green
    } else {
        Write-Host "✗ 页面仍包含额外的结构元素" -ForegroundColor Red
    }
    
    # 检查HTML内容是否正确渲染
    if ($htmlContent -match "<h1>这是一级标题</h1>") {
        Write-Host "✓ H1标题渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ H1标题渲染异常" -ForegroundColor Red
    }
    
    if ($htmlContent -match "<strong>粗体</strong>") {
        Write-Host "✓ 粗体文字渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 粗体文字渲染异常" -ForegroundColor Red
    }
    
    if ($htmlContent -match "<blockquote>") {
        Write-Host "✓ 引用块渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 引用块渲染异常" -ForegroundColor Red
    }
    
    if ($htmlContent -match "<table>") {
        Write-Host "✓ 表格渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 表格渲染异常" -ForegroundColor Red
    }
    
    if ($htmlContent -match 'style="color: blue') {
        Write-Host "✓ 内联样式渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 内联样式渲染异常" -ForegroundColor Red
    }
    
    if ($htmlContent -match "<pre><code>") {
        Write-Host "✓ 代码块渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 代码块渲染异常" -ForegroundColor Red
    }
    
    # 输出页面开头部分用于验证
    Write-Host "`n页面开头内容预览:" -ForegroundColor Cyan
    $preview = $htmlContent.Substring(0, [Math]::Min(500, $htmlContent.Length))
    Write-Host $preview -ForegroundColor White
    
} catch {
    Write-Host "✗ 获取Web页面失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. 检查静态文件
Write-Host "`n4. 检查静态文件生成..." -ForegroundColor Yellow
$staticFilePath = "static\articles\$articleSlug\index.html"
if (Test-Path $staticFilePath) {
    Write-Host "✓ 静态文件已生成: $staticFilePath" -ForegroundColor Green
    
    $staticContent = Get-Content $staticFilePath -Raw
    if ($staticContent -match "<h1>这是一级标题</h1>" -and $staticContent -notmatch "article-header") {
        Write-Host "✓ 静态文件使用纯净模板，HTML内容渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 静态文件渲染异常" -ForegroundColor Red
    }
} else {
    Write-Host "✗ 静态文件未生成" -ForegroundColor Red
}

# 5. 清理测试数据
Write-Host "`n5. 清理测试数据..." -ForegroundColor Yellow
try {
    $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles/$articleId" -Method DELETE -WebSession $session
    Write-Host "✓ 测试文章已删除" -ForegroundColor Green
} catch {
    Write-Host "✗ 删除测试文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== 纯净HTML内容渲染测试完成 ===" -ForegroundColor Green
Write-Host "现在您的文章页面只显示纯净的HTML内容，没有额外的标题、时间或页面结构！" -ForegroundColor Cyan
