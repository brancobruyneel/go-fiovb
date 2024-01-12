package teec

/*
#include "tee_client_api.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void fillTempMemoryReference(TEEC_Parameter *parameter, void *buffer, size_t size)
{
	parameter->tmpref.buffer = buffer;
	parameter->tmpref.size = size;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type ParameterTypes uint32

const (
	NONE                  ParameterTypes = 0x00000000
	VALUE_INPUT           ParameterTypes = 0x00000001
	VALUE_OUTPUT          ParameterTypes = 0x00000002
	VALUE_INOUT           ParameterTypes = 0x00000003
	MEMREF_TEMP_INPUT     ParameterTypes = 0x00000005
	MEMREF_TEMP_OUTPUT    ParameterTypes = 0x00000006
	MEMREF_TEMP_INOUT     ParameterTypes = 0x00000007
	MEMREF_WHOLE          ParameterTypes = 0x0000000C
	MEMREF_PARTIAL_INPUT  ParameterTypes = 0x0000000D
	MEMREF_PARTIAL_OUTPUT ParameterTypes = 0x0000000E
	MEMREF_PARTIAL_INOUT  ParameterTypes = 0x0000000F
)

type TEEC struct {
	context *C.TEEC_Context
	session *C.TEEC_Session
}

type UUID struct {
	TimeLow          uint32
	TimeMid          uint16
	TimeHiAndVersion uint16
	ClockSeqAndNode  [8]byte
}

type Parameter struct {
	Buffer []byte
	Size   uint32
}

type Operation struct {
	ParamTypes [4]ParameterTypes
	Params     [4]Parameter
}

func New() *TEEC {
	return &TEEC{}
}

func (teec *TEEC) Initialize() error {
	teec.context = (*C.TEEC_Context)(C.calloc(1, C.sizeof_TEEC_Context))

	res := C.TEEC_InitializeContext(nil, teec.context)

	if res != C.TEEC_SUCCESS {
		teec.freeContext()
		return teec.resultToErrorMessage("Initialize", res)
	}

	return nil
}

func (teec *TEEC) OpenSession(destination UUID) error {
	teec.session = (*C.TEEC_Session)(C.calloc(1, C.sizeof_TEEC_Session))

	dest := (*C.TEEC_UUID)(C.calloc(1, C.sizeof_TEEC_UUID))
	defer C.free(unsafe.Pointer(dest))

	dest.timeLow = C.uint32_t(destination.TimeLow)
	dest.timeMid = C.uint16_t(destination.TimeMid)
	dest.timeHiAndVersion = C.uint16_t(destination.TimeHiAndVersion)
	C.memcpy(unsafe.Pointer(&dest.clockSeqAndNode), unsafe.Pointer((*C.uint8_t)(&destination.ClockSeqAndNode[0])), 8)

	res := C.TEEC_OpenSession(teec.context, teec.session, dest, 0, nil, nil, nil)

	if res != C.TEEC_SUCCESS {
		teec.freeSession()
		return teec.resultToErrorMessage("OpenSession", res)
	}

	return nil
}

func (teec *TEEC) InvokeCommand(commandID uint32, operation *Operation, origin *uint32) error {
	if teec.session == nil {
		return fmt.Errorf("session is nil")
	}

	if operation == nil {
		return fmt.Errorf("operation is nil")
	}

	var op C.TEEC_Operation

	op.paramTypes = teec.parameterTypes(operation)

	var buffers [4]*C.char

	for i, parameter := range operation.Params {
		buffers[i] = C.CString(string(parameter.Buffer))
	}

	defer func() {
		for _, buffer := range buffers {
			C.free(unsafe.Pointer(buffer))
		}
	}()

	for i, parameter := range operation.Params {
		switch operation.ParamTypes[i] {
		case NONE:
			// Do nothing
			break
		case MEMREF_TEMP_INPUT:
			C.fillTempMemoryReference(&op.params[i], unsafe.Pointer(buffers[i]), C.size_t(parameter.Size))
		case MEMREF_TEMP_OUTPUT:
			C.fillTempMemoryReference(&op.params[i], unsafe.Pointer(buffers[i]), C.size_t(parameter.Size))
		case MEMREF_TEMP_INOUT:
			C.fillTempMemoryReference(&op.params[i], unsafe.Pointer(buffers[i]), C.size_t(parameter.Size))
		case VALUE_INPUT:
			return fmt.Errorf("VALUE_INPUT not implemented")
		case VALUE_OUTPUT:
			return fmt.Errorf("VALUE_OUTPUT not implemented")
		case VALUE_INOUT:
			return fmt.Errorf("VALUE_INOUT not implemented")
		case MEMREF_WHOLE:
			return fmt.Errorf("MEMREF_WHOLE not implemented")
		case MEMREF_PARTIAL_INPUT:
			return fmt.Errorf("MEMREF_PARTIAL_INPUT not implemented")
		case MEMREF_PARTIAL_OUTPUT:
			return fmt.Errorf("MEMREF_PARTIAL_OUTPUT not implemented")
		case MEMREF_PARTIAL_INOUT:
			return fmt.Errorf("MEMREF_PARTIAL_INOUT not implemented")
		}
	}

	res := C.TEEC_InvokeCommand(teec.session, (C.uint32_t)(commandID), &op, (*C.uint32_t)(origin))

	if res != C.TEEC_SUCCESS {
		return teec.resultToErrorMessage("InvokeCommand", res)
	}

	for i, buffer := range buffers {
		switch operation.ParamTypes[i] {
		case NONE:
			// Do nothing
		case MEMREF_TEMP_INPUT:
			// Do nothing
		case MEMREF_TEMP_OUTPUT:
			operation.Params[i].Buffer = []byte(C.GoString(buffer))
		case MEMREF_TEMP_INOUT:
			operation.Params[i].Buffer = []byte(C.GoString(buffer))
		case VALUE_INPUT:
			return fmt.Errorf("VALUE_INPUT not implemented")
		case VALUE_OUTPUT:
			return fmt.Errorf("VALUE_OUTPUT not implemented")
		case VALUE_INOUT:
			return fmt.Errorf("VALUE_INOUT not implemented")
		case MEMREF_WHOLE:
			return fmt.Errorf("MEMREF_WHOLE not implemented")
		case MEMREF_PARTIAL_INPUT:
			return fmt.Errorf("MEMREF_PARTIAL_INPUT not implemented")
		case MEMREF_PARTIAL_OUTPUT:
			return fmt.Errorf("MEMREF_PARTIAL_OUTPUT not implemented")
		case MEMREF_PARTIAL_INOUT:
			return fmt.Errorf("MEMREF_PARTIAL_INOUT not implemented")
		}
	}

	return nil
}

func (teec *TEEC) CloseSession() error {
	if teec.session == nil {
		return fmt.Errorf("session is nil")
	}

	C.TEEC_CloseSession(teec.session)
	teec.freeSession()

	return nil
}

func (teec *TEEC) Finalize() error {
	if teec.context == nil {
		return fmt.Errorf("context is nil")
	}

	C.TEEC_FinalizeContext(teec.context)
	teec.freeContext()

	return nil
}

func (teec *TEEC) parameterTypes(operation *Operation) C.uint32_t {
	var paramTypes C.uint32_t

	paramTypes = 0

	for i, param := range operation.ParamTypes {
		paramTypes |= (C.uint32_t)(param) << (i * 4)
	}

	return paramTypes
}

func (teec *TEEC) resultToErrorMessage(name string, res C.TEEC_Result) error {

	switch res {
	case C.TEEC_ERROR_STORAGE_NOT_AVAILABLE:
		return fmt.Errorf("%s: TEEC_ERROR_STORAGE_NOT_AVAILABLE", name)
	case C.TEEC_ERROR_GENERIC:
		return fmt.Errorf("%s: TEEC_ERROR_GENERIC", name)
	case C.TEEC_ERROR_ACCESS_DENIED:
		return fmt.Errorf("%s: TEEC_ERROR_ACCESS_DENIED", name)
	case C.TEEC_ERROR_CANCEL:
		return fmt.Errorf("%s: TEEC_ERROR_CANCEL", name)
	case C.TEEC_ERROR_ACCESS_CONFLICT:
		return fmt.Errorf("%s: TEEC_ERROR_ACCESS_CONFLICT", name)
	case C.TEEC_ERROR_EXCESS_DATA:
		return fmt.Errorf("%s: TEEC_ERROR_EXCESS_DATA", name)
	case C.TEEC_ERROR_BAD_FORMAT:
		return fmt.Errorf("%s: TEEC_ERROR_BAD_FORMAT", name)
	case C.TEEC_ERROR_BAD_PARAMETERS:
		return fmt.Errorf("%s: TEEC_ERROR_BAD_PARAMETERS", name)
	case C.TEEC_ERROR_BAD_STATE:
		return fmt.Errorf("%s: TEEC_ERROR_BAD_STATE", name)
	case C.TEEC_ERROR_ITEM_NOT_FOUND:
		return fmt.Errorf("%s: TEEC_ERROR_ITEM_NOT_FOUND", name)
	case C.TEEC_ERROR_NOT_IMPLEMENTED:
		return fmt.Errorf("%s: TEEC_ERROR_NOT_IMPLEMENTED", name)
	case C.TEEC_ERROR_NOT_SUPPORTED:
		return fmt.Errorf("%s: TEEC_ERROR_NOT_SUPPORTED", name)
	case C.TEEC_ERROR_NO_DATA:
		return fmt.Errorf("%s: TEEC_ERROR_NO_DATA", name)
	case C.TEEC_ERROR_OUT_OF_MEMORY:
		return fmt.Errorf("%s: TEEC_ERROR_OUT_OF_MEMORY", name)
	case C.TEEC_ERROR_BUSY:
		return fmt.Errorf("%s: TEEC_ERROR_BUSY", name)
	case C.TEEC_ERROR_COMMUNICATION:
		return fmt.Errorf("%s: TEEC_ERROR_COMMUNICATION", name)
	case C.TEEC_ERROR_SECURITY:
		return fmt.Errorf("%s: TEEC_ERROR_SECURITY", name)
	case C.TEEC_ERROR_SHORT_BUFFER:
		return fmt.Errorf("%s: TEEC_ERROR_SHORT_BUFFER", name)
	case C.TEEC_ERROR_EXTERNAL_CANCEL:
		return fmt.Errorf("%s: TEEC_ERROR_EXTERNAL_CANCEL", name)
	case C.TEEC_ERROR_TARGET_DEAD:
		return fmt.Errorf("%s: TEEC_ERROR_TARGET_DEAD", name)
	case C.TEEC_ERROR_STORAGE_NO_SPACE:
		return fmt.Errorf("%s: TEEC_ERROR_STORAGE_NO_SPACE", name)
	}

	return fmt.Errorf("unknown error")
}

func (teec *TEEC) freeContext() {
	C.free(unsafe.Pointer(teec.context))
	teec.context = nil
}

func (teec *TEEC) freeSession() {
	C.free(unsafe.Pointer(teec.session))
	teec.session = nil
}
