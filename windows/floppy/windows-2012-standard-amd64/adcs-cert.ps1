Import-Module ServerManager
Add-WindowsFeature Adcs-Cert-Authority
$secpasswd = ConvertTo-SecureString "Nanocloud123+" -AsPlainText -Force
$mycreds = New-Object System.Management.Automation.PSCredential ("Administrator", $secpasswd)
Install-AdcsCertificationAuthority -CAType "EnterpriseRootCa" -Credential:$mycreds -force:$true
cmd.exe /c "certutil -ca.cert ad.cer"
cmd.exe /c "certutil -encode ad.cer ad2012.cer"