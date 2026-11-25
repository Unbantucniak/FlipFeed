@echo off
setlocal EnableExtensions
cd /d "%~dp0" || goto :failed

echo Stopping FlipFeed...
docker compose -p flipfeed -f apps/docker-compose.yml -f apps/compose.windows.yaml stop || goto :failed

echo FlipFeed stopped. Double-click start-flipfeed.cmd to run it again.
exit /b 0

:failed
echo Could not stop FlipFeed. See the error above.
pause
exit /b 1
