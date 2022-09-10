package globals

import (
	"github.com/PretendoNetwork/plogger-go"
	"github.com/bwmarrin/snowflake"
)

var Logger = plogger.NewLogger()
var ConnectedUsers map[uint32]*ConnectedUser
var SnowflakeNodes []*snowflake.Node
