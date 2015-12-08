Install-windowsfeature AD-domain-services
Import-Module ADDSDeployment
$pwd=ConvertTo-SecureString 'Nanocloud123+' -asplaintext -force
Install-ADDSForest -CreateDnsDelegation:$false -DatabasePath "C:\Windows\NTDS" -DomainMode "Win2012R2" -DomainName "intra.localdomain.com" -SafeModeAdministratorPassword:$pwd -DomainNetbiosName "INTRA" -ForestMode "Win2012R2" -InstallDns:$true -LogPath "C:\Windows\NTDS" -NoRebootOnCompletion:$true -SysvolPath "C:\Windows\SYSVOL" -Force:$true