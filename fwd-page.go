package swos_client

import "fmt"

/**
{
ir:[0x00000000,0x00000000,0x00000000,0x00000000,0x00000000,0x00000000],
or:[0x00000000,0x00000000,0x00000000,0x00000000,0x00000000,0x00000000],
fp1:0x3e,
fp2:0x3d,
fp3:0x3b,
fp4:0x37,
fp5:0x2f,
fp6:0x1f,
lck:0x00,
lckf:0x00,
imr:0x00,
omr:0x00,
mrto:0x01,
vlan:[0x01,0x01,0x01,0x01,0x01,0x01],
vlnh:[0x00,0x00,0x00,0x00,0x00,0x00],
vlni:[0x00,0x00,0x00,0x00,0x00,0x00],
fvid:0x00,
dvid:[0x0001,0x0001,0x0001,0x0001,0x0001,0x0001],
srt:[0x00,0x00,0x00,0x00,0x00,0x00],
suni:0x00}
*/

type VlanMode int

const (
	VlanModeDisabled VlanMode = 0
	VlanModeOptional          = 1
	VlanModeEnabled           = 2
	VlanModeStrict            = 2
)

type VlanHeader int

const (
	VlanHeaderLeaveAsIs  VlanHeader = 0
	VlanHeaderStrip                 = 1
	VlanHeaderAddMissing            = 2
)

type VlanReceive int

const (
	VlanReceiveAny      VlanReceive = 0
	VlanReceiveTagged               = 1
	VlanReceiveUntagged             = 2
)

type fwdStatus struct {
	Ir   []string `json:"ir"` // TODO: implement
	Or   []string `json:"or"`
	Fp1  string   `json:"fp1"`
	Fp2  string   `json:"fp2"`
	Fp3  string   `json:"fp3"`
	Fp4  string   `json:"fp4"`
	Fp5  string   `json:"fp5"`
	Fp6  string   `json:"fp6"`
	Lck  string   `json:"lck"`
	Lckf string   `json:"lckf"`
	Imr  string   `json:"imr"`
	Omr  string   `json:"omr"`
	Mrto string   `json:"mrto"`
	Vlan []string `json:"vlan"`
	Vlnh []string `json:"vlnh"`
	Vlni []string `json:"vlni"`
	Fvid string   `json:"fvid"`
	Dvid []string `json:"dvid"`
	Srt  []string `json:"srt"`  // TODO: implement
	Suni string   `json:"suni"` // TODO: implement
}

/**
{
	fp1:0x3e,
	fp2:0x3d,
	fp3:0x3b,
	fp4:0x33,
	fp5:0x2f,
	fp6:0x1f,
	lck:0x00,
	lckf:0x00,
	imr:0x00,
	omr:0x00,
	mrto:0x01,
	or:[0x00,0x00,0x00,0x00,0x00,0x00]
	vlan:[0x01,0x00,0x02,0x03,0x01,0x01],
	vlni:[0x00,0x00,0x00,0x00,0x00,0x00],
	dvid:[0x01,0x01,0x01,0x01,0x01,0x01],
	fvid:0x00,
	vlnh:[0x00,0x00,0x00,0x00,0x00,0x00]
}
*/

type fwdChange struct {
	Fp1  int
	Fp2  int
	Fp3  int
	Fp4  int
	Fp5  int
	Fp6  int
	Lck  int
	Lckf int
	Imr  int
	Omr  int
	Mrto int
	Or   []int
	Vlan []int
	Vlni []int
	Dvid []int
	Fvid int
	Vlnh []int
}

type FwdPage struct {
	PortLock      []bool
	LockOnFirst   []bool
	MirrorIngress []bool
	MirrorEgress  []bool
	MirrorTo      int
	ForwardTable  [][]bool
	VlanMode      []VlanMode
	ForceVlanId   []bool
	VlanHeader    []VlanHeader
	DefaultVlanId []int
	VlanReceive   []VlanReceive
	EgressRate    []int

	numPorts int
}

func (f *FwdPage) url() string {
	return "/fwd.b"
}

func (f *FwdPage) load(in fwdStatus) error {
	var err error
	f.PortLock, err = bitMaskToArray(in.Lck, f.numPorts)
	if err != nil {
		return err
	}
	f.LockOnFirst, err = bitMaskToArray(in.Lckf, f.numPorts)
	if err != nil {
		return err
	}
	f.MirrorIngress, err = bitMaskToArray(in.Imr, f.numPorts)
	if err != nil {
		return err
	}
	f.MirrorEgress, err = bitMaskToArray(in.Omr, f.numPorts)
	if err != nil {
		return err
	}
	f.MirrorTo, err = parseInt(in.Mrto)
	if err != nil {
		return err
	}
	f.ForwardTable = make([][]bool, f.numPorts)
	f.ForwardTable[0], err = bitMaskToArray(in.Fp1, f.numPorts)
	if err != nil {
		return err
	}
	f.ForwardTable[1], err = bitMaskToArray(in.Fp2, f.numPorts)
	if err != nil {
		return err
	}
	f.ForwardTable[2], err = bitMaskToArray(in.Fp3, f.numPorts)
	if err != nil {
		return err
	}
	f.ForwardTable[3], err = bitMaskToArray(in.Fp4, f.numPorts)
	if err != nil {
		return err
	}
	f.ForwardTable[4], err = bitMaskToArray(in.Fp5, f.numPorts)
	if err != nil {
		return err
	}
	f.ForwardTable[5], err = bitMaskToArray(in.Fp6, f.numPorts)
	if err != nil {
		return err
	}
	if len(in.Vlan) != f.numPorts {
		return fmt.Errorf("invalid vlan count")
	}
	f.VlanMode = make([]VlanMode, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		vm, err := parseInt(in.Vlan[i])
		if err != nil {
			return err
		}
		f.VlanMode[i] = VlanMode(vm)
	}
	f.ForceVlanId, err = bitMaskToArray(in.Fvid, f.numPorts)
	if err != nil {
		return err
	}

	if len(in.Vlnh) != f.numPorts {
		return fmt.Errorf("invalid vlan header count")
	}
	f.VlanHeader = make([]VlanHeader, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		vh, err := parseInt(in.Vlnh[i])
		if err != nil {
			return err
		}
		f.VlanHeader[i] = VlanHeader(vh)
	}

	if len(in.Dvid) != f.numPorts {
		return fmt.Errorf("invalid vlan id count")
	}
	f.DefaultVlanId = make([]int, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		f.DefaultVlanId[i], err = parseInt(in.Dvid[i])
		if err != nil {
			return err
		}
	}

	if len(in.Vlni) != f.numPorts {
		return fmt.Errorf("invalid vlan in count")
	}
	f.VlanReceive = make([]VlanReceive, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		vr, err := parseInt(in.Vlni[i])
		if err != nil {
			return err
		}
		f.VlanReceive[i] = VlanReceive(vr)
	}

	if len(in.Or) != f.numPorts {
		return fmt.Errorf("invalid or count")
	}
	f.EgressRate = make([]int, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		f.EgressRate[i], err = parseInt(in.Or[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FwdPage) store() fwdChange {
	vlanMode := make([]int, f.numPorts)
	vlanReceive := make([]int, f.numPorts)
	vlanHeader := make([]int, f.numPorts)
	for i := 0; i < f.numPorts; i++ {
		vlanMode[i] = int(f.VlanMode[i])
		vlanReceive[i] = int(f.VlanReceive[i])
		vlanHeader[i] = int(f.VlanHeader[i])
	}
	return fwdChange{
		Fp1:  arrayToBitMask(f.ForwardTable[0]),
		Fp2:  arrayToBitMask(f.ForwardTable[1]),
		Fp3:  arrayToBitMask(f.ForwardTable[2]),
		Fp4:  arrayToBitMask(f.ForwardTable[3]),
		Fp5:  arrayToBitMask(f.ForwardTable[4]),
		Fp6:  arrayToBitMask(f.ForwardTable[5]),
		Lck:  arrayToBitMask(f.PortLock),
		Lckf: arrayToBitMask(f.LockOnFirst),
		Imr:  arrayToBitMask(f.MirrorIngress),
		Omr:  arrayToBitMask(f.MirrorEgress),
		Mrto: f.MirrorTo,
		Or:   f.EgressRate,
		Vlan: vlanMode,
		Vlni: vlanReceive,
		Dvid: f.DefaultVlanId,
		Fvid: arrayToBitMask(f.ForceVlanId),
		Vlnh: vlanHeader,
	}
}

var _ swOsPage[fwdStatus] = (*FwdPage)(nil)

var _ writebleSwOsPage[fwdChange, fwdStatus] = (*FwdPage)(nil)
