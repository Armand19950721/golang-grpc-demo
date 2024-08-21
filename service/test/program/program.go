package program

import (
	"service/test/testUtils"
	"service/utils"

	"service/protos/ArContent"
	"service/protos/Common"
	"service/protos/Program"
	"service/protos/WebServices"
)

var (
	ProgramFree = &Program.ProgramModel{
		Name:  "program_free",
		State: Common.ProgramState_ON_SALE,
		Seats: 3,

		Categories: []ArContent.ArContentCategoryEnum{
			ArContent.ArContentCategoryEnum_CATRGORY_ALL,
		},
		Types: []ArContent.ArContentTypeEnum{
			ArContent.ArContentTypeEnum_TYPE_ALL,
		},
		Templates: []ArContent.ArContentTemplateEnum{
			ArContent.ArContentTemplateEnum_TEMPLATE_ALL,
		},
		EffectTools:         []Common.EffectTool{},
		ArInteractModules:   []Common.ArInteractModule{},
		ArEditWindowModules: []Common.ArEditWindowModule{},
	}

	ProgramAdmin = &Program.ProgramModel{
		Name:  testUtils.RandomIntString() + "_program",
		State: Common.ProgramState_ON_SALE,
		Seats: 3,

		Categories: []ArContent.ArContentCategoryEnum{
			ArContent.ArContentCategoryEnum_CATRGORY_ALL,
		},
		Types: []ArContent.ArContentTypeEnum{
			ArContent.ArContentTypeEnum_TYPE_ALL,
		},
		Templates: []ArContent.ArContentTemplateEnum{
			ArContent.ArContentTemplateEnum_TEMPLATE_ALL,
		},
		EffectTools:         []Common.EffectTool{Common.EffectTool_BEAUTIFUL_SKIN},
		ArInteractModules:   []Common.ArInteractModule{},
		ArEditWindowModules: []Common.ArEditWindowModule{Common.ArEditWindowModule_AR_EDIT_WINDOW_MODULE_NONE, Common.ArEditWindowModule_UI_EDIT},
	}
	ProgramLimit = &Program.ProgramModel{
		Name:  testUtils.RandomIntString() + "_program",
		State: Common.ProgramState_ON_SALE,
		Seats: 3,

		Categories: []ArContent.ArContentCategoryEnum{
			ArContent.ArContentCategoryEnum_CATRGORY_HEAD,
		},
		Types: []ArContent.ArContentTypeEnum{
			ArContent.ArContentTypeEnum_TYPE_GLASSES,
		},
		Templates: []ArContent.ArContentTemplateEnum{
			ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES,
		},
		EffectTools:         []Common.EffectTool{Common.EffectTool_BEAUTIFUL_SKIN},
		ArInteractModules:   []Common.ArInteractModule{},
		ArEditWindowModules: []Common.ArEditWindowModule{Common.ArEditWindowModule_AR_EDIT_WINDOW_MODULE_NONE, Common.ArEditWindowModule_UI_EDIT},
	}

	c = WebServices.NewProgramServicesClient(testUtils.Conn)
)

func CreateProgram(model *Program.ProgramModel) {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("Program create")

	model.Name = testUtils.RandomIntString() + "_program"

	r, err := c.CreateProgram(ctx, &Program.CreateProgramRequest{
		Model: model,
	})
	testUtils.DisplayResult(r, err, false)

	basicInfo := testUtils.GetBasicInfo()
	basicInfo.ProgramId = r.ProgramId
	testUtils.SetBasicInfo(basicInfo)

	utils.PrintTitle("Program Get List")

	rR, err := c.GetProgramList(ctx, &Program.GetProgramListRequest{})
	testUtils.DisplayResult(len(rR.Models), err, false)
}

func GetFreeProgram() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("GetFreeProgram")

	r, err := c.GetFreeProgram(ctx, &Program.GetFreeProgramRequest{})
	testUtils.DisplayResult(r, err, false)

	basicInfo := testUtils.GetBasicInfo()
	basicInfo.ProgramId = r.ProgramId
	testUtils.SetBasicInfo(basicInfo)
}

func CreateFreeProgram(model *Program.ProgramModel) {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("Program create")

	r, err := c.CreateProgram(ctx, &Program.CreateProgramRequest{
		Model: model,
	})
	testUtils.DisplayResult(r, err, false)

	basicInfo := testUtils.GetBasicInfo()
	basicInfo.ProgramId = r.ProgramId
	testUtils.SetBasicInfo(basicInfo)

	utils.PrintTitle("Program Get List")

	rR, err := c.GetProgramList(ctx, &Program.GetProgramListRequest{})
	testUtils.DisplayResult(len(rR.Models), err, false)
}
