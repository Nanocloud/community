::
:: Install Active Directory

::

::
:: Install Active Directory Certification Services
::

:: Add AD CS feature

::%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Install-WindowsFeature AD-Certificate


::%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Install-AdcsCertificationAuthority -force

:: Add AD feature

%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe -ImportSystemModules Add-WindowsFeature AD-Domain-Services



%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe -ImportSystemModules Add-WindowsFeature RSAT-AD-AdminCenter



:: Configure and Promote AD

::set intradomain="intra.localdomain.com"
::%SystemRoot%\System32\dcpromo.exe /unattend /NewDomain:forest /ReplicaOrNewDomain:Domain /NewDomainDNSName:%intradomain% /DomainLevel:4 /ForestLevel:4 /SafeModeAdminPassword:"Nanocloud123+"

::%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Import-Module ADDSDeployment
::%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Install-ADDSForest -CreateDnsDelegation:$false -DatabasePath "C:\Windows\NTDS" -DomainMode "Win2012R2" -DomainName "intra.localdomain.com" -DomainNetbiosName "INTRA" -ForestMode "Win2012R2" -InstallDns:$true -LogPath "C:\Windows\NTDS" -SysvolPath "C:\Windows\SYSVOL" -Force:$true
