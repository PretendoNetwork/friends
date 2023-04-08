package database

import "github.com/PretendoNetwork/friends-secure/globals"

func initPostgres3DS() {
	var err error

	_, err = Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS "3ds"`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[3DS] Postgres schema created")

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS "3ds".user_data (
		pid integer PRIMARY KEY,
		show_online boolean DEFAULT true,
		show_current_game boolean DEFAULT true,
		comment text,
		comment_changed bigint,
		last_online bigint,
		favorite_title bigint,
		favorite_title_version integer,
		mii_name text,
		mii_data bytea,
		mii_changed bigint,
		region integer,
		area integer,
		language integer
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS "3ds".friendships (
		id bigserial PRIMARY KEY,
		user1_pid integer,
		user2_pid integer,
		type integer,
		date bigint,
		UNIQUE (user1_pid, user2_pid)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("[3DS] Postgres tables created")
}
