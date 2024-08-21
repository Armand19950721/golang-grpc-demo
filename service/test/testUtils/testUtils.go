package testUtils

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"service/protos/Common"
	"service/protos/ThirdPartyCommon"
	"service/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	// ip = "127.0.0.1" // local
	ip = "18.180.240.3" // dev
	// ip      = "54.64.52.208" // prod
	addr              = flag.String("addr", ip+":20000", "the address to connect to")
	Conn, _           = grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	addrThirdParty    = flag.String("addrThirdParty", ip+":20100", "the address to connect to")
	ConnThirdParty, _ = grpc.Dial(*addrThirdParty, grpc.WithTransportCredentials(insecure.NewCredentials()))
)

var (
	basicInfoPath    = "./testUtils/basicInfo.json"
	basicInfoPathDev = "./test/testUtils/basicInfo.json"
)

type AdminBasicInfo struct {
	AdminId     string
	Account     string
	Pwd         string
	Token       string
	ProgramId   string
	ArContentId string
	Child       ChildInfo
	Remark      string
	TestFaild   bool
}

type ChildInfo struct {
	Account string
	Pwd     string
	Token   string
}

func GetCtx(params ...string) context.Context {
	if len(params) >= 1 {
		return metadata.NewOutgoingContext(context.Background(), GetMedataObject("use Child Token"))
	} else {
		return metadata.NewOutgoingContext(context.Background(), GetMedataObject())
	}
}

func GetThirdPartyCtx(params ...string) context.Context {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))

	info := AdminBasicInfo{}

	if len(params) > 0 {
		info = GetBasicInfo(params[0])
	} else {
		info = GetBasicInfo()
	}

	fmt.Println("info.Token", info.Token)
	md.Set("authorization", info.Token)

	return metadata.NewOutgoingContext(context.Background(), md)
}

func SetBasicInfo(info AdminBasicInfo, params ...string) {
	// utils.PrintObj(info, "SetBasicInfo")
	byteData := []byte(utils.ToJson(info))
	path := basicInfoPath // target from different startup location
	if len(params) >= 1 {
		path = basicInfoPathDev
	}

	err := os.WriteFile(path, byteData, 0644)

	if err != nil {
		utils.PrintObj(err, "err setBasicInfo")
	}
}

func GetBasicInfo(params ...string) AdminBasicInfo {
	path := basicInfoPath // target from different startup location

	if len(params) > 0 {
		path = basicInfoPathDev
	}

	data, err := os.ReadFile(path)
	if err != nil {
		utils.PrintObj(err, "err ReadBasicInfo")
	}

	res, err := utils.ParseJsonWithType[AdminBasicInfo](string(data))
	if err != nil {
		DisplayResult("", err, false)
	}

	// utils.PrintObj(res, "ReadBasicInfo")
	return res
}

func ParseErrorCode(err error) (Common.ErrorCodes, string, bool) {
	if strings.Contains(err.Error(), "{") {
		arr := strings.Split(err.Error(), "{")
		// utils.PrintObj(arr)
		json := "{" + arr[1]
		// utils.PrintObj(utils.ValidJson(json))
		if utils.ValidJson(json) {
			parse, errParse := utils.ParseJsonWithType[Common.ErrorReply](json)
			// utils.PrintObj(err)
			// utils.PrintObj(parse.Code.String())
			if errParse == nil {
				return parse.Code, parse.Message, true
			} else {
				utils.PrintObj(errParse.Error())
			}
		}
	}

	return Common.ErrorCodes_ErrorCodes_NONE, "", false
}

type Stateful interface {
	GetState() ThirdPartyCommon.StatusCode
}

func DisplayResultThirdParty(obj Stateful, err error, expectErr bool) {
	resultJson, _ := json.Marshal(obj)
	panicMsg := "Encountered unknown error. Test stoped."

	if err != nil {
		panic("internal fail:" + err.Error())
	}

	// display res
	fmt.Println(string(resultJson))
	fmt.Println(obj.GetState().String())

	if obj.GetState() == ThirdPartyCommon.StatusCode_SUCCESS {
		fmt.Println("--> success")
		// success
		if expectErr {
			fmt.Println("=====================> not expected <=====================")
			panic(panicMsg)
		}
	} else {
		fmt.Println("--> fail")
		// fail
		if !expectErr {
			fmt.Println("=====================> not expected <=====================")
			panic(panicMsg)
		}
	}

}

func DisplayResult(obj interface{}, err error, expectErr bool) {
	mdJson, _ := json.Marshal(obj)
	panicMsg := "Encountered unknown error. Test stoped."

	if err != nil {
		fmt.Println("--> fail")

		parse, msg, parseOk := ParseErrorCode(err)

		if parseOk {
			fmt.Println(parse.String())
			fmt.Println(msg)
		} else {
			fmt.Println(err.Error())
		}

		fmt.Println(string(mdJson))

		if !expectErr {
			fmt.Println("=====================> not expected <=====================")
			panic(panicMsg)
		}
	} else {
		fmt.Println("--> success")
		fmt.Println(string(mdJson))

		if expectErr {
			fmt.Println("=====================> not expected <=====================")
			panic(panicMsg)
		}
	}
}

var (
	SuccessCount = 0
	ErrorCount   = 0
)

func CountResult(obj interface{}, err error, expectErr bool) {
	if err != nil || obj == nil {
		ErrorCount++
		// fmt.Println("--> fail")
		// fmt.Println(ErrorCount)

	} else {
		SuccessCount++
		// fmt.Println("--> success")
		// fmt.Println(SuccessCount)
	}
}

func RandomIntString() string {
	uuidNew := uuid.New().String()
	splitRes := strings.Split(uuidNew, "-")
	return splitRes[0]
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	rdnNum := min + rand.Intn(max-min)
	return rdnNum
}

func GetMedataObject(params ...string) metadata.MD {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	info := GetBasicInfo()
	// actionString := license_server_protobuf.LICENSE_ACTION_INIT.String()
	// actionInt := int(license_server_protobuf.LICENSE_ACTION_value[actionString])

	if len(params) >= 1 {
		md.Set("authorization", info.Child.Token)
	} else {
		md.Set("authorization", info.Token)
	}
	// md.Set("license_id", "ebc536af-a0a7-4934-859c-baaa18a80596")
	// md.Set("scret_key", "a4faf819c1b9f03ca5042a8106e16f2b9b51086eea35921c4ee7f0fb082cf574")
	// md.Set("application_name", "ios38336")
	// md.Set("action", utils.ToString(actionInt))
	// md.Set("license_action_raw", utils.ToJson("fake raw"))
	// utils.PrintObj(md, "SetMedataObject")
	return md
}

// file content is a struct which contains a file's name, its type and its data.
type FileContent struct {
	name  string
	field string
	data  []byte
}

type IUploadFile struct {
	Message    string
	Code       int
	UploadName string
}

func UploadFile(filename string) string {
	path := "./static/"
	// filename := "thumbnail-1.png"
	data, err := os.ReadFile(path + filename)
	if err != nil {
		utils.PrintObj(err.Error())
		return ""
	}

	c := FileContent{
		name:  filename,
		field: "file",
		data:  data,
	}

	basicInfo := GetBasicInfo()
	res, err := sendPostRequest("http://"+ip+":8080/api/upload_file", basicInfo.Token, c)
	// res, err := sendPostRequest("https://dev-api-v2.metarcommerce.com/api/upload_file", basicInfo.Token, c)

	if err != nil {
		utils.PrintObj(err.Error(), "UploadFile")
		return ""
	}

	return string(res)
}

func sendPostRequest(url string, token string, files ...FileContent) ([]byte, error) {
	var (
		buf    = new(bytes.Buffer)
		writer = multipart.NewWriter(buf)
	)

	for _, file := range files {
		part, err := writer.CreateFormFile(file.field, filepath.Base(file.name))
		if err != nil {
			return []byte{}, err
		}
		part.Write(file.data)
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("authorization", token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	utils.PrintObj(string(cnt), "sendPostRequest cnt")
	return cnt, nil
}
