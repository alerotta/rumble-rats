package protocol

// this files only defines the envelope common to all the messages
// and the types of messages that are transmitted and their structures.

import "encoding/json"

type Envelope struct {
	Type string
	Seq uint32
	Ts int64
	data json.RawMessage
} 

type ClientHello struct{
	Version string
	Name string
}

type ClientInput struct{

	// to be completed
}

type ServerWelcome struct {
	PlayerID string
	TickRate int
}

type ServerError struct {
	Code  string
	Message string 
}

type SnapshotMessage struct {
	Tich uint64
	State any //to be replaced 
}