:: Install RDS
:: http://www.virtualizationadmin.com/articles-tutorials/vdi-articles/general/using-powershell-control-rds-windows-server-2012.html
:: Add AD feature
:: %SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe -ImportSystemModules Add-WindowsFeature AD-Domain-Services
%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Import-Module RemoteDesktop
%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe New-RDSessionDeployment -ConnectionBroker "adapps.intra.nanocloud.com" -WebAccessServer "adapps.intra.nanocloud.com" -SessionHost "adapps.intra.nanocloud.com"
::%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe New-RDSessionDeployment -ConnectionBroker "adapps.intra.localdomain.com" -WebAccessServer "adapps.intra.localdomain.com" -SessionHost "adapps.intra.localdomain.com"
