package NameNode

import (

	"context"
	"fmt"
	"net"
	"log"
	"math/rand"
	"os"
	"time"
	"bufio"
	"strings"

	"google.golang.org/grpc"
	"Tarea/Msgpb"
)

var path = "NameNode/DATA.txt"

var (
	cc_Cremator,err1 = grpc.Dial("dist066:50062",grpc.WithInsecure())
	c_Cremator = msg.NewGuardarDatoClient(cc_Cremator)
	cc_Grunt,err2 = grpc.Dial("dist067:50063",grpc.WithInsecure())
	c_Grunt = msg.NewGuardarDatoClient(cc_Grunt)
	cc_Synth,err3 = grpc.Dial("dist065:50055",grpc.WithInsecure())
	c_Synth = msg.NewGuardarDatoClient(cc_Synth)
)

type server struct{
	msg.UnimplementedGuardarDatoServer
}

func ID_rep(id string)(ok string){
	
	if _,err := os.Stat(path); os.IsNotExist(err){
		return "0"
	}
	
	file,err := os.Open(path)

	if err != nil{
		log.Fatalf("Error: %v",err)
	}

	scanner := bufio.NewScanner(file) //Devuelve un scanner del archivo
	scanner.Split(bufio.ScanLines)

	for scanner.Scan(){
		text := string(scanner.Text())
		text_split := strings.Split(text,":")

		if id == text_split[1]{
			return "1"
		}
	}
	return "0"
}

func (*server) Guardar(ctx context.Context, req *msg.EnvioNombre) (*msg.ConfSave,error){
	fmt.Println("Desde el server de NameNode, recibiendo Request de Combine")

	ok := ID_rep(req.GetMsgId())

	res := &msg.ConfSave{
		Confirmacion: "1",
	}

	if ok != "1"{
		res = &msg.ConfSave{
			Confirmacion: "0",
		}
		msg_type := req.GetMsgType()
		msg_id := req.GetMsgId()
		msg_msg := req.GetMsgMsg()
		
		var tipo string

		switch msg_type{
		case "1":
			tipo = "MILITAR"
			break
		case "2":
			tipo = "FINANCIERA"
			break
		case "3":
			tipo = "LOGISTICA"
			break
		default:
			break
		}

		req_to_Name := &msg.EnvioNombre{
			MsgType: tipo,
			MsgId: msg_id,
			MsgMsg: msg_msg,
		}
		destino := ""
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		aleatorio := r.Intn(3)

		switch aleatorio {
		case 0:
			_,err := c_Cremator.Guardar(context.Background(),req_to_Name)
			if err != nil{
				log.Fatalf("Error: %v",err)
			}
			destino = "Cremator"
			break
		case 1:
			_,err := c_Grunt.Guardar(context.Background(),req_to_Name)
			if err != nil{
				log.Fatalf("Error: %v",err)
			}
			destino = "Grunt"
			break
		case 2:
			_,err := c_Synth.Guardar(context.Background(),req_to_Name)
			if err != nil{
				log.Fatalf("Error: %v",err)
			}
			destino = "Synth"
			break
		default:
			break
		}

		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
		    panic(err)
		}
		defer f.Close()
		informacion  := tipo + ":" + msg_id + ":" + destino + "\n"
		if _, err = f.WriteString(informacion); err != nil {
		    panic(err)
		}
		fmt.Println("Guardando información en Namenode ",destino)
	}

	return res,nil
}

func (*server) ObtenerInfoName(req *msg.Peticion,stream msg.GuardarDato_ObtenerInfoNameServer) error{
	fmt.Println("Recibiendo Request desde Rebeldes en NameNode")

	tipo_info := req.GetPet()

	file,err := os.Open(path)
	if err != nil{
		log.Fatalf("error: %v",err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text string
	var text_split []string
	for scanner.Scan(){
		text = string(scanner.Text())
		text_split = strings.Split(text,":")
		if text_split[0] == tipo_info{
			req := &msg.Peticion{
				Pet: text_split[1],
			}
			var res *msg.EnvioNombre
			switch text_split[2]{
			case "Cremator":
				res,err = c_Cremator.ObtenerInfoData(context.Background(),req)
				break
			case "Grunt":
				res,err = c_Grunt.ObtenerInfoData(context.Background(),req)
				break
			case "Synth":
				res,err = c_Synth.ObtenerInfoData(context.Background(),req)
				break
			default:
				break
			}
			reStream := &msg.EnvioNombre{
				MsgType: res.GetMsgType(),
				MsgId: res.GetMsgId(),
				MsgMsg: res.GetMsgMsg(),
			}
			stream.Send(reStream)
		}
	}
	return nil
}

func (*server) Cierre(ctx context.Context,req *msg.Peticion) (*msg.Peticion, error){
	fmt.Println("Cerrando NameNode")

	res_s ,_ := c_Synth.Cierre(context.Background(),req)
	 _,_ = c_Cremator.Cierre(context.Background(),req)
	_ ,_ = c_Grunt.Cierre(context.Background(),req)

	os.Exit(0)

	return res_s,nil
}



func NameNode(){

	fmt.Println("NameNode Funcionando")

	lis,err := net.Listen("tcp",":50061")

	if err != nil{
		log.Fatalf("Error listener: %v",err)
	}

	s := grpc.NewServer()

	msg.RegisterGuardarDatoServer(s,&server{})

	if err := s.Serve(lis); err != nil{
		log.Fatalf("Error server: %v",err)
	}

	defer cc_Cremator.Close()
	defer cc_Grunt.Close()
	defer cc_Synth.Close()
}
