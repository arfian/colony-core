package colonycore

import (
	"errors"
	"fmt"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/sshclient"
	"golang.org/x/crypto/ssh"
	"strings"
)

type Server struct {
	orm.ModelBase
	ID            string           `json:"_id", bson:"_id"`
	OS            string           `json:"os", bson:"os"`
	AppPath       string           `json:"appPath",bson:"appPath"`
	DataPath      string           `json:"dataPath",bson:"dataPath"`
	Host          string           `json:"host", bson:"host"`
	ServerType    string           `json:"serverType", bson:"serverType"`
	SSHType       string           `json:"sshtype", bson:"sshtype"`
	SSHFile       string           `json:"sshfile", bson:"sshfile"`
	SSHUser       string           `json:"sshuser", bson:"sshuser"`
	SSHPass       string           `json:"sshpass", bson:"sshpass"`
	CmdExtract    string           `json:"cmdextract", bson:"cmdextract"`
	CmdNewFile    string           `json:"cmdnewfile", bson:"cmdnewfile"`
	CmdCopy       string           `json:"cmdcopy", bson:"cmdcopy"`
	CmdMkDir      string           `json:"cmdmkdir", bson:"cmdmkdir"`
	HostAlias     []*HostAlias     `json:"hostAlias", bson:"hostAlias"`
	InstalledLang []*InstalledLang `json:"installedLang", bson:"installedLang"`
}

type HostAlias struct {
	IP       string `json:"ip", bson:"ip"`
	HostName string `json:"hostName", bson:"hostName"`
}

func (s *Server) TableName() string {
	return "servers"
}

func (s *Server) RecordID() interface{} {
	return s.ID
}

func (s *Server) Connect() (sshclient.SshSetting, *ssh.Client, error) {
	client := sshclient.SshSetting{}
	client.SSHHost = s.Host

	if s.SSHType == "File" {
		client.SSHAuthType = sshclient.SSHAuthType_Certificate
		client.SSHKeyLocation = s.SSHFile
	} else {
		client.SSHAuthType = sshclient.SSHAuthType_Password
		client.SSHUser = s.SSHUser
		client.SSHPassword = s.SSHPass
	}

	theClient, err := client.Connect()

	return client, theClient, err
}

func (s *Server) IsCommandExists(cmd string) (bool, string, error) {
	setting, _, err := s.Connect()
	if err != nil {
		return false, "", err
	}

	res, err := setting.RunCommandSshAsMap(fmt.Sprintf("which %s", cmd))
	if err != nil {
		return false, "", err
	}

	output := strings.TrimSpace(res[0].Output)
	if output == "" {
		return false, output, errors.New("command not found")
	}

	if resOutput, err := setting.RunCommandSshAsMap(cmd); err == nil {
		output = resOutput[0].Output
	}

	return true, output, nil
}

func (s *Server) Ping() (bool, error) {
	if _, _, err := s.Connect(); err != nil {
		return false, err
	}

	return true, nil
}

type ServerLanguage struct {
	ServerID   string
	ServerOS   string
	ServerHost string
	Languages  []*InstalledLang
}

type InstalledLang struct {
	Lang        string
	Version     string
	IsInstalled bool
}
