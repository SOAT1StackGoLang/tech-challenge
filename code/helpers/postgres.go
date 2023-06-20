package helpers

//
//var (
//	postgresUsername, postgresPassword, postgresHost string
//	postgresPort                                     int
//)
//
//type PostgresParams struct {
//	Host, User, Pwd, DbName string
//	Port                    int
//}
//
//func InitPostgresWithFlagSet(flagSet *flag.FlagSet) {
//	flagSet.StringVar(&postgresUsername, "postgresusr", GetEnvOrDefault("DB_USER", "postgres"), "")
//	flagSet.StringVar(&postgresUsername, "postgrespass", GetEnvOrDefault("DB_PASSWORD", "postgres"), "")
//	flagSet.StringVar(&postgresUsername, "postgresusr", GetEnvOrDefault("DB_USER", "postgres"), "")
//	flagSet.StringVar(&postgresUsername, "postgresusr", GetEnvOrDefault("DB_USER", "postgres"), "")
//}
//
//func GetEnvOrDefault(key, defaultValue string) string {
//	val := os.Getenv(key)
//	if val == "" {
//		return defaultValue
//	}
//	return val
//}
//
//func GetEnvOrDefault(key string, defaultValue int) int {
//	val := os.Getenv(key)
//	if val == "" {
//		return defaultValue
//	}
//	i, _ := strconv.ParseInt(val, 10, 32)
//	return int(i)
//}
