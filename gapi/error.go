package gapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
