package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    pb "grpc-ecommerce/pb"
)

var (
    client pb.OrderServiceClient
    token  string
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Erro ao conectar: %v", err)
    }
    defer conn.Close()

    client = pb.NewOrderServiceClient(conn)
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("1. Login")
    fmt.Println("2. Registrar")
    fmt.Print("Escolha: ")
    opcao, _ := reader.ReadString('\n')
    opcao = strings.TrimSpace(opcao)

    switch opcao {
    case "2":
        fmt.Print("Novo usuário: ")
        username, _ := reader.ReadString('\n')
        fmt.Print("Nova senha: ")
        password, _ := reader.ReadString('\n')

        username = strings.TrimSpace(username)
        password = strings.TrimSpace(password)

        res, err := client.Register(context.Background(), &pb.RegisterRequest{
            Username: username,
            Password: password,
        })
        if err != nil {
            log.Fatalf("Erro ao registrar: %v", err)
        }

        fmt.Println(res.Message)
        return

    case "1":
        fmt.Print("Usuário: ")
        username, _ := reader.ReadString('\n')
        fmt.Print("Senha: ")
        password, _ := reader.ReadString('\n')

        username = strings.TrimSpace(username)
        password = strings.TrimSpace(password)

        res, err := client.Login(context.Background(), &pb.LoginRequest{
            Username: username,
            Password: password,
        })
        if err != nil {
            log.Fatalf("Login falhou: %v", err)
        }
        token = res.Token
        fmt.Println("Login bem-sucedido!")

    default:
        fmt.Println("Opção inválida.")
        return
    }

    // Menu principal
    for {
        fmt.Println("\n=== MENU PEDIDOS ===")
        fmt.Println("1. Criar pedido")
        fmt.Println("2. Listar pedidos")
        fmt.Println("3. Buscar pedido por ID")
        fmt.Println("4. Atualizar status do pedido")
        fmt.Println("5. Deletar pedido")
        fmt.Println("0. Sair")
        fmt.Print("Escolha: ")
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            createOrder(reader)
        case "2":
            listOrders()
        case "3":
            getOrder(reader)
        case "4":
            updateOrder(reader)
        case "5":
            deleteOrder(reader)
        case "0":
            fmt.Println("Encerrando.")
            return
        default:
            fmt.Println("Opção inválida.")
        }
    }
}

func authCtx() context.Context {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    go func() {
        <-ctx.Done()
        cancel()
    }()
    return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func createOrder(reader *bufio.Reader) {
    fmt.Print("Nome do cliente: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    res, err := client.CreateOrder(authCtx(), &pb.CreateOrderRequest{CustomerName: name})
    if err != nil {
        log.Println("Erro:", err)
        return
    }
    fmt.Printf("Pedido criado! ID: %s | Cliente: %s | Status: %s\n", res.Id, res.CustomerName, res.Status)
}

func listOrders() {
    stream, err := client.ListOrders(authCtx(), &pb.Empty{})
    if err != nil {
        log.Println("Erro ao listar pedidos:", err)
        return
    }
    fmt.Println("=== Pedidos ===")
    for {
        order, err := stream.Recv()
        if err != nil {
            break
        }
        fmt.Printf("- [%s] %s (%s)\n", order.Id, order.CustomerName, order.Status)
    }
}

func getOrder(reader *bufio.Reader) {
    fmt.Print("ID do pedido: ")
    id, _ := reader.ReadString('\n')
    id = strings.TrimSpace(id)

    res, err := client.GetOrder(authCtx(), &pb.GetOrderRequest{Id: id})
    if err != nil {
        log.Println("Erro:", err)
        return
    }
    fmt.Printf("Pedido: %s | Cliente: %s | Status: %s\n", res.Id, res.CustomerName, res.Status)
}

func updateOrder(reader *bufio.Reader) {
    fmt.Print("ID do pedido: ")
    id, _ := reader.ReadString('\n')
    id = strings.TrimSpace(id)

    fmt.Print("Novo status: ")
    status, _ := reader.ReadString('\n')
    status = strings.TrimSpace(status)

    res, err := client.UpdateOrder(authCtx(), &pb.UpdateOrderRequest{Id: id, Status: status})
    if err != nil {
        log.Println("Erro:", err)
        return
    }
    fmt.Printf("Pedido atualizado: %s | Status atual: %s\n", res.Id, res.Status)
}

func deleteOrder(reader *bufio.Reader) {
    fmt.Print("ID do pedido: ")
    id, _ := reader.ReadString('\n')
    id = strings.TrimSpace(id)

    res, err := client.DeleteOrder(authCtx(), &pb.DeleteOrderRequest{Id: id})
    if err != nil {
        log.Println("Erro:", err)
        return
    }
    fmt.Println(res.Message)
}
