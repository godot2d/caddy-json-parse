# PowerShell script to test the JSON parse plugin

# Start Caddy server in background (uncomment if needed)
# Start-Process -FilePath "caddy" -ArgumentList "run", "--config", "example.json" -NoNewWindow

# Wait a moment for server to start
Start-Sleep -Seconds 2

# Test JSON request with attach field containing ServerName
$jsonBody = @'
{
  "game": "com.arrow.defense3d",
  "orderId": "order-id",
  "uid": "test",
  "amount": "100",
  "tradeState": "SUCCESS",
  "timestamp": 123456789,
  "platform": "ios",
  "sku": "whatever",
  "attach": "{ \"roleId\": 100016, \"gameServerId\": 1, \"shopId\": 1, \"merchandiseId\": 1001 , \"ServerName\" : \"dev\"}"
}
'@

Write-Host "Sending test request..."
Write-Host "Request body: $jsonBody"
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080" -Method POST -Body $jsonBody -ContentType "application/json"
    Write-Host "Response:"
    Write-Host $response
} catch {
    Write-Host "Error: $($_.Exception.Message)"
    Write-Host "Make sure Caddy server is running with: caddy run --config example.json"
}

Write-Host ""
Write-Host "Test completed."