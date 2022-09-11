package database

func Connect() {
	connectMongo()
	connectPostgres()
}
