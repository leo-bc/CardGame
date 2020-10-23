package model

// Updateable :
type Updateable struct {
	UpdateID int
}

// UpdateInfo :
type UpdateInfo struct {
	IsUpdated bool
	NewID     int
}

// SetUpdated :
func (updateable *Updateable) SetUpdated() {
	updateable.UpdateID++
}

// TriggerUpdate :
func (updateable *Updateable) TriggerUpdate(updateID int) UpdateInfo {
	if updateable.UpdateID == updateID {
		return UpdateInfo{false, updateable.UpdateID}
	}
	return UpdateInfo{true, updateable.UpdateID}
}
