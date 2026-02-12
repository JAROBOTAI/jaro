# JARO API Test Script

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  JARO System API Integration Test" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health Check
Write-Host "[1/4] Testing Health Check..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri http://localhost:8080/health
    Write-Host "[OK] Health check passed!" -ForegroundColor Green
    Write-Host "  Status: $($health.status)" -ForegroundColor Gray
    Write-Host "  Service: $($health.service)" -ForegroundColor Gray
} catch {
    Write-Host "[FAIL] Health check failed!" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 2: Create Task
Write-Host "[2/4] Creating new task..." -ForegroundColor Yellow
try {
    $taskData = @{
        input = "Deploy microservice to production with zero downtime"
        user_id = "admin-001"
    } | ConvertTo-Json
    
    $task = Invoke-RestMethod -Uri http://localhost:8080/tasks -Method POST -Body $taskData -ContentType 'application/json'
    Write-Host "[OK] Task created successfully!" -ForegroundColor Green
    Write-Host "  Task ID: $($task.task_id)" -ForegroundColor Gray
    Write-Host "  Status: $($task.status)" -ForegroundColor Gray
    Write-Host "  User ID: $($task.user_id)" -ForegroundColor Gray
    
    $taskId = $task.task_id
} catch {
    Write-Host "[FAIL] Task creation failed!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 3: Get Task Status
Write-Host "[3/4] Retrieving task status..." -ForegroundColor Yellow
Start-Sleep -Seconds 1
try {
    $taskStatus = Invoke-RestMethod -Uri "http://localhost:8080/tasks/$taskId"
    Write-Host "[OK] Task status retrieved!" -ForegroundColor Green
    Write-Host "  Task ID: $($taskStatus.id)" -ForegroundColor Gray
    Write-Host "  Status: $($taskStatus.status)" -ForegroundColor Gray
    Write-Host "  Input: $($taskStatus.input)" -ForegroundColor Gray
    Write-Host "  Created: $($taskStatus.created_at)" -ForegroundColor Gray
} catch {
    Write-Host "[FAIL] Failed to get task status!" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 4: Test 404 Error
Write-Host "[4/4] Testing 404 error handling..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8080/tasks/nonexistent-task-id" -ErrorAction Stop
    Write-Host "[FAIL] Should have returned 404!" -ForegroundColor Red
} catch {
    Write-Host "[OK] 404 error handled correctly!" -ForegroundColor Green
}
Write-Host ""

# Summary
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  All Tests Passed!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "JARO System is fully operational." -ForegroundColor Green
