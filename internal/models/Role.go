package models

type Role string

const (
	ADMIN  Role = "admin"
	WORKER Role = "worker"
	ALL    Role = "all"
)

var ADMIN_ROLE = []Role{ADMIN}
var WORKER_ROLE = []Role{WORKER}
var ALL_ROLE = []Role{ALL}
