package main

import (
	"Tarea/Rebeldes"
	"Tarea/DataNode/Synth"
)

func main(){
	go Synth.Synth()
	Rebeldes.Rebeldes()
}