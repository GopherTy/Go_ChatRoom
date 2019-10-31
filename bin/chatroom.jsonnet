local Millisecond = 1;
local Second = 1000 * Millisecond;
local Minute = 60 * Second;
local Hour = 60 * Minute;
local Day = 24 * Hour;
{
	GRPC:{
		Addr:":6000",
		// x509 if empty use h2c
		CertFile:"test.pem",
		KeyFile:"test.key",
	},
	Logger:{
		// zap http
		//HTTP:"localhost:20000",
		// log name
		//Filename:"logs/chatroom.log",
		// MB
		MaxSize:    100, 
		// number of files
		MaxBackups: 3,
		// day
		MaxAge:     28,
		// level : debug info warn error dpanic panic fatal
		Level :"debug",
		// 是否要 輸出 代碼位置
        Caller:true,
	},
}