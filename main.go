package main

import (
	"bytes"
	"encoding/binary"
	"io"

	// "encoding/hex"
	"bufio"
	"errors"
	"fmt"
	"os"

	// "github.com/deadsy/rvda"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

func LoadFile(mu uc.Unicorn, fn string, base uint64) uint64 {
	dat, err := os.ReadFile(fn)
	if err != nil {
		fmt.Println("Not loaded program", err)
	}
	mu.MemWrite(base, dat)
	return uint64(len(dat))
}

func LoadBigFile(mu uc.Unicorn, fn string, base uint64) {
	file, err := os.Open(fn)
	if err != nil {
		panic("Bad big file")
	}
	defer file.Close()

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(file)
	buf := make([]byte, 0, 4)
	address_base := base
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
		}

		address_base = address_base + (4)
		mu.MemWrite(address_base, buf)
		nBytes += int64(len(buf))
		nChunks++
	}
	fmt.Println("bytes:", nBytes)
	fmt.Println("chunks:", nChunks)
}

func FileLength(fn string) uint64 {
	dat, err := os.ReadFile(fn)
	if err != nil {
		fmt.Println("Not loaded program", err)
	}

	return uint64(len(dat))
}
func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return x
}

func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return x
}
func Int64ToBytes(n uint64) []byte {
	x := uint64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, binary.LittleEndian, x)

	return bytesBuffer.Bytes()
}
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)

	return bytesBuffer.Bytes()
}

func pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errors.New("invalid blocksize")
	}
	if b == nil || len(b) == 0 {
		return nil, errors.New("invalid blocksize")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(0)}, n))
	return pb, nil
}
func main() {
	create_binary_blob()

	INPUT_PATH := "./blob.bin"
	// PROGRAM_PATH := "./program/startup/infinite.bin"
	PROGRAM_PATH := "./riscv-from-scratch/risc64-in-c.bin"
	// PROGRAM_PATH := "./program/program.bin"

	input_len := FileLength(INPUT_PATH)
	fmt.Println("input length: ", input_len)
	fmt.Println("input bytes length", len(Int64ToBytes(input_len)))

	program_len := FileLength(PROGRAM_PATH)
	fmt.Println("program length: ", program_len)

	ADDRESS_START := uint64(0x8000000)
	fmt.Printf("Start Address: 0x%x\n", ADDRESS_START)

	// PROGRAM_SLOT := uint64(0x10000)
	// PROGRAM_ADDRESS := uint64(0x1000000)

	OUTPUT_SLOT := uint64(0x100000)
	// OUTPUT_ADDRESS := INPUT_ADDRESS + INPUT_ADDRESS_SLOT + INPUT_SLOT
	OUTPUT_ADDRESS := uint64(0x40000000)

	INPUT_SLOT := uint64(0x40000000)
	// INPUT_ADDRESS_SLOT := uint64(8)
	// // INPUT_ADDRESS := ADDRESS_START + PROGRAM_SLOT
	// INPUT_ADDRESS := uint64(0x40000000) + PROGRAM_SLOT
	INPUT_LENGTH_ADDRESS := OUTPUT_ADDRESS + OUTPUT_SLOT
	INPUT_DATA_ADDRESS := INPUT_LENGTH_ADDRESS + 8

	// ADDRESS_END := OUTPUT_ADDRESS + OUTPUT_SLOT + uint64(0x1000000)
	ADDRESS_END := INPUT_DATA_ADDRESS + INPUT_SLOT
	fmt.Printf("End Address: 0x%x\n", ADDRESS_END)

	fmt.Printf("Input Address: 0x%x\n", INPUT_DATA_ADDRESS)
	fmt.Printf("Output Address: 0x%x\n", OUTPUT_ADDRESS)
	// isa, _ := rvda.New(32, rvda.RV64gc)
	mu, _ := uc.NewUnicorn(uc.ARCH_RISCV, uc.MODE_RISCV64)

	// mu.HookAdd(uc.HOOK_CODE, func(mu uc.Unicorn, addr uint64, size uint32) {
	// 	fmt.Printf("code: addr 0x%x size %d\n", addr, size)
	//
	// 	ins_bytes, _ := mu.MemRead(addr, uint64(size))
	// 	normalized_ins_bytes, _ := pad(ins_bytes, 4)
	// 	normalized_ins := BytesToInt32(normalized_ins_bytes)
	// 	fmt.Printf("ins: 0x%v\n", hex.EncodeToString(normalized_ins_bytes))
	// 	normalized_da := isa.Disassemble(uint(addr), uint(normalized_ins))
	// 	fmt.Printf("normalized decode: %#v\n", normalized_da)
	//
	// 	reg_a4_raw, _ := mu.RegRead(uc.RISCV_REG_A4)
	// 	reg_a5_raw, _ := mu.RegRead(uc.RISCV_REG_A5)
	//
	// 	fmt.Println("reg A4: ", reg_a4_raw)
	// 	fmt.Println("reg A5: ", reg_a5_raw)
	// }, ADDRESS_START, ADDRESS_END)

	// mu.HookAdd(uc.HOOK_INTR, func(mu uc.Unicorn, intno uint32) {
	//
	// 	fmt.Println("intr no:", intno)
	// 	x1, _ := mu.RegRead(uc.RISCV_REG_X1)
	// 	x2, _ := mu.RegRead(uc.RISCV_REG_X2)
	// 	x3, _ := mu.RegRead(uc.RISCV_REG_X3)
	// 	x4, _ := mu.RegRead(uc.RISCV_REG_X4)
	// 	x5, _ := mu.RegRead(uc.RISCV_REG_X5)
	// 	x6, _ := mu.RegRead(uc.RISCV_REG_X6)
	// 	x7, _ := mu.RegRead(uc.RISCV_REG_X7)
	// 	x8, _ := mu.RegRead(uc.RISCV_REG_X8)
	// 	x9, _ := mu.RegRead(uc.RISCV_REG_X9)
	//
	// 	fmt.Println("x1: ", x1)
	// 	fmt.Println("x2: ", x2)
	// 	fmt.Println("x3: ", x3)
	// 	fmt.Println("x4: ", x4)
	// 	fmt.Println("x5: ", x5)
	// 	fmt.Println("x6: ", x6)
	// 	fmt.Println("x7: ", x7)
	// 	fmt.Println("x8: ", x8)
	// 	fmt.Println("x9: ", x9)
	// }, 0, OUTPUT_ADDRESS+OUTPUT_ADDRESS_OFFSET)
	// fmt.Println(mu)

	// mu.HookAdd(uc.HOOK_MEM_READ, func(mu uc.Unicorn, access int, addr uint64, size int, value int64) {
	// 	if value != 0 {
	// 		// 	// addr > INPUT_ADDRESS &&
	// 		// 	// 	addr < INPUT_ADDRESS+INPUT_ADDRESS_SLOT
	// 		// 	// {
	// 		//
	// 		fmt.Printf("0x%x 0x%x  0x%x\n", INPUT_ADDRESS, addr, INPUT_ADDRESS+INPUT_ADDRESS_SLOT)
	// 		fmt.Printf("mem read: @0x%x, 0x%x = 0x%x\n", addr, size, value)
	// 	}
	//
	// }, 0, ADDRESS_END)
	// mu.HookAdd(uc.HOOK_MEM_WRITE, func(mu uc.Unicorn, access int, addr uint64, size int, value int64) {
	// 	fmt.Printf("mem write: @0x%x, 0x%x = 0x%x\n", addr, size, value)
	// }, 0, ADDRESS_END)

	mu.MemMap(0, uint64(32*1024*1024*1024))

	mu.MemWrite(OUTPUT_ADDRESS, Int64ToBytes(0))

	mu.MemWrite(INPUT_LENGTH_ADDRESS, Int64ToBytes(input_len))
	LoadBigFile(mu, INPUT_PATH, INPUT_DATA_ADDRESS)
	// LoadFile(mu, INPUT_PATH, INPUT_ADDRESS+4)

	LoadFile(mu, PROGRAM_PATH, ADDRESS_START)

	if err := mu.Start(ADDRESS_START, 0); err != nil {
		panic(err)
	}

	input_value_bytes, _ := mu.MemRead(INPUT_LENGTH_ADDRESS, 8)
	input_value := BytesToInt64(input_value_bytes)
	fmt.Println("Input value: ", input_value)

	output_value_bytes, _ := mu.MemRead(OUTPUT_ADDRESS, 8)
	output_value := BytesToInt64(output_value_bytes)
	fmt.Println("Output value: ", output_value)

}
