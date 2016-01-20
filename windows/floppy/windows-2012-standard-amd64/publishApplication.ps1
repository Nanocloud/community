# Poweshell script to publish a new application

Function Get-FileName($initialDirectory)
{
    [System.Reflection.Assembly]::LoadWithPartialName("System.windows.forms") | Out-Null

    $OpenFileDialog = New-Object System.Windows.Forms.OpenFileDialog
    $OpenFileDialog.initialDirectory = $initialDirectory
    $OpenFileDialog.filter = "All files (*.*)| *.*"
    $OpenFileDialog.ShowDialog() | Out-Null
    $OpenFileDialog.filename
} #end function Get-FileName

# First import the right module (does nothing if the module is already loaded)
Import-Module C:\windows\system32\windowspowershell\v1.0\Modules\RemoteDesktop\RemoteDesktop.psd1

# *** Get our publish app filename ***
if ($args[0]) { $filename = $args[0] } else { $filename = Get-FileName -initialDirectory "C:\cygwin64\home\Administrator\" }

# *** Find out the collection name ***
$apps = Get-RDRemoteApp
$collectionName = $apps[0].CollectionName
$displayName = [io.path]::GetFileNameWithoutExtension($filename)

# *** Publish the application ***
New-RDRemoteApp -CollectionName $collectionName -DisplayName $displayName -FilePath $filename
