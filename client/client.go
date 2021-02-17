package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	pb "oss.navercorp.com/metis/metis-server/api"
)

const (
	rpcAddr = "localhost:10118"
	timeout = 10 * time.Second
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MetisClient
}

func New() (*Client, error) {
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewMetisClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateModel(ctx context.Context, modelName string) (*pb.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.CreateModel(ctx, &pb.CreateModelRequest{
		ModelName: modelName,
	})
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

func (c *Client) CreateProject(ctx context.Context, name string) (*pb.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.CreateProject(ctx, &pb.CreateProjectRequest{
		ProjectName: name,
	})
	if err != nil {
		return nil, err
	}

	return res.Project, nil
}

func (c *Client) ListProjects(ctx context.Context) ([]*pb.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := c.client.ListProjects(ctx, &pb.ListProjectsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Projects, nil
}
