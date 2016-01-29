Import-Module ServerManager
Add-WindowsFeature Adcs-Cert-Authority
$secpasswd = ConvertTo-SecureString "Nanocloud123+" -AsPlainText -Force
$mycreds = New-Object System.Management.Automation.PSCredential ("Administrator", $secpasswd)
Install-AdcsCertificationAuthority -CAType "EnterpriseRootCa" -Credential:$mycreds -force:$true
