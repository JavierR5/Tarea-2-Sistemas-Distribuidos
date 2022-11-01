package main

import (
	"Tarea/DataNode/Cremator"
	"Tarea/Combine"
)

func main(){
	go Cremator.Cremator()
	Combine.Combine()
}