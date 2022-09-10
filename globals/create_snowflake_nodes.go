package globals

import (
	"runtime"

	"github.com/bwmarrin/snowflake"
)

func CreateSnowflakeNodes() {
	SnowflakeNodes = make([]*snowflake.Node, 0)

	for corenum := 0; corenum < runtime.NumCPU(); corenum++ {
		node, err := snowflake.NewNode(int64(corenum))
		if err != nil {
			// TODO: Handle error
			Logger.Critical(err.Error())
			return
		}
		SnowflakeNodes = append(SnowflakeNodes, node)
	}
}
