# Springboot API

Spring Boot service that exposes the seeded contents of `viewData.json` through a single REST endpoint for my internship test.

## Prerequisites

- Java 25 (matching `pom.xml`'s `java.version`)
- Maven 3.9+ (the project ships with `mvnw` wrappers if Maven is not installed)
- PostgreSQL 16 (or compatible) running locally on the default port

## Database setup

1. Ensure PostgreSQL is running and reachable at `jdbc:postgresql://localhost:5432/viewdata` with user `postgres` and password `postgres` (as configured in `src/main/resources/application.properties`).
   - **Docker option** (recommended for local dev):
     ```bash
     docker run --name pg \
       -e POSTGRES_PASSWORD=postgres \
       -e POSTGRES_DB=viewdata \
       -p 5432:5432 -d postgres:16
     ```
     Stop with `docker stop pg` and start again using `docker start pg`.
   - **Native install**: create the database manually if it does not exist:
     ```bash
     createdb -U postgres viewdata
     ```
2. Adjust credentials/URL in `application.properties` if your environment differs.

> On startup the application seeds data into the database from `src/main/resources/viewData.json`. Existing rows are preserved by the `ddl-auto=update` setting, but you can drop the schema manually if you want a fresh load.

## Running the application

Use the Maven wrapper (recommended so the correct Maven version is used):

```bash
./mvnw spring-boot:run
```

The app starts on `http://localhost:8080`.

To stop the server, press `Ctrl + C` in the terminal.

## API

- **GET** `/api/view-data`
  - Returns the entire dataset with two arrays: `data` (transactions) and `status` (reference data).
  - Example response (truncated):
    ```json
    {
      "data": [
        {
          "id": 1372,
          "productID": "10001",
          "productName": "Test 1",
          "amount": "1000",
          "customerName": "abc",
          "status": 0,
          "transactionDate": "2022-07-10 11:14:52",
          "createBy": "abc",
          "createOn": "2022-07-10 11:14:52"
        }
      ],
      "status": [
        { "id": 0, "name": "SUCCESS" },
        { "id": 1, "name": "FAILED" }
      ]
    }
    ```

## Calling the API from the CLI

```bash
curl -s http://localhost:8080/api/view-data | jq
```

- Remove the pipe to `jq` if it is not installed.
- Add `-v` to see request/response headers.

## Calling the API with Postman

1. Start the application so it is listening on `http://localhost:8080`.
2. Open Postman and create a new **GET** request tab.
3. Set the request URL to `http://localhost:8080/api/view-data`.
4. Click **Send**. The JSON response should appear in the `Body` tab.
5. (Optional) Save the request to a collection for future use.

## Running tests

```bash
./mvnw test
```

# IT Logical Test Programming

The answer are inside of the folder `LogicalTestAnswer`
