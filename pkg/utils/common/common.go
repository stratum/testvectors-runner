/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import "encoding/binary"

func GetUint32(i uint32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, i)
	return byteSlice
}
