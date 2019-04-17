package main

import (
	"encoding/json"
	"fmt"
	"github.com/autom8ter/goproxy"
	"github.com/autom8ter/objectify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var (
	util    = objectify.Default()
	pconfig = &goproxy.ProxyConfig{
		Configs: []*goproxy.Config{
			{
				PathPrefix: "/twilio",
				TargetUrl:  "https://api.twilio.com/2010-04-01",
			},
		},
	}
	addr   string
	config string
)

var root = &cobra.Command{
	Use: "GoProxy",
	Long: fmt.Sprintf(`
              (                         
 (            )\ )                      
 )\ )        (()/( (            )  (    
(()/(     (   /(_)))(    (   ( /(  )\ ) 
 /(_))_   )\ (_)) (()\   )\  )\())(()/( 
(_)) __| ((_)| _ \ ((_) ((_)((_)\  )(_))
  | (_ |/ _ \|  _/| '_|/ _ \\ \ / | || |
   \___|\___/|_|  |_|  \___//_\_\  \_, |
                                   |__/

Current Config:
%s
`, viper.AllSettings()),
	Version: "v1.0",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the GoProxy server",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i, p := range pconfig.Configs {
			fmt.Printf("index: %v | registered proxy config: %s\n", i, util.MarshalJSON(p))
		}
		util.Entry().Debugf("starting GoProxy server: %s", addr)
		return http.ListenAndServe(addr, goproxy.New(pconfig))
	},
}

func main() {
	if err := root.Execute(); err != nil {
		util.Fatalln(util.WrapErr(err, "failed to run GoProxy"))
	}
}

func init() {
	root.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml", "relative path to file containing proxy configuration")
	root.PersistentFlags().StringVarP(&addr, "addr", "a", ":8080", "address to run server on")
	root.AddCommand(serveCmd)
	viper.SetConfigFile(config)
	_ = viper.BindPFlags(root.PersistentFlags())
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		util.Fatalln(util.WrapErr(err, "a config file is required to run GoProxy").Error())
	}
	if err := json.Unmarshal(util.MarshalJSON(viper.AllSettings()), pconfig); err != nil {
		util.Fatalln(err.Error())
	}
	if len(pconfig.Configs) == 0 {
		util.Fatalln(util.NewError("0 proxy configs registered from config file"))
	}
}
