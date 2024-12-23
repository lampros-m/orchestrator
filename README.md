# orchestrator-api

1. [ Overview ](#overview)
2. [ How to run - API ](#runapi)
3. [ How to run - CLI ](#runcli)
4. [ Configuration ](#configuration)
5. [ Swagger ](#swagger)
6. [ Mock Services ](#mock)
7. [ TODO ](#todo)

<a name="overview"></a>
## 1. Overview

Creates a mechanism that can handle and monitor running processes

<a name="runapi"></a>
## 2. How to run - API
```
make build run-server
```
or
```
./bin/orchestratorserver
```

<a name="runcli"></a>
## 3. How to run - CLI
```
make build run-cli
```
or
```
./bin/orchestratorcli
```

<a name="configuration"></a>
## 4. How to configure
Edit the `executables.json` file. The format of configuration is in `executables_example.json`.

Example:

```
"name": "Service Beta",
"binary_path": "/usr/local/go/bin/go",
"working_dir": "/home/bayman/repositories/invisiblez/orchestrator/mockservices/serviceb/cmd",
"log_dir": "/home/bayman/repositories/invisiblez/orchestrator/mockservices/serviceb",
"arguments": ["run", "main.go"],
"log_file_name": "out",
"error_file_name": "errors",
"auto_restart": false
```

The `.env` file keeps info about `server port` and where the `executables.json` file is located.

<a name="swagger"></a>
## 5. Swagger
In order to update swagger documenation, run `make swag`
After running the API Server, the documentation of the API can be found at: `http://localhost:8090/swagger/index.html`


<a name="mock"></a>
## 6. Mock Services
The mock services package is for mock-test purposes only.

<a name="todo"></a>
## 7. TODO
TODO file for future implementation
