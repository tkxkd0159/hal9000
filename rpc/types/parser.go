package types

import (
	"strings"
)

type TypedEventParser struct {
	protoPkg   string
	protoMsg   string
	protoField string
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
	return strings.Join([]string{p.protoPkg, p.protoMsg, protoFieldName}, ".")
}

func NewTypedEventParser(pkgName, msgName string, fieldName ...string) *TypedEventParser {
	if len(fieldName) != 1 {
		panic("You should assign only onw field")
	}
	if fieldName == nil {
		return &TypedEventParser{protoPkg: pkgName, protoMsg: msgName}
	}
	return &TypedEventParser{protoPkg: pkgName, protoMsg: msgName, protoField: fieldName[0]}
}
