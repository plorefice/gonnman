package connman

import "github.com/godbus/dbus"

type Technology struct {
	Path      dbus.ObjectPath
	Name      string
	Type      string
	Powered   bool
	Connected bool
	Tethering bool
}

func (t *Technology) Enable() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}
	return db.Set("Powered", true)
}

func (t *Technology) Disable() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}
	return db.Set("Powered", false)
}

func (t *Technology) Scan() error {
	db, err := DBusTechnology(t.Path)
	if err != nil {
		return err
	}

	_, err = db.Call("Scan")
	return err
}
