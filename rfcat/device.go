package rfcat

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/gousb"
)

const (
	appGeneric = 0x01
	appDebug   = 0xfe
	appSystem  = 0xff
	appNIC     = 0x42
	appSpecan  = 0x43

	nicRecv         = 0x1
	nicXmit         = 0x2
	nicSetID        = 0x3
	nicSetRecvLarge = 0x5
	nicSetAESMode   = 0x6
	nicGetAESMode   = 0x7
	nicSetAESIV     = 0x8
	nicSetAESKey    = 0x9
	nicSetAMPMode   = 0xa
	nicGetAMPMode   = 0xb
	nicXmitLong     = 0xc
	nicXmitLongMore = 0xd

	sysCmdPeek       = 0x80
	sysCmdPoke       = 0x81
	sysCmdPing       = 0x82
	sysCmdStatus     = 0x83
	sysCmdPokeReg    = 0x84
	sysCmdGetClock   = 0x85
	sysCmdBuildtype  = 0x86
	sysCmdBootloader = 0x87
	sysCmdRfmode     = 0x88
	sysCmdCompiler   = 0x89
	sysCmdPartnum    = 0x8e
	sysCmdReset      = 0x8f
	sysCmdClearCodes = 0x90
)

const (
	rfMaxTXBlock = 255
	rfMaxTXChunk = 240 // must match MAX_TX_MSGLEN in firmware/include/FHSS.h and be divisible by 16 for crypto operations
	rfMaxTXLong  = 65535
	rfMaxRXBlock = 512 // must match BUFFER_SIZE definition in firmware/include/cc1111rf.h
)

const (
	ep5outMaxPacketSize = 64
	ep5inMaxPacketSize  = 64
	// ep5outBufferSize must match firmware/include/chipcon_usb.h definition
	ep5outBufferSize = 516
)

const rxChannelSize = 32

var ErrNoDevice = errors.New("rfcat: device not available")

type DeviceDesc struct {
	Bus       int
	Address   int
	VendorID  int
	ProductID int
}

type Device struct {
	mhz     int
	chip    Chip
	ctx     *gousb.Context
	dev     *gousb.Device
	cfg     *gousb.Config
	ifc     *gousb.Interface
	ine     *gousb.InEndpoint
	ins     *gousb.ReadStream
	oute    *gousb.OutEndpoint
	closeCh chan struct{}
	closeWG sync.WaitGroup
	rxMu    sync.Mutex
	rxQu    map[uint16]chan []byte
	txMu    sync.Mutex
	rcfgMu  sync.RWMutex
	rcfg    []byte
	rfmode  byte
}

func Open(index int) (*Device, error) {
	devs, err := ListDevices()
	if err != nil {
		return nil, err
	}
	if index >= len(devs) {
		return nil, ErrNoDevice
	}
	return devs[index].Open()
}

func ListDevices() ([]*DeviceDesc, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	inBootloaderMode := false

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		switch desc.Vendor {
		case 0x0451:
			if desc.Product == 0x4715 {
				return true
			}
		case 0x1d50:
			switch desc.Product {
			case 0x6047, 0x6048, 0x605b:
				return true
			case 0x6049, 0x604a:
				inBootloaderMode = true
				return false
			}
		}
		return false
	})
	// All Devices returned from OpenDevices must be closed.
	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if inBootloaderMode {
		return nil, errors.New("rfcat: already in bootloader mode")
	}

	devDesc := make([]*DeviceDesc, len(devs))
	for i, d := range devs {
		devDesc[i] = &DeviceDesc{
			Bus:       d.Desc.Bus,
			Address:   d.Desc.Address,
			VendorID:  int(d.Desc.Vendor),
			ProductID: int(d.Desc.Product),
		}
		// s, err := d.SerialNumber()
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Println(s)
	}
	return devDesc, nil
}

func (dd *DeviceDesc) Open() (*Device, error) {
	usbCtx := gousb.NewContext()
	devs, err := usbCtx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Bus == dd.Bus && desc.Address == dd.Address
	})
	if err != nil {
		// Just in case any devices were opened
		for _, d := range devs {
			d.Close()
		}
		usbCtx.Close()
		return nil, err
	}
	if len(devs) == 0 {
		usbCtx.Close()
		return nil, ErrNoDevice
	}
	// Sanity check
	if len(devs) != 1 {
		for _, d := range devs {
			d.Close()
		}
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: expected to open 1 devices not %d", len(devs))
	}
	dev := devs[0]
	// if err := dev.SetAutoDetach(true); err != nil {
	// 	dev.Close()
	// 	usbCtx.Close()
	// 	return nil, fmt.Errorf("rfcat: reset failed: %s", err)
	// }
	if err := dev.Reset(); err != nil {
		dev.Close()
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: reset failed: %s", err)
	}
	cfg, err := dev.Config(1)
	if err != nil {
		dev.Close()
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: failed to get config 1: %s", err)
	}
	ifc, err := cfg.Interface(0, 0)
	if err != nil {
		cfg.Close()
		dev.Close()
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: failed to get interface 0: %s", err)
	}
	ine, err := ifc.InEndpoint(0x85)
	if err != nil {
		ifc.Close()
		cfg.Close()
		dev.Close()
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: failed to get in endpoint: %s", err)
	}
	oute, err := ifc.OutEndpoint(0x05)
	if err != nil {
		ifc.Close()
		cfg.Close()
		dev.Close()
		usbCtx.Close()
		return nil, fmt.Errorf("rfcat: failed to get out endpoint: %s", err)
	}
	// ins, err := ine.NewStream(ine.Desc.MaxPacketSize, 4)
	// if err != nil {
	// 	ifc.Close()
	// 	cfg.Close()
	// 	dev.Close()
	// 	usbCtx.Close()
	// 	return nil, fmt.Errorf("rfcat: failed to get out endpoint: %s", err)
	// }
	d := &Device{
		ctx:  usbCtx,
		dev:  dev,
		cfg:  cfg,
		ifc:  ifc,
		ine:  ine,
		oute: oute,
		// ins:     ins,
		closeCh: make(chan struct{}),
		rxQu:    make(map[uint16]chan []byte),
	}
	// Make sure the number added to the close waitgroup matches how many routines we'll wait on when closing the device.
	d.closeWG.Add(1)
	go d.recvLoop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chip, err := d.partNum(ctx)
	if err != nil {
		d.Close()
		return nil, err
	}
	d.chip = chip
	d.mhz = chip.Mhz()
	if err := d.updateRadioConfig(ctx); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

func (d *Device) Close() error {
	// Wait for goroutines to stop or timeout
	chWait := make(chan struct{})
	go func() {
		// TODO: maybe wait for writes to clear d.txMu.Lock()
		d.closeWG.Wait()
		close(chWait)
	}()
	close(d.closeCh)
	select {
	case <-chWait:
	case <-time.After(time.Second * 2):
		println("Timeout")
	}
	// Close all receive mboxes
	d.rxMu.Lock()
	for _, ch := range d.rxQu {
		close(ch)
	}
	d.rxMu.Unlock()
	d.ifc.Close() // No error returned
	errCfg := d.cfg.Close()
	errDev := d.dev.Close()
	errCtx := d.ctx.Close()
	// TODO: return all non-nil errors
	if errCfg != nil {
		return errCfg
	}
	if errDev != nil {
		return errDev
	}
	if errCtx != nil {
		return errCtx
	}
	return nil
}

func (d *Device) Ping(ctx context.Context, count int, buf []byte) (time.Duration, error) {
	if count <= 0 {
		count = 10
	}
	if buf == nil {
		buf = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	start := time.Now()
	for x := 0; x < count; x++ {
		r, err := d.send(ctx, appSystem, sysCmdPing, buf)
		if err != nil {
			return 0, err
		}
		if !bytes.Equal(buf, r) {
			return 0, fmt.Errorf("rfcat: ping response does not match")
		}
	}
	return time.Since(start), nil
}

func (d *Device) MdmModulation() Modulation {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	return Modulation(d.rcfg[cfgMdmcfg2] & mdmcfg2ModFormat)
}

func (d *Device) MdmChanSpacing() float64 {
	d.rcfgMu.RLock()
	m := d.rcfg[cfgMdmcfg0]
	e := d.rcfg[cfgMdmcfg1] & 3
	d.rcfgMu.RUnlock()
	return 1000000.0 * float64(d.mhz) / (1 << 18) * float64(256+int(m)) * float64(uint(1)<<uint(e))
}

// SetMdmChanSpacing calculates the appropriate exponent and mantissa and updates the correct registers
// chanspc is in kHz
func (d *Device) SetMdmChanSpacing(ctx context.Context, chanspc float64) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()

	// calculates the appropriate exponent and mantissa and updates the correct registers
	// chanspc is in kHz.  if you prefer, you may set the chanspc_m and chanspc_e settings
	// directly.
	// only use one or the other:
	// * chanspc
	// * chanspc_m and chanspc_e
	// def setMdmChanSpc(self, chanspc=None, chanspc_m=None, chanspc_e=None, mhz=24, radiocfg=None):
	// if (chanspc != None):
	// 	for e in range(4):
	// 		m = int(((chanspc * pow(2,18) / (1000000.0 * mhz * pow(2,e)))-256) +.5)    # rounded evenly
	// 		if m < 256:
	// 			chanspc_e = e
	// 			chanspc_m = m
	// 			break
	// if chanspc_e is None or chanspc_m is None:
	// 	raise(Exception("ChanSpc does not translate into acceptable parameters.  Should you be changing this?"))
	// #chanspc = 1000000.0 * mhz/pow(2,18) * (256 + chanspc_m) * pow(2, chanspc_e)
	// #print "chanspc_e: %x   chanspc_m: %x   chanspc: %f hz" % (chanspc_e, chanspc_m, chanspc)
	// radiocfg.mdmcfg1 &= ~MDMCFG1_CHANSPC_E            # clear out old exponent value
	// radiocfg.mdmcfg1 |= chanspc_e
	// radiocfg.mdmcfg0 = chanspc_m
	// self.setRFRegister(MDMCFG1, (radiocfg.mdmcfg1))
	// self.setRFRegister(MDMCFG0, (radiocfg.mdmcfg0))
	return errors.New("rfcat: SetMdmChanSpacing not implemented")
}

func (d *Device) FreqOffsetEst() int {
	d.rcfgMu.RLock()
	freqEst := d.rcfg[cfgFreqest]
	d.rcfgMu.RUnlock()
	return int(freqEst)
}

func (d *Device) Channel() int {
	d.rcfgMu.RLock()
	c := d.rcfg[cfgChannr]
	d.rcfgMu.RUnlock()
	return int(c)
}

func (d *Device) SetChannel(ctx context.Context, c int) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgChannr] = byte(c)
	return d.setRFRegister(ctx, regChannr, d.rcfg[cfgChannr], false)
}

func (d *Device) FsIF() int {
	d.rcfgMu.RLock()
	freqIF := int64(d.rcfg[cfgFsctrl1]&0x1f) * 1000000 * int64(d.mhz) / (1 << 10)
	d.rcfgMu.RUnlock()
	return int(freqIF)
}

func (d *Device) FsOffset() int {
	d.rcfgMu.RLock()
	fo := d.rcfg[cfgFsctrl0]
	d.rcfgMu.RUnlock()
	return int(fo)
}

func (d *Device) MdmDeviatn() float64 {
	d.rcfgMu.RLock()
	e := d.rcfg[cfgDeviatn] >> 4
	m := d.rcfg[cfgDeviatn] & deviatnDeviationM
	d.rcfgMu.RUnlock()
	return 1000000.0 * float64(d.mhz) * (8 + float64(m)) * float64(uint(1)<<e) / (1 << 17)
}

// SetMdmDeviatn configure the deviation settings for the given modulation scheme
func (d *Device) SetMdmDeviatn(ctx context.Context, hz int) (float64, error) {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	devE := -1
	devM := -1
	for e := 0; e < 8; e++ {
		m := int((float64(hz)*(1<<17)/(float64(uint(1)<<uint(e))*(float64(d.mhz)*1000000.0)) - 8) + .5) // rounded evenly
		if m < 8 {
			devE = e
			devM = m
			break
		}
	}
	if devE < 0 {
		return 0, errors.New("rfcat: deviation does not translate into acceptable parameters")
	}

	d.rcfg[cfgDeviatn] = byte((devE << 4) | devM)
	if err := d.setRFRegister(ctx, regDeviatn, d.rcfg[cfgDeviatn], false); err != nil {
		return 0, err
	}
	return 1000000.0 * float64(d.mhz) * (8 + float64(devM)) * float64(uint(1)<<uint(devE)) / (1 << 17), nil
}

func (d *Device) MdmChanBW() float64 {
	d.rcfgMu.RLock()
	e := (d.rcfg[cfgMdmcfg4] >> 6) & 3
	m := (d.rcfg[cfgMdmcfg4] >> 4) & 3
	d.rcfgMu.RUnlock()
	return 1000000.0 * float64(d.mhz) / (8.0 * (4 + float64(m)) * float64(uint(1)<<e))
}

// TODO
// def setMdmChanBW(self, bw, mhz=24, radiocfg=None):
//     '''
//     For best performance, the channel filter
//     bandwidth should be selected so that the
//     signal bandwidth occupies at most 80% of the
//     channel filter bandwidth. The channel centre
//     tolerance due to crystal accuracy should also
//     be subtracted from the signal bandwidth. The
//     following example illustrates this:

//         With the channel filter bandwidth set to 500
//         kHz, the signal should stay within 80% of 500
//         kHz, which is 400 kHz. Assuming 915 MHz
//         frequency and +/-20 ppm frequency uncertainty
//         for both the transmitting device and the
//         receiving device, the total frequency
//         uncertainty is +/-40 ppm of 915 MHz, which is
//         +/-37 kHz. If the whole transmitted signal
//         bandwidth is to be received within 400 kHz, the
//         transmitted signal bandwidth should be
//         maximum 400 kHz - 2*37 kHz, which is 326
//         kHz.

//     DR:1.2kb Dev:5.1khz Mod:GFSK RXBW:63kHz sensitive     fsctrl1:06 mdmcfg:e5 a3 13 23 11 dev:16 foc/bscfg:17/6c agctrl:03 40 91 frend:56 10
//     DR:1.2kb Dev:5.1khz Mod:GFSK RXBW:63kHz lowpower      fsctrl1:06 mdmcfg:e5 a3 93 23 11 dev:16 foc/bscfg:17/6c agctrl:03 40 91 frend:56 10    (DEM_DCFILT_OFF)
//     DR:2.4kb Dev:5.1khz Mod:GFSK RXBW:63kHz sensitive     fsctrl1:06 mdmcfg:e6 a3 13 23 11 dev:16 foc/bscfg:17/6c agctrl:03 40 91 frend:56 10
//     DR:2.4kb Dev:5.1khz Mod:GFSK RXBW:63kHz lowpower      fsctrl1:06 mdmcfg:e6 a3 93 23 11 dev:16 foc/bscfg:17/6c agctrl:03 40 91 frend:56 10    (DEM_DCFILT_OFF)
//     DR:38.4kb Dev:20khz Mod:GFSK RXBW:94kHz sensitive     fsctrl1:08 mdmcfg:ca a3 13 23 11 dev:36 foc/bscfg:16/6c agctrl:43 40 91 frend:56 10    (IF changes, Deviation)
//     DR:38.4kb Dev:20khz Mod:GFSK RXBW:94kHz lowpower      fsctrl1:08 mdmcfg:ca a3 93 23 11 dev:36 foc/bscfg:16/6c agctrl:43 40 91 frend:56 10    (.. DEM_DCFILT_OFF)

//     DR:250kb Dev:129khz Mod:GFSK RXBW:600kHz sensitive    fsctrl1:0c mdmcfg:1d 55 13 23 11 dev:63 foc/bscfg:1d/1c agctrl:c7 00 b0 frend:b6 10    (IF_changes, Deviation)

//     DR:500kb            Mod:MSK  RXBW:750kHz sensitive    fsctrl1:0e mdmcfg:0e 55 73 43 11 dev:00 foc/bscfg:1d/1c agctrl:c7 00 b0 frend:b6 10    (IF_changes, Modulation of course, Deviation has different meaning with MSK)
//     '''
//     if radiocfg==None:
//         self.getRadioConfig()
//         radiocfg = self.radiocfg

//     chanbw_e = None
//     chanbw_m = None
//     for e in range(4):
//         m = int(((mhz*1000000.0 / (bw *pow(2,e) * 8.0 )) - 4) + .5)        # rounded evenly
//         if m < 4:
//             chanbw_e = e
//             chanbw_m = m
//             break
//     if chanbw_e is None:
//         raise(Exception("ChanBW does not translate into acceptable parameters.  Should you be changing this?"))

//     bw = 1000.0*mhz / (8.0*(4+chanbw_m) * pow(2,chanbw_e))
//     #print "chanbw_e: %x   chanbw_m: %x   chanbw: %f kHz" % (e, m, bw)

//     radiocfg.mdmcfg4 &= ~(MDMCFG4_CHANBW_E | MDMCFG4_CHANBW_M)
//     radiocfg.mdmcfg4 |= ((chanbw_e<<6) | (chanbw_m<<4))
//     self.setRFRegister(MDMCFG4, (radiocfg.mdmcfg4))

//     # from http://www.cs.jhu.edu/~carlson/download/datasheets/ask_ook_settings.pdf
//     if bw > 102e3:
//         self.setRFRegister(FREND1, 0xb6)
//     else:
//         self.setRFRegister(FREND1, 0x56)

//     if bw > 325e3:
//         self.setRFRegister(TEST2, 0x88)
//         self.setRFRegister(TEST1, 0x31)
//     else:
//         self.setRFRegister(TEST2, 0x81)
//         self.setRFRegister(TEST1, 0x35)

// SetMdmNumPreamble sets the minimum number of preamble bits to be transmitted (default: MFMCFG1_NUM_PREAMBLE_4)
func (d *Device) SetMdmNumPreamble(ctx context.Context, n NumPreamble) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgMdmcfg1] &^= mfmcfg1NumPreamble
	d.rcfg[cfgMdmcfg1] |= byte(n)
	return d.setRFRegister(ctx, mdmcfg1, d.rcfg[cfgMdmcfg1], false)
}

// MdmNumPreamble returns the minimum number of preamble bits to be transmitted.
func (d *Device) MdmNumPreamble() NumPreamble {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgMdmcfg1] & mfmcfg1NumPreamble
	d.rcfgMu.RUnlock()
	// TODO: shifted by 4 bits??
	return NumPreamble(v)
}

// BSLimit returns the saturation point for the data rate offset compensation algorithm.
func (d *Device) BSLimit() BSLimit {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgBscfg] & bscfgBSLimit
	d.rcfgMu.RUnlock()
	return BSLimit(v)
}

// MdmDCFilterEnabled returns true if the DC filter is enabled.
func (d *Device) MdmDCFilterEnabled() bool {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgMdmcfg2]&0x80 == 0
	d.rcfgMu.RUnlock()
	return v
}

func (d *Device) SetMdmDCFilterEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgMdmcfg2] &^= mdmcfg2DemDCFiltOff
	} else {
		d.rcfg[cfgMdmcfg2] |= mdmcfg2DemDCFiltOff
	}
	v := d.rcfg[cfgMdmcfg2]
	return d.setRFRegister(ctx, mdmcfg2, v, false)
}

// SetMdmManchesterEnabled enables or disables manchester coding.
func (d *Device) SetMdmManchesterEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgMdmcfg2] |= mdmcfg2ManchesterEn
	} else {
		d.rcfg[cfgMdmcfg2] &^= mdmcfg2ManchesterEn
	}
	return d.setRFRegister(ctx, mdmcfg2, d.rcfg[cfgMdmcfg2], false)
}

// MdmManchesterEnabled returns true if manchester coding is enable
func (d *Device) MdmManchesterEnabled() bool {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgMdmcfg2] >> 3
	d.rcfgMu.RUnlock()
	return v&1 != 0
}

// PktDataWhiteningEnabled returns true if data whitening is enabledc
func (d *Device) PktDataWhiteningEnabled() bool {
	d.rcfgMu.RLock()
	v := (d.rcfg[cfgPktctrl0]>>6)&1 != 0
	d.rcfgMu.RUnlock()
	return v
}

func (d *Device) SetPktDataWhiteningEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgPktctrl0] |= pktctrl0WhiteData
	} else {
		d.rcfg[cfgPktctrl0] &^= pktctrl0WhiteData
	}
	return d.setRFRegister(ctx, pktctrl0Reg, d.rcfg[cfgPktctrl0], false)
}

// MdmFECEnabled returns true if FEC is enabled
func (d *Device) MdmFECEnabled() bool {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgMdmcfg1] >> 7
	d.rcfgMu.RUnlock()
	return v&1 != 0
}

func (d *Device) SetMdmFECEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgMdmcfg1] |= mfmcfg1FECEn
	} else {
		d.rcfg[cfgMdmcfg1] &^= mfmcfg1FECEn
	}
	return d.setRFRegister(ctx, mdmcfg1, d.rcfg[cfgMdmcfg1], false)
}

func (d *Device) configAddr() byte {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgAddr]
	d.rcfgMu.RUnlock()
	return v
}

func (d *Device) PktAppendStatusEnabled() bool {
	d.rcfgMu.RLock()
	v := (d.rcfg[cfgPktctrl1]>>2)&1 != 0
	d.rcfgMu.RUnlock()
	return v
}

// SetPktAppendStatusEnabled enables or disables append status bytes. two bytes will be
// appended to the payload of the packet, containing RSSI and LQI values as well as CRC OK.
func (d *Device) SetPktAppendStatusEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgPktctrl1] |= pktctrl1AppendStatus
	} else {
		d.rcfg[cfgPktctrl1] &^= pktctrl1AppendStatus
	}
	return d.setRFRegister(ctx, pktctrl1Reg, d.rcfg[cfgPktctrl1], false)
}

func (d *Device) AddressCheck() AddressCheckType {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgPktctrl1] & 3
	d.rcfgMu.RUnlock()
	return AddressCheckType(v)
}

func (d *Device) packetFormat() PacketFormat {
	d.rcfgMu.RLock()
	v := (d.rcfg[cfgPktctrl0] >> 5) & 3
	d.rcfgMu.RUnlock()
	return PacketFormat(v)
}

func (d *Device) SetPktCRCEnabled(ctx context.Context, enabled bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	if enabled {
		d.rcfg[cfgPktctrl0] |= pktctrl0CRCEn
	} else {
		d.rcfg[cfgPktctrl0] &^= pktctrl0CRCEn
	}
	return d.setRFRegister(ctx, pktctrl0Reg, d.rcfg[cfgPktctrl0], false)
}

func (d *Device) PktCRCEnabled() bool {
	d.rcfgMu.RLock()
	v := (d.rcfg[cfgPktctrl0]>>2)&1 != 0
	d.rcfgMu.RUnlock()
	return v
}

func (d *Device) Freq() float64 {
	d.rcfgMu.RLock()
	num := d.freqNum()
	d.rcfgMu.RUnlock()
	return float64(num) / d.freqMult()
}

func (d *Device) freqNum() int {
	return (int(d.rcfg[cfgFreq2]) << 16) + (int(d.rcfg[cfgFreq1]) << 8) + int(d.rcfg[cfgFreq0])
}

// SetFreq sets the frequency and returns what the actual frequency that was tuned to which
// may differ from the frequency requested.
func (d *Device) SetFreq(ctx context.Context, freq float64) (float64, error) {
	num := int(freq * d.freqMult())

	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgFreq2] = byte(num >> 16)
	d.rcfg[cfgFreq1] = byte((num >> 8) & 0xff)
	d.rcfg[cfgFreq0] = byte(num & 0xff)

	if (freq > freqEdge900 && freq < freqMid900) || (freq > freqEdge400 && freq < freqMid400) || (freq < freqMid300) {
		// select low VCO
		d.rcfg[cfgFscal2] = 0x0a
	} else if freq < 1e9 && ((freq > freqMid900) || (freq > freqMid400) || (freq > freqMid300)) {
		// select high VCO
		d.rcfg[cfgFscal2] = 0x2a
	}

	marcstate := d.rcfg[cfgMarcstate]
	if marcstate != marcStateIdle {
		if err := d.strobeModeIDLE(ctx); err != nil {
			return 0, err
		}
	}
	if _, err := d.poke(ctx, freq2, []byte{d.rcfg[cfgFreq2], d.rcfg[cfgFreq1], d.rcfg[cfgFreq0]}); err != nil {
		return 0, err
	}
	if _, err := d.poke(ctx, fscal2, []byte{d.rcfg[cfgFscal2]}); err != nil {
		return 0, err
	}

	if err := d.strobeModeReturn(ctx, marcstate); err != nil {
		return 0, err
	}
	// #if (radiocfg.marcstate == MARC_STATE_RX):
	// 	#self.strobeModeRX()
	// #elif (radiocfg.marcstate == MARC_STATE_TX):
	// 	#self.strobeModeTX()

	return float64(d.freqNum()) / d.freqMult(), nil
}

func (d *Device) SetMdmModulation(ctx context.Context, mod Modulation, invert bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()

	d.rcfg[cfgMdmcfg2] &^= mdmcfg2ModFormat
	d.rcfg[cfgMdmcfg2] |= byte(mod)

	power := -1
	// ASK_OOK needs to flip power table
	if mod == ModASKOOK && !invert {
		if d.rcfg[cfgPaTable1] == 0x00 && d.rcfg[cfgPaTable0] != 0x00 {
			power = int(d.rcfg[cfgPaTable0])
		}
	} else {
		if d.rcfg[cfgPaTable0] == 0x00 && d.rcfg[cfgPaTable1] != 0x00 {
			power = int(d.rcfg[cfgPaTable1])
		}
	}

	if err := d.setRFRegister(ctx, mdmcfg2, d.rcfg[cfgMdmcfg2], false); err != nil {
		return err
	}
	return d.setPower(ctx, power, invert)
}

// SetPower sets 'standard' power - for more complex power shaping this will need to be done manually
func (d *Device) SetPower(ctx context.Context, power int, invert bool) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	return d.setPower(ctx, power, invert)
}

func (d *Device) setPower(ctx context.Context, power int, invert bool) error {
	if power > 255 {
		power = 255
	}

	mod := d.MdmModulation()

	// we may be only changing PA_POWER, not power levels
	if power >= 0 {
		if mod == ModASKOOK && !invert {
			d.rcfg[cfgPaTable0] = 0x00
			d.rcfg[cfgPaTable1] = byte(power)
		} else {
			d.rcfg[cfgPaTable0] = byte(power)
			d.rcfg[cfgPaTable1] = 0x00
		}
		if err := d.setRFRegister(ctx, paTable0, d.rcfg[cfgPaTable0], false); err != nil {
			return err
		}
		if err := d.setRFRegister(ctx, paTable1, d.rcfg[cfgPaTable1], false); err != nil {
			return err
		}
	}

	d.rcfg[cfgFrend0] &^= frend0PaPower
	if mod == ModASKOOK {
		d.rcfg[cfgFrend0] |= 0x01
	}

	return d.setRFRegister(ctx, frend0, d.rcfg[cfgFrend0], false)
}

func (d *Device) MdmDRate() float64 {
	d.rcfgMu.RLock()
	rate := d.mdmDRate()
	d.rcfgMu.RUnlock()
	return rate
}

func (d *Device) mdmDRate() float64 {
	e := d.rcfg[cfgMdmcfg4] & 0xf
	m := d.rcfg[cfgMdmcfg3]
	return 1000000.0 * float64(d.mhz) * (256 + float64(m)) * float64(uint(1)<<e) / (1 << 28)
}

// SetMdmDRate sets the baud of data being modulated through the radio and returns the actual value.
func (d *Device) SetMdmDRate(ctx context.Context, drate int) (float64, error) {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()

	drateE := -1
	drateM := -1
	for e := 0; e < 16; e++ {
		m := int((float64(drate)*(1<<28)/(float64(uint(1)<<uint(e))*(float64(d.mhz)*1000000.0)) - 256) + 0.5) // rounded evenly
		if m < 256 {
			drateE = e
			drateM = m
			break
		}
	}
	if drateE < 0 {
		return 0, errors.New("rfcat: DRate does not translate into acceptable parameters")
	}

	// drate = 1000000.0 * float64(d.mhz) * (256+drateM) * (1<<drateE) / pow(2,28)
	// if self._debug: print "drate_e: %x   drate_m: %x   drate: %f Hz" % (drate_e, drate_m, drate)

	d.rcfg[cfgMdmcfg3] = byte(drateM)
	d.rcfg[cfgMdmcfg4] &^= mdmcfg4DRateE
	d.rcfg[cfgMdmcfg4] |= byte(drateE)
	if err := d.setRFRegister(ctx, mdmcfg3, d.rcfg[cfgMdmcfg3], false); err != nil {
		return 0, err
	}
	if err := d.setRFRegister(ctx, mdmcfg4, d.rcfg[cfgMdmcfg4], false); err != nil {
		return 0, err
	}
	return d.mdmDRate(), nil
}

func (d *Device) MakePktFLEN(ctx context.Context, flen int) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()

	if flen <= 0 {
		flen = rfMaxTXBlock
	}

	if flen > ep5outBufferSize-4 {
		return fmt.Errorf("rfcat: packet too large (%d bytes). Maximum fixed length packet is %d bytes", flen, ep5outBufferSize-6)
	}

	d.rcfg[cfgPktctrl0] &^= pktctrl0LengthConfig
	// if we're sending a large block, pktlen is dealt with by the firmware
	// using 'infinite' mode
	if flen > rfMaxTXBlock {
		d.rcfg[cfgPktlen] = 0x00
	} else {
		d.rcfg[cfgPktlen] = byte(flen)
	}
	if err := d.setRFRegister(ctx, pktctrl0Reg, d.rcfg[cfgPktctrl0], false); err != nil {
		return err
	}
	return d.setRFRegister(ctx, pktlenReg, d.rcfg[cfgPktlen], false)
}

func (d *Device) PktLen() int {
	d.rcfgMu.RLock()
	n := int(d.rcfg[cfgPktlen])
	d.rcfgMu.RUnlock()
	return n
}

func (d *Device) PktLenConfig() LengthConfig {
	d.rcfgMu.RLock()
	m := LengthConfig(d.rcfg[cfgPktctrl0] & pktctrl0LengthConfig)
	d.rcfgMu.RUnlock()
	return m
}

// SetMdmSyncWord sets the sync word.
func (d *Device) SetMdmSyncWord(ctx context.Context, word uint16) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgSync1] = byte(word >> 8)
	d.rcfg[cfgSync0] = byte(word & 0xff)
	if err := d.setRFRegister(ctx, sync1, d.rcfg[cfgSync1], false); err != nil {
		return err
	}
	return d.setRFRegister(ctx, sync0, d.rcfg[cfgSync0], false)
}

// MdmSyncWord returns the currently set sync word.
func (d *Device) MdmSyncWord() uint16 {
	d.rcfgMu.RLock()
	v := (uint16(d.rcfg[cfgSync1]) << 8) | uint16(d.rcfg[cfgSync0])
	d.rcfgMu.RUnlock()
	return v
}

// SetMdmSyncMode sets the sync mode.
func (d *Device) SetMdmSyncMode(ctx context.Context, mode SyncMode) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgMdmcfg2] &^= mdmcfg2SyncMode
	d.rcfg[cfgMdmcfg2] |= byte(mode)
	return d.setRFRegister(ctx, mdmcfg2, d.rcfg[cfgMdmcfg2], false)
}

// MdmSyncMode returns the currently set sync mode.
func (d *Device) MdmSyncMode() SyncMode {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgMdmcfg2] & mdmcfg2SyncMode
	d.rcfgMu.RUnlock()
	return SyncMode(v)
}

// SetPktPQT sets the Preamble Quality Threshold (PQT).
func (d *Device) SetPktPQT(ctx context.Context, num int) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	d.rcfg[cfgPktctrl1] &^= 7 << 5
	d.rcfg[cfgPktctrl1] |= byte((num & 7) << 5)
	return d.setRFRegister(ctx, pktctrl1Reg, d.rcfg[cfgPktctrl1], false)
}

// PktPQT returns the currently set Preamble Quality Threshold (PQT).
func (d *Device) PktPQT() byte {
	d.rcfgMu.RLock()
	v := d.rcfg[cfgPktctrl1]
	d.rcfgMu.RUnlock()
	return (v >> 5) & 7
}

func (d *Device) RSSI(ctx context.Context) (int, error) {
	b, err := d.peek(ctx, rssi, 1)
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, errors.New("rfcat: rssi returned empty response")
	}
	return int(b[0]), nil
}

func (d *Device) LQI(ctx context.Context) (int, error) {
	b, err := d.peek(ctx, lqi, 1)
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, errors.New("rfcat: lqi returned empty response")
	}
	return int(b[0]), nil
}

// RFxmit transmits a packet. set repeat & offset to optionally repeat tx of a section
// of the data block. repeat of 65535 or -1 means 'forever'
func (d *Device) RFxmit(ctx context.Context, data []byte, repeat, offset int) error {
	// TODO
	// if len(data) > rfMaxTXBlock {
	// 	if repeat or offset:
	// 		return PY_TX_BLOCKSIZE_INCOMPAT
	// 	return self.RFxmitLong(data, doencoding=False)
	// }

	// waitlen = len(data)
	// waitlen += repeat * (len(data) - offset)
	// wait = USB_TX_WAIT * ((waitlen / RF_MAX_TX_BLOCK) + 1)
	if repeat < 0 {
		repeat = 65535
	}
	buf := make([]byte, len(data)+6)
	buf[0] = byte(len(data) & 0xff)
	buf[1] = byte(len(data) >> 8)
	buf[2] = byte(repeat & 0xff)
	buf[3] = byte(repeat >> 8)
	buf[4] = byte(offset & 0xff)
	buf[5] = byte(offset >> 8)
	copy(buf[6:], data)
	_, err := d.send(ctx, appNIC, nicXmit, data)
	return err
}

// TODO
// def RFxmitLong(self, data, doencoding=True):
//     if len(data) > RF_MAX_TX_LONG:
//         return PY_TX_BLOCKSIZE_TOO_LARGE

//     datalen = len(data)

//     # calculate wait time
//     waitlen = len(data)
//     wait = USB_TX_WAIT * ((waitlen / RF_MAX_TX_BLOCK) + 1)

//     # load chunk buffers
//     chunks = []
//     for x in range(datalen / RF_MAX_TX_CHUNK):
//         chunks.append(data[x * RF_MAX_TX_CHUNK:(x + 1) * RF_MAX_TX_CHUNK])
//     if datalen % RF_MAX_TX_CHUNK:
//         chunks.append(data[-(datalen % RF_MAX_TX_CHUNK):])

//     preload = RF_MAX_TX_BLOCK / RF_MAX_TX_CHUNK
//     retval, ts = self.send(APP_NIC, NIC_XMIT_LONG, "%s" % struct.pack("<HB",datalen,preload)+data[:RF_MAX_TX_CHUNK * preload], wait=wait*preload)
//     #sys.stderr.write('=' + repr(retval))
//     error = struct.unpack("<B", retval[0])[0]
//     if error:
//         return error

//     chlen = len(chunks)
//     for chidx in range(preload, chlen):
//         chunk = chunks[chidx]
//         error = RC_TEMP_ERR_BUFFER_NOT_AVAILABLE
//         while error == RC_TEMP_ERR_BUFFER_NOT_AVAILABLE:
//             retval,ts = self.send(APP_NIC, NIC_XMIT_LONG_MORE, "%s" % struct.pack("B", len(chunk))+chunk, wait=wait)
//             error = struct.unpack("<B", retval[0])[0]
//         if error:
//             return error
//             #if error == RC_TEMP_ERR_BUFFER_NOT_AVAILABLE:
//             #    sys.stderr.write('.')
//         #sys.stderr.write('+')
//     # tell dongle we've finished
//     retval,ts = self.send(APP_NIC, NIC_XMIT_LONG_MORE, "%s" % struct.pack("B", 0), wait=wait)
//     return struct.unpack("<b", retval[0])[0]

// RFrecv blocks until a single packet is received or the context times out or is cancelled.
func (d *Device) RFrecv(ctx context.Context) ([]byte, error) {
	// TODO:
	// set blocksize to larger than 255 to receive large blocks or 0 to revert to normal
	// if blocksize > EP5OUT_BUFFER_SIZE:
	// 	raise(Exception("Blocksize too large. Maximum %d") % EP5OUT_BUFFER_SIZE)
	// self.send(APP_NIC, NIC_SET_RECV_LARGE, "%s" % struct.pack("<H",blocksize))
	return d.recv(ctx, appNIC, nicRecv)
}

// RecvChannel returns a channel that all received RF packets will be send to.
// The channel will be closed when the device is.
func (d *Device) RecvChannel() <-chan []byte {
	return d.rxChannel(appNIC, nicRecv)
}

type AESMode struct {
	ENCCSMode ENCCSMode
	RFInput   Crypto
	RFOutput  Crypto
}

// AESMode returns the currently set AES co-processor mode
func (d *Device) AESMode(ctx context.Context) (AESMode, error) {
	b, err := d.send(ctx, appNIC, nicGetAESMode, nil)
	if err != nil {
		return AESMode{}, err
	}
	if len(b) == 0 {
		return AESMode{}, errors.New("rfcat: empty response to AESMode")
	}
	cryptoRFInput := CryptoOff
	cryptoRFOutput := CryptoOff
	if b[0]&aesCryptoInEnable != 0 {
		if b[0]&aesCryptoInType == 0 {
			cryptoRFInput = CryptoDecrypt
		} else {
			cryptoRFInput = CryptoEncrypt
		}
	}
	if b[0]&aesCryptoOutEnable != 0 {
		if b[0]&aesCryptoOutType == 0 {
			cryptoRFOutput = CryptoDecrypt
		} else {
			cryptoRFOutput = CryptoEncrypt
		}
	}
	return AESMode{
		ENCCSMode: ENCCSMode((b[0] >> 4) & 7),
		RFInput:   cryptoRFInput,
		RFOutput:  cryptoRFOutput,
	}, nil
}

func (d *Device) partNum(ctx context.Context) (Chip, error) {
	b, err := d.send(ctx, appSystem, sysCmdPartnum, nil)
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, errors.New("rfcat: empty response to part number request")
	}
	return Chip(b[0]), nil
}

func (d *Device) BuildInfo(ctx context.Context) (string, error) {
	b, err := d.send(ctx, appSystem, sysCmdBuildtype, nil)
	return string(b), err
}

func (d *Device) CompilerInfo(ctx context.Context) (string, error) {
	b, err := d.send(ctx, appSystem, sysCmdCompiler, nil)
	return string(b), err
}

func (d *Device) Describe(ctx context.Context) (string, error) {
	build, err := d.BuildInfo(ctx)
	if err != nil {
		return "", err
	}
	compiler, err := d.CompilerInfo(ctx)
	if err != nil {
		return "", err
	}
	aesMode, err := d.AESMode(ctx)
	if err != nil {
		return "", err
	}
	// TODO:
	// # see if we have a bootloader by loooking for it's recognition semaphores
	// # in SFR I2SCLKF0 & I2SCLKF1
	// if(self.peek(0xDF46,1) == '\xF0' and self.peek(0xDF47,1) == '\x0D'):
	// 	output.append("Bootloader:          CC-Bootloader")
	// else:
	// 	output.append("Bootloader:          Not installed")
	lines := []string{
		"== Hardware ==",
		fmt.Sprintf("Dongle:              %s", strings.Split(build, " ")[0]),
		fmt.Sprintf("Firmware rev:        %s", strings.Split(build, "r")[1]),
		fmt.Sprintf("Compiler:            %s", compiler),
		fmt.Sprintf("Chip:                %s", d.chip),
		"",
		"== Frequency Configuration ==",
		fmt.Sprintf("Frequency:           %f hz", d.Freq()),
		fmt.Sprintf("Channel:             %d", d.Channel()),
		fmt.Sprintf("Intermediate freq:   %d hz", d.FsIF()),
		fmt.Sprintf("Frequency Offset:    %d +/-", d.FsOffset()),
		fmt.Sprintf("Est. Freq Offset:    %d", d.FreqOffsetEst()),
		"",
		"== Modem Configuration ==",
		fmt.Sprintf("Modulation:          %s", d.MdmModulation()),
		fmt.Sprintf("DRate:               %f hz", d.MdmDRate()),
		fmt.Sprintf("ChanBW:              %f hz", d.MdmChanBW()),
		fmt.Sprintf("DEVIATION:           %f hz", d.MdmDeviatn()),
		fmt.Sprintf("Sync Mode:           %s", d.MdmSyncMode()),
		fmt.Sprintf("Min TX Preamble:     %s", d.MdmNumPreamble()),
		fmt.Sprintf("Chan Spacing:        %f hz", d.MdmChanSpacing()),
		fmt.Sprintf("BSLimit:             %s", d.BSLimit()),
		fmt.Sprintf("DC Filter:           %t", d.MdmDCFilterEnabled()),
		fmt.Sprintf("Manchester Encoding: %t", d.MdmManchesterEnabled()),
		fmt.Sprintf("Fwd Err Correct:     %t", d.MdmFECEnabled()),
		"",
		"== Packet Configuration ==",
		fmt.Sprintf("Sync Word:           0x%04x", d.MdmSyncWord()),
		fmt.Sprintf("Packet Length:       %d", d.PktLen()),
		fmt.Sprintf("Length Config:       %s", d.PktLenConfig()),
		fmt.Sprintf("Configured Address:  0x%02x", d.configAddr()),
		fmt.Sprintf("Preamble Quality Threshold: 4 * %d", d.PktPQT()),
		fmt.Sprintf("Append Status:       %t", d.PktAppendStatusEnabled()),
		fmt.Sprintf("Rcvd Packet Check:   %s", d.AddressCheck()),
		fmt.Sprintf("Data Whitening:      %t", d.PktDataWhiteningEnabled()), // TODO ("off", "ON (but only with cc2400_en==0)")[whitedata])
		fmt.Sprintf("Packet Format:       %s", d.packetFormat()),
		fmt.Sprintf("CRC:                 %t", d.PktCRCEnabled()),
		"",
		"== AES Crypto Configuration ==",
		fmt.Sprintf("AES Mode:            %s", aesMode.ENCCSMode),
		fmt.Sprintf("Crypt RF Input:      %s", aesMode.RFInput),
		fmt.Sprintf("Crypt RF Output:     %s", aesMode.RFOutput),
		"",
		// "== Radio Test Signal Configuration ==",
		// // output.append("TEST2:               0x%x"%radiocfg.test2)
		// // output.append("TEST1:               0x%x"%radiocfg.test1)
		// // output.append("TEST0:               0x%x"%(radiocfg.test0&0xfd))
		// // output.append("VCO_SEL_CAL_EN:      0x%x"%((radiocfg.test2>>1)&1))
		// "",
		// "== Radio State ==",
		// // output.append("     MARCSTATE:      %s (%x)" % (self.getMARCSTATE(radiocfg)))
		// // output.append("     DONGLE RESPONDING:  mode :%x, last error# %d"%(self.getDebugCodes()))
		// "",
		// "== Client State ==",
		// //
	}

	return strings.Join(lines, "\n"), nil
}

func (d *Device) freqMult() float64 {
	return (0x10000 / 1000000.0) / float64(d.mhz)
}

// updateRadioConfig fetches the radio configuration from the device and updates the local cached version.
func (d *Device) updateRadioConfig(ctx context.Context) error {
	rcfg, err := d.peek(context.Background(), 0xdf00, 0x3e)
	if err != nil {
		return err
	}
	d.rcfgMu.Lock()
	d.rcfg = rcfg
	d.rcfgMu.Unlock()
	return nil
}

func (d *Device) peek(ctx context.Context, addr, bytecount uint16) ([]byte, error) {
	if bytecount <= 0 {
		bytecount = 1
	}
	b := make([]byte, 4)
	b[0] = byte(bytecount & 0xff)
	b[1] = byte(bytecount >> 8)
	b[2] = byte(addr & 0xff)
	b[3] = byte(addr >> 8)
	return d.send(ctx, appSystem, sysCmdPeek, b)
}

func (d *Device) poke(ctx context.Context, addr uint16, data []byte) ([]byte, error) {
	b := make([]byte, len(data)+2)
	b[0] = byte(addr & 0xff)
	b[1] = byte(addr >> 8)
	copy(b[2:], data)
	return d.send(ctx, appSystem, sysCmdPoke, b)
}

func (d *Device) pokeReg(ctx context.Context, addr uint16, data []byte) ([]byte, error) {
	b := make([]byte, len(data)+2)
	b[0] = byte(addr & 0xff)
	b[1] = byte(addr >> 8)
	copy(b[2:], data)
	return d.send(ctx, appSystem, sysCmdPokeReg, b)
}

// setRFRegister sets the radio register 'regaddr' to 'value' (first setting RF state
// to IDLE, then returning to RX/TX) value is always considered a 1-byte value
// if 'suppress' the radio state (RX/TX/IDLE) is not modified.
// d.rcfgMu should be locked during this call.
func (d *Device) setRFRegister(ctx context.Context, regaddr uint16, value byte, suppress bool) error {
	if suppress {
		_, err := d.poke(ctx, regaddr, []byte{value})
		return err
	}

	marcState := d.rcfg[cfgMarcstate]
	if marcState != marcStateIdle {
		if err := d.strobeModeIDLE(ctx); err != nil {
			return err
		}
	}

	if _, err := d.poke(ctx, regaddr, []byte{value}); err != nil {
		return err
	}

	return d.strobeModeReturn(ctx, marcState)
}

// setRfMode sets the radio state to "rfmode", and makes
func (d *Device) setRfMode(ctx context.Context, rfmode byte) error {
	d.rfmode = rfmode
	_, err := d.send(ctx, appSystem, sysCmdRfmode, []byte{rfmode})
	return err
}

// set standard radio state to TX/RX/IDLE (TX is pretty much only good for jamming).
// TX/RX modes are set to return to whatever state you choose here.

// SetModeTX sets radio to TX state
func (d *Device) SetModeTX(ctx context.Context) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	return d.setRfMode(ctx, rfstSTX)
}

// SetModeRX sets radio to RX state
func (d *Device) SetModeRX(ctx context.Context) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	return d.setRfMode(ctx, rfstSRX)
}

// SetModeIDLE sets radio to IDLE state
func (d *Device) SetModeIDLE(ctx context.Context) error {
	d.rcfgMu.RLock()
	defer d.rcfgMu.RUnlock()
	return d.setRfMode(ctx, rfstSIDLE)
}

// send raw state change to radio (doesn't update the return state for after RX/TX occurs)

// strobeModeTX sets radio to TX state (transient)
func (d *Device) strobeModeTX(ctx context.Context) error {
	_, err := d.poke(ctx, xRFST, []byte{rfstSTX})
	return err
}

// strobeModeRX sets radio to RX state (transient)
func (d *Device) strobeModeRX(ctx context.Context) error {
	_, err := d.poke(ctx, xRFST, []byte{rfstSRX})
	return err
}

// strobeModeIDLE sets radio to IDLE state (transient)
func (d *Device) strobeModeIDLE(ctx context.Context) error {
	_, err := d.poke(ctx, xRFST, []byte{rfstSIDLE})
	return err
}

// strobeModeFSTXON sets radio to FSTXON state (transient)
func (d *Device) strobeModeFSTXON(ctx context.Context) error {
	_, err := d.poke(ctx, xRFST, []byte{rfstSFSTXON})
	return err
}

// strobeModeCAL sets radio to CAL state (will return to whichever
// state is configured (via setMode functions)
func (d *Device) strobeModeCAL(ctx context.Context) error {
	_, err := d.poke(ctx, xRFST, []byte{rfstSCAL})
	return err
}

// strobeModeReturn attempts to return the the correct mode after configuring some radio register(s).
func (d *Device) strobeModeReturn(ctx context.Context, marcstate byte) error {
	// #if marcstate is None:
	// 	#marcstate = self.radiocfg.marcstate
	// #if self._debug: print("MARCSTATE: %x   returning to %x" % (marcstate, MARC_STATE_MAPPINGS[marcstate][2]) )
	// #self.poke(X_RFST, "%c"%MARC_STATE_MAPPINGS[marcstate][2])
	_, err := d.poke(ctx, xRFST, []byte{d.rfmode})
	return err
}

func (d *Device) send(ctx context.Context, app, cmd byte, buf []byte) ([]byte, error) {
	if len(buf) > 65535 {
		return nil, fmt.Errorf("rfcat: trying to send large message %d", len(buf))
	}
	txBuf := make([]byte, len(buf)+4)
	txBuf[0] = app
	txBuf[1] = cmd
	txBuf[2] = byte(len(buf) & 0xff)
	txBuf[3] = byte(len(buf) >> 8)
	copy(txBuf[4:], buf)
	// fmt.Printf("TX: %s\n", hex.Dump(txBuf))
	timeout := time.Second
	if t, ok := ctx.Deadline(); ok {
		dt := time.Until(t)
		if dt <= 0 {
			return nil, errors.New("rfcat: tx timeout")
		}
	}
	d.txMu.Lock()
	d.oute.Timeout = timeout
	n, err := d.oute.Write(txBuf)
	d.txMu.Unlock()
	if err != nil {
		return nil, err
	}
	if n != len(txBuf) {
		// TODO: loop to write everything?
		return nil, fmt.Errorf("rfcat: partial write %d of %d", n, len(txBuf))
	}
	return d.recv(ctx, app, cmd)
}

func (d *Device) recv(ctx context.Context, app, cmd byte) ([]byte, error) {
	select {
	case data := <-d.rxChannel(app, cmd):
		return data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *Device) sendLoop() {
	defer d.closeWG.Done()
	for {
		select {
		case <-d.closeCh:
			return
		default:
		}
	}
}

func (d *Device) rxChannel(app, cmd byte) chan []byte {
	key := (uint16(app) << 8) | uint16(cmd)
	d.rxMu.Lock()
	ch := d.rxQu[key]
	if ch == nil {
		ch = make(chan []byte, rxChannelSize)
		d.rxQu[key] = ch
	}
	d.rxMu.Unlock()
	return ch
}

func (d *Device) recvLoop() {
	defer d.closeWG.Done()
	buf := make([]byte, 65536)
	off := 0
	for {
		select {
		case <-d.closeCh:
			return
		default:
		}

		if off == len(buf) {
			log.Printf("RX overflow")
			off = 0
		}
		d.ine.Timeout = time.Millisecond * 100
		n, err := d.ine.Read(buf[off:])
		if err != nil {
			switch err := err.(type) {
			case gousb.TransferStatus:
				if err != gousb.TransferTimedOut {
					log.Printf("Error reading: %s", err)
				}
			default:
				log.Printf("Error reading: %s", err)
			}
			continue
		}
		if n == 0 {
			continue
		}
		off += n

		// Package format: '@' app:byte cmd:byte len:uint16le data:...
		for off >= 5 {
			b := buf[:off]
			if b[0] != '@' {
				ix := bytes.IndexByte(b, '@')
				if ix < 0 {
					log.Printf("Bad recv, no @")
					off = 0
					break
				}
				log.Printf("%d garbage received", ix)
				copy(buf, buf[ix:off])
				off -= ix
				continue
			}
			ln := int(b[3]) | (int(b[4]) << 8)
			if 5+ln > len(b) {
				break
			}
			app := b[1]
			cmd := b[2]
			data := b[5 : 5+ln]

			// fmt.Printf("app=%02x cmd=%02x\n%s\n", app, cmd, hex.Dump(data))

			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)
			ch := d.rxChannel(app, cmd)
			select {
			case ch <- dataCopy:
			default:
				log.Printf("Mbox for app %02x cmd %02x full", app, cmd)
				// Discard one item and try again
				select {
				case <-ch:
				default:
				}
				select {
				case ch <- dataCopy:
				default:
					// Shouldn't ever fail since we're the only producers insert into the channel
					log.Fatal("Mbox broken")
				}
			}

			if 5+ln == len(b) {
				off = 0
				break
			}
			copy(buf, buf[5+ln:off])
			off -= 5 + ln
		}
	}
}
