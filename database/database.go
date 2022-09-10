package database

func Connect() {
	connectMongo()
	connectCassandra()
}
