package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func (idw *IdWorker) InitIdWorker(workerId, datacenterId int64) error {
	var baseValue int64 = -1
	idw.startTime = 1463834116272
	idw.workerIdBits = 5
	idw.datacenterIdBits = 5
	idw.maxWorkerId = baseValue ^ (baseValue << idw.workerIdBits)
	idw.maxDatacenterId = baseValue ^ (baseValue << idw.datacenterIdBits)
	idw.sequenceBits = 12
	idw.workerIdLeftShift = idw.sequenceBits
	idw.datacenterIdLeftShift = idw.workerIdBits + idw.workerIdLeftShift
	idw.timestampLeftShift = idw.datacenterIdBits + idw.datacenterIdLeftShift
	idw.sequenceMask = baseValue ^ (baseValue << idw.sequenceBits)
	idw.sequence = 0
	idw.lastTimestamp = -1
	idw.signMask = ^baseValue + 1
	idw.idLock = &sync.Mutex{}
	if idw.workerId < 0 || idw.workerId > idw.maxWorkerId {
		return fmt.Errorf("workerId[%v] is less than 0 or greater than maxWorkerId[%v].",
			workerId, datacenterId)
	}
	if idw.datacenterId < 0 || idw.datacenterId > idw.maxDatacenterId {
		return fmt.Errorf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d].",
			workerId, datacenterId)
	}
	idw.workerId = workerId
	idw.datacenterId = datacenterId
	return nil
}

// NextId 返回一个唯一的 INT64 ID
func (idw *IdWorker) NextId() (int64, error) {
	idw.idLock.Lock()
	timestamp := time.Now().UnixNano()
	if timestamp < idw.lastTimestamp {
		return -1, fmt.Errorf(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds",
			idw.lastTimestamp-timestamp))
	}
	if timestamp == idw.lastTimestamp {
		idw.sequence = (idw.sequence + 1) & idw.sequenceMask
		if idw.sequence == 0 {
			timestamp = idw.tilNextMillis()
			idw.sequence = 0
		}
	} else {
		idw.sequence = 0
	}
	idw.lastTimestamp = timestamp
	idw.idLock.Unlock()
	id := ((timestamp - idw.startTime) << idw.timestampLeftShift) |
		(idw.datacenterId << idw.datacenterIdLeftShift) |
		(idw.workerId << idw.workerIdLeftShift) |
		idw.sequence
	if id < 0 {
		id = -id
	}
	return id, nil
}

// tilNextMillis
func (idw *IdWorker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= idw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}

func ID64() (int64, error) {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return 0, err
	}
	return currWorker.NextId()
}

func ID() int64 {
	id, _ := ID64()
	return id
}

func IDStr() string {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return ""
	}
	id, err := currWorker.NextId()
	if err != nil {
		return ""
	}
	return Int64ToStr(id)
}

func IDMd5() string {
	return Get16MD5Encode(IDStr())
}

func IDShort() int64 {
	currWoker := &IdWorker{}
	currWoker.InitIdWorker(1000, 2)
	a, _ := currWoker.NextId()
	a1 := strconv.FormatInt(a, 10)
	sa := a1[len(a1)-5:]
	b := time.Now().UnixNano()
	b1 := strconv.FormatInt(b, 10)
	sb := b1[len(b1)-10 : len(b1)-6]
	sc := sa + sb
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(147483647)
	ids, _ := strconv.Atoi(sc)
	return int64(ids + x)
}
