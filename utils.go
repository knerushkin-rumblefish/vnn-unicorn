package main

import (
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

func WriteRam(ram map[uint32](uint32), addr uint32, value uint32) {
	// we no longer delete from ram, since deleting from tries is hard
	if value == 0 && false {
		delete(ram, addr)
	} else {
		/*if addr < 0xc0000000 {
			fmt.Printf("store %x = %x\n", addr, value)
		}*/
		ram[addr] = value
	}
}

var heap_start uint64 = 0
var REG_OFFSET uint32 = 0xc0000000
var REG_PC uint32 = REG_OFFSET + 0x20*4
var REG_HEAP uint32 = REG_OFFSET + 0x23*4

func SyncRegs(mu uc.Unicorn, ram map[uint32](uint32)) {
	pc, _ := mu.RegRead(uc.MIPS_REG_PC)
	//fmt.Printf("%d uni %x\n", step, pc)
	WriteRam(ram, 0xc0000080, uint32(pc))

	addr := uint32(0xc0000000)
	for i := uc.MIPS_REG_ZERO; i < uc.MIPS_REG_ZERO+32; i++ {
		reg, _ := mu.RegRead(i)
		WriteRam(ram, addr, uint32(reg))
		addr += 4
	}

	reg_hi, _ := mu.RegRead(uc.MIPS_REG_HI)
	reg_lo, _ := mu.RegRead(uc.MIPS_REG_LO)
	WriteRam(ram, REG_OFFSET+0x21*4, uint32(reg_hi))
	WriteRam(ram, REG_OFFSET+0x22*4, uint32(reg_lo))

	WriteRam(ram, REG_HEAP, uint32(heap_start))
}

func ZeroRegisters(ram map[uint32](uint32)) {
	for i := uint32(0xC0000000); i < 0xC0000000+36*4; i += 4 {
		WriteRam(ram, i, 0)
	}
}
