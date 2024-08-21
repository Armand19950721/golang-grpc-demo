package main

import (
	"service/protos/ArContent"
	"service/test/accountSettings"
	arContentTest "service/test/arContent"
	"service/test/auth"
	"service/test/echo"
	"service/test/permission"
	"service/test/program"
	"service/test/s3"

	"flag"
	"service/test/thirdParty"
	"service/test/unittest"
	"service/test/user"
	"service/utils"
)

func main() {
	flag.Parse()

	// cmd
	switch flag.Arg(0) {
	case "init":
		initProcess()
	case "dev-init":
		devInit()
	case "tmp":
		utils.InitDB()
		program.GetFreeProgram()
		auth.RegAndEmailFlow()
	case "echo":
		echo.Flow()
	default:
		basicTestFlow()
	}
}

func basicTestFlow() {
	utils.InitDB()
	unittest.Flow()
	thirdParty.MainProcess()
	regularAdminTesting()
	adminArContentMainPageFlow()
	adminArContentFlow()
}


func adminArContentMainPageFlow() {
	utils.PrintObj("", "adminArContentMainPageFlow")
	// createNewUser()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES, false)
	arContentTest.TemplateGlassesUpdateFlow()
	arContentTest.UpdateViewSetting()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES, false)
	arContentTest.TemplateContactLensesUpdateFlow()
	arContentTest.UpdateViewSetting()

	arContentTest.ArContentMainPageFlow()
}

func adminArContentFlow() {
	utils.PrintObj("", "adminArContentFlow")
	// program.CreateProgram(program.ProgramAdmin)
	// auth.Flow()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES, false)
	arContentTest.TemplateGlassesUpdateFlow()
	arContentTest.UpdateViewSetting()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_EARRING, false)
	arContentTest.TemplateEarringUpdateFlow()
	arContentTest.UpdateViewSetting()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES, false)
	arContentTest.TemplateContactLensesUpdateFlow()
	arContentTest.UpdateViewSetting()

	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_HAIR, false)
	arContentTest.TemplateHairUpdateFlow()
	arContentTest.UpdateViewSetting()

	// test limited program flow until Type
	program.CreateProgram(program.ProgramLimit)
	user.UpdateUserProgram()
	arContentTest.CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_EARRING, true)
}

func createNewUser() {
	program.GetFreeProgram()
	program.CreateProgram(program.ProgramAdmin)
	auth.Flow()
}

func regularAdminTesting() {
	utils.PrintObj("", "regularAdminTesting")

	createNewUser()
	accountSettings.Flow()
	permission.Flow()
	user.ChangePasswordFlow()
	user.UserChildFlow()
	s3.Flow()
}

func initProcess() {
	utils.PrintObj("", "initProcess")
	program.CreateFreeProgram(program.ProgramFree)
}

func devInit() {
	utils.PrintObj("", "devInit")

	program.CreateProgram(program.ProgramAdmin)
	auth.RegDevGuest()
	auth.LoginDevGuest()
}
func pinkoiInit() {
	utils.PrintObj("", "pinkoiInit")

	program.CreateProgram(program.ProgramAdmin)
	auth.RegDevPinkoi()
	auth.LoginDevPinkoi()
}
