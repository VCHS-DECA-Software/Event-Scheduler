package users

import (
	"main/components/dbmanager"
	"main/components/encryption"

	"main/components/events"

	uuid "github.com/satori/go.uuid"
)

type Admin struct {
	ID       string `storm:"id"`
	Username string `storm:"unique"`
	Name     string
	Password string
}

func CreateAdmin(username, name, password string) error {
	id := uuid.NewV4().String()
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return err
	}
	err = dbmanager.Save(&Admin{ID: id, Username: username, Name: name, Password: hashedPassword})
	return err
}

func GetAdmin(id string) (Admin, error) {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	return admin, err
}

func UpdateAdmin(id, username, name, password string) error {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	if err != nil {
		return err
	}
	admin.Username = username
	admin.Name = name
	hashedPassword, _ := encryption.HashPassword(password)
	admin.Password = hashedPassword
	err = dbmanager.Update(&admin)
	return err
}

func DeleteAdmin(id string) error {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&admin)
	return err
}

func (admin Admin) CreateEvent(name, description, eventType, startTime, endTime, location string) error {
	id := uuid.NewV4().String()
	err := dbmanager.Save(&events.Event{ID: id, Name: name, AdminID: admin.ID, Description: description, EventType: eventType, StartTime: startTime, EndTime: endTime, Location: location})
	return err
}

func (admin Admin) DeleteEvent(id string) error {
	var event events.Event
	err := dbmanager.Query("ID", id, &event)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&event)
	return err
}
