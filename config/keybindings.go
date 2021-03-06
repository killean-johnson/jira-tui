package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Email    string
	APIToken string
	JiraURL  string
	Project  []LayoutStruct
	Issue    []LayoutStruct
}

type LayoutStruct struct {
	View string
	Keys []Keybinding
}

type Keybinding struct {
	Name        string
	Key         string
	Description string
}

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func (kb *Config) LoadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, ".config", "jira-tui")
	configFile := filepath.Join(configPath, "config.json")

	// If the file doesn't exist, make the folder structure and a default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFile), 0770); err != nil {
			log.Fatal("Error creating config file")
		}

		createDefaultConfig(configFile)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&kb)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultConfig(configFile string) error {
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}

	defaultConfigString := []byte(`{
    "email": "",
    "apitoken": "",
    "jiraurl": "",
    "project": [
        {
            "view": "projectlist",
            "keys": [
                {
                    "name": "plcursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "plcursorup",
                    "key": "k",
                    "description": "Cursor Up"
                },
                {
                    "name": "plselect",
                    "key": "<ENTER>",
                    "description": "Select Board"
                },
                {
                    "name": "plquit",
                    "key": "q",
                    "description": "Quit"
                }
            ]
        },
        {
            "view": "messagebox",
            "keys": [
                {
                    "name": "mbclear",
                    "key": "<ENTER>",
                    "description": "Clear Messagebox"
                }
            ]
        }
    ],
    "issue": [
        {
            "view": "issuelist",
            "keys": [
                {
                    "name": "ilcursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "ilcursorup",
                    "key": "k",
                    "description": "Cursor Up"
                },
                {
                    "name": "ilselectissue",
                    "key": "<ENTER>",
                    "description": "Select Issue"
                },
                {
                    "name": "ileditdescription",
                    "key": "d",
                    "description": "Edit Description"
                },
                {
                    "name": "ileditstatus",
                    "key": "s",
                    "description": "Change Status"
                },
                {
                    "name": "ileditassignee",
                    "key": "a",
                    "description": "Change Assignee"
                },
                {
                    "name": "iladdissue",
                    "key": "A",
                    "description": "Add Issue"
                },
                {
                    "name": "ilquit",
                    "key": "q",
                    "description": "Quit"
                }
            ]
        },
        {
            "view": "issueview",
            "keys": [
                {
                    "name": "ivcursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "ivcursorup",
                    "key": "k",
                    "description": "Cursor Up"
                }
            ]
        },
        {
            "view": "editdesc",
            "keys": [
                {
                    "name": "edconfirm",
                    "key": "<C-s>",
                    "description": "Save Changes"
                },
                {
                    "name": "edcancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel Edit"
                }
            ]
        },
        {
            "view": "editstatus",
            "keys": [
                {
                    "name": "escursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "escursorup",
                    "key": "k",
                    "description": "Cursor Up"
                },
                {
                    "name": "esconfirm",
                    "key": "<ENTER>",
                    "description": "Set Status"
                },
                {
                    "name": "escancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel"
                }
            ]
        },
        {
            "view": "editassignee",
            "keys": [
                {
                    "name": "eacursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "eacursorup",
                    "key": "k",
                    "description": "Cursor Up"
                },
                {
                    "name": "eaconfirm",
                    "key": "<ENTER>",
                    "description": "Set Assignee"
                },
                {
                    "name": "eacancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel"
                }
            ]
        },
        {
            "view": "createisssummary",
            "keys": [
                {
                    "name": "ciscycle",
                    "key": "<TAB>",
                    "description": "Cycle Widgets"
                },
                {
                    "name": "cisconfirm",
                    "key": "<C-s>",
                    "description": "Create Issue"
                },
                {
                    "name": "ciscancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel"
                }
            ]
        },
        {
            "view": "createissassignee",
            "keys": [
                {
                    "name": "ciacursordown",
                    "key": "j",
                    "description": "Cursor Down"
                },
                {
                    "name": "ciacursorup",
                    "key": "k",
                    "description": "Cursor Up"
                },
                {
                    "name": "ciasetassignee",
                    "key": "<ENTER>",
                    "description": "Set Assignee"
                },
                {
                    "name": "ciacycle",
                    "key": "<TAB>",
                    "description": "Cycle Widgets"
                },
                {
                    "name": "ciaconfirm",
                    "key": "<C-s>",
                    "description": "Create Issue"
                },
                {
                    "name": "ciacancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel"
                }
            ]
        },
        {
            "view": "createissdesc",
            "keys": [
                {
                    "name": "cidcycle",
                    "key": "<TAB>",
                    "description": "Cycle Widgets"
                },
                {
                    "name": "cidconfirm",
                    "key": "<C-s>",
                    "description": "Create Issue"
                },
                {
                    "name": "cidcancel",
                    "key": "<ESCAPE>",
                    "description": "Cancel"
                }
            ]
        },
        {
            "view": "messagebox",
            "keys": [
                {
                    "name": "mbclear",
                    "key": "<ENTER>",
                    "description": "Clear Messagebox"
                }
            ]
        }
    ]
}`)
	_, err = file.Write(defaultConfigString)
	if err != nil {
		return err
	}
	return nil
}
