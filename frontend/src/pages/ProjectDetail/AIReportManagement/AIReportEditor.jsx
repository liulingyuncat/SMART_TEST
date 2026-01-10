import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Space, message, Spin, Empty, Modal, Dropdown, Input } from 'antd';
import { SaveOutlined, DownloadOutlined, EditOutlined, CloseOutlined } from '@ant-design/icons';
import MarkdownIt from 'markdown-it';
import MdEditor from 'react-markdown-editor-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import ChartRenderer from './ChartRenderer';
import 'react-markdown-editor-lite/lib/index.css';
import './AIReportEditor.css';

const AIReportEditor = ({ 
  report, 
  projectName,
  onSave, 
  onContentChange,
  onNameChange,
  loading = false
}) => {
  const { t } = useTranslation();
  const [isEditing, setIsEditing] = useState(false);
  const [content, setContent] = useState('');
  const [editingName, setEditingName] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [tocItems, setTocItems] = useState([]);
  const mdParser = useRef(new MarkdownIt({
    html: true,  // 允许 HTML 标签通过
    breaks: true,  // 将换行符转换为 <br>
    linkify: true  // 自动链接 URL
  }));
  const previewRef = useRef(null);

  // 初始化日志
  console.log('[AIReportEditor] Component initialized with props:', { 
    reportId: report?.id, 
    reportName: report?.name,
    projectName, 
    loading 
  });

  // 将标题转换为ID
  const titleToId = useCallback((title) => {
    return 'heading-' + title
      .toLowerCase()
      .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }, []);

  // 解析Markdown生成目录
  const generateTOC = useCallback((markdown) => {
    if (!markdown) return [];
    
    const headings = [];
    const lines = markdown.split('\n');
    
    lines.forEach((line) => {
      const match = line.match(/^(#{1,6})\s+(.+)$/);
      if (match) {
        const level = match[1].length;
        const title = match[2].trim();
        const id = titleToId(title);
        headings.push({ level, title, id });
      }
    });
    
    return headings;
  }, [titleToId]);

  // 初始化内容和名称
  useEffect(() => {
    if (report) {
      const newContent = report.content || '';
      setContent(newContent);
      setEditingName(report.name || '');
      setTocItems(generateTOC(newContent));
    }
  }, [report, generateTOC]);

  // 内容变化处理
  const handleContentChange = (e) => {
    const newContent = e.text;
    setContent(newContent);
    if (onContentChange) {
      onContentChange(newContent);
    }
  };

  // 保存
  const handleSave = async () => {
    if (!report) return;

    if (!editingName.trim()) {
      message.error(t('aiReport.reportNameRequired', { defaultValue: '请输入报告名称' }));
      return;
    }

    setIsSaving(true);
    try {
      // 保存内容
      await onSave(report.id, content);
      // 如果名称有变化，也保存名称
      if (editingName !== report.name && onNameChange) {
        await onNameChange(report.id, editingName);
      }
      message.success(t('aiReport.updateSuccess'));
      setIsEditing(false);
    } catch (error) {
      message.error(error.message || t('aiReport.updateFailed'));
    } finally {
      setIsSaving(false);
    }
  };

  // 取消编辑
  const handleCancel = () => {
    if (report) {
      setContent(report.content || '');
      setEditingName(report.name || '');
    }
    setIsEditing(false);
  };

  // 处理 Markdown 下载
  const handleDownloadMarkdown = () => {
    console.log('[AIReportEditor] Markdown download button clicked');
    downloadAsMarkdown();
  };

  // 将 chart 代码块转换为 HTML 表格（用于下载显示）
  const convertChartToTable = (chartType, config) => {
    try {
      if (!config.data || config.data.length === 0) {
        return '<p style="color: #999;">图表数据为空</p>';
      }

      const dataKey = config.dataKey || 'value';
      const headers = Object.keys(config.data[0]);
      
      let html = `<table style="border-collapse: collapse; width: 100%; margin: 16px 0;">
        <thead>
          <tr style="background-color: #f6f8fa;">`;
      
      headers.forEach(header => {
        html += `<th style="border: 1px solid #dfe2e5; padding: 8px; text-align: left; font-weight: 600;">${header}</th>`;
      });
      
      html += `</tr>
        </thead>
        <tbody>`;
      
      config.data.forEach((row, idx) => {
        html += `<tr style="background-color: ${idx % 2 === 0 ? '#fff' : '#f9f9f9'};">`;
        headers.forEach(header => {
          html += `<td style="border: 1px solid #dfe2e5; padding: 8px;">${row[header] || '-'}</td>`;
        });
        html += `</tr>`;
      });
      
      html += `</tbody>
      </table>`;
      
      return html;
    } catch (error) {
      console.error('[AIReportEditor] Chart to table conversion error:', error);
      return '<p style="color: #f00;">图表转换失败</p>';
    }
  };

  // 预处理 markdown，将 chart 代码块转换为表格
  const preprocessMarkdownForDownload = (md) => {
    let result = md;
    const chartRegex = /```chart:(\w+)\s*\n([\s\S]*?)\n```/g;
    
    result = result.replace(chartRegex, (match, chartType, configStr) => {
      try {
        const config = JSON.parse(configStr.trim());
        const title = config.title ? `<h4 style="margin: 20px 0 10px 0;">${config.title}</h4>` : '';
        const table = convertChartToTable(chartType, config);
        // 使用 HTML raw block 格式，让 markdown-it 识别为 raw HTML
        return `\n${title}\n<div style="margin: 16px 0;">${table}</div>\n`;
      } catch (e) {
        console.error('[AIReportEditor] Failed to convert chart:', e);
        return match; // 保留原始代码块
      }
    });
    
    return result;
  };

  // 生成 HTML 内容
  const generateHtmlContent = () => {
    const processedContent = preprocessMarkdownForDownload(content);
    const renderedContent = mdParser.current.render(processedContent);
    const timestamp = new Date().toLocaleString('zh-CN');
    
    // 构建完整的 HTML 文档
    let htmlContent = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${report.name}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: rgba(0, 0, 0, 0.85);
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 40px 20px;
            background-color: white;
        }
        .content {
            padding: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.35;
        }
        h1 { font-size: 32px; }
        h2 { font-size: 28px; border-bottom: 1px solid #eaecef; padding-bottom: 10px; }
        h3 { font-size: 24px; }
        h4 { font-size: 20px; }
        h5 { font-size: 16px; }
        h6 { font-size: 14px; color: #666; }
        p {
            margin-bottom: 16px;
        }
        code {
            background-color: #f6f7f9;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
            font-size: 0.85em;
        }
        pre {
            background-color: #f6f7f9;
            border: 1px solid #eaecef;
            border-radius: 3px;
            padding: 16px;
            overflow: auto;
            margin-bottom: 16px;
        }
        pre code {
            background-color: transparent;
            padding: 0;
            font-size: 14px;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: 16px;
        }
        table th, table td {
            border: 1px solid #dfe2e5;
            padding: 6px 13px;
            text-align: left;
        }
        table th {
            background-color: #f6f8fa;
            font-weight: 600;
        }
        table tr:nth-child(2n) {
            background-color: #f6f8fa;
        }
        blockquote {
            border-left: 4px solid #dfe2e5;
            padding-left: 16px;
            margin-bottom: 16px;
            color: #666;
        }
        a {
            color: #0969da;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        ul, ol {
            margin-left: 2em;
            margin-bottom: 16px;
        }
        li {
            margin-bottom: 8px;
        }
        .header {
            text-align: center;
            padding: 40px 0;
            border-bottom: 2px solid #eaecef;
            margin-bottom: 40px;
        }
        .header h1 {
            margin-top: 0;
            margin-bottom: 10px;
        }
        .metadata {
            color: #666;
            font-size: 14px;
        }
        .toc {
            background-color: #f6f8fa;
            border: 1px solid #eaecef;
            border-radius: 3px;
            padding: 16px;
            margin-bottom: 32px;
        }
        .toc-title {
            font-weight: 600;
            margin-bottom: 12px;
        }
        .toc-list {
            list-style: none;
            padding: 0;
        }
        .toc-list li {
            margin-bottom: 6px;
        }
        .toc-list a {
            color: #0969da;
            font-size: 14px;
        }
        .toc-indent-2 { margin-left: 16px; }
        .toc-indent-3 { margin-left: 32px; }
        @media (max-width: 768px) {
            .container {
                padding: 20px 10px;
            }
            h1 { font-size: 24px; }
            h2 { font-size: 20px; }
            h3 { font-size: 18px; }
        }
        @page {
            margin: 2cm;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>${report.name}</h1>
            <div class="metadata">
                <p>项目: ${projectName}</p>
                <p>生成时间: ${timestamp}</p>
            </div>
        </div>
        <div class="content">
`;
    
    // 添加渲染后的内容
    htmlContent += renderedContent;
    
    // 关闭 HTML 文档
    htmlContent += `
        </div>
    </div>
</body>
</html>`;
    
    return htmlContent;
  };

  // 下载 HTML - 完整保留SVG图表和样式
  const downloadAsHTML = () => {
    console.log('[AIReportEditor] Starting HTML download');
    try {
      if (!report) {
        console.error('[AIReportEditor] Report is not available');
        message.error('报告不可用');
        return;
      }

      if (!previewRef.current) {
        message.error('无法获取报告内容');
        return;
      }

      console.log('[AIReportEditor] Getting preview content');
      
      // 获取完整的HTML内容（包括所有部分，不仅仅是当前可见区域）
      const reportContent = previewRef.current.innerHTML;
      
      // 收集所有样式（包括Ant Design样式）
      const styles = [];
      
      // 获取所有link标签的样式
      for (const link of document.querySelectorAll('link[rel="stylesheet"]')) {
        styles.push(`<link rel="stylesheet" href="${link.href}">`);
      }
      
      // 获取所有style标签的内容
      for (const styleTag of document.querySelectorAll('style')) {
        styles.push(`<style>${styleTag.innerHTML}</style>`);
      }

      // 创建完整的HTML文档
      const htmlContent = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${report.name || '报告'}</title>
  ${styles.join('\n')}
  <style>
    body {
      margin: 20px;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
      background: #f5f5f5;
    }
    .report-container {
      background: white;
      padding: 40px;
      border-radius: 4px;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
      max-width: 1200px;
      margin: 0 auto;
    }
    .report-header {
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 2px solid #f0f0f0;
    }
    .report-title {
      font-size: 28px;
      font-weight: bold;
      color: #000;
      margin: 0 0 10px 0;
    }
    .report-meta {
      color: #666;
      font-size: 14px;
    }
    @media print {
      body { background: white; margin: 0; }
      .report-container { box-shadow: none; }
    }
  </style>
</head>
<body>
  <div class="report-container">
    <div class="report-header">
      <h1 class="report-title">${report.name || '报告'}</h1>
      <div class="report-meta">
        <p>项目: ${projectName}</p>
        <p>生成时间: ${new Date().toLocaleString('zh-CN')}</p>
      </div>
    </div>
    <div class="report-content">
      ${reportContent}
    </div>
  </div>
  
  <script>
    console.log('HTML 报告已加载');
    // 支持打印功能
    document.addEventListener('keydown', (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
        window.print();
      }
    });
  </script>
</body>
</html>`;

      console.log('[AIReportEditor] HTML content prepared, size:', htmlContent.length);

      // 创建Blob并下载
      const blob = new Blob([htmlContent], { type: 'text/html;charset=utf-8' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      const dateStr = new Date().toISOString().split('T')[0].replace(/-/g, '');
      const filename = `${projectName}-${report.name}-${dateStr}.html`;
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);

      message.success('✅ HTML 报告下载成功（支持打印为PDF）');
      console.log('[AIReportEditor] HTML file downloaded:', filename);

    } catch (error) {
      console.error('[AIReportEditor] HTML download failed:', error);
      message.error('HTML 下载失败: ' + (error.message || '未知错误'));
    }
  };

  // 处理 HTML 下载
  const handleDownloadHTML = () => {
    console.log('[AIReportEditor] HTML download button clicked');
    downloadAsHTML();
  };

  // 下载 Markdown
  const downloadAsMarkdown = () => {
    console.log('[AIReportEditor] Starting Markdown download');
    try {
      if (!report) {
        console.error('[AIReportEditor] Report is not available');
        message.error('报告不可用');
        return;
      }
      
      console.log('[AIReportEditor] Content length:', content.length);
      
      const now = new Date();
      const dateStr = now.toISOString().split('T')[0].replace(/-/g, '');
      const filename = `${projectName}-${report.name}-${dateStr}.md`;
      console.log('[AIReportEditor] Download filename:', filename);
      
      const element = document.createElement('a');
      element.setAttribute('href', 'data:text/markdown;charset=utf-8,' + encodeURIComponent(content));
      element.setAttribute('download', filename);
      element.style.display = 'none';
      console.log('[AIReportEditor] Download element created');
      
      document.body.appendChild(element);
      console.log('[AIReportEditor] Element appended to body');
      
      element.click();
      console.log('[AIReportEditor] Click triggered, download should start');
      
      document.body.removeChild(element);
      console.log('[AIReportEditor] Element removed from body');
      
      message.success('Markdown 报告下载成功');
      console.log('[AIReportEditor] Markdown download completed successfully');
    } catch (error) {
      console.error('[AIReportEditor] Markdown download failed:', error);
      message.error('Markdown 下载失败: ' + error.message);
      throw error;
    }
  };

  if (loading) {
    return (
      <div className="ai-report-editor">
        <Spin />
      </div>
    );
  }

  if (!report) {
    return (
      <div className="ai-report-editor">
        <Empty description={t('aiReport.selectOrCreateReport')} />
      </div>
    );
  }

  return (
    <div className="ai-report-editor">
      <div className="editor-toolbar">
        {isEditing ? (
          <Input
            value={editingName}
            onChange={(e) => setEditingName(e.target.value)}
            className="report-title-input"
            placeholder={t('aiReport.enterReportName', { defaultValue: '请输入报告名称' })}
          />
        ) : (
          <div className="report-title">{report.name}</div>
        )}
        <Space>
          {isEditing ? (
            <>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={handleSave}
                loading={isSaving}
              >
                {t('aiReport.save')}
              </Button>
              <Button icon={<CloseOutlined />} onClick={handleCancel}>
                {t('aiReport.cancel')}
              </Button>
            </>
          ) : (
            <>
              <Button icon={<EditOutlined />} onClick={() => setIsEditing(true)}>
                {t('aiReport.edit')}
              </Button>
              <Dropdown
                menu={{
                  items: [
                    {
                      key: 'html',
                      label: 'HTML 格式',
                      onClick: () => handleDownloadHTML()
                    },
                    {
                      key: 'markdown',
                      label: 'Markdown 格式',
                      onClick: () => handleDownloadMarkdown()
                    }
                  ]
                }}
              >
                <Button icon={<DownloadOutlined />}>
                  {t('aiReport.download')}
                </Button>
              </Dropdown>
            </>
          )}
        </Space>
      </div>

      {/* 下载格式选择已集成到 handleDownload 函数中 */}

      <div className="editor-content">
        {isEditing ? (
          <MdEditor
            value={content}
            style={{ height: '100%' }}
            renderHTML={(text) => (
              <ReactMarkdown
                children={text}
                remarkPlugins={[remarkGfm]}
              />
            )}
            onChange={handleContentChange}
            config={{
              view: {
                menu: true,
                md: true,
                html: false,
              },
              table: {
                maxRow: 5,
                maxCol: 6,
              },
            }}
          />
        ) : (
          <div className="readonly-container">
            {/* 左侧目录 */}
            {tocItems.length > 0 && (
              <div className="toc-sidebar">
                <div className="toc-title">目录</div>
                <div className="toc-list">
                  {tocItems.map((item, index) => (
                    <div
                      key={index}
                      className={`toc-item toc-level-${item.level}`}
                      onClick={() => {
                        console.log('Clicking TOC item:', item.title, 'ID:', item.id);
                        const element = document.getElementById(item.id);
                        console.log('Found element:', element);
                        if (element) {
                          element.scrollIntoView({ behavior: 'smooth', block: 'start' });
                        } else {
                          console.warn('Element not found with ID:', item.id);
                        }
                      }}
                    >
                      {item.title}
                    </div>
                  ))}
                </div>
              </div>
            )}
            
            {/* 右侧内容 */}
            <div className="markdown-preview" ref={previewRef}>
              {content ? (
                <ReactMarkdown
                  children={content}
                  remarkPlugins={[remarkGfm]}
                  allowHtml={true}
                  components={{
                    // SVG 标签支持
                    svg: ({ node, ...props }) => {
                      return <svg {...props} style={{ maxWidth: '100%', height: 'auto', display: 'block', margin: '16px 0', border: '1px solid #eaecef', borderRadius: '3px', padding: '12px', backgroundColor: '#fafbfc' }} />;
                    },
                    // 代码块支持 - 检测并渲染 SVG 和 Recharts 图表
                    code({ node, inline, className, children, ...props }) {
                      const language = className ? className.replace(/language-/, '') : '';
                      const codeString = String(children).trim();
                      
                      console.log('[AIReportEditor] Code block detected:', { 
                        language, 
                        codeLength: codeString.length, 
                        inline 
                      });
                      
                      // 检测是否是 Recharts 图表 (chart:type 格式)
                      if (language.startsWith('chart:')) {
                        const chartType = language.replace('chart:', '');
                        console.log('[AIReportEditor] Chart detected:', chartType);
                        try {
                          console.log('[AIReportEditor] Parsing JSON config for chart:', chartType);
                          const config = JSON.parse(codeString);
                          console.log('[AIReportEditor] Chart config parsed successfully:', { 
                            type: chartType, 
                            dataPoints: config.data?.length 
                          });
                          return <ChartRenderer type={chartType} config={config} />;
                        } catch (error) {
                          console.error('[AIReportEditor] Chart parsing error:', error, 'Code:', codeString.substring(0, 100));
                          return (
                            <div style={{ padding: '16px', margin: '16px 0', backgroundColor: '#fff7e6', border: '1px solid #ffc069', borderRadius: '3px', color: '#ad6800' }}>
                              <strong>图表配置解析失败:</strong> {error.message}
                              <div style={{ marginTop: '8px', fontSize: '12px', fontFamily: 'monospace', maxHeight: '200px', overflow: 'auto' }}>
                                {codeString.substring(0, 500)}
                              </div>
                            </div>
                          );
                        }
                      }
                      
                      // 检测是否是 SVG 代码
                      if ((language === 'html' || language === 'svg' || language === 'xml') && codeString.includes('<svg')) {
                        console.log('[AIReportEditor] SVG code block detected');
                        try {
                          return (
                            <div
                              style={{
                                maxWidth: '100%',
                                height: 'auto',
                                display: 'block',
                                margin: '16px 0',
                                border: '1px solid #eaecef',
                                borderRadius: '3px',
                                padding: '12px',
                                backgroundColor: '#fafbfc',
                                overflow: 'auto'
                              }}
                              dangerouslySetInnerHTML={{ __html: codeString }}
                            />
                          );
                        } catch (error) {
                          console.error('[AIReportEditor] SVG rendering error:', error);
                          // 降级处理：显示为普通代码块
                        }
                      }
                      
                      // 普通代码块处理
                      return (
                        <code className={className} {...props}>
                          {children}
                        </code>
                      );
                    },
                    h1: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h1 id={id} {...props}>{children}</h1>;
                    },
                    h2: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h2 id={id} {...props}>{children}</h2>;
                    },
                    h3: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h3 id={id} {...props}>{children}</h3>;
                    },
                    h4: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h4 id={id} {...props}>{children}</h4>;
                    },
                    h5: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h5 id={id} {...props}>{children}</h5>;
                    },
                    h6: ({ node, children, ...props }) => {
                      const extractText = (ch) => {
                        if (typeof ch === 'string') return ch;
                        if (Array.isArray(ch)) return ch.map(extractText).join('');
                        if (ch?.props?.children) return extractText(ch.props.children);
                        return '';
                      };
                      const id = titleToId(extractText(children));
                      return <h6 id={id} {...props}>{children}</h6>;
                    },
                  }}
                />
              ) : (
                <Empty description={t('aiReport.emptyContent')} />
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default AIReportEditor;
