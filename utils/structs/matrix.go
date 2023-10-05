package structs

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tuneinsight/lattigo/v4/utils/buffer"
)

// Vector is a struct wrapping a doube slice of components of type T.
// T can be:
// - uint, uint64, uint32, uint16, uint8/byte, int, int64, int32, int16, int8, float64, float32.
// - Or any object that implements CopyNewer, CopyNewer, BinarySizer, io.WriterTo or io.ReaderFrom 
//   depending on the method called.
type Matrix[T any] [][]T

// CopyNew returns a deep copy of the object.
// If T is a struct, this method requires that T implements CopyNewer.
func (m Matrix[T]) CopyNew() (mcpy Matrix[T]) {

	var t T
	switch any(t).(type) {
	case uint, uint64, uint32, uint16, uint8, int, int64, int32, int16, int8, float64, float32:

		mcpy = Matrix[T](make([][]T, len(m)))

		for i := range m {

			mcpy[i] = make([]T, len(m[i]))
			copy(mcpy[i], m[i])
		}

	default:
		if _, isCopiable := any(t).(CopyNewer[T]); !isCopiable {
			panic(fmt.Errorf("matrix component of type %T does not comply to %T", t, new(CopyNewer[T])))
		}

		mcpy = Matrix[T](make([][]T, len(m)))

		for i := range m {

			mcpy[i] = make([]T, len(m[i]))

			for j := range m[i] {
				mcpy[i][j] = *any(&m[i][j]).(CopyNewer[T]).CopyNew()
			}
		}
	}

	return
}

// BinarySize returns the serialized size of the object in bytes.
// If T is a struct, this method requires that T implements BinarySizer.
func (m Matrix[T]) BinarySize() (size int) {

	size += 8

	for _, v := range m {
		/* #nosec G601 -- Implicit memory aliasing in for loop acknowledged */
		size += (*Vector[T])(&v).BinarySize()
	}

	return
}

// WriteTo writes the object on an io.Writer. It implements the io.WriterTo
// interface, and will write exactly object.BinarySize() bytes on w.
//
// If T is a struct, this method requires that T implements io.WriterTo.
//
// Unless w implements the buffer.Writer interface (see lattigo/utils/buffer/writer.go),
// it will be wrapped into a bufio.Writer. Since this requires allocations, it
// is preferable to pass a buffer.Writer directly:
//
//   - When writing multiple times to a io.Writer, it is preferable to first wrap the
//     io.Writer in a pre-allocated bufio.Writer.
//   - When writing to a pre-allocated var b []byte, it is preferable to pass
//     buffer.NewBuffer(b) as w (see lattigo/utils/buffer/buffer.go).
func (m Matrix[T]) WriteTo(w io.Writer) (n int64, err error) {

	switch w := w.(type) {
	case buffer.Writer:

		var inc int64
		if inc, err = buffer.WriteAsUint64[int](w, len(m)); err != nil {
			return inc, fmt.Errorf("buffer.WriteAsUint64[int]: %w", err)
		}
		n += inc

		for _, v := range m {
			if inc, err = Vector[T](v).WriteTo(w); err != nil {
				var t T
				return n + inc, fmt.Errorf("structs.Vector[%T].WriteTo: %w", t, err)
			}
			n += inc
		}

		return n, w.Flush()

	default:
		return m.WriteTo(bufio.NewWriter(w))
	}
}

// ReadFrom reads on the object from an io.Writer. It implements the
// io.ReaderFrom interface.
//
// If T is a struct, this method requires that T implements io.ReaderFrom.
//
// Unless r implements the buffer.Reader interface (see see lattigo/utils/buffer/reader.go),
// it will be wrapped into a bufio.Reader. Since this requires allocation, it
// is preferable to pass a buffer.Reader directly:
//
//   - When reading multiple values from a io.Reader, it is preferable to first
//     first wrap io.Reader in a pre-allocated bufio.Reader.
//   - When reading from a var b []byte, it is preferable to pass a buffer.NewBuffer(b)
//     as w (see lattigo/utils/buffer/buffer.go).
func (m *Matrix[T]) ReadFrom(r io.Reader) (n int64, err error) {

	switch r := r.(type) {
	case buffer.Reader:

		var inc int64

		var size int

		if n, err = buffer.ReadAsUint64[int](r, &size); err != nil {
			return int64(n), fmt.Errorf("buffer.ReadAsUint64[int]: %w", err)
		}

		if cap(*m) < size {
			*m = make([][]T, size)
		}

		*m = (*m)[:size]

		for i := range *m {
			if inc, err = (*Vector[T])(&(*m)[i]).ReadFrom(r); err != nil {
				var t T
				return n + inc, fmt.Errorf("structs.Vector[%T].ReadFrom: %w", t, err)
			}
			n += inc
		}

		return n, nil

	default:
		return m.ReadFrom(bufio.NewReader(r))
	}
}

// MarshalBinary encodes the object into a binary form on a newly allocated slice of bytes.
// If T is a struct, this method requires that T implements io.WriterTo.
func (m Matrix[T]) MarshalBinary() (p []byte, err error) {
	buf := buffer.NewBufferSize(m.BinarySize())
	_, err = m.WriteTo(buf)
	return buf.Bytes(), err
}

// UnmarshalBinary decodes a slice of bytes generated by
// MarshalBinary or WriteTo on the object.
// If T is a struct, this method requires that T implements io.ReaderFrom.
func (m *Matrix[T]) UnmarshalBinary(p []byte) (err error) {
	_, err = m.ReadFrom(buffer.NewBuffer(p))
	return
}

// Equal performs a deep equal.
// If T is a struct, this method requires that T implements Equatable.
func (m Matrix[T]) Equal(other Matrix[T]) bool {
	
	var t T
	switch any(t).(type) {
	case uint, uint64, uint32, uint16, uint8, int, int64, int32, int16, int8, float64, float32:
	default:
		if _, isEquatable := any(t).(Equatable[T]); !isEquatable {
			panic(fmt.Errorf("matrix component of type %T does not comply to %T", t, new(Equatable[T])))
		}
	}

	isEqual := true
	for i := range m {
		isEqual = isEqual && Vector[T](m[i]).Equal(Vector[T](other[i]))
	}

	return isEqual
}
