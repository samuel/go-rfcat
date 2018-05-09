package rfcat

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestDevice(t *testing.T) {
	dev, err := Open(0)
	if err != nil {
		t.Fatal(err)
	}
	defer dev.Close()

	ctx := context.Background()

	updateConfigAndCompare := func(t *testing.T) {
		dev.rcfgMu.Lock()
		oldCfg := make([]byte, len(dev.rcfg))
		copy(oldCfg, dev.rcfg)
		dev.rcfgMu.Unlock()
		if err := dev.updateRadioConfig(ctx); err != nil {
			t.Fatal(err)
		}
		dev.rcfgMu.Lock()
		newCfg := make([]byte, len(dev.rcfg))
		copy(newCfg, dev.rcfg)
		dev.rcfgMu.Unlock()
		if !bytes.Equal(oldCfg, newCfg) {
			var diffs []string
			for i, b := range oldCfg {
				// Skip some registers that can change without explicit requests
				// TODO: check marcstate for anything interesting
				switch i {
				case cfgFscal3, cfgFscal1, cfgMarcstate, cfgVCOVCDAC:
					continue
				case cfgFscal2:
					// TODO: setFreq changes this but the value comes back different
					continue
				}
				if b != newCfg[i] {
					diffs = append(diffs, fmt.Sprintf("[%02x] cache %02x != remote %02x", i, b, newCfg[i]))
				}
			}
			if len(diffs) != 0 {
				t.Fatalf("config fetched from radio does not match local cache\n%s", strings.Join(diffs, "\n"))
			}
		}
	}

	t.Run("ping", func(t *testing.T) {
		dt, err := dev.Ping(ctx, 5, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Ping response %s (%s per)", dt, dt/5)
	})
	t.Run("describe", func(t *testing.T) {
		desc, err := dev.Describe(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(desc)
	})
	t.Run("freq", func(t *testing.T) {
		freq, err := dev.SetFreq(ctx, 433000000)
		if err != nil {
			t.Fatal(err)
		}
		if freq2 := dev.Freq(); freq != freq2 {
			t.Fatalf("Expected %f from local cache got %f", freq, freq2)
		}
		updateConfigAndCompare(t)
		freq, err = dev.SetFreq(ctx, 920000000)
		if err != nil {
			t.Fatal(err)
		}
		if freq2 := dev.Freq(); freq != freq2 {
			t.Fatalf("Expected %f from local cache got %f", freq, freq2)
		}
		updateConfigAndCompare(t)
	})
	t.Run("deviation", func(t *testing.T) {
		hz, err := dev.SetMdmDeviatn(ctx, 26350)
		if err != nil {
			t.Fatal(err)
		}
		if hz2 := dev.MdmDeviatn(); hz != hz2 {
			t.Fatalf("Expected %f from local cache got %f", hz, hz2)
		}
		updateConfigAndCompare(t)
		hz, err = dev.SetMdmDeviatn(ctx, 15000)
		if err != nil {
			t.Fatal(err)
		}
		if hz2 := dev.MdmDeviatn(); hz != hz2 {
			t.Fatalf("Expected %f from local cache got %f", hz, hz2)
		}
		updateConfigAndCompare(t)
	})
	t.Run("modulation", func(t *testing.T) {
		if err := dev.SetMdmModulation(ctx, ModASKOOK, false); err != nil {
			t.Fatal(err)
		}
		if mod := dev.MdmModulation(); mod != ModASKOOK {
			t.Fatalf("Expected %s from local cache got %s", ModASKOOK, mod)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmModulation(ctx, Mod2FSK, false); err != nil {
			t.Fatal(err)
		}
		if mod := dev.MdmModulation(); mod != Mod2FSK {
			t.Fatalf("Expected %s from local cache got %s", Mod2FSK, mod)
		}
		updateConfigAndCompare(t)
	})
	t.Run("datarate", func(t *testing.T) {
		rate, err := dev.SetMdmDRate(ctx, 4800)
		if err != nil {
			t.Fatal(err)
		}
		if rate2 := dev.MdmDRate(); rate != rate2 {
			t.Fatalf("Expected %f from local cache got %f", rate, rate2)
		}
		updateConfigAndCompare(t)
		rate, err = dev.SetMdmDRate(ctx, 9600)
		if err != nil {
			t.Fatal(err)
		}
		if rate2 := dev.MdmDRate(); rate != rate2 {
			t.Fatalf("Expected %f from local cache got %f", rate, rate2)
		}
		updateConfigAndCompare(t)
	})
	t.Run("pktlen", func(t *testing.T) {
		if err := dev.MakePktFLEN(ctx, 32); err != nil {
			t.Fatal(err)
		}
		if ln := dev.PktLen(); ln != 32 {
			t.Fatalf("Expected %d from local cache got %d", 32, ln)
		}
		updateConfigAndCompare(t)
		if err := dev.MakePktFLEN(ctx, 64); err != nil {
			t.Fatal(err)
		}
		if ln := dev.PktLen(); ln != 64 {
			t.Fatalf("Expected %d from local cache got %d", 64, ln)
		}
		updateConfigAndCompare(t)
	})
	t.Run("syncword", func(t *testing.T) {
		if err := dev.SetMdmSyncWord(ctx, 0x0a0e); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmSyncWord(); v != 0x0a0e {
			t.Fatalf("Expected %04x from local cache got %04x", 0x0a04, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmSyncWord(ctx, 0xa5e5); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmSyncWord(); v != 0xa5e5 {
			t.Fatalf("Expected %04x from local cache got %04x", 0xa5e5, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("syncmode", func(t *testing.T) {
		if err := dev.SetMdmSyncMode(ctx, SyncModeCarrier); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmSyncMode(); v != SyncModeCarrier {
			t.Fatalf("Expected %s from local cache got %s", SyncModeCarrier, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmSyncMode(ctx, SyncMode15Of16); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmSyncMode(); v != SyncMode15Of16 {
			t.Fatalf("Expected %s from local cache got %s", SyncMode15Of16, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("pktpqt", func(t *testing.T) {
		if err := dev.SetPktPQT(ctx, 0); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktPQT(); v != 0 {
			t.Fatalf("Expected %d from local cache got %d", 0, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetPktPQT(ctx, 4); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktPQT(); v != 4 {
			t.Fatalf("Expected %d from local cache got %d", 4, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("manchester", func(t *testing.T) {
		if err := dev.SetMdmManchesterEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmManchesterEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmManchesterEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmManchesterEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("fec", func(t *testing.T) {
		if err := dev.SetMdmFECEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmFECEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmFECEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmFECEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("dcfilter", func(t *testing.T) {
		if err := dev.SetMdmDCFilterEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmDCFilterEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetMdmDCFilterEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.MdmDCFilterEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("datawhitening", func(t *testing.T) {
		if err := dev.SetPktDataWhiteningEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktDataWhiteningEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetPktDataWhiteningEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktDataWhiteningEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("appendstatus", func(t *testing.T) {
		if err := dev.SetPktAppendStatusEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktAppendStatusEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetPktAppendStatusEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktAppendStatusEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
	})
	t.Run("crc", func(t *testing.T) {
		if err := dev.SetPktCRCEnabled(ctx, true); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktCRCEnabled(); v != true {
			t.Fatalf("Expected %t from local cache got %t", true, v)
		}
		updateConfigAndCompare(t)
		if err := dev.SetPktCRCEnabled(ctx, false); err != nil {
			t.Fatal(err)
		}
		if v := dev.PktCRCEnabled(); v != false {
			t.Fatalf("Expected %t from local cache got %t", false, v)
		}
		updateConfigAndCompare(t)
	})
}
