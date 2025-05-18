# Simple Golang Clean Archtecture
Golang rest api with echo and mongo database

## Getting Started

#### Requirements

- Database: `PostgreSQL`
- Viper
- Gin
- Zerolog
- Dependency injection: `sarulabs`
- Redis

#### Install & Run
Clone this project:
```shell script
git clone https://github.com/HasanNugroho/gin-clean.git
```

Setup project:
```shell script
make setup
```

Build database with docker 
```shell script
make env-up
```