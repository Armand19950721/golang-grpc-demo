package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"service/model"
	"service/protos/ArContent"
	"service/protos/Common"
	"service/protos/Program"

	valid "github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	// envFileName    = ".env-docker"
	envFileName     = ".env"
	loadEnvSuccess  = InitEnv()
	DateLayout      = "2006-01-02T15:04:05-0700"
	DateLayoutEcpay = "2006/01/02 15:04:05"
)

func InitEnv() bool {
	loadEnvErr := godotenv.Load("./" + envFileName)

	if loadEnvErr != nil {
		PrintObj("load env err. try another locaion")
		loadEnvErr = godotenv.Load("../" + envFileName) // only test use

		if loadEnvErr != nil {
			PrintObj("load env err. fail")
			return false
		}
	}

	PrintObj("load env suceess")
	return true
}

var (
	RecordExpireHours = 24
	UserEmailLimit    = 100
)

type UserTempRecord struct {
	EmailSendCount int
}

func CheckUserEmailLimit(userId uuid.UUID) (bool, error) {
	redisKey := "UserTempRecord:" + userId.String()

	// get redis record
	res := GetRedis(redisKey)

	if res != "" {
		userTempRecord, err := ParseJsonWithType[UserTempRecord](res)
		if err != nil {
			return false, err
		}

		if userTempRecord.EmailSendCount >= UserEmailLimit {
			// out of limit
			return false, nil
		} else {
			// add count and save
			userTempRecord.EmailSendCount++
			SetRedis(redisKey, ToJson(userTempRecord), int16(RecordExpireHours))

			return true, nil
		}
	}

	// create new record if not exist
	userTempRecord := UserTempRecord{EmailSendCount: 1}
	SetRedis(redisKey, ToJson(userTempRecord), int16(RecordExpireHours))

	return true, nil
}

func ConvertStartAndEndDateIfEmpty(startDate, endDate string) (string, string) {

	// 如果日期為空
	// 開始日會被換算成1900年
	// 結束日會被換算成明天
	if startDate == "" {
		startDate = GetMostEarlyDateIsoString()
	}

	if endDate == "" {
		endDate = GetTomorrowDateIsoString()
	}

	return startDate, endDate
}

func GetMostEarlyDateIsoString() string {
	return "1900-01-01T00:00:00.000Z"
}

func GetTomorrowDateIsoString() string {
	return ParseDateToIsoString(time.Now().Add(1 * 24 * time.Hour))
}

func ReturnError(code Common.ErrorCodes, internalMsg, returnMsg string) error {
	return status.Errorf(codes.Internal, GetError(ErrorType{
		Code:        code,
		InternalMsg: internalMsg,
		ReturnMsg:   returnMsg,
	}))
}

func ReturnUnKnownError(err error) error {
	return status.Errorf(codes.Internal, GetError(ErrorType{
		Code:        Common.ErrorCodes_UNKNOWN_ERROR,
		InternalMsg: err.Error(),
	}))
}

func ReturnBasicInfoError() error {
	return status.Errorf(codes.Internal, GetError(ErrorType{
		Code:        Common.ErrorCodes_UNKNOWN_ERROR,
		InternalMsg: "basicInfo error",
	}))
}

func IsErrorNotFound(err error) bool {
	res := errors.Is(err, gorm.ErrRecordNotFound)
	PrintObj(res, "IsErrorNotFound")

	return res
}

func GetRandomShortString() string {
	uuidNew := uuid.New().String()
	splitRes := strings.Split(uuidNew, "-")
	return splitRes[0]
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	rdnNum := min + rand.Intn(max-min)
	return rdnNum
}

func GetSuffixByFileName(name string) string {
	split := strings.Split(name, ".")
	return "." + split[len(split)-1]
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func RemoveFromArrayString(arr []string, val string) []string {
	idx := -1
	for i := 0; i < len(arr); i++ {
		if arr[i] == val {
			idx = i
		}
	}
	if idx != -1 {
		arr = RemoveIndex[string](arr, idx)
	}

	return arr
}

func IsTimeBetween(timeTarget, timeStart, timeEnd time.Time) bool {
	timeStart = timeStart.Add(time.Second * -1)
	timeEnd = timeEnd.Add(time.Second * 1)

	if timeTarget.Before(timeEnd) && timeTarget.After(timeStart) {
		PrintObj(true, "isTimeBetween")
		return true
	}

	PrintObj(false, "isTimeBetween")
	return false
}

func CheckExpire(expireTime time.Duration, updateDate time.Time) (bool, time.Time) {
	// check expire date
	endDate := updateDate.Add(expireTime)
	now := time.Now()

	PrintObj([]time.Time{endDate, now}, "endDate,now")

	if now.Before(endDate) {
		return false, endDate
	}

	return true, endDate
}

func IsUserChild(user *model.User) bool {
	return user.ParentId.Valid
}

func GetParentId(data BasicInfo) uuid.UUID {
	if IsUserChild(&data.User) {
		PrintObj("parent", "GetParentId")
		return data.User.ParentId.UUID
	} else {
		PrintObj("id", "GetParentId")
		return data.User.Id
	}
}

func GetNullableString(str string) sql.NullString {
	if len(str) == 0 {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
}

func SetNullableUUID(str string) uuid.NullUUID {
	// PrintObj(str, "SetNullableUUID")

	if len(str) == 0 || str == "" {
		return uuid.NullUUID{Valid: false}
	}

	return uuid.NullUUID{UUID: ParseUUID(str), Valid: true}
}

func ArrayToJson[T any](arr []T) string {
	newArr := []interface{}{}

	for i := 0; i < len(arr); i++ {
		newArr = append(newArr, arr[i])
	}

	return ToJson(newArr)
}

func SaveLog(key string, value interface{}) {
	valueJson := ToJson(value)

	res := DatabaseManager.Create(&model.Log{
		Key:   key,
		Value: valueJson,
	})

	if res.Error != nil {
		PrintObj(res.Error.Error(), "SaveLog error")
	}
}

func GetTodayDate() time.Time {
	t := time.Now()
	y, m, d := t.Date()

	//設定時區為UTC
	localLocation, _ := time.LoadLocation("UTC")
	utc := time.Date(y, m, d, 0, 0, 0, 0, localLocation)

	return utc
}

func GetTodayDateTime() time.Time {
	return time.Now()
}

func ParseTimeToUTC8(t time.Time) time.Time {
	return t.Add(time.Hour * 8)
}

func ValidIsoDate(date string, params ...string) bool {

	nullable := validParamsParse(params)

	PrintObj(date, "ValidIsoDate")
	// PrintObj(nullable, "nullable ValidIsoDate")

	if nullable && date == "" {
		return true
	}

	_, err := time.Parse(time.RFC3339, date)

	if err != nil {
		PrintObj(err.Error(), "ValidDate err")
		return false
	}

	return true
}

func CheckStringDateIsIncludingToday(startDateStr, endDateStr string) bool {
	PrintObj("CheckStringDateIsIncludingToday")

	includeToday := false

	// 判斷字串日期是否包含今天
	if ValidIsoDate(startDateStr) && ValidIsoDate(endDateStr) {
		start, err := ToIsoDate(startDateStr)
		if err != nil {
			PrintObj(err.Error(), "CheckStringDateIsIncludingToday startDateStr")
			return false
		}

		end, err := ToIsoDate(endDateStr)
		if err != nil {
			PrintObj(err.Error(), "CheckStringDateIsIncludingToday endDateStr")
			return false
		}

		PrintObj([]time.Time{start, end}, "start,end")
		includeToday = IsTimeBetween(GetTodayDate(), start, end)
	}

	PrintObj(includeToday, "includeToday")

	return includeToday
}

func ToIsoDate(date string) (time.Time, error) {
	res, err := time.Parse(time.RFC3339, date)

	if err != nil {
		PrintObj(err.Error(), "ToIsoDate err")
	}

	return res, err
}

func GetRedisFolderDateFormat(day int) string {
	// return "2022-09-21"
	return GetTodayDate().Add(time.Duration(day) * 24 * time.Hour).String()[0:10] + ":"
}

func fixStringLen(str string) string {
	limit := 500

	if len(str) > limit {
		return str[0:limit] + "....more"
	}
	return str
}

type BasicInfo struct {
	AdminId uuid.UUID
	User    model.User
	Program *Program.ProgramModel
	Token   string
	Success bool
}

func GetBasicInfo(ctx context.Context) BasicInfo {
	basicInfo := BasicInfo{
		User:    model.User{},
		Program: &Program.ProgramModel{},
		Success: false,
	}

	// parse metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		PrintObj("parse fail", "get basic info md err")
		return basicInfo
	}

	// get data and parse
	userDataJson := GetMetaDataField(md, "user_data_json")
	programDataJson := GetMetaDataField(md, "program_data_json")
	tokenStr := GetMetaDataField(md, "authorization")

	PrintObj([]int{len(userDataJson), len(programDataJson), len(tokenStr)}, "get basic info")

	if len(userDataJson) == 0 || len(programDataJson) == 0 || len(tokenStr) == 0 {
		PrintObj("data length invalid", "get basic info err")
		return basicInfo
	}

	// convert
	user, err := ParseJsonWithType[model.User](userDataJson)
	if err != nil {
		PrintObj(err.Error(), "user")
		basicInfo.Success = false
		return basicInfo
	}
	programTable, err := ParseJsonWithType[model.Program](programDataJson)
	if err != nil {
		PrintObj(err.Error(), "programTable")
		basicInfo.Success = false
		return basicInfo
	}
	arContentCategoryEnum, err := ParseJsonWithType[[]ArContent.ArContentCategoryEnum](programTable.Categories)
	if err != nil {
		PrintObj(err.Error(), "arContentCategoryEnum")
		basicInfo.Success = false
		return basicInfo
	}
	arContentTypeEnum, err := ParseJsonWithType[[]ArContent.ArContentTypeEnum](programTable.Types)
	if err != nil {
		PrintObj(err.Error(), "arContentTypeEnum")
		basicInfo.Success = false
		return basicInfo
	}
	arContentTemplateEnum, err := ParseJsonWithType[[]ArContent.ArContentTemplateEnum](programTable.Templates)
	if err != nil {
		PrintObj(err.Error(), "arContentTemplateEnum")
		basicInfo.Success = false
		return basicInfo
	}
	effectTool, err := ParseJsonWithType[[]Common.EffectTool](programTable.EffectTools)
	if err != nil {
		PrintObj(err.Error(), "effectTool")
		basicInfo.Success = false
		return basicInfo
	}
	arInteractModule, err := ParseJsonWithType[[]Common.ArInteractModule](programTable.ArInteractModules)
	if err != nil {
		PrintObj(err.Error(), "arInteractModule")
		basicInfo.Success = false
		return basicInfo
	}
	arEditWindowModule, err := ParseJsonWithType[[]Common.ArEditWindowModule](programTable.ArEditWindowModules)
	if err != nil {
		PrintObj(err.Error(), "arEditWindowModule")
		basicInfo.Success = false
		return basicInfo
	}

	basicInfo.User = user
	// parse to proto model
	model := &Program.ProgramModel{
		Id:                  programTable.Id.String(),
		Name:                programTable.Name,
		State:               programTable.State,
		Seats:               programTable.Seats,
		Categories:          arContentCategoryEnum,
		Types:               arContentTypeEnum,
		Templates:           arContentTemplateEnum,
		EffectTools:         effectTool,
		ArInteractModules:   arInteractModule,
		ArEditWindowModules: arEditWindowModule,
	}

	// set admin id
	if IsUserChild(&basicInfo.User) {
		basicInfo.AdminId = basicInfo.User.ParentId.UUID
	} else {
		basicInfo.AdminId = basicInfo.User.Id
	}

	// set return
	basicInfo.Program = model
	basicInfo.Token = tokenStr
	basicInfo.Success = true

	PrintObj(basicInfo, "basicInfo")

	return basicInfo
}

func GetMetaDataField(md metadata.MD, name string) string {
	arr := md[name]
	if len(arr) < 1 {
		return ""
	}
	// PrintObj("length:"+ToString(len(arr[0])), name)
	// PrintObj(arr[0], name)
	return arr[0]
}

func ParseDateToString(datetime time.Time) string {
	return datetime.UTC().Format(DateLayout)
}

func ParseDateToIsoString(datetime time.Time) string {
	return datetime.UTC().Format(time.RFC3339)
}

func ParseDate(datetime time.Time, layout string) string {
	return datetime.UTC().Format(layout)
}

func FileOrFolderExist(location string) bool {
	if _, err := os.Stat(location); err != nil {
		return false
	}

	return true
}

func CreateFolder(location string) bool {
	PrintObj("Path:"+location, "CreateFolder")

	if FileOrFolderExist(location) {
		PrintObj("CreateFolder File Exist", "")
		return true
	}

	err := os.MkdirAll(location, 0755)
	if err != nil {
		PrintObj(err.Error())
		return false
	}

	return true
}

func ParseUUID(str string) uuid.UUID {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil
	}

	return id
}

func Copy(oldLocation, newLocation string) bool {
	sourceFileStat, err := os.Stat(oldLocation)
	if err != nil {
		PrintObj(err.Error(), "copy err:")
		return false
	}

	if !sourceFileStat.Mode().IsRegular() {
		PrintObj("", "copy err: oldLocation is not a regular file")
		return false
	}

	source, err := os.Open(oldLocation)
	if err != nil {
		PrintObj(err.Error(), "copy err")
		return false
	}
	defer source.Close()

	destination, err := os.Create(newLocation)
	if err != nil {
		PrintObj(err.Error(), "copy err")
		return false
	}

	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		PrintObj(err.Error(), "copy err")
		return false
	}

	PrintObj("", "copy file ok")
	PrintObj(oldLocation)
	PrintObj(newLocation)

	return true
}

func chosePath(uploadType string) string {
	switch uploadType {
	case "model_object":
		return GetEnv("ASSET_MODEL_OBJECT_PATH")
	case "thumbnail":
		return GetEnv("ASSET_THUMBNAIL_PATH")
	default:
		PrintObj("chosePath: uploadType Invalid")
		return ""
	}
}

func DeleteFile(location string) bool {
	PrintObj(location, "deleteFile file")

	err := os.Remove(location)
	if err != nil {
		PrintObj(err.Error(), "")
		return false
	}

	return true
}

func EncodePwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}

func CheckPwd(encodePwd, incomingPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePwd), []byte(incomingPwd))
	if err != nil {
		fmt.Println("pw wrong")
		return false
	} else {
		fmt.Println("pw ok")
		return true
	}
}

func validParamsParse(params []string) bool {

	if len(params) > 0 {
		if params[0] != "" {
			return true
		}
	}

	return false
}

func ValidId(str string, params ...string) bool {
	nullable := validParamsParse(params)

	// PrintObj(str, "str ValidId")
	// PrintObj(nullable, "nullable ValidId")

	if nullable && str == "" {
		return true
	}

	if str == "00000000-0000-0000-0000-000000000000" {
		return false
	}

	return valid.IsUUID(str)
}

func ValidEmail(str string, params ...string) bool {
	nullable := validParamsParse(params)
	// PrintObj(str, "str ValidString")
	// PrintObj(GetLength(str))
	// PrintObj(nullable)

	if nullable && str == "" {
		return true
	}

	if !ValidString(str, 1, 100) {
		return false
	}

	if !valid.IsEmail(str) {
		return false
	}

	return true
}

func ValidPassword(str string) bool {
	return ValidString(str, 6, 100)
}

func GetLength(str string) int {
	return utf8.RuneCountInString(str)
}

func ToString(val int) string {
	return strconv.Itoa(val)
}

func ToInt(number string) int {
	val, err := strconv.Atoi(number)
	if err != nil {
		return -1
	}

	return val
}

func AutoRename(str string) string {
	defaultReturn := str + "(1)"
	arr := strings.Split(str, "")

	// check final str is ")" if not then simple add "(1)"
	if arr[len(str)-1] != ")" {
		PrintObj(defaultReturn, "add number")
		return defaultReturn
	}

	// rename number
	// get final "(" idx
	result := ""
	targetIdx := -1
	for idx, item := range arr {
		if item == "(" {
			targetIdx = idx
		}
	}

	PrintObj(targetIdx, "targetIdx")

	if targetIdx == -1 {
		PrintObj(defaultReturn, "not found number")
		return defaultReturn
	}

	// try to add num in string
	for idx, item := range arr {
		result += item

		if idx == targetIdx {
			numStrRes := ""
			idxCount := 1

			// get all num in "()"
			for {
				nextStr := arr[idx+idxCount]
				num := ToInt(nextStr)

				if num == -1 {
					break
				}

				numStrRes += nextStr
				idxCount++
			}

			if numStrRes == "" {
				PrintObj(defaultReturn, "cant found num")
				return defaultReturn
			}

			// add 1
			numRes := ToInt(numStrRes)
			if numRes == -1 {
				PrintObj(defaultReturn, "num cant parse to int")
				return defaultReturn
			}

			numRes += 1

			// build result string
			result += ToString(numRes) + ")"

			break
		}
	}

	PrintObj(result, "result")
	return result
}

func ValidNumber(num, max, min int, params ...string) bool {
	nullable := validParamsParse(params)

	// PrintObj(num, "num ValidNumber")
	// PrintObj(nullable, "nullable ValidNumber")

	if nullable && num == 0 {
		return true
	}

	if num <= min && num >= max {
		return true
	}

	return false
}

func ValidModelObjectType(num uint32) bool {
	return num > 0 //model object type start from 1
}

func ValidString(str string, min, max int, params ...string) bool {
	nullable := validParamsParse(params)
	// PrintObj(str, "str ValidString")
	// PrintObj(GetLength(str))
	// PrintObj(nullable)

	if nullable && str == "" {
		return true
	}

	if max == -1 { //ignore max
		if GetLength(str) >= min {
			return true
		}
	} else {
		if GetLength(str) <= max && GetLength(str) >= min {
			return true
		}
	}

	return false
}

// func HmacSha256(key, data string) string {
// 	//Hash-based Message Authentication Code
// 	hash := hmac.New(sha256.New, []byte(key)) //using sha256
// 	hash.Write([]byte(data))
// 	return hex.EncodeToString(hash.Sum([]byte("")))
// }

func GenerateBearerToken() string {
	secretKey := "mysecretkey"
	issuer := "myapp"
	subject := "12345"
	expirationTime := time.Now().Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"exp": expirationTime.Unix(),
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic("GenerateBearerToken fail")
	}

	// return "Bearer " + signedToken
	return signedToken
}

func PrintTitle(text string) {
	fmt.Println("")
	fmt.Println("-------------- " + text + " --------------")
}

func ParseBoolToString(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func ToJson(obj interface{}) string {
	mdJson, err := json.Marshal(obj)

	if err != nil {
		fmt.Println("to json err")
		return ""
	}

	return string(mdJson)
}

func ExtractToken(authHeader string) string {
	// 將字串按空格分割為一個字串陣列
	parts := strings.Split(authHeader, " ")

	// 如果陣列長度不為 2，或者第一個元素不是 "Bearer"，則返回空字串
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	// 返回第二個元素，即 Token 字串
	return parts[1]
}

func ValidJson(str string) bool {
	return json.Valid([]byte(str))
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func ValidPageInfo(info *Common.PageInfoRequest) *Common.PageInfoRequest {
	// PrintObj(info, "ValidPageInfo")

	if info.CurrentPageNum < 0 {
		// PrintObj("reset CurrentPageNum")
		info.CurrentPageNum = 1
	}

	if !ValidNumber(int(info.PageItemCount), 1, 50) {
		// PrintObj("reset PageItemCount")
		info.PageItemCount = 10
	}

	return info
}

type ErrorType struct {
	Code        Common.ErrorCodes
	ReturnMsg   string
	InternalMsg string
}

func GetError(err ErrorType) string {
	if err.InternalMsg != "" {
		PrintObj(err.InternalMsg, "internalMsg")
	}

	model := &Common.ErrorReply{
		Code:    err.Code,
		Message: err.ReturnMsg,
	}

	return ToJson(model)
}

func GetFolderPath(imageType ArContent.ArContentImageType) (string, error) {
	var path string

	switch imageType {
	case ArContent.ArContentImageType_TEMP:
		path = GetEnv("AWS_S3_PATH_TEMP")
	case ArContent.ArContentImageType_THUMBNAIL:
		path = GetEnv("AWS_S3_PATH_THUMBNAIL")
	case ArContent.ArContentImageType_TEMPLATE_IMAGE:
		path = GetEnv("AWS_S3_PATH_TEMPLATE_IMAGE")
	case ArContent.ArContentImageType_VIEWER_IMAGE:
		path = GetEnv("AWS_S3_PATH_VIEWER_IMAGE")
	case ArContent.ArContentImageType_STATIC_IMAGE:
		path = GetEnv("AWS_S3_PATH_STATIC_IMAGE")
	default:
		return "", errors.New("image type not support")
	}

	return path, nil
}

func GetImagePath(imageName string, _type ArContent.ArContentImageType) string {
	folderPath, err := GetFolderPath(_type)
	if err != nil {
		PrintObj(err.Error(), "GetImagePath")
		return ""
	}

	return GetDomainAPI() + GetEnv("GIN_IMAGE_ROUTE") + folderPath + imageName
}

func GetDomainAPI() string {
	return GetEnv("DOMAIN_API")
}

func GetDomain() string {
	return GetEnv("DOMAIN")
}

func GetErrorGin(err ErrorType) gin.H {
	if err.InternalMsg != "" {
		PrintObj(err.InternalMsg, "internalMsg")
	}

	model := gin.H{
		"Code":    err.Code,
		"Message": err.ReturnMsg,
	}

	return model
}

func GetNewFileName(file *multipart.FileHeader) string {
	return uuid.New().String() + GetSuffix(file)
}

func GetSuffix(file *multipart.FileHeader) string {
	uploadFileNameWithSuffix := path.Base(file.Filename)
	return path.Ext(uploadFileNameWithSuffix)
}

func ParseJsonWithType[T any](str string) (res T, err error) {

	if len(str) == 0 {
		PrintObj("input invalid:", "ParseJsonWithType err")
		return res, errors.New("input invalid")
	}

	err = json.Unmarshal([]byte(str), &res)
	if err != nil {
		PrintObj(err.Error(), "ParseJsonWithType err")
		return res, err
	}

	return res, nil
}

func ParseJson(str string) interface{} {
	var result interface{}

	err := json.Unmarshal([]byte(str), &result)

	if err != nil {
		fmt.Println("parse json err")
		return ""
	}

	return result
}

func GetSqlLikeString(str string) string {
	return "%" + str + "%"
}

func PrintObj(obj interface{}, params ...string) {
	// print
	json, _ := json.Marshal(obj)
	key := ""

	if len(params) == 1 {
		if params[0] != "" {
			key = params[0]
			fmt.Println("=== " + key + " ===")
		}
	}

	if obj != "" {
		fmt.Println(string(json))
	}
}

func IsEmpty(str string) (res bool) {
	return str == ""
}

func GetEnv(str string) string {
	if !loadEnvSuccess {
		PrintObj("load env err pls check the .env", "GetEnv err")
	}

	val := os.Getenv(str)
	// PrintObj(val, "GetEnv "+str)
	return val
}

func ToBool(str string) (error, bool) {
	boolValue, err := strconv.ParseBool(str)
	if err != nil {
		PrintObj(err.Error(), "ToBool")
	}

	return nil, boolValue
}
