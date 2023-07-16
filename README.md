## How to launch CAIS (campaigns and items service)?\
First of all, you should create a .env file in the root directory of the project,
which will contain the following information:
```dotenv
HOST=
PORT=
USER=
PASSWD=
DBNAME=
SSLMODE=
```

## Drivers and tools
### Postgresql
**1. Client:**
```
"github.com/jackc/pgconn"
"github.com/jackc/pgx/v4"
```

**2. Migrations:**
```
github.com/golang-migrate/migrate
```

### Packages:

**Environment variables:**
```
"github.com/joho/godotenv"
```

**Errors:**
```
"github.com/pkg/errors"
```
**Logging:**
```
github.com/sirupsen/logrus
```

**Router:**
```
"github.com/julienschmidt/httprouter"
```