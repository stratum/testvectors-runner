/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import "encoding/binary"

import "strconv"

func GetUint32(i uint32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, i)
	return byteSlice
}

func GetInt(s []byte) int {
	var res int
	for _, v := range s {
		res <<= 4
		res |= int(v)
	}
	return res
}

func GetStr(s []byte) string {
	return strconv.Itoa(GetInt(s))
}
