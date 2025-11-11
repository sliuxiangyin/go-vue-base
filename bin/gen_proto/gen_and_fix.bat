@echo off
REM 生成 Proto 文件并自动修复 Python gRPC 导入

echo ==========================================
echo 生成 Protocol Buffer 文件
echo ==========================================
echo.

cd /d "%~dp0"

buf generate

if %ERRORLEVEL% EQU 0 (
    echo.
    echo √ Proto 文件生成成功
    echo.
    echo ==========================================
    echo 修复 Python gRPC 导入语句
    echo ==========================================
    echo.
    
    python ..\..\python\learn_en\fix_grpc_imports.py ..\..\python\learn_en\proto
    
    echo.
    echo ==========================================
    echo 全部完成!
    echo ==========================================
) else (
    echo.
    echo × Proto 文件生成失败
    exit /b 1
)

pause
