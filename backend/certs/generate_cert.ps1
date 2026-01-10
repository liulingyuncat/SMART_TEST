# 生成自签名证书用于HTTPS
$certPath = ".\server.crt"
$keyPath = ".\server.key"

# 创建证书配置
$params = @{
    DnsName = "localhost", "127.0.0.1", "webtest.local"
    CertStoreLocation = "Cert:\CurrentUser\My"
    NotAfter = (Get-Date).AddYears(1)
    KeyAlgorithm = "RSA"
    KeyLength = 2048
    KeyExportPolicy = "Exportable"
}

# 生成证书
$cert = New-SelfSignedCertificate @params

Write-Host "证书已生成，Thumbprint: $($cert.Thumbprint)"
Write-Host "位置: Cert:\CurrentUser\My\$($cert.Thumbprint)"
Write-Host ""
Write-Host "证书将保存到当前目录的 server.pfx"
Write-Host "Go 服务器将直接加载系统证书存储中的证书"
Write-Host ""
Write-Host "如需导出PFX："
Write-Host '$password = ConvertTo-SecureString -String "webtest123" -Force -AsPlainText'
Write-Host "Export-PfxCertificate -Cert Cert:\CurrentUser\My\$($cert.Thumbprint) -FilePath .\server.pfx -Password `$password"
