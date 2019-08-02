---
title:
- neon-cli; neon frakework v2
author:
- Rg
theme:
- Copenhagen
---

# Current Problem

- Directory Structs

- Bff Interfaces

- Service Interfaces

- Freedom

- Package Manager

- Doc Generate

- Validate Generate

- Mock Server

- SQL Parser

# Solve ?

![](http://images.xuejuzi.cn/1707/1_170709095348_1.jpg)

# Directory structs

``` shell
xx_foo_system
	+ services
	+ bff
	+ protocol
	+ doc
    + vendor
    - go.mod
    - go.mod
	- README.md
```

# Services Directory 
``` shell
xx_foo_system
	+ services
		+ demo
			+ config
                - config.go
				- config.yaml
				- config_test.yaml
			+ impls
			  + server.go
			  + group
			    - impl_change_permission_in_group.go
			    - server.go
			+ hook
			  - hook.go
			- main.go
			- .k8s.yaml
			- .k8s_test.yaml
			- Dockerfile
            - Makefile
```

# Bff Directory

``` shell
	+ bff
		+ admin
		  + codes
		  	- error_code.go
		  + config
            - config.go
		    - config.yml
		    - config_test.yml
		+ impls
		  + group
		  	- impl_change_permission_in_group.go
		  	- types.go
		  + department
		    - impl_get_department_list.go
		  - impl_login.go
		+ router
		 - router_base.go
		+ hook
		 - hook.go
		- main.go
		- .k8s.yaml
		- .k8s_test.yaml
		- Dockerfile
        - Makefile
```

# Protocol

``` shell
    + protocol
		+ demo
         - demo.proto
		 - group.go
		 - login.go
```

**Note**

- `+`: Directory

- `-`: File

# Annotation

**Note**

> structs uniq

> not cross directory

# Annotation
## Type

- s.i | service.interface
- s.i.r | serivice.interface.request
- s.i.r | service.interface.response
- b.i | bff.interface
- b.i.rt | bff.interface.request
- b.i.re | bff.interface.response

# Annotation
## Bff Annotation
### interface
```go
// @type: b.i.rt 
// @name: login 
// @login: Y
// @page: xxxx | yyyy
// @uri: /api/admin/v1/login
// @des: xxxx
func LoginHandler(state *bff.State) {
}
```
*Option*: des

# Annotation
## Bff Annotation
### request
```go
// @type: b.i.r
// @interface: LoginHandler
// @des: xxxx
type LoginRequest struct {
	GroupId      int64 `binding:"gte=1"` // User Group id | Y | 0 |
}
```

# Annotation
## Bff Annotation
### response
```go
// @type: b.i.re
// @interface: LoginHandler
// @des: xxxx
type LoginResponse struct {
}
```

# Annotation
## Services Annotation
### interface
```go
// @type: s.i
// @path: demo
// @des: Demo
service Demo {
	rpc	Ping(PingRequest) returns (PingResponse);
}
```
# Annotation
## Services Annotation
### PingRequest

``` protocol
// @type: s.i.rt
// @des: Demo
message PingRequest {}
```
# Annotation
## Services Annotation
### PingResponse

``` protocol
// @type: s.i.re
message PingResponse {
}
```

# Protocol

```
+ pb
	+ demo
		- login.go
        - login.proto
	+ foo
        - ping.proto
		- ping.go
```

# Mocker Server

- uri

# Datatabse Operate

- Create Table Mapping to Struct

- Add, Delete, Query, Update

- Split Page
