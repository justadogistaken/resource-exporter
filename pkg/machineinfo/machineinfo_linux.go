/*
Copyright 2021 The Volcano Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machineinfo

import (
	"os"

	"github.com/google/cadvisor/fs"
	"github.com/google/cadvisor/machine"
	"github.com/google/cadvisor/utils/sysfs"
	"k8s.io/klog"
)

func init() {
	fsContext := fs.Context{}
	sysFs := sysfs.NewRealSysFs()

	fsInfo, err := fs.NewFsInfo(fsContext)
	if err != nil {
		klog.Fatalf("failed to initiate FsInfo, err: %v", err)
		return
	}

	inHostNamespace := false
	if _, err = os.Stat("/rootfs/proc"); os.IsNotExist(err) {
		inHostNamespace = true
	}

	machineInfo, err := machine.Info(sysFs, fsInfo, inHostNamespace)
	if err != nil {
		klog.Fatalf("failed to initiate machine info, err: %v", err)
		return
	}
	gMachineInfo = machineInfo
}

func LoadMachineInfo() {

}
