package main

import (

)

type MeanCalculation interface {
	Mean() (float32)
}

type MedianCalculation interface {
	Median() (float32)
}


