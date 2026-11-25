@echo off
setlocal EnableExtensions
cd /d "%~dp0" || goto :failed

where docker.exe >nul 2>&1
if errorlevel 1 (
    echo Docker Desktop is required. Install it and run this script again.
    goto :failed
)

docker info --format "{{.ServerVersion}}" >nul 2>&1
if errorlevel 1 (
    echo Starting Docker Desktop...
    docker desktop start --timeout 180 || goto :failed
)

echo Starting FlipFeed...
docker compose -p flipfeed -f apps/docker-compose.yml -f apps/compose.windows.yaml up -d --wait || goto :failed

set "web=http://127.0.0.1:5174"
echo Checking the home page and timeline...
for %%U in ("%web%/" "%web%/api/feed-items?scene=timeline") do (
    curl.exe --noproxy "*" --fail --silent --show-error --retry 5 --retry-all-errors --retry-delay 1 --max-time 10 --output NUL "%%~U" || goto :failed
)

echo FlipFeed is ready at %web%/
start "" "%web%/"
exit /b 0

:failed
echo FlipFeed startup failed. See the error above.
pause
exit /b 1
