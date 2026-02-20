package dto

import (
	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
)

type ViewMarketsRequest struct {
	UserRoles []string
}

func (vmr *ViewMarketsRequest) DtoToProto() *pb.ViewMarketsRequest {
	pbVmr := pb.ViewMarketsRequest{}
	for _, u := range vmr.UserRoles {
		pbVmr.UserRoles = append(pbVmr.UserRoles, u)
	}
	return &pbVmr
}

func (vmr *ViewMarketsRequest) UserRolesToProto(roles []string) *ViewMarketsRequest {
	for _, r := range roles {
		vmr.UserRoles = append(vmr.UserRoles, r)
	}
	return vmr
}
