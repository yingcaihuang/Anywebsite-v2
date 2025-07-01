# HTML渲染功能测试脚本
$baseUrl = "http://localhost:8080"

Write-Host "=== HTML内容渲染测试 ===" -ForegroundColor Green

# 1. 首先登录管理员
Write-Host "1. 管理员登录..." -ForegroundColor Yellow
$loginResponse = Invoke-RestMethod -Uri "$baseUrl/admin/login" -Method POST -Body @{
    username = "admin"
    password = "password"
} -ContentType "application/x-www-form-urlencoded" -SessionVariable session

Write-Host "登录响应: $loginResponse" -ForegroundColor Cyan

# 2. 创建包含HTML内容的文章
Write-Host "2. 创建包含HTML内容的测试文章..." -ForegroundColor Yellow

$htmlContent = @"
<h2>这是一个测试标题</h2>
<p>这是一个包含<strong>粗体文字</strong>和<em>斜体文字</em>的段落。</p>
<blockquote>
这是一个引用块，用来测试HTML渲染。
</blockquote>
<ul>
<li>这是列表项1</li>
<li>这是列表项2</li>
<li>这是<a href="https://example.com">带链接的列表项</a></li>
</ul>
<p>下面是一段代码:</p>
<pre><code>function hello() {
    console.log("Hello, World!");
}</code></pre>
<p style="color: blue;">这是一个带样式的段落。</p>
"@

$createBody = @{
    title = "HTML渲染测试文章"
    content = $htmlContent
    status = "published"
    slug = "html-rendering-test"
}

try {
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles" -Method POST -Body ($createBody | ConvertTo-Json) -ContentType "application/json" -WebSession $session
    Write-Host "文章创建成功！文章ID: $($createResponse.article.id)" -ForegroundColor Green
    $articleId = $createResponse.article.id
    $articleSlug = $createResponse.article.slug
} catch {
    Write-Host "创建文章失败: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "响应内容: $($_.Exception.Response)" -ForegroundColor Red
    exit 1
}

# 3. 通过API获取文章内容
Write-Host "3. 通过API获取文章内容..." -ForegroundColor Yellow
try {
    $apiResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles/$articleId" -Method GET -WebSession $session
    Write-Host "API返回的文章内容:" -ForegroundColor Cyan
    Write-Host $apiResponse.article.content -ForegroundColor White
} catch {
    Write-Host "获取文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. 通过Web页面获取渲染后的HTML
Write-Host "4. 测试Web页面HTML渲染..." -ForegroundColor Yellow
try {
    $webResponse = Invoke-WebRequest -Uri "$baseUrl/articles/$articleSlug" -Method GET
    Write-Host "Web页面状态码: $($webResponse.StatusCode)" -ForegroundColor Cyan
    
    # 检查HTML内容是否正确渲染
    $htmlBody = $webResponse.Content
    
    if ($htmlBody -match "<h2>这是一个测试标题</h2>") {
        Write-Host "✓ H2标题渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ H2标题渲染异常" -ForegroundColor Red
    }
    
    if ($htmlBody -match "<strong>粗体文字</strong>") {
        Write-Host "✓ 粗体文字渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 粗体文字渲染异常" -ForegroundColor Red
    }
    
    if ($htmlBody -match "<blockquote>") {
        Write-Host "✓ 引用块渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 引用块渲染异常" -ForegroundColor Red
    }
    
    if ($htmlBody -match "<ul>") {
        Write-Host "✓ 列表渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 列表渲染异常" -ForegroundColor Red
    }
    
    if ($htmlBody -match "<pre><code>") {
        Write-Host "✓ 代码块渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 代码块渲染异常" -ForegroundColor Red
    }
    
    if ($htmlBody -match 'style="color: blue;"') {
        Write-Host "✓ 内联样式渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 内联样式渲染异常" -ForegroundColor Red
    }
    
} catch {
    Write-Host "获取Web页面失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. 测试文章更新功能
Write-Host "5. 测试文章更新功能..." -ForegroundColor Yellow

$updatedContent = @"
<h2>这是更新后的标题</h2>
<p>这是<strong>更新后的内容</strong>，包含<em>新的HTML标签</em>。</p>
<div style="background: yellow; padding: 10px;">
<p>这是一个带样式的div容器。</p>
</div>
"@

$updateBody = @{
    title = "更新后的HTML渲染测试文章"
    content = $updatedContent
    status = "published"
}

try {
    $updateResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles/$articleId" -Method PUT -Body ($updateBody | ConvertTo-Json) -ContentType "application/json" -WebSession $session
    Write-Host "文章更新成功！" -ForegroundColor Green
    
    # 验证更新后的内容
    Start-Sleep -Seconds 1
    $updatedApiResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles/$articleId" -Method GET -WebSession $session
    if ($updatedApiResponse.article.content -match "更新后的标题") {
        Write-Host "✓ 文章内容更新成功" -ForegroundColor Green
    } else {
        Write-Host "✗ 文章内容更新失败" -ForegroundColor Red
    }
    
} catch {
    Write-Host "更新文章失败: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "错误详情: $($_.ErrorDetails.Message)" -ForegroundColor Red
}

# 6. 检查静态文件生成
Write-Host "6. 检查静态文件生成..." -ForegroundColor Yellow
$staticFilePath = "static\articles\$articleSlug\index.html"
if (Test-Path $staticFilePath) {
    Write-Host "✓ 静态文件已生成: $staticFilePath" -ForegroundColor Green
    
    $staticContent = Get-Content $staticFilePath -Raw
    if ($staticContent -match "更新后的标题") {
        Write-Host "✓ 静态文件HTML内容渲染正确" -ForegroundColor Green
    } else {
        Write-Host "✗ 静态文件HTML内容渲染异常" -ForegroundColor Red
    }
} else {
    Write-Host "✗ 静态文件未生成" -ForegroundColor Red
}

# 7. 清理测试数据
Write-Host "7. 清理测试数据..." -ForegroundColor Yellow
try {
    $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/api/articles/$articleId" -Method DELETE -WebSession $session
    Write-Host "✓ 测试文章已删除" -ForegroundColor Green
} catch {
    Write-Host "删除测试文章失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "=== HTML渲染测试完成 ===" -ForegroundColor Green
