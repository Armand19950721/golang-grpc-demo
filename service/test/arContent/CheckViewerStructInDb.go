package arContent

import (
	"service/model"
	"service/protos/ArContent"
	"service/utils"
)

func MainProcess() {
	datas := []model.ArContent{}
	res := utils.DatabaseManagerSlave.Debug().Raw(`select * from ar_content`).Scan(&datas)

	if res.Error != nil {
		utils.PrintObj(res.Error.Error(), "err")
		return
	}

	utils.PrintObj(len(datas), "datas len")

	for _, item := range datas {
		viewerSetting, err := utils.ParseJsonWithType[ArContent.ArViewerSetting](item.ViewerSetting)

		utils.PrintObj(&viewerSetting, "viewerSetting")

		// if had err then confirm is new struct
		if err != nil {
			utils.PrintObj(err.Error(), " err")
			utils.PrintObj(item, " err item")
			panic("")
		}

	}
}
