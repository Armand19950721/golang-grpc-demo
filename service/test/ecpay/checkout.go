package ecpay

import (
	"errors"
	"service/utils"
	"service/utils/httpHelper"
	"strings"

	"github.com/google/uuid"
)

func GetRandomString(length int) (error, string) {

	if length > 32 {
		return errors.New("length invalid. max length is 32"), ""
	}

	randomStr := uuid.New().String()
	randomStr = strings.Replace(randomStr, "-", "", 4)
	randomStr = randomStr[0:length]

	return nil, randomStr
}

type CheckOutModel struct {
	MerchantTradeNo   string // ex:'f0a0d7e9fae1bb72bc93', //請帶20碼uid
	MerchantTradeDate string // ex:'2017/02/13 15:45:30', //交易時間
	TotalAmount       string // ex:'100',
	TradeDesc         string // ex:'測試交易描述',
	ItemName          string // ex:'測試商品等',
	ReturnURL         string // ex:'http://192.168.0.1',
}

type IssueInvoiceModel struct {
	RelateNumber       string // ex: "asc1232aaasadFY", // 請帶30碼uid, ex: werntfg9os48trhw34etrwerh8ew2r
	CustomerID         string // ex: "12124", // 客戶代號，長度為20字元
	CustomerIdentifier string // ex: "", // 統一編號，長度為8字元
	CustomerName       string // ex: "綠先生", // 客戶名稱，長度為20字元
	CustomerAddr       string // ex: "台北市南港區三重路19-2號6-2樓()", // 客戶地址，長度為100字元
	CustomerPhone      string // ex: "0930597532", // 客戶電話，長度為20字元
	CustomerEmail      string // ex: "asdfnmb24040@gmail.com", // 客戶信箱，長度為80字元
	ClearanceMark      string // ex: "", // 通關方式，僅可帶入'1'、'2'、''
	Print              string // ex: "1", // 列印註記，僅可帶入'0'、'1'
	Donation           string // ex: "0", // 捐贈註記，僅可帶入'1'、'0'
	LoveCode           string // ex: "", // 愛心碼，長度為7字元
	CarruerType        string // ex: "", // 載具類別，僅可帶入'1'、'2'、'3'、''
	CarruerNum         string // ex: "", // 載具編號，當載具類別為'2'時，長度為16字元，當載具類別為'3'時，長度為7字元
	TaxType            string // ex: "1", // 課稅類別，僅可帶入'1'、'2'、'3'、'9'
	SalesAmount        string // ex: "200", // 發票金額
	InvoiceRemark      string // ex: "", // 備註
	ItemName           string // ex: "洗衣精|洗髮乳", // 商品名稱，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	ItemCount          string // ex: "1|1", // 商品數量，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	ItemWord           string // ex: "瓶|罐", // 商品單位，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	ItemPrice          string // ex: "100|100", // 商品價格，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	ItemTaxType        string // ex: "", // 商品課稅別，如果超過一樣商品時請以｜(為半形不可使用全形)分隔，如果TaxType為9請帶值，其餘為空
	ItemAmount         string // ex: "100|100", // 商品合計，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	ItemRemark         string // ex: "test item|test item", // 商品備註，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
	InvType            string // ex: "07", // 字軌類別，、'07'一般稅額
	Vat                string // ex: "1" // 商品單價是否含稅，'1'為含稅價'、'2'為未稅價
}

type EcpayConfigModel struct {
	OperationMode       string
	MercProfile         MercProfileModel
	IgnorePayment       []string
	IsProjectContractor bool
}
type MercProfileModel struct {
	MerchantID string
	HashKey    string
	HashIV     string
}

func GetEcpayConfig() EcpayConfigModel {

	err, isProjectContractor := utils.ToBool(utils.GetEnv("ECPAY_IS_PROJECT_CONTRACTOR"))
	if err != nil {
		panic("get ECPAY_IS_PROJECT_CONTRACTOR err:" + err.Error())
	}

	config := EcpayConfigModel{
		OperationMode: utils.GetEnv("ECPAY_OPERATION_MODE"),
		MercProfile: MercProfileModel{
			MerchantID: utils.GetEnv("ECPAY_MERCHANT_ID"),
			HashKey:    utils.GetEnv("ECPAY_HASH_KEY"),
			HashIV:     utils.GetEnv("ECPAY_HASH_IV"),
		},
		IgnorePayment: []string{
			"WebATM",
			"ATM",
			"CVS",
			"BARCODE",
			"AndroidPay",
		},
		IsProjectContractor: isProjectContractor,
	}

	return config
}

func CheckOut() {

	// prepare data
	err, merchantTradeNo := GetRandomString(20)
	if err != nil {
		panic(err.Error())
	}

	merchantTradeDate := utils.ParseDate(utils.ParseTimeToUTC8(utils.GetTodayDateTime()), utils.DateLayoutEcpay)

	checkOutData := CheckOutModel{
		MerchantTradeNo:   merchantTradeNo,   //交易uid 20碼
		MerchantTradeDate: merchantTradeDate, //交易日 ex:"2006/01/02 15:04:05"
		TotalAmount:       "100",
		TradeDesc:         "測試交易描述",
		ItemName:          "測試商品等",
		ReturnURL:         "https://www.google.com/",
	}

	body := map[string][]string{
		"checkOutData": {utils.ToJson(checkOutData)},
		"ecpayConfig":  {utils.ToJson(GetEcpayConfig())},
		//pkey
	}

	header := map[string][]string{
		"auth": {utils.GetEnv("SERVICE_PRIVATE_KEY")},
	}

	// post to nodejs ecpay service
	res, err := httpHelper.SendPostRequest(
		utils.GetEnv("ECPAY_SERVICE_DOMAIN")+utils.GetEnv("ECPAY_SERVICE_ROUTE_CHECK_OUT"),
		header, body)
	if err != nil {
		// ecpay node service is dead maybe
		panic(err.Error())
	}

	parseRes, err := utils.ParseJsonWithType[EcpayServiceReturn](res)
	if err != nil {
		//check out response format err
		panic(err.Error())
	}

	if !parseRes.Success {
		//check out fail
	}

	//return success

	utils.PrintObj(parseRes, "parseRes")
}

func InvoiceIssue() {
	// prepare data
	err, relateNumber := GetRandomString(30)
	if err != nil {
		panic(err.Error())
	}

	customerEmail := "asdfnmb24040@gmail.com"
	salesAmount := "100"
	itemName := "洗衣精"
	itemCount := "1"
	itemWord := "瓶"
	itemPrice := "100"
	itemAmount := "100"
	itemRemark := "test item remark"
	invType := "07"

	invoice := IssueInvoiceModel{
		RelateNumber:       relateNumber,  // 請帶30碼uid, ex: werntfg9os48trhw34etrwerh8ew2r
		CustomerID:         "",            // 客戶代號，長度為20字元
		CustomerIdentifier: "",            // 統一編號，長度為8字元
		CustomerName:       "",            // 客戶名稱，長度為20字元
		CustomerAddr:       "",            // 客戶地址，長度為100字元
		CustomerPhone:      "",            // 客戶電話，長度為20字元
		CustomerEmail:      customerEmail, // 客戶信箱，長度為80字元
		ClearanceMark:      "",            // 通關方式，僅可帶入'1'、'2'、''
		Print:              "0",           // 列印註記，僅可帶入'0'、'1'
		Donation:           "0",           // 捐贈註記，僅可帶入'1'、'0'
		LoveCode:           "",            // 愛心碼，長度為7字元
		CarruerType:        "",            // 載具類別，僅可帶入'1'、'2'、'3'、''
		CarruerNum:         "",            // 載具編號，當載具類別為'2'時，長度為16字元，當載具類別為'3'時，長度為7字元
		TaxType:            "1",           // 課稅類別，僅可帶入'1'、'2'、'3'、'9'
		SalesAmount:        salesAmount,   // 發票金額
		InvoiceRemark:      "",            // 備註
		ItemName:           itemName,      // 商品名稱，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		ItemCount:          itemCount,     // 商品數量，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		ItemWord:           itemWord,      // 商品單位，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		ItemPrice:          itemPrice,     // 商品價格，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		ItemTaxType:        "",            // 商品課稅別，如果超過一樣商品時請以｜(為半形不可使用全形)分隔，如果TaxType為9請帶值，其餘為空
		ItemAmount:         itemAmount,    // 商品合計，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		ItemRemark:         itemRemark,    // 商品備註，如果超過一樣商品時請以｜(為半形不可使用全形)分隔
		InvType:            invType,       // 字軌類別，、'07'一般稅額
		Vat:                "1",           // 商品單價是否含稅，'1'為含稅價'、'2'為未稅價
	}

	body := map[string][]string{
		"invoice": {utils.ToJson(invoice)},
		//pkey
	}

	header := map[string][]string{
		"auth": {utils.GetEnv("SERVICE_PRIVATE_KEY")},
	}

	// post to nodejs ecpay service
	res, err := httpHelper.SendPostRequest(
		utils.GetEnv("ECPAY_SERVICE_DOMAIN")+utils.GetEnv("ECPAY_SERVICE_ROUTE_INVOICE_ISSUE"),
		header, body)
	if err != nil {
		// ecpay node service is dead maybe
		panic(err.Error())
	}

	parseRes, err := utils.ParseJsonWithType[EcpayServiceReturn](res)
	if err != nil {
		//check out response format err
		utils.PrintObj(parseRes)
		panic(err.Error())
	}

	if !parseRes.Success {
		//check out fail
	}
}

type EcpayServiceReturn struct {
	Success bool
	Msg     string
	Data    string
}
