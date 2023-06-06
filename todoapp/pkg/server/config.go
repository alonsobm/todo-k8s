package server

type Config struct {
	Environment      string `json:"environment"`
	DatabaseHost     string `json:"database_host"`
	DatabasePort     int    `json:"database_port"`
	DatabaseName     string `json:"database_name"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	AuthKey          string `json:"auth_key"`
	Port             string `json:"port"`
}
