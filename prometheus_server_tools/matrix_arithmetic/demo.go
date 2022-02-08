package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func matPrint(X mat.Matrix){
	fa := mat.Formatted(X,mat.Prefix(""),mat.Squeeze())
	fmt.Printf("%v\n",fa)
}

func main() {
	v := make([]float64,12)

	for i := 0;i<12;i++{
		v[i] = float64(i)
	}

	A := mat.NewDense(3,4,v)
	println("A:")
	matPrint(A)

	a := A.At(0,2)
	println("A[0,2]:",a)

	A.Set(0,2,-1.5)
	matPrint(A)

	fmt.Println("Row 1 of A:")
	matPrint(A.RowView(1)) // vector -> column vector

	fmt.Println("Column 0 of A:")
	matPrint(A.ColView(0))

	row := []float64{10,9,8,7}
	A.SetRow(0,row)
	matPrint(A)

	col := []float64{3,2,1}
	A.SetCol(0,col)
	matPrint(A)

	B := mat.NewDense(3,4,nil)
	B.Add(A,A)
	fmt.Println("B:")
	matPrint(B)

	C:= mat.NewDense(3,4,nil)
	C.Sub(A,B)
	fmt.Println("A-B:")
	matPrint(C)

	// We can scale all elements of the matrix by a constant
	C.Scale(-3.5,C)
	fmt.Println("-3.5*C")
	matPrint(C)

	//Transposing a matrix is a little funky.
	fmt.Println("A'")
	matPrint(A.T())

	// Multiplication is pretty straightforward
	D := mat.NewDense(3,3,nil)
	D.Product(A,B.T())
	println("A*B'")
	matPrint(D)

	//We can use Product to multiply as many matrices as we want
	//provided the receiver has the appropriate dimensions
	// The order of operations is optimized to reduce operations
	D.Product(D,A,B.T(),D)
	println("D*A*B'*D")
	matPrint(D)

	//We can also apply a function to elements of the matrix
	//This function must take two integers and a float64,
	//representing the row and column indices and the value in the
	//input matrix. It must return a float. See sumOfIndices below
	E := mat.NewDense(3,4,nil)
	E.Apply(sumOfIndices,A)
	println("E:")
	matPrint(E)

	//Once again, we have some functions that return scalar values
	//For example, we can compute the determinant
	E2 := A.Slice(0,3,0,3)
	d := mat.Det(E2)
	println("det(E)=",d)
	fmt.Println(A.Dims())
	var aa,bb int
	aa,bb = A.Caps()
	fmt.Println(aa,bb)
	A.Dims()
	//row, _ := A.Dims()
	fmt.Println(A.Caps())

	//compute the trace:
	t := mat.Trace(D)
	fmt.Println("tr(D)=",t)


}

func sumOfIndices(i,j int, v float64) float64{
	return float64(i+j)
}
