package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"justforfunc/31grpc/todo"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type taskServer struct {
	todo.UnimplementedTasksServer
}

func main() {
	server := grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(server, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("could not listen to :8888, %v", err)
	}
	log.Fatal(server.Serve(l))
}

type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)

var endianness = binary.LittleEndian

func (taskServer) List(ctx context.Context, _ *todo.Void) (*todo.TaskList, error) {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var tasks todo.TaskList
	for {
		if len(b) == 0 {
			return &tasks, nil
		} else if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]
		tasks.Tasks = append(tasks.Tasks, &task)
	}
}

func (taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {

	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("Could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("Could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return &todo.Task{}, fmt.Errorf("Could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("Could not write to db %s: %v", dbPath, err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("Could not close %s: %v", dbPath, err)
	}

	return task, nil
}
