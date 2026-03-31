package dto

import (
	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen/v1"
	m "github.com/erdedan1/shared/mapper"
)

type ViewMarketsRequest struct {
	UserRole string
}

func (v *ViewMarketsRequest) FromProto(request *pb.ViewMarketsRequest) *ViewMarketsRequest {
	v.UserRole = request.UserRole.String()
	return v
}

func (v *ViewMarketsRequest) ToProto() *pb.ViewMarketsRequest {
	return &pb.ViewMarketsRequest{
		UserRole: m.UserRoleToProtoSpot(v.UserRole),
	}
}
