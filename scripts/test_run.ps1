param(
  [switch]$ApiOnly,
  [switch]$GodotOnly
)

$RootDir = Split-Path -Parent $PSScriptRoot
$ApiDir = Join-Path $RootDir "api"
$GodotDir = Join-Path $RootDir "godot"

$env:PORT = "8080"
$env:DATABASE_URL ??= "postgres://postgres:postgres@localhost:5432/shinjuku_lunch?sslmode=disable"

# --- kill leftover process on port ---
$old = netstat -ano | Select-String "LISTENING" | Select-String ":$($env:PORT)\s"
if ($old) {
  $pid = $old -replace '.+?(\d+)$','$1'
  Write-Host "Port $env:PORT in use by PID $pid — killing..." -ForegroundColor Yellow
  Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
  Start-Sleep 1
}

# --- build & start API ---
if (-not $GodotOnly) {
  Write-Host "=== Building Go API ===" -ForegroundColor Cyan
  Push-Location $ApiDir
  go build -o server.exe . 2>&1
  if ($LASTEXITCODE -ne 0) { Write-Host "Build failed" -ForegroundColor Red; exit 1 }
  Write-Host "=== Starting API on port $env:PORT ===" -ForegroundColor Green

  $apiJob = Start-Job -ScriptBlock {
    param($d, $p)
    $env:DATABASE_URL = $d
    $env:PORT = $p
    Set-Location $using:ApiDir
    ./server.exe
  } -ArgumentList $env:DATABASE_URL, $env:PORT

  Start-Sleep 2
}

# --- start Godot ---
if (-not $ApiOnly) {
  $GodotExe = Get-ChildItem "$env:LOCALAPPDATA\Programs\Godot\Godot_v*-stable_win64.exe" |
    Select-Object -First 1 -ExpandProperty FullName
  if (-not $GodotExe -or -not (Test-Path $GodotExe)) {
    Write-Host "Godot not found — open godot/project.godot manually" -ForegroundColor Yellow
  } else {
    Write-Host "=== Starting Godot ===" -ForegroundColor Green
    Start-Process -FilePath $GodotExe -ArgumentList "--path", $GodotDir
  }
}

# --- info ---
if (-not $GodotOnly) {
  Write-Host "`nWeb → http://localhost:$env:PORT/ (redirects to /api/list)" -ForegroundColor Cyan
  Write-Host "API → http://localhost:$env:PORT/api" -ForegroundColor Cyan
  Write-Host "`nPress Ctrl+C to stop" -ForegroundColor Yellow
  Wait-Job $apiJob | Receive-Job
}
