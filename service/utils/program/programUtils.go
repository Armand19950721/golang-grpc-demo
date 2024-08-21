package programUtils

import (
	"service/model"
	"service/protos/ArContent"
	"service/protos/Common"
	"service/protos/Program"
	"service/repositories"
	"service/utils"

	"github.com/google/uuid"
)

func ConvertModelTablesToProtos(programs []model.Program) []*Program.ProgramModel {
	var models []*Program.ProgramModel
	for _, item := range programs {
		models = append(models, ConvertModelTableToProto(item))
	}

	return models
}

func ConvertModelTableToProto(program model.Program) *Program.ProgramModel {

	types, err := utils.ParseJsonWithType[[]ArContent.ArContentTypeEnum](program.Types)
	if err != nil {
		panic(err.Error())
	}

	effectTools, err := utils.ParseJsonWithType[[]Common.EffectTool](program.EffectTools)
	if err != nil {
		panic(err.Error())
	}

	arInteractModules, err := utils.ParseJsonWithType[[]Common.ArInteractModule](program.ArInteractModules)
	if err != nil {
		panic(err.Error())
	}

	arEditWindowModules, err := utils.ParseJsonWithType[[]Common.ArEditWindowModule](program.ArEditWindowModules)
	if err != nil {
		panic(err.Error())
	}

	model := &Program.ProgramModel{
		Id:                  program.Id.String(),
		Name:                program.Name,
		State:               program.State,
		Seats:               program.Seats,
		Types:               types,
		EffectTools:         effectTools,
		ArInteractModules:   arInteractModules,
		ArEditWindowModules: arEditWindowModules,
	}

	return model
}

func ProgramTableJsonToProtoModel(json string) (result *Program.ProgramModel, success bool) {
	if !utils.ValidJson(json) {
		return &Program.ProgramModel{}, false
	}

	tableModel, err := utils.ParseJsonWithType[model.Program](json)
	if err != nil {
		return &Program.ProgramModel{}, false
	}

	return ConvertModelTableToProto(tableModel), true
}

func GetProgramDataJson(user *model.User) (string, Common.ErrorCodes) {
	defaultVal := ""
	var programId uuid.UUID

	// check admin or child (get admin program)
	if utils.IsUserChild(user) {
		res, admin := repositories.QueryUser(model.User{Id: user.ParentId.UUID})
		if res.Error != nil {
			utils.PrintObj(res.Error.Error(), "GetProgramDataJson QueryAdmin")
		}

		if admin.ProgramId.Valid {
			programId = admin.ProgramId.UUID
		} else {
			utils.PrintObj(admin, "GetProgramDataJson admin has not program")
			return defaultVal, Common.ErrorCodes_UNKNOWN_ERROR
		}

	} else {
		programId = user.ProgramId.UUID
	}

	// get program data
	program := model.Program{
		Id: programId,
	}

	result, program := repositories.QueryProgram(program)
	if result.Error != nil {
		utils.PrintObj(result.Error.Error(), "GetProgramDataJson QueryProgram")

		return defaultVal, Common.ErrorCodes_UNKNOWN_ERROR
	}

	// utils.PrintObj(utils.ToJson(program), "GetProgramDataJson")

	return utils.ToJson(program), Common.ErrorCodes_SUCCESS
}
