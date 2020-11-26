package handlers

const (
	// 中间件服务
	MiddlewareConfig    = "config"
	MiddlewareLoginDB = "test"
	MiddlewareLoginDBORM = "dborm"
	MiddlewareLoginDBSQLX = "dbsqlx"
	MiddlewareLoginREDIS = "redis"
)

//Json返回结果
type Result struct {
	Message string `json:"msg"`
	Status int `json:"status"`
}
