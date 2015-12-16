import-module RemoteDesktop
ac -Encoding UTF8  C:\Windows\system32\drivers\etc\hosts "127.0.0.1 adapps.intra.localdomain.com"
New-RDSessionDeployment -ConnectionBroker adapps.intra.localdomain.com -WebAccessServer adapps.intra.localdomain.com -SessionHost adapps.intra.localdomain.com
New-RDSessionCollection -CollectionName Collection -SessionHost adapps.intra.localdomain.com -CollectionDescription "This Collection Does stuff" -ConnectionBroker adapps.intra.localdomain.com
New-RDRemoteApp -CollectionName Collection -DisplayName hapticPowershell -FilePath 'C:\Windows\system32\WindowsPowerShell\v1.0\powershell.exe' -Alias hapticPowershell -CommandLineSetting Require -RequiredCommandLine '-ExecutionPolicy Bypass c:\publishApplication.ps1'
