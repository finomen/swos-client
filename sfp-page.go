package swos_client

// {
//vnd:'4d696b726f54696b2020202020202020',
//pnr:'58532b44413030303120202020202020',
//rev:'312e3020',
//ser:'53323530343036363434343932202020',
//dat:'32352d30342d3131',
//typ:'316d20636f70706572',
//wln:0x00000000,
//tmp:0xffffff80,
//vcc:0x0000,
//tbs:0x0000,
//tpw:0x0000,
//rpw:0x0000}

type sfpStatus struct {
	Vnd string `json:"vnd"`
	Pnr string `json:"pnr"`
	Rev string `json:"rev"`
	Ser string `json:"ser"`
	Dat string `json:"dat"`
	Typ string `json:"typ"`
	Wln string `json:"wln"`
	Tmp string `json:"tmp"`
	Vcc string `json:"vcc"`
	Tbs string `json:"tbs"`
	Tpw string `json:"tpw"`
	Rpw string `json:"rpw"`
}

type SfpPage struct {
	Vendor      string
	PartNumber  string
	Revision    string
	Serial      string
	Date        string
	Type        string
	Temperature int
	Voltage     int
	TxBias      int
	TxPower     int
	RxPower     int
}

func (s *SfpPage) url() string {
	return "/sfp.b"
}

func (s *SfpPage) load(in sfpStatus) error {
	v, err := stringFromMikrotik(in.Vnd)
	if err != nil {
		return err
	}
	s.Vendor = v

	v, err = stringFromMikrotik(in.Pnr)
	if err != nil {
		return err
	}
	s.PartNumber = v

	v, err = stringFromMikrotik(in.Rev)
	if err != nil {
		return err
	}
	s.Revision = v

	v, err = stringFromMikrotik(in.Ser)
	if err != nil {
		return err
	}
	s.Serial = v

	v, err = stringFromMikrotik(in.Dat)
	if err != nil {
		return err
	}
	s.Date = v

	v, err = stringFromMikrotik(in.Typ)
	if err != nil {
		return err
	}
	s.Type = v

	i, err := parseInt(in.Tmp)
	if err != nil {
		return err
	}
	s.Temperature = i

	i, err = parseInt(in.Vcc)
	if err != nil {
		return err
	}
	s.Voltage = i

	i, err = parseInt(in.Tbs)
	if err != nil {
		return err
	}
	s.TxBias = i

	i, err = parseInt(in.Tpw)
	if err != nil {
		return err
	}
	s.TxPower = i

	i, err = parseInt(in.Rpw)
	if err != nil {
		return err
	}
	s.RxPower = i

	return nil
}

var _ swOsPage[sfpStatus] = (*SfpPage)(nil)
