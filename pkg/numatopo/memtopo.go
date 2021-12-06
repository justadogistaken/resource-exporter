package numatopo

import (
	"io/ioutil"
	"reflect"

	"volcano.sh/apis/pkg/apis/nodeinfo/v1alpha1"
	"volcano.sh/resource-exporter/pkg/args"

	"k8s.io/klog"
	memstate "k8s.io/kubernetes/pkg/kubelet/cm/memorymanager/state"
)

type MemoryNumaInfo struct {
	MemoryDetail map[int]v1alpha1.MemoryNode
}

func (m *MemoryNumaInfo) Name() string {
	return TopoMemory
}

func (m *MemoryNumaInfo) Update(opt *args.Argument) NumaInfo {
	newInfo := NewMemoryNumaInfo()
	newInfo.updateMemoryDetail(opt.MemoryMngState)
	if !reflect.DeepEqual(newInfo, m) {
		return newInfo
	}
	return nil
}

func (m *MemoryNumaInfo) GetResourceInfoMap() v1alpha1.ResourceInfo {
	return v1alpha1.ResourceInfo{}
}

func (m *MemoryNumaInfo) GetResTopoDetail() interface{} {
	return m.MemoryDetail
}

func (m *MemoryNumaInfo) updateMemoryDetail(memMngStatePath string) {
	mDetail := getMachineState(memMngStatePath)
	for k, v := range mDetail {
		if v == nil {
			continue
		}

		mn := v1alpha1.MemoryNode{}
		mn.Cells = v.Cells
		for rn, mt := range v.MemoryMap {
			mn.MemoryMap[string(rn)] = v1alpha1.MemoryTable{
				TotalMemSize:   mt.TotalMemSize,
				SystemReserved: mt.SystemReserved,
				Allocatable:    mt.Allocatable,
				Reserved:       mt.Reserved,
				Free:           mt.Free,
			}
		}

		m.MemoryDetail[k] = mn
	}
}

func getMachineState(memMngStatePath string) map[int]*memstate.NUMANodeState {
	data, err := ioutil.ReadFile(memMngStatePath)
	if err != nil {
		klog.Errorf("failed to read memory_manager_state, err: %v", err)
		return nil
	}

	checkpoint := memstate.NewMemoryManagerCheckpoint()
	checkpoint.UnmarshalCheckpoint(data)

	return checkpoint.MachineState
}

func NewMemoryNumaInfo() *MemoryNumaInfo {
	return &MemoryNumaInfo{
		MemoryDetail: make(map[int]v1alpha1.MemoryNode),
	}
}
