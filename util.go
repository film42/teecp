package main

import (
	"io"
)

type sinkReadWriter struct {
	io.ReadWriter
}

func (s *sinkReadWriter) Read(b []byte) (int, error) {
	return len(b), nil
}

func (s *sinkReadWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
