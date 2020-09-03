// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The code in this file was largely written by Damian Gryski as part of
// https://github.com/dgryski/go-tsz and published under the license below.
// It received minor modifications to suit Prometheus's needs.

// Copyright (c) 2015,2016 Damian Gryski <damian@gryski.com>
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// * Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package chunkenc

import "io"

// bstream is a stream of bits.
type bstream struct {
	stream []byte // the data stream
	count  uint8  // how many bits are valid in current byte, 当前字节有多少bit可用，0-8
}

func newBReader(b []byte) bstream {
	return bstream{stream: b, count: 8}
}

//bstream getter
func (b *bstream) bytes() []byte {
	return b.stream
}

type bit bool

const (
	zero bit = false
	one  bit = true
)

func (b *bstream) writeBit(bit bit) { //写入一个bit
	if b.count == 0 {
		b.stream = append(b.stream, 0) //把stream后新的一个完整字节全部变为0, 因为切片种类为byte， 所以0只有一个字节大小
		b.count = 8
	}

	i := len(b.stream) - 1 //最后一个字节的位置

	if bit { //如果bit为0则不需要更改
		b.stream[i] |= 1 << (b.count - 1) //按位或后赋值， 放入1到具体位置
	}

	b.count--
}

func (b *bstream) writeByte(byt byte) { //写入一个byte
	if b.count == 0 {
		b.stream = append(b.stream, 0) //把stream后新的一个完整字节全部变为0
		b.count = 8
	}

	i := len(b.stream) - 1 //最后一个字节的位置

	// fill up b.b with b.count bits from byt
	b.stream[i] |= byt >> (8 - b.count) //剩下b.count位放入剩余空的位置中

	b.stream = append(b.stream, 0) //把stream后新的一个完整字节全部变为0
	i++                            //调整最后的位置
	b.stream[i] = byt << b.count   //加入的字节前几位为剩余的data，后面的则为空（0）
}

func (b *bstream) writeBits(u uint64, nbits int) { //写入多个bits
	u <<= (64 - uint(nbits)) //	左移后赋值，把打他放到最左侧
	for nbits >= 8 {         //需一个一个字节写入
		byt := byte(u >> 56) //需写入的字节
		b.writeByte(byt)
		u <<= 8 //移除写入的一个字节
		nbits -= 8
	}

	for nbits > 0 {
		b.writeBit((u >> 63) == 1) //写入一个bit
		u <<= 1                    //移除写入的bit
		nbits--
	}
}

//从头开始读，count记录第一字节剩余没有被阅读的bits，bit从byte最左边开始读取
func (b *bstream) readBit() (bit, error) {
	if len(b.stream) == 0 {
		return false, io.EOF
	}

	if b.count == 0 {
		b.stream = b.stream[1:]

		if len(b.stream) == 0 {
			return false, io.EOF
		}
		b.count = 8
	}

	d := (b.stream[0] << (8 - b.count)) & 0x80 // 先移位，与1000 0000(uint)比较， 每一位比较均为1对应结果位数为1， 即只看第一位
	b.count--
	return d != 0, nil
}

func (b *bstream) ReadByte() (byte, error) {
	return b.readByte()
}

func (b *bstream) readByte() (byte, error) {
	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	if b.count == 0 {
		b.stream = b.stream[1:]

		if len(b.stream) == 0 {
			return 0, io.EOF
		}
		return b.stream[0], nil
	}

	if b.count == 8 {
		b.count = 0
		return b.stream[0], nil
	}

	byt := b.stream[0] << (8 - b.count)
	b.stream = b.stream[1:]

	if len(b.stream) == 0 {
		return 0, io.EOF
	}

	// We just advanced the stream and can assume the shift to be 0.
	byt |= b.stream[0] >> b.count

	return byt, nil
}

func (b *bstream) readBits(nbits int) (uint64, error) {
	var u uint64

	for nbits >= 8 {
		byt, err := b.readByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		nbits -= 8
	}

	if nbits == 0 {
		return u, nil
	}

	if nbits > int(b.count) {
		u = (u << uint(b.count)) | uint64((b.stream[0]<<(8-b.count))>>(8-b.count)) //只选出b.count的位数并放到最后
		nbits -= int(b.count)
		b.stream = b.stream[1:]

		if len(b.stream) == 0 {
			return 0, io.EOF
		}
		b.count = 8
	}

	u = (u << uint(nbits)) | uint64((b.stream[0]<<(8-b.count))>>(8-uint(nbits)))
	b.count -= uint8(nbits)
	return u, nil
}
