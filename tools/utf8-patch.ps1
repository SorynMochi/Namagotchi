$Utf8 = [System.Text.Encoding]::UTF8
$Utf8NoBom = New-Object System.Text.UTF8Encoding($false)

function Read-Utf8Text {
  param([Parameter(Mandatory = $true)][string]$Path)

  $fullPath = [System.IO.Path]::GetFullPath((Join-Path (Get-Location) $Path))
  return [System.IO.File]::ReadAllText($fullPath, $Utf8)
}

function Write-Utf8NoBomText {
  param(
    [Parameter(Mandatory = $true)][string]$Path,
    [Parameter(Mandatory = $true)][string]$Content
  )

  $fullPath = [System.IO.Path]::GetFullPath((Join-Path (Get-Location) $Path))
  $folder = Split-Path $fullPath -Parent

  if ($folder) {
    New-Item -ItemType Directory -Force -Path $folder | Out-Null
  }

  [System.IO.File]::WriteAllText($fullPath, $Content, $Utf8NoBom)
}