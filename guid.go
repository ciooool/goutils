package goutils

import (
	"fmt"
	"sync"
	"time"
)

/*
生成全局唯一ID的雪花算法原理
	雪花算法是一种用于生成全局唯一ID的算法，最初由Twitter开发，用于解决分布式系统中生成ID的问题。
	其核心思想是将一个64位的长整型ID划分成多个部分，每个部分用于表示不同的信息，确保了生成的ID在分布式环境下的唯一性。

ID结构
	符号位（1位）：始终为0，用于保证ID为正数。
	时间戳（41位）：表示生成ID的时间戳，精确到毫秒级。
	工作节点ID（10位）：表示生成ID的机器的唯一标识。
	序列号（12位）：表示在同一毫秒内生成的多个ID的序列号。
生成步骤
	获取当前时间戳，精确到毫秒级。
	如果当前时间小于上次生成ID的时间，或者在同一毫秒内生成的ID数量超过最大值，等待下一毫秒再继续生成。
	如果当前时间等于上次生成ID的时间，序列号自增1。
	如果当前时间大于上次生成ID的时间，序列号重新从0开始。
	将各个部分的值组合，得到最终的64位ID。
*/

const (
	workerBits  = 10
	seqBits     = 12
	workerMax   = -1 ^ (-1 << workerBits)
	seqMask     = -1 ^ (-1 << seqBits)
	timeShift   = workerBits + seqBits
	workerShift = seqBits
	epoch       = 1609459200000
)

type Snowflake struct {
	lastTime int64
	workerID int64
	sequence int64
	mu       sync.Mutex
}

func NewSnowflake(workerID int64) *Snowflake {
	if workerID < 0 || workerID > workerMax {
		panic(fmt.Sprintf("worker ID must be between 0 and %d", workerMax))
	}

	return &Snowflake{
		lastTime: time.Now().UnixNano() / 1e6,
		workerID: workerID,
		sequence: 0,
	}
}

func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := time.Now().UnixNano() / 1e6
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) % seqMask

		if s.sequence == 0 {
			for currentTime <= s.lastTime {
				currentTime = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = currentTime

	id := (currentTime-epoch)<<timeShift | s.workerID<<workerShift | s.sequence
	return id
}
