package swos_client

/*
	{
	    en:0x3f,
	    lnk:0x20,
	    dpx:0x20,
	    fct:0x3f,
	    an:0x3f,
	    nm:['506f727431','506f727432','506f727433','506f727434','506f727435','534650'],
	    poe:[0x01,0x01,0x01,0x01,0x01,0x01],
	    prio:[0x00,0x00,0x01,0x02,0x03,0x00],
	    poes:[0x00,0x02,0x02,0x02,0x02,0x00],
	    spd:[0x03,0x03,0x03,0x03,0x03,0x02],
	    spdc:[0x00,0x00,0x00,0x00,0x00,0x00],
	    curr:[0x0000,0x0000,0x0000,0x0000,0x0000,0x0000],
	    pwr:[0x0000,0x0000,0x0000,0x0000,0x0000,0x0000],
	    dpxc:0x3f,
	}
*/
type linkStatus struct {
	Nm   []string `json:"nm"`
	En   string   `json:"en"`
	Lnk  string   `json:"lnk"`
	Spd  []string `json:"spd"`
	Dpx  string   `json:"dpx"`
	An   string   `json:"an"`
	Spdc []string `json:"spdc"`
	Dpxc string   `json:"dpxc"`
	Fct  string   `json:"fct"`
	Poe  []string `json:"poe"`
	Prio []string `json:"prio"`
	Poes []string `json:"poes"`
	Curr []string `json:"curr"`
	Pwr  []string `json:"pwr"`
}

//	{
//		nm:['506f727431','506f727432','506f727433','506f727434','506f727435','534650'],
//		spdc:[0x00,0x00,0x00,0x00,0x00,0x00],
//		prio:[0x00,0x00,0x00,0x00,0x00,0x00],
//		poe:[0x01,0x01,0x01,0x02,0x03,0x00],
//		en:0x63,an:0x63,
//		dpxc:0x63,
//		fct:0x63,
//	}
type linkChange struct {
	Nm   []string `json:"nm"`
	En   int      `json:"en"`
	An   int      `json:"an"`
	Spdc []int    `json:"spdc"`
	Dpxc int      `json:"dpxc"`
	Fct  int      `json:"fct"`
	Prio []int    `json:"prio"`
	Poe  []int    `json:"poe"`
}

type PoeMode int

const (
	Off   PoeMode = 0
	Auto          = 1
	On            = 2
	Calib         = 3
)

type PoeStatus int

const (
	Unavailable    PoeStatus = 0
	Disabled                 = 1
	WaitingForLoad           = 2
	Active                   = 3
)

type Link struct {
	Name            string
	Enabled         bool
	LinkUp          bool
	Duplex          bool
	DuplexControl   bool
	FlowControl     bool
	AutoNegotiation bool
	PoeMode         PoeMode
	PoePrio         int
	PoeStatus       PoeStatus
	SpeedControl    int
	Power           int
	Current         int
}

type LinkPage struct {
	Links []*Link
}

var _ swOsPage[linkStatus] = (*LinkPage)(nil)
var _ writebleSwOsPage[linkChange, linkStatus] = (*LinkPage)(nil)

func (link *LinkPage) url() string {
	return "/link.b"
}

func (link *LinkPage) load(status linkStatus) error {
	numPorts := len(status.Nm)

	en, err := bitMaskToArray(status.En, numPorts)
	if err != nil {
		return err
	}
	an, err := bitMaskToArray(status.An, numPorts)
	if err != nil {
		return err
	}
	lnk, err := bitMaskToArray(status.Lnk, numPorts)
	if err != nil {
		return err
	}
	dpx, err := bitMaskToArray(status.Dpx, numPorts)
	if err != nil {
		return err
	}
	dpxc, err := bitMaskToArray(status.Dpxc, numPorts)
	if err != nil {
		return err
	}
	fct, err := bitMaskToArray(status.Fct, numPorts)
	if err != nil {
		return err
	}

	for i := 0; i < numPorts; i++ {
		name, err := stringFromMikrotik(status.Nm[i])
		if err != nil {
			return err
		}

		poeMode, err := parseInt(status.Poe[i])
		if err != nil {
			return err
		}
		poePrio, err := parseInt(status.Prio[i])
		if err != nil {
			return err
		}
		poeStatus, err := parseInt(status.Poes[i])
		if err != nil {
			return err
		}
		speedControl, err := parseInt(status.Spdc[i])
		if err != nil {
			return err
		}
		power, err := parseInt(status.Pwr[i])
		if err != nil {
			return err
		}
		current, err := parseInt(status.Curr[i])
		if err != nil {
			return err
		}

		link.Links = append(link.Links, &Link{
			Name:            name,
			Enabled:         en[i],
			LinkUp:          lnk[i],
			Duplex:          dpx[i],
			DuplexControl:   dpxc[i],
			FlowControl:     fct[i],
			AutoNegotiation: an[i],
			PoeMode:         PoeMode(poeMode),
			PoePrio:         poePrio,
			PoeStatus:       PoeStatus(poeStatus),
			SpeedControl:    speedControl,
			Power:           power,
			Current:         current,
		})
	}

	return nil
}

func (link *LinkPage) store() linkChange {
	var names []string
	var en []bool
	var an = []bool{}
	var spdc []int
	var dpxc []bool
	var fct = []bool{}
	var poe = []int{}
	prio := []int{}

	for _, link := range link.Links {
		names = append(names, stringToMikrotik(link.Name))
		en = append(en, link.Enabled)
		an = append(an, link.AutoNegotiation)
		spdc = append(spdc, link.SpeedControl)
		dpxc = append(dpxc, link.DuplexControl)
		fct = append(fct, link.FlowControl)
		poe = append(poe, int(link.PoeMode))
		prio = append(prio, link.PoePrio)
	}

	return linkChange{
		Nm:   names,
		En:   arrayToBitMask(en),
		An:   arrayToBitMask(an),
		Spdc: spdc,
		Dpxc: arrayToBitMask(dpxc),
		Fct:  arrayToBitMask(fct),
		Prio: prio,
		Poe:  poe,
	}

}
