package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
)

func main() {
	//func (m *Dense) UnmarshalBinaryFrom(r io.Reader) (int, error)
//	UnmarshalBinaryFrom decodes the binary form into the receiver and returns the number of bytes read and an error if any. It panics if the receiver is a non-empty Dense matrix.
	//
	//See MarshalBinary for the on-disk layout.
	//
	//Limited checks on the validity of the binary input are performed:
	//
	//- matrix.ErrShape is returned if the number of rows or columns is negative,
	//- an error is returned if the resulting Dense matrix is too
	//big for the current architecture (e.g. a 16GB matrix written by a
	//64b application and read back from a 32b application.)
//	UnmarshalBinary does not limit the size of the unmarshaled matrix, and so it should not be used on untrusted data.

//	func (m *Dense) Zero()
//	Zero sets all of the matrix elements to zero.

//	// DiagDense represents a diagonal matrix in dense storage format.
	//type DiagDense struct {
	//	mat blas64.Vector
	//}
//	DiagDense represents a diagonal matrix in dense storage format.

//	func NewDiagDense(n int, data []float64) *DiagDense
//	NewDiagDense creates a new Diagonal matrix with n rows and n columns. The length of data must be n or data must be nil, otherwise NewDiagDense will panic.
	//	NewDiagDense will panic if n is zero

//	func (d *DiagDense) At(i, j int) float64
//	At returns the element at row i, column j.

//	func (d *DiagDense) Bandwidth() (kl, ku int)
//	Bandwidth returns the upper and lower bandwidths of the matrix.
	//	These values are always zero for diagonal matrices.

//	func (d *DiagDense) Diag() int
//	Diag returns the dimension of the receiver.

//	func (d *DiagDense) DiagFrom(m Matrix)
//	DiagFrom copies the diagonal of m into the receiver.
	//	The receiver must be min(r, c) long or empty,
	//	otherwise DiagFrom will panic.

//	func (d *DiagDense) DiagView() Diagonal
//	DiagView returns the diagonal as a matrix backed by the original data.

//	func (d *DiagDense) Dims() (r, c int)
//	Dims returns the dimensions of the matrix.

//	func (d *DiagDense) IsEmpty() bool
//	IsEmpty returns whether the receiver is empty.
	//	Empty matrices can be the receiver for size-restricted operations.
	//	The receiver can be emptied using Reset.

//	func (d *DiagDense) RawBand() blas64.Band
//	RawBand returns the underlying data
	//	used by the receiver represented as a blas64.Band.
	//	Changes to elements in the receiver
	//	following the call will be reflected in returned blas64.Band.

//	func (d *DiagDense) RawSymBand() blas64.SymmetricBand
//	RawSymBand returns the underlying data
	//	used by the receiver represented as a blas64.SymmetricBand.
	//	Changes to elements in the receiver
	//	following the call will be reflected in returned blas64.Band.

//	func (d *DiagDense) Reset()
//	Reset empties the matrix so that
	//	it can be reused as the receiver of a dimensionally restricted operation.
//	Reset should not be used when the matrix shares backing data.
	//	See the Reseter interface for more information.

//	func (d *DiagDense) SetDiag(i int, v float64)
//	SetDiag sets the element at row i, column i to the value v.
	//	It panics if the location is outside the appropriate region of the matrix.

//	func (d *DiagDense) SymBand() (n, k int)
//	SymBand returns the number of rows/columns in the matrix, and the size of the bandwidth.

//	func (d *DiagDense) Symmetric() int
//	Symmetric implements the Symmetric interface.

//	func (d *DiagDense) T() Matrix
//	T returns the transpose of the matrix.

//	func (d *DiagDense) TBand() Banded
//	TBand performs an implicit transpose by returning the receiver inside a TransposeBand.

//	func (d *DiagDense) TTri() Triangular
//	TTri returns the transpose of the matrix. Note that Diagonal matrices are Upper by default.

//	func (d *DiagDense) TTriBand() TriBanded
//	TTriBand performs an implicit transpose by returning the receiver inside a TransposeTriBand.
	//	Note that Diagonal matrices are Upper by default.

//	func (d *DiagDense) Trace() float64
//	Trace returns the trace.

//	func (d *DiagDense) TriBand() (n, k int, kind TriKind)
//	TriBand returns the number of rows/columns in the matrix, the size of the bandwidth, and the orientation.
	//	Note that Diagonal matrices are Upper by default.

//	func (d *DiagDense) Triangle() (int, TriKind)
//	Triangle implements the Triangular interface.

//	func (d *DiagDense) Zero()
//	Zero sets all of the matrix elements to zero.

//	type Diagonal interface {
	//	Matrix
	//	// Diag returns the number of rows/columns in the matrix.
	//	Diag() int
	//
	//	// Bandwidth and TBand are included in the Diagonal interface
	//	// to allow the use of Diagonal types in banded functions.
	//	// Bandwidth will always return (0, 0).
	//	Bandwidth() (kl, ku int)
	//	TBand() Banded
	//
	//	// Triangle and TTri are included in the Diagonal interface
	//	// to allow the use of Diagonal types in triangular functions.
	//	Triangle() (int, TriKind)
	//	TTri() Triangular
	//
	//	// Symmetric and SymBand are included in the Diagonal interface
	//	// to allow the use of Diagonal types in symmetric and banded symmetric
	//	// functions respectively.
	//	Symmetric() int
	//	SymBand() (n, k int)
//	// TriBand and TTriBand are included in the Diagonal interface
	//	// to allow the use of Diagonal types in triangular banded functions.
	//	TriBand() (n, k int, kind TriKind)
	//	TTriBand() TriBanded
	//}
//	Diagonal represents a diagonal matrix, that is a square matrix that only has non-zero terms on the diagonal.

//	// Eigen is a type for creating and using the eigenvalue decomposition of a dense matrix.
	//type Eigen struct {
	//	n int // The size of the factorized matrix.
	//
	//	kind EigenKind
	//
	//	values   []complex128
	//	rVectors *CDense
	//	lVectors *CDense
	//}
//	Eigen is a type for creating and using
	//	the eigenvalue decomposition of a dense matrix.

	a14 := mat.NewDense(2, 2, []float64{
		1, -1,
		1, 1,
	})
	fmt.Printf("A = %v\n\n", mat.Formatted(a14, mat.Prefix("    ")))

	var eig mat.Eigen

	ok := eig.Factorize(a14,mat.EigenLeft)

	if !ok {
		log.Fatal("Eigendecomposition failed")
	}
	fmt.Printf("Eigenvalues of A:\n%v\n", eig.Values(nil))

//	func (e *Eigen) Factorize(a Matrix, kind EigenKind) (ok bool)
//Factorize computes the eigenvalues of the square matrix a, and optionally the eigenvectors.
	//A right eigenvalue/eigenvector combination is defined by:
//	A * x_r = λ * x_r
//where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.where x_r is the column vector called an eigenvector, and λ is the corresponding eigenvalue.
//Similarly, a left eigenvalue/eigenvector combination is defined by
//x_l * A = λ * x_l

//The eigenvalues, but not the eigenvectors, are the same for both decompositions.
//Typically eigenvectors refer to right eigenvectors.
//In all cases, Factorize computes the eigenvalues of the matrix. kind specifies which of the eigenvectors, if any, to compute. See the EigenKind documentation for more information.
//Eigen panics if the input matrix is not square.

//Factorize returns whether the decomposition succeeded.
//If the decomposition failed,
//methods that require a successful factorization will panic.


//func (e *Eigen) Kind() EigenKind
//Kind returns the EigenKind of the decomposition. If no decomposition has been computed, Kind returns -1.

//func (e *Eigen) LeftVectorsTo(dst *CDense)
//LeftVectorsTo stores the left eigenvectors of the decomposition into the columns of dst.
//The computed eigenvectors are normalized to have Euclidean norm equal to 1 and largest component real.

//If dst is empty, LeftVectorsTo will resize dst to be n×n.
//When dst is non-empty, LeftVectorsTo will panic if dst is not n×n.
//LeftVectorsTo will also panic
//if the left eigenvectors were not computed during the factorization,
//or if the receiver does not contain a successful factorization

//func (e *Eigen) Values(dst []complex128) []complex128
//Values extracts the eigenvalues of the factorized matrix.
//If dst is non-nil, the values are stored in-place into dst.
//In this case dst must have length n, otherwise Values will panic.
//If dst is nil, then a new slice
//will be allocated of the proper length and filed with the eigenvalues.
//Values panics if the Eigen decomposition was not successful.

//func (e *Eigen) VectorsTo(dst *CDense)
//VectorsTo stores the right eigenvectors of the decomposition
//into the columns of dst.
//The computed eigenvectors are normalized to have Euclidean norm equal to 1
//and largest component real.

//If dst is empty, VectorsTo will resize dst to be n×n.
//When dst is non-empty, VectorsTo will panic if dst is not n×n.
//VectorsTo will also panic if the eigenvectors were not computed
//during the factorization,
//or if the receiver does not contain a successful factorization.

//type EigenKind int
//EigenKind specifies the computation of eigenvectors during factorization.
//const (
	//	// EigenNone specifies to not compute any eigenvectors.
	//	EigenNone EigenKind = 0
	//	// EigenLeft specifies to compute the left eigenvectors.
	//	EigenLeft EigenKind = 1 << iota
	//	// EigenRight specifies to compute the right eigenvectors.
	//	EigenRight
	//	// EigenBoth is a convenience value for computing both eigenvectors.
	//	EigenBoth EigenKind = EigenLeft | EigenRight
	//)

//	// EigenSym is a type for
//	creating and manipulating the Eigen decomposition of
	//// symmetric matrices.
	//type EigenSym struct {
	//	vectorsComputed bool
	//
	//	values  []float64
	//	vectors *Dense
	//}

//	EigenSym is a type for creating and manipulating the Eigen decomposition of symmetric matrices
	a15 := mat.NewSymDense(2, []float64{
		7, 0.5,
		0.5, 1,
	})
	fmt.Printf("A = %v\n\n", mat.Formatted(a15, mat.Prefix("    ")))

	//var eigsym mat.EigenSym
	//ok15 := eigsym.Factorize(a15, true)
	var eigsym mat.EigenSym
	ok15 := eigsym.Factorize(a15,true)
	if !ok15 {
		log.Fatal("Symmetric eigendecomposition failed")
	}
	fmt.Printf("Eigenvalues of A:\n%1.3f\n\n", eigsym.Values(nil))

	var ev mat.Dense
	eigsym.VectorsTo(&ev)

	fmt.Printf("Eigenvectors of A:\n%1.3f\n\n", mat.Formatted(&ev))

//	func (e *EigenSym) Factorize(a Symmetric, vectors bool) (ok bool)
//Factorize computes the eigenvalue decomposition of the symmetric matrix a.
//The Eigen decomposition is defined as:
//A = P * D * P^-1

//where D is a diagonal matrix containing the eigenvalues of the matrix,
//and P is a matrix of the eigenvectors of A.
//Factorize computes the eigenvalues in ascending order.
//If the vectors input argument is false,
//the eigenvectors are not computed.

//Factorize returns whether the decomposition succeeded.
//If the decomposition failed,
//methods that require a successful factorization will panic.

//func (e *EigenSym) Values(dst []float64) []float64
//Values extracts the eigenvalues of the factorized matrix.
//If dst is non-nil, the values are stored in-place into dst.
//In this case dst must have length n, otherwise Values will panic.
//If dst is nil, then a new slice
//will be allocated of the proper length and filled with the eigenvalues.
//Values panics if the Eigen decomposition was not successful.

//func (e *EigenSym) VectorsTo(dst *Dense)
//VectorsTo stores the eigenvectors of the decomposition into the columns of dst.
//If dst is empty, VectorsTo will resize dst to be n×n.
//When dst is non-empty, VectorsTo will panic if dst is not n×n.
//VectorsTo will also panic if the eigenvectors were not computed
//during the factorization,
//or if the receiver does not contain a successful factorization.

//type Error struct {
	//	// contains filtered or unexported fields
	//}
	//type Error struct{ string }
	//
	//func (err Error) Error() string { return err.string }
//	Error represents matrix handling errors.
//	These errors can be recovered by Maybe wrappers.

//func (err Error) Error() string
//type ErrorStack struct {
	//	Err error
	//
	//	// StackTrace is the stack trace
	//	// recovered by Maybe, MaybeFloat
	//	// or MaybeComplex.
	//	StackTrace string
	//}

//
//ErrorStack represents matrix handling errors that have been recovered by Maybe wrappers
//func (err ErrorStack) Error() string









}
