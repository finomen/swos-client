package swos_client

import "net"

//	{
//		mac:'f41e575c867a',
//		sid:'4848363041394e36565650',
//		id:'4d696b726f54696b',
//		ver:'322e3138',
//		brd:'4353533130362d31472d34502d3153',
//		rmac:'f41e575c867a',
//		upt:0x0001e2f2,
//		ip:0xfe0110ac,
//		bld:0x6864f9b7,
//		wdt:0x01,
//		dsc:0x01,
//		pdsc:0x3f,
//		ivl:0x00,
//		alla:0x00000000,
//		allm:0x00,
//		allp:0x3f,
//		avln:0x0000,
//		prio:0x8000,
//		cost:0x00,
//		frmc:0x00,
//		rpr:0x8000,
//		igmp:0x00,
//		igmq:0x01,
//		sip:0x0158a8c0,
//		iptp:0x00,
//		volt:0x00ed,
//		temp:0x0000001e,
//		lcbl:0x00,
//		upgr:0x00,
//		igfl:0x00,
//		igve:0x01
//	}
type sysStatus struct {
	Upt  string `json:"upt"`
	Ip   string `json:"ip"`
	Mac  string `json:"mac"`
	Sid  string `json:"sid"`
	Id   string `json:"id"`
	Ver  string `json:"ver"`
	Brd  string `json:"brd"`
	Bld  string `json:"bld"`
	Wdt  string `json:"wdt"`
	Dsc  string `json:"dsc"`
	Pdsc string `json:"pdsc"`
	Ivl  string `json:"ivl"`
	Alla string `json:"alla"`
	Allm string `json:"allm"`
	Allp string `json:"allp"`
	Avln string `json:"avln"`
	Prio string `json:"prio"`
	Cost string `json:"cost"`
	Frmc string `json:"frmc"`
	Rpr  string `json:"rpr"`
	Rmac string `json:"rmac"`
	Igmp string `json:"igmp"`
	Igmq string `json:"igmq"`
	Sip  string `json:"sip"`
	Iptp string `json:"iptp"`
	Volt string `json:"volt"`
	Temp string `json:"temp"`
	Lcbl string `json:"lcbl"`
	Upgr string `json:"upgr"`
	Igfl string `json:"igfl"`
	Igve string `json:"igve"`
}

/*
	{
		// id:'4d696b726f54696b32',
		//iptp:0x00,
		//sip:0x0158a8c0,
		//alla:0x00,
		//allm:0x00,
		//allp:0x3f,
		//avln:0x00,
		//ivl:0x00,
		//igmp:0x00,
		//igmq:0x00,
		//igfl:0x00,
		//igve:0x01,
		//pdsc:0x3f,
		//lcbl:0x00
		// prio: 0x8000,
		// cost: 0x00,
		// frmc: 0x00,
	}
*/
type sysChange struct {
	Iptp int
	Sip  int
	Id   string
	Alla int
	Allm int
	Allp int
	Avln int
	Ivl  int
	Igmp int
	Igmq int
	Igfl int
	Igve int
	Pdsc int
	Lcbl int
}

type SysPage struct {
	Mac                       net.HardwareAddr
	SerialNumber              string
	Identity                  string
	Version                   string
	BoardName                 string
	RootBridgeMac             net.HardwareAddr
	Uptime                    int64
	Ip                        net.IP
	Build                     int // FIXME: wtf?
	Dsc                       int // FIXME: wtf?
	Wdt                       int // FIXME: wtf?
	MikrotikDiscoveryProtocol []bool
	IndependentVlanLookup     bool
	AllowFrom                 net.IP
	Allm                      int // WTF
	AllowFromPorts            []bool
	AllowFromVlan             int
	IgmpSnooping              bool
	IgmpQuerier               bool
	LongPoeCable              bool
	IgmpFastLeave             []bool
	IgmpVersion               int //TODO: enum
	Voltage                   int
	Temperature               int
	BridgePriority            int
	PortCostMode              int //TODO: enum
	ForwardReservedMulticast  bool
	AddressAquisition         int // TODO: enum 0 dhcp+fb, 1 static, 2 dhcp only
	StaticIpAddress           net.IP

	numPorts int
}

func (s *SysPage) store() sysChange {
	return sysChange{
		Iptp: s.AddressAquisition,
		Sip:  ipToMikrotik(s.StaticIpAddress),
		Id:   stringToMikrotik(s.Identity),
		Alla: ipToMikrotik(s.AllowFrom),
		Allm: s.Allm,
		Allp: arrayToBitMask(s.AllowFromPorts),
		Avln: s.AllowFromVlan,
		Ivl:  bootToInt(s.IndependentVlanLookup),
		Igmp: bootToInt(s.IgmpSnooping),
		Igmq: bootToInt(s.IgmpQuerier),
		Igfl: arrayToBitMask(s.IgmpFastLeave),
		Igve: s.IgmpVersion,
		Pdsc: arrayToBitMask(s.MikrotikDiscoveryProtocol),
		Lcbl: bootToInt(s.LongPoeCable),
	}
}

func (s *SysPage) url() string {
	return "/sys.b"
}

func (s *SysPage) load(sys sysStatus) error {
	mac, err := macFromMikrotik(sys.Mac)
	if err != nil {
		return err
	}
	s.Mac = mac

	str, err := stringFromMikrotik(sys.Sid)
	if err != nil {
		return err
	}
	s.SerialNumber = str

	str, err = stringFromMikrotik(sys.Id)
	if err != nil {
		return err
	}
	s.Identity = str

	str, err = stringFromMikrotik(sys.Ver)
	if err != nil {
		return err
	}
	s.Version = str

	str, err = stringFromMikrotik(sys.Brd)
	if err != nil {
		return err
	}
	s.BoardName = str

	mac, err = macFromMikrotik(sys.Rmac)
	if err != nil {
		return err
	}
	s.RootBridgeMac = mac

	i64, err := parseInt64(sys.Upt)
	if err != nil {
		return err
	}
	s.Uptime = i64

	ip, err := ipFromMikrotik(sys.Ip)
	if err != nil {
		return err
	}
	s.Ip = ip

	i, err := parseInt(sys.Bld)
	if err != nil {
		return err
	}
	s.Build = i

	i, err = parseInt(sys.Wdt)
	if err != nil {
		return err
	}
	s.Wdt = i

	i, err = parseInt(sys.Dsc)
	if err != nil {
		return err
	}
	s.Dsc = i

	fl, err := bitMaskToArray(sys.Pdsc, s.numPorts)
	if err != nil {
		return err
	}
	s.MikrotikDiscoveryProtocol = fl

	i, err = parseInt(sys.Ivl)
	if err != nil {
		return err
	}
	s.IndependentVlanLookup = i != 0

	ip, err = ipFromMikrotik(sys.Alla)
	if err != nil {
		return err
	}
	s.AllowFrom = ip

	i, err = parseInt(sys.Allm)
	if err != nil {
		return err
	}
	s.Allm = i

	fl, err = bitMaskToArray(sys.Allp, s.numPorts)
	if err != nil {
		return err
	}
	s.AllowFromPorts = fl

	i, err = parseInt(sys.Avln)
	if err != nil {
		return err
	}
	s.AllowFromVlan = i

	i, err = parseInt(sys.Igmp)
	if err != nil {
		return err
	}
	s.IgmpSnooping = i != 0

	i, err = parseInt(sys.Igmq)
	if err != nil {
		return err
	}
	s.IgmpQuerier = i != 0

	i, err = parseInt(sys.Volt)
	if err != nil {
		return err
	}
	s.Voltage = i

	i, err = parseInt(sys.Temp)
	if err != nil {
		return err
	}
	s.Temperature = i

	i, err = parseInt(sys.Lcbl)
	if err != nil {
		return err
	}
	s.LongPoeCable = i != 0

	fl, err = bitMaskToArray(sys.Igfl, s.numPorts)
	if err != nil {
		return err
	}
	s.IgmpFastLeave = fl

	i, err = parseInt(sys.Igve)
	if err != nil {
		return err
	}
	s.IgmpVersion = i

	i, err = parseInt(sys.Prio)
	if err != nil {
		return err
	}
	s.BridgePriority = i

	i, err = parseInt(sys.Cost)
	if err != nil {
		return err
	}
	s.PortCostMode = i

	i, err = parseInt(sys.Frmc)
	if err != nil {
		return err
	}
	s.ForwardReservedMulticast = i != 0

	i, err = parseInt(sys.Iptp)
	if err != nil {
		return err
	}
	s.AddressAquisition = i

	ip, err = ipFromMikrotik(sys.Sip)
	if err != nil {
		return err
	}
	s.StaticIpAddress = ip

	// TODO: Looks like rpr is the same as prio. Should it be different field?
	//TODO: pare upgrade field?

	return nil
}

var _ swOsPage[sysStatus] = (*SysPage)(nil)
var _ writebleSwOsPage[sysChange, sysStatus] = (*SysPage)(nil)
