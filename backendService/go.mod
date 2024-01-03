module github.com/ntua-el19707/DistributedSystems/backendService

go 1.21.3
replace services v0.0.0 => ./services
require services v0.0.0
replace entitys v0.0.0 => ./entitys
require entitys v0.0.0
require github.com/joho/godotenv v1.5.1 // indirect
