# Template Go
Repositorio base para la creacion de un servicio REST en GO. 

# Old Docs

se cambio de sqlc a pgx solamente y scany, no me gusto para nada sqlc y tendria q hacer mucho boilerplate xd.

to do: 
    agregar inicializacion de bd a main.go



El REST tiene:
 * Implementaciones para Mongo y Postgres.
 * Uso de JWT (simetrico y asimetrico con detección automática) y acceso restringido a servicios con roles.
 * Hasheo de passwords con bcrypt.  
 * Documentacion en Swagger.

 # Importante
Es clave mencionar que la lectura completa de este documento es de gran utilidad si es la primera vez que está integrando este proyecto.

# Tabla de contenidos
- [Backend golang template](#Backend-golang-template)
  - [Tabla de contenidos](#Tabla-de-contenidos)
- [Importante](#Importante)
  - [Guía rápida inicio](#Guía-rápida-inicio)
    - [Cambiar nombre paquete Notice](#Cambiar-nombre-paquete)
    - [Ejecución](#Ejecución)
        - [Docker](#Docker)
        - [Sin Docker](#Sin-Docker)
    - [Importante detalles de ejecución](#Importante-detalles-de-ejecución)
  - [Gin](#Gin)
  - [CORS](#Cors)
  - [JWT](#JWT)
    - [¿Qué es autenticar?](#¿Qué-es-autenticar?)
    - [¿Qué es autorizar?](#¿Qué-es-autorizar?)
    - [Contenido de un JWT](#Contenido-de-un-JWT)
    - [Implementación en Go](#Implementación-en-Go)
        - [jwt-go](#jwt-go)
        - [Gin JWT Middleware (appleboy)](#Gin-JWT-Middleware-(appleboy))
    - [Log](#Log)
        - [Loggers conocidos en Go](#Loggers-conocidos-en-Go)
            - [Lumberjack's Logger](#Lumberjack's-Logger)
            - [Zap's Logger](#Zap's-Logger)
              - [Niveles de Logging](#Niveles-de-Logging)
    - [Bases de datos](#Bases-de-datos)
        - [Mongo Driver](#Mongo-Driver)
        - [GORM](#GORM)
        - [Migrate](#Migrate)
    - [Swagger](#Swagger)
    - [Estructura del proyecto](#Estructura-del-proyecto)
        - [Arquitectura](#Arquitectura)
            - [Models](#Models)
            - [Repositories](#Repositories)
            - [Services](#Services)
            - [Controllers](#Controllers)
            - [Routes](#Routes)
            - [Middleware](#Middleware)
            - [Utils](#Utils)
            - [Autenticación](#Autenticación)
        - [Para registrar un nuevo servicio](#Para-registrar-un-nuevo-servicio)
    - [Referencias](#Referencias)

# Guía rápida inicio 
## Cambiar nombre paquete
Lo primero a realizar es cambiar el nombre del paquete de go para comenzar a trabajar. El nombre actual del paquete es citiaps/golang-backend-template. Lo que se debe hacer entonces, es cambiar el nombre del repositorio en github en donde va a residir el proyecto.

## Importante detalles de ejecución
Es de suma importancia destacar que el proyecto funciona tanto para Mongo como para Postgres, el único pero de esta implementación es que se debería comentar/eliminar el código de la implementación que *no* vaya a utilizar y *debe* tener previamente la base de datos a utilizar levantada correctamente.

Es decir, si utilizará MongoDB y su implementación, comente/elimine el código respectivo de la implementación no utilizada. Esto debido a que ahora están ambos funcionando si y solo si, ambas bases de datos están arriba, debido a que se inicializan las conexiones para ambas bases de datos en el momento de arrancar en el siguiente pedazo de código.

```bash 
func InitRepositories() {
	catRepo = repositories.NewCatRepository()
	userRepo = repositories.NewUserRepository()

	catRepoPostgres = repositories.NewCatRepositoryPostgres()
	userRepoPostgres = repositories.NewUserRepositoryPostgres()
}
```

## Ejecución
### Docker
1) Defina las variables de entorno en el archivo .env de la carpeta deployment, puede hacer uso como guía del archivo .env.example.
2) Para enfoques de producción o desarrollo debe dirigirse al path golang-backend-template/deployment:
    2.1) Producción
        Ejecute docker compose up -d --build
    2.2) Desarrollo
        Ejecute docker compose -f docker-compose-dev.yml up -d --build

### Sin Docker
1) Defina las variables de entorno en el archivo .env de la carpeta deployment, puede hacer uso como guía del archivo .env.example.
3) Levante una base de datos. (MongoDB o Postgres)
2) Dirijase a backend/ y ejecute go run main.go.

# Gin
Es un framework escrito en Golang utilizado para crear aplicaciones HTTP/API REST. Es bastante rápido, ofrece soporte para middlewares, validaciones json, gran control de rutas y mucho más.

# CORS
Cross Origin Resource Sharing, es básicamente una medida de control para restringir qué origenes (http://localhost:8000, http://localhost:3000, etc.) pueden acceder a recursos.

# JWT
Json Web Token (JWT) es una idea que aparece al querer transmitir/intercambiar información de manera segura, son objetos que nos permiten *Autenticar* y *Autorizar* a un usuario.

Su contenido es encriptado con opciones a 2 tipos de algoritmos, simétricos y asimétricos, el primero es basicamente: encripto todo con una llave privada y por otro lado, el segundo encripto con una llave privada y verifico esta autenticidad con una llave pública.

## ¿Qué es autenticar?
Este término es básicamente la acción de comprobar un inicio de sesión, es decir, la persona que esta diciendo que es, es realmente esa persona?. 

## ¿Qué es autorizar?
Este término entra en juego al querer acceder a diversos recursos, por ejemplo, si un usuario con rol normal intenta acceder a un recurso que tiene rol admin ¿debería autorizar ese acceso?. 

## Contenido de un JWT
Dentro de este apartado existe algo bien conocido por "payload", esto es básicamente el contenido de este objeto, en el puede ir información variada pero en general suele tener:

    - user id
    - token 
    - expire
    - etc.

Es importante mencionar que este contenido, no debe incluir información sensible respecto del usuario, ya que alguien podría romper nuestro JWT y poder acceder a esta.

## Implementación en Go
En el contexto de Go, existen 2 grandes implementaciones de esta herramienta que son:

### jwt-go 
Proporciona funciones básicas de JWT, encarga al desarrollador de manejar la integración y flujos.

### Gin JWT Middleware (appleboy)
Es un middleware que utiliza jwt-go, proporciona integración directa con Gin que lo hace más rápido y amigable directamente. 

Dentro de las funciones que provee/necesita este middleware están las siguientes:

    Proporcionada: LoginHandler -> sigue un flujo de login con las funciones de más abajo, es decir va pasando por cada una de ellas.
    Requerida: Authenticator -> es básicamente el "login" de nuestra aplicación.
    Opcional: PayloadFunc -> es la función que obtiene las claims una vez hecho el login.
    Opcional: LoginResponse -> es una personalizacion de la respuesta al hacer un login exitoso.
    Proporcionada: MiddlewareFunc -> es el middleware que gin, debe ser usado en cada llamada a endpoints que requieran un token jwt.
    Opcional: IdentityHandler -> es una función que le pasa las claims al authorizator.
    Opcional: Authorizator -> es una función que debería verificar si el usuario está/no está autorizado para acceder a X recurso.
    Opcional: LogoutHandler -> limpia cookies solamente.
    Opcional: LogoutResponse -> es una respuesta personalizada del logout.
    Proporcionada: RefreshHandler -> obtiene un nuevo token, es el concepto de "refrescar/refresh".
    Opcional: RefreshResponse -> es una respuesta personalizada del refresh.
    Opcional Unauthorized ->  es una respuesta personalizada en caso de errores, token invalido, etc.

# Log
Un log es una herramienta de "debug" (en muchas ocasiones, en otras es utilizada como una herramienta informativa como un info, warning, etc.) que nos permite "seguir el código" en pantalla (qué pasa cuando algo sucede mal, bien, etc.).

## Loggers conocidos en Go
### Lumberjack's Logger
Es un conocido logger ampliamente utilizado con el objetivo de escribir Logs en archivos específicos. 

Tiene varias personalizaciones con respecto a cuánto es el tiempo de vida de los archivos, tamaño máximo, etc. Pero eso no se abordará en esta documentación. 

IMPORTANTE (dicho por el mismo creador): este logger no es un TODO, es una gran herramienta pero debe ser acompañada (que veremos adelante) por otras para poder sacar partido a su potencial.

### Zap's Logger
Es un conocido logger desarrollado por Uber, entrega logging rápido, estructurado y por niveles.

#### Niveles de Logging (Leveled Logging)
Este término hace mención a que (valga la redundancia) existen distintos niveles para poder mostrar cosas por pantalla, esto puede ser algo que ya mencionamos anteriormente como:

    - Warn (advertencia y continua el codigo)
    - Info (info solamente y continua el codigo)
    - Fatal (error fatal y se detiene la ejecucion del codigo)
    - Debug (info solamente y continua el codigo)
    - Error (error y continua el codigo)
    - Etc.

Es bastante útil ya que nos entrega una estructura a seguir respecto a cuándo usar cada uno.

Esta herramienta se combina de buena manera con Lumberjack's Logger debido a que es un gran combo, niveles, rápido, estructurado y escritura en archivos.
 
# Bases de datos
El proyecto tiene a disposición dos tipos de bases de datos a utilizar, PostgresSQL y MongoDB.

La mayor diferencia entre estas viene dada por la implementación de modelos, repositorios, servicios y más.

Por ejemplo, Mongo tiene un tipo de dato específico para el _id, Postgres, también, así, estos no son compatibles entre sí, por ende significa hacer cambios en la implementación, más adelante en la descripción de capas existen ejemplos de estas diferencias.

Es importante mencionar que al menos se debe tener Postgres o Mongo levantado para poder utilizar las implementaciones respectivas.

## Mongo Driver 
Es la implementación escrita en Golang para utilizar Mongo, ofrece un catálogo gigante de funciones útiles a integrar en el código, conexión a la base de datos, ping, insertar, filtrar, etc.

## GORM 
Es la implementación escrita en Golang para utilizar múltiples bases de datos SQL (en este caso utilizada para postgres), ofrece un catálogo gigante de funciones útiles a integrar en el código, conexión a la base de datos, ping, insertar, filtrar, etc.

## Migrate
Es una librería de Go que nos permite migrar bases de datos de manera sencilla, utiliza 2 tipos de archivos, up y down, donde uno se utiliza para "crear" y el otro para "eliminar los cambios".

En el directorio migrations/ podrá encontrar ejemplos tanto para mongo y postgres.

Para poder utilizar los archivos up basta con utilizar:

migrate -database 'URI-DATABASE' -path 'PATH-UP-FILES' up 

Para poder utilizar los archivos down basta con utilizar:

migrate -database 'URI-DATABASE' -path 'PATH-DOWN-FILES' down

## Swagger
Swagger es una gran herramienta para poder documentar de manera rápida y sencilla a través de comentarios. Ofrece una primer organización en el archivo main, siendo esta la siguiente:

``` bash
// Docs
docs.SwaggerInfo.BasePath = "/api/v1/x"
docs.SwaggerInfo.Version = "v1"
docs.SwaggerInfo.Host = "localhost:8080"
// Route docs
app.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

```

Donde se especifican versiones, host, path para las llamadas http, y una ruta la cual mostrará la documentación creada.

Luego, en el mismo archivo se puede agregar la primera capa de documentación, siendo esta la siguiente:

```bash
// @title          Template
// @version        1.0
// @description    Template description
// @termsOfService http://swagger.io/terms/

// @contact.name  CITIAPS
// @contact.url   https://citiaps.cl
// @contact.email citiaps@usach.cl

// lincense.name  Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @tag.name        main
// @tag.description main description

// @host     localhost:8080
// @BasePath /api/v1/x

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                BearerJWTToken in Authorization Header

// @accept  json
// @produce json

// @schemes http https
```

Luego, la documentación específica se realiza en la capa de controladores, donde siguen la siguiente estructura:

```bash
// CreateUserController
// @Title CreateUserController
// @Description Permite crear un usuario en el sistema con manejo de errores, etc.
// @Summary Crea un usuario
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200 {object} forms.UserForm "Usuario creado con exito."
// @Router /user/ [post]
func CreateUserController(c *gin.Context) {}
```

Una generalidad de esto viene dado por:

    - Nombre de función: {{func}} 
    - Resumen de la funcionalidad: @Summary
    - Descripción de la funcionalidad: @Description
    - Etiqueta: @Tags
    - MIME(s) que acepta: @Accept
    - MIME(s) que produce: @Produce
    - Estado de exito: @Success
    - Posibles respuestas de error: @Failure
    - Path y verbo: @Router

Para poder visualizar la documentación generada por swagger basta con acceder a la siguiente url una vez levantado el backend:
    - http://localhost:8000/swagger/index.html#/

# Estructura del proyecto

El proyecto está orientado de la siguiente manera:

```bash

├── assets
│   └── .keep
├── config
│   ├── database
│        ├── mongoConnection.go
│        └── postgresConnection.go
│   └── config.go
├── controllers
│   ├── Cat.go
│   ├── Auth.go
│   └── User.go
├── docs
│   └── docs.go
├── forms
│   ├── CatForm.go
│   └── UserForm.go
├── keys
│   └── .gitkeep
├── logs
├── mailer
│    ├── templates
│         └── confirmCreatedUser.html
│    └── mailer.go
├── middleware
│    ├── authentication.go
│    └── cors.go
├── models
│    ├── Cat.go
│    ├── Login.go
│    └── User.go
├── repositories
│    ├── Cat.go
│    ├── config.go
│    └── User.go
├── routes
│    ├── Auth.go
│    ├── Cat.go
│    ├── routes.go
│    └── User.go
├── services
│    ├── Cat.go
│    ├── index.go
│    └── User.go
├── storage
├── utils
│    ├── env.go
│    ├── files.go
│    ├── hash.go
│    ├── logger.go
│    └── response.go
├── Dockerfile.dev
├── Dockerfile.prod
├── go.mod
├── go.sum
├── main.go
├── deployment
│    ├── .env
│    ├── .env.example
│    ├── docker-compose.yml
│    └── docker-compose-dev.yml
├── .dockerignore
├── .gitignore
├── Dockerfile
├── Dockerfile.dev
└── README.md
```
## Arquitectura 
El proyecto esta orientado a una arquitectura por capas, donde:

### Models
Tiene los modelos de cada uno de los recursos del sistema. En este solo se debería implementar estructuras de las entidades a utilizar.

Por ejemplo, si quisieramos modelar un Gato, se entiende que este debería tener un identificador (si no como lo diferenciamos?), un nombre, edad, y más.

Así, la definición se haría de la siguiente manera:

```bash
type Cat struct {
	ID    primitive.ObjectID # o uint si fuera postgres
	Name  tipoDeDatoName             
	Age   tipoDeDatoAge               
	Owner tipodeDatoOwner
}
```
Sabemos que en mongo los id vienen dados por primitive.ObjectID, en postgres son uint, entonces varía dependiendo de la implementación.

### Repositories
Tiene únicamente operaciones que requieren interactuar con la base de datos y nada más que eso, no existen validaciones, lógica de negocio, etc.

Por ejemplo, si quisieramos insertar un Gato vendría ser el siguiente trozo de código si estuviesemos trabajando en Mongo:

```bash
func (catRepo *CatRepository) InsertOne(cat *models.Cat) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := catRepo.Collection.InsertOne(ctx, cat)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
```

La gran diferencia entre esta implementación y una de postgres, son directamente las funciones que el driver de la base de datos nos ofrece, en este caso usando GORM y Postgres la función insertar viene dada por:

```bash
func (catRepoPostgres *CatRepositoryPostgres) InsertOnePostgres(cat *models.CatPostgres) (uint, error) {
	result := catRepoPostgres.DB.Create(cat)

	if result.Error != nil {
		return 0, result.Error
	}

	return cat.ID, nil
}
```

### Services
Tienen validaciones de datos y lógica de negocios en general, en este caso, para ambas implementaciones, Postgres y Mongo se deberían hacer las mismas validaciones el servicio es prácticamente igual. La única diferencia es en el comentario de abajo:

```bash
func CreateCatService(cat *models.CatPostgres) (*models.CatPostgres, error) {
    # Muchas validaciones que el desarrollador desee incluir

    id, err := catRepo.InsertOne(cat) # caso mongo
    # id, err := catRepoPostgres.InsertOnePostgres(cat) caso postgres

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	cat.ID = id
	return cat, nil
}
```

### Controller
Tiene los controladores de cada uno de los recursos del sistema. En este solo se debería manejar la lógica, recepción y respuesta de las consultas.

En este caso, si quisieramos por ejemplo crear un gato, el controlador deberia verificar que el objeto que viene en la solicitud es efectivamente un gato. De esa forma, Gin nos ayuda de la siguiente manera:

```bash
func CreateCatController(c *gin.Context) {
	var newFormCat *forms.CatForm
	err := c.BindJSON(&newFormCat)

    # mucho código eliminado 

	newCat := &models.Cat{
		Name:  newFormCat.Name,
		Age:   newFormCat.Age,
		Owner: userID,
	}

    # intento de creacion de gato
	returnedCat, err := services.CreateCatService(newCat)
	
    # feedback
	utils.Info("Creacion exitosa de gato del usuario con id " + userIDStr + " desde ip: " + c.ClientIP())
	utils.JsonResponse(c, 201, "Gato creado con exito.", returnedCat)
}
```

Siguiendo la línea de Mongo o Postgres, directamente deberíamos utilizar en vez del modelo usual de Cat, CatPostgres y respectivamente el servicio para poder insertar un gato en la base de datos Postgres.

### Routes
Tiene las rutas de los end-points de la aplicación, divididas según entidad. Las rutas son registradas en la aplicación en el archivo routes.go

Por ejemplo si quisieramos registrar las rutas vinculadas a Cat, sería:

```bash
func InitCatRoutes(r *gin.RouterGroup) {
	catGroup := r.Group("/cat")
	{
		catGroup.POST("/", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.CreateCatController)
		catGroup.GET("/", middleware.SetRoles(models.ADMIN), middleware.LoadJWTAuth().MiddlewareFunc(), controllers.GetAllCatsController)

		// postgres
		catGroup.POST("/postgres", controllers.CreateCatControllerPostgres)
		catGroup.GET("/postgres", controllers.GetAllCatsControllerPostgres)
	}
}
```

Es importante mencionar que en cada ruta se pueden agregar middlewares respectivos, además, que estas rutas deben ser utilizas por el motor de Gin, así, también deben ser registradas y utilizadas en Gin.

### Middleware
Tiene middlewares definidos como la implementación de JWT y la configuración de CORS, son utilizados para proteger los servicios a través de roles, restricciones de origenes, etc.

### Utils
Paquete que tiene funciones utilizables a lo largo de toda la aplicación.

Por ejemplo, 

```bash
func LoadEnv() {
	_ = godotenv.Load("../deployment/.env")
}
```

Una funcion reutilizable que únicamente lee variables de entorno.

### Autenticación
Los roles son definidos en [models/UserModel.go](models/UserModel.go) como constantes en el archivo.

Para exigir autenticación en alguna ruta, utilizar el middleware correspondiente ```middleware.LoadJWTAuth().MiddlewareFunc()```:

Para exigir autenticación con roles en alguna ruta, utilizar el middleware correspondiente ```middleware.SetRoles(example_rol), middleware.LoadJWTAuth().MiddlewareFunc()```:


Por ejemplo, para proteger una ruta que solo pueda ser utilizada por un usuario logueado
```golang
catGroup.POST("/", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.CreateCatController)
```

Para proteger una ruta dependiendo de rol del usuario utilizar el middleware antes de la autorización.

```golang
catGroup.GET("/", middleware.SetRoles(models.ADMIN), middleware.LoadJWTAuth().MiddlewareFunc(), controllers.GetAllCatsController)
```

Si se quiere permitir a mas de un rol, se puede agregar como otro parametro.

Hay más ejemplos disponibles al interior del directorio [routes](backend/routes/)

## Para registrar un nuevo servicio
1. Crear modelos como struct en el directorio models
2. Crear repositorio que maneja los accesos a la base de datos
3. Crear servicio que maneje la logica de negocios
4. Crear controlador que reciba y gestione los parametros https
5. Crear ruta que ejecute la funcionalidad creada en el controlador, agregar restricciones de autenticación y roles si es necesario.

# Referencias
    - https://github.com/natefinch/lumberjack
    - https://github.com/appleboy/gin-jwt?tab=readme-ov-file#gin-jwt-middleware
    - https://github.com/golang-jwt/jwt
