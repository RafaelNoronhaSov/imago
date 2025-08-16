Write-Host "Cleaning up previous build artifacts and data..."
Remove-Item "..\node\node_server.exe" -ErrorAction SilentlyContinue
Write-Host "Building the node server inside ../node..."

go -C ..\node build -o .\node_server.exe .

$ports = @("50052", "50053", "50054")

foreach ($port in $ports) {
    Write-Host "Starting node on port $port..."
    $env:NODE_PORT=$port
    Start-Process -FilePath "..\node\node_server.exe"
}

Write-Host "All nodes started in separate windows."
