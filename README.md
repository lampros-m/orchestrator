# orchestrator-api

1. [ Overview ](#overview)
2. [ How to run - API ](#runapi)
3. [ Configuration ](#configuration)
4. [ Swagger ](#swagger)
5. [ Mock Services ](#mock)
6. [ TODO ](#todo)

<a name="overview"></a>
## 1. Overview

Creates a mechanism that can handle and monitor running processes

<a name="runapi"></a>
## 2. How to run - API
```
make build run-server
```
or execute the already built file
```
./bin/orchestratorserver
```

<a name="configuration"></a>
## 3. How to configure
Edit the `executables.json` file. The format of configuration is in `executables_example.json`.

Example:

```
"name": "Service Charlie",
"binary_path": "/home/bayman/repositories/invisiblez/orchestrator/mockservices/servicec/cmd/main",
"working_dir": "/home/bayman/repositories/invisiblez/orchestrator/mockservices/servicec/cmd",
"log_dir": "/home/bayman/repositories/invisiblez/orchestrator/mockservices/servicec",
"arguments": [],
"log_file_name": "out",
"error_file_name": "errors",
"auto_restart": true,
"group": 2
```

The `.env` file keeps info about:
- `server port` - Server port
- `executables.json` Where the file with executables is located
- `setup` and `run`: If the executables set and run will be applied automatically after the start of the server

<a name="swagger"></a>
## 4. Swagger
In order to update swagger documenation, run `make swag`
After running the API Server, the documentation of the API can be found at: `http://localhost:8090/swagger/index.html`


<a name="mock"></a>
## 5. Mock Services
The mock services package is for mock-test purposes only.

<a name="todo"></a>
## 6. TODO
TODO file for future implementation
