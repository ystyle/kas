package model

import (
	"github.com/labstack/gommon/log"
	"time"
)

type Drive struct {
	ID        string `storm:"id"`
	Count     uint
	CreatedAt time.Time
	Last      time.Time
}

func Statistics(driveid string) {
	if driveid == "" {
		driveid = "unknow"
	}
	var drive Drive
	err := store.One("ID", driveid, &drive)
	if err != nil {
		log.Error(err)
		return
	}
	if drive.ID == "" {
		drive.ID = driveid
		drive.CreatedAt = time.Now()
		drive.Last = time.Now()
		drive.Count = 1
		store.Save(&drive)
	} else {
		drive.Count += 1
		store.Update(&drive)
	}

}
