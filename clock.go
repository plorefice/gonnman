package connman

import (
	"fmt"
	"os/exec"
	"time"
)

type Clock struct {
	Method  string `json:"method"`
	Hours   int    `json:"hours"`
	Minutes int    `json:"mins"`
	Year    int    `json:"year"`
	Month   int    `json:"month"`
	Day     int    `json:"day"`
}

func (c *Clock) FromTime(t time.Time) {
	var m time.Month

	c.Hours, c.Minutes, _ = t.Clock()
	c.Year, m, c.Day = t.Date()
	c.Month = int(m)
}

func (c *Clock) ToTime() (time.Time, error) {
	str := fmt.Sprintf("%v/%v/%v %v:%v", c.Day, c.Month, c.Year, c.Hours, c.Minutes)
	tm, err := time.Parse("2/1/2006 15:4", str)
	return tm, err
}

func (c *Clock) Parse(hours, mins, year, month, day string) error {
	str := fmt.Sprintf("%v/%v/%v %v:%v", day, month, year, hours, mins)
	tm, err := time.Parse("2/1/2006 15:4", str)
	if err != nil {
		return err
	}
	c.FromTime(tm)
	return nil
}

func GetTime() *Clock {
	c := &Clock{}
	c.FromTime(time.Now())
	c.Method, _ = GetTimeMethod()
	return c
}

func SetManualTime(c Clock) error {
	time, err := c.ToTime()
	if err != nil {
		return err
	}

	ck, err := DBusClock()
	if err != nil {
		return err
	}
	if err := ck.Set("TimeUpdates", "manual"); err != nil {
		return err
	}
	if err := ck.Set("Time", uint64(time.Unix())); err != nil {
		return err
	}

	// Save the time to HW clock
	return exec.Command("hwclock", "-w").Run()
}

func SetAutoTime() error {
	ck, err := DBusClock()
	if err != nil {
		return err
	}

	err = ck.Set("TimeUpdates", "auto")
	if err != nil {
		return err
	}

	// Force time sync
	exec.Command("systemctl", "stop", "ntpd").Run()
	exec.Command("ntpdate", "-s", "0.develer.pool.ntp.org").Run()
	exec.Command("systemctl", "start", "ntpd").Run()

	return nil
}

func GetTimeMethod() (string, error) {
	ck, err := DBusClock()
	if err != nil {
		return "", nil
	}
	method, err := ck.Get("TimeUpdates")
	return method.(string), err
}
