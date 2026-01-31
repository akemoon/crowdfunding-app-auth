package config

type Config struct {
	Port           int      `json:"port"`
	JwtSecret      string   `json:"jwt_secret"`
	ClientSecret   string   `json:"client_secret"`
	Databases      Database `json:"databases"`
	ClusterClients Cluster  `json:"cluster"`
}

type Database struct {
	Redis Redis `json:"redis"`
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Cluster struct {
	UserService string `json:"user_service_url"`
}
