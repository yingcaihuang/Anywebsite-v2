#!/bin/bash

# 数据库备份和恢复脚本

set -e

BACKUP_DIR="./backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.sql"

# 创建备份目录
mkdir -p $BACKUP_DIR

case "$1" in
    "backup")
        echo "🗃️ 创建数据库备份..."
        docker-compose exec mysql mysqldump -u app -ppassword static_hosting > $BACKUP_FILE
        echo "✅ 备份完成: $BACKUP_FILE"
        
        # 保留最近10个备份文件
        ls -t $BACKUP_DIR/backup_*.sql | tail -n +11 | xargs -r rm
        echo "📂 清理旧备份，保留最近10个文件"
        ;;
        
    "restore")
        if [ -z "$2" ]; then
            echo "❌ 请指定要恢复的备份文件"
            echo "用法: $0 restore <backup_file>"
            echo "可用备份:"
            ls -la $BACKUP_DIR/backup_*.sql 2>/dev/null || echo "无可用备份"
            exit 1
        fi
        
        if [ ! -f "$2" ]; then
            echo "❌ 备份文件不存在: $2"
            exit 1
        fi
        
        echo "⚠️ 这将删除现有数据并恢复到备份状态"
        read -p "确认继续？(y/N): " confirm
        if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
            echo "❌ 操作已取消"
            exit 1
        fi
        
        echo "🔄 恢复数据库备份..."
        docker-compose exec -T mysql mysql -u app -ppassword static_hosting < "$2"
        echo "✅ 数据库恢复完成"
        ;;
        
    "list")
        echo "📂 可用备份文件:"
        ls -la $BACKUP_DIR/backup_*.sql 2>/dev/null || echo "无可用备份"
        ;;
        
    "clean")
        echo "🧹 清理所有备份文件..."
        read -p "确认删除所有备份？(y/N): " confirm
        if [ "$confirm" == "y" ] || [ "$confirm" == "Y" ]; then
            rm -f $BACKUP_DIR/backup_*.sql
            echo "✅ 备份文件已清理"
        else
            echo "❌ 操作已取消"
        fi
        ;;
        
    *)
        echo "数据库备份和恢复工具"
        echo ""
        echo "用法:"
        echo "  $0 backup          - 创建数据库备份"
        echo "  $0 restore <file>  - 从备份文件恢复数据库"
        echo "  $0 list           - 列出所有备份文件"
        echo "  $0 clean          - 清理所有备份文件"
        echo ""
        echo "示例:"
        echo "  $0 backup"
        echo "  $0 restore ./backups/backup_20250701_120000.sql"
        ;;
esac
