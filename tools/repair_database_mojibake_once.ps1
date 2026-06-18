$ErrorActionPreference = "Stop"

try {
  [System.Text.Encoding]::RegisterProvider([System.Text.CodePagesEncodingProvider]::Instance)
} catch {
}

$Utf8 = [System.Text.Encoding]::UTF8
$Utf8NoBom = New-Object System.Text.UTF8Encoding($false)
$Win1252 = [System.Text.Encoding]::GetEncoding(1252)

$DbPath = ".\internal\database\database.go"
$FullPath = [System.IO.Path]::GetFullPath((Join-Path (Get-Location) $DbPath))

if (-not [System.IO.File]::Exists($FullPath)) {
  throw "database.go was not found at $FullPath"
}

$Text = [System.IO.File]::ReadAllText($FullPath, $Utf8)

if ($null -eq $Text) {
  throw "database.go could not be read."
}

if ($Text.Length -eq 0) {
  throw "database.go is empty. Refusing to write anything."
}

$BackupPath = "$FullPath.mojibake-backup"
[System.IO.File]::WriteAllText($BackupPath, $Text, $Utf8NoBom)

$MarkerChars = @(
  [char]0x00C3, # Ã
  [char]0x00C2, # Â
  [char]0x00E2, # â
  [char]0x00F0  # ð
)

function HasMarker($Value) {
  if ($null -eq $Value) {
    return $false
  }

  foreach ($marker in $MarkerChars) {
    if ($Value.IndexOf($marker) -ge 0) {
      return $true
    }
  }

  return $false
}

function RepairLine($Value) {
  if ($null -eq $Value) {
    return $Value
  }

  $current = [string]$Value

  for ($i = 0; $i -lt 8; $i++) {
    if (-not (HasMarker $current)) {
      break
    }

    try {
      $candidate = $Utf8.GetString($Win1252.GetBytes($current))
    } catch {
      break
    }

    if ($candidate -eq $current) {
      break
    }

    $current = $candidate
  }

  return $current
}

$Lines = $Text -split "`r?`n", -1

Write-Host ""
Write-Host "Scanning database.go for mojibake-looking markers..."

$BeforeCount = 0
for ($i = 0; $i -lt $Lines.Count; $i++) {
  if (HasMarker $Lines[$i]) {
    $BeforeCount++
    Write-Host ("Before line {0}: {1}" -f ($i + 1), $Lines[$i])
  }
}

$FixedLines = New-Object System.Collections.Generic.List[string]

foreach ($line in $Lines) {
  if (HasMarker $line) {
    [void]$FixedLines.Add((RepairLine $line))
  } else {
    [void]$FixedLines.Add($line)
  }
}

$Fixed = [string]::Join("`r`n", $FixedLines)

# Keep restored punctuation ASCII-safe in Go source.
$Fixed = $Fixed.Replace(([string][char]0x2019), "'")
$Fixed = $Fixed.Replace(([string][char]0x2018), "'")
$Fixed = $Fixed.Replace(([string][char]0x201C), '"')
$Fixed = $Fixed.Replace(([string][char]0x201D), '"')
$Fixed = $Fixed.Replace(([string][char]0x2026), "...")

[System.IO.File]::WriteAllText($FullPath, $Fixed, $Utf8NoBom)

gofmt -w $DbPath

$After = [System.IO.File]::ReadAllText($FullPath, $Utf8)
$AfterLines = $After -split "`r?`n", -1

$AfterCount = 0
for ($i = 0; $i -lt $AfterLines.Count; $i++) {
  if (HasMarker $AfterLines[$i]) {
    $AfterCount++
    Write-Host ("After line {0}: {1}" -f ($i + 1), $AfterLines[$i])
  }
}

Write-Host ""
Write-Host ("Mojibake-looking lines before repair: {0}" -f $BeforeCount)
Write-Host ("Mojibake-looking lines after repair:  {0}" -f $AfterCount)
Write-Host ("Backup written to: {0}" -f $BackupPath)

if ($AfterCount -gt 0) {
  throw "database.go still has mojibake-looking markers. Review the printed 'After line' entries."
}

go test ./...
