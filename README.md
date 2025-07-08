## 用于生成 gf 的 代码

0. 安装 `gf` 命令行, 用于后续代码生成

	```shell
	go install github.com/gogf/gf/cmd/gf/v2@latest
	```

1. 使用 `gf init myapp` 初始化项目

	生成的目录结构如下

	```shell
	tree -L 2 myapp

	myapp
	├── api
	│   └── hello
	├── go.mod
	├── go.sum
	├── hack
	│   ├── config.yaml
	│   ├── hack-cli.mk
	│   └── hack.mk
	├── internal
	│   ├── cmd
	│   ├── consts
	│   ├── controller
	│   ├── dao
	│   ├── logic
	│   ├── model
	│   ├── packed
	│   └── service
	├── main.go
	├── Makefile
	├── manifest
	│   ├── config
	│   ├── deploy
	│   ├── docker
	│   ├── i18n
	│   └── protobuf
	├── README.MD
	├── resource
	│   ├── public
	│   └── template
	└── utility

	```

2. 创建好数据库，并创建好表. 比如 `user`

    ```sql
    CREATE TABLE `user` (
        `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
        `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
        `deleted_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '软删除',
        PRIMARY KEY (`id`)
    ) COMMENT='用户表';
    ```

3. 在 `hack/config.yaml` 中配置好数据库连接和数据库表信息

	```yaml
	gfcli:
	  gen:
	    dao:
	    - link: "mysql:root:123456@tcp(127.0.0.1:3306)/test"
	      descriptionTag: true
	      tables: "user" # 这里可指定要生成的表. 其他参数参考官方文档 https://goframe.org/docs/cli/gen-dao
	```

4. 生成 dao 层代码

	```shell
    # 使用官方 gf 命令行生成
    make dao # 或者 `gf gen dao` 
    ```

    结果变动如下

    ```shell
    tree internal/model

    internal/model
    ├── do
    │   └── user.go
    └── entity
        └── user.go

    tree internal/dao
    internal/dao
    ├── internal
    │   └── user.go
    └── user.go

    ```

5. 安装本工具

	```shell
	go install github.com/magalab/gfx@v0.0.1
	```

6. 补充生成 `model` 内 `user.go`

	```shell
    # 使用本工具生成
	gfx gen model -s user
	```

	结果如下. 目前仅生成模板, 需手动补充实际业务字段. 后续再补充自动生成字段的功能

	```go
	// internal/model/user.go
	package model

	import (
	    "myapp/internal/model/entity"
	)

	type UserModel struct {
	    entity.User
	}

	type UserAddReq struct {
	}

	type UserDeleteReq struct {
        UserId uint64 `json:"user_id"`
	}

	type UserUpdateReq struct {
        UserId uint64 `json:"user_id"`
	}

	type UserListReq struct {
        PageNum  int `json:"page_num"`
	    PageSize int `json:"page_size"`
	}

	type UserGetReq struct {
        UserId uint64 `json:"user_id"`
	}

	type UserItem struct {
	}

	func (m *UserModel) ToItem() *UserItem {

	    return &UserItem{}
	}
	```

7. 补充生成 `dao` 层的相关查询方法

    ```shell
    # 使用本工具生成, 手动按需删除或更新
    gfx gen dao -s user
    ```

    结果变动如下
    其中 `UserId`, `DeletedAt` 字段来自数据库定义. `DeletedAt` 一般定义为 `int64` 类型, 用于处理软删除场景. goframe 框架在查询时会自动过滤

    ```go
    // internal/dao/user.go

    func (d *userDao) GetByUserId(ctx context.Context, req *model.UserGetReq) (*model.UserModel, error) {
        var item *model.UserModel
        if err := d.Ctx(ctx).
            Where(d.Columns().UserId, req.UserId).
            Scan(&item); err != nil {
            return nil, err
        }
        return item, nil
    }

    // UpdateByUserId
    func (d *userDao) UpdateByUserId(ctx context.Context, req *model.UserUpdateReq) error {
        updater := g.Map{
            // TODO
        }
        if _, err := d.Ctx(ctx).
            Where(d.Columns().UserId, req.UserId).
            UpdateAndGetAffected(updater); err != nil {
            return err
        }
        return nil
    }

    // DeleteByUserId
    func (d *userDao) DeleteByUser(ctx context.Context, req *model.UserDeleteReq) error {
        updater := g.Map{
            d.Columns().DeletedAt: gtime.Now().Unix(),
        }
        if _, err := d.Ctx(ctx).
            Where(d.Columns().UserId, req.UserId).
            UpdateAndGetAffected(updater); err != nil {
            return err
        }
        return nil
    }

    // UserList
    func (d *userDao) UserList(ctx context.Context, req *model.UserListReq) ([]*model.UserModel, int, error) {
        db := d.Ctx(ctx)
        count, err := db.Count()
        if err != nil {
            return nil, 0, err
        }
        if count == 0 {
            return nil, 0, nil
        }
        var items []*model.UserModel
        // TODO order
        if err = db.Page(req.PageNum, req.PageSize).
            Scan(&items); err != nil {
            return nil, 0, err
        }

        return items, count, nil
    }
    ```

7. 生成 `api` 接口定义

	```shell
    # 使用本工具生成, 手动按需删除或更新
	gfx gen api -s user
	```

	结果如下. 其中 `model.UserXXXReq` 来自上一步骤的 `model` 补充定义
	```go
	// api/user/v1/user.go
	package v1

	import (
	    "github.com/gogf/gf/v2/frame/g"

	    "myapp/internal/model"
	)

	type UserAddReq struct {
	    g.Meta `path:"/user" tags:"User" method:"post" sm:"新增"`
	    *model.UserAddReq
	}
	type UserAddRes struct {
	}

	type UserGetReq struct {
	    g.Meta `path:"/user" tags:"User" method:"get" sm:"单条"`
	    *model.UserGetReq
	}
	type UserGetRes struct {
	    *model.UserItem
	}

	type UserListReq struct {
	    g.Meta `path:"/users" tags:"User" method:"get" sm:"列表"`
	    *model.UserListReq
	}
	type UserListRes struct {
	    Items []*model.UserItem `json:"items"`
	    Total int               `json:"total"`
	}

	type UserUpdateReq struct {
	    g.Meta `path:"/user" tags:"User" method:"put" sm:"更新"`
	    *model.UserUpdateReq
	}
	type UserUpdateRes struct {
	}

	type UserDeleteReq struct {
	    g.Meta `path:"/user" tags:"User" method:"delete" sm:"删除"`
	    *model.UserDeleteReq
	}
	type UserDeleteRes struct {
	}

	```

8. 生成 `controller` 层代码

    ```shell
    # 使用官方 `gf` 命令行生成
    make ctrl # 或者 gf gen ctrl
    ```

    相关变动. 会生成的 `controller` 文件以及 `api` 中补充的接口定义
    ```shell
    tree api/user -L 1
    api/user
    ├── user.go # 这个
    └── v1

    tree internal/controller

    internal/controller/user
    ├── user_new.go
    ├── user_v1_user_add.go
    ├── user_v1_user_delete.go
    ├── user_v1_user_get.go
    ├── user_v1_user_list.go
    ├── user_v1_user_update.go
    └── user.go
    ```

9. 生成 `logic` 层代码

    ```shell
    # 使用本工具生成, 手动按需删除或更新
    gfx gen logic -s user
    ```

    相关变动. 仅模板结构, 具体逻辑自行实现. `dao` 层的查询之前的步骤已生成简单场景的
    ```shell
    tree internal/logic/user

    internal/logic/user
    ├── user_add.go
    ├── user_delete.go
    ├── user_get.go
    ├── user_list.go
    ├── user_update.go
    └── user.go # 这里会有个报错, 因为 service 层还没有生成
    ```

10. 生成 `service` 层代码

    ```shell
    # 使用官方 `gf` 命令行生成
    make service # 或者 gf gen service
    ```

    相关变动

    ```shell
    tree internal/service
    internal/service
    └── user.go

    tree internal/logic

    internal/logic
    ├── logic.go # 这个
    └── user
    ```

11. 更新 `ctrl` 层. 调用 `service` 层

    ```go
    // internal/controller/user/user_v1_user_get.go
    func (c *ControllerV1) UserGet(ctx context.Context, req *v1.UserGetReq) (res *v1.UserGetRes, err error) {
        // 这里可能要处理 req 中的 userId 的来源问题. 比如登录时注入到 ctx 中
        // req.UserId = 1
        item, err := service.User().UserGet(ctx, req.UserGetReq)
        if err != nil {
            return nil, err
        }

        return &v1.UserGetRes{UserItem: item}, nil
    }
    ```

12. 新增路由定义

    ```go
    // internal/cmd/cmd.go

    import (

        // ...
        // 引入 mysql 驱动包. 或者在外层的 main.go 中引入
        _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	    // _ "github.com/gogf/gf/contrib/nosql/redis/v2"

        "myapp/internal/controller/user" // 智能补全别补错, 需要是 controller 包下的
        // ...
    )

    // ...
    s := g.Server()
    // xxx 用于区分业务
    // v1Group := s.Group("/api/xxx/v1")
    v1Group := s.Group("/api/v1")
    v1Group.Middleware(
        // ...
    )
    v1Group.Bind(
        user.NewV1(), // 批量注册路由
        
    )

    s.Run()
    // ...

    ```

13. 调整运行时配置文件

    ```yaml
    # manifest/config/config.yaml
    server:
        address:     ":8000"
        openapiPath: "/api.json"
        logPath: "logs"
        accessLogEnabled: true
        accessLogPattern: "access.{Ymd}.log"
        mode: "dev"
    logger:
        path: "logs"
        level : "all"
        stdout: true
        StStatus: 1
        stack: false
        file: "{Y-m-d}.log"
        rotateSize: "10M"
        rotateBackupLimit: 2
        rotateBackupExpire: "7d"
        rotateBackupCompress:  9 
    database:
        logger:
            path:    "logs"
            file: "db-{Y-m-d}.log"
            level:   "all"
            stdout:  true
            rotateSize: "10M"
            rotateBackupLimit: 2
            rotateBackupExpire: "7d"
            rotateBackupCompress:  9
        default:
        - link:  "mysql:root:123456@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
          debug: true
    ```

14. 启动

    ```shell
    go run .
    ```

15. 重复上述步骤新增功能. 其中本工具命令的功能仅支持生成一次