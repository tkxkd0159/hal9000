package types

import (
	"strings"
)

type TypedEventParser struct {
	Events     TmEvents
	protoPkg   string
	protoMsg   string
	protoField string
}

func (p *TypedEventParser) SetEvents(evts TmEvents) {
	p.Events = evts
}

func (p *TypedEventParser) SetProtoPkg(protoPkgName string) {
	p.protoPkg = protoPkgName
}

func (p *TypedEventParser) SetProtoMsg(protoMsgName string) {
	p.protoMsg = protoMsgName
}

func (p *TypedEventParser) SetProtoField(protoFieldName string) {
	p.protoField = protoFieldName
}

func (p *TypedEventParser) Event() string {
	if p.protoField == "" {
		panic("You should add proto field name")
	}
	return p.EventWithFieldName(p.protoField)
}

func (p *TypedEventParser) EventWithFieldName(protoFieldName string) string {
	evtKey := strings.Join([]string{p.protoPkg, p.protoMsg, protoFieldName}, ".")
	if len(p.Events[evtKey]) > 0 {
		return p.Events[evtKey][0]
	}
	return ""
}

func NewTypedEventParser(pkgName, msgName string, fieldName ...string) *TypedEventParser {
	if len(fieldName) > 1 {
		panic("You should assign only one field")
	}
	if fieldName == nil {
		return &TypedEventParser{protoPkg: pkgName, protoMsg: msgName}
	}
	return &TypedEventParser{protoPkg: pkgName, protoMsg: msgName, protoField: fieldName[0]}
}
