package mysql

type DBType string

const (
	Reader DBType = "reader"
	Writer DBType = "writer"
)
