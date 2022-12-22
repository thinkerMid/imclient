package json

import (
	"github.com/goccy/go-json"
)

//var sonicAPI sonic.API
//
//func init() {
//	sonicConfig := sonic.Config{
//		CopyString:           true,
//		NoQuoteTextMarshaler: true,
//	}
//
//	sonicAPI = sonicConfig.Froze()
//	sonic.ConfigDefault = sonicAPI
//}
//
//// Marshal returns the JSON encoding bytes of v.
//func Marshal(val interface{}) ([]byte, error) {
//	return sonicAPI.Marshal(val)
//}
//
//// MarshalString returns the JSON encoding string of v.
//func MarshalString(val interface{}) (string, error) {
//	return sonicAPI.MarshalToString(val)
//}
//
//// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
//// NOTICE: This API copies given buffer by default,
//// if you want to pass JSON more efficiently, use UnmarshalString instead.
//func Unmarshal(buf []byte, val interface{}) error {
//	return sonicAPI.Unmarshal(buf, val)
//}
//
//// UnmarshalString is like Unmarshal, except buf is a string.
//func UnmarshalString(buf string, val interface{}) error {
//	return sonicAPI.UnmarshalFromString(buf, val)
//}
//
//// MarshalIndent returns the JSON encoding bytes with indent and prefix.
//func MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
//	return sonicAPI.MarshalIndent(val, prefix, indent)
//}

// Marshal .
var Marshal = json.Marshal

// Unmarshal .
var Unmarshal = json.Unmarshal

// MarshalIndent .
var MarshalIndent = json.MarshalIndent
