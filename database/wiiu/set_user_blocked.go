package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// SetUserBlocked marks a blocked PID as blocked on a blocker PID block list
func SetUserBlocked(blockerPID uint32, blockedPID uint32, titleID uint64, titleVersion uint16) error {
	date := types.NewDateTime(0).Now()

	_, err := database.Manager.Exec(`
	INSERT INTO wiiu.blocks (blocker_pid, blocked_pid, title_id, title_version, date)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (blocker_pid, blocked_pid)
	DO UPDATE SET
	date = $5`, blockerPID, blockedPID, titleID, titleVersion, date.Value())
	if err != nil {
		return err
	}

	return nil
}
