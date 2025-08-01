// Copyright © 2022 Hengqi Chen
package event

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

type inner struct {
	TimestampNS uint64   `json:"timestamp"`
	Pid         uint32   `json:"pid"`
	Tid         uint32   `json:"tid"`
	Len         int32    `json:"Len"`
	PayloadType uint8    `json:"payloadType"`
	Comm        [16]byte `json:"Comm"`
}

type GoTLSEvent struct {
	inner
	Data []byte `json:"data"`
}

func (ge *GoTLSEvent) Decode(payload []byte) error {
	r := bytes.NewBuffer(payload)
	err := binary.Read(r, binary.LittleEndian, &ge.inner)
	if err != nil {
		return err
	}
	if ge.Len > 0 {
		ge.Data = make([]byte, ge.Len)
		if err = binary.Read(r, binary.LittleEndian, &ge.Data); err != nil {
			return err
		}
	} else {
		ge.Len = 0
	}
	decodedKtime, err := DecodeKtime(int64(ge.TimestampNS), true)
	if err == nil {
		ge.TimestampNS = uint64(decodedKtime.Unix())
	}

	return err
}

func (ge *GoTLSEvent) String() string {
	s := fmt.Sprintf("PID: %d, Comm: %s, TID: %d, PayloadType:%d, Payload: %s\n", ge.Pid, string(ge.Comm[:]), ge.Tid, ge.inner.PayloadType, string(ge.Data[:ge.Len]))
	eb := new(Base)
	eb.UUID = ge.GetUUID()
	eb.PID = int32(ge.Pid)
	eb.PName = string(ge.Comm[:])
	eb.Type = uint32(ge.inner.PayloadType)
	eb.PayloadBase64 = base64.StdEncoding.EncodeToString(ge.Data[:])
	p, e := eb.Encode()
	if e != nil {
		return s
	}
	return string(p)
}

func (ge *GoTLSEvent) StringHex() string {
	perfix := COLORGREEN
	b := dumpByteSlice(ge.Data[:ge.Len], perfix)
	b.WriteString(COLORRESET)
	s := fmt.Sprintf("PID: %d, Comm: %s, TID: %d, PayloadType:%d, Payload: \n%s\n", ge.Pid, string(ge.Comm[:]), ge.Tid, ge.inner.PayloadType, b.String())
	return s
}

func (ge *GoTLSEvent) Clone() IEventStruct {
	return &GoTLSEvent{}
}

func (ge *GoTLSEvent) EventType() EventType {
	return EventTypeOutput
}

func (ge *GoTLSEvent) GetUUID() string {
	return fmt.Sprintf("%d_%d_%s", ge.Pid, ge.Tid, ge.Comm)
}

func (ge *GoTLSEvent) Payload() []byte {
	return ge.Data[:ge.Len]
}

func (ge *GoTLSEvent) PayloadLen() int {
	return int(ge.Len)
}
