package config

import (
	"encoding/json"
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/spf13/viper"
)

type ConfigStruct struct {
    Layouts map[string]LayoutStruct
}

type LayoutStruct struct {
    Views map[string]ViewStruct
}

type ViewStruct struct {
    Bindings map[string][]Keybinding
}

type Keybindings struct {
    binding []map[string]Keybinding
}

type Keybinding struct {
    key gocui.Key
    description string
}

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func (kb *Keybindings) LoadKeybindings() error {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    viper.SetConfigType("json")

    viper.AutomaticEnv()

    err := viper.ReadInConfig()
    if err != nil {
        return err
    }

    MarshalPrint(viper.GetStringMap("keybindings"))

    var keybs ConfigStruct
    viper.Unmarshal(&keybs)

    MarshalPrint(keybs)

    return nil
}

// thing := map[string]interface {}{
// 	"board":map[string]interface {}{
// 		"":[]interface {}{}, 
// 		"boardlist":[]interface {}{
//             map[string]interface {}{
//                 "description":"Cursor Down",
//                 "key":"j"
//             },
//             map[string]interface {}{
//                 "description":"Cursor Up",
//                 "key":"k"
//             }
// 	    }
//     }, 
//     "issue":map[string]interface {}{
//         "":[]interface {}{},
//         "editdesc":[]interface {}{},
//         "issuelist":[]interface {}{
//             map[string]interface {}{
//                 "description":"Cursor Down", 
//                 "key":"j"
//             }, map[string]interface {}{
//                 "description":"Cursor Up", 
//                 "key":"k"
//             }
//         }, 
//         "issueview":[]interface {}{}
//     }
// }

