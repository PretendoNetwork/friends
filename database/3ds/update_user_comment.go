package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// UpdateUserComment updates a user's comment
func UpdateUserComment(pid uint32, message string) error {
	changed := types.NewDateTime(0).Now()

	_, err := database.Manager.Exec(`
		INSERT INTO "3ds".user_data (pid, comment, comment_changed)
		VALUES ($1, $2, $3)
		ON CONFLICT (pid)
		DO UPDATE SET 
		comment = $2,
		comment_changed = $3`, pid, message, uint64(changed))

	if err != nil {
		return err
	}

	return nil
}
