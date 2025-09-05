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

type PortForward struct {
	PortLock      bool
	LockOnFirst   bool
	MirrorIngress bool
	MirrorEgress  bool
	ForwardTable  []bool
	VlanMode      VlanMode
	ForceVlanId   bool
	VlanHeader    VlanHeader
	DefaultVlanId int
	VlanReceive   VlanReceive
	EgressRate    int
}

type FwdPage struct {
	PortForward []PortForward
	MirrorTo    int

	numPorts int
}

func (f *FwdPage) url() string {
	return "/fwd.b"
}

func (f *FwdPage) load(in fwdStatus) error {
	var err error

	f.PortForward = make([]PortForward, f.numPorts)

	portLock, err := bitMaskToArray(in.Lck, f.numPorts)
	if err != nil {
		return err
	}

	lockOnFirst, err := bitMaskToArray(in.Lckf, f.numPorts)
	if err != nil {
		return err
	}

	mirrorIngress, err := bitMaskToArray(in.Imr, f.numPorts)
	if err != nil {
		return err
	}

	mirrorEgress, err := bitMaskToArray(in.Omr, f.numPorts)
	if err != nil {
		return err
	}

	forceVlanId, err := bitMaskToArray(in.Fvid, f.numPorts)
	if err != nil {
		return err
	}

	for i := 0; i < f.numPorts; i++ {
		f.PortForward[i].PortLock = portLock[i]
		f.PortForward[i].LockOnFirst = lockOnFirst[i]
		f.PortForward[i].MirrorIngress = mirrorIngress[i]
		f.PortForward[i].MirrorEgress = mirrorEgress[i]
		f.PortForward[i].ForceVlanId = forceVlanId[i]
	}

	f.MirrorTo, err = parseInt(in.Mrto)
	if err != nil {
		return err
	}

	f.PortForward[0].ForwardTable, err = bitMaskToArray(in.Fp1, f.numPorts)
	if err != nil {
		return err
	}
	f.PortForward[1].ForwardTable, err = bitMaskToArray(in.Fp2, f.numPorts)
	if err != nil {
		return err
	}
	f.PortForward[2].ForwardTable, err = bitMaskToArray(in.Fp3, f.numPorts)
	if err != nil {
		return err
	}
	f.PortForward[3].ForwardTable, err = bitMaskToArray(in.Fp4, f.numPorts)
	if err != nil {
		return err
	}
	f.PortForward[4].ForwardTable, err = bitMaskToArray(in.Fp5, f.numPorts)
	if err != nil {
		return err
	}
	f.PortForward[5].ForwardTable, err = bitMaskToArray(in.Fp6, f.numPorts)
	if err != nil {
		return err
	}
	if len(in.Vlan) != f.numPorts {
		return fmt.Errorf("invalid vlan count")
	}

	for i := 0; i < f.numPorts; i++ {
		vm, err := parseInt(in.Vlan[i])
		if err != nil {
			return err
		}
		f.PortForward[i].VlanMode = VlanMode(vm)
	}

	if len(in.Vlnh) != f.numPorts {
		return fmt.Errorf("invalid vlan header count")
	}

	for i := 0; i < f.numPorts; i++ {
		vh, err := parseInt(in.Vlnh[i])
		if err != nil {
			return err
		}
		f.PortForward[i].VlanHeader = VlanHeader(vh)
	}

	if len(in.Dvid) != f.numPorts {
		return fmt.Errorf("invalid vlan id count")
	}

	for i := 0; i < f.numPorts; i++ {
		f.PortForward[i].DefaultVlanId, err = parseInt(in.Dvid[i])
		if err != nil {
			return err
		}
	}

	if len(in.Vlni) != f.numPorts {
		return fmt.Errorf("invalid vlan in count")
	}

	for i := 0; i < f.numPorts; i++ {
		vr, err := parseInt(in.Vlni[i])
		if err != nil {
			return err
		}
		f.PortForward[i].VlanReceive = VlanReceive(vr)
	}

	if len(in.Or) != f.numPorts {
		return fmt.Errorf("invalid or count")
	}

	for i := 0; i < f.numPorts; i++ {
		f.PortForward[i].EgressRate, err = parseInt(in.Or[i])
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
	portLock := make([]bool, f.numPorts)
	lockOnFirst := make([]bool, f.numPorts)
	mirrorIngress := make([]bool, f.numPorts)
	mirrorEgress := make([]bool, f.numPorts)
	forceVlanId := make([]bool, f.numPorts)
	defaultVlanId := make([]int, f.numPorts)
	egressRate := make([]int, f.numPorts)

	for i := 0; i < f.numPorts; i++ {
		vlanMode[i] = int(f.PortForward[i].VlanMode)
		vlanReceive[i] = int(f.PortForward[i].VlanReceive)
		vlanHeader[i] = int(f.PortForward[i].VlanHeader)
		portLock[i] = f.PortForward[i].PortLock
		lockOnFirst[i] = f.PortForward[i].LockOnFirst
		mirrorIngress[i] = f.PortForward[i].MirrorIngress
		mirrorEgress[i] = f.PortForward[i].MirrorEgress
		forceVlanId[i] = f.PortForward[i].ForceVlanId
		defaultVlanId[i] = f.PortForward[i].DefaultVlanId
		egressRate[i] = f.PortForward[i].EgressRate
	}
	return fwdChange{
		Fp1:  arrayToBitMask(f.PortForward[0].ForwardTable),
		Fp2:  arrayToBitMask(f.PortForward[1].ForwardTable),
		Fp3:  arrayToBitMask(f.PortForward[2].ForwardTable),
		Fp4:  arrayToBitMask(f.PortForward[3].ForwardTable),
		Fp5:  arrayToBitMask(f.PortForward[4].ForwardTable),
		Fp6:  arrayToBitMask(f.PortForward[5].ForwardTable),
		Lck:  arrayToBitMask(portLock),
		Lckf: arrayToBitMask(lockOnFirst),
		Imr:  arrayToBitMask(mirrorIngress),
		Omr:  arrayToBitMask(mirrorEgress),
		Mrto: f.MirrorTo,
		Or:   egressRate,
		Vlan: vlanMode,
		Vlni: vlanReceive,
		Dvid: defaultVlanId,
		Fvid: arrayToBitMask(forceVlanId),
		Vlnh: vlanHeader,
	}
}

var _ swOsPage[fwdStatus] = (*FwdPage)(nil)

var _ writebleSwOsPage[fwdChange, fwdStatus] = (*FwdPage)(nil)
