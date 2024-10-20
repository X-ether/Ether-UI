package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "unsafe"

	"x-ui/config"
	"x-ui/database"
	"x-ui/logger"
	"x-ui/sub"
	"x-ui/web"
	"x-ui/web/global"
	"x-ui/web/service"

	"github.com/op/go-logging"
)

func redirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果请求不是 HTTPS
		if r.TLS == nil {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		// 是 HTTPS 请求，继续处理
		next.ServeHTTP(w, r)
	})
}


func runWebServer() {
	log.Printf("%v %v", config.GetName(), config.GetVersion())

	switch config.GetLogLevel() {
	case config.Debug:
		logger.InitLogger(logging.DEBUG)
	case config.Info:
		logger.InitLogger(logging.INFO)
	case config.Warn:
		logger.InitLogger(logging.WARNING)
	case config.Error:
		logger.InitLogger(logging.ERROR)
	default:
		log.Fatal("unknown log level:", config.GetLogLevel())
	}

	err := database.InitDB(config.GetDBPath())
	if err != nil {
		log.Fatal(err)
	}

	var server *web.Server

	server = web.NewServer()
	global.SetWebServer(server)
	err = server.Start()
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()
	wrappedMux := redirectToHTTPS(mux)

	err = server.Start(wrappedMux) // 将包裹的 mux 传入 Start 函数
	if err != nil {
		log.Println(err)
		return
	}

	var subServer *sub.Server
	subServer = sub.NewServer()
	global.SetSubServer(subServer)

	err = subServer.Start()
	if err != nil {
		log.Println(err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	// Trap shutdown signals
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM)
	for {
		sig := <-sigCh

		switch sig {
		case syscall.SIGHUP:
			err := server.Stop()
			if err != nil {
				logger.Warning("stop server err:", err)
			}
			err = subServer.Stop()
			if err != nil {
				logger.Warning("stop server err:", err)
			}

			server = web.NewServer()
			global.SetWebServer(server)
			err = server.Start()
			if err != nil {
				log.Println(err)
				return
			}

			subServer = sub.NewServer()
			global.SetSubServer(subServer)

			err = subServer.Start()
			if err != nil {
				log.Println(err)
				return
			}
		default:
			server.Stop()
			subServer.Stop()
			return
		}
	}
}

func resetSetting() {
	err := database.InitDB(config.GetDBPath())
	if err != nil {
		fmt.Println(err)
		return
	}

	settingService := service.SettingService{}
	err = settingService.ResetSettings()
	if err != nil {
		fmt.Println("reset setting failed:", err)
	} else {
		fmt.Println("reset setting success")
	}
}

func showSetting(show bool) {
	if show {
		settingService := service.SettingService{}
		port, err := settingService.GetPort()
		if err != nil {
			fmt.Println("get current port failed, error info:", err)
		}

		webBasePath, err := settingService.GetBasePath()
		if err != nil {
			fmt.Println("get webBasePath failed, error info:", err)
		}

		userService := service.UserService{}
		userModel, err := userService.GetFirstUser()
		if err != nil {
			fmt.Println("get current user info failed, error info:", err)
		}

		username := userModel.Username
		userpasswd := userModel.Password
		if username == "" || userpasswd == "" {
			fmt.Println("current username or password is empty")
		}

		fmt.Println("current panel settings as follows:")
		fmt.Println("username:", username)
		fmt.Println("password:", userpasswd)
		fmt.Println("port:", port)
		if webBasePath != "" {
			fmt.Println("webBasePath:", webBasePath)
		} else {
			fmt.Println("webBasePath is not set")
		}
	}
}

func updateTgbotEnableSts(status bool) {
	settingService := service.SettingService{}
	currentTgSts, err := settingService.GetTgbotenabled()
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Infof("current enabletgbot status[%v],need update to status[%v]", currentTgSts, status)
	if currentTgSts != status {
		err := settingService.SetTgbotenabled(status)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			logger.Infof("SetTgbotenabled[%v] success", status)
		}
	}
}

func updateTgbotSetting(tgBotToken string, tgBotChatid string, tgBotRuntime string) {
	err := database.InitDB(config.GetDBPath())
	if err != nil {
		fmt.Println(err)
		return
	}

	settingService := service.SettingService{}

	if tgBotToken != "" {
		err := settingService.SetTgBotToken(tgBotToken)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			logger.Info("updateTgbotSetting tgBotToken success")
		}
	}

	if tgBotRuntime != "" {
		err := settingService.SetTgbotRuntime(tgBotRuntime)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			logger.Infof("updateTgbotSetting tgBotRuntime[%s] success", tgBotRuntime)
		}
	}

	if tgBotChatid != "" {
		err := settingService.SetTgBotChatId(tgBotChatid)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			logger.Info("updateTgbotSetting tgBotChatid success")
		}
	}
}

func updateSetting(port int, username string, password string, webBasePath string) {
	err := database.InitDB(config.GetDBPath())
	if err != nil {
		fmt.Println("Database initialization failed:", err)
		return
	}

	settingService := service.SettingService{}
	userService := service.UserService{}

	if port > 0 {
		err := settingService.SetPort(port)
		if err != nil {
			fmt.Println("Failed to set port:", err)
		} else {
			fmt.Printf("Port set successfully: %v\n", port)
		}
	}

	if username != "" || password != "" {
		err := userService.UpdateFirstUser(username, password)
		if err != nil {
			fmt.Println("Failed to update username and password:", err)
		} else {
			fmt.Println("Username and password updated successfully")
		}
	}

	if webBasePath != "" {
		err := settingService.SetBasePath(webBasePath)
		if err != nil {
			fmt.Println("Failed to set base URI path:", err)
		} else {
			fmt.Println("Base URI path set successfully")
		}
	}
}
func updateCert(publicKey string, privateKey string) {
	err := database.InitDB(config.GetDBPath())
	if err != nil {
		fmt.Println(err)
		return
	}

	if (privateKey != "" && publicKey != "") || (privateKey == "" && publicKey == "") {
		settingService := service.SettingService{}
		err = settingService.SetCertFile(publicKey)
		if err != nil {
			fmt.Println("set certificate public key failed:", err)
		} else {
			fmt.Println("set certificate public key success")
		}

		err = settingService.SetKeyFile(privateKey)
		if err != nil {
			fmt.Println("set certificate private key failed:", err)
		} else {
			fmt.Println("set certificate private key success")
		}
	} else {
		fmt.Println("both public and private key should be entered.")
	}
}

func migrateDb() {
	inboundService := service.InboundService{}

	err := database.InitDB(config.GetDBPath())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Start migrating database...")
	inboundService.MigrateDB()
	fmt.Println("Migration done!")
}

func main() {
	if len(os.Args) < 2 {
		runWebServer()
		return
	}

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version")

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	settingCmd := flag.NewFlagSet("setting", flag.ExitOnError)
	var port int
	var username string
	var password string
	var webBasePath string
	var webCertFile string
	var webKeyFile string
	var tgbottoken string
	var tgbotchatid string
	var enabletgbot bool
	var tgbotRuntime string
	var reset bool
	var show bool
	settingCmd.BoolVar(&reset, "reset", false, "reset all settings")
	settingCmd.BoolVar(&show, "show", false, "show current settings")
	settingCmd.IntVar(&port, "port", 0, "set panel port")
	settingCmd.StringVar(&username, "username", "", "set login username")
	settingCmd.StringVar(&password, "password", "", "set login password")
	settingCmd.StringVar(&webBasePath, "webBasePath", "", "Set base path for Panel")
	settingCmd.StringVar(&webCertFile, "webCert", "", "Set path to public key file for panel")
	settingCmd.StringVar(&webKeyFile, "webCertKey", "", "Set path to private key file for panel")
	settingCmd.StringVar(&tgbottoken, "tgbottoken", "", "Set token for Telegram bot")
	settingCmd.StringVar(&tgbotRuntime, "tgbotRuntime", "", "set telegram bot cron time")
	settingCmd.StringVar(&tgbotchatid, "tgbotchatid", "", "set telegram bot chat id")
	settingCmd.BoolVar(&enabletgbot, "enabletgbot", false, "enable telegram bot notify")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("    run            run web panel")
		fmt.Println("    migrate        migrate form other/old x-ui")
		fmt.Println("    setting        set settings")
	}

	flag.Parse()
	if showVersion {
		fmt.Println(config.GetVersion())
		return
	}

	switch os.Args[1] {
	case "run":
		err := runCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
		runWebServer()
	case "migrate":
		migrateDb()
	case "setting":
		err := settingCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
		if reset {
			resetSetting()
		} else {
			updateSetting(port, username, password, webBasePath)
		}
		if show {
			showSetting(show)
		}
		if (tgbottoken != "") || (tgbotchatid != "") || (tgbotRuntime != "") {
			updateTgbotSetting(tgbottoken, tgbotchatid, tgbotRuntime)
		}
		if enabletgbot {
			updateTgbotEnableSts(enabletgbot)
		}
	case "cert":
		err := settingCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
		if reset {
			updateCert("", "")
		} else {
			updateCert(webCertFile, webKeyFile)
		}

	default:
		fmt.Println("Invalid subcommands")
		fmt.Println()
		runCmd.Usage()
		fmt.Println()
		settingCmd.Usage()
	}
}
