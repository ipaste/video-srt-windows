package app

import (
	"videosrt/app/aliyun"
	"videosrt/app/datacache"
	"videosrt/app/translate"
)

var RootDir string

var oss,translates,engine,setings *datacache.AppCache

//输出文件类型
const(
	OUTPUT_SRT = 1 //字幕SRT文件
	OUTPUT_STRING = 2 //普通文本
	OUTPUT_LRC = 3 //LRC文本
)

//输出文件编码
const(
	OUTPUT_ENCODE_UTF8 = 1 //文件编码 utf-8
	OUTPUT_ENCODE_UTF8_BOM = 2 //文件编码 utf-8 带 BOM
)

func init()  {
	RootDir = GetAppRootDir()
	if RootDir == "" {
		panic("应用根目录获取失败")
	}

	oss = datacache.NewAppCahce(RootDir , "oss")
	translates = datacache.NewAppCahce(RootDir , "translate")
	engine = datacache.NewAppCahce(RootDir , "engine")
	setings = datacache.NewAppCahce(RootDir , "setings")
}




//设置表单
type OperateFrom struct {
	EngineId int
	
	OutputSrt bool
	OutputLrc bool
	OutputTxt bool

	OutputType int
	OutputEncode int //输出文件编码
	SoundTrack int //输出音轨
}

func (from *OperateFrom) Init(setings *AppSetings)  {
	if setings.CurrentEngineId != 0 {
		from.EngineId = setings.CurrentEngineId
	}
	if setings.OutputType == 0 {
		from.OutputType = OUTPUT_SRT
		from.OutputSrt = true
	} else {
		from.OutputType = setings.OutputType
		if setings.OutputType == OUTPUT_SRT {
			from.OutputSrt = true
		}
		if setings.OutputType == OUTPUT_STRING {
			from.OutputTxt = true
		}
		if setings.OutputType == OUTPUT_LRC {
			from.OutputLrc = true
		}
	}
	if setings.OutputEncode == 0 {
		from.OutputEncode = OUTPUT_ENCODE_UTF8 //默认编码
	} else {
		from.OutputEncode = setings.OutputEncode
	}
	from.SoundTrack = setings.SoundTrack
}

func (from *OperateFrom) LoadOutputType(t int) {
	if OUTPUT_SRT != t {
		from.OutputSrt = false
	}
	if OUTPUT_LRC != t {
		from.OutputLrc = false
	}
	if OUTPUT_STRING != t {
		from.OutputTxt = false
	}

	if from.OutputSrt {
		from.OutputType = OUTPUT_SRT
	} else if from.OutputLrc {
		from.OutputType = OUTPUT_LRC
	} else if from.OutputTxt {
		from.OutputType = OUTPUT_STRING
	} else {
		from.OutputType = 0
	}
}


//引擎选项
type EngineSelects struct {
	Id   int
	Name string
}

//输出类型选项
type OutputSelects struct {
	Id   int
	Name string
}

//输出音轨类型选项
type SoundTrackSelects struct {
	Id   int
	Name string
}

//阿里云OSS - 缓存结构
type AliyunOssCache struct {
	aliyun.AliyunOss
}

//百度翻译账号认证类型选项
type BaiduAuthTypeSelects struct {
	Id   int
	Name string
}

//软件翻译接口 - 缓存结构
type TranslateCache struct {
	translate.BaiduTranslate //百度翻译

	AutoTranslation bool //中英互译（默认关闭）
	BilingualSubtitles bool //输出双语字幕（默认关闭）
}

//阿里云语音识别引擎 - 缓存结构
type AliyunEngineCache struct {
	aliyun.AliyunClound
	Id int //Id
	Alias string //别名
}

//阿里云语音识别引擎 - 列表缓存结构
type AliyunEngineListCache struct {
	Engine [] *AliyunEngineCache
}

//应用配置 - 缓存结构
type AppSetings struct {
	CurrentEngineId int //目前使用引擎Id
	MaxConcurrency int //任务最大处理并发数
	OutputType int //输出文件类型
	OutputEncode int //输出文件编码
	SrtFileDir string //Srt文件输出目录
	SoundTrack int //输出音轨

	CloseNewVersionMessage bool //关闭软件新版本提醒（默认开启）
}

//任务文件列表 - 结构
type TaskHandleFile struct {
	Files [] string
}


//获取 阿里云OSS 缓存数据
func GetCacheAliyunOssData() *AliyunOssCache {
	data := new(AliyunOssCache)
	vdata := oss.Get(data)
	if v, ok := vdata.(*AliyunOssCache); ok {
		return v
	}
	return data
}

//设置 阿里云OSS 缓存
func SetCacheAliyunOssData(data *AliyunOssCache) {
	oss.Set(data)
}


//获取 软件翻译接口 设置缓存
func GetCacheTranslateSettings() *TranslateCache {
	data := new(TranslateCache)
	vdata := translates.Get(data)
	if v, ok := vdata.(*TranslateCache); ok {
		return v
	}
	return data
}

//设置 软件翻译接口 缓存
func SetCacheTranslateSettings(data *TranslateCache)  {
	translates.Set(data)
}


//获取 阿里云语音识别引擎 缓存数据
func GetCacheAliyunEngineListData() *AliyunEngineListCache {
	data := new(AliyunEngineListCache)
	vdata := engine.Get(data)
	if v, ok := vdata.(*AliyunEngineListCache); ok {
		return v
	}
	return data
}

//设置 阿里云语音识别引擎 缓存
func SetCacheAliyunEngineListData(data *AliyunEngineListCache)  {
	engine.Set(data)
}

//根据id 删除 阿里云语音识别引擎 缓存数据
func RemoveCacheAliyunEngineData(id int) (bool) {
	var ok = false
	var newEngine = make([]*AliyunEngineCache , 0)
	origin := GetCacheAliyunEngineListData()

	total := len(origin.Engine)
	for i,engine := range origin.Engine	{
		if engine.Id == id {
			if i == (total - 1) {
				newEngine = origin.Engine[:i]
			} else {
				newEngine = append(origin.Engine[:i] , origin.Engine[i+1:]...)
			}
			ok = true
			break
		}
	}
	if ok {
		origin.Engine = newEngine
		//更新缓存数据
		SetCacheAliyunEngineListData(origin)
	}
	return ok
}


//获取 应用配置
func GetCacheAppSetingsData() *AppSetings {
	data := new(AppSetings)
	vdata := setings.Get(data)
	if v, ok := vdata.(*AppSetings); ok {
		return v
	}
	return data
}

//设置 应用配置
func SetCacheAppSetingsData(data *AppSetings)  {
	setings.Set(data)
}


//获取 引擎选项列表
func GetEngineOtionsSelects() []*EngineSelects {
	engines := make([]*EngineSelects , 0)
	//获取数据
	data := GetCacheAliyunEngineListData()

	for _,v := range data.Engine {
		engines = append(engines , &EngineSelects{
			Id:v.Id,
			Name:v.Alias,
		})
	}
	return engines
}


//根据选择引擎id 获取 引擎数据
func GetEngineById(id int) (*AliyunEngineCache , bool) {
	//获取数据
	data := GetCacheAliyunEngineListData()
	for _,v := range data.Engine {
		if id == v.Id {
			return v , true
		}
	}
	return nil , false
}


//获取 当前引擎id 下标
func GetCurrentIndex(data []*EngineSelects , id int) int {
	for index,v := range data {
		if v.Id == id {
			return index
		}
	}
	return -1
}


//获取 百度翻译账号认证类型
func GetBaiduTranslateAuthenTypeOptionsSelects() []*BaiduAuthTypeSelects {
	return []*BaiduAuthTypeSelects{
		&BaiduAuthTypeSelects{Id:translate.ACCOUNT_COMMON_AUTHEN , Name:"标准版"},
		&BaiduAuthTypeSelects{Id:translate.ACCOUNT_SENIOR_AUTHEN , Name:"高级版"},
	}
}


//获取 输出文件选项列表
func GetOutputOptionsSelects() []*OutputSelects {
	return []*OutputSelects{
		&OutputSelects{Id:OUTPUT_SRT , Name:"字幕文件"},
		&OutputSelects{Id:OUTPUT_STRING , Name:"普通文本"},
	}
}

//获取 输出文件编码选项列表
func GetOutputEncodeOptionsSelects() []*OutputSelects {
	return []*OutputSelects{
		&OutputSelects{Id:OUTPUT_ENCODE_UTF8 , Name:"UTF-8"},
		&OutputSelects{Id:OUTPUT_ENCODE_UTF8_BOM , Name:"UTF-8-BOM"},
	}
}

//获取 输出音轨选项列表
func GetSoundTrackSelects() []*SoundTrackSelects {
	return []*SoundTrackSelects{
		&SoundTrackSelects{Id:0 , Name:"全部"},
		&SoundTrackSelects{Id:1 , Name:"音轨一"},
		&SoundTrackSelects{Id:2 , Name:"音轨二"},
	}
}