# 运行MCP服务器并捕获日志到文件
$logFile = "mcp-server.log"
taskkill /F /IM mcp-server.exe 2>&1 | Out-Null
Start-Sleep -Seconds 1
cmd /c "mcp-server.exe 2>$logFile"
