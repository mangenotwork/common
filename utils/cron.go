package utils

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	cron "github.com/robfig/cron/v3"
)

// 定时任务 https://crontab.guru/

type Cron struct {
	C *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		C: cron.New(),
	}
}

func (c *Cron) Run() {
	c.C.Start()
}

func (c *Cron) add(spec string, f func()) {
	enterId, err := c.C.AddFunc(spec, f)
	if err != nil {
		panic(err)
	}
	log.InfoF("添加任务id是 %d \n", enterId)
}

// AddAtMinute 每多少分钟 i=0为每分钟
func (c *Cron) AddAtMinute(f func(), i int) {
	spec := "* * * * *"
	if i > 0 && i < 60 {
		spec = fmt.Sprintf("*/%d * * * *", i)
	}
	c.add(spec, f)
}

// AddAtHours 每多少小时 i=0为每小时
func (c *Cron) AddAtHours(f func(), i int) {
	spec := "0 * * * *"
	if i > 0 && i < 24 {
		spec = fmt.Sprintf("0 */%d * * *", i)
	}
	c.add(spec, f)
}

// AddAtDayWhatTime 每天的几点 i=0为每天
func (c *Cron) AddAtDayWhatTime(f func(), i int) {
	spec := "0 0 * * *"
	if i > 0 && i < 24 {
		spec = fmt.Sprintf("0 %d * * *", i)
	}
	c.add(spec, f)
}

// AddAtSunday  每周日
func (c *Cron) AddAtSunday(f func()) {
	c.add("0 0 * * SUN", f)
}

// AddAtMonday  每周一
func (c *Cron) AddAtMonday(f func()) {
	c.add("0 0 * * MON", f)
}

// AddAtTuesday 每周二
func (c *Cron) AddAtTuesday(f func()) {
	c.add("0 0 * * TUE", f)
}

// AddAtWednesday 每周三
func (c *Cron) AddAtWednesday(f func()) {
	c.add("0 0 * * WED", f)
}

// AddAtThursday 每周四
func (c *Cron) AddAtThursday(f func()) {
	c.add("0 0 * * THU", f)
}

// AddAtFriday 每周五
func (c *Cron) AddAtFriday(f func()) {
	c.add("0 0 * * FRI", f)
}

// AddAtSaturday 每周六
func (c *Cron) AddAtSaturday(f func()) {
	c.add("0 0 * * SAT", f)
}
