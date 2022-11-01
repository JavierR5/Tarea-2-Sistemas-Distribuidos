package Rebeldes

import (
	"fmt"
	"io"
	"log"
	"context"

	"google.golang.org/grpc"
	"Tarea/Msgpb"
)

func Rebeldes(){
	fmt.Println("Bienvenido Rebelde")

	cc,err := grpc.Dial("dist068:50061",grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Error conexion: %v",err)
	}

	defer cc.Close()

	c := msg.NewGuardarDatoClient(cc)

	var peticion string
	dummy := 0

	for {
		fmt.Println("\nSobre que información quieres consultar: \n1-Militar\n2-Logisitca\n3-Financiera\n4-Salir")
		fmt.Scanln(&peticion)
		
		var tipo string

		switch peticion {
		case "1":
			tipo = "MILITAR"
			break
		case "2":
			tipo = "LOGISTICA"
			break
		case "3":
			tipo = "FINANCIERA"
			break
		case "4":
			dummy = 1
		default:
				break
		}
		if dummy == 1{
			req_cierre := &msg.Peticion{
				Pet: "OK",
			}
			_,_ = c.Cierre(context.Background(),req_cierre)
			break
		}

		req:= &msg.Peticion{
			Pet: tipo,
		}

		reStream,err := c.ObtenerInfoName(context.Background(),req)

		if err != nil{
			log.Fatalf("Error: %v",err)
		}
		for {
			msgs,err := reStream.Recv()
			if err == io.EOF{
				break
			}
			if err != nil{
				log.Fatalf("Error: %v",err)
			}
			mensaje := msgs.GetMsgType() + ":" + msgs.GetMsgId() + ":" + msgs.GetMsgMsg()
			fmt.Println(mensaje)
		}
	}
}
