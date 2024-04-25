package simpleaudiodecoder

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -L./ -lsimple_audio_decoder_rs
#include "simple_audio_decoder_rs.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type AudioClip struct {
	ptr unsafe.Pointer
}

type Pcm struct {
	ptr unsafe.Pointer
}

func NewPcmFromFile(file string) (*Pcm, error) {
	ptr := C.pcm_new_from_file(C.CString(file))

	if ptr == nil {
		return nil, fmt.Errorf("Failed to create pcm")
	}

	return &Pcm{ptr}, nil
}

func NewPcmFromData(data []byte) (*Pcm, error) {
	ptr := C.pcm_new_from_data(unsafe.Pointer(&data[0]), C.size_t(len(data)))

	if ptr == nil {
		return nil, fmt.Errorf("Failed to create pcm")
	}

	return &Pcm{ptr}, nil
}

func (pcm *Pcm) Free() {
	C.pcm_free(pcm.ptr)
}

func NewAudioClip(pcm *Pcm, target_sample_rate int, chunk_size int) (*AudioClip, error) {
	ptr := C.audio_clip_new(pcm.ptr, C.size_t(target_sample_rate), C.size_t(chunk_size))

	if ptr == nil {
		return nil, fmt.Errorf("Failed to create audio clip")
	}

	return &AudioClip{ptr}, nil
}

func (audioClip *AudioClip) Free() {
	C.audio_clip_free(audioClip.ptr)
}

func (audioClip *AudioClip) ResampleNext() (*[][]float32, bool, error) {
	ptr := C.audio_clip_resample_next(audioClip.ptr)
	if ptr == nil {
		return nil, true, fmt.Errorf("Failed to resample next chunk")
	}

	resampleResult := (*C.CResampleResult)(unsafe.Pointer(ptr))

	channels := int(resampleResult.channels)
	frames := int(resampleResult.frames)
	isDone := bool(resampleResult.is_done)

	infoSlice := unsafe.Slice((*uint)(unsafe.Pointer(resampleResult.buffer)), 3*channels)

	result := make([][]float32, channels)
	for i := 0; i < channels; i++ {
		pointer := unsafe.Pointer(uintptr(infoSlice[3*i+1]))
		result[i] = unsafe.Slice((*float32)(pointer), frames)
	}

	C.resample_result_free(ptr)

	return &result, isDone, nil
}
