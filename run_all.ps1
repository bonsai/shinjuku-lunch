param(
  [switch]$NoHotReload,
  [switch]$ApiOnly,
  [switch]$FrontendOnly
)

$RootDir = Split-Path -Parent $PSCommandPath
$ApiDir = Join-Path $RootDir "api"
$FrontendDir = Join-Path $RootDir "frontend"

$env:PORT = "8080"
$env:FRONTEND_PORT = "3000"
if (-not $env:DATABASE_URL) { $env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/shinjuku_lunch?sslmode=disable" }

# --- kill leftover process on ports ---
$ports = @($env:PORT, $env:FRONTEND_PORT)
foreach ($p in $ports) {
  $old = netstat -ano | Select-String "LISTENING" | Select-String ":$($p)\s"
  if ($old) {
    $procId = ($old -split '\s+')[-1]
    if ($procId -and $procId -match '^\d+$') {
      Write-Host "Port $p in use by PID $procId — killing..."
      Stop-Process -Id $procId -Force -ErrorAction SilentlyContinue
      Start-Sleep 1
    }
  }
}

# --- build & start API ---
if (-not $FrontendOnly) {
  Write-Host "=== Starting Go API (port $env:PORT) ===" -ForegroundColor Cyan

  if ($NoHotReload) {
    Push-Location $ApiDir
    go build -o server.exe . 2>&1
    if ($LASTEXITCODE -ne 0) { Write-Host "Build failed" -ForegroundColor Red; exit 1 }
    $apiJob = Start-Job -ScriptBlock {
      param($d, $p)
      $env:DATABASE_URL = $d
      $env:PORT = $p
      Set-Location $using:ApiDir
      ./server.exe
    } -ArgumentList $env:DATABASE_URL, $env:PORT
    Pop-Location
  } else {
    $global:airExe = Get-Command "air" -ErrorAction SilentlyContinue
    if (-not $airExe) {
      Write-Host "air not found — installing..." -ForegroundColor Yellow
      go install github.com/air-verse/air@latest
      $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
    }
    $apiJob = Start-Job -ScriptBlock {
      param($d, $p)
      $env:DATABASE_URL = $d
      $env:PORT = $p
      Set-Location $using:ApiDir
      air
    } -ArgumentList $env:DATABASE_URL, $env:PORT
  }

  Start-Sleep 3
}

# --- start Next.js frontend ---
if (-not $ApiOnly) {
  Write-Host "=== Starting Next.js (port $env:FRONTEND_PORT) ===" -ForegroundColor Green
  $frontendJob = Start-Job -ScriptBlock {
    param($d, $p)
    $env:NEXT_PUBLIC_API_URL = $d
    Set-Location $using:FrontendDir
    npx next dev -p $p
  } -ArgumentList "http://localhost:$env:PORT", $env:FRONTEND_PORT
}

# --- info ---
Write-Host "`n===================================" -ForegroundColor Cyan
Write-Host "  Frontend → http://localhost:$env:FRONTEND_PORT" -ForegroundColor Cyan
Write-Host "  API       → http://localhost:$env:PORT/api" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
if (-not $FrontendOnly -and -not $NoHotReload) {
  Write-Host "Go hot-reload active — edit .go files to auto-restart" -ForegroundColor Green
}
Write-Host "Next.js hot-reload active — edit .tsx files for instant refresh" -ForegroundColor Green
Write-Host "`nPress Ctrl+C to stop" -ForegroundColor Yellow
Wait-Job $apiJob, $frontendJob | Receive-Job
