package main

import (

)

/*
	MEAN DEFINITION
	****************************************
	The mean is also known as the average. To calculate the mean of
	a range of numbers, simply add the values in the set, then divide by the number of values.
	In this example, we add 14 twice, 11 five times, a three, four sevens, a four, and an eight.
	Divide that sum by 14, the total of numbers in the set, and in this case, the mean equals nine.
	*****************************************
*/
type MeanCalculation interface {
	Mean() (float32)
}


/*
	MEDIAN DEFINITION
	*****************************************
	The median is the number halfway between all the values in a sorted range of
	values. Think of the median as a median strip of a road. It always marks the center
	of the road. To calculate the median, first sort the numbers from lowest to highest.
	For an odd number of values, just take the middle number. For an even number of values,
	calculate the average of the two central numbers.
	******************************************
*/
type MedianCalculation interface {
	Median() (float32)
}


