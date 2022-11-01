package Combine

import (
	"fmt"
	"context"
	"log"

	"google.golang.org/grpc"
	"Tarea/Msgpb"
)


func Combine(){

	var info string
	var ID string
	var mensaje string

	fmt.Println("Bienvenidos a la consola de los Combine.")
	
	cc,err := grpc.Dial("dist068:50061",grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Error conexion: %v",err)
	}

	defer cc.Close()

	c := msg.NewGuardarDatoClient(cc)
	var tipo string
	for {
		fmt.Println("Ingrese la acción a realizar:\n1-Ingresar Información\n2-Salir")
		fmt.Scanln(&info)
		if info == "2"{
			break
		} else{
			fmt.Println("\nIngrese el tipo de información:\n1-Militar\n2-Financiera\n3-Logistica")
			fmt.Scanln(&tipo)
			fmt.Println("\nIngrese el ID de la información")
			fmt.Scanln(&ID)
			fmt.Println("\nIngrese el mensaje a guardar")
			fmt.Scanln(&mensaje)
			req := &msg.EnvioNombre{
				MsgType: tipo,
				MsgId: ID,
				MsgMsg: mensaje,
			}
			for{
				res, err:= c.Guardar(context.Background(),req)

				if err != nil{
					log.Fatalf("Error: %v",err)
				}
				if res.GetConfirmacion() == "1"{
					fmt.Println("Error al guardar la información. ID ya existe, asi que elige otro.")
					fmt.Println("Ingrese nuevo ID")
					fmt.Scanln(&ID)
					req = &msg.EnvioNombre{
						MsgType: tipo,
						MsgId: ID,
						MsgMsg: mensaje,
					}
				} else{
					break
				}
			}

		}

		fmt.Println("\nMensaje Guardado!!!\n\n")
	}

}
