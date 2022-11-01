package Synth

import (
	
	"context"
	"log"
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"

	"google.golang.org/grpc"
	"Tarea/Msgpb"

)

var path = "DataNode/Synth/DATA.txt"

type server struct{
	msg.UnimplementedGuardarDatoServer
}

func (*server) Guardar(ctx context.Context, req *msg.EnvioNombre) (*msg.ConfSave,error){
	fmt.Println(fmt.Println("DataNode Synth recibe mensaje desde NameNode\nMensaje: ",req))
	fmt.Println("\nSobre que información quieres consultar: \n1-Militar\n2-Logisitca\n3-Financiera\n4-Salir")

	msg_type := req.GetMsgType()
	msg_id := req.GetMsgId()
	msg_msg := req.GetMsgMsg()

	informacion := msg_type + ":" + msg_id + ":" + msg_msg + "\n"

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
	    panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(informacion); err != nil {
	    panic(err)
	}


	

	res := &msg.ConfSave{
		Confirmacion: "1",
	}

	return res,nil
}

func (*server) ObtenerInfoData(ctx context.Context,req *msg.Peticion) (*msg.EnvioNombre, error) {
	fmt.Println("DataNode Synth: Leyendo peticion ID desde NameNode")
	fmt.Println("\nSobre que información quieres consultar: \n1-Militar\n2-Logisitca\n3-Financiera\n4-Salir")

	ID := req.GetPet()

	file,err := os.Open(path)

	if err != nil{
		log.Fatalf("Error: %v",err)
	}

	scanner := bufio.NewScanner(file) //Devuelve un scanner del archivo
	scanner.Split(bufio.ScanLines)

	var text string
	var text_split []string

	for scanner.Scan() {
		text = string(scanner.Text())
		text_split = strings.Split(text,":")
		if text_split[1] == ID{
			break
		}
	}
	res := &msg.EnvioNombre{
		MsgType: text_split[0],
		MsgId: text_split[1],
		MsgMsg: text_split[2],
	}
	return res,nil
}

func (*server) Cierre(ctx context.Context,req *msg.Peticion) (*msg.Peticion, error){
	fmt.Println("Cerrando DataNode Synth")

	res := &msg.Peticion{
		Pet: "OK",
	}

	os.Exit(0)

	return res,nil
}

func Synth(){

	fmt.Println("DataNode funcionando")

	lis,err := net.Listen("tcp",":50055")
	
	if err != nil{
		log.Fatalf("Error listener: %v",err)
	}

	s := grpc.NewServer()

	msg.RegisterGuardarDatoServer(s,&server{})

	if err := s.Serve(lis); err != nil{
		log.Fatalf("Error server: %v",err)
	}
}