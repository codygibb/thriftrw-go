// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ast

import "fmt"

// Type unifies the different types representing Thrift field types.
type Type interface {
	fmt.Stringer

	fieldType()
}

// BaseTypeID is an identifier for primitive types supported by Thrift.
type BaseTypeID int

//go:generate stringer -type=BaseTypeID

// IDs of the base types supported by Thrift.
const (
	BoolTypeID   BaseTypeID = iota + 1 // bool
	ByteTypeID                         // byte
	I16TypeID                          // i16
	I32TypeID                          // i32
	I64TypeID                          // i64
	DoubleTypeID                       // double
	StringTypeID                       // string
	BinaryTypeID                       // binary
)

// BaseType is a reference to a Thrift base type.
//
// 	bool, byte, i16, i32, i64, double, string, binary
//
// All references to base types in the document may be followed by type
// annotations.
//
// 	bool (go.type = "int")
type BaseType struct {
	// ID of the base type.
	ID BaseTypeID

	// Type annotations associated with this reference.
	Annotations []*Annotation
}

func (BaseType) fieldType() {}

func (bt BaseType) String() string {
	var name string

	switch bt.ID {
	case BoolTypeID:
		name = "bool"
	case ByteTypeID:
		name = "byte"
	case I16TypeID:
		name = "i16"
	case I32TypeID:
		name = "i32"
	case I64TypeID:
		name = "i64"
	case DoubleTypeID:
		name = "double"
	case StringTypeID:
		name = "string"
	case BinaryTypeID:
		name = "binary"
	default:
		panic(fmt.Sprintf("unknown base type %v", bt))
	}

	return appendAnnotations(name, bt.Annotations)
}

// MapType is a reference to a the Thrift map type.
//
// 	map<k, v>
//
// All references to map types may be followed by type annotations.
//
// 	map<string, list<i32>> (java.type = "MultiMap")
type MapType struct {
	KeyType, ValueType Type
	Annotations        []*Annotation
}

func (MapType) fieldType() {}

func (mt MapType) String() string {
	return appendAnnotations(
		fmt.Sprintf("map<%s, %s>", mt.KeyType, mt.ValueType),
		mt.Annotations,
	)
}

// ListType is a reference to the Thrift list type.
//
// 	list<a>
//
// All references to list types may be followed by type annotations.
//
// 	list<i64> (cpp.type = "vector")
type ListType struct {
	ValueType   Type
	Annotations []*Annotation
}

func (ListType) fieldType() {}

func (lt ListType) String() string {
	return appendAnnotations(
		fmt.Sprintf("list<%s>", lt.ValueType.String()),
		lt.Annotations,
	)
}

// SetType is a reference to the Thrift set type.
//
// 	set<a>
//
// All references to set types may be followed by type annotations.
//
// 	set<string> (js.type = "list")
type SetType struct {
	ValueType   Type
	Annotations []*Annotation
}

func (SetType) fieldType() {}

func (st SetType) String() string {
	return appendAnnotations(
		fmt.Sprintf("set<%s>", st.ValueType.String()),
		st.Annotations,
	)
}

// TypeReference references a user-defined type.
type TypeReference struct {
	Name string
	Line int
}

func (TypeReference) fieldType() {}

func (tr TypeReference) String() string {
	return tr.Name
}

func appendAnnotations(s string, anns []*Annotation) string {
	if a := FormatAnnotations(anns); len(a) > 0 {
		return s + " " + a
	}
	return s
}