package rfcat

import "fmt"

type Chip byte

const (
	ChipCC1110 Chip = 0x01
	ChipCC1111 Chip = 0x11
	ChipCC2510 Chip = 0x81
	ChipCC2511 Chip = 0x91
)

func (c Chip) Mhz() int {
	switch c {
	case ChipCC1110:
		return 26
	case ChipCC1111:
		return 24
	case ChipCC2510:
		return 26
	case ChipCC2511:
		return 24
	}
	return 24
}

func (c Chip) String() string {
	switch c {
	case ChipCC1110:
		return "CC1110"
	case ChipCC1111:
		return "CC1111"
	case ChipCC2510:
		return "CC2510"
	case ChipCC2511:
		return "CC2511"
	}
	return fmt.Sprintf("Chip(%02x)", c)
}

type Crypto byte

const (
	CryptoOff     Crypto = 0
	CryptoEncrypt Crypto = 1
	CryptoDecrypt Crypto = 2
)

func (c Crypto) String() string {
	switch c {
	case CryptoOff:
		return "Off"
	case CryptoEncrypt:
		return "Encrypt"
	case CryptoDecrypt:
		return "Decrypt"
	}
	return "InvalidCryptoType"
}

type Modulation byte

const (
	Mod2FSK       Modulation = 0x00
	ModGFSK       Modulation = 0x10
	ModASKOOK     Modulation = 0x30
	ModMSK        Modulation = 0x70
	ModManchester Modulation = 0x08
)

func (m Modulation) String() string {
	switch m {
	case Mod2FSK:
		return "2FSK"
	case ModGFSK:
		return "GFSK"
	case ModASKOOK:
		return "ASK/OOK"
	case ModMSK:
		return "MSK"
	case Mod2FSK | ModManchester:
		return "2FSK/Manchester encoding"
	case ModGFSK | ModManchester:
		return "GFSK/Manchester encoding"
	case ModASKOOK | ModManchester:
		return "ASK/OOK/Manchester encoding"
	case ModMSK | ModManchester:
		return "MSK/Manchester encoding"
	}
	return fmt.Sprintf("Modulation(%02x)", m)
}

type SyncMode byte

const (
	SyncModeNone          SyncMode = 0
	SyncMode15Of16        SyncMode = 1
	SyncMode16Of16        SyncMode = 2
	SyncMode30Of32        SyncMode = 3
	SyncModeCarrier       SyncMode = 4
	SyncModeCarrier15Of16 SyncMode = 5
	SyncModeCarrier16Of16 SyncMode = 6
	SyncModeCarrier30Of32 SyncMode = 7
)

func (s SyncMode) String() string {
	switch s {
	case SyncModeNone:
		return "None"
	case SyncMode15Of16:
		return "15 of 16 bits must match"
	case SyncMode16Of16:
		return "16 of 16 bits must match"
	case SyncMode30Of32:
		return "30 of 32 sync bits must match"
	case SyncModeCarrier:
		return "Carrier Detect"
	case SyncModeCarrier15Of16:
		return "Carrier Detect and 15 of 16 sync bits must match"
	case SyncModeCarrier16Of16:
		return "Carrier Detect and 16 of 16 sync bits must match"
	case SyncModeCarrier30Of32:
		return "Carrier Detect and 30 of 32 sync bits must match"
	}
	return fmt.Sprintf("SyncMode(%d)", s)
}

type BSLimit byte

const (
	BSLimit0  BSLimit = 0x00
	BSLimit3  BSLimit = 0x01
	BSLimit6  BSLimit = 0x02
	BSLimit12 BSLimit = 0x03
)

func (b BSLimit) String() string {
	switch b {
	case BSLimit0:
		return "No data rate offset compensation performed"
	case BSLimit3:
		return "+/- 3.125% data rate offset"
	case BSLimit6:
		return "+/- 6.25% data rate offset"
	case BSLimit12:
		return "+/- 12.5% data rate offset"
	}
	return fmt.Sprintf("BSLimit(%d)", b)
}

type LengthConfig byte

const (
	LengthConfigFixed    LengthConfig = 0
	LengthConfigVariable LengthConfig = 1
)

func (l LengthConfig) String() string {
	switch l {
	case LengthConfigFixed:
		return "Fixed Packet Mode"
	case LengthConfigVariable:
		return "Variable Packet Mode (len=first byte after sync word)"
	}
	return fmt.Sprintf("LengthConfig(%d)", l)
}

type AddressCheckType byte

const (
	AddressCheckTypeNone        AddressCheckType = 0
	AddressCheckTypeNoBroadcast AddressCheckType = 1
	AddressCheckType00          AddressCheckType = 2
	AddressCheckType00AndFF     AddressCheckType = 3
)

func (a AddressCheckType) String() string {
	switch a {
	case AddressCheckTypeNone:
		return "No address check"
	case AddressCheckTypeNoBroadcast:
		return "Address Check, No Broadcast"
	case AddressCheckType00:
		return "Address Check, 0x00 is broadcast"
	case AddressCheckType00AndFF:
		return "Address Check, 0x00 and 0xff are broadcast"
	}
	return fmt.Sprintf("AddressCheckType(%d)", a)
}

type PacketFormat byte

const (
	PacketFormatNormal   PacketFormat = 0
	PacketFormatRandomTX PacketFormat = 2
)

func (p PacketFormat) String() string {
	switch p {
	case PacketFormatNormal:
		return "Normal mode"
	case PacketFormatRandomTX:
		return "Random TX mode"
	}
	return fmt.Sprintf("PacketFormat(%d)", p)
}

type ENCCSMode byte

const (
	ENCCSModeCBC    ENCCSMode = 0
	ENCCSModeCBCMAC ENCCSMode = 1
	ENCCSModeCFB    ENCCSMode = 2
	ENCCSModeCTR    ENCCSMode = 3
	ENCCSModeECB    ENCCSMode = 4
	ENCCSModeOFB    ENCCSMode = 5
)

func (e ENCCSMode) String() string {
	switch e {
	case ENCCSModeCBC:
		return "CBC - Cipher Block Chaining"
	case ENCCSModeCBCMAC:
		return "CBC-MAC - Cipher Block Chaining Message Authentication Code"
	case ENCCSModeCFB:
		return "CFB - Cipher Feedback"
	case ENCCSModeCTR:
		return "CTR - Counter"
	case ENCCSModeECB:
		return "ECB - Electronic Codebook"
	case ENCCSModeOFB:
		return "OFB - Output Feedback"
	}
	return fmt.Sprintf("ENCCSMode(%d)", e)
}

type NumPreamble byte

const (
	NumPreamble2Bytes  NumPreamble = 0
	NumPreamble3Bytes  NumPreamble = 1
	NumPreamble4Bytes  NumPreamble = 2
	NumPreamble6Bytes  NumPreamble = 3
	NumPreamble8Bytes  NumPreamble = 4
	NumPreamble12Bytes NumPreamble = 5
	NumPreamble16Bytes NumPreamble = 6
	NumPreamble24Bytes NumPreamble = 7
)

var numPreamble = [...]int{2, 3, 4, 6, 8, 12, 16, 24}

func (n NumPreamble) Bytes() int {
	ni := int(n) >> 4
	if int(ni) < len(numPreamble) {
		return numPreamble[ni]
	}
	return -1
}

func (n NumPreamble) String() string {
	ni := int(n) >> 4
	if int(ni) < len(numPreamble) {
		return fmt.Sprintf("%d bytes", numPreamble[ni])
	}
	return fmt.Sprintf("NumPreamble(%d)", n)
}

const (
	cfgSync1 = iota
	cfgSync0
	cfgPktlen
	cfgPktctrl1
	cfgPktctrl0
	cfgAddr
	cfgChannr
	cfgFsctrl1
	cfgFsctrl0
	cfgFreq2
	cfgFreq1
	cfgFreq0
	cfgMdmcfg4
	cfgMdmcfg3
	cfgMdmcfg2
	cfgMdmcfg1
	// 0x10
	cfgMdmcfg0
	cfgDeviatn
	cfgMcsm2
	cfgMcsm1
	cfgMcsm0
	cfgFoccfg
	cfgBscfg
	cfgAgcctrl2
	cfgAgcctrl1
	cfgAgcctrl0
	cfgFrend1
	cfgFrend0
	cfgFscal3
	cfgFscal2
	cfgFscal1
	cfgFscal0
	// 0x20
	cfgZ0
	cfgZ1
	cfgZ2
	cfgTest2
	cfgTest1
	cfgTest0
	cfgZ3
	cfgPaTable7
	cfgPaTable6
	cfgPaTable5
	cfgPaTable4
	cfgPaTable3
	cfgPaTable2
	cfgPaTable1
	cfgPaTable0
	cfgIocfg2
	// 0x30
	cfgIocfg1
	cfgIocfg0
	cfgZ4
	cfgZ5
	cfgZ6
	cfgZ7
	cfgPartnum
	cfgChipID
	cfgFreqest
	cfgLQI
	cfgRSSI
	cfgMarcstate
	cfgPkstatus
	cfgVCOVCDAC
)

const (
	// band limits in Hz
	freqMin300 = 281000000
	freqMax300 = 361000000
	freqMin400 = 378000000
	freqMax400 = 481000000
	freqMin900 = 749000000
	freqMax900 = 962000000

	// band transition points in Hz
	freqEdge400 = 369000000
	freqEdge900 = 615000000

	// VCO transition points in Hz
	freqMid300 = 318000000
	freqMid400 = 424000000
	freqMid900 = 848000000
)

const (
	// AC                             = 64
	// ACC                            = 0xE0
	// ACC_0                          = 1
	// ACC_1                          = 2
	// ACC_2                          = 4
	// ACC_3                          = 8
	// ACC_4                          = 16
	// ACC_5                          = 32
	// ACC_6                          = 64
	// ACC_7                          = 128
	// ACTIVE                         = 1
	// ADCCFG                         = 0xF2
	// ADCCFG_0                       = 0x01
	// ADCCFG_1                       = 0x02
	// ADCCFG_2                       = 0x04
	// ADCCFG_3                       = 0x08
	// ADCCFG_4                       = 0x10
	// ADCCFG_5                       = 0x20
	// ADCCFG_6                       = 0x40
	// ADCCFG_7                       = 0x80
	// ADCCON1                        = 0xB4
	// ADCCON1_EOC                    = 0x80
	// ADCCON1_RCTRL                  = 0x0C
	// ADCCON1_RCTRL0                 = 0x04
	// ADCCON1_RCTRL1                 = 0x08
	// ADCCON1_RCTRL_COMPL            = (0x00 << 2)
	// ADCCON1_RCTRL_LFSR13           = (0x01 << 2)
	// ADCCON1_ST                     = 0x40
	// ADCCON1_STSEL                  = 0x30
	// ADCCON1_STSEL0                 = 0x10
	// ADCCON1_STSEL1                 = 0x20
	// ADCCON2                        = 0xB5
	// ADCCON2_ECH                    = 0x0F
	// ADCCON2_ECH0                   = 0x01
	// ADCCON2_ECH1                   = 0x02
	// ADCCON2_ECH2                   = 0x04
	// ADCCON2_ECH3                   = 0x08
	// ADCCON2_SCH                    = 0x0F
	// ADCCON2_SCH0                   = 0x01
	// ADCCON2_SCH1                   = 0x02
	// ADCCON2_SCH2                   = 0x04
	// ADCCON2_SCH3                   = 0x08
	// ADCCON2_SCH_AIN0               = (0x00)
	// ADCCON2_SCH_AIN0_1             = (0x08)
	// ADCCON2_SCH_AIN1               = (0x01)
	// ADCCON2_SCH_AIN2               = (0x02)
	// ADCCON2_SCH_AIN2_3             = (0x09)
	// ADCCON2_SCH_AIN3               = (0x03)
	// ADCCON2_SCH_AIN4               = (0x04)
	// ADCCON2_SCH_AIN4_5             = (0x0A)
	// ADCCON2_SCH_AIN5               = (0x05)
	// ADCCON2_SCH_AIN6               = (0x06)
	// ADCCON2_SCH_AIN6_7             = (0x0B)
	// ADCCON2_SCH_AIN7               = (0x07)
	// ADCCON2_SCH_GND                = (0x0C)
	// ADCCON2_SCH_POSVOL             = (0x0D)
	// ADCCON2_SCH_TEMPR              = (0x0E)
	// ADCCON2_SCH_VDD_3              = (0x0F)
	// ADCCON2_SDIV                   = 0x30
	// ADCCON2_SDIV0                  = 0x10
	// ADCCON2_SDIV1                  = 0x20
	// ADCCON2_SDIV_128               = (0x01 << 4)
	// ADCCON2_SDIV_256               = (0x02 << 4)
	// ADCCON2_SDIV_512               = (0x03 << 4)
	// ADCCON2_SDIV_64                = (0x00 << 4)
	// ADCCON2_SREF                   = 0xC0
	// ADCCON2_SREF0                  = 0x40
	// ADCCON2_SREF1                  = 0x80
	// ADCCON2_SREF_1_25V             = (0x00 << 6)
	// ADCCON2_SREF_AVDD              = (0x02 << 6)
	// ADCCON2_SREF_P0_6_P0_7         = (0x03 << 6)
	// ADCCON2_SREF_P0_7              = (0x01 << 6)
	// ADCCON3                        = 0xB6
	// ADCCON3_ECH_AIN0               = (0x00)
	// ADCCON3_ECH_AIN0_1             = (0x08)
	// ADCCON3_ECH_AIN1               = (0x01)
	// ADCCON3_ECH_AIN2               = (0x02)
	// ADCCON3_ECH_AIN2_3             = (0x09)
	// ADCCON3_ECH_AIN3               = (0x03)
	// ADCCON3_ECH_AIN4               = (0x04)
	// ADCCON3_ECH_AIN4_5             = (0x0A)
	// ADCCON3_ECH_AIN5               = (0x05)
	// ADCCON3_ECH_AIN6               = (0x06)
	// ADCCON3_ECH_AIN6_7             = (0x0B)
	// ADCCON3_ECH_AIN7               = (0x07)
	// ADCCON3_ECH_GND                = (0x0C)
	// ADCCON3_ECH_POSVOL             = (0x0D)
	// ADCCON3_ECH_TEMPR              = (0x0E)
	// ADCCON3_ECH_VDD_3              = (0x0F)
	// ADCCON3_EDIV                   = 0x30
	// ADCCON3_EDIV0                  = 0x10
	// ADCCON3_EDIV1                  = 0x20
	// ADCCON3_EDIV_128               = (0x01 << 4)
	// ADCCON3_EDIV_256               = (0x02 << 4)
	// ADCCON3_EDIV_512               = (0x03 << 4)
	// ADCCON3_EDIV_64                = (0x00 << 4)
	// ADCCON3_EREF                   = 0xC0
	// ADCCON3_EREF0                  = 0x40
	// ADCCON3_EREF1                  = 0x80
	// ADCCON3_EREF_1_25V             = (0x00 << 6)
	// ADCCON3_EREF_AVDD              = (0x02 << 6)
	// ADCCON3_EREF_P0_6_P0_7         = (0x03 << 6)
	// ADCCON3_EREF_P0_7              = (0x01 << 6)
	// ADCH                           = 0xBB
	// ADCIE                          = 2
	// ADCIF                          = 32
	// ADCL                           = 0xBA
	// ADC_VECTOR                     = 1    #  ADC End of Conversion
	// ADDR                           = 0xDF05
	// ADR_CHK_0_255_BRDCST           = (0x03)
	// ADR_CHK_0_BRDCST               = (0x02)
	// ADR_CHK_NONE                   = (0x00)
	// ADR_CHK_NO_BRDCST              = (0x01)
	// AGCCTRL0                       = 0xDF19
	// AGCCTRL0_AGC_FREEZE            = 0x0C
	// AGCCTRL0_FILTER_LENGTH         = 0x03
	// AGCCTRL0_HYST_LEVEL            = 0xC0
	// AGCCTRL0_WAIT_TIME             = 0x30
	// AGCCTRL1                       = 0xDF18
	// AGCCTRL1_AGC_LNA_PRIORITY      = 0x40
	// AGCCTRL1_CARRIER_SENSE_ABS_THR = 0x0F
	// AGCCTRL1_CARRIER_SENSE_REL_THR = 0x30
	// AGCCTRL2                       = 0xDF17
	// AGCCTRL2_MAGN_TARGET           = 0x07
	// AGCCTRL2_MAX_DVGA_GAIN         = 0xC0
	// AGCCTRL2_MAX_LNA_GAIN          = 0x38
	// B                              = 0xF0
	// BSCFG                          = 0xDF16
	bscfgBSLimit = 0x03
	// BSCFG_BS_LIMIT0                = 0x01
	// BSCFG_BS_LIMIT1                = 0x02
	// BSCFG_BS_POST_KI               = 0x08
	// BSCFG_BS_POST_KP               = 0x04
	// BSCFG_BS_PRE_KI                = 0xC0
	// BSCFG_BS_PRE_KI0               = 0x40
	// BSCFG_BS_PRE_KI1               = 0x80
	// BSCFG_BS_PRE_KI_1K             = (0x00 << 6)
	// BSCFG_BS_PRE_KI_2K             = (0x01 << 6)
	// BSCFG_BS_PRE_KI_3K             = (0x02 << 6)
	// BSCFG_BS_PRE_KI_4K             = (0x03 << 6)
	// BSCFG_BS_PRE_KP                = 0x30
	// BSCFG_BS_PRE_KP0               = 0x10
	// BSCFG_BS_PRE_KP1               = 0x20
	// BSCFG_BS_PRE_KP_1K             = (0x00 << 4)
	// BSCFG_BS_PRE_KP_2K             = (0x01 << 4)
	// BSCFG_BS_PRE_KP_3K             = (0x02 << 4)
	// BSCFG_BS_PRE_KP_4K             = (0x03 << 4)
	// B_0                            = 1
	// B_1                            = 2
	// B_2                            = 4
	// B_3                            = 8
	// B_4                            = 16
	// B_5                            = 32
	// B_6                            = 64
	// B_7                            = 128
	regChannr = 0xDF06
	// CLKCON                         = 0xC6
	// CLKCON_CLKSPD                  = 0x07  # bit mask, for the clock speed
	// CLKCON_CLKSPD0                 = 0x01  # bit mask, for the clock speed
	// CLKCON_CLKSPD1                 = 0x02  # bit mask, for the clock speed
	// CLKCON_CLKSPD2                 = 0x04  # bit mask, for the clock speed
	// CLKCON_OSC                     = 0x40  # bit mask, for the system clock oscillator
	// CLKCON_OSC32                   = 0x80  # bit mask, for the slow 32k clock oscillator
	// CLKCON_TICKSPD                 = 0x38  # bit mask, for timer ticks output setting
	// CLKCON_TICKSPD0                = 0x08  # bit mask, for timer ticks output setting
	// CLKCON_TICKSPD1                = 0x10  # bit mask, for timer ticks output setting
	// CLKCON_TICKSPD2                = 0x20  # bit mask, for timer ticks output setting
	// CLKSPD_DIV_1                   = (0x00)
	// CLKSPD_DIV_128                 = (0x07)
	// CLKSPD_DIV_16                  = (0x04)
	// CLKSPD_DIV_2                   = (0x01)
	// CLKSPD_DIV_32                  = (0x05)
	// CLKSPD_DIV_4                   = (0x02)
	// CLKSPD_DIV_64                  = (0x06)
	// CLKSPD_DIV_8                   = (0x03)
	// CY                             = 128
	regDeviatn = 0xDF11
	// DEVIATN_DEVIATION_E            = 0x70
	// DEVIATN_DEVIATION_E0           = 0x10
	// DEVIATN_DEVIATION_E1           = 0x20
	// DEVIATN_DEVIATION_E2           = 0x40
	deviatnDeviationM = 0x07
	// DEVIATN_DEVIATION_M0           = 0x01
	// DEVIATN_DEVIATION_M1           = 0x02
	// DEVIATN_DEVIATION_M2           = 0x04
	// DMA0CFGH                       = 0xD5
	// DMA0CFGL                       = 0xD4
	// DMA1CFGH                       = 0xD3
	// DMA1CFGL                       = 0xD2
	// DMAARM                         = 0xD6
	// DMAARM0                        = 0x01
	// DMAARM1                        = 0x02
	// DMAARM2                        = 0x04
	// DMAARM3                        = 0x08
	// DMAARM4                        = 0x10
	// DMAARM_ABORT                   = 0x80
	// DMAIE                          = 1
	// DMAIF                          = 1
	// DMAIRQ                         = 0xD1
	// DMAIRQ_DMAIF0                  = 0x01
	// DMAIRQ_DMAIF1                  = 0x02
	// DMAIRQ_DMAIF2                  = 0x04
	// DMAIRQ_DMAIF3                  = 0x08
	// DMAIRQ_DMAIF4                  = 0x10
	// DMAREQ                         = 0xD7
	// DMAREQ0                        = 0x01
	// DMAREQ1                        = 0x02
	// DMAREQ2                        = 0x04
	// DMAREQ3                        = 0x08
	// DMAREQ4                        = 0x10
	// DMA_VECTOR                     = 8    #  DMA Transfer Complete
	// DPH0                           = 0x83
	// DPH1                           = 0x85
	// DPL0                           = 0x82
	// DPL1                           = 0x84
	// DPS                            = 0x92
	// DPS_VDPS                       = 0x01
	// DSM_IP_OFF_OS_OFF              = (0x03)    # Interpolator & output shaping disabled
	// DSM_IP_OFF_OS_ON               = (0x02)    # Interpolator disabled & output shaping enabled
	// DSM_IP_ON_OS_OFF               = (0x01)    # Interpolator enabled & output shaping disabled
	// DSM_IP_ON_OS_ON                = (0x00)    # Interpolator & output shaping enabled
	// EA                             = 128
	// ENCCS                          = 0xB3
	// ENCCS_CMD                      = 0x06
	// ENCCS_CMD0                     = 0x02
	// ENCCS_CMD1                     = 0x04
	// ENCCS_CMD_DEC                  = (0x01 << 1)
	// ENCCS_CMD_ENC                  = (0x00 << 1)
	// ENCCS_CMD_LDIV                 = (0x03 << 1)
	// ENCCS_CMD_LDKEY                = (0x02 << 1)
	// ENCCS_MODE                     = 0x70
	// ENCCS_MODE0                    = 0x10
	// ENCCS_MODE1                    = 0x20
	// ENCCS_MODE2                    = 0x40
	// ENCCS_MODE_CBC                 = (0x00 << 4)
	// ENCCS_MODE_CBCMAC              = (0x05 << 4)
	// ENCCS_MODE_CFB                 = (0x01 << 4)
	// ENCCS_MODE_CTR                 = (0x03 << 4)
	// ENCCS_MODE_ECB                 = (0x04 << 4)
	// ENCCS_MODE_OFB                 = (0x02 << 4)
	// ENCCS_RDY                      = 0x08
	// ENCCS_ST                       = 0x01
	// ENCDI                          = 0xB1
	// ENCDO                          = 0xB2
	// ENCIE                          = 16
	// ENCIF_0                        = 1
	// ENCIF_1                        = 2
	// ENC_VECTOR                     = 4    #  AES Encryption/Decryption Complete
	// EP_STATE_IDLE                  = 0
	// EP_STATE_RX                    = 2
	// EP_STATE_STALL                 = 3
	// EP_STATE_TX                    = 1
	// ERR                            = 8
	// F0                             = 32
	// F1                             = 2
	// FADDRH                         = 0xAD
	// FADDRL                         = 0xAC
	// FCTL                           = 0xAE
	// FCTL_BUSY                      = 0x80
	// FCTL_CONTRD                    = 0x10
	// FCTL_ERASE                     = 0x01
	// FCTL_SWBSY                     = 0x40
	// FCTL_WRITE                     = 0x02
	// FE                             = 16
	// FOCCFG                         = 0xDF15
	// FOCCFG_FOC_BS_CS_GATE          = 0x20
	// FOCCFG_FOC_LIMIT               = 0x03
	// FOCCFG_FOC_LIMIT0              = 0x01
	// FOCCFG_FOC_LIMIT1              = 0x02
	// FOCCFG_FOC_POST_K              = 0x04
	// FOCCFG_FOC_PRE_K               = 0x18
	// FOCCFG_FOC_PRE_K0              = 0x08
	// FOCCFG_FOC_PRE_K1              = 0x10
	// FOC_LIMIT_0                    = (0x00)
	// FOC_LIMIT_DIV2                 = (0x03)
	// FOC_LIMIT_DIV4                 = (0x02)
	// FOC_LIMIT_DIV8                 = (0x01)
	// FOC_PRE_K_1K                   = (0x00 << 3)
	// FOC_PRE_K_2K                   = (0x02 << 3)
	// FOC_PRE_K_3K                   = (0x03 << 3)
	// FOC_PRE_K_4K                   = (0x04 << 3)
	frend0 = 0xDF1B
	// FREND0_LODIV_BUF_CURRENT_TX    = 0x30
	frend0PaPower = 0x07
	// FREND1                         = 0xDF1A
	// FREND1_LNA2MIX_CURRENT         = 0x30
	// FREND1_LNA_CURRENT             = 0xC0
	// FREND1_LODIV_BUF_CURRENT_RX    = 0x0C
	// FREND1_MIX_CURRENT             = 0x03
	freq0              = 0xDF0B
	freq1              = 0xDF0A
	freq2              = 0xDF09
	freqest            = 0xDF38
	fscal0             = 0xDF1F
	fscal1             = 0xDF1E
	fscal2             = 0xDF1D
	fscal2Fscal2       = 0x1F
	fscal2VCOCodeHEn   = 0x20
	fscal3             = 0xDF1C
	fscal3ChpCurrCalEn = 0x30
	fscal3Fscal3       = 0xC0
	fsctrl0            = 0xDF08
	fsctrl1            = 0xDF07
	// FS_AUTOCAL_4TH_TO_IDLE         = (0x03 << 4)
	// FS_AUTOCAL_FROM_IDLE           = (0x01 << 4)
	// FS_AUTOCAL_NEVER               = (0x00 << 4)
	// FS_AUTOCAL_TO_IDLE             = (0x02 << 4)
	// FWDATA                         = 0xAF
	// FWT                            = 0xAB
	// I2SCFG0                        = 0xDF40
	// I2SCFG0_ENAB                   = 0x01
	// I2SCFG0_MASTER                 = 0x02
	// I2SCFG0_RXIEN                  = 0x40
	// I2SCFG0_RXMONO                 = 0x04
	// I2SCFG0_TXIEN                  = 0x80
	// I2SCFG0_TXMONO                 = 0x08
	// I2SCFG0_ULAWC                  = 0x10
	// I2SCFG0_ULAWE                  = 0x20
	// I2SCFG1                        = 0xDF41
	// I2SCFG1_IOLOC                  = 0x01
	// I2SCFG1_TRIGNUM                = 0x06
	// I2SCFG1_TRIGNUM0               = 0x02
	// I2SCFG1_TRIGNUM1               = 0x04
	// I2SCFG1_TRIGNUM_IOC_1          = (0x02 << 1)
	// I2SCFG1_TRIGNUM_NO_TRIG        = (0x00 << 1)
	// I2SCFG1_TRIGNUM_T1_CH0         = (0x03 << 1)
	// I2SCFG1_TRIGNUM_USB_SOF        = (0x01 << 1)
	// I2SCFG1_WORDS                  = 0xF8
	// I2SCFG1_WORDS0                 = 0x08
	// I2SCFG1_WORDS1                 = 0x10
	// I2SCFG1_WORDS2                 = 0x20
	// I2SCFG1_WORDS3                 = 0x40
	// I2SCFG1_WORDS4                 = 0x80
	// I2SCLKF0                       = 0xDF46
	// I2SCLKF1                       = 0xDF47
	// I2SCLKF2                       = 0xDF48
	// I2SDATH                        = 0xDF43
	// I2SDATL                        = 0xDF42
	// I2SSTAT                        = 0xDF45
	// I2SSTAT_RXIRQ                  = 0x04
	// I2SSTAT_RXLR                   = 0x10
	// I2SSTAT_RXOVF                  = 0x40
	// I2SSTAT_TXIRQ                  = 0x08
	// I2SSTAT_TXLR                   = 0x20
	// I2SSTAT_TXUNF                  = 0x80
	// I2SSTAT_WCNT                   = 0x03
	// I2SSTAT_WCNT0                  = 0x01
	// I2SSTAT_WCNT1                  = 0x02
	// I2SSTAT_WCNT_10BIT             = (0x02)
	// I2SSTAT_WCNT_9BIT              = (0x01)
	// I2SSTAT_WCNT_9_10BIT           = (0x02)
	// I2SWCNT                        = 0xDF44
	// IEN0                           = 0xA8
	// IEN1                           = 0xB8
	// IEN2                           = 0x9A
	// IEN2_I2STXIE                   = 0x08
	// IEN2_P1IE                      = 0x10
	// IEN2_P2IE                      = 0x02
	// IEN2_RFIE                      = 0x01
	// IEN2_USBIE                     = 0x02
	// IEN2_UTX0IE                    = 0x04
	// IEN2_UTX1IE                    = 0x08
	// IEN2_WDTIE                     = 0x20
	// IOCFG0                         = 0xDF31
	// IOCFG0_GDO0_CFG                = 0x3F
	// IOCFG0_GDO0_INV                = 0x40
	// IOCFG1                         = 0xDF30
	// IOCFG1_GDO1_CFG                = 0x3F
	// IOCFG1_GDO1_INV                = 0x40
	// IOCFG1_GDO_DS                  = 0x80
	// IOCFG2                         = 0xDF2F
	// IOCFG2_GDO2_CFG                = 0x3F
	// IOCFG2_GDO2_INV                = 0x40
	// IP0                            = 0xA9
	// IP0_IPG0                       = 0x01
	// IP0_IPG1                       = 0x02
	// IP0_IPG2                       = 0x04
	// IP0_IPG3                       = 0x08
	// IP0_IPG4                       = 0x10
	// IP0_IPG5                       = 0x20
	// IP1                            = 0xB9
	// IP1_IPG0                       = 0x01
	// IP1_IPG1                       = 0x02
	// IP1_IPG2                       = 0x04
	// IP1_IPG3                       = 0x08
	// IP1_IPG4                       = 0x10
	// IP1_IPG5                       = 0x20
	// IRCON                          = 0xC0
	// IRCON2                         = 0xE8
	// IT0                            = 1
	// IT1                            = 4
	lqi = 0xDF39
	// MARCSTATE                      = 0xDF3B
	// MARCSTATE_MARC_STATE           = 0x1F
	// MARC_STATE_BWBOOST             = 0x09
	// MARC_STATE_ENDCAL              = 0x0C
	// MARC_STATE_FSTXON              = 0x12
	// MARC_STATE_FS_LOCK             = 0x0A
	marcStateIdle = 0x01
	// MARC_STATE_IFADCON             = 0x0B
	// MARC_STATE_MANCAL              = 0x05
	// MARC_STATE_REGON               = 0x07
	// MARC_STATE_REGON_MC            = 0x04
	marcStateRX = 0x0D
	// MARC_STATE_RXTX_SWITCH         = 0x15
	// MARC_STATE_RX_END              = 0x0E
	// MARC_STATE_RX_OVERFLOW         = 0x11
	// MARC_STATE_RX_RST              = 0x0F
	// MARC_STATE_SLEEP               = 0x00
	// MARC_STATE_STARTCAL            = 0x08
	marcStateTX = 0x13
	// MARC_STATE_TXRX_SWITCH         = 0x10
	// MARC_STATE_TX_END              = 0x14
	// MARC_STATE_TX_UNDERFLOW        = 0x16
	// MARC_STATE_VCOON               = 0x06
	// MARC_STATE_VCOON_MC            = 0x03
	// MCSM0                          = 0xDF14
	// MCSM0_FS_AUTOCAL               = 0x30
	// MCSM1                          = 0xDF13
	// MCSM1_CCA_MODE                 = 0x30
	// MCSM1_CCA_MODE0                = 0x10
	// MCSM1_CCA_MODE1                = 0x20
	// MCSM1_CCA_MODE_ALWAYS          = (0x00 << 4)
	// MCSM1_CCA_MODE_PACKET          = (0x02 << 4)
	// MCSM1_CCA_MODE_RSSI0           = (0x01 << 4)
	// MCSM1_CCA_MODE_RSSI1           = (0x03 << 4)
	// MCSM1_RXOFF_MODE               = 0x0C
	// MCSM1_RXOFF_MODE0              = 0x04
	// MCSM1_RXOFF_MODE1              = 0x08
	// MCSM1_RXOFF_MODE_FSTXON        = (0x01 << 2)
	// MCSM1_RXOFF_MODE_IDLE          = (0x00 << 2)
	// MCSM1_RXOFF_MODE_RX            = (0x03 << 2)
	// MCSM1_RXOFF_MODE_TX            = (0x02 << 2)
	// MCSM1_TXOFF_MODE               = 0x03
	// MCSM1_TXOFF_MODE0              = 0x01
	// MCSM1_TXOFF_MODE1              = 0x02
	// MCSM1_TXOFF_MODE_FSTXON        = (0x01 << 0)
	// MCSM1_TXOFF_MODE_IDLE          = (0x00 << 0)
	// MCSM1_TXOFF_MODE_RX            = (0x03 << 0)
	// MCSM1_TXOFF_MODE_TX            = (0x02 << 0)
	// MCSM2                          = 0xDF12
	// MCSM2_RX_TIME                  = 0x07
	// MCSM2_RX_TIME_QUAL             = 0x08
	// MCSM2_RX_TIME_RSSI             = 0x10
	mdmcfg0             = 0xDF10
	mdmcfg1             = 0xDF0F
	mdmcfg1ChanspcE     = 0x03
	mdmcfg2             = 0xDF0E
	mdmcfg2DemDCFiltOff = 0x80
	mdmcfg2ManchesterEn = 0x08
	mdmcfg2ModFormat    = 0x70
	mdmcfg2ModFormat0   = 0x10
	mdmcfg2ModFormat1   = 0x20
	mdmcfg2ModFormat2   = 0x40
	mdmcfg2SyncMode     = 0x07
	mdmcfg2SyncMode0    = 0x01
	mdmcfg2SyncMode1    = 0x02
	mdmcfg2SyncMode2    = 0x04

	mdmcfg3 = 0xDF0D
	mdmcfg4 = 0xDF0C
	// MDMCFG4_CHANBW_E               = 0xC0
	// MDMCFG4_CHANBW_M               = 0x30
	mdmcfg4DRateE = 0x0F
	// MFMCFG1_CHANSPC_E              = 0x03
	// MFMCFG1_CHANSPC_E0             = 0x01
	// MFMCFG1_CHANSPC_E1             = 0x02
	mfmcfg1FECEn       = 0x80
	mfmcfg1NumPreamble = 0x70

	// MFMCFG1_NUM_PREAMBLE0          = 0x10
	// MFMCFG1_NUM_PREAMBLE1          = 0x20
	// MFMCFG1_NUM_PREAMBLE2          = 0x40
	// MFMCFG1_NUM_PREAMBLE_12        = (0x05 << 4)
	// MFMCFG1_NUM_PREAMBLE_16        = (0x06 << 4)
	// MFMCFG1_NUM_PREAMBLE_2         = (0x00 << 4)
	// MFMCFG1_NUM_PREAMBLE_24        = (0x07 << 4)
	// MFMCFG1_NUM_PREAMBLE_3         = (0x01 << 4)
	// MFMCFG1_NUM_PREAMBLE_4         = (0x02 << 4)
	// MFMCFG1_NUM_PREAMBLE_6         = (0x03 << 4)
	// MFMCFG1_NUM_PREAMBLE_8         = (0x04 << 4)
	// MDMCTRL0H                      = 0xDF02
	// MEMCTR                         = 0xC7
	// MEMCTR_CACHD                   = 0x02
	// MEMCTR_PREFD                   = 0x01
	// MODE                           = 128
	// MOD_FORMAT_2_FSK               = (0x00 << 4)
	// MOD_FORMAT_GFSK                = (0x01 << 4)
	// MOD_FORMAT_MSK                 = (0x07 << 4)
	// MPAGE                          = 0x93
	// OV                             = 4
	// OVFIM                          = 64
	// P                              = 1
	// P0                             = 0x80
	// P0DIR                          = 0xFD
	// P0IE                           = 32
	// P0IF                           = 32
	// P0IFG                          = 0x89
	// P0IFG_USB_RESUME               = 0x80    #rw0
	// P0INP                          = 0x8F
	// P0INT_VECTOR                   = 13   #  Port 0 Inputs
	// P0SEL                          = 0xF3
	// P0_0                           = 1
	// P0_1                           = 2
	// P0_2                           = 4
	// P0_3                           = 8
	// P0_4                           = 16
	// P0_5                           = 32
	// P0_6                           = 64
	// P0_7                           = 128
	// P1                             = 0x90
	// P1DIR                          = 0xFE
	// P1IEN                          = 0x8D
	// P1IF                           = 8
	// P1IFG                          = 0x8A
	// P1INP                          = 0xF6
	// P1INT_VECTOR                   = 15   #  Port 1 Inputs
	// P1SEL                          = 0xF4
	// P1_0                           = 1
	// P1_1                           = 2
	// P1_2                           = 4
	// P1_3                           = 8
	// P1_4                           = 16
	// P1_5                           = 32
	// P1_6                           = 64
	// P1_7                           = 128
	// P2                             = 0xA0
	// P2DIR                          = 0xFF
	// P2DIR_0PRIP0                   = 0x40
	// P2DIR_1PRIP0                   = 0x80
	// P2DIR_DIRP2                    = 0x1F
	// P2DIR_DIRP2_0                  = (0x01)
	// P2DIR_DIRP2_1                  = (0x02)
	// P2DIR_DIRP2_2                  = (0x04)
	// P2DIR_DIRP2_3                  = (0x08)
	// P2DIR_DIRP2_4                  = (0x10)
	// P2DIR_PRIP0                    = 0xC0
	// P2DIR_PRIP0_0                  = (0x00 << 6)
	// P2DIR_PRIP0_1                  = (0x01 << 6)
	// P2DIR_PRIP0_2                  = (0x02 << 6)
	// P2DIR_PRIP0_3                  = (0x03 << 6)
	// P2IF                           = 1
	// P2IFG                          = 0x8B
	// P2INP                          = 0xF7
	// P2INP_MDP2                     = 0x1F
	// P2INP_MDP2_0                   = (0x01)
	// P2INP_MDP2_1                   = (0x02)
	// P2INP_MDP2_2                   = (0x04)
	// P2INP_MDP2_3                   = (0x08)
	// P2INP_MDP2_4                   = (0x10)
	// P2INP_PDUP0                    = 0x20
	// P2INP_PDUP1                    = 0x40
	// P2INP_PDUP2                    = 0x80
	// P2INT_VECTOR                   = 6    #  Port 2 Inputs
	// P2SEL                          = 0xF5
	// P2SEL_PRI0P1                   = 0x08
	// P2SEL_PRI1P1                   = 0x10
	// P2SEL_PRI2P1                   = 0x20
	// P2SEL_PRI3P1                   = 0x40
	// P2SEL_SELP2_0                  = 0x01
	// P2SEL_SELP2_3                  = 0x02
	// P2SEL_SELP2_4                  = 0x04
	// P2_0                           = 1
	// P2_1                           = 2
	// P2_2                           = 4
	// P2_3                           = 8
	// P2_4                           = 16
	// P2_5                           = 32
	// P2_6                           = 64
	// P2_7                           = 128
	// PARTNUM                        = 0xDF36
	paTable0 = 0xDF2E
	paTable1 = 0xDF2D
	paTable2 = 0xDF2C
	paTable3 = 0xDF2B
	paTable4 = 0xDF2A
	paTable5 = 0xDF29
	paTable6 = 0xDF28
	paTable7 = 0xDF27
	// PCON                           = 0x87
	// PCON_IDLE                      = 0x01
	// PERCFG                         = 0xF1
	// PERCFG_T1CFG                   = 0x40
	// PERCFG_T3CFG                   = 0x20
	// PERCFG_T4CFG                   = 0x10
	// PERCFG_U0CFG                   = 0x01
	// PERCFG_U1CFG                   = 0x02
	// PICTL                          = 0x8C
	// PICTL_P0ICON                   = 0x01
	// PICTL_P0IENH                   = 0x10
	// PICTL_P0IENL                   = 0x08
	// PICTL_P1ICON                   = 0x02
	// PICTL_P2ICON                   = 0x04
	// PICTL_P2IEN                    = 0x20
	// PICTL_PADSC                    = 0x40
	pktctrl0Reg = 0xDF04
	// PKTCTRL0_CC2400_EN             = 0x08
	pktctrl0CRCEn         = 0x04
	pktctrl0LengthConfig  = 0x03
	pktctrl0LengthConfig0 = 0x01
	// PKTCTRL0_LENGTH_CONFIG_FIX     = (0x00)
	// PKTCTRL0_LENGTH_CONFIG_VAR     = (0x01)
	// PKTCTRL0_PKT_FORMAT            = 0x30
	// PKTCTRL0_PKT_FORMAT0           = 0x10
	// PKTCTRL0_PKT_FORMAT1           = 0x20
	pktctrl0WhiteData = 0x40
	pktctrl1Reg       = 0xDF03
	// PKTCTRL1_ADR_CHK               = 0x03
	// PKTCTRL1_ADR_CHK0              = 0x01
	// PKTCTRL1_ADR_CHK1              = 0x02
	pktctrl1AppendStatus = 0x04
	// PKTCTRL1_PQT                   = 0xE0
	// PKTCTRL1_PQT0                  = 0x20
	// PKTCTRL1_PQT1                  = 0x40
	// PKTCTRL1_PQT2                  = 0x80
	pktlenReg    = 0xDF02
	pktstatusReg = 0xDF3C
	// PKT_FORMAT_NORM                = (0x00)
	// PKT_FORMAT_RAND                = (0x02)
	// PSW                            = 0xD0
	// RE                             = 64
	// RFD                            = 0xD9
	// RFIF                           = 0xE9
	// RFIF_IRQ_CCA                   = 0x02
	// RFIF_IRQ_CS                    = 0x08
	// RFIF_IRQ_DONE                  = 0x10
	// RFIF_IRQ_PQT                   = 0x04
	// RFIF_IRQ_RXOVF                 = 0x40
	// RFIF_IRQ_SFD                   = 0x01
	// RFIF_IRQ_TIMEOUT               = 0x20
	// RFIF_IRQ_TXUNF                 = 0x80
	// RFIM                           = 0x91
	// RFIM_IM_CCA                    = 0x02
	// RFIM_IM_CS                     = 0x08
	// RFIM_IM_DONE                   = 0x10
	// RFIM_IM_PQT                    = 0x04
	// RFIM_IM_RXOVF                  = 0x40
	// RFIM_IM_SFD                    = 0x01
	// RFIM_IM_TIMEOUT                = 0x20
	// RFIM_IM_TXUNF                  = 0x80
	// RFST                           = 0xE1
	rfstSCAL    = 0x01
	rfstSFSTXON = 0x00
	rfstSIDLE   = 0x04
	rfstSNOP    = 0x05
	rfstSRX     = 0x02
	rfstSTX     = 0x03

	// RFTXRXIE                       = 1
	// RFTXRXIF                       = 2
	// RFTXRX_VECTOR                  = 0    #  RF TX done / RX ready
	// RF_VECTOR                      = 16   #  RF General Interrupts
	// RNDH                           = 0xBD
	// RNDL                           = 0xBC
	// RS0                            = 8
	// RS1                            = 16
	rssi = 0xDF3A
	// RX_BYTE                        = 4
	// S0CON                          = 0x98
	// S1CON                          = 0x9B
	// S1CON_RFIF_0                   = 0x01
	// S1CON_RFIF_1                   = 0x02
	// SLAVE                          = 32
	// SLEEP                          = 0xBE
	// SLEEP_HFRC_S                   = 0x20
	// SLEEP_MODE                     = 0x03
	// SLEEP_MODE0                    = 0x01
	// SLEEP_MODE1                    = 0x02
	// SLEEP_MODE_PM0                 = (0x00)
	// SLEEP_MODE_PM1                 = (0x01)
	// SLEEP_MODE_PM2                 = (0x02)
	// SLEEP_MODE_PM3                 = (0x03)
	// SLEEP_OSC_PD                   = 0x04
	// SLEEP_RST                      = 0x18
	// SLEEP_RST0                     = 0x08
	// SLEEP_RST1                     = 0x10
	// SLEEP_RST_EXT                  = (0x01 << 3)
	// SLEEP_RST_POR_BOD              = (0x00 << 3)
	// SLEEP_RST_WDT                  = (0x02 << 3)
	// SLEEP_USB_EN                   = 0x80
	// SLEEP_XOSC_S                   = 0x40
	// SP                             = 0x81
	// STIE                           = 32
	// STIF                           = 128
	// STSEL_FULL_SPEED               = (0x01 << 4)
	// STSEL_P2_0                     = (0x00 << 4)
	// STSEL_ST                       = (0x03 << 4)
	// STSEL_T1C0_CMP_EVT             = (0x02 << 4)
	// ST_VECTOR                      = 5    #  Sleep Timer Compare
	sync0 = 0xDF01
	sync1 = 0xDF00
	// SYNC_MODE_15_16                = (0x01)
	// SYNC_MODE_15_16_CS             = (0x05)
	// SYNC_MODE_16_16                = (0x02)
	// SYNC_MODE_16_16_CS             = (0x06)
	// SYNC_MODE_30_32                = (0x03)
	// SYNC_MODE_30_32_CS             = (0x07)
	// SYNC_MODE_NO_PRE               = (0x00)
	// SYNC_MODE_NO_PRE_CS            = (0x04)    # CS = carrier-sense above threshold
	// T1C0_BOTH_EDGE                 = (0x03)    # Capture on both edges
	// T1C0_CLR_CMP_UP_SET_0          = (0x04 << 3)    # Set output on compare
	// T1C0_CLR_ON_CMP                = (0x01 << 3)    # Set output on compare-up clear on 0
	// T1C0_FALL_EDGE                 = (0x02)    # Capture on falling edge
	// T1C0_NO_CAP                    = (0x00)    # No capture
	// T1C0_RISE_EDGE                 = (0x01)    # Capture on rising edge
	// T1C0_SET_CMP_UP_CLR_0          = (0x03 << 3)    # Clear output on compare
	// T1C0_SET_ON_CMP                = (0x00 << 3)    # Clear output on compare-up set on 0
	// T1C0_TOG_ON_CMP                = (0x02 << 3)    # Toggle output on compare
	// T1C1_BOTH_EDGE                 = (0x03)    # Capture on both edges
	// T1C1_CLR_C1_SET_C0             = (0x06 << 3)  # Clear when equal to T1CC1, set when equal to T1CC0
	// T1C1_CLR_CMP_UP_SET_0          = (0x04 << 3)  # Clear output on compare-up set on 0
	// T1C1_CLR_ON_CMP                = (0x01 << 3)  # Clear output on compare
	// T1C1_DSM_MODE                  = (0x07 << 3)  # DSM mode
	// T1C1_FALL_EDGE                 = (0x02)    # Capture on falling edge
	// T1C1_NO_CAP                    = (0x00)     # No capture
	// T1C1_RISE_EDGE                 = (0x01)    # Capture on rising edge
	// T1C1_SET_C1_CLR_C0             = (0x05 << 3)  # Set when equal to T1CC1, clear when equal to T1CC0
	// T1C1_SET_CMP_UP_CLR_0          = (0x03 << 3)  # Set output on compare-up clear on 0
	// T1C1_SET_ON_CMP                = (0x00 << 3)  # Set output on compare
	// T1C1_TOG_ON_CMP                = (0x02 << 3)  # Toggle output on compare
	// T1C2_BOTH_EDGE                 = (0x03)    # Capture on both edges
	// T1C2_CLR_C2_SET_C0             = (0x06 << 3)  # Clear when equal to T1CC2, set when equal to T1CC0
	// T1C2_CLR_CMP_UP_SET_0          = (0x04 << 3)  # Clear output on compare-up set on 0
	// T1C2_CLR_ON_CMP                = (0x01 << 3)  # Clear output on compare
	// T1C2_FALL_EDGE                 = (0x02)    # Capture on falling edge
	// T1C2_NO_CAP                    = (0x00)     # No capture
	// T1C2_RISE_EDGE                 = (0x01)    # Capture on rising edge
	// T1C2_SET_C2_CLR_C0             = (0x05 << 3)  # Set when equal to T1CC2, clear when equal to T1CC0
	// T1C2_SET_CMP_UP_CLR_0          = (0x03 << 3)  # Set output on compare-up clear on 0
	// T1C2_SET_ON_CMP                = (0x00 << 3)  # Set output on compare
	// T1C2_TOG_ON_CMP                = (0x02 << 3)  # Toggle output on compare
	// T1CC0H                         = 0xDB
	// T1CC0L                         = 0xDA
	// T1CC1H                         = 0xDD
	// T1CC1L                         = 0xDC
	// T1CC2H                         = 0xDF
	// T1CC2L                         = 0xDE
	// T1CCTL0                        = 0xE5
	// T1CCTL0_CAP                    = 0x03
	// T1CCTL0_CAP0                   = 0x01
	// T1CCTL0_CAP1                   = 0x02
	// T1CCTL0_CMP                    = 0x38
	// T1CCTL0_CMP0                   = 0x08
	// T1CCTL0_CMP1                   = 0x10
	// T1CCTL0_CMP2                   = 0x20
	// T1CCTL0_CPSEL                  = 0x80    # Timer 1 channel 0 capture select
	// T1CCTL0_IM                     = 0x40    # Channel 0 Interrupt mask
	// T1CCTL0_MODE                   = 0x04    # Capture or compare mode
	// T1CCTL1                        = 0xE6
	// T1CCTL1_CAP                    = 0x03
	// T1CCTL1_CAP0                   = 0x01
	// T1CCTL1_CAP1                   = 0x02
	// T1CCTL1_CMP                    = 0x38
	// T1CCTL1_CMP0                   = 0x08
	// T1CCTL1_CMP1                   = 0x10
	// T1CCTL1_CMP2                   = 0x20
	// T1CCTL1_CPSEL                  = 0x80    # Timer 1 channel 1 capture select
	// T1CCTL1_DSM_SPD                = 0x04
	// T1CCTL1_IM                     = 0x40    # Channel 1 Interrupt mask
	// T1CCTL1_MODE                   = 0x04    # Capture or compare mode
	// T1CCTL2                        = 0xE7
	// T1CCTL2_CAP                    = 0x03
	// T1CCTL2_CAP0                   = 0x01
	// T1CCTL2_CAP1                   = 0x02
	// T1CCTL2_CMP                    = 0x38
	// T1CCTL2_CMP0                   = 0x08
	// T1CCTL2_CMP1                   = 0x10
	// T1CCTL2_CMP2                   = 0x20
	// T1CCTL2_CPSEL                  = 0x80    # Timer 1 channel 2 capture select
	// T1CCTL2_IM                     = 0x40    # Channel 2 Interrupt mask
	// T1CCTL2_MODE                   = 0x04    # Capture or compare mode
	// T1CNTH                         = 0xE3
	// T1CNTL                         = 0xE2
	// T1CTL                          = 0xE4
	// T1CTL_CH0IF                    = 0x20 # Timer 1 channel 0 interrupt flag
	// T1CTL_CH1IF                    = 0x40 # Timer 1 channel 1 interrupt flag
	// T1CTL_CH2IF                    = 0x80 # Timer 1 channel 2 interrupt flag
	// T1CTL_DIV                      = 0x0C
	// T1CTL_DIV0                     = 0x04
	// T1CTL_DIV1                     = 0x08
	// T1CTL_DIV_1                    = (0x00 << 2) # Divide tick frequency by 1
	// T1CTL_DIV_128                  = (0x03 << 2) # Divide tick frequency by 128
	// T1CTL_DIV_32                   = (0x02 << 2) # Divide tick frequency by 32
	// T1CTL_DIV_8                    = (0x01 << 2) # Divide tick frequency by 8
	// T1CTL_MODE                     = 0x03
	// T1CTL_MODE0                    = 0x01
	// T1CTL_MODE1                    = 0x02
	// T1CTL_MODE_FREERUN             = (0x01)    # Free Running mode
	// T1CTL_MODE_MODULO              = (0x02)    # Modulo
	// T1CTL_MODE_SUSPEND             = (0x00)    # Operation is suspended (halt)
	// T1CTL_MODE_UPDOWN              = (0x03)    # Up/down
	// T1CTL_OVFIF                    = 0x10 # Timer 1 counter overflow interrupt flag
	// T1IE                           = 2
	// T1IF                           = 2
	// T1_VECTOR                      = 9    #  Timer 1 (16-bit) Capture/Compare/Overflow
	// T2CT                           = 0x9C
	// T2CTL                          = 0x9E
	// T2CTL_INT                      = 0x10    # Enable Timer 2 interrupt
	// T2CTL_TEX                      = 0x40
	// T2CTL_TIG                      = 0x04    # Tick generator mode
	// T2CTL_TIP                      = 0x03
	// T2CTL_TIP0                     = 0x01
	// T2CTL_TIP1                     = 0x02
	// T2CTL_TIP_1024                 = (0x03)
	// T2CTL_TIP_128                  = (0x01)
	// T2CTL_TIP_256                  = (0x02)
	// T2CTL_TIP_64                   = (0x00)
	// T2IE                           = 4
	// T2IF                           = 4
	// T2PR                           = 0x9D
	// T2_VECTOR                      = 10   #  Timer 2 (MAC Timer) Overflow
	// T3C0_CLR_CMP_SET_0             = (0x06 << 3)  # Clear when equal to T3CC0, set on 0
	// T3C0_CLR_CMP_UP_SET_0          = (0x04 << 3)  # Clear output on compare-up set on 0
	// T3C0_CLR_ON_CMP                = (0x01 << 3)  # Clear output on compare
	// T3C0_SET_CMP_CLR_255           = (0x05 << 3)  # Set when equal to T3CC0, clear on 255
	// T3C0_SET_CMP_UP_CLR_0          = (0x03 << 3)  # Set output on compare-up clear on 0
	// T3C0_SET_ON_CMP                = (0x00 << 3)  # Set output on compare
	// T3C0_TOG_ON_CMP                = (0x02 << 3)  # Toggle output on compare
	// T3C1_CLR_CMP_SET_C0            = (0x06 << 3)  # Clear when equal to T3CC1, set when equal to T3CC0
	// T3C1_CLR_CMP_UP_SET_0          = (0x04 << 3)  # Clear output on compare-up set on 0
	// T3C1_CLR_ON_CMP                = (0x01 << 3)  # Clear output on compare
	// T3C1_SET_CMP_CLR_C0            = (0x05 << 3)  # Set when equal to T3CC1, clear when equal to T3CC0
	// T3C1_SET_CMP_UP_CLR_0          = (0x03 << 3)  # Set output on compare-up clear on 0
	// T3C1_SET_ON_CMP                = (0x00 << 3)  # Set output on compare
	// T3C1_TOG_ON_CMP                = (0x02 << 3)  # Toggle output on compare
	// T3CC0                          = 0xCD
	// T3CC1                          = 0xCF
	// T3CCTL0                        = 0xCC
	// T3CCTL0_CMP                    = 0x38
	// T3CCTL0_CMP0                   = 0x08
	// T3CCTL0_CMP1                   = 0x10
	// T3CCTL0_CMP2                   = 0x20
	// T3CCTL0_IM                     = 0x40
	// T3CCTL0_MODE                   = 0x04
	// T3CCTL1                        = 0xCE
	// T3CCTL1_CMP                    = 0x38
	// T3CCTL1_CMP0                   = 0x08
	// T3CCTL1_CMP1                   = 0x10
	// T3CCTL1_CMP2                   = 0x20
	// T3CCTL1_IM                     = 0x40
	// T3CCTL1_MODE                   = 0x04
	// T3CH0IF                        = 2
	// T3CH1IF                        = 4
	// T3CNT                          = 0xCA
	// T3CTL                          = 0xCB
	// T3CTL_CLR                      = 0x04
	// T3CTL_DIV                      = 0xE0
	// T3CTL_DIV0                     = 0x20
	// T3CTL_DIV1                     = 0x40
	// T3CTL_DIV2                     = 0x80
	// T3CTL_DIV_1                    = (0x00 << 5)
	// T3CTL_DIV_128                  = (0x07 << 5)
	// T3CTL_DIV_16                   = (0x04 << 5)
	// T3CTL_DIV_2                    = (0x01 << 5)
	// T3CTL_DIV_32                   = (0x05 << 5)
	// T3CTL_DIV_4                    = (0x02 << 5)
	// T3CTL_DIV_64                   = (0x06 << 5)
	// T3CTL_DIV_8                    = (0x03 << 5)
	// T3CTL_MODE                     = 0x03
	// T3CTL_MODE0                    = 0x01
	// T3CTL_MODE1                    = 0x02
	// T3CTL_MODE_DOWN                = (0x01)
	// T3CTL_MODE_FREERUN             = (0x00)
	// T3CTL_MODE_MODULO              = (0x02)
	// T3CTL_MODE_UPDOWN              = (0x03)
	// T3CTL_OVFIM                    = 0x08
	// T3CTL_START                    = 0x10
	// T3IE                           = 8
	// T3IF                           = 8
	// T3OVFIF                        = 1
	// T3_VECTOR                      = 11   #  Timer 3 (8-bit) Capture/Compare/Overflow
	// T4CC0                          = 0xED
	// T4CC1                          = 0xEF
	// T4CCTL0                        = 0xEC
	// T4CCTL0_CLR_CMP_SET_0          = (0x06 << 3)
	// T4CCTL0_CLR_CMP_UP_SET_0       = (0x04 << 3)
	// T4CCTL0_CLR_ON_CMP             = (0x01 << 3)
	// T4CCTL0_CMP                    = 0x38
	// T4CCTL0_CMP0                   = 0x08
	// T4CCTL0_CMP1                   = 0x10
	// T4CCTL0_CMP2                   = 0x20
	// T4CCTL0_IM                     = 0x40
	// T4CCTL0_MODE                   = 0x04
	// T4CCTL0_SET_CMP_CLR_255        = (0x05 << 3)
	// T4CCTL0_SET_CMP_UP_CLR_0       = (0x03 << 3)
	// T4CCTL0_SET_ON_CMP             = (0x00 << 3)
	// T4CCTL0_TOG_ON_CMP             = (0x02 << 3)
	// T4CCTL1                        = 0xEE
	// T4CCTL1_CLR_CMP_SET_C0         = (0x06 << 3)
	// T4CCTL1_CLR_CMP_UP_SET_0       = (0x04 << 3)
	// T4CCTL1_CLR_ON_CMP             = (0x01 << 3)
	// T4CCTL1_CMP                    = 0x38
	// T4CCTL1_CMP0                   = 0x08
	// T4CCTL1_CMP1                   = 0x10
	// T4CCTL1_CMP2                   = 0x20
	// T4CCTL1_IM                     = 0x40
	// T4CCTL1_MODE                   = 0x04
	// T4CCTL1_SET_CMP_CLR_C0         = (0x05 << 3)
	// T4CCTL1_SET_CMP_UP_CLR_0       = (0x03 << 3)
	// T4CCTL1_SET_ON_CMP             = (0x00 << 3)
	// T4CCTL1_TOG_ON_CMP             = (0x02 << 3)
	// T4CH0IF                        = 16
	// T4CH1IF                        = 32
	// T4CNT                          = 0xEA
	// T4CTL                          = 0xEB
	// T4CTL_CLR                      = 0x04
	// T4CTL_DIV                      = 0xE0
	// T4CTL_DIV0                     = 0x20
	// T4CTL_DIV1                     = 0x40
	// T4CTL_DIV2                     = 0x80
	// T4CTL_DIV_1                    = (0x00 << 5)
	// T4CTL_DIV_128                  = (0x07 << 5)
	// T4CTL_DIV_16                   = (0x04 << 5)
	// T4CTL_DIV_2                    = (0x01 << 5)
	// T4CTL_DIV_32                   = (0x05 << 5)
	// T4CTL_DIV_4                    = (0x02 << 5)
	// T4CTL_DIV_64                   = (0x06 << 5)
	// T4CTL_DIV_8                    = (0x03 << 5)
	// T4CTL_MODE                     = 0x03
	// T4CTL_MODE0                    = 0x01
	// T4CTL_MODE1                    = 0x02
	// T4CTL_MODE_DOWN                = (0x01)
	// T4CTL_MODE_FREERUN             = (0x00)
	// T4CTL_MODE_MODULO              = (0x02)
	// T4CTL_MODE_UPDOWN              = (0x03)
	// T4CTL_OVFIM                    = 0x08
	// T4CTL_START                    = 0x10
	// T4IE                           = 16
	// T4IF                           = 16
	// T4OVFIF                        = 8
	// T4_VECTOR                      = 12   #  Timer 4 (8-bit) Capture/Compare/Overflow
	// TCON                           = 0x88
	// TEST0                          = 0xDF25
	// TEST1                          = 0xDF24
	// TEST2                          = 0xDF23
	// TICKSPD_DIV_1                  = (0x00 << 3)
	// TICKSPD_DIV_128                = (0x07 << 3)
	// TICKSPD_DIV_16                 = (0x04 << 3)
	// TICKSPD_DIV_2                  = (0x01 << 3)
	// TICKSPD_DIV_32                 = (0x05 << 3)
	// TICKSPD_DIV_4                  = (0x02 << 3)
	// TICKSPD_DIV_64                 = (0x06 << 3)
	// TICKSPD_DIV_8                  = (0x03 << 3)
	// TIMIF                          = 0xD8
	// TX_BYTE                        = 2
	// U0BAUD                         = 0xC2
	// U0CSR                          = 0x86
	// U0CSR_ACTIVE                   = 0x01
	// U0CSR_ERR                      = 0x08
	// U0CSR_FE                       = 0x10
	// U0CSR_MODE                     = 0x80
	// U0CSR_RE                       = 0x40
	// U0CSR_RX_BYTE                  = 0x04
	// U0CSR_SLAVE                    = 0x20
	// U0CSR_TX_BYTE                  = 0x02
	// U0DBUF                         = 0xC1
	// U0GCR                          = 0xC5
	// U0GCR_BAUD_E                   = 0x1F
	// U0GCR_BAUD_E0                  = 0x01
	// U0GCR_BAUD_E1                  = 0x02
	// U0GCR_BAUD_E2                  = 0x04
	// U0GCR_BAUD_E3                  = 0x08
	// U0GCR_BAUD_E4                  = 0x10
	// U0GCR_CPHA                     = 0x40
	// U0GCR_CPOL                     = 0x80
	// U0GCR_ORDER                    = 0x20
	// U0UCR                          = 0xC4
	// U0UCR_BIT9                     = 0x10
	// U0UCR_D9                       = 0x20
	// U0UCR_FLOW                     = 0x40
	// U0UCR_FLUSH                    = 0x80
	// U0UCR_PARITY                   = 0x08
	// U0UCR_SPB                      = 0x04
	// U0UCR_START                    = 0x01
	// U0UCR_STOP                     = 0x02
	// U1BAUD                         = 0xFA
	// U1CSR                          = 0xF8
	// U1CSR_ACTIVE                   = 0x01
	// U1CSR_ERR                      = 0x08
	// U1CSR_FE                       = 0x10
	// U1CSR_MODE                     = 0x80
	// U1CSR_RE                       = 0x40
	// U1CSR_RX_BYTE                  = 0x04
	// U1CSR_SLAVE                    = 0x20
	// U1CSR_TX_BYTE                  = 0x02
	// U1DBUF                         = 0xF9
	// U1GCR                          = 0xFC
	// U1GCR_BAUD_E                   = 0x1F
	// U1GCR_BAUD_E0                  = 0x01
	// U1GCR_BAUD_E1                  = 0x02
	// U1GCR_BAUD_E2                  = 0x04
	// U1GCR_BAUD_E3                  = 0x08
	// U1GCR_BAUD_E4                  = 0x10
	// U1GCR_CPHA                     = 0x40
	// U1GCR_CPOL                     = 0x80
	// U1GCR_ORDER                    = 0x20
	// U1UCR                          = 0xFB
	// U1UCR_BIT9                     = 0x10
	// U1UCR_D9                       = 0x20
	// U1UCR_FLOW                     = 0x40
	// U1UCR_FLUSH                    = 0x80
	// U1UCR_PARITY                   = 0x08
	// U1UCR_SPB                      = 0x04
	// U1UCR_START                    = 0x01
	// U1UCR_STOP                     = 0x02
	// URX0IE                         = 4
	// URX0IF                         = 8
	// URX0_VECTOR                    = 2    #  USART0 RX Complete
	// URX1IE                         = 8
	// URX1IF                         = 128
	// URX1_VECTOR                    = 3    #  USART1 RX Complete
	// USBADDR                        = 0xDE00
	// USBADDR_UPDATE                 = 0x80    #r
	// USBCIE                         = 0xDE0B
	// USBCIE_RESUMEIE                = 0x02    #rw
	// USBCIE_RSTIE                   = 0x04    #rw
	// USBCIE_SOFIE                   = 0x08    #rw
	// USBCIE_SUSPENDIE               = 0x01    #rw
	// USBCIF                         = 0xDE06
	// USBCIF_RESUMEIF                = 0x02    #r h0
	// USBCIF_RSTIF                   = 0x04    #r h0
	// USBCIF_SOFIF                   = 0x08    #r h0
	// USBCIF_SUSPENDIF               = 0x01    #r h0
	// USBCNT0                        = 0xDE16
	// USBCNTH                        = 0xDE17
	// USBCNTL                        = 0xDE16
	// USBCS0                         = 0xDE11
	// USBCS0_CLR_OUTPKT_RDY          = 0x40    #rw h0
	// USBCS0_CLR_SETUP_END           = 0x80    #rw h0
	// USBCS0_DATA_END                = 0x08    #rw h0
	// USBCS0_INPKT_RDY               = 0x02    #rw h0
	// USBCS0_OUTPKT_RDY              = 0x01    #r
	// USBCS0_SEND_STALL              = 0x20    #rw h0
	// USBCS0_SENT_STALL              = 0x04    #rw h1
	// USBCS0_SETUP_END               = 0x10    #r
	// USBCSIH                        = 0xDE12
	// USBCSIH_AUTOSET                = 0x80    #rw
	// USBCSIH_FORCE_DATA_TOG         = 0x08    #rw
	// USBCSIH_IN_DBL_BUF             = 0x01    #rw
	// USBCSIH_ISO                    = 0x40    #rw
	// USBCSIL                        = 0xDE11
	// USBCSIL_CLR_DATA_TOG           = 0x40    #rw h0
	// USBCSIL_FLUSH_PACKET           = 0x08    #rw h0
	// USBCSIL_INPKT_RDY              = 0x01    #rw h0
	// USBCSIL_PKT_PRESENT            = 0x02    #r
	// USBCSIL_SEND_STALL             = 0x10    #rw
	// USBCSIL_SENT_STALL             = 0x20    #rw
	// USBCSIL_UNDERRUN               = 0x04    #rw
	// USBCSOH                        = 0xDE15
	// USBCSOH_AUTOCLEAR              = 0x80    #rw
	// USBCSOH_ISO                    = 0x40    #rw
	// USBCSOH_OUT_DBL_BUF            = 0x01    #rw
	// USBCSOL                        = 0xDE14
	// USBCSOL_CLR_DATA_TOG           = 0x80    #rw h0
	// USBCSOL_DATA_ERROR             = 0x08    #r
	// USBCSOL_FIFO_FULL              = 0x02    #r
	// USBCSOL_FLUSH_PACKET           = 0x10    #rw
	// USBCSOL_OUTPKT_RDY             = 0x01    #rw
	// USBCSOL_OVERRUN                = 0x04    #rw
	// USBCSOL_SEND_STALL             = 0x20    #rw
	// USBCSOL_SENT_STALL             = 0x40    #rw
	// USBF0                          = 0xDE20
	// USBF1                          = 0xDE22
	// USBF2                          = 0xDE24
	// USBF3                          = 0xDE26
	// USBF4                          = 0xDE28
	// USBF5                          = 0xDE2A
	// USBFRMH                        = 0xDE0D
	// USBFRML                        = 0xDE0C
	// USBIF                          = 1
	// USBIIE                         = 0xDE07
	// USBIIE_EP0IE                   = 0x01    #rw
	// USBIIE_INEP1IE                 = 0x02    #rw
	// USBIIE_INEP2IE                 = 0x04    #rw
	// USBIIE_INEP3IE                 = 0x08    #rw
	// USBIIE_INEP4IE                 = 0x10    #rw
	// USBIIE_INEP5IE                 = 0x20    #rw
	// USBIIF                         = 0xDE02
	// USBIIF_EP0IF                   = 0x01    #r h0
	// USBIIF_INEP1IF                 = 0x02    #r h0
	// USBIIF_INEP2IF                 = 0x04    #r h0
	// USBIIF_INEP3IF                 = 0x08    #r h0
	// USBIIF_INEP4IF                 = 0x10    #r h0
	// USBIIF_INEP5IF                 = 0x20    #r h0
	// USBINDEX                       = 0xDE0E
	// USBMAXI                        = 0xDE10
	// USBMAXO                        = 0xDE13
	// USBOIE                         = 0xDE09
	// USBOIE_EP1IE                   = 0x02    #rw
	// USBOIE_EP2IE                   = 0x04    #rw
	// USBOIE_EP3IE                   = 0x08    #rw
	// USBOIE_EP4IE                   = 0x10    #rw
	// USBOIE_EP5IE                   = 0x20    #rw
	// USBOIF                         = 0xDE04
	// USBOIF_OUTEP1IF                = 0x02    #r h0
	// USBOIF_OUTEP2IF                = 0x04    #r h0
	// USBOIF_OUTEP3IF                = 0x08    #r h0
	// USBOIF_OUTEP4IF                = 0x10    #r h0
	// USBOIF_OUTEP5IF                = 0x20    #r h0
	// USBPOW                         = 0xDE01
	// USBPOW_ISO_WAIT_SOF            = 0x80    #rw
	// USBPOW_RESUME                  = 0x04    #rw
	// USBPOW_RST                     = 0x08    #r
	// USBPOW_SUSPEND                 = 0x02    #r
	// USBPOW_SUSPEND_EN              = 0x01    #rw
	// USB_BM_REQTYPE_DIRMASK         = 0x80
	// USB_BM_REQTYPE_DIR_IN          = 0x80
	// USB_BM_REQTYPE_DIR_OUT         = 0x00
	// USB_BM_REQTYPE_TGTMASK         = 0x1f
	// USB_BM_REQTYPE_TGT_DEV         = 0x00
	// USB_BM_REQTYPE_TGT_EP          = 0x02
	// USB_BM_REQTYPE_TGT_INTF        = 0x01
	// USB_BM_REQTYPE_TYPEMASK        = 0x60
	// USB_BM_REQTYPE_TYPE_CLASS      = 0x20
	// USB_BM_REQTYPE_TYPE_RESERVED   = 0x60
	// USB_BM_REQTYPE_TYPE_STD        = 0x00
	// USB_BM_REQTYPE_TYPE_VENDOR     = 0x40
	// USB_CLEAR_FEATURE              = 0x01
	// USB_DESC_CONFIG                = 0x02
	// USB_DESC_DEVICE                = 0x01
	// USB_DESC_ENDPOINT              = 0x05
	// USB_DESC_INTERFACE             = 0x04
	// USB_DESC_STRING                = 0x03
	// USB_ENABLE_PIN                 = P1_0
	// USB_GET_CONFIGURATION          = 0x08
	// USB_GET_DESCRIPTOR             = 0x06
	// USB_GET_INTERFACE              = 0x0a
	// USB_GET_STATUS                 = 0x00
	// USB_SET_ADDRESS                = 0x05
	// USB_SET_CONFIGURATION          = 0x09
	// USB_SET_DESCRIPTOR             = 0x07
	// USB_SET_FEATURE                = 0x03
	// USB_SET_INTERFACE              = 0x11
	// USB_STATE_BLINK                = 0xff
	// USB_STATE_IDLE                 = 0
	// USB_STATE_RESET                = 3
	// USB_STATE_RESUME               = 2
	// USB_STATE_SUSPEND              = 1
	// USB_STATE_WAIT_ADDR            = 4
	// USB_SYNCH_FRAME                = 0x12
	// UTX0IF                         = 2
	// UTX0_VECTOR                    = 7    #  USART0 TX Complete
	// UTX1IF                         = 4
	// UTX1_VECTOR                    = 14   #  USART1 TX Complete
	// VCO_VC_DAC                     = 0xDF3D
	// VERSION                        = 0xDF37
	// WDCTL                          = 0xC9
	// WDCTL_CLR                      = 0xF0
	// WDCTL_CLR0                     = 0x10
	// WDCTL_CLR1                     = 0x20
	// WDCTL_CLR2                     = 0x40
	// WDCTL_CLR3                     = 0x80
	// WDCTL_EN                       = 0x08
	// WDCTL_INT                      = 0x03
	// WDCTL_INT0                     = 0x01
	// WDCTL_INT1                     = 0x02
	// WDCTL_INT1_MSEC_250            = (0x01)
	// WDCTL_INT2_MSEC_15             = (0x02)
	// WDCTL_INT3_MSEC_2              = (0x03)
	// WDCTL_INT_SEC_1                = (0x00)
	// WDCTL_MODE                     = 0x04
	// WDTIF                          = 16
	// WDT_VECTOR                     = 17   #  Watchdog Overflow in Timer Mode
	// WORCTL_WOR_RES                 = 0x03
	// WORCTL_WOR_RES0                = 0x01
	// WORCTL_WOR_RES1                = 0x02
	// WORCTL_WOR_RESET               = 0x04
	// WORCTL_WOR_RES_1               = (0x00)
	// WORCTL_WOR_RES_1024            = (0x02)
	// WORCTL_WOR_RES_32              = (0x01)
	// WORCTL_WOR_RES_32768           = (0x03)
	// WORCTRL                        = 0xA2
	// WOREVT0                        = 0xA3
	// WOREVT1                        = 0xA4
	// WORIRQ                         = 0xA1
	// WORIRQ_EVENT0_FLAG             = 0x01
	// WORIRQ_EVENT0_MASK             = 0x10
	// WORTIME0                       = 0xA5
	// WORTIME1                       = 0xA6
	// X_ADCCFG                       = 0xDFF2
	// X_ADCCON1                      = 0xDFB4
	// X_ADCCON2                      = 0xDFB5
	// X_ADCCON3                      = 0xDFB6
	// X_ADCH                         = 0xDFBB
	// X_ADCL                         = 0xDFBA
	// X_CLKCON                       = 0xDFC6
	// X_DMA0CFGH                     = 0xDFD5
	// X_DMA0CFGL                     = 0xDFD4
	// X_DMA1CFGH                     = 0xDFD3
	// X_DMA1CFGL                     = 0xDFD2
	// X_DMAARM                       = 0xDFD6
	// X_DMAIRQ                       = 0xDFD1
	// X_DMAREQ                       = 0xDFD7
	// X_ENCCS                        = 0xDFB3
	// X_ENCDI                        = 0xDFB1
	// X_ENCDO                        = 0xDFB2
	// X_FADDRH                       = 0xDFAD
	// X_FADDRL                       = 0xDFAC
	// X_FCTL                         = 0xDFAE
	// X_FWDATA                       = 0xDFAF
	// X_FWT                          = 0xDFAB
	// X_MEMCTR                       = 0xDFC7
	// X_MPAGE                        = 0xDF93
	// X_P0DIR                        = 0xDFFD
	// X_P0IFG                        = 0xDF89
	// X_P0INP                        = 0xDF8F
	// X_P0SEL                        = 0xDFF3
	// X_P1DIR                        = 0xDFFE
	// X_P1IEN                        = 0xDF8D
	// X_P1IFG                        = 0xDF8A
	// X_P1INP                        = 0xDFF6
	// X_P1SEL                        = 0xDFF4
	// X_P2DIR                        = 0xDFFF
	// X_P2IFG                        = 0xDF8B
	// X_P2INP                        = 0xDFF7
	// X_P2SEL                        = 0xDFF5
	// X_PERCFG                       = 0xDFF1
	// X_PICTL                        = 0xDF8C
	// X_RFD                          = 0xDFD9
	// X_RFIF                         = 0xDFE9
	// X_RFIM                         = 0xDF91
	xRFST = 0xDFE1

	// X_RNDH                         = 0xDFBD
	// X_RNDL                         = 0xDFBC
	// X_SLEEP                        = 0xDFBE
	// X_T1CC0H                       = 0xDFDB
	// X_T1CC0L                       = 0xDFDA
	// X_T1CC1H                       = 0xDFDD
	// X_T1CC1L                       = 0xDFDC
	// X_T1CC2H                       = 0xDFDF
	// X_T1CC2L                       = 0xDFDE
	// X_T1CCTL0                      = 0xDFE5
	// X_T1CCTL1                      = 0xDFE6
	// X_T1CCTL2                      = 0xDFE7
	// X_T1CNTH                       = 0xDFE3
	// X_T1CNTL                       = 0xDFE2
	// X_T1CTL                        = 0xDFE4
	// X_T2CT                         = 0xDF9C
	// X_T2CTL                        = 0xDF9E
	// X_T2PR                         = 0xDF9D
	// X_T3CC0                        = 0xDFCD
	// X_T3CC1                        = 0xDFCF
	// X_T3CCTL0                      = 0xDFCC
	// X_T3CCTL1                      = 0xDFCE
	// X_T3CNT                        = 0xDFCA
	// X_T3CTL                        = 0xDFCB
	// X_T4CC0                        = 0xDFED
	// X_T4CC1                        = 0xDFEF
	// X_T4CCTL0                      = 0xDFEC
	// X_T4CCTL1                      = 0xDFEE
	// X_T4CNT                        = 0xDFEA
	// X_T4CTL                        = 0xDFEB
	// X_TIMIF                        = 0xDFD8
	// X_U0BAUD                       = 0xDFC2
	// X_U0CSR                        = 0xDF86
	// X_U0DBUF                       = 0xDFC1
	// X_U0GCR                        = 0xDFC5
	// X_U0UCR                        = 0xDFC4
	// X_U1BAUD                       = 0xDFFA
	// X_U1CSR                        = 0xDFF8
	// X_U1DBUF                       = 0xDFF9
	// X_U1GCR                        = 0xDFFC
	// X_U1UCR                        = 0xDFFB
	// X_WDCTL                        = 0xDFC9
	// X_WORCTRL                      = 0xDFA2
	// X_WOREVT0                      = 0xDFA3
	// X_WOREVT1                      = 0xDFA4
	// X_WORIRQ                       = 0xDFA1
	// X_WORTIME0                     = 0xDFA5
	// X_WORTIME1                     = 0xDFA6
	// _NA_ACC                        = 0xDFE0
	// _NA_B                          = 0xDFF0
	// _NA_DPH0                       = 0xDF83
	// _NA_DPH1                       = 0xDF85
	// _NA_DPL0                       = 0xDF82
	// _NA_DPL1                       = 0xDF84
	// _NA_DPS                        = 0xDF92
	// _NA_IEN0                       = 0xDFA8
	// _NA_IEN1                       = 0xDFB8
	// _NA_IEN2                       = 0xDF9A
	// _NA_IP0                        = 0xDFA9
	// _NA_IP1                        = 0xDFB9
	// _NA_IRCON                      = 0xDFC0
	// _NA_IRCON2                     = 0xDFE8
	// _NA_P0                         = 0xDF80
	// _NA_P1                         = 0xDF90
	// _NA_P2                         = 0xDFA0
	// _NA_PCON                       = 0xDF87
	// _NA_PSW                        = 0xDFD0
	// _NA_S0CON                      = 0xDF98
	// _NA_S1CON                      = 0xDF9B
	// _NA_SP                         = 0xDF81
	// _NA_TCON                       = 0xDF88
	// _SFR8E                         = 0x8E
	// _SFR94                         = 0x94
	// _SFR95                         = 0x95
	// _SFR96                         = 0x96
	// _SFR97                         = 0x97
	// _SFR99                         = 0x99
	// _SFR9F                         = 0x9F
	// _SFRA7                         = 0xA7
	// _SFRAA                         = 0xAA
	// _SFRB0                         = 0xB0
	// _SFRB7                         = 0xB7
	// _SFRBF                         = 0xBF
	// _SFRC3                         = 0xC3
	// _SFRC8                         = 0xC8
	// _XPAGE                         = 0x93
	// _XREGDF20                      = 0xDF20
	// _XREGDF21                      = 0xDF21
	// _XREGDF22                      = 0xDF22
	// _XREGDF26                      = 0xDF26
	// _XREGDF32                      = 0xDF32
	// _XREGDF33                      = 0xDF33
	// _XREGDF34                      = 0xDF34
	// _XREGDF35                      = 0xDF35
	// _X_SFR8E                       = 0xDF8E
	// _X_SFR94                       = 0xDF94
	// _X_SFR95                       = 0xDF95
	// _X_SFR96                       = 0xDF96
	// _X_SFR97                       = 0xDF97
	// _X_SFR99                       = 0xDF99
	// _X_SFR9F                       = 0xDF9F
	// _X_SFRA7                       = 0xDFA7
	// _X_SFRAA                       = 0xDFAA
	// _X_SFRB0                       = 0xDFB0
	// _X_SFRB7                       = 0xDFB7
	// _X_SFRBF                       = 0xDFBF
	// _X_SFRC3                       = 0xDFC3
	// _X_SFRC8                       = 0xDFC8

	// # AES co-processor
	// # defines for specifying desired crypto operations.
	// # AES_CRYPTO is in two halves:
	// #    upper 4 bits mirror CC1111 mode (ENCCS_MODE_CBC etc.)
	// #    lower 4 bits are switches
	// # AES_CRYPTO[7:4]     ENCCS_MODE...
	// # AES_CRYPTO[3]       OUTBOUND 0 == OFF, 1 == ON
	// # AES_CRYPTO[2]       OUTBOUND 0 == Decrypt, 1 == Encrypt
	// # AES_CRYPTO[1]       INBOUND  0 == OFF, 1 == ON
	// # AES_CRYPTO[0]       INBOUND  0 == Decrypt, 1 == Encrypt
	// # bitfields
	aesCryptoMode       = 0xF0
	aesCryptoOut        = 0x0C
	aesCryptoOutEnable  = 0x08
	aesCryptoOutOn      = 0x01 << 3
	aesCryptoOutOff     = 0x00 << 3
	aesCryptoOutType    = 0x04
	aesCryptoOutDecrypt = 0x00 << 2
	aesCryptoOutEncrypt = 0x01 << 2
	aesCryptoIN         = 0x03
	aesCryptoInEnable   = 0x02
	aesCryptoInOn       = 0x01 << 1
	aesCryptoInOff      = 0x00 << 1
	aesCryptoInType     = 0x01
	aesCryptoInDecrypt  = 0x00 << 0
	aesCryptoInEncrypt  = 0x01 << 0
	aesCryptoNone       = 0x00
	aesCryptoDefault    = ENCCSModeCBC | aesCryptoOutOn | aesCryptoOutEncrypt | aesCryptoInOn | aesCryptoInDecrypt

// # flags
// AES_DISABLE               = 0x00
// AES_ENABLE                = 0x01
// AES_DECRYPT               = 0x00
// AES_ENCRYPT               = 0x01

// ADCCON1S = {}
// ADCCON2S = {}
// ADCCON3S = {}
// AGCCTRL0S = {}
// AGCCTRL1S = {}
// BSCFGS = {}
// CLKCONS = {}
// CLKSPDS = {}
// DEVIATNS = {}
// IEN0S = {}
// IEN1S = {}
// IEN2S = {}
// IOCFG0S = {}
// IOCFG1S = {}
// IOCFG2S = {}
// MARC_STATES = {}
// MCSM0S = {}
// MCSM1S = {}
// MCSM2S = {}
// MDMCFG2S = {}
// PKTCTRL0S = {}
// PKTCTRL1S = {}
// RFIFS = {}
// RFIMS = {}

// for key,val in globals().items():
//     if key.startswith("RFIF_"):
//         RFIFS[val] = key
//     elif key.startswith("RFIM_"):
//         RFIMS[val] = key
//     elif key.startswith("ADCCON1_"):
//         ADCCON1S[val] = key
//     elif key.startswith("ADCCON2_"):
//         ADCCON2S[val] = key
//     elif key.startswith("ADCCON3_"):
//         ADCCON3S[val] = key
//     elif key.startswith("AGCCTRL0_"):
//         AGCCTRL0S[val] = key
//     elif key.startswith("AGCCTRL1_"):
//         AGCCTRL1S[val] = key
//     elif key.startswith("BSCFG_"):
//         BSCFGS[val] = key
//     elif key.startswith("CLKCON_"):
//         CLKCONS[val] = key
//     elif key.startswith("CLKSPD_"):
//         CLKSPDS[val] = key
//     elif key.startswith("DEVIATN_"):
//         DEVIATNS[val] = key
//     elif key.startswith("IEN0_"):
//         IEN0S[val] = key
//     elif key.startswith("IEN1_"):
//         IEN1S[val] = key
//     elif key.startswith("IEN2_"):
//         IEN2S[val] = key
//     elif key.startswith("IOCFG0_"):
//         IOCFG0S[val] = key
//     elif key.startswith("IOCFG1_"):
//         IOCFG1S[val] = key
//     elif key.startswith("IOCFG2_"):
//         IOCFG2S[val] = key
//     elif key.startswith("MARC_STATE_"):
//         MARC_STATES[val] = key
//     elif key.startswith("MCSM0_"):
//         MCSM0S[val] = key
//     elif key.startswith("MCSM1_"):
//         MCSM1S[val] = key
//     elif key.startswith("MCSM2_"):
//         MCSM2S[val] = key
//     elif key.startswith("MDMCFG2_"):
//         MDMCFG2S[val] = key
//     elif key.startswith("PKTCTRL0_"):
//         PKTCTRL0S[val] = key
//     elif key.startswith("PKTCTRL1_"):
//         PKTCTRL1S[val] = key
)
