package numatopo

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	"github.com/pkg/errors"

	memorystate "k8s.io/kubernetes/pkg/kubelet/cm/memorymanager/state"
	"volcano.sh/apis/pkg/apis/nodeinfo/v1alpha1"
	"volcano.sh/resource-exporter/pkg/args"
)

type memoryNumaInfo struct {
	// key is memory node id
	topologyInfo map[int]*v1alpha1.MemoryNode
}

func (info *memoryNumaInfo) Name() numaTopoName {
	return memoryNumaTopoName
}

func (info *memoryNumaInfo) Update(opt *args.Argument) NumaInfo {
	newInfo := NewMemoryNumaInfo()
	err := newInfo.updateInfoFromMemMngState(opt.MemoryMngState)
	if err != nil {
		klog.Infof("failed to update memory numa info, err: %v", err)
		return nil
	}

	if !reflect.DeepEqual(newInfo, info) {
		return newInfo
	}

	return nil
}

func (info *memoryNumaInfo) GetResourceInfoMap() v1alpha1.ResourceInfo {
	resourceInfo := v1alpha1.ResourceInfo{}
	memoryInfo := info.topologyInfo
	capa := uint64(0)
	allocatable := uint64(0)
	for _, mNode := range memoryInfo {
		table := mNode.MemoryMap[v1.ResourceMemory]
		if table == nil {
			continue
		}
		allocatable += table.Allocatable
		capa += table.TotalMemSize
	}

	resourceInfo.Allocatable = strconv.FormatUint(allocatable, 10)
	resourceInfo.Capacity = int(capa)
	return resourceInfo
}

// GetResTopoDetail return the cpu capability topology info
func (info *memoryNumaInfo) GetResTopoDetail() interface{} {
	return info.topologyInfo
}

func (info *memoryNumaInfo) updateInfoFromMemMngState(memMngState string) error {
	bytes, err := ioutil.ReadFile(memMngState)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read file memory_manager_state from %s", memMngState))
	}

	checkpoint := memorystate.NewMemoryManagerCheckpoint()
	err = checkpoint.UnmarshalCheckpoint(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal checkpoint for memory manager")
	}

	for node, s := range checkpoint.MachineState {
		if s == nil {
			info.topologyInfo[node] = nil
			continue
		}

		memoryNode := &v1alpha1.MemoryNode{}
		info.topologyInfo[node] = memoryNode
		memoryNode.NumberOfAssignments = s.NumberOfAssignments
		memoryNode.Cells = append([]int{}, s.Cells...)

		if s.MemoryMap == nil {
			continue
		}

		memoryNode.MemoryMap = map[v1.ResourceName]*v1alpha1.MemoryTable{}
		for memoryType, memoryTable := range s.MemoryMap {
			memoryNode.MemoryMap[memoryType] = &v1alpha1.MemoryTable{
				Allocatable:    memoryTable.Allocatable,
				Free:           memoryTable.Free,
				Reserved:       memoryTable.Reserved,
				SystemReserved: memoryTable.SystemReserved,
				TotalMemSize:   memoryTable.TotalMemSize,
			}
		}
	}

	return nil
}

func NewMemoryNumaInfo() *memoryNumaInfo {
	return &memoryNumaInfo{
		topologyInfo: map[int]*v1alpha1.MemoryNode{},
	}
}
