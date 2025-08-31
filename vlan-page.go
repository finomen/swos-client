package swos_client

import "fmt"

/*
*
[

	{
		vid:0x0a,
		ivl:0x00,
		igmp:0x00,
		prt:[0x00,0x00,0x00,0x00,0x00,0x00]
	}

]
*/
type vlanStatus struct {
	Vid  string   `json:"vid"`
	Ivl  string   `json:"ivl"`
	Igmp string   `json:"igmp"`
	Prt  []string `json:"prt"`
}

type vlanChange struct {
	Vid  int
	Ivl  bool
	Igmp bool
	Prt  []int
}

type VlanPortMode int

const (
	VlanPortModeLeaveAsIs    VlanPortMode = 0
	VlanPortModeAlwaysStrip               = 1
	VlanPortModeAddIfMissing              = 2
	VlanPortModeNotAMember                = 3
)

type Vlan struct {
	Id                    int
	IndependentVlanLookup bool
	IgmpSnooping          bool
	PortMode              []VlanPortMode
}

type VlanPage struct {
	Vlans    []Vlan
	numPorts int
}

func (v *VlanPage) url() string {
	return "/vlan.b"
}

func (v *VlanPage) load(in []vlanStatus) error {
	v.Vlans = make([]Vlan, len(in))
	var err error
	for i, vl := range in {
		v.Vlans[i].Id, err = parseInt(vl.Vid)
		if err != nil {
			return err
		}
		v.Vlans[i].IndependentVlanLookup, err = parseBool(vl.Ivl)
		if err != nil {
			return err
		}
		v.Vlans[i].IgmpSnooping, err = parseBool(vl.Igmp)
		if err != nil {
			return err
		}
		if len(vl.Prt) != v.numPorts {
			return fmt.Errorf("invalid number of ports in vlan")
		}
		v.Vlans[i].PortMode = make([]VlanPortMode, v.numPorts)
		for j := 0; j < v.numPorts; j++ {
			pm, err := parseInt(vl.Prt[j])
			if err != nil {
				return err
			}
			v.Vlans[i].PortMode[j] = VlanPortMode(pm)
		}
	}
	return nil
}

func (v *VlanPage) store() []vlanChange {
	res := make([]vlanChange, len(v.Vlans))
	for i, vl := range v.Vlans {
		prt := make([]int, v.numPorts)
		for j := 0; j < v.numPorts; j++ {
			prt[j] = int(vl.PortMode[j])
		}
		res[i] = vlanChange{
			Vid:  vl.Id,
			Ivl:  vl.IndependentVlanLookup,
			Igmp: vl.IgmpSnooping,
			Prt:  prt,
		}
	}

	return res
}

func (v *VlanPage) AddVlan(id int) (*Vlan, error) {
	for _, vl := range v.Vlans {
		if vl.Id == id {
			return nil, fmt.Errorf("vlan already exists")
		}
	}
	v.Vlans = append(v.Vlans, Vlan{
		Id:       id,
		PortMode: make([]VlanPortMode, v.numPorts),
	})
	return &v.Vlans[len(v.Vlans)-1], nil
}

func (v *VlanPage) GetVlan(id int) (*Vlan, error) {
	for i, vl := range v.Vlans {
		if vl.Id == id {
			return &v.Vlans[i], nil
		}
	}
	return nil, fmt.Errorf("vlan %v does not exists", id)
}

func (v *VlanPage) DeleteVlan(id int) {
	wi := 0
	for i, vl := range v.Vlans {
		v.Vlans[wi] = v.Vlans[i]
		if vl.Id != id {
			wi++
		}
	}
	v.Vlans = v.Vlans[:wi]
}

var _ swOsPage[[]vlanStatus] = (*VlanPage)(nil)

var _ writebleSwOsPage[[]vlanChange, []vlanStatus] = (*VlanPage)(nil)
