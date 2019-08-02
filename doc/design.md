# System Design 

## Directory Struct

system level as the top directory

```xquery
xx_foo_system
	+ services
		+ demo
			+ config
				- config.yaml
				- config_test.yaml
			+ impls
			  + group
			    - impl_change_permission_in_group.go
			    - impl_change_user_in_group.go
			    - server.go
				- impl_delete_group_by_id.go
				- impl_demo.go
				- server.go
			+ hook
			  - hook.go
			- main.go
			- .k8s.yaml
			- .k8s_test.yaml
			- Dockerfile
	+ bff
		+ admin
		  + codes
		  	- error_code.go
		  + config
		    - config.yml
		    - config_test.yml
		+ impls
		  + group
		  	- impl_change_permission_in_group.go
		  	- impl_change_user_in_group.go
		  	- types.go
		  + department
		    - impl_get_department_list.go
		    - impl_get_department_list_with_user.go
		  - impl_login.go
		  - impl_logout.go
		  - impl_demo.go
		+ router
		 - router_base.go
		+ hook
		 - hook.go
		- main.go
		- .k8s.yaml
		- .k8s_test.yaml
		- Dockerfile
	+ pb
		+ demo
		 - group.go
		 - login.go
	- doc
	- README.md
	- Makefile
	- go.mod
	- go.sum
```

### TODO

- [ ] .k8s_yaml
- [ ] .k8s_test.yaml
- [ ] config_test.yaml
- [ ] Dockerfile
- [ ] service impls name、file organization
- [ ] bff  impls name、file organization
- [ ] pb impls name、file organization
- [ ] Errcode impls name、file organization



## Annotation

- Type 
  - service.interface
  - serivice.interface.request
  - service.interface.response
  - bff.interface
  - bff.interface.request
  - bff.interface.response

## Bff

- [ ] Mutil Bff
- [ ] Annotation

- interface

  ```go
  // @type: bff.interface
  // @name: 更改用户组权限
  // @login: Y
  // @page: xxxx
  // @des: xxxx
  func ChangePermissionInGroupHandler(state *bff.State) {}
  ```

  *Option*: des

- interface.request

  ```go
  // @type: bff.interface.request
  // @interface: ChangePermissionInGroupHandler
  // @des: xxxx
  type ChangePermissionInGroupItem struct {
  	GroupId      int64 `binding:"gte=1"`       // 用户组id | Y | 0 |
  	PermissionId int64 `binding:"gte=1"`       // 权限id | Y | 0 |
  	Operate      int64 `binding:"gte=0,lte=1"` // 操作 | Y | 0 | 新增为1，取消为0
  }
  ```

  *Option*: des

- interface.response

  ```go
  // @type: bff.interface.response
  // @interface: ChangePermissionInGroupHandler
  // @des: xxxx
  type ChangePermissionInGroupRequest struct {
  	List []*ChangePermissionInGroupItem // 要新增或者取消的list | Y | [] |
  }
  ```

**Note**

> - `bff`可以有多层结构，但是接口、请求参数、返回参数不能都出现重复现象
>
> - `interface.request`、`interface.response` 注解中的`interface`应该和接口名字一一对应；`interface.request`、`interface.response`内嵌结构也需做到`bff`包下全局唯一

## Service

- [ ] Mutil service

- [ ] Annotation

  from xxx.proto

  - interface

    ```go
    // @type: service.interface
    // @path: demo
    // @des: Demo
    service Demo {
    	rpc	Ping(PingRequest) returns (PingResponse);
    }
    ```

  - interface.request

    ```go
    // @type: service.interface.request
    // @des: Demo
    message PingRequest {}
    ```

  - interface.response

    ```go
    // @type: service.interface.response
    message PingResponse {}
    ```

## Protocol

```sequence
Protocol Buffer ->PB: Generate
Protocol Buffer ->Service: Generate
```



### protocol buffer

```yaml
protcol
	- xxx_demo_system
		- demo
			- ping.proto
			- login.proto
		- bar
			- ping.proto
	- xxx_foo_system
		- demo
			- ping.proto
			- login.proto
```

### pb

```
	+ pb
		- demo
		  - ping.go
		  - login.go
		- bar
			- ping.go
```

### PB



## SQL parser



