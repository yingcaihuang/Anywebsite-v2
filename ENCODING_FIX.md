# Windows 中文乱码解决方案

## 问题描述
在Windows系统中运行batch文件时，中文字符显示为乱码或问号。

## 解决方案

### 方案一：修改后的启动脚本 (推荐)
使用修改后的 `start.bat` 或 `start-enhanced.bat`，这些脚本已经包含了编码修复：

```bat
@echo off
chcp 65001 >nul
REM 设置UTF-8编码
```

### 方案二：手动设置控制台编码
如果仍然出现乱码，可以在运行脚本前手动设置编码：

1. 打开命令提示符 (cmd)
2. 输入：`chcp 65001`
3. 然后运行：`start.bat`

### 方案三：使用PowerShell (推荐)
PowerShell对中文支持更好，可以直接使用：

```powershell
# 进入项目目录
cd "d:\Anywebsite-v2"

# 检查Docker
docker --version
docker-compose --version

# 创建目录
New-Item -ItemType Directory -Force -Path "static", "uploads", "certs"

# 启动服务
docker-compose up -d --build

# 等待启动
Start-Sleep -Seconds 10

# 显示结果
Write-Host "服务启动完成！" -ForegroundColor Green
Write-Host "管理后台: http://localhost:8080/admin" -ForegroundColor Cyan
Write-Host "默认账号: admin / admin123" -ForegroundColor Yellow
```

### 方案四：Windows Terminal (最佳体验)
如果安装了Windows Terminal，建议使用它运行脚本：
1. 安装 Windows Terminal (从Microsoft Store)
2. 在Windows Terminal中运行脚本
3. Windows Terminal默认支持UTF-8编码

## 字体设置
如果控制台字体不支持中文，可能需要：
1. 右键点击控制台标题栏
2. 选择"属性"
3. 在"字体"选项卡中选择支持中文的字体（如：宋体、微软雅黑）

## 系统区域设置
确保Windows系统区域设置正确：
1. 打开"控制面板" > "区域"
2. 在"管理"标签页点击"更改系统区域设置"
3. 选择"中文(简体，中国)"
4. 重启计算机

## 验证编码修复
运行以下命令验证编码是否正确：
```cmd
echo 测试中文显示
```

如果显示正常，说明编码问题已解决。

## 常见问题

### Q: 仍然显示乱码怎么办？
A: 尝试以下步骤：
1. 使用Windows Terminal
2. 检查字体设置
3. 使用PowerShell替代cmd
4. 确保文件保存为UTF-8编码

### Q: 在不同Windows版本上有差异吗？
A: 是的，Windows 10和11对UTF-8支持更好。较老版本可能需要额外设置。

### Q: 可以永久设置UTF-8编码吗？
A: 可以，在注册表中设置，但不推荐。建议在脚本中使用`chcp 65001`。

## 技术说明

修复原理：
- `chcp 65001`：设置控制台代码页为UTF-8
- `>nul`：隐藏命令输出，避免干扰
- 使用英文标签替代特殊符号，确保兼容性

修改前后对比：
```bat
# 修改前（可能乱码）
echo 🚀 启动静态网页托管服务器...

# 修改后（兼容性好）  
echo 启动静态网页托管服务器...
```
