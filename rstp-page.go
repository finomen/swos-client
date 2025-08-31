package swos_client

import "fmt"

/**
{
rpc:[0x00000000,0x00000000,0x00000000,0x00000000,0x00000000,0x00000000],
cst:[0x00000000,0x00000000,0x00000000,0x00000000,0x00000000,0x00000004],
ena:0x3d,
rstp:0x3f,
p2p:0x3f,
edge:0x22,
lrn:0x22,
fwd:0x22,
role:[0x03,0x00,0x03,0x03,0x03,0x03]
}
*/

type rstpStatus struct {
	Rpc  []string `json:"rpc"`
	Cst  []string `json:"cst"`
	Ena  string   `json:"ena"`
	Rstp string   `json:"rstp"`
	P2p  string   `json:"p2p"`
	Edge string   `json:"edge"`
	Lrn  string   `json:"lrn"`
	Fwd  string   `json:"fwd"`
	Role []string `json:"role"`
}

type rstpChange struct {
	Ena int
}

type RstpPage struct {
	RstpEnabled   []bool
	Mode          int
	Role          []int
	RootPathCoast []int
	Type          []int
	State         int

	numPorts int
}

func (r *RstpPage) url() string {
	return "/rstp.b"
}

func (r *RstpPage) load(in rstpStatus) error {
	var err error
	r.RstpEnabled, err = bitMaskToArray(in.Ena, r.numPorts)
	if err != nil {
		return err
	}
	r.Mode, err = parseInt(in.Rstp)
	if err != nil {
		return err
	}

	r.Role = make([]int, r.numPorts)
	r.RootPathCoast = make([]int, r.numPorts)
	r.Type = make([]int, r.numPorts)

	if len(in.Role) != r.numPorts {
		return fmt.Errorf("invalid number of ports in rstp status")
	}

	if len(in.Cst) != r.numPorts {
		return fmt.Errorf("invalid number of ports in rstp status")
	}

	if len(in.Rpc) != r.numPorts {
		return fmt.Errorf("invalid number of ports in rstp status")
	}

	for i := 0; i < r.numPorts; i++ {
		r.Role[i], err = parseInt(in.Role[i])
		if err != nil {
			return err
		}
		r.Type[i], err = parseInt(in.Cst[i])
		if err != nil {
			return err
		}
		r.RootPathCoast[i], err = parseInt(in.Rpc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RstpPage) store() rstpChange {
	return rstpChange{
		Ena: arrayToBitMask(r.RstpEnabled),
	}
}

var _ swOsPage[rstpStatus] = (*RstpPage)(nil)
var _ writebleSwOsPage[rstpChange, rstpStatus] = (*RstpPage)(nil)
