package UniMessCore

import (
	"unsafe"
	"runtime"
)

/*
#cgo LDFLAGS: -lunimess
void unimess_init();
void * unimess_config_chain_generate(unsigned int size);
unsigned char * unimess_config_chain_dump(void *cc, unsigned int *len_out);
void * unimess_config_chain_load(unsigned char *data, unsigned int len);
void unimess_config_chain_destroy(void *cc);
void * unimess_config_chain_get_protocol_chain(void *cc);
unsigned char * unimess_protocol_chain_encode_packet(void *pc, unsigned int *len_out, unsigned char *pkt, unsigned int len);
unsigned char * unimess_protocol_chain_decode_packet(void *pc, unsigned int *len_out, unsigned char *pkt, unsigned int len);
void unimess_protocol_chain_destroy(void *pc);
void unimess_binary_buffer_destroy(unsigned char *buf);
*/
import "C"

type ConfigChain struct {
	handle unsafe.Pointer
}

type ProtocolChain struct {
	handle unsafe.Pointer
}

func LoadConfigChain(data []byte) *ConfigChain {
	handle := C.unimess_config_chain_load(
		(*C.uchar)(&data[0]),
		(C.uint)(len(data)),
	)
	if handle == nil {
		return nil
	}

	ret := &ConfigChain {
		handle: handle,
	}
	runtime.SetFinalizer(ret, _destroyConfigChain)
	return ret
}

func GenerateConfigChain(size int) *ConfigChain {
	handle := C.unimess_config_chain_generate(C.uint(size))
	ret := &ConfigChain {
		handle: handle,
	}
	runtime.SetFinalizer(ret, _destroyConfigChain)
	return ret
}

func _destroyConfigChain(cc *ConfigChain) {
	if cc.handle != nil {
		C.unimess_config_chain_destroy(cc.handle)
		cc.handle = nil
	}
}

func (cc *ConfigChain) Dump() []byte {
	var out_len uint32
	data := C.unimess_config_chain_dump(cc.handle, (*C.uint)(&out_len))
	ret := C.GoBytes(unsafe.Pointer(data), C.int(out_len))
	C.unimess_binary_buffer_destroy(data)
	return ret
}

func (cc *ConfigChain) GetProtocolChain() *ProtocolChain {
	handle := C.unimess_config_chain_get_protocol_chain(cc.handle)
	ret := &ProtocolChain {
		handle: handle,
	}
	runtime.SetFinalizer(ret, _destroyProtocolChain)
	return ret
}

func (pc *ProtocolChain) EncodePacket(pkt []byte) []byte {
	var out_len uint32
	data := C.unimess_protocol_chain_encode_packet(
		pc.handle,
		(*C.uint)(&out_len),
		(*C.uchar)(&pkt[0]),
		(C.uint)(len(pkt)),
	)
	if data == nil {
		return nil
	}
	ret := C.GoBytes(unsafe.Pointer(data), C.int(out_len))
	C.unimess_binary_buffer_destroy(data)
	return ret
}

func (pc *ProtocolChain) DecodePacket(pkt []byte) []byte {
	var out_len uint32
	data := C.unimess_protocol_chain_decode_packet(
		pc.handle,
		(*C.uint)(&out_len),
		(*C.uchar)(&pkt[0]),
		(C.uint)(len(pkt)),
	)
	if data == nil {
		return nil
	}
	ret := C.GoBytes(unsafe.Pointer(data), C.int(out_len))
	C.unimess_binary_buffer_destroy(data)
	return ret
}

func _destroyProtocolChain(pc *ProtocolChain) {
	if pc.handle != nil {
		C.unimess_protocol_chain_destroy(pc.handle)
		pc.handle = nil
	}
}

func init() {
	C.unimess_init()
}
