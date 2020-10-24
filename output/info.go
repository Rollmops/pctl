package output

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

type Info struct {
	Name           string
	Comment        string
	ConfigCommand  []string
	RunningCommand []string
	IsRunning      bool
	DirtyCommand   bool
	DirtyMd5Hashes []string
	RunningInfo    *gopsutil.Process
	Dirty          bool
}

func CreateInfos(processConfigs []*config.ProcessConfig) ([]*Info, error) {
	var infos []*Info
	for _, processConfig := range processConfigs {
		infoEntry := &Info{
			Name:           processConfig.Name,
			ConfigCommand:  processConfig.Command,
			RunningCommand: processConfig.Command,
		}

		runningInfo, err := process.FindRunningInfo(processConfig.Name)
		if err != nil {
			return nil, err
		}
		if runningInfo == nil {
			infoEntry.IsRunning = false
			infoEntry.ConfigCommand = processConfig.Command
		} else {

			dirtyHashes, err := collectDirtyHashes(&processConfig.Command, runningInfo)
			if err != nil {
				return nil, err
			}
			infoEntry.DirtyMd5Hashes = *dirtyHashes

			infoEntry.IsRunning = true
			infoEntry.ConfigCommand = runningInfo.Config.Command
			infoEntry.DirtyCommand = !common.CompareStringSlices(infoEntry.ConfigCommand, processConfig.Command)
			infoEntry.Dirty = infoEntry.DirtyCommand || len(infoEntry.DirtyMd5Hashes) > 0

			p := process.Process{Config: processConfig}
			err = p.SynchronizeWithPid(runningInfo.Pid)
			if err != nil {
				return nil, err
			}
			infoEntry.RunningInfo, err = p.GetPsutilProcess()
			if err != nil {
				return nil, err
			}
			infoEntry.RunningCommand, err = infoEntry.RunningInfo.CmdlineSlice()
			if err != nil {
				return nil, err
			}
		}

		infos = append(infos, infoEntry)
	}
	return infos, nil
}

func collectDirtyHashes(command *[]string, runningInfo *process.RunningEnvironInfo) (*[]string, error) {
	logrus.Tracef("Collecting dirty file hashes from command '%v'", *command)
	var returnDirtyHashes []string
	md5hashes, err := common.CreateFileHashesFromCommand(*command)
	if err != nil {
		return nil, err
	}
	for arg, hash := range *md5hashes {
		runningHash := runningInfo.Md5Hashes[arg]
		if runningHash != hash {
			logrus.Tracef("Found dirty hash for arg '%s': %s != %s", arg, runningHash, hash)
			returnDirtyHashes = append(returnDirtyHashes, arg)
		}
	}
	return &returnDirtyHashes, nil
}
