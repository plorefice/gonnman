package connman

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/godbus/dbus"
)

type DBusInterface struct {
	Connection *dbus.Conn
	Object     dbus.BusObject
	Interface  string
}

func (db *DBusInterface) Call(name string, args ...interface{}) ([]interface{}, error) {
	call := db.Object.Call(db.Interface+"."+name, 0, args...)
	return call.Body, call.Err
}

func (db *DBusInterface) Set(name string, value interface{}) error {
	v := dbus.MakeVariant(value)
	return db.Object.Call(db.Interface+".SetProperty", 0, name, v).Err
}

func (db *DBusInterface) Get(name string) (interface{}, error) {
	call := db.Object.Call(db.Interface+".GetProperties", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	props := call.Body[0].(map[string]dbus.Variant)
	if prop, ok := props[name]; !ok {
		return nil, errors.New("Invalid property")
	} else {
		return prop.Value(), nil
	}
}

func (db *DBusInterface) Done() error {
	return db.Connection.Close()
}

func DBus(service string, path dbus.ObjectPath, ifname string) (*DBusInterface, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	obj := conn.Object(service, path)
	if obj == nil {
		conn.Close()
		return nil, err
	}

	return &DBusInterface{
		Connection: conn,
		Object:     obj,
		Interface:  ifname,
	}, nil
}

func DBusClock() (*DBusInterface, error) {
	return DBus("net.connman", "/", "net.connman.Clock")
}

func DBusWifi() (*DBusInterface, error) {
	return DBus(
		"net.connman",
		"/net/connman/technology/wifi",
		"net.connman.Technology")
}

func DBusManager() (*DBusInterface, error) {
	return DBus("net.connman", "/", "net.connman.Manager")
}

func DBusService(svc dbus.ObjectPath) (*DBusInterface, error) {
	return DBus("net.connman", svc, "net.connman.Service")
}

func DBusTechnology(tech dbus.ObjectPath) (*DBusInterface, error) {
	return DBus("net.connman", tech, "net.connman.Technology")
}

/*
 *	D-Bus to Go mapping functions
 */

func setField(dst interface{}, key string, val dbus.Variant) error {
	key = strings.Replace(key, ".", "", -1)

	sv := reflect.ValueOf(dst).Elem()
	sfv := sv.FieldByName(key)

	if !sfv.IsValid() {
		return fmt.Errorf("No such field %s in structure", key)
	}

	if !sfv.CanSet() {
		return fmt.Errorf("Cannot set %s field value", key)
	}

	sft := sfv.Type()
	v := reflect.ValueOf(val.Value())
	vt := reflect.TypeOf(val.Value())

	switch vt {
	case reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(val)):
		subf := reflect.New(sft)
		if err := dictToStruct(v.Interface().(map[string]dbus.Variant), subf.Interface()); err != nil {
			return err
		}
		sfv.Set(reflect.Indirect(subf))

	default:
		if sft != vt {
			return fmt.Errorf("Value type (%v) does not match field type (%s) : %v",
				vt, sft.Name(), v)
		}
		sfv.Set(v)
	}

	return nil
}

func dictToStruct(dict map[string]dbus.Variant, dst interface{}) error {
	for k, v := range dict {
		err := setField(dst, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func structToDict(src interface{}) (map[string]dbus.Variant, error) {
	st := reflect.TypeOf(src)
	sv := reflect.ValueOf(src)

	for st.Kind() == reflect.Ptr {
		sv = reflect.Indirect(reflect.ValueOf(src))
		st = sv.Type()
	}

	if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Expected struct, found %v", st.Kind())
	}

	ret := make(map[string]dbus.Variant)

	for i := 0; i < sv.NumField(); i++ {
		sft := st.Field(i).Type
		sfn := st.Field(i).Name
		sfv := sv.Field(i).Interface()

		switch sft.Kind() {
		case reflect.Struct:
			sub, err := structToDict(sfv)
			if err != nil {
				return nil, err
			}
			ret[sfn] = dbus.MakeVariant(sub)
			break

		default:
			ret[sfn] = dbus.MakeVariant(sfv)
		}
	}

	return ret, nil
}
