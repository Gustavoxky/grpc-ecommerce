package server

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    pb "grpc-ecommerce/pb"

)

type orderService struct {
    pb.UnimplementedOrderServiceServer
    mu sync.Mutex
}

func (s *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    id := fmt.Sprintf("%d", len(req.CustomerName)+1)
    status := "Pendente"

    err := InsertOrder(id, req.CustomerName, status)
    if err != nil {
        return nil, err
    }

    return &pb.OrderResponse{Id: id, CustomerName: req.CustomerName, Status: status}, nil
}

func (s *orderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    id, customerName, status, err := GetOrderById(req.Id)
    if err != nil {
        return nil, fmt.Errorf("Pedido não encontrado")
    }

    return &pb.OrderResponse{Id: id, CustomerName: customerName, Status: status}, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
    err := UpdateOrderStatus(req.Id, req.Status)
    if err != nil {
        return nil, fmt.Errorf("Erro ao atualizar status: %v", err)
    }

    id, customerName, status, _ := GetOrderById(req.Id)
    return &pb.OrderResponse{Id: id, CustomerName: customerName, Status: status}, nil
}

func (s *orderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteResponse, error) {
    err := DeleteOrder(req.Id)
    if err != nil {
        return nil, fmt.Errorf("Erro ao deletar pedido: %v", err)
    }

    return &pb.DeleteResponse{Message: "Pedido deletado com sucesso"}, nil
}

func (s *orderService) ListOrders(_ *pb.Empty, stream pb.OrderService_ListOrdersServer) error {
    orders, err := ListOrders()
    if err != nil {
        return err
    }

    for _, o := range orders {
        stream.Send(&pb.OrderResponse{
            Id:            o["id"],
            CustomerName:  o["customer_name"],
            Status:        o["status"],
        })
    }
    return nil
}

func (s *orderService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    if !ValidateUser(req.Username, req.Password) {
        return nil, fmt.Errorf("usuário ou senha inválidos")
    }

    token, err := GenerateToken(req.Username)
    if err != nil {
        return nil, err
    }

    return &pb.LoginResponse{Token: token}, nil
}



func StartServer() {
    InitDB()
    defer CloseDB()

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Falha ao iniciar servidor: %v", err)
    }

    s := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            LoggingUnaryInterceptor(),
            JWTUnaryInterceptor(),
        ),
    )
    
    pb.RegisterOrderServiceServer(s, &orderService{})

    fmt.Println("Servidor gRPC rodando na porta 50051...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Erro ao rodar servidor gRPC: %v", err)
    }
}

func (s *orderService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    if req.Username == "" || req.Password == "" {
        return nil, fmt.Errorf("usuário e senha são obrigatórios")
    }

    err := RegisterUser(req.Username, req.Password)
    if err != nil {
        return nil, fmt.Errorf("erro ao registrar usuário: %v", err)
    }

    return &pb.RegisterResponse{Message: "Usuário registrado com sucesso"}, nil
}

