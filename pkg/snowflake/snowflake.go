package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

// Snowflake 雪花算法
type Snowflake struct {
	Node *sf.Node
}

func NewSnowflake(startTime string, machineID int64) (*Snowflake, error) {
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.L().Debug("解析时间失败", zap.String("startTime", startTime), zap.Error(err))
		return nil, err
	}
	nano := parse.UnixNano()
	sf.Epoch = nano // 1000000
	node, err := sf.NewNode(machineID)
	if err != nil {
		zap.L().Debug("创建雪花算法节点失败", zap.Int64("machineID", machineID), zap.Error(err))
		return nil, err
	}
	s := &Snowflake{
		Node: node,
	}
	return s, nil
}
func (s *Snowflake) GetID() int64 {
	return s.Node.Generate().Int64()
}
