package main

import (
	"encoding/binary"
	"log"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

// GetTemp - Measure air temperature
func GetTemp() float32 {
	_, err := host.Init()
	if err != nil {
		log.Println("host.Init()")
		log.Fatal(err)
	}

	p, err := i2creg.Open("")
	if err != nil {
		log.Println("i2creg.Open(\"\")")
		log.Fatal(err)
	}
	defer p.Close()

	d := &i2c.Dev{Addr: 0x40, Bus: p}

	write := []byte{0xE3}
	read := make([]byte, 2)
	if err := d.Tx(write, read); err != nil {
		log.Println("d.Tx(write, read)")
		log.Fatal(err)
	}

	return CalculateCelsius(read)
}

func BytesToFloat(bytes []byte) float32 {
	float := float32(binary.BigEndian.Uint16(bytes))
	return float
}

func CalculateCelsius(bytes []byte) float32 {
	celsius := -46.85 + (175.72 * (BytesToFloat(bytes) / 65536))
	return celsius
}
