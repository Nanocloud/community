import-module RemoteDesktop
Import-module ServerManager
Add-WindowsFeature -Name RDS-RD-Server -IncludeAllSubFeature
Add-WindowsFeature -Name RDS-Web-Access -IncludeAllSubFeature
Add-WindowsFeature -Name RDS-Connection-Broker -IncludeAllSubFeature
Install-windowsfeature RSAT-AD-AdminCenter