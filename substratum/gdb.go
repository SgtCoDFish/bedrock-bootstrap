package substratum

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/cyrus-and/gdb"
)

// GdbConnection wraps a gdb connection. By wrapping the connection we can add additional features, like maintaining
// a per-connection map of registers to GDB-internal register numbers which must be used to fetch their values.
type GdbConnection struct {
	// Conn is the underlying connection to GDB which can be used manually to make requests
	Conn *gdb.Gdb

	// Logger is a simple log.Logger which is used to log informational messages
	Logger *log.Logger

	// registerToGDBNumber is a map of register names (both ABI like "zero" and number like "x0") to their underlying
	// GDB register numbers. It includes the program counter.
	registerToGDBNumber map[string]int

	// abiRegistersToGDBNumbers is similar to registerToGDBNumber except only ABI names are included
	abiRegistersToGDBNumbers map[string]int

	// gdbNumberToABIRegister maps GDB-internal register numbers to ABI names
	gdbNumberToABIRegister map[int]string

	// abiNames is a slice of ABI register names, in the same order as abiRegNumbers
	abiNames []string

	// abiRegNumbers is a slice of GDB-internal register numbers, in the same order as abiNames
	abiRegNumbers []int

	// allRegNumbers is the concatenation of every reg number in the same order as abiRegNumbers; it is used
	// to cache the string and avoid recreating it for every register dump
	allRegNumbers string
}

type gdbRegisterNamesResponse struct {
	Payload gdbRegisterNamesResponsePayload `json:"payload" mapstructure:"payload"`
}

type gdbRegisterNamesResponsePayload struct {
	RegisterNames []string `json:"register-names" mapstructure:"register-names"`
}

// NewGdbConnection creates a GdbConnection with given parameters, and initialises that connection
// to use the RISC-V rv32 architecture.
func NewGdbConnection(logger *log.Logger, gdbPath string, remoteTarget string) (*GdbConnection, error) {
	conn, err := gdb.NewCmd([]string{gdbPath, "--nx", "--quiet", "--interpreter=mi2", "-ex", "set architecture riscv:rv32"}, nil)
	if err != nil {
		return nil, err
	}

	_, err = conn.CheckedSend(fmt.Sprintf("target-select remote %s", remoteTarget))
	if err != nil {
		return nil, err
	}

	resp, err := conn.CheckedSend("data-list-register-names")
	if err != nil {
		return nil, err
	}

	var registerNamesResponse gdbRegisterNamesResponse
	err = mapstructure.Decode(resp, &registerNamesResponse)
	if err != nil {
		return nil, err
	}

	regList := append(GetRegisterList(), "pc")

	gdbConn := GdbConnection{
		Logger: logger,
		Conn:   conn,

		registerToGDBNumber:      make(map[string]int),
		abiRegistersToGDBNumbers: make(map[string]int),
		gdbNumberToABIRegister:   make(map[int]string),

		abiNames:      make([]string, len(regList)),
		abiRegNumbers: make([]int, len(regList)),
	}

	for regIdx, regName := range regList {
		for i, foundReg := range registerNamesResponse.Payload.RegisterNames {
			if regName == foundReg {
				gdbConn.registerToGDBNumber[regName] = i
				gdbConn.abiRegistersToGDBNumbers[regName] = i

				gdbConn.abiNames[regIdx] = regName
				gdbConn.abiRegNumbers[regIdx] = i
				gdbConn.gdbNumberToABIRegister[i] = regName

				if regName == "pc" {
					// there's no "number register" for the program counter
					break
				}

				numberRegName, err := GetNumberRegisterForABIName(regName)

				if err != nil {
					logger.Printf("[warning] couldn't find number register for %s", regName)
				} else {
					gdbConn.registerToGDBNumber[numberRegName] = i
				}

				break
			}
		}
	}

	numberStringBuilder := new(strings.Builder)
	for _, i := range gdbConn.abiRegNumbers {
		numberStringBuilder.WriteString(strconv.Itoa(i))
		numberStringBuilder.WriteString(" ")
	}

	gdbConn.allRegNumbers = numberStringBuilder.String()

	return &gdbConn, nil
}

type gdbRegisterValuesResponse struct {
	Payload gdbRegisterValuesResponsePayload `json:"payload" mapstructure:"payload"`
}

type gdbRegisterValue struct {
	RegNumber string `json:"number" mapstructure:"number"`
	Value     string `json:"value" mapstructure:"value"`
}

type gdbRegisterValuesResponsePayload struct {
	RegisterValues []gdbRegisterValue `json:"register-values" mapstructure:"register-values"`
}

// GDBRegisterFrame holds the typed values of every regular RISC-V rv32 integer register along with the PC
type GDBRegisterFrame struct {
	Zero uint32 `json:"zero" mapstructure:"zero"`
	Ra   uint32 `json:"ra"  mapstructure:"ra"`
	Sp   uint32 `json:"sp"  mapstructure:"sp"`
	Gp   uint32 `json:"gp"  mapstructure:"gp"`
	Tp   uint32 `json:"tp"  mapstructure:"tp"`
	T0   uint32 `json:"t0"  mapstructure:"t0"`
	T1   uint32 `json:"t1"  mapstructure:"t1"`
	T2   uint32 `json:"t2"  mapstructure:"t2"`
	Fp   uint32 `json:"fp"  mapstructure:"fp"`
	S1   uint32 `json:"s1"  mapstructure:"s1"`
	A0   uint32 `json:"a0"  mapstructure:"a0"`
	A1   uint32 `json:"a1"  mapstructure:"a1"`
	A2   uint32 `json:"a2"  mapstructure:"a2"`
	A3   uint32 `json:"a3"  mapstructure:"a3"`
	A4   uint32 `json:"a4"  mapstructure:"a4"`
	A5   uint32 `json:"a5"  mapstructure:"a5"`
	A6   uint32 `json:"a6"  mapstructure:"a6"`
	A7   uint32 `json:"a7"  mapstructure:"a7"`
	S2   uint32 `json:"s2"  mapstructure:"s2"`
	S3   uint32 `json:"s3"  mapstructure:"s3"`
	S4   uint32 `json:"s4"  mapstructure:"s4"`
	S5   uint32 `json:"s5"  mapstructure:"s5"`
	S6   uint32 `json:"s6"  mapstructure:"s6"`
	S7   uint32 `json:"s7"  mapstructure:"s7"`
	S8   uint32 `json:"s8"  mapstructure:"s8"`
	S9   uint32 `json:"s9"  mapstructure:"s9"`
	S10  uint32 `json:"s10" mapstructure:"s10"`
	S11  uint32 `json:"s11" mapstructure:"s11"`
	T3   uint32 `json:"t3"  mapstructure:"t3"`
	T4   uint32 `json:"t4"  mapstructure:"t4"`
	T5   uint32 `json:"t5"  mapstructure:"t5"`
	T6   uint32 `json:"t6"  mapstructure:"t6"`
	PC   uint32 `json:"pc"  mapstructure:"pc"`
}

// AsMap returns the given register frame as a map of register names to register values
// In some cases - such as iterating over all registers - a map is much easier to use
func (f GDBRegisterFrame) AsMap() map[string]uint32 {
	asMap := make(map[string]uint32)
	jsonEnc, _ := json.Marshal(f)

	_ = json.Unmarshal(jsonEnc, &asMap)
	return asMap
}

// Dump uses the given logger to write prettified contents of every register except `zero`
func (f GDBRegisterFrame) Dump(logger *log.Logger) {
	asMap := f.AsMap()
	regList := append(GetRegisterList(), "pc")

	for _, regName := range regList {
		if regName == "zero" {
			continue
		}

		value := asMap[regName]
		if value == 0 {
			logger.Printf("%4.4s:          0", regName)
		} else {
			logger.Printf("%4.4s: 0x%8.8x", regName, asMap[regName])
		}
	}
}

// FetchRegister makes a GDB call equivalent to `i r n` or `info registers n` to get the value of "n", which can be
// any RISC-V rv32i integer register including the PC.
func (s *GdbConnection) FetchRegister(name string) (uint32, error) {
	gdbRegNumber, ok := s.registerToGDBNumber[strings.ToLower(name)]
	if !ok {
		return 0, fmt.Errorf("couldn't find internal GDB register number for '%s'", name)
	}

	resp, err := s.Conn.CheckedSend(fmt.Sprintf("data-list-register-values x %d", gdbRegNumber))
	if err != nil {
		return 0, err
	}

	var valuesResponse gdbRegisterValuesResponse
	err = mapstructure.Decode(resp, &valuesResponse)
	if err != nil {
		return 0, err
	}

	for _, registerValue := range valuesResponse.Payload.RegisterValues {
		fetchedRegNumber, err := strconv.Atoi(registerValue.RegNumber)
		if err != nil {
			return 0, fmt.Errorf("invalid internal register number: '%s'", registerValue.RegNumber)
		}

		if gdbRegNumber == fetchedRegNumber {
			value, err := strconv.ParseUint(registerValue.Value, 0, 32)
			if err != nil {
				return 0, err
			}

			return uint32(value), nil
		}
	}

	return 0, fmt.Errorf("couldn't find register value in response for %s", name)
}

// FetchPC is a convenience function to return the current value of the PC.
func (s *GdbConnection) FetchPC() (uint32, error) {
	return s.FetchRegister("pc")
}

// FetchRegisterFrame returns, if possible, a snapshot of the state of all current registers as a GDBRegisterFrame
func (s *GdbConnection) FetchRegisterFrame() (GDBRegisterFrame, error) {
	var registerFrame GDBRegisterFrame

	resp, err := s.Conn.CheckedSend(fmt.Sprintf("data-list-register-values x %s", s.allRegNumbers))
	if err != nil {
		return registerFrame, err
	}

	var valuesResponse gdbRegisterValuesResponse
	err = mapstructure.Decode(resp, &valuesResponse)
	if err != nil {
		return registerFrame, err
	}

	outMap := make(map[string]uint32)

	for _, registerValue := range valuesResponse.Payload.RegisterValues {
		gdbRegNumber, err := strconv.Atoi(registerValue.RegNumber)
		if err != nil {
			return registerFrame, fmt.Errorf("invalid internal register number: '%s'", registerValue.RegNumber)
		}

		regName, ok := s.gdbNumberToABIRegister[gdbRegNumber]
		if !ok {
			return registerFrame, fmt.Errorf("corresponding register not found for register number: %s", registerValue.RegNumber)
		}

		value, err := strconv.ParseUint(registerValue.Value, 0, 32)
		if err != nil {
			return registerFrame, err
		}

		outMap[regName] = uint32(value)
	}

	err = mapstructure.Decode(outMap, &registerFrame)
	if err != nil {
		return registerFrame, err
	}

	return registerFrame, nil
}
