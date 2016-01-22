package connman

import "github.com/godbus/dbus"

func GetServices() ([]*Service, error) {
	var resp []struct {
		Path dbus.ObjectPath
		Map  map[string]dbus.Variant
	}

	db, err := DBusManager()
	if err != nil {
		return nil, err
	}

	ret, err := db.Call("GetServices")
	if err != nil {
		return nil, err
	}
	dbus.Store(ret, &resp)

	var svcs []*Service
	for _, s := range resp {
		var svc = &Service{}

		svc.Path = s.Path
		if err := dictToStruct(s.Map, svc); err != nil {
			return nil, err
		}
		svcs = append(svcs, svc)
	}

	return svcs, nil
}

func GetTechnologies() ([]*Technology, error) {
	var resp []struct {
		Path dbus.ObjectPath
		Map  map[string]dbus.Variant
	}

	db, err := DBusManager()
	if err != nil {
		return nil, err
	}

	ret, err := db.Call("GetTechnologies")
	if err != nil {
		return nil, err
	}
	dbus.Store(ret, &resp)

	var techs []*Technology
	for _, t := range resp {
		var tech = &Technology{}

		tech.Path = t.Path
		if err := dictToStruct(t.Map, tech); err != nil {
			return nil, err
		}
		techs = append(techs, tech)
	}

	return techs, nil
}

func RegisterAgent(a *Agent) error {
	db, err := DBusManager()
	if err != nil {
		return err
	}

	_, err = db.Call("RegisterAgent", a.Path)
	return err
}

func UnregisterAgent(a *Agent) error {
	db, err := DBusManager()
	if err != nil {
		return err
	}

	_, err = db.Call("UnregisterAgent", a.Path)
	return err
}
