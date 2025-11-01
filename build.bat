@echo off
echo Building frontend...
cd frontend
call npm install
call npm run build
cd ..

echo Building backend...
go build -o databaseAi.exe .

echo Build completed!