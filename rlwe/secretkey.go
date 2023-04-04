package rlwe

import (
	"io"

	"github.com/google/go-cmp/cmp"
	"github.com/tuneinsight/lattigo/v4/rlwe/ringqp"
)

// SecretKey is a type for generic RLWE secret keys.
// The Value field stores the polynomial in NTT and Montgomery form.
type SecretKey struct {
	Value ringqp.Poly
}

// NewSecretKey generates a new SecretKey with zero values.
func NewSecretKey(params Parameters) *SecretKey {
	return &SecretKey{Value: *params.RingQP().NewPoly()}
}

func (sk *SecretKey) Equal(other *SecretKey) bool {
	return cmp.Equal(sk.Value, other.Value)
}

// LevelQ returns the level of the modulus Q of the target.
func (sk *SecretKey) LevelQ() int {
	return sk.Value.Q.Level()
}

// LevelP returns the level of the modulus P of the target.
// Returns -1 if P is absent.
func (sk *SecretKey) LevelP() int {
	if sk.Value.P != nil {
		return sk.Value.P.Level()
	}

	return -1
}

// CopyNew creates a deep copy of the receiver secret key and returns it.
func (sk *SecretKey) CopyNew() *SecretKey {
	if sk == nil {
		return nil
	}
	return &SecretKey{*sk.Value.CopyNew()}
}

// BinarySize returns the size in bytes that the object once marshalled into a binary form.
func (sk *SecretKey) BinarySize() (dataLen int) {
	return sk.Value.BinarySize()
}

// MarshalBinary encodes the object into a binary form on a newly allocated slice of bytes.
func (sk *SecretKey) MarshalBinary() (data []byte, err error) {
	data = make([]byte, sk.BinarySize())
	if _, err = sk.Read(data); err != nil {
		return nil, err
	}
	return
}

// WriteTo writes the object on an io.Writer.
// To ensure optimal efficiency and minimal allocations, the user is encouraged
// to provide a struct implementing the interface buffer.Writer, which defines
// a subset of the method of the bufio.Writer.
// If w is not compliant to the buffer.Writer interface, it will be wrapped in
// a new bufio.Writer.
// For additional information, see lattigo/utils/buffer/writer.go.
func (sk *SecretKey) WriteTo(w io.Writer) (n int64, err error) {
	return sk.Value.WriteTo(w)
}

// Read encodes the object into a binary form on a preallocated slice of bytes
// and returns the number of bytes written.
func (sk *SecretKey) Read(data []byte) (ptr int, err error) {
	return sk.Value.Read(data)
}

// UnmarshalBinary decodes a slice of bytes generated by MarshalBinary
// or Read on the object.
func (sk *SecretKey) UnmarshalBinary(data []byte) (err error) {
	_, err = sk.Write(data)
	return
}

// ReadFrom reads on the object from an io.Writer.
// To ensure optimal efficiency and minimal allocations, the user is encouraged
// to provide a struct implementing the interface buffer.Reader, which defines
// a subset of the method of the bufio.Reader.
// If r is not compliant to the buffer.Reader interface, it will be wrapped in
// a new bufio.Reader.
// For additional information, see lattigo/utils/buffer/reader.go.
func (sk *SecretKey) ReadFrom(r io.Reader) (n int64, err error) {
	return sk.Value.ReadFrom(r)
}

// Write decodes a slice of bytes generated by MarshalBinary or
// Read on the object and returns the number of bytes read.
func (sk *SecretKey) Write(data []byte) (ptr int, err error) {
	return sk.Value.Write(data)
}
