$env:DATABASE_URL = "postgres://postgres:santiago09@localhost:5432/restaurantos?sslmode=disable"
$env:JWT_SECRET   = "super-secreto-local"
$env:PORT         = "8080"
Set-Location "$PSScriptRoot\backend"
Write-Host "Iniciando RestaurantOS Backend en http://localhost:8080" -ForegroundColor Green
.\server.exe