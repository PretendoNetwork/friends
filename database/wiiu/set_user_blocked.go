package database_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
)

func SetUserBlocked(blockerPID uint32, blockedPID uint32, titleId uint64, titleVersion uint16) {
	date := nex.NewDateTime(0)
	date.FromTimestamp(time.Now())

	_, err := database.Postgres.Exec(`
	INSERT INTO wiiu.blocks (blocker_pid, blocked_pid, title_id, title_version, date)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (blocker_pid, blocked_pid)
	DO UPDATE SET
	date = $5`, blockerPID, blockedPID, titleId, titleVersion, date.Value())
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
