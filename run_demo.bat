@echo off
REM This script demonstrates how to run a cluster of distributed auth system nodes
REM and interact with them using the client on Windows.

REM Function to clean up background processes on exit
:cleanup
echo Cleaning up...
taskkill /F /PID %NODE1_PID% /PID %NODE2_PID% /PID %NODE3_PID% 2>nul
exit /b

REM Build the system
echo Building the distributed auth system...
go build -o auth-system.exe

REM Build the client
echo Building the client...
cd client
go mod tidy
go build -o auth-client.exe
cd ..

REM Check if Redis is running
echo Checking if Redis is running...
redis-cli ping > nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo Starting Redis...
    start /B redis-server
    timeout /t 2 > nul
)

REM Create data directories
if not exist data\node1 mkdir data\node1
if not exist data\node2 mkdir data\node2
if not exist data\node3 mkdir data\node3

REM Start three nodes in the background
echo Starting node 1...
start /B cmd /c auth-system.exe -node=node1 -port=:50051 -db=./data -redis=localhost:6379
for /f "tokens=2" %%a in ('tasklist /fi "imagename eq auth-system.exe" /fo list ^| find "PID:"') do set NODE1_PID=%%a
timeout /t 2 > nul

echo Starting node 2...
start /B cmd /c auth-system.exe -node=node2 -port=:50052 -db=./data -redis=localhost:6379
for /f "tokens=2" %%a in ('tasklist /fi "imagename eq auth-system.exe" /fo list ^| find "PID:"') do set NODE2_PID=%%a
timeout /t 2 > nul

echo Starting node 3...
start /B cmd /c auth-system.exe -node=node3 -port=:50053 -db=./data -redis=localhost:6379
for /f "tokens=2" %%a in ('tasklist /fi "imagename eq auth-system.exe" /fo list ^| find "PID:"') do set NODE3_PID=%%a
timeout /t 2 > nul

echo All nodes are running!
echo Node 1 PID: %NODE1_PID%
echo Node 2 PID: %NODE2_PID%
echo Node 3 PID: %NODE3_PID%

REM Demonstrate client operations
echo.
echo === Client Demo ===
echo 1. Storing user data...
client\auth-client.exe -server=localhost:50051 -op=store -key=john.doe -value=secure123

echo.
echo 2. Retrieving user data...
client\auth-client.exe -server=localhost:50052 -op=get -key=john.doe

echo.
echo 3. Authenticating user...
client\auth-client.exe -server=localhost:50053 -op=auth -key=john.doe -value=secure123

echo.
echo Demo completed! The nodes are still running.
echo Press Ctrl+C to stop all nodes and clean up.

REM Wait for user to press Ctrl+C
pause
goto cleanup
