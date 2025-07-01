# 安全改进测试脚本
Write-Host "======================================" -ForegroundColor Green
Write-Host "   安全改进测试脚本" -ForegroundColor Green  
Write-Host "======================================" -ForegroundColor Green
Write-Host ""

# 测试1: 未登录访问管理后台
Write-Host "[1] 测试未登录访问管理后台..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/admin/dashboard" -MaximumRedirection 0 -ErrorAction SilentlyContinue
    if ($response.StatusCode -eq 302) {
        Write-Host "✅ 正确重定向到登录页面" -ForegroundColor Green
    }
} catch {
    if ($_.Exception.Response.StatusCode -eq 302) {
        Write-Host "✅ 正确重定向到登录页面" -ForegroundColor Green
    } else {
        Write-Host "❌ 意外的响应" -ForegroundColor Red
    }
}

Write-Host ""

# 测试2: UUID文章创建
Write-Host "[2] 测试UUID文章创建..." -ForegroundColor Yellow
$headers = @{
    "X-API-Key" = "test-key-12345"
    "Content-Type" = "application/json"
}

$jsonBody = '{"title":"UUID测试文章","content":"这是一篇测试UUID的文章","status":"published"}'

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/articles" -Method POST -Headers $headers -Body $jsonBody
    if ($response.success) {
        $articleId = $response.data.id
        Write-Host "✅ 文章创建成功" -ForegroundColor Green
        Write-Host "   文章ID: $articleId" -ForegroundColor Cyan
        Write-Host "   ID长度: $($articleId.Length) 字符" -ForegroundColor Cyan
        
        if ($articleId.Length -eq 36 -and $articleId -match "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$") {
            Write-Host "✅ ID格式符合UUID标准" -ForegroundColor Green
        } else {
            Write-Host "❌ ID格式不符合UUID标准" -ForegroundColor Red
        }

        # 测试3: 通过UUID访问文章
        Write-Host ""
        Write-Host "[3] 测试通过UUID访问文章..." -ForegroundColor Yellow
        try {
            $getResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/articles/$articleId" -Headers $headers
            if ($getResponse.success) {
                Write-Host "✅ 通过UUID成功访问文章" -ForegroundColor Green
                Write-Host "   文章标题: $($getResponse.data.title)" -ForegroundColor Cyan
            } else {
                Write-Host "❌ 无法通过UUID访问文章" -ForegroundColor Red
            }
        } catch {
            Write-Host "❌ 访问文章时发生错误: $($_.Exception.Message)" -ForegroundColor Red
        }

        # 清理测试数据
        Write-Host ""
        Write-Host "[4] 清理测试数据..." -ForegroundColor Yellow
        try {
            Invoke-RestMethod -Uri "http://localhost:8080/api/articles/$articleId" -Method DELETE -Headers $headers
            Write-Host "✅ 测试数据清理完成" -ForegroundColor Green
        } catch {
            Write-Host "⚠️  清理测试数据时发生错误，可能需要手动清理" -ForegroundColor Yellow
        }

    } else {
        Write-Host "❌ 文章创建失败: $($response.error)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 创建文章时发生错误: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "======================================" -ForegroundColor Green
Write-Host "   测试完成！" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
Write-Host "现在您可以访问:" -ForegroundColor White
Write-Host "• 管理后台: http://localhost:8080/admin" -ForegroundColor Cyan
Write-Host "• 登录凭据: admin / admin123" -ForegroundColor Cyan
Write-Host ""
Write-Host "按任意键退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
