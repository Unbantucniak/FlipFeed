@echo off
setlocal EnableExtensions
cd /d "%~dp0" || goto :failed

echo Stopping GCFeed...
docker compose -p gcfeed -f apps/docker-compose.yml -f apps/compose.windows.yaml stop || goto :failed

echo GCFeed stopped. Double-click start-gcfeed.cmd to run it again.
exit /b 0

:failed
echo Could not stop GCFeed. See the error above.
pause
exit /b 1
