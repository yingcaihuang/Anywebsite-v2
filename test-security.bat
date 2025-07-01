@echo off
chcp 65001 >nul
echo ======================================
echo   安全改进测试脚本
echo ======================================
echo.

echo [1] 测试未登录访问管理后台（应该重定向到登录页）
curl -s -I http://localhost:8080/admin/dashboard | findstr "Location"
echo.

echo [2] 测试登录功能
echo 正在测试登录...
curl -s -c cookies.txt -d "username=admin&password=admin123" http://localhost:8080/admin/login > nul
echo 登录完成，检查重定向...
curl -s -I -b cookies.txt http://localhost:8080/admin/dashboard | findstr "200 OK"
echo.

echo [3] 测试UUID文章创建（通过API）
echo 创建测试文章...
for /f %%i in ('curl -s -X POST http://localhost:8080/api/articles -H "X-API-Key: test-key-12345" -H "Content-Type: application/json" -d "{\"title\":\"UUID测试文章\",\"content\":\"这是一篇测试UUID的文章\",\"status\":\"published\"}" ^| jq -r ".data.id"') do set ARTICLE_ID=%%i

echo 创建的文章ID: %ARTICLE_ID%
echo ID长度: 
echo %ARTICLE_ID% | powershell -Command "$input | Measure-Object -Character | Select-Object -ExpandProperty Characters"

echo.
echo [4] 测试UUID文章访问
curl -s "http://localhost:8080/api/articles/%ARTICLE_ID%" -H "X-API-Key: test-key-12345" | jq ".data.title"

echo.
echo [5] 检查数据库中的UUID格式
docker exec static-hosting-mysql mysql -u app -ppassword static_hosting -e "SELECT id, title FROM articles LIMIT 1;" 2>nul

echo.
echo [6] 清理测试数据
curl -s -X DELETE "http://localhost:8080/api/articles/%ARTICLE_ID%" -H "X-API-Key: test-key-12345" > nul
del cookies.txt > nul 2>&1

echo.
echo ======================================
echo   测试完成！
echo ======================================
echo.
echo 测试结果说明：
echo - 未登录访问应显示 "Location: /admin/login"
echo - 登录后访问应显示 "200 OK"  
echo - 文章ID应为36字符长度的UUID格式
echo - 能够通过UUID正常访问文章
echo.
pause
