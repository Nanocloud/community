:: Install Active Directory
::
:: Add AD CS feature
%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Install-WindowsFeature AD-Certificate

%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe Install-AdcsCertificationAuthority -force
