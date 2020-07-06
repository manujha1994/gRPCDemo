package main

import (
	"context"
	"fmt"
	grpcdemopb "gRPCDemo/gRPCProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func sendData(stream grpcdemopb.AdditionService_BiDirectionalStreamingCalculationClient) {
	numberSlice := make([] *grpcdemopb.BiDirectionalStreamCalculationRequest,0)
	for i:= int32(1); i<=25; i++ {
		numberSlice = append(numberSlice,&grpcdemopb.BiDirectionalStreamCalculationRequest{
			Number: i,
		})
	}
	for _, req := range numberSlice {
		stream.Send(req)
	}
	err := stream.CloseSend()
	if err != nil {
		showError(err, `sendData`)
	}
	wg.Done()
}

func receiveData(stream grpcdemopb.AdditionService_BiDirectionalStreamingCalculationClient) {
	for {
		response, err:= stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			showError(err, `receiveData`)
		}
		fmt.Println(`Receiving the sum of range as BiDirectional stream `, response.GetResult())
	}
	wg.Done()
}

func getBiDirectionalResult(c grpcdemopb.AdditionServiceClient) {
	fmt.Println(`**************************************`)
	fmt.Println(`Starting Bi Directional Stream Client`)
	fmt.Println(`**************************************`)
	stream, err := c.BiDirectionalStreamingCalculation(context.Background())
	if err!= nil {
		showError(err, `getBiDirectionalResult`)
	}
	wg.Add(2)
	go sendData(stream)
	go receiveData(stream)
	wg.Wait()
}

func getClientStreamResult(c grpcdemopb.AdditionServiceClient) {
	fmt.Println(`********************************`)
	fmt.Println(`Starting Client Stream Client`)
	fmt.Println(`********************************`)
	numberSlice := make([] *grpcdemopb.ClientStreamCalculationRequest,0)
	for i:= int32(1); i<=20; i++ {
		numberSlice = append(numberSlice, &grpcdemopb.ClientStreamCalculationRequest{Number: i})
	}
	stream, err := c.ClientStreamingCalculation(context.Background())
	if err!= nil {
		showError(err, `getClientStreamResult ClientStreamingCalculation`)
	}
	for _, req := range numberSlice {
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err!= nil {
		showError(err, `getClientStreamResult`)
	}
	fmt.Println(`Receiving the sum of range as client stream `, res.GetResult())
}

func getServerStreamResult(c grpcdemopb.AdditionServiceClient) {
	fmt.Println(`******************************`)
	fmt.Println(`Starting Server Stream Client`)
	fmt.Println(`******************************`)
	request := &grpcdemopb.ServerStreamCalculationRequest{
		Number: &grpcdemopb.Numbers{
			StartingValue: 1,
			EndingValue: 10,
		},
	}
	response, err:= c.ServerStreamingCalculation(context.Background(), request)
	if err!= nil {
		showError(err, `getServerStreamResult ServerStreamingCalculation`)
	}
	for {
		sum, err:= response.Recv()
		if err == io.EOF {
			fmt.Println(`Calculation complete for Server Streaming`)
			break
		}
		if err != nil {
			showError(err, `getServerStreamResult`)
		}
		fmt.Println(`Receiving the sum of range as server stream `, sum.GetResult())
	}
}

func getUnaryResult(c grpcdemopb.AdditionServiceClient, timeout time.Duration) {
	fmt.Println(`Starting Unary Client`)
	request := &grpcdemopb.CalculationUnaryRequest{
		Number: &grpcdemopb.Numbers{
			StartingValue: 5,
			EndingValue: 500,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := c.UnaryCalculation(ctx, request)
	if err!= nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		}
		return	}
	fmt.Println(`Sum received using Unary API is`, response.GetResult())
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
	fmt.Println(`***********************`)
	fmt.Println(`Starting gRPC client`)
	fmt.Println(`***********************`)
	conn, err := grpc.Dial(`localhost:50051`, grpc.WithInsecure())
	if err != nil {
		showError(err, `main`)
	}
	defer conn.Close()
	c := grpcdemopb.NewAdditionServiceClient(conn)
	getUnaryResult(c, 2 * time.Second)
	getServerStreamResult(c)
	getClientStreamResult(c)
	getBiDirectionalResult(c)
}

