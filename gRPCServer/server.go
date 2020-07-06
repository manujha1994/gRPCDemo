package main

import (
	"context"
	"fmt"
	grpcdemopb "gRPCDemo/gRPCProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"time"
)

type server struct {}

func (s *server) BiDirectionalStreamingCalculation(stream grpcdemopb.AdditionService_BiDirectionalStreamingCalculationServer) error {
	time.Sleep(1 * time.Second)
	fmt.Println(`Starting gRPC Bi Directional streaming Calculation`)
	var sum int32 = 0
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			showError(err, `BiDirectionalStreamingCalculation stream.Recv`)
			return err
		}
		sum += request.GetNumber()
		time.Sleep(1 * time.Second)
		err = stream.Send(&grpcdemopb.BiDirectionalStreamCalculationResponse{
			Result: sum,
		})
		if err != nil {
			showError(err, `BiDirectionalStreamingCalculation stream.Send`)
			return err
		}
	}
}

func (s *server) ClientStreamingCalculation(stream grpcdemopb.AdditionService_ClientStreamingCalculationServer) error {
	time.Sleep(1 * time.Second)
	fmt.Println(`Starting gRPC Client streaming Calculation`)
	var sum int32 = 0
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			showError(err, `ClientStreamingCalculation stream.Recv`)
			return err
		}
		sum += request.GetNumber()
	}
	err := stream.SendAndClose(&grpcdemopb.ClientStreamCalculationResponse{
		Result: sum,
	})
	if err != nil {
		showError(err, `ClientStreamingCalculation stream.SendAndClose`)
		return err
	}
	return nil
}

func (s *server) ServerStreamingCalculation(request *grpcdemopb.ServerStreamCalculationRequest, stream grpcdemopb.AdditionService_ServerStreamingCalculationServer) error {
	time.Sleep(1 * time.Second)
	fmt.Println(`Starting gRPC Server streaming Calculation`)
	var sum int32 = 0
	startingValue := request.GetNumber().StartingValue
	endingValue   := request.GetNumber().EndingValue
	for i:= startingValue; i<= endingValue; i++ {
		sum+=i
		response := &grpcdemopb.ServerStreamCalculationResponse{
			Result: sum,
		}
		err := stream.Send(response)
		if err != nil {
			showError(err, `ServerStreamingCalculation`)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) UnaryCalculation(ctx context.Context, request *grpcdemopb.CalculationUnaryRequest) (*grpcdemopb.CalculationUnaryResponse, error) {
	time.Sleep(1 * time.Second)
	fmt.Println(`Starting gRPC UnaryCalculation`)
	var sum int32 = 0
	startingValue := request.GetNumber().StartingValue
	endingValue   := request.GetNumber().EndingValue
	for i:= startingValue; i<= endingValue; i++ {
		sum += i
	}
	response := &grpcdemopb.CalculationUnaryResponse{
		Result: sum,
	}
	return response,nil
}

func showError (err error, methodName string) {
	respError, ok := status.FromError(err)
	if ok {
		fmt.Println(respError.Message())
		fmt.Println(respError.Code())
	} else {
		log.Fatalf(`Falling back due to error %v in method %v`, err, methodName)
	}
}

func main() {
	fmt.Println(`Staring gRPC Demo`)
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln(`Falling back due to error`, err)
	}
	s := grpc.NewServer()
	grpcdemopb.RegisterAdditionServiceServer(s, &server{})
	if err = s.Serve(listener); err != nil {
		showError(err, `main`)
	}
}
