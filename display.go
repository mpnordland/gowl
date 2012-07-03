
package gowl

import (
	"bytes"
)

type Display struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Display, msg []byte)
}

//// Requests
func (d *Display) Bind (name uint32, iface string, version uint32, id Object ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, name)
	writeString(buf, []byte(iface))
	writeInteger(buf, version)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(d, 0, buf.Bytes())
}

func (d *Display) Sync (callback *Callback ) {
	buf := new(bytes.Buffer)
	appendObject(callback)
	writeInteger(buf, callback.Id())

	sendmsg(d, 1, buf.Bytes())
}

//// Events
func (d *Display) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type DisplayError struct {
	Object_id Object
	Code uint32
	Message string
}

func (d *Display) AddErrorListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func display_error(d *Display, msg []byte) {
	printEvent("error", msg)
	var data DisplayError
	buf := bytes.NewBuffer(msg)

	object_idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	object_id := new(Object)
	object_id = getObject(object_idid).(Object)
	data.Object_id = object_id

	code,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Code = code

	_,message,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Message = message

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
}

type DisplayGlobal struct {
	Name uint32
	Iface string
	Version uint32
}

func (d *Display) AddGlobalListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func display_global(d *Display, msg []byte) {
	printEvent("global", msg)
	var data DisplayGlobal
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	_,iface,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Iface = iface

	version,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Version = version

	for _,channel := range d.listeners[1] {
		go func () { channel <- data }()
	}
}

type DisplayGlobal_remove struct {
	Name uint32
}

func (d *Display) AddGlobal_removeListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func display_global_remove(d *Display, msg []byte) {
	printEvent("global_remove", msg)
	var data DisplayGlobal_remove
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
}

type DisplayDelete_id struct {
	Id uint32
}

func (d *Display) AddDelete_idListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
}

func display_delete_id(d *Display, msg []byte) {
	printEvent("delete_id", msg)
	var data DisplayDelete_id
	buf := bytes.NewBuffer(msg)

	id,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	for _,channel := range d.listeners[3] {
		go func () { channel <- data }()
	}
}

func NewDisplay() (d *Display) {
	d = new(Display)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, display_error)
	d.events = append(d.events, display_global)
	d.events = append(d.events, display_global_remove)
	d.events = append(d.events, display_delete_id)
	return
}

func (d *Display) SetId(id int32) {
	d.id = id
}

func (d *Display) Id() int32 {
	return d.id
}