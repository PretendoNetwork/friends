package database_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go"
)

// SetUserBlocked marks a blocked PID as blocked on a bloker PID block list
func SetUserBlocked(blockerPID uint32, blockedPID uint32, titleId uint64, titleVersion uint16) error {
	date := nex.NewDateTime(0)
	date.FromTimestamp(time.Now())

	_, err := database.Postgres.Exec(`
	INSERT INTO wiiu.blocks (blocker_pid, blocked_pid, title_id, title_version, date)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (blocker_pid, blocked_pid)
	DO UPDATE SET
	date = $5`, blockerPID, blockedPID, titleId, titleVersion, date.Value())
	if err != nil {
		return err
	}

	return nil
}
