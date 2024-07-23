package rgsw

import (
	"bufio"
	"io"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/utils/buffer"
)

// Ciphertext is a generic type for RGSW ciphertext.
type Ciphertext struct {
	Value [2]rlwe.GadgetCiphertext
}

// LevelQ returns the level of the modulus Q of the target.
func (ct Ciphertext) LevelQ() int {
	return ct.Value[0].LevelQ()
}

// LevelP returns the level of the modulus P of the target.
func (ct Ciphertext) LevelP() int {
	return ct.Value[0].LevelP()
}

// NewCiphertext allocates a new RGSW [Ciphertext] in the NTT domain.
func NewCiphertext(params rlwe.Parameters, levelQ, levelP, BaseTwoDecomposition int) (ct *Ciphertext) {
	return &Ciphertext{
		Value: [2]rlwe.GadgetCiphertext{
			*rlwe.NewGadgetCiphertext(params, 1, levelQ, levelP, BaseTwoDecomposition),
			*rlwe.NewGadgetCiphertext(params, 1, levelQ, levelP, BaseTwoDecomposition),
		},
	}
}

// BinarySize returns the serialized size of the object in bytes.
func (ct Ciphertext) BinarySize() int {
	return ct.Value[0].BinarySize() + ct.Value[1].BinarySize()
}

// WriteTo writes the object on an [io.Writer]. It implements the [io.WriterTo]
// interface, and will write exactly object.BinarySize() bytes on w.
//
// Unless w implements the [buffer.Writer] interface (see lattigo/utils/buffer/writer.go),
// it will be wrapped into a [bufio.Writer]. Since this requires allocations, it
// is preferable to pass a [buffer.Writer] directly:
//
//   - When writing multiple times to a [io.Writer], it is preferable to first wrap the
//     [io.Writer] in a pre-allocated [bufio.Writer].
//   - When writing to a pre-allocated var b []byte, it is preferable to pass
//     buffer.NewBuffer(b) as w (see lattigo/utils/buffer/buffer.go).
func (ct Ciphertext) WriteTo(w io.Writer) (n int64, err error) {
	switch w := w.(type) {
	case buffer.Writer:

		if n, err = ct.Value[0].WriteTo(w); err != nil {
			return
		}

		inc, err := ct.Value[1].WriteTo(w)

		return n + inc, err

	default:
		return ct.WriteTo(bufio.NewWriter(w))
	}
}

// ReadFrom reads on the object from an [io.Writer]. It implements the
// [io.ReaderFrom] interface.
//
// Unless r implements the [buffer.Reader] interface (see see lattigo/utils/buffer/reader.go),
// it will be wrapped into a [bufio.Reader]. Since this requires allocation, it
// is preferable to pass a [buffer.Reader] directly:
//
//   - When reading multiple values from a [io.Reader], it is preferable to first
//     first wrap [io.Reader] in a pre-allocated [bufio.Reader].
//   - When reading from a var b []byte, it is preferable to pass a buffer.NewBuffer(b)
//     as w (see lattigo/utils/buffer/buffer.go).
func (ct *Ciphertext) ReadFrom(r io.Reader) (n int64, err error) {
	switch r := r.(type) {
	case buffer.Reader:

		if n, err = ct.Value[0].ReadFrom(r); err != nil {
			return
		}

		inc, err := ct.Value[1].ReadFrom(r)

		return n + inc, err

	default:
		return ct.ReadFrom(bufio.NewReader(r))
	}
}

// MarshalBinary encodes the object into a binary form on a newly allocated slice of bytes.
func (ct Ciphertext) MarshalBinary() (p []byte, err error) {
	buf := buffer.NewBufferSize(ct.BinarySize())
	_, err = ct.WriteTo(buf)
	return buf.Bytes(), err
}

// UnmarshalBinary decodes a slice of bytes generated by
// [Ciphertext.MarshalBinary] or [Ciphertext.WriteTo] on the object.
func (ct *Ciphertext) UnmarshalBinary(p []byte) (err error) {
	_, err = ct.ReadFrom(buffer.NewBuffer(p))
	return
}

// Plaintext stores an RGSW plaintext value.
type Plaintext rlwe.GadgetPlaintext

// NewPlaintext creates a new RGSW plaintext from value, which can be either uint64, int64 or *[ring.Poly].
// Plaintext is returned in the NTT and Montgomery domain.
func NewPlaintext(params rlwe.Parameters, value interface{}, levelQ, levelP, BaseTwoDecomposition int) (*Plaintext, error) {
	gct, err := rlwe.NewGadgetPlaintext(params, value, levelQ, levelP, BaseTwoDecomposition)
	return &Plaintext{Value: gct.Value}, err
}
