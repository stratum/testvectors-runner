/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package common implements generic functions used in the repo
*/
package common

import (
	"encoding/binary"
	"strconv"
)

//GetUint32 converts uint32 to byte slice
func GetUint32(i uint32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, i)
	return byteSlice
}

//GetUint16 converts uint32 to byte slice
func GetUint16(i uint32) []byte {
	i16 := uint16(i)
	byteSlice := make([]byte, 2)
	binary.BigEndian.PutUint16(byteSlice, i16)
	return byteSlice
}

//GetInt converts byte slice to int
func GetInt(s []byte) int {
	var res int
	for _, v := range s {
		res <<= 4
		res |= int(v)
	}
	return res
}

//GetStr converts various interfaces to string
func GetStr(i interface{}) string {
	switch v := i.(type) {
	case uint16:
		return strconv.Itoa(int(v))
	case uint32:
		return strconv.Itoa(int(v))
	case []byte:
		return strconv.Itoa(GetInt(v))
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}
