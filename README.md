link: [Packages, drivers and utilities that were used when writing the project](./packages.md)

### How to launch CAIS (campaigns and items service)?
First of all, you should create a .env file in the root directory of the project,
which will contain the following information:
```yaml
PG_HOST=
PG_PORT=
PG_USER=
PG_PASSWD=
PG_DBNAME=
PG_SSLMODE=
REDIS_ADDR=
REDIS_PASS= 
REDIS_DB=
```

### Project Description:
The project is deploying a service on Golang that uses Postgres as a database, 
Clickhouse for logging changes, Nats as a message broker and Redis for temporary caching.

**The main entities of the project: item and campaign.**

### API Methods of the service:
| API Method               | Functional                             | Works |
|--------------------------|----------------------------------------|-------|
| **GetAllItems**          | Retrieves the item list                | ✔     |
| **GetItemsByCampaignId** | Retrieves the item list by campaign id | ✔     |
| **CreateNewItem**        | Create new item in storage             | ✔     |
| **DeleteItem**           | Delete item in storage by id           | ✔     |
| **UpdateItem**           | Update item fields in storage          | ✔     |


### Features:
1. Implemented CRUD methods for the Postgres database
in the items table.

2. When editing data in Postgres (create, update, delete), 
an isolation level is set to block read-write and all requests with manipulations are wrapped in a transaction.

3. When GETTING data requests from Postgres, the data is cached in Redis for a minute. With subsequent GET requests, the service tries to get data first from Redis, if there are none, go to postgresql and puts them in REDIS.

4. When adding, editing or deleting an entry in Postgres, a log is written in batches to Clickhouse via the Nats queue.