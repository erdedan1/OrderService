package dto

import (
	"fmt"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen/v2"
)

type ViewMarketsRequest struct {
	UserRoles []string
}

func (v *ViewMarketsRequest) FromProto(request *pb.ViewMarketsRequest) *ViewMarketsRequest {
	fmt.Println(request.UserRoles)
	v.UserRoles = request.GetUserRoles()
	return v
}

func (v *ViewMarketsRequest) ToProto() *pb.ViewMarketsRequest {
	fmt.Println(v.UserRoles)
	return &pb.ViewMarketsRequest{
		UserRoles: v.UserRoles,
	}
}
