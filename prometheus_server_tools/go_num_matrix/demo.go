package main

import (
	"fmt"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
	"gonum.org/v1/gonum/mat"
	"log"
	"math/rand"
)

func matPrint(X mat.Matrix){
	fa := mat.Formatted(X,mat.Prefix(""),mat.Squeeze())
	fmt.Printf("%v\n",fa)
}

type row []float64
type column []float64

const(
	//CondNorm is the matrix norm used for computing the condition numer by routines in the matrix packages
	CondNorm = lapack.MaxRowSum

	//CondTrams is the norm used to compute on A^T to get the same result as computing CondNorm on A.
	CondNormTrans = lapack.MaxColumnSum
)

const ConditionTolerance = 1e16
//conditionTolerance is the tolerance limit of the condition number. If the condition
//number is above this value, the matrix is considered singluar.

//var (
//	ErrNegativeDimension   = Error{"mat: negative dimension"}
//	ErrIndexOutOfRange     = Error{"mat: index out of range"}
//	ErrReuseNonEmpty       = Error{"mat: reuse of non-empty matrix"}
//	ErrRowAccess           = Error{"mat: row index out of range"}
//	ErrColAccess           = Error{"mat: column index out of range"}
//	ErrVectorAccess        = Error{"mat: vector index out of range"}
//	ErrZeroLength          = Error{"mat: zero length in matrix dimension"}
//	ErrRowLength           = Error{"mat: row length mismatch"}
//	ErrColLength           = Error{"mat: col length mismatch"}
//	ErrSquare              = Error{"mat: expect square matrix"}
//	ErrNormOrder           = Error{"mat: invalid norm order for matrix"}
//	ErrSingular            = Error{"mat: matrix is singular"}
//	ErrShape               = Error{"mat: dimension mismatch"}
//	ErrIllegalStride       = Error{"mat: illegal stride"}
//	ErrPivot               = Error{"mat: malformed pivot list"}
//	ErrTriangle            = Error{"mat: triangular storage mismatch"}
//	ErrTriangleSet         = Error{"mat: triangular set out of bounds"}
//	ErrBandwidth           = Error{"mat: bandwidth out of range"}
//	ErrBandSet             = Error{"mat: band set out of bounds"}
//	ErrDiagSet             = Error{"mat: diagonal set out of bounds"}
//	ErrSliceLengthMismatch = Error{"mat: input slice length mismatch"}
//	ErrNotPSD              = Error{"mat: input not positive symmetric definite"}
//	ErrFailedEigen         = Error{"mat: eigendecomposition not successful"}
//)

var(
	ErrNegativeDimension = mat.ErrNegativeDimension //	ErrNegativeDimension   = Error{"mat: negative dimension"}
	ErrIndexOutOfRange = mat.ErrIndexOutOfRange //	ErrIndexOutOfRange     = Error{"mat: index out of range"}
	ErrReuseNonEmpty = mat.ErrReuseNonEmpty
	ErrRowAccess = mat.ErrRowAccess //	ErrRowAccess           = Error{"mat: row index out of range"}
	ErrColAccess = mat.ErrColAccess //	ErrColAccess           = Error{"mat: column index out of range"}
	ErrVectorAccess = mat.ErrVectorAccess //	ErrVectorAccess        = Error{"mat: vector index out of range"}
	ErrZeroLength = mat.ErrZeroLength //	ErrZeroLength          = Error{"mat: zero length in matrix dimension"}
	ErrRowLength = mat.ErrRowLength //	ErrRowLength           = Error{"mat: row length mismatch"}
	ErrColLength = mat.ErrColLength //	ErrColLength           = Error{"mat: col length mismatch"}
	ErrSquare = mat.ErrSquare //	ErrSquare              = Error{"mat: expect square matrix"}
	ErrNormOrder = mat.ErrNormOrder //	ErrNormOrder           = Error{"mat: invalid norm order for matrix"}
	ErrSingular = mat.ErrSingular //	ErrSingular            = Error{"mat: matrix is singular"}
	ErrShape = mat.ErrShape //	ErrShape               = Error{"mat: dimension mismatch"}
	ErrTriangle = mat.ErrTriangle //	ErrTriangle            = Error{"mat: triangular storage mismatch"}
	ErrTriangleSet = mat.ErrTriangleSet //	ErrTriangleSet         = Error{"mat: triangular set out of bounds"}
	ErrBandwidth = mat.ErrBandwidth  //	ErrBandwidth           = Error{"mat: bandwidth out of range"}
	ErrBandSet = mat.ErrBandSet //	ErrBandSet             = Error{"mat: band set out of bounds"}
	ErrDiagSet = mat.ErrDiagSet //	ErrDiagSet             = Error{"mat: diagonal set out of bounds"}
	ErrSliceLengthMismatch = mat.ErrSliceLengthMismatch
	ErrNotPSD = mat.ErrNotPSD //	ErrNotPSD              = Error{"mat: input not positive symmetric definite"}
	ErrFailedEigen = mat.ErrFailedEigen //	ErrFailedEigen         = Error{"mat: eigendecomposition not successful"}
)

func main() {
	zero := mat.NewDense(3,5,nil) //allocate a zeroed real matrix of size 3x5
	matPrint(zero)

	//generate a 6x6 matrix of random values
	data := make([]float64,36)
	for i := range data{
		data[i] = rand.NormFloat64()
	}

	a := mat.NewDense(6,6,data)
	matPrint(a)



	tr := mat.Trace(a)
	fmt.Println(tr)

	//implemented as methods when the operations modifies the receiver
	matPrint(zero)
	zero.Copy(a)
	matPrint(zero)

	var c mat.Dense// the zero value of mat.Dense is the empty matrix and the size can be assumed to be correct for the results
	//var d *mat.Dense
	//(*d).Mul(a,a)
	c.Mul(a,a)
	fmt.Println(c)
	//fmt.Println(*d)
	//matPrint(c)

	//var lu mat.LU
	//lu.Factorize(a)
	//fmt.Println(lu)

	// perform the cross product of [1,2,3,4] and [1,2,3]
	r := row{1,2,3,4}
	c2 := column{1,2,3}

	var m mat.Dense
	m.Mul(c2,r)

	matPrint(&m)
	fmt.Println(mat.Formatted(&m))

	//mat.CEqual(r,c2) //whether the matrices a and b have the same size and are element-wise equal.

	//mat.CEqualApprox(a,b CMatrix,epsilon float64) bool // whether the matrices a and b have the same size
	// and contain all equal elements with tolerance for element-wise equality specified by epsilon. Matrices with non-equal shapes are not equal

	//func Col(dst[]float64, j int, a mat.Matrix) []float64 // copies the elements in the jth column of the matrix into the slice dst.
	// The length of the provided slice must equal the number of rows, unless the slice is nil in which case a new slice is first allocated.
	//This example copies the second colum of a matrix into col, allocating a new slice of float64
	m2 := mat.NewDense(3,3,[]float64{
		2.0,9.0,3.0,
		4.5,6.7,8.0,
		1.2,3.0,6.0,
	})

	col := mat.Col(nil,1,m2)
	fmt.Printf("col = %#v",col)
//	col = []float64{9, 6.7, 3}

	//func Cond(a mat.Matrix,norm float64) float64
	// returns the condition number of the given matrix under the given norm
	// The condition number must be based on the 1-norm,2-norm or inf-norm
	//Cond will panic with matrix.ErrShape if the matrix has zero size

	//BUG(btracey): the computation of the 1-norm and inf-norm for non-square matrics is inaccurate
	//altough is typically the right order of magnitude. While the value returned will
	// change with the resolution of this bug. the result from Cond will match the condition number
	// used internally.

	//func Det(a mat.Matrix) float64
	// Det returns the determiant of the matrix a. In many expressions using LogDet will be more
	// numerically stable.
	//func Dot(a,b mat.Vector) float64
	// Dot returns the sum of the element-wise product of a and b. Dot panics if the
	// matrix sizes are unequal.

	//func Equal(a,b mat.Matrix) bool

	//Equal returns whether the matrices a and b have the same size and are element-size equal.
	//func EqualApprox(a,b mat.Matrix, epsilon float64) bool
	// EqualApprox returns whether the matrices a and b have the same size and contain all
	// equal elements with tolerance for element-wise equality specified by epsilon. Matrices with
	// non-equal shapes are not equal.

	//func Formatted(m mat.Matrix,options ...mat.FormatOption) fmt.Formatter
	//Formatted returns a fmt.Formatter for the matrix m using the given options

	a2 := mat.NewDense(3,3,[]float64{1,2,3,0,4,5,0,0,6})
	//Create a matrix formatting value with a prefix and calculating each column
	// with individually...
	fa := mat.Formatted(a2,mat.Prefix(" "),mat.Squeeze())

	// and then point with and without zero value elements.
	fmt.Printf("with all values:a =\n %v\n\n",fa)
	fmt.Printf("with only non-zero values:a=\n % v\n\n",fa)

	//Modify the matrix...
	a.Set(0,2,0)

	// and print it without zero value elements.
	fmt.Printf("after modification with only non-zero values:a =\n % v\n\n",fa)

	//Modify the matrix again
	a.Set(0,2,123.456)

	// and print it using scientific notation for large exponents.
	fmt.Printf("after modification with scientific notation:a=\n %.2g\n\n",fa)
	//see golang.org/pkg/fmt/floating-point verbs for a comprehensive list.

	//func Inner(x Vector, a mat.Matrix,y mat.Vector) float64
	//computes: x^T A y
	//between the vectors x and y with matrix A, where x and y are treated as columns vectors
	// this is only a true inner product if A is symmetric positive definite, though the operation
	// works for any matrix A.
	// Inner panics if x.Len != m or y.Len != n when A is an mxn matrix.
	//func LogDet(a mat.Matrix) (det float64,sign float64)
	// LogDet returns the log of the determinant and the sign of the determinant of the matrix
	// that has been factorized. Numerical stability in product and division expressions is generally improved by working in log space
	//func Max(a mat.Matrix) float64
	// Matrix returns the largest element value of the matrix A. Max will panic with matrix.ErrShape if the matrix has zero Size.

	//func Maybe(fn func()) (err error)
	// Maybe will recover a panic with a type mat.Error from fn, and return this error as the Err field of an ErrorStack.
	// The stack trace for the panicking function will be recovered and placed in the StackTrace field.
	// Any other error is re-panicked.

	//func MaybeComplex(fn func() complex128) (f complex128,err error)
	// MaybeComplex will recover a panic with a type mat.Errir from fn, and return
	// this error as the Err field of an ErrorStack. The stack trace for the panicking
	// function will be recovered and placed in the StackTrace field. Any other error is
	// re-panicked.

	//func Maybefloat(fn func() float64) (f float64,err error)

	// func Min(a matrix) float64
	// Min returns the smallest element value of the matrix A. Min will panic with matrix.Errorshape if the matrix hae zero size.

	//func Norm(a Matrix,norm float64) float64

	// Norm returns the specified norm of the matrix A. Vaild norms are:
	//1-norm: the maximum absolute column sum
	// 2-norm: The Frobenius norm, the square root of the sum of the squares of the elements
	//Inf: the maximum absolute row sum
	//Norm will panic sith ErrorNormOrder if an illegal norm order is specified and with ErrShape if the matrix has zero size.
	// func ROw(dst []float64, i int, a Matrix) []float64
	// Row copies the elements in the i-th row of the matrix into the slice dst.
	// The length of the provided slice must equal the number of columns, unless
	// the slice is nil in which case a new slice is first allocated.
	m3 := mat.NewDense(3,3,[]float64{
		2.0,9.0,3.0,
		4.5,6.7,8.0,
		1.2,3.0,6.0,
	})

	row := mat.Row(nil,2,m3)
	fmt.Printf("row = %#v",row)

	// func Sum(a Matrix) float64
	// sum returns the sum of the elements of the matrix
	//func Trace(a Matrix) float64
	// Trace returns the trace of the matix. Trace will panic if the matrix is not square
	// if a is a Tracer, its Trace method will be used  to calculate the matrix trace.

	//Types:
//	1. type BandCholesky struct {
	//	// The chol pointer must never be retained as a pointer outside the Cholesky
	//	// struct, either by returning chol outside the struct or by setting it to
	//	// a pointer coming from outside. The same prohibition applies to the data
	//	// slice within chol.
	//	chol *TriBandDense
	//	cond float64
	//}
	//BandCholesky is a symmetric positive-definite band matrix represented by its Cholesky decomposition.
	// Note that this matrix representation is useful for certain operations, in particular finding solutions to linear equations.
	// It is very inefficient at other operations, in particular At is slow.

	// BandCholesky methods may only be called on a value that hase been successfully initialized by a call to Factorize that has returned true.
	// Calls to methods of an unsuccessful Cholesky factorization will panic

	//func (ch *mat.BandCholesky) At(i,j int) float64
	// At returns the element at row i, column j

	//func (ch *mat.BandCholesky) Bandwidth() (kl,ku int)
	//Bandwidth returns the lower and upper bandwidth values for the matrix. The total bandwidth of the matrix is kl+ku+1.

	// func (ch *BandCholesky) Cond() float64
	// Cond returns the condition number of the factorized matrix.

	// func (ch *BandCholesky) Det() float64
	// Det returns the determinant of the matrix that hase been factorized.

	//func (ch *mat.BandCholesky) Dims() float64
	// Dims returns the dimensions of the matrix.

	//func(ch *mat.BandCholesky) Factorize(a mat.SymBanded) (ok bool)
	// Factorize calculates the Cholesky decomposition of the Matrix A and returns whether
	// the matrix is positive definite. If Factorize returns false, the factorization must not be used.

	//func(ch *mat.BandCholesky) IsEmpty() bool
	// IsEmpty returns whether the receiver is empty. Empty matrices can be the receiver for dimensionally restricted operations
	// The receiver can be emptied using Reset.

	//func(ch *mat.BandCholesky) LogDet() float64
	// LogDet returns the log of the determinant of the matrix that hse been factorized.


	//func(ch *mat.BandCholesky)Reset()
	// Reset resets tje factorization so that if can be reused as the receiver of a dimensionally restricted operations.

	// func(ch *BandCholesky) SolveTo(dst *Dense, b Matrix)error
	// SolveTO finds the Matrix X that solves: A*X = B where A is represented by the
	// Cholesky decomposition. The result is stored in-place into dst.
	// If the Cholesky decomposition is singular or near-signular a Condition error is returned
	// See the documentation for Condition for more information

	//func(ch *mat.BandCholesky) SolveVecTo (dst *VecDense, b Vector) error
	//SolveVecTo finds the vector x that solves A*x = b where A is represented by the
	// Cholesky decomposition. The result is stored in-place into dst. If the Cholesky
	// decomposition is singular or near-singular a Condition error is returned. See the documentation
	// for Condition for more information.

	//func (ch *mat.BandCholesky) SymBand() (n,k int)
	// SymBand returns the number of rows/columns in the matrix, and the size of the bandwidth. The total
	// bandwidth of the matrix is 2*k+1

	//func (ch *mat.BandCholesky) Symmetric() int
	// Symmetric implements the Symmetric interface and returns the number of rows
	// in the matrix (this is also the number of columns).

	// func (ch *BandCholesky) T() Matrix
	// T returns the receiver, the transpose of a symmetric matrix

	//func (ch *mat.BandCholesky) TBand() blas64.Band
	// TBand returns the receiver, the  transpose of a symmetric band matrix.

//	func NewBandDense(r, c, kl, ku int, data []float64) *BandDense
//NewBandDense creates a new Band matrix with r rows and c columns.
//If data == nil, a new slice is allocated for the backing slice.
//If len(data) == min(r, c+kl)*(kl+ku+1), data is used as the backing slice,
//and changes to the elements of the returned BandDense will be reflected in data.
//If neither of these is true, NewBandDense will panic.
//kl must be at least zero and less r, and ku must be at least zero and less than c,
//otherwise NewBandDense will panic.
//NewBandDense will panic if either r or c is zero.
//The data must be arranged in row-major order constructed by removing the zeros
//from the rows outside the band and aligning the diagonals.
//For example, the matrix
//1  2  3  0  0  0
//4  5  6  7  0  0
//0  8  9 10 11  0
//0  0 12 13 14 15
//0  0  0 16 17 18
//0  0  0  0 19 20

//	func NewDiagonalRect(r, c int, data []float64) *BandDense
//NewDiagonalRect is a convenience function that
//returns a diagonal matrix represented by a BandDense.
//The length of data must be min(r, c) otherwise NewDiagonalRect will panic.

//func (b *BandDense) At(i, j int) float64
// At returns the element at row i, column j.

//func (b *BandDense) Bandwidth() (kl, ku int)
//Bandwidth returns the upper and lower bandwidths of the matrix.

//func (b *BandDense) DiagView() Diagonal
//DiagView returns the diagonal as a matrix backed by the original data.

// func (b *BandDense) Dims() (r, c int)
//Dims returns the number of rows and columns in the matrix.

//func (b *BandDense) DoColNonZero(j int, fn func(i, j int, v float64))
//DoColNonZero calls the function fn for each of the non-zero elements of column j of b.
//The function fn takes a row/column index and the element value of b at (i, j).

//func (b *BandDense) DoNonZero(fn func(i, j int, v float64))
//DoNonZero calls the function fn for each of the non-zero elements of b.
//The function fn takes a row/column index and the element value of b at (i, j).

//func (b *BandDense) DoRowNonZero(i int, fn func(i, j int, v float64))
//DoRowNonZero calls the function fn for each of the non-zero elements of row i of b.
//The function fn takes a row/column index and the element value of b at (i, j).

//func (b *BandDense) IsEmpty() bool
//IsEmpty returns whether the receiver is empty. Empty matrices can be the receiver for size-restricted operations.
//The receiver can be zeroed using Reset

//func (b *BandDense) MulVecTo(dst *VecDense, trans bool, x Vector)
//MulVecTo computes B⋅x or Bᵀ⋅x storing the result into dst.

//func (b *BandDense) RawBand() blas64.Band
//RawBand returns the underlying blas64.Band used by the receiver.
//Changes to elements in the receiver following the call will be reflected in returned blas64.Band.

//func (b *BandDense) Reset()
//Reset empties the matrix so that it can be reused as the receiver of a dimensionally restricted operation.
//Reset should not be used when the matrix shares backing data. See the Reseter interface for more information.

//func (b *BandDense) SetBand(i, j int, v float64)
//SetBand sets the element at row i, column j to the value v.
//It panics if the location is outside the appropriate region of the matrix.

//func (b *BandDense) SetRawBand(mat blas64.Band)
//SetRawBand sets the underlying blas64.Band used by the receiver.
//Changes to elements in the receiver following the call will be reflected in the input.

//func (b *BandDense) T() Matrix
//T performs an implicit transpose by returning the receiver inside a Transpose

//func (b *BandDense) TBand() Banded
//TBand performs an implicit transpose by returning the receiver inside a TransposeBand.

//func (b *BandDense) Trace() float64
//Trace computes the trace of the matrix.

//func (b *BandDense) Zero()
//Zero sets all of the matrix elements to zero

//type BandWidther interface {
	//	BandWidth() (k1, k2 int)
	//}
	//A BandWidther represents a banded matrix and can return the left and right half-bandwidths, k1 and k2.

//	type Banded interface {
	//	Matrix
	//	// Bandwidth returns the lower and upper bandwidth values for
	//	// the matrix. The total bandwidth of the matrix is kl+ku+1.
	//	Bandwidth() (kl, ku int)
	//
	//	// TBand is the equivalent of the T() method in the Matrix
	//	// interface but guarantees the transpose is of banded type.
	//	TBand() Banded
	//Banded is a band matrix representation.

//	type CDense struct {
	//	// contains filtered or unexported fields
	//}
//	CDense is a dense matrix representation with complex data.

// func NewCDense(r, c int, data []complex128) *CDense
//NewCDense creates a new complex Dense matrix with r rows and c columns.
//If data == nil, a new slice is allocated for the backing slice.
//If len(data) == r*c, data is used as the backing slice,
//and changes to the elements of the returned CDense will be reflected in data.
//If neither of these is true, NewCDense will panic.
//NewCDense will panic if either r or c is zero.
// The data must be arranged in row-major order, i.e. the (i*c + j)-th element
//in the data slice is the {i, j}-th element in the matrix.

//func (m *CDense) At(i, j int) complex128
//At returns the element at row i, column j.

//func (m *CDense) Caps() (r, c int)
// Caps returns the number of rows and columns in the backing matrix.

//func (m *CDense) Conj(a CMatrix)
//Conj calculates the element-wise conjugate of a and stores the result in the receiver.
//Conj will panic if m and a do not have the same dimension unless m is empty.

//func (m *CDense) Copy(a CMatrix) (r, c int)
//Copy makes a copy of elements of a into the receiver.
//It is similar to the built-in copy; it copies as much as the overlap
//between the two matrices and returns the number of rows and columns it copied.
//If a aliases the receiver and is a transposed Dense or VecDense,
//with a non-unitary increment, Copy will panic.
//See the Copier interface for more information.

//func (m *CDense) Dims() (r, c int)
//Dims returns the number of rows and columns in the matrix.

//func (m *CDense) Grow(r, c int) CMatrix
//Grow returns the receiver expanded by r rows and c columns.
//If the dimensions of the expanded matrix
//are outside the capacities of the receiver a new allocation is made,
//otherwise not.
//Note the receiver itself is not modified during the call to Grow.

//func (m *CDense) H() CMatrix
//H performs an implicit conjugate transpose by returning the receiver inside a ConjTranspose.

//func (m *CDense) IsEmpty() bool
//IsEmpty returns whether the receiver is empty.
//Empty matrices can be the receiver for size-restricted operations.
//The receiver can be zeroed using Reset.

//func (m *CDense) RawCMatrix() cblas128.General
//RawCMatrix returns the underlying cblas128.General used by the receiver.
//Changes to elements in the receiver following the call will be reflected in
//returned cblas128.General.

//func (m *CDense) Reset()
//Reset zeros the dimensions of the matrix so that
//it can be reused as the receiver of a dimensionally restricted operation.
//Reset should not be used when the matrix shares backing data.
//See the Reseter interface for more information.

//func (m *CDense) ReuseAs(r, c int)
//ReuseAs changes the receiver if it IsEmpty() to be of size r×c.
//ReuseAs re-uses the backing data slice if it has sufficient capacity,
//otherwise a new slice is allocated.
//The backing data is zero on return.
//ReuseAs panics if the receiver is not empty,
//and panics if the input sizes are less than one.
//To empty the receiver for re-use, Reset should be used.

//func (m *CDense) Set(i, j int, v complex128)
//Set sets the element at row i, column j to the value v.

//func (m *CDense) SetRawCMatrix(b cblas128.General)
//SetRawCMatrix sets the underlying cblas128.General used by the receiver.
//Changes to elements in the receiver following the call will be reflected in b.

//func (m *CDense) Slice(i, k, j, l int) CMatrix
//Slice returns a new CMatrix that shares backing data with the receiver.
//The returned matrix starts at {i,j} of the receiver
//and extends k-i rows and l-j columns.
//The final row in the resulting matrix is k-1 and the final column is l-1.
//Slice panics with ErrIndexOutOfRange
//if the slice is outside the capacity of the receiver.

//func (m *CDense) T() CMatrix
//T performs an implicit transpose by returning the receiver inside a CTranspose.

//func (m *CDense) Zero()
//Zero sets all of the matrix elements to zero.

//type CMatrix interface {
	//	// Dims returns the dimensions of a CMatrix.
	//	Dims() (r, c int)
	//
	//	// At returns the value of a matrix element at row i, column j.
	//	// It will panic if i or j are out of bounds for the matrix.
	//	At(i, j int) complex128
	//
	//	// H returns the conjugate transpose of the CMatrix. Whether H
	//	// returns a copy of the underlying data is implementation dependent.
	//	// This method may be implemented using the ConjTranspose type, which
	//	// provides an implicit matrix conjugate transpose.
	//	H() CMatrix
	//
	//	// T returns the transpose of the CMatrix. Whether T returns a copy of the
	//	// underlying data is implementation dependent.
	//	// This method may be implemented using the CTranspose type, which
	//	// provides an implicit matrix transpose.
	//	T() CMatrix
	//}
//	CMatrix is the basic matrix interface type for complex matrices.

//type CTranspose struct {
	//	CMatrix CMatrix
	//}
//	CTranspose is a type for performing an implicit matrix conjugate transpose.
//	It implements the CMatrix interface,
//	returning values from the conjugate transpose of the matrix within.

//func (t CTranspose) At(i, j int) complex128
//At returns the value of the element
//at row i and column j of the conjugate transposed matrix,
//that is, row j and column i of the CMatrix field.

//func (t CTranspose) Dims() (r, c int)
//Dims returns the dimensions of the transposed matrix.
//The number of rows returned is the number of columns in the CMatrix field,
//and the number of columns is the number of rows in the CMatrix field.

//func (t CTranspose) H() CMatrix
//H performs an implicit transpose
//by returning the receiver inside a ConjTranspose.

//func (t CTranspose) T() CMatrix
//T performs an implicit conjugate transpose by returning the CMatrix field.

//func (t CTranspose) Untranspose() CMatrix
//Untranspose returns the CMatrix field.

//type CUntransposer interface {
	//	// Untranspose returns the underlying CMatrix stored for the implicit
	//	// transpose.
	//	Untranspose() CMatrix
	//}
//	CUntransposer is a type that can undo an implicit transpose.

//	type Cholesky struct {
//	// The chol pointer must never be retained as a pointer outside the Cholesky
//	// struct, either by returning chol outside the struct or by setting it to
//	// a pointer coming from outside. The same prohibition applies to the data
//	// slice within chol.
//	chol *TriDense
//	cond float64
//}
//Cholesky is a symmetric positive definite matrix
//represented by its Cholesky decomposition.
//The decomposition can be constructed using the Factorize method.
//The factorization itself can be extracted using the UTo or LTo methods,
//and the original symmetric matrix can be recovered with ToSym.

//Note that this matrix representation is useful for certain operations,
//in particular finding solutions to linear equations.
//It is very inefficient at other operations, in particular At is slow.

//Cholesky methods may only be called on a value
//that has been successfully initialized by a call to Factorize
//that has returned true.
//Calls to methods of an unsuccessful Cholesky factorization will panic.

// Construct a symmetric positive definite matrix.
	tmp := mat.NewDense(4,4,[]float64{
	2, 6, 8, -4,
		1, 8, 7, -2,
		2, 2, 1, 7,
		8, -2, -2, 1,
	})
var a3 mat.SymDense
a3.SymOuterK(1, tmp)
fmt.Printf("a = %0.4v\n", mat.Formatted(&a3, mat.Prefix("    ")))

// Compute the cholesky factorization.
	var chol mat.Cholesky

	if ok := chol.Factorize(&a3);!ok{
		fmt.Println("a matrix is not positive semi-definite.")
	}

	// Find the determinant.
	fmt.Printf("\nThe determinant of a is %0.4g\n\n", chol.Det())

	// Use the factorization to solve the system of equations a * x = b.
	b3 := mat.NewVecDense(4, []float64{1, 2, 3, 4})
	var x mat.VecDense
	if err := chol.SolveVecTo(&x,b3);err != nil{
		fmt.Println("Matrix is near singular: ", err)
	}

	fmt.Println("Solve a * x = b")
	fmt.Printf("x = %0.4v\n", mat.Formatted(&x, mat.Prefix("    ")))

	// Extract the factorization and check that it equals the original matrix.
	var t mat.TriDense
	chol.LTo(&t)
	var test mat.Dense
	test.Mul(&t,t.T())
	fmt.Println()
	fmt.Printf("L * Lᵀ = %0.4v\n", mat.Formatted(&test, mat.Prefix("         ")))

//func (c *Cholesky) At(i, j int) float64
//At returns the element at row i, column j.

//func (c *Cholesky) Clone(chol *Cholesky)
//Clone makes a copy of the input Cholesky into the receiver, overwriting the previous value of the receiver. Clone does not place any restrictions on receiver shape. Clone panics if the input Cholesky is not the result of a valid decomposition.

//func (c *Cholesky) Cond() float64
//Cond returns the condition number of the factorized matrix.
//Cond returns the condition number of the factorized matrix.

//func (c *Cholesky) Det() float64
//Det returns the determinant of the matrix that has been factorized.

//func (ch *Cholesky) Dims() (r, c int)
//Dims returns the dimensions of the matrix.

//func (c *Cholesky) ExtendVecSym(a *Cholesky, v Vector) (ok bool)
//ExtendVecSym computes the Cholesky decomposition of the original matrix A,
//whose Cholesky decomposition is in a,
//extended by a the n×1 vector v according to:
//[A  w]
//[w' k]
//where k = v[n-1] and w = v[:n-1].
//The result is stored into the receiver.
//In order for the updated matrix to be positive definite,
//it must be the case that k > w' A^-1 w.
//If this condition does not hold then ExtendVecSym will return false
//and the receiver will not be updated.
//ExtendVecSym will panic if v.Len() != a.Symmetric()+1 or
//if a does not contain a valid decomposition.

//func (c *Cholesky) Factorize(a Symmetric) (ok bool)
//Factorize calculates the Cholesky decomposition of the matrix A and returns
//whether the matrix is positive definite.
//If Factorize returns false, the factorization must not be used.

//func (c *Cholesky) InverseTo(dst *SymDense) error
//InverseTo computes the inverse of the matrix represented
//by its Cholesky factorization and stores the result into s.
//If the factorized matrix is ill-conditioned,
//a Condition error will be returned.
//Note that matrix inversion is numerically unstable,
//and should generally be avoided where possible,
//for example by using the Solve routines.

//func (c *Cholesky) IsEmpty() bool
//IsEmpty returns whether the receiver is empty.
//Empty matrices can be the receiver for size-restricted operations.
//The receiver can be emptied using Reset.

//func (c *Cholesky) LTo(dst *TriDense)
//LTo stores into dst the n×n lower triangular matrix L
//from a Cholesky decomposition:
//A = L * Lᵀ.
//If dst is empty, it is resized to be an n×n lower triangular matrix.
//When dst is non-empty, LTo panics if dst is not n×n or not Lower.
//LTo will also panic
//if the receiver does not contain a successful factorization.

//func (c *Cholesky) LogDet() float64
//LogDet returns the log of the determinant of the matrix that has been factorized.

//func (c *Cholesky) RawU() Triangular
//RawU returns the Triangular matrix used to store the Cholesky decomposition of the original matrix A.
//The returned matrix should not be modified.
//If it is modified, the decomposition is invalid and should not be used.

//func (c *Cholesky) Reset()
//Reset resets the factorization
//so that it can be reused as
//the receiver of a dimensionally restricted operation.

//func (c *Cholesky) Scale(f float64, orig *Cholesky)
//Scale multiplies the original matrix A
//by a positive constant using its Cholesky decomposition,
//storing the result in-place into the receiver.
//That is, if the original Cholesky factorization is
	//Uᵀ * U = A
	//the updated factorization is
	//U'ᵀ * U' = f A = A'
//Scale panics if the constant is non-positive, or if the receiver is non-empty and is of a different size from the input.

//func (c *Cholesky) SetFromU(t Triangular)
//SetFromU sets the Cholesky decomposition from the given triangular matrix.
//SetFromU panics if t is not upper triangular.
//If the receiver is empty it is resized to be n×n, the size of t.
//If dst is non-empty, SetFromU panics if c is not of size n×n.
//Note that t is copied into, not stored inside, the receiver.

//func (a *Cholesky) SolveCholTo(dst *Dense, b *Cholesky) error
//SolveCholTo finds the matrix X that solves A * X = B
//where A and B are represented by their Cholesky decompositions a and b.
//The result is stored in-place into dst.
//If the Cholesky decomposition is singular or near-singular
//a Condition error is returned.
//See the documentation for Condition for more information

//func (c *Cholesky) SolveTo(dst *Dense, b Matrix) error
//SolveTo finds the matrix X that solves A * X = B
//where A is represented by the Cholesky decomposition.
//The result is stored in-place into dst.
//If the Cholesky decomposition is singular or near-singular
//a Condition error is returned.
//See the documentation for Condition for more information.

//func (c *Cholesky) SolveVecTo(dst *VecDense, b Vector) error
//SolveVecTo finds the vector x that solves A * x = b
//where A is represented by the Cholesky decomposition.
//The result is stored in-place into dst.
//If the Cholesky decomposition is singular or near-singular
//a Condition error is returned.
//See the documentation for Condition for more information.

//func (c *Cholesky) SymRankOne(orig *Cholesky, alpha float64, x Vector) (ok bool)
//SymRankOne performs a rank-1 update of the original matrix A
//and refactorizes its Cholesky factorization, storing the result into the receiver.
//That is, if in the original Cholesky factorization
//raw: Uᵀ * U = A,
//in the updated factorization: U'ᵀ * U' = A + alpha * x * xᵀ = A'.
//Note that when alpha is negative, the updating problem may be ill-conditioned
//and the results may be inaccurate,
//or the updated matrix A' may not be positive definite
//and not have a Cholesky factorization.
//SymRankOne returns whether the updated matrix A' is positive definite.
//If the update fails the receiver is left unchanged.
//SymRankOne updates a Cholesky factorization in O(n²) time.
//The Cholesky factorization computation from scratch is O(n³).

a4 := mat.NewSymDense(4,[]float64{
	1, 1, 1, 1,
		0, 2, 3, 4,
		0, 0, 6, 10,
		0, 0, 0, 20,
})
fmt.Printf("A = %0.4v\n", mat.Formatted(a4, mat.Prefix("    ")))
// Compute the Cholesky factorization.
var chol2 mat.Cholesky
if ok := chol2.Factorize(a4);!ok{
	fmt.Println("matrix a is not positive definite.")
}

x3 := mat.NewVecDense(4,[]float64{0,0,0,1})
fmt.Printf("\nx = %0.4v\n", mat.Formatted(x3, mat.Prefix("    ")))

// Rank-1 update the factorization.
chol2.SymRankOne(&chol2,1,x3)

//// Rank-1 update the matrix a.
a4.SymRankOne(a4,1,x3)
var au mat.SymDense

chol2.ToSym(&au)
// Print the matrix that was updated directly.
	fmt.Printf("\nA' =        %0.4v\n", mat.Formatted(a4, mat.Prefix("            ")))
	// Print the matrix recovered from the factorization.
	fmt.Printf("\nU'ᵀ * U' =  %0.4v\n", mat.Formatted(&au, mat.Prefix("            ")))

//func (c *Cholesky) Symmetric() int
//Symmetric implements the Symmetric interface and returns the number of rows in the matrix (this is also the number of columns).

//func (c *Cholesky) T() Matrix
//T returns the receiver, the transpose of a symmetric matrix.

//func (c *Cholesky) ToSym(dst *SymDense)
//ToSym reconstructs the original positive definite matrix from its Cholesky decomposition,
//storing the result into dst. If dst is empty it is resized to be n×n.
//If dst is non-empty, ToSym panics if dst is not of size n×n.
//ToSym will also panic if the receiver does not contain a successful factorization.

//func (c *Cholesky) UTo(dst *TriDense)
//UTo stores into dst the n×n upper triangular matrix U from a Cholesky decomposition
//A = Uᵀ * U.
//If dst is empty, it is resized to be an n×n upper triangular matrix.
//When dst is non-empty, UTo panics if dst is not n×n or not Upper.
//UTo will also panic if the receiver does not contain a successful factorization.

//type ClonerFrom interface {
	//	CloneFrom(a Matrix)
	//}
//	A ClonerFrom can make a copy of a into the receiver,
//	overwriting the previous value of the receiver.
//	The clone operation does not make any restriction on shape
//	and will not cause shadowing.

//type ColNonZeroDoer interface {
	//	DoColNonZero(j int, fn func(i, j int, v float64))
	//}
//	A ColNonZeroDoer can call a function for each non-zero element of
//	a column of the receiver.
//	The parameters of the function are the element indices and its value.

//type ColViewer interface {
	//	ColView(j int) Vector
	//}
//	A ColViewer can return a Vector reflecting a column
//	that is backed by the matrix data.
//	The Vector returned will have length equal to the number of rows.

//type Condition float64
//Condition is the condition number of a matrix.
//The condition number is defined as |A| * |A^-1|
//One important use of Condition is during linear solve routines (finding x such that A * x = b).
//The condition number of A indicates the accuracy of the computed solution.
//A Condition error will be returned if the condition number of A is sufficiently large.
//If A is exactly singular to working precision,
//Condition == ∞, and the solve algorithm may have completed early.
//If Condition is large and finite the solve algorithm will be performed,
//but the computed solution may be inaccurate.
//Due to the nature of finite precision arithmetic,
//the value of Condition is only an approximate test of singularity.

//func (c Condition) Error() string

//type ConjTranspose struct {
	//	CMatrix CMatrix
	//}
//	ConjTranspose is a type for performing an implicit matrix conjugate transpose. It implements the CMatrix interface,
//	returning values from the conjugate transpose of the matrix within.

//func (t ConjTranspose) At(i, j int) complex128
//At returns the value of the element at row i and column j of the conjugate transposed matrix, that is, row j and column i of the CMatrix field.

//func (t ConjTranspose) Dims() (r, c int)
//Dims returns the dimensions of the transposed matrix. The number of rows returned is the number of columns in the CMatrix field, and the number of columns is the number of rows in the CMatrix field.

//func (t ConjTranspose) H() CMatrix
//H performs an implicit conjugate transpose by returning the CMatrix field.

//func (t ConjTranspose) T() CMatrix
//T performs an implicit transpose by returning the receiver inside a CTranspose.

//func (t ConjTranspose) UnConjTranspose() CMatrix
//UnConjTranspose returns the CMatrix field.

//type Copier interface {
	//	Copy(a Matrix) (r, c int)
	//}
//	A Copier can make a copy of elements of a into the receiver.
//	The submatrix copied starts at row and column 0
//	and has dimensions equal to the minimum dimensions of the two matrices.
//	The number of row and columns copied is returned.
//	Copy will copy from a source that aliases the receiver
//	unless the source is transposed;
//	an aliasing transpose copy will panic with the exception for a special case
//	when the source data has a unitary increment or stride.

//type Dense struct {
//	mat blas64.General
//
//	capRows, capCols int
//}
//Dense is a dense matrix representation.

//func DenseCopyOf(a Matrix) *Dense
//DenseCopyOf returns a newly allocated copy of the elements of a.

//func NewDense(r, c int, data []float64) *Dense
//NewDense creates a new Dense matrix with r rows and c columns.
//If data == nil, a new slice is allocated for the backing slice.
//If len(data) == r*c, data is used as the backing slice,
//and changes to the elements of the returned Dense will be reflected in data.
//If neither of these is true,
//NewDense will panic. NewDense will panic if either r or c is zero.
//The data must be arranged in row-major order, i.e. the (i*c + j)-th element in the data slice is the {i, j}-th element in the matrix.

//func (m *Dense) Add(a, b Matrix)
//Add adds a and b element-wise, placing the result in the receiver.
//Add will panic if the two matrices do not have the same shape.

a5 := mat.NewDense(2,2,[]float64{
	1, 0,
		1, 0,
})

b5 := mat.NewDense(2, 2, []float64{
		0, 1,
		0, 1,
	})

// Add a and b, placing the result into c.
	// Notice that the size is automatically adjusted
	// when the receiver is empty (has zero size).

	var c5 mat.Dense
c5.Add(a5,b5)
// Print the result using the formatter.
fc5 := mat.Formatted(&c5,mat.Prefix("    "),mat.Squeeze())

fmt.Printf("c = %v\n", fc5)

//func (m *Dense) Apply(fn func(i, j int, v float64) float64, a Matrix)
//Apply applies the function fn to each of the elements of a,
//placing the resulting matrix in the receiver.
//The function fn takes a row/column index
//and element value and returns some function of that tuple.

//func (m *Dense) At(i, j int) float64
//At returns the element at row i, column j.



//func (m *Dense) Augment(a, b Matrix)
//Augment creates the augmented matrix of a and b,
//where b is placed in the greater indexed columns.
//Augment will panic if the two input matrices
//do not have the same number of rows or the constructed augmented matrix
//is not the same shape as the receiver.

//func (m *Dense) Caps() (r, c int)
//Caps returns the number of rows and columns in the backing matrix.

//func (m *Dense) CloneFrom(a Matrix)
//CloneFrom makes a copy of a into the receiver,
//overwriting the previous value of the receiver.
//The clone from operation does not make any restriction on shape
//and will not cause shadowing.
//See the ClonerFrom interface for more information.

//func (m *Dense) ColView(j int) Vector
//ColView returns a Vector reflecting the column j, backed by the matrix data.
	//See ColViewer for more information.

//	func (m *Dense) Copy(a Matrix) (r, c int)
//Copy makes a copy of elements of a into the receiver.
//It is similar to the built-in copy;
//it copies as much as the overlap between the two matrices
//and returns the number of rows and columns it copied.
//If a aliases the receiver and is a transposed Dense or VecDense,
//with a non-unitary increment, Copy will panic.

//See the Copier interface for more information.

//func (m *Dense) DiagView() Diagonal
//DiagView returns the diagonal as a matrix backed by the original data.

//func (m *Dense) Dims() (r, c int)
//Dims returns the number of rows and columns in the matrix.

//func (m *Dense) DivElem(a, b Matrix)
//DivElem performs element-wise division of a by b, placing the result in the receiver.
//DivElem will panic if the two matrices do not have the same shape.

a6 := mat.NewDense(2, 2, []float64{
		5, 10,
		15, 20,
	})

b6 := mat.NewDense(2, 2, []float64{
		5, 5,
		5, 5,
	})

//// Divide the elements of a by b, placing the result into a.
a6.DivElem(a6,b6)
// Print the result using the formatter.
fa6 := mat.Formatted(a6, mat.Prefix("    "), mat.Squeeze())
fmt.Printf("a = %v\n", fa6)

//func (m *Dense) Exp(a Matrix)
//Exp calculates the exponential of the matrix a, e^a, placing the result in the receiver. Exp will panic with matrix.ErrShape if a is not square.

a7 := mat.NewDense(2, 2, []float64{
		1, 0,
		0, 1,
	})

// Take the exponential of the matrix and place the result in m.
var m7 mat.Dense

m7.Exp(a7)
matPrint(&m7)

//func (m *Dense) Grow(r, c int) Matrix
//Grow returns the receiver expanded by r rows and c columns.
//If the dimensions of the expanded matrix
//are outside the capacities of the receiver a new allocation is made,
//otherwise not.
//Note the receiver itself is not modified during the call to Grow.

//func (m *Dense) Inverse(a Matrix) error
//Inverse computes the inverse of the matrix a,
//storing the result into the receiver.
//If a is ill-conditioned, a Condition error will be returned.
//Note that matrix inversion is numerically unstable,
//and should generally be avoided where possible,
//for example by using the Solve routines.

// Initialize a matrix A.
a8 := mat.NewDense(2, 2, []float64{
		2, 1,
		6, 4,
	})

//// Compute the inverse of A.

var aInv8 mat.Dense

err8 := aInv8.Inverse(a8)
if err8 != nil {
		log.Fatalf("A is not invertible: %v", err8)
	}
	// Print the result using the formatter.
	fa8 := mat.Formatted(&aInv8, mat.Prefix("       "), mat.Squeeze())
	fmt.Printf("aInv = %.2g\n\n", fa8)

//	// Confirm that A * A^-1 = I.
var I mat.Dense
I.Mul(a8,&aInv8)
fi := mat.Formatted(&I, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("I = %v\n\n", fi)
//// The Inverse operation, however, should typically be avoided. If the
	//	// goal is to solve a linear system
	//	//  A * X = B,
	//	// then the inverse is not needed and computing the solution as
	//	// X = A^{-1} * B is slower and has worse stability properties than
	//	// solving the original problem. In this case, the SolveVec method of
	//	// VecDense (if B is a vector) or Solve method of Dense (if B is a
	//	// matrix) should be used instead of computing the Inverse of A.

	b8 := mat.NewDense(2, 2, []float64{
		2, 3,
		1, 2,
	})

	var x8 mat.Dense
	err8 = x8.Solve(a8,b8)
	if err8 != nil {
		log.Fatalf("no solution: %v", err8)
	}

	// Print the result using the formatter.
	fx8 := mat.Formatted(&x8, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("x = %.1f\n", fx8)

//	func (m *Dense) IsEmpty() bool
//IsEmpty returns whether the receiver is empty.
//Empty matrices can be the receiver for size-restricted operations.
//The receiver can be emptied using Reset.

//func (m *Dense) Kronecker(a, b Matrix)
//Kronecker calculates the Kronecker product of a and b, placing the result in the receiver.

//func (m Dense) MarshalBinary() ([]byte, error)
//MarshalBinary encodes the receiver into a binary form and returns the result.

//Dense is little-endian encoded as follows:
//0 -  3  Version = 1          (uint32)
	// 4       'G'                  (byte)
	// 5       'F'                  (byte)
	// 6       'A'                  (byte)
	// 7       0                    (byte)
	// 8 - 15  number of rows       (int64)
	//16 - 23  number of columns    (int64)
	//24 - 31  0                    (int64)
	//32 - 39  0                    (int64)
	//40 - ..  matrix data elements (float64)
	//         [0,0] [0,1] ... [0,ncols-1]
	//         [1,0] [1,1] ... [1,ncols-1]
	//         ...
	//         [nrows-1,0] ... [nrows-1,ncols-1]

//	func (m Dense) MarshalBinaryTo(w io.Writer) (int, error)
//MarshalBinaryTo encodes the receiver into a binary form and writes it into w.
//MarshalBinaryTo returns the number of bytes written into w
//and an error, if any.
//See MarshalBinary for the on-disk layout.

//func (m *Dense) Mul(a, b Matrix)
//Mul takes the matrix product of a and b, placing the result in the receiver. If the number of columns in a does not equal the number of rows in b, Mul will panic.

a9 := mat.NewDense(2, 2, []float64{
		4, 0,
		0, 4,
	})
	b9 := mat.NewDense(2, 3, []float64{
		4, 0, 0,
		0, 0, 4,
	})

	// Take the matrix product of a and b and place the result in c.
	var c9 mat.Dense

	c9.Mul(a9,b9)

	// Print the result using the formatter.
	fc9 := mat.Formatted(&c9, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("c = %v\n", fc9)

//	func (m *Dense) MulElem(a, b Matrix)
//MulElem performs element-wise multiplication of a and b,
//placing the result in the receiver.
//MulElem will panic if the two matrices do not have the same shape.

// Initialize two matrices, a and b.
	a10 := mat.NewDense(2, 2, []float64{
		1, 2,
		3, 4,
	})
	b10 := mat.NewDense(2, 2, []float64{
		1, 2,
		3, 4,
	})

//	// Multiply the elements of a and b, placing the result into a.
	a10.MulElem(a10,b10)
	// Print the result using the formatter.
	fa10 := mat.Formatted(a10, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("a = %v\n", fa10)

//	func (m *Dense) Outer(alpha float64, x, y Vector)
//Outer calculates the outer product of the vectors x and y, where x and y are treated as column vectors,
//and stores the result in the receiver.
//m = alpha * x * yᵀ
//In order to update an existing matrix, see RankOne.

//func (m *Dense) Permutation(r int, swaps []int)
//Permutation constructs an r×r permutation matrix with the given row swaps.
//A permutation matrix has exactly one element equal to one
//in each row and column and all other elements equal to zero.
//swaps[i] specifies the row with which i will be swapped,
//which is equivalent to the non-zero column of row i.

//func (m *Dense) Pow(a Matrix, n int)
//Pow calculates the integral power of the matrix a to n, placing the result in the receiver.
//Pow will panic if n is negative or if a is not square.

// Initialize a matrix with some data.
	a11 := mat.NewDense(2, 2, []float64{
		4, 4,
		4, 4,
	})

	// Take the second power of matrix a and place the result in m.
	var m11 mat.Dense
	m11.Pow(a11,2)

	// Print the result using the formatter.
	fm11 := mat.Formatted(&m11, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("m = %v\n\n", fm11)

//	// Take the zeroth power of matrix a and place the result in n.
	//	// We expect an identity matrix of the same size as matrix a.
	var n11 mat.Dense
	n11.Pow(a11,0)

	// Print the result using the formatter.
	fn11 := mat.Formatted(&n11, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("n = %v\n", fn11)

//	func (m *Dense) Product(factors ...Matrix)
//Product calculates the product of the given factors and places the result in the receiver.
//The order of multiplication operations is optimized to minimize the number of floating point operations on the basis that all matrix multiplications are general.

//func (m *Dense) RankOne(a Matrix, alpha float64, x, y Vector)
//RankOne performs a rank-one update to the matrix a with the vectors x and y, where x and y are treated as column vectors. The result is stored in the receiver.
//The Outer method can be used instead of RankOne if a is not needed.
//m = a + alpha * x * yᵀ

//func (m *Dense) RawMatrix() blas64.General
//RawMatrix returns the underlying blas64.General used by the receiver.
//Changes to elements in the receiver
//following the call will be reflected in returned blas64.General.

//func (m *Dense) RawRowView(i int) []float64
//RawRowView returns a slice backed by the same array as backing the receiver.

//func (m *Dense) Reset()
//Reset empties the matrix so that it can be reused
//as the receiver of a dimensionally restricted operation.
//Reset should not be used when the matrix shares backing data.
//See the Reseter interface for more information.

//func (m *Dense) ReuseAs(r, c int)
//ReuseAs changes the receiver if it IsEmpty() to be of size r×c.
	//ReuseAs re-uses the backing data slice if it has sufficient capacity,
	//otherwise a new slice is allocated.
	//The backing data is zero on return.
	//ReuseAs panics if the receiver is not empty,
	//and panics if the input sizes are less than one.
	//To empty the receiver for re-use, Reset should be used.

//	func (m *Dense) RowView(i int) Vector
//RowView returns row i of the matrix data represented as a column vector,
//backed by the matrix data.
	//See RowViewer for more information.

//	func (m *Dense) Scale(f float64, a Matrix)
//Scale multiplies the elements of a by f, placing the result in the receiver.
	//See the Scaler interface for more information.

	// Initialize a matrix with some data.
	a12 := mat.NewDense(2, 2, []float64{
		4, 4,
		4, 4,
	})
	// Scale the matrix by a factor of 0.25 and place the result in m.
	var m12 mat.Dense
	m12.Scale(0.25,a12)

	// Print the result using the formatter.
	fm12 := mat.Formatted(&m12, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("m = %4.3f", fm12)

//	func (m *Dense) Set(i, j int, v float64)
//Set sets the element at row i, column j to the value v.

//func (m *Dense) SetCol(j int, src []float64)
//SetCol sets the values in the specified column of the matrix to the values in src.
//len(src) must equal the number of rows in the receiver.

//func (m *Dense) SetRawMatrix(b blas64.General)
//SetRawMatrix sets the underlying blas64.General used by the receiver.
//Changes to elements in the receiver
//following the call will be reflected in b.

//func (m *Dense) SetRow(i int, src []float64)
//SetRow sets the values in the specified rows of the matrix to the values in src.
//len(src) must equal the number of columns in the receiver.

//func (m *Dense) Slice(i, k, j, l int) Matrix
//Slice returns a new Matrix that shares backing data with the receiver.
//The returned matrix starts at {i,j} of the receiver and extends k-i rows and l-j columns.
//The final row in the resulting matrix is k-1 and the final column is l-1.
//Slice panics with ErrIndexOutOfRange
//if the slice is outside the capacity of the receiver.

//func (m *Dense) Solve(a, b Matrix) error
//Solve solves the linear least squares problem
//minimize over x |b - A*x|_2
//where A is an m×n matrix A, b is a given m element vector and x is n element solution vector.
//Solve assumes that A has full rank, that is
//rank(A) = min(m,n)
//If m >= n, Solve finds the unique least squares solution of an overdetermined system.
//If m < n, there is an infinite number of solutions that satisfy b-A*x=0. In this case Solve finds the unique solution of an underdetermined system that minimizes |x|_2.
//Several right-hand side vectors b and solution vectors x can be handled in a single call. Vectors b are stored in the columns of the m×k matrix B.
//Vectors x will be stored in-place into the n×k receiver.
//If A does not have full rank, a Condition error is returned.
//See the documentation for Condition for more information.

//func (m *Dense) Stack(a, b Matrix)
//Stack appends the rows of b onto the rows of a, placing the result into the receiver
//with b placed in the greater indexed rows.
//Stack will panic if the two input matrices do not have the same number of columns
//or the constructed stacked matrix is not the same shape as the receiver.

//func (m *Dense) Sub(a, b Matrix)
//Sub subtracts the matrix b from a, placing the result in the receiver. Sub will panic if the two matrices do not have the same shape.

// Initialize two matrices, a and b.
	a13 := mat.NewDense(2, 2, []float64{
		1, 1,
		1, 1,
	})

	b13 := mat.NewDense(2, 2, []float64{
		1, 0,
		0, 1,
	})

//	// Subtract b from a, placing the result into a.
a13.Sub(a13,b13)
// Print the result using the formatter.
	fa13 := mat.Formatted(a13, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("a = %v\n", fa13)

//	func (m *Dense) T() Matrix
//T performs an implicit transpose by returning the receiver inside a Transpose.

//func (m *Dense) Trace() float64
//Trace returns the trace of the matrix.
//The matrix must be square or Trace will panic.

//func (m *Dense) UnmarshalBinary(data []byte) error
//UnmarshalBinary decodes the binary form into the receiver. It panics if the receiver is a non-empty Dense matrix.
	//See MarshalBinary for the on-disk layout.
//	Limited checks on the validity of the binary input are performed:
//- matrix.ErrShape is returned if the number of rows or columns is negative,
	//- an error is returned if the resulting Dense matrix is too
	//big for the current architecture (e.g. a 16GB matrix written by a
	//64b application and read back from a 32b application.)
//	UnmarshalBinary does not limit the size of the unmarshaled matrix, and so it should not be used on untrusted data.

//func UnmarshalBinaryFrom



}

func (v row) Dims()(r,c int){
	return 1,len(v)
}

func (v row) At(_,j int) float64{
	return v[j]
}

func (v row)T() mat.Matrix  {
	return column(v)
}

// RawVector allows fast path computation with the vector.
func (v row)RawVector() blas64.Vector {
	return blas64.Vector{N:len(v),Data: v,Inc: 1}
}

func (v column) Dims()(r,c int){
	return len(v),1
}

func (v column) At(i,_ int) float64{
	return v[i]
}

func (v column) T() mat.Matrix {
	return row(v)
}

// RawVector allows fast path computation with the vector
func (v column) RawVector() blas64.Vector{
	return blas64.Vector{N: len(v),Data: v,Inc: 1}
}
