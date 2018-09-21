package errors

type BaseErrMsg struct {
	ErrMsg   string `json:"err_msg"`
	ErrMsgEn string `json:"err_msg_en"`
}

type ErrCode int

var (
	baseErrors = map[ErrCode]BaseErrMsg{}
)

const (
	ErrSuccess ErrCode = iota
	ErrUnknownError
	ErrUnstableNetwork
	ErrPermissionDeny
	ErrServiceUnderMaintaining
	ErrTooMuchRequest
	ErrServiceNotFound

	ErrNeedLogin    ErrCode = 40100
	ErrTokenExpired ErrCode = 41101
)

func init() {
	// err code definition rule:
	// 	0				System Reserved: Success
	//  1 - 99			System Reserved: Basic Error
	//  100 - 299		System Reserved: Api Gateway
	//  300 - 999  		System Reserved: Reserved

	//  1000 - 9999  	Business Logic
	//
	//  	1000 - 1999		Module: User
	//  	2000 - 2999		Module: Post
	//  	3000 - 3999 	Module: Demand/Quote
	//  	4000 - 4999     Module: HomePage
	//  	5000 - 5999     Module: DropShipping
	//  	6000 - 6999 	Module: WxHtml5
	//  	7000 - 7999     Module: WxMiniProgram
	//             7200 - 7299        Module:order_service_mp
	//		8000 - 9999 	Reserved
	//		8600 - 8699 	Module: DealerService
	//      8700 - 8799     Module: DsIndex

	//  41101 / 40100   System Reserved: Need Login

	baseErrors[0] = BaseErrMsg{ErrMsg: "成功", ErrMsgEn: "success"}
	
	baseErrors[1] = BaseErrMsg{ErrMsg: "未知错误", ErrMsgEn: "unknown error"}
	baseErrors[2] = BaseErrMsg{ErrMsg: "网络波动, 请重新尝试", ErrMsgEn: "unstable network connection, please retry"}
	baseErrors[3] = BaseErrMsg{ErrMsg: "权限不足", ErrMsgEn: "permission denied"}
	baseErrors[4] = BaseErrMsg{ErrMsg: "服务维护中, 请稍后", ErrMsgEn: "service maintaining, please wait"}
	baseErrors[5] = BaseErrMsg{ErrMsg: "访问量过大, 请稍后重试", ErrMsgEn: "router under flow control"}
	baseErrors[6] = BaseErrMsg{ErrMsg: "请求的服务不存在", ErrMsgEn: "service not found"}
	baseErrors[7] = BaseErrMsg{ErrMsg: "缺少参数: %s", ErrMsgEn: "lack of parameter"}
	baseErrors[8] = BaseErrMsg{ErrMsg: "缺少参数", ErrMsgEn: "lack of parameter"}
	baseErrors[9] = BaseErrMsg{ErrMsg: "参数非法: %s", ErrMsgEn: "invalid parameter"}

	baseErrors[30] = BaseErrMsg{ErrMsg: "'%s'长度不得大于%d", ErrMsgEn: "field too long"}

	baseErrors[80] = BaseErrMsg{ErrMsg: "网易云信接口调用失败: %s", ErrMsgEn: "netease request failed"}
	baseErrors[81] = BaseErrMsg{ErrMsg: "地址信息有误: %s", ErrMsgEn: "amap request failed"}
	baseErrors[82] = BaseErrMsg{ErrMsg: "短信验证码过期, 请重新获取", ErrMsgEn: "verify sms code expired"}

	// router operating error
	baseErrors[100] = BaseErrMsg{ErrMsg: "路由无效", ErrMsgEn: "routing table not exists"}
	baseErrors[120] = BaseErrMsg{ErrMsg: "服务已存在", ErrMsgEn: "service already exists"}
	baseErrors[121] = BaseErrMsg{ErrMsg: "服务不存在", ErrMsgEn: "service not exists"}
	baseErrors[122] = BaseErrMsg{ErrMsg: "终端服务已存在", ErrMsgEn: "endpoint already exists"}
	baseErrors[123] = BaseErrMsg{ErrMsg: "终端服务不存在", ErrMsgEn: "endpoint not exists"}

	baseErrors[124] = BaseErrMsg{ErrMsg: "路由规则使用中, 不可移除外层API", ErrMsgEn: "router is online, can not remove frontend api"}
	baseErrors[125] = BaseErrMsg{ErrMsg: "路由规则不存在", ErrMsgEn: "router not exist"}
	baseErrors[126] = BaseErrMsg{ErrMsg: "缺少外层API定义, 无法创建路由规则", ErrMsgEn: "api obj not completed"}
	baseErrors[127] = BaseErrMsg{ErrMsg: "缺少服务定义, 无法创建路由规则", ErrMsgEn: "service obj not completed"}
	baseErrors[128] = BaseErrMsg{ErrMsg: "路由规则已存在", ErrMsgEn: "router already exists"}
	baseErrors[129] = BaseErrMsg{ErrMsg: "绑定终端服务过程中出现错误, 无法创建路由规则", ErrMsgEn: "error raised when add endpoint to endpoint-table"}
	baseErrors[130] = BaseErrMsg{ErrMsg: "绑定的终端服务无一在线, 无法创建路由规则", ErrMsgEn: "no server online, can not create router"}
	baseErrors[131] = BaseErrMsg{ErrMsg: "无法找到路由规则", ErrMsgEn: "can not find router by name"}

	baseErrors[132] = BaseErrMsg{ErrMsg: "路由规则在线, 无法被注销", ErrMsgEn: "router is online, can not unregister it"}
	baseErrors[133] = BaseErrMsg{ErrMsg: "绑定的终端服务无一在线, 路由无法被设置为在线", ErrMsgEn: "all server are not online, this router should not be set to online"}
	baseErrors[134] = BaseErrMsg{ErrMsg: "绑定终端服务过程中出现错误, 路由无法被设置为在线", ErrMsgEn: "error raised when add server to server-table"}
	baseErrors[135] = BaseErrMsg{ErrMsg: "路由规则与在线路由表中的规则不符", ErrMsgEn: "data mapping error"}
	baseErrors[136] = BaseErrMsg{ErrMsg: "路由规则已在线", ErrMsgEn: "router already online"}
	baseErrors[137] = BaseErrMsg{ErrMsg: "无法找到服务", ErrMsgEn: "can not find service by name"}
	baseErrors[138] = BaseErrMsg{ErrMsg: "服务链接到某个已存在的路由规则, 无法被移除", ErrMsgEn: "service is linked to an online router, can not be removed"}
	baseErrors[139] = BaseErrMsg{ErrMsg: "无法找到终端服务", ErrMsgEn: "can not find endpoint by name"}
	baseErrors[140] = BaseErrMsg{ErrMsg: "轮询终端服务队列时发生异常: 节点数据错误", ErrMsgEn: "the Ring of endpoints contains an invalid value"}
	baseErrors[141] = BaseErrMsg{ErrMsg: "轮询终端服务队列时发生异常: 无一节点在线", ErrMsgEn: "all node of rings is not online"}
	baseErrors[142] = BaseErrMsg{ErrMsg: "外层API不存在", ErrMsgEn: "frontend api not found"}
	baseErrors[143] = BaseErrMsg{ErrMsg: "路由规则不在线", ErrMsgEn: "router is not online"}
	baseErrors[144] = BaseErrMsg{ErrMsg: "服务对应的终端节点无一在线", ErrMsgEn: "all node of rings is not online"}

	baseErrors[150] = BaseErrMsg{ErrMsg: "不正确的状态数值", ErrMsgEn: "invalid status code"}
	baseErrors[151] = BaseErrMsg{ErrMsg: "终端节点属性值不全", ErrMsgEn: "uncompleted attribute assigned"}

	baseErrors[160] = BaseErrMsg{ErrMsg: "缺少健康检查配置", ErrMsgEn: "health check not defined"}
	baseErrors[161] = BaseErrMsg{ErrMsg: "健康检查失败", ErrMsgEn: "health check failed"}

	baseErrors[170] = BaseErrMsg{ErrMsg: "无法生成Token令牌", ErrMsgEn: "can not generate token bucket"}

	baseErrors[310] = BaseErrMsg{ErrMsg: "context传递参数失败", ErrMsgEn: "can not get variable from context"}
	baseErrors[311] = BaseErrMsg{ErrMsg: "数据库连接类型有误", ErrMsgEn: "database connection type error"}

	// mini program
	baseErrors[7000] = BaseErrMsg{ErrMsg: "微信小程序登录凭证校验失败: %s", ErrMsgEn: "wx login check failed"}
	baseErrors[7001] = BaseErrMsg{ErrMsg: "缺少地理信息, 需要重新授权", ErrMsgEn: "need geo info"}
	baseErrors[7010] = BaseErrMsg{ErrMsg: "小程序无法从header上获取UserID", ErrMsgEn: "can not get user id from header"}
	baseErrors[7011] = BaseErrMsg{ErrMsg: "过期的小程序码: %s", ErrMsgEn: "invalid share code"}
	baseErrors[7012] = BaseErrMsg{ErrMsg: "用户与匠商绑定关系有误: %s", ErrMsgEn: "invalid relation between user and shop"}
	baseErrors[7013] = BaseErrMsg{ErrMsg: "用户%.2f公里内没有找到匠商", ErrMsgEn: "can not find shop owner"}
	baseErrors[7014] = BaseErrMsg{ErrMsg: "无法识别的ShareCode前缀", ErrMsgEn: "unrecognized share code prefix"}
	baseErrors[7050] = BaseErrMsg{ErrMsg: "小程序查询数据库失败: %s", ErrMsgEn: "database query failed"}
	baseErrors[7051] = BaseErrMsg{ErrMsg: "无法查询到用户信息: %s", ErrMsgEn: "can not find user info"}
	baseErrors[7052] = BaseErrMsg{ErrMsg: "无法查询到匠商信息: %s", ErrMsgEn: "can not find shop owner info"}
	baseErrors[7053] = BaseErrMsg{ErrMsg: "新增匠商信息失败: %s", ErrMsgEn: "error raised when insert new shop"}
	baseErrors[7060] = BaseErrMsg{ErrMsg: "分享图片正在生成中, 请稍后重试", ErrMsgEn: "share image still not available"}
	baseErrors[7070] = BaseErrMsg{ErrMsg: "gps信息上报过程中发生错误: %s", ErrMsgEn: "error raised when upload gps info"}
	baseErrors[7071] = BaseErrMsg{ErrMsg: "高德地图IP定位错误: %s", ErrMsgEn: "amap ip location error"}
	baseErrors[7072] = BaseErrMsg{ErrMsg: "visit record信息上报过程中发生错误: %s", ErrMsgEn: "error raised when upload record"}
	baseErrors[7100] = BaseErrMsg{ErrMsg: "需要重新登录", ErrMsgEn: "need login"}
	baseErrors[7101] = BaseErrMsg{ErrMsg: "Session不存在或过期, 需要重新登录", ErrMsgEn: "session expired, need login"}
	baseErrors[7102] = BaseErrMsg{ErrMsg: "需要授权获取微信用户信息", ErrMsgEn: "need wx user info"}
	baseErrors[7103] = BaseErrMsg{ErrMsg: "匠商已经存在, 无需注册", ErrMsgEn: "shop already exists"}
	baseErrors[7104] = BaseErrMsg{ErrMsg: "匠商缺少信息: %s", ErrMsgEn: "shop need more information"}
	baseErrors[7105] = BaseErrMsg{ErrMsg: "匠商审核失败: %s", ErrMsgEn: "shop audit failed"}

	// order_service_mp
	baseErrors[7200] = BaseErrMsg{ErrMsg: "获取产品信息失败:%s", ErrMsgEn: "get the product info fail"}
	baseErrors[7201] = BaseErrMsg{ErrMsg: "获取订单信息失败", ErrMsgEn: "get the order info fail"}
	baseErrors[7202] = BaseErrMsg{ErrMsg: "产品不存在购物车中", ErrMsgEn: "get the order info fail"}
	baseErrors[7203] = BaseErrMsg{ErrMsg: "产品区域信息出错:%s", ErrMsgEn: "get the order info fail"}
	

	// dealer service
	baseErrors[8600] = BaseErrMsg{ErrMsg: "服务暂时不可用, 请稍后再试", ErrMsgEn: "service down, please try again"}
	baseErrors[8601] = BaseErrMsg{ErrMsg: "生成UUID失败", ErrMsgEn: "uuid generate failed"}
	baseErrors[8602] = BaseErrMsg{ErrMsg: "条件类型非法", ErrMsgEn: "invalid condition value type"}
	baseErrors[8603] = BaseErrMsg{ErrMsg: "操作标识符无法识别", ErrMsgEn: "unknown operation symbol"}
	baseErrors[8604] = BaseErrMsg{ErrMsg: "允许类型标识符无法识别", ErrMsgEn: "unknown allow type symbol"}
	baseErrors[8605] = BaseErrMsg{ErrMsg: "db.Rows() panic: %s", ErrMsgEn: "db.Rows() panic"}
	baseErrors[8606] = BaseErrMsg{ErrMsg: "获取经销商详情失败", ErrMsgEn: "panic when get detail of dealer"}
	baseErrors[8607] = BaseErrMsg{ErrMsg: "获取经销商详情失败: %s", ErrMsgEn: "panic when get detail of dealer"}
	baseErrors[8608] = BaseErrMsg{ErrMsg: "续期时间不得小于当前时间", ErrMsgEn: "renewal time must larger than current time"}
	baseErrors[8609] = BaseErrMsg{ErrMsg: "动作标识符有误", ErrMsgEn: "action string invalid"}
	baseErrors[8610] = BaseErrMsg{ErrMsg: "时间格式有误: %s", ErrMsgEn: "time format invalid"}
	baseErrors[8611] = BaseErrMsg{ErrMsg: "操作失败", ErrMsgEn: "can not renewal"}
	baseErrors[8612] = BaseErrMsg{ErrMsg: "缺少原因", ErrMsgEn: "lack reason"}
	baseErrors[8613] = BaseErrMsg{ErrMsg: "缺少续期时间", ErrMsgEn: "lack renewal time"}
	baseErrors[8614] = BaseErrMsg{ErrMsg: "该经销商资格记录不存在", ErrMsgEn: "dealer record not exists"}
	baseErrors[8615] = BaseErrMsg{ErrMsg: "无法修改该经销商资格记录的状态, 原因: 需要[%s]状态", ErrMsgEn: "incorrect dealer status"}

	// ds index
	baseErrors[8700] = BaseErrMsg{ErrMsg: "搜索关键字不能为空", ErrMsgEn: "search keyword cannot ''"}

	baseErrors[40100] = BaseErrMsg{ErrMsg: "需要重新登录", ErrMsgEn: "need login"}
	baseErrors[41101] = BaseErrMsg{ErrMsg: "需要重新登录", ErrMsgEn: "need login"}
}
