# Sistema de Gestión de Candidatos

Este proyecto expone un conjunto de APIs (CRUD) para gestionar candidatos en un proceso de selección y contratación. Incluye autenticación JWT, migraciones con Flyway y documentación Swagger.

## Índice

1. [Tecnologías / Dependencias](#tecnologías--dependencias)  
2. [Estructura de Directorios](#estructura-de-directorios)  
3. [Configuración del Proyecto](#configuración-del-proyecto)  
    - [Inicializar el Módulo Go](#1-inicializar-el-módulo-go)  
    - [Instalar Dependencias](#2-instalar-dependencias)  
    - [Conexión a la Base de Datos MySQL](#3-conexión-a-la-base-de-datos-mysql)  
    - [Ejecutar Migraciones con Flyway](#4-ejecutar-migraciones-con-flyway)  
4. [Arquitectura del Proyecto](#arquitectura-del-proyecto)  
    - [Capa Domain](#capa-domain)  
    - [Capa Repository](#capa-repository)  
    - [Capa Service](#capa-service)  
    - [Capa Handler (Controladores)](#capa-handler-controladores)  
    - [Middleware de Seguridad (JWT)](#middleware-de-seguridad-jwt)  
5. [Pruebas Unitarias](#pruebas-unitarias)  
6. [Documentación con Swagger](#documentación-con-swagger)  
7. [Ejecutar la Aplicación](#ejecutar-la-aplicación)  
8. [Autor / Créditos](#autor--créditos)

---

## Tecnologías / Dependencias

- **Go** (versión 1.19 o superior)
- **MySQL** (5.7 / 8.x)
- **Gin** (framework web en Go)
- **Flyway** (migraciones de BD)
- **JWT** para autenticación
- **Swagger** (Swaggo) para la documentación
- **Testify** y **go-sqlmock** para pruebas

---

## Estructura de Directorios

```bash
.
├── cmd
│   └── main.go               # Punto de entrada de la aplicación (contiene configuración general y rutas)
├── internal
│   ├── config
│   │   └── database.go       # Configuración y conexión a MySQL
│   ├── domain
│   │   └── candidate.go      # Modelo de dominio (Candidate)
│   ├── handler
│   │   ├── auth_handler.go   # Endpoint para /login (generar token JWT)
│   │   └── candidate_handler.go # Endpoints CRUD de Candidatos
│   ├── repository
│   │   └── candidate_repository.go
│   ├── security
│   │   └── auth_middleware.go  # Middleware de JWT
│   └── service
│       └── candidate_service.go
├── migrations
│   ├── V1__create_table_candidates.sql
│   └── V2__initial_data_candidates.sql
├── test
│   ├── repository
│   │   └── candidate_repository_test.go
│   └── service
│       └── candidate_service_test.go
├── docs                      # Generado por swag init (contiene swagger.json, swagger.yaml, docs.go)
├── go.mod
├── go.sum
└── README.md
```

---

## Configuración del Proyecto

## Inicializar el Módulo Go

1.- Dentro de la directorio raíz del proyecto

```bash
go mod init github.com/torvictorvic/seek-v2/
```

Esto creará un archivo go.mod con la ruta del módulo.

## Instalar Dependencias

2- Instalar Dependencias

En go.mod se tiene:

```bash
require (
    github.com/gin-gonic/gin v1.9.0
    github.com/go-sql-driver/mysql v1.6.0
    github.com/golang-jwt/jwt/v4 v4.4.2
    github.com/stretchr/testify v1.8.2
    github.com/swaggo/gin-swagger v1.8.3
    github.com/swaggo/files v1.8.3
    github.com/stretchr/testify/mock v1.8.2
    github.com/DATA-DOG/go-sqlmock v1.5.0
)
```

Si se genera al instalar puedes hacer esto:

```bash
go get github.com/gin-gonic/gin
go get github.com/go-sql-driver/mysql
go get github.com/golang-jwt/jwt/v4
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go mod tidy
```


## Conexión a la Base de Datos MySQL
3.- Conexión a la Base de Datos MySQL

Define tu variable de entorno DB_URL o ajusta en internal/config/database.go. 
Ejemplo:

```bash
root:password@tcp(localhost:3306)/candidates_db?parseTime=true
```



## Ejecutar Migraciones con Flyway
4.- Ejecutar Migraciones con Flyway

4.1.- Instalar Flyway (la mejor opción en con Docker).

4.2.- Ubícate en la directorio del proyecto.

4.3.- Ejecutar, dependiendo de la instalacion

```bash
docker run --rm -v $(pwd)/migrations:/flyway/sql flyway/flyway \
  -url="jdbc:mysql://localhost:3306/seek" \
  -user="root" \
  -password="my-secret-pw" \
  migrate
```

o

```bash
flyway -user=<USER> -password=<PASSWORD> \
  -url="jdbc:mysql://localhost:3306/seek?useSSL=false" migrate
```

4.4.- Flyway ejeutará los scripts de .sql que existen en directorio migrations/


---

## Arquitectura del Proyecto

## Capa Domain

Representa las entidades del negocio. 

internal/domain/candidate.go:

```bash
package domain

import "time"

type Candidate struct {
    ID             int       `json:"id"`
    Name           string    `json:"name"`
    Email          string    `json:"email"`
    Gender         string    `json:"gender"`
    SalaryExpected float64   `json:"salary_expected"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

```



## Capa Repository

Comunica la aplicación con la base de datos (infraestructura).

candidate_repository.go:

```bash
type CandidateRepository interface {
    Create(domain.Candidate) (int, error)
    // ...
}

type candidateRepositoryImpl struct {
    db *sql.DB
}

func NewCandidateRepository(db *sql.DB) CandidateRepository {
    return &candidateRepositoryImpl{db: db}
}

// (Implementación CRUD)


```



## Capa Service

Representa las entidades del negocio. 

candidate_service.go:

```bash
type CandidateService interface {
    CreateCandidate(domain.Candidate) (int, error)
    // ...
}

type candidateServiceImpl struct {
    repo repository.CandidateRepository
}

func NewCandidateService(r repository.CandidateRepository) CandidateService {
    return &candidateServiceImpl{repo: r}
}

func (s *candidateServiceImpl) CreateCandidate(c domain.Candidate) (int, error) {
    // validaciones
    if c.Name == "" || c.Email == "" {
        return 0, fmt.Errorf("los campos 'Name' y 'Email' son obligatorios")
    }
    return s.repo.Create(c)
}

```





## Capa Handler (Controladores)

Exponen los endpoints. 

candidate_handler.go:

```bash
// CreateCandidate godoc
// @Summary Crear un nuevo candidato
// @Router /candidates [post]
// @Security Bearer
func (h *CandidateHandler) CreateCandidate(c *gin.Context) {
    var candidate domain.Candidate
    if err := c.ShouldBindJSON(&candidate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
        return
    }
    id, err := h.service.CreateCandidate(candidate)
    // ...
}

```








## Middleware de Seguridad (JWT)

Exponen los endpoints. 

internal/security/auth_middleware.go:

```bash
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Lee "Authorization: Bearer <token>"
        // Valida token JWT
        // Si es inválido => 401
        c.Next()
    }
}

```









---

## Pruebas Unitarias

En la raíz del proyecto ejecutar 

```bash
go test -v ./...
```

Debe retornar el resultado de las pruebas

```bash
=== RUN   TestCreateCandidate
--- PASS: TestCreateCandidate (0.00s)
=== RUN   TestGetByIDCandidate_Found
--- PASS: TestGetByIDCandidate_Found (0.00s)
=== RUN   TestGetByIDCandidate_NotFound
--- PASS: TestGetByIDCandidate_NotFound (0.00s)
=== RUN   TestUpdateCandidate
--- PASS: TestUpdateCandidate (0.00s)
=== RUN   TestDeleteCandidate
--- PASS: TestDeleteCandidate (0.00s)
PASS
ok  	github.com/torvictorvic/seek-v2/test/repository	0.004s
=== RUN   TestCreateCandidate_Success
--- PASS: TestCreateCandidate_Success (0.00s)
=== RUN   TestCreateCandidate_MissingData
--- PASS: TestCreateCandidate_MissingData (0.00s)
=== RUN   TestGetCandidateByID_Success
--- PASS: TestGetCandidateByID_Success (0.00s)
=== RUN   TestGetCandidateByID_NotFound
--- PASS: TestGetCandidateByID_NotFound (0.00s)
=== RUN   TestUpdateCandidate_Success
--- PASS: TestUpdateCandidate_Success (0.00s)
=== RUN   TestUpdateCandidate_RepoError
--- PASS: TestUpdateCandidate_RepoError (0.00s)
=== RUN   TestDeleteCandidate_Success
--- PASS: TestDeleteCandidate_Success (0.00s)
=== RUN   TestDeleteCandidate_RepoError
--- PASS: TestDeleteCandidate_RepoError (0.00s)
PASS
ok  	github.com/torvictorvic/seek-v2/test/service	0.003s

```








---

## Documentación con Swagger

Instalar la herramienta

```bash
go install github.com/swaggo/swag/cmd/swag@latest

```

En el archivo princiapl main.go se tienen estas anotaciones que son típicas de Swagger


```bash
// @title Sistema de Gestión de Candidatos
// @version 1.0
// @description Esta API gestiona candidatos en un proceso de reclutamiento
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

```

Y así para internal/handler/candidate_handler.go tiene esas anotaciones.

Para ejecutar la documentacion se corre:


```bash
swag init -g cmd/main.go --parseDependency --parseInternal

```

Esto generará archivos en docs/


```bash
.
├── docs
│   └── docs.go
│   └── swagger.json
│   └── swagger.yaml   
```

Luego en el siguiente paso, se correrá el proyecto y se puede abrir en en un navegador la siguiente dirección


```bash
http://localhost:8080/swagger/index.html
```









---

## Ejecutar la Aplicación

7-1.- Configura las variables de entorno (o ajustar database.go)

```bash
export DB_URL="root:password@tcp(localhost:3306)/seek?parseTime=true"
export JWT_SECRET="MiSecretoSuperSeguroXXXTTYYYY"

```

7.2.- Compila y ejecutar

```bash
go run ./cmd

```

Debe mostrar esta salida


```bash
Conexión exitosa a la base de datos
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /login                    --> github.com/torvictorvic/seek-v2/internal/handler.GenerateToken (3 handlers)
[GIN-debug] POST   /api/candidates           --> github.com/torvictorvic/seek-v2/internal/handler.(*CandidateHandler).CreateCandidate-fm (4 handlers)
[GIN-debug] GET    /api/candidates/:id       --> github.com/torvictorvic/seek-v2/internal/handler.(*CandidateHandler).GetCandidateByID-fm (4 handlers)
[GIN-debug] GET    /api/candidates           --> github.com/torvictorvic/seek-v2/internal/handler.(*CandidateHandler).GetAllCandidates-fm (4 handlers)
[GIN-debug] PUT    /api/candidates/:id       --> github.com/torvictorvic/seek-v2/internal/handler.(*CandidateHandler).UpdateCandidate-fm (4 handlers)
[GIN-debug] DELETE /api/candidates/:id       --> github.com/torvictorvic/seek-v2/internal/handler.(*CandidateHandler).DeleteCandidate-fm (4 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (3 handlers)
2024/12/27 08:08:37 Server run http://localhost:8080
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2024/12/27 - 08:08:50 | 200 |     268.767µs |       127.0.0.1 | GET      "/swagger/index.html"
```


7.3.- Probar con cURL o Postman:

Leer el token generado

```bash
POST http://localhost:8080/login
```

Con el token generado, usar este servicio y en Authorization colocar Bearer {TOKEN}

```bash
GET http://localhost:8080/api/candidates
```








---

## Autor / Créditos

Nombre: Victor Manuel Suárez / torvictorvic

Contacto: victormst@gmail.com

Descripción: Este proyecto fue desarrollado con práctica de Go + MySQL + JWT + Swagger + Hexagonal Architecture para un test técnico de SEEK.