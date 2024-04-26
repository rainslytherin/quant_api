package models

import (
	"quant_api/database"
)

type SyncStatus struct {
	ID         int    `json:"id" gorm:"primary_key"`
	ClientID   string `json:"client_id"`
	UpdateTime int64  `json:"update_time"`
}

func (s SyncStatus) Create() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	res, err := db.Exec("insert into sync_status (client_id, update_time) values (?, ?)", s.ClientID, s.UpdateTime)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)
	return nil
}

func (s SyncStatus) Save() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("update sync_status set update_time = ? where client_id = ?", s.UpdateTime, s.ClientID)
	return err
}

func LoadSyncStatus(clientID string) (syncStatus *SyncStatus, err error) {
	db, err := database.GetGlobalDB()
	if err != nil {
		return syncStatus, err
	}

	if err := db.Select(syncStatus, "select * from sync_status where client_id = ?", clientID); err != nil {
		return syncStatus, err
	}

	if syncStatus != nil {
		return syncStatus, nil
	}

	// create new sync status and save
	syncStatus = &SyncStatus{
		ClientID:   clientID,
		UpdateTime: 0,
	}

	if err := syncStatus.Create(); err != nil {
		return syncStatus, err
	}

	return syncStatus, nil
}

func LoadAllSyncStatus() (syncStatus []*SyncStatus, err error) {
	db, err := database.GetGlobalDB()
	if err != nil {
		return syncStatus, err
	}

	if err := db.Select(syncStatus, "select * from sync_status"); err != nil {
		return syncStatus, err
	}

	return syncStatus, nil
}
