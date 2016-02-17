Invoke-WebRequest https://github.com/PowerShell/Win32-OpenSSH/releases/download/latest/OpenSSH-Win64.zip -OutFile C:\openssh.zip
$zip = "C:\openssh.zip"
$dest = "C:\Program Files"
Add-Type -assembly "system.io.compression.filesystem"
[io.compression.zipfile]::ExtractToDirectory($zip, $dest)
cd "C:\Program Files\OpenSSH-Win64"
.\ssh-keygen.exe -A
New-NetFirewallRule -Protocol TCP -LocalPort 22 -Direction Inbound -Action Allow -DisplayName SSH
.\sshd.exe install
Start-Service sshd
Set-Service sshd -StartupType Automatic