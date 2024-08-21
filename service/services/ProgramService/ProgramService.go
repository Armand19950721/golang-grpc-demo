package ProgramService

import (
	"service/model"
	"service/protos/Common"
	"service/protos/Program"

	"context"
	"service/repositories"
	"service/utils"
	programUtils "service/utils/program"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateProgram(ctx context.Context, req *Program.CreateProgramRequest) (*Program.CreateProgramReply, error) {
	reply := &Program.CreateProgramReply{}

	// check param
	if !utils.ValidString(req.Model.Name, 1, 100) ||
		!utils.ValidNumber(int(req.Model.Seats), 1, 10) {
		return reply, status.Errorf(codes.InvalidArgument, "param err")
	}

	// insert
	program := model.Program{
		Name:                req.Model.Name,
		State:               req.Model.State,
		Seats:               req.Model.Seats,
		Categories:          utils.ArrayToJson(req.Model.Categories),
		Types:               utils.ArrayToJson(req.Model.Types),
		Templates:           utils.ArrayToJson(req.Model.Templates),
		EffectTools:         utils.ArrayToJson(req.Model.EffectTools),
		ArInteractModules:   utils.ArrayToJson(req.Model.ArInteractModules),
		ArEditWindowModules: utils.ArrayToJson(req.Model.ArEditWindowModules),
	}

	res, program := repositories.CreateProgram(program)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.ProgramId = program.Id.String()

	return reply, status.Errorf(codes.OK, "")
}

func GetProgramList(ctx context.Context, req *Program.GetProgramListRequest) (*Program.GetProgramListReply, error) {
	reply := &Program.GetProgramListReply{}

	// query
	program := &model.Program{
		State: Common.ProgramState_ON_SALE,
	}

	result, programs, _ := repositories.GetProgramList(*program)

	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, result.Error.Error())
	}

	reply.Models = programUtils.ConvertModelTablesToProtos(programs)
	return reply, status.Errorf(codes.OK, "")
}

func GetFreeProgram(ctx context.Context, req *Program.GetFreeProgramRequest) (*Program.GetFreeProgramReply, error) {
	reply := &Program.GetFreeProgramReply{}

	res, program := repositories.QueryProgram(model.Program{
		Name: "program_free",
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.ProgramId = program.Id.String()

	return reply, status.Errorf(codes.OK, "")
}
