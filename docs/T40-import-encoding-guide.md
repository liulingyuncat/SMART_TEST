# T40缺陷管理 - CSV导入编码指南

## 如何避免导入时出现乱码

### 问题原因

CSV导入时可能出现乱码的原因：
1. **文件编码不正确**：CSV文件不是UTF-8编码
2. **Excel保存时编码问题**：Excel默认保存为ANSI编码（中文Windows上是GBK）
3. **BOM标记缺失**：文件没有UTF-8 BOM标记，导致读取器无法正确识别编码

### 后端已实现的保护措施

#### 1. 导出时添加BOM
```go
// 添加UTF-8 BOM头，解决Excel乱码问题
buf.Write([]byte{0xEF, 0xBB, 0xBF})
```
- 所有导出的CSV自动添加UTF-8 BOM
- Excel打开时会自动识别为UTF-8编码

#### 2. 导入时检测和移除BOM
```go
// 检测并移除UTF-8 BOM (EF BB BF)
if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
    data = data[3:]
}
// 清理表头（移除可能的BOM残留）
headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
```
- 自动检测UTF-8 BOM并正确处理
- 防止BOM字符混入表头字段

### 用户操作指南

#### ✅ 推荐做法（不会出现乱码）

**方法1：直接使用导出的CSV**
1. 从系统导出CSV文件
2. 在Excel中编辑（会自动识别UTF-8编码）
3. 保存时选择"UTF-8 CSV (*.csv)"
4. 重新导入

**方法2：使用文本编辑器**
1. 用VS Code、Notepad++等文本编辑器打开CSV
2. 确保编码设置为"UTF-8 with BOM"
3. 编辑后保存
4. 导入系统

**方法3：使用Google Sheets**
1. 导入CSV到Google Sheets
2. 在Google Sheets中编辑
3. 下载为CSV格式（自动UTF-8编码）
4. 导入系统

#### ❌ 避免的做法（可能出现乱码）

**不推荐的操作：**
1. ❌ Excel中"另存为"选择"CSV (逗号分隔)(*.csv)"
   - 这会使用ANSI/GBK编码，中文会乱码
   
2. ❌ 在记事本中保存为ANSI编码
   
3. ❌ 使用老版本Excel（2016之前）编辑CSV

### Excel正确保存CSV的步骤

#### Excel 2019/365/2021

1. 打开CSV文件
2. 编辑内容
3. **关键步骤**：点击"文件" → "另存为"
4. 在"保存类型"下拉菜单中选择：
   - ✅ "CSV UTF-8 (逗号分隔)(*.csv)"
   - ❌ 不要选择"CSV (逗号分隔)(*.csv)"
5. 点击"保存"

#### Excel 2016及更早版本

Excel 2016没有直接保存为UTF-8的选项，推荐使用以下方法：

**方法A：通过记事本转换**
1. Excel保存为CSV
2. 用记事本打开该CSV
3. 点击"文件" → "另存为"
4. 编码选择"UTF-8"
5. 保存

**方法B：使用Google Sheets**
1. 上传CSV到Google Sheets
2. 编辑后下载为CSV

### 验证文件编码

#### 使用Notepad++
1. 打开CSV文件
2. 查看右下角显示的编码
3. 应该显示"UTF-8 BOM"或"UTF-8"

#### 使用VS Code
1. 打开CSV文件
2. 查看右下角状态栏
3. 应该显示"UTF-8 with BOM"

#### 使用命令行
```powershell
# PowerShell中查看文件前几个字节
Get-Content -Path "defects_export.csv" -Encoding Byte -TotalCount 3
# 如果返回 239 187 191，说明有UTF-8 BOM
```

### 编码转换工具

#### PowerShell脚本转换编码
```powershell
# 将CSV转换为UTF-8 with BOM
$content = Get-Content -Path "input.csv" -Encoding Default
$utf8BOM = New-Object System.Text.UTF8Encoding $true
[System.IO.File]::WriteAllLines("output.csv", $content, $utf8BOM)
```

#### 批量转换脚本
```powershell
# 转换目录下所有CSV
Get-ChildItem -Path . -Filter *.csv | ForEach-Object {
    $content = Get-Content $_.FullName -Encoding Default
    $utf8BOM = New-Object System.Text.UTF8Encoding $true
    [System.IO.File]::WriteAllLines($_.FullName, $content, $utf8BOM)
    Write-Host "Converted: $($_.Name)"
}
```

### 导入时的错误提示

如果导入出现乱码，系统会：
1. 记录日志：`[Defect Import] UTF-8 BOM detected and removed`
2. 自动处理BOM标记
3. 清理字段中的特殊字符

如果仍然出现乱码：
- 检查原始文件编码
- 使用推荐的编辑和保存方法
- 重新从系统导出模板

### 常见编码问题对照表

| 症状 | 原因 | 解决方案 |
|------|------|----------|
| 中文显示为"ģ��" | GBK编码被当作UTF-8读取 | 转换文件为UTF-8 with BOM |
| 中文显示为"???" | 编码丢失 | 使用原始导出文件重新编辑 |
| 表头乱码 | BOM处理不当 | 后端已自动处理 |
| 部分字段乱码 | 混合编码 | 统一使用UTF-8编码 |

### 最佳实践总结

1. **导出**：系统自动添加UTF-8 BOM ✅
2. **编辑**：使用Excel 2019+或文本编辑器 ✅
3. **保存**：选择"CSV UTF-8"格式 ✅
4. **导入**：系统自动检测和处理BOM ✅

遵循以上指南，可以完全避免CSV导入时的中文乱码问题。

### 技术细节

#### UTF-8 BOM字节序列
- 十六进制：`EF BB BF`
- 十进制：`239 187 191`
- Unicode字符：`U+FEFF` (Zero Width No-Break Space)

#### 为什么需要BOM？
- Excel在Windows上默认使用系统编码（中文Windows是GBK）
- 没有BOM时，Excel无法判断文件是UTF-8编码
- 添加BOM后，Excel会自动切换到UTF-8解码

#### Go语言中的处理
```go
// 导出时添加BOM
buf.Write([]byte{0xEF, 0xBB, 0xBF})

// 导入时检测BOM
if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
    data = data[3:] // 跳过BOM
}

// 清理BOM字符
headers[0] = strings.TrimPrefix(headers[0], "\ufeff")
```

### 故障排除

如果导入后数据库中已经出现乱码：

**SQL修复脚本：**
```sql
-- 查看乱码数据
SELECT defect_id, subject, hex(subject) FROM defects 
WHERE subject LIKE '%?%' OR subject LIKE '%ģ%';

-- 批量修复为正确的模块名称
UPDATE defects 
SET subject = '模块1' 
WHERE subject LIKE '%?%' OR subject LIKE '%ģ%';
```

**手动编辑：**
1. 在前端编辑页面重新选择Subject和Phase
2. 保存后数据会自动更新为正确的UTF-8编码

### 参考资料

- [RFC 4180 - CSV格式规范](https://tools.ietf.org/html/rfc4180)
- [UTF-8 BOM说明](https://en.wikipedia.org/wiki/Byte_order_mark)
- [Excel CSV编码问题](https://support.microsoft.com/zh-cn/office)
