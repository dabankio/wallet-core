/*
SLIP-0044 : Registered coin types for BIP-0044
Number:  SLIP-0044
Title:   Registered coin types for BIP-0044
Type:    Standard
Status:  Draft
Authors: Pavol Rusnak <stick@satoshilabs.com>
         Marek Palatinus <slush@satoshilabs.com>
Created: 2014-07-09
*/

package bip44

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/pkg/errors"
)

const (
	PathFormat      = "m/44'/%d'"
	FullPathFormat  = "m/44'/%d'/0'/0/0"
	FullPathFormat4 = "m/44'/%d'/%d'/%d/%d" //symbol,account,0|1 external|internal(change), index
	Password        = "DASafe360"

	// ChangeTypeExternal 通常用于收款，对外部可见
	ChangeTypeExternal = 0
	// ChangeTypeInternal 通常用于找零，通常不对外可见
	ChangeTypeInternal = 1
)

func init() {
	combineCoinType()
}

func combineCoinType() {
	for coin := range customCoinType {
		if _, exist := registeredCoinType[coin]; exist {
			fmt.Printf("repetitive definition in coin type definitions: (%s)", coin)
		} else {
			if customCoinType[coin] != 0 {
				registeredCoinType[coin] = customCoinType[coin]
			}
		}
	}
}

// GetCoinType get bip44 id for symbol,
func GetCoinType(symbol string) (coinType uint32, err error) {
	if strings.Compare(strings.ToUpper(symbol), symbol) != 0 {
		// fmt.Printf("symbol has been converted to uppercase. (%s) -> (%s)", symbol, strings.ToUpper(symbol))
		symbol = strings.ToUpper(symbol)
	}
	coinType, exist := registeredCoinType[symbol]
	if !exist {
		err = errors.Errorf("unregistered coin type: %s", symbol)
	} else {
		coinType -= hdkeychain.HardenedKeyStart
	}
	return
}

// GetCoinDerivationPath get actual path for coin
func GetCoinDerivationPath(path, symbol string) (derivationPath accounts.DerivationPath, err error) {
	coinType, err := GetCoinType(symbol)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bip44 id for: %s", symbol)
	}
	return GetDerivePath(path, coinType, nil)
}

// AdditionalDeriveParam 额外的推导参数
type AdditionalDeriveParam struct {
	AccountIndex, ChangeType, Index int
}

// GetDerivePath provide path, symbole id (and other param ) to get bip44 derivation path
func GetDerivePath(path string, symbolID uint32, ap *AdditionalDeriveParam) (accounts.DerivationPath, error) {
	count := strings.Count(path, "%d")
	switch count {
	case 1:
		return accounts.ParseDerivationPath(fmt.Sprintf(path, symbolID))
	case 4:
		if ap == nil {
			return nil, errors.Errorf("bip44 additional param not provided")
		}
		if ap.AccountIndex < 0 || ap.Index < 0 {
			return nil, errors.Errorf("invalid account index or index")
		}
		if ap.ChangeType != ChangeTypeExternal && ap.ChangeType != ChangeTypeInternal {
			return nil, errors.Errorf("invalid change type %d", ap.ChangeType)
		}
		return accounts.ParseDerivationPath(fmt.Sprintf(path, symbolID, ap.AccountIndex, ap.ChangeType, ap.Index))
	default:
		return nil, errors.Errorf("bip44 derive path unknown format: %s", path)
	}
}

var customCoinType = map[string]uint32{
	"LMC": (1 << 31) + 0x8102, // LomoCoin
	"MGD": (1 << 31) + 0x8103, // MassGrid
	"WCG": (1 << 31) + 0x8104, // World Crypto Gold
}

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
var registeredCoinType = map[string]uint32{
	"BIGBANG CORE":  0x80000224,
	"MARKETFINANCE": 0x80000225,

	"AC":         0x80000033,
	"ACC":        0x800000a1,
	"ACM":        0x800000e4,
	"ACT":        0x8000029a,
	"ADA":        0x80000717,
	"ADF":        0x80000376,
	"AIB":        0x80000037,
	"AION":       0x800001a9,
	"AKA":        0x80030fb1,
	"ANON":       0x800000dc,
	"AOA":        0x80000a0a,
	"AQUA":       0x83adbc39,
	"ARA":        0x80000138,
	"ARG":        0x8000002d,
	"ARK":        0x8000006f,
	"ASK":        0x800000df,
	"ATH":        0x80000654,
	"ATN":        0x800000d0,
	"ATP":        0x800000ce,
	"AUR":        0x80000055,
	"AXE":        0x80001092,
	"BANANO":     0x800000c6,
	"BBC":        0x80000457,
	"BCA":        0x800000b9,
	"BCD":        0x800003e7,
	"BCH":        0x80000091,
	"BCO":        0x80501949,
	"BCS":        0x8000022b,
	"BCX":        0x80000698,
	"BEET":       0x80000320,
	"BELA":       0x80000049,
	"BHD":        0x8050194a,
	"BIFI":       0x800000c9,
	"BIGUP":      0x80000064,
	"BIO":        0x8000009a,
	"BIQ":        0x80000061,
	"BIS":        0x800000d1,
	"BITB":       0x80000058,
	"BITG":       0x800000de,
	"BKT":        0x800003ea,
	"BLK":        0x8000000a,
	"BLOCK":      0x80000148,
	"BOPO":       0x800000d3,
	"BOXY":       0x800000d7,
	"BPA":        0x80001a0a,
	"BRIT":       0x80000046,
	"BRK":        0x8000007e,
	"BSD":        0x8000005b,
	"BSQ":        0x8000008e,
	"BTA":        0x80000059,
	"BTC":        0x80000000,
	"BTC2X":      0x8000009d,
	"BTCC":       0x80000101,
	"BTCD":       0x80000034,
	"BTCP":       0x800000b7,
	"BTCS":       0x80008456,
	"BTCZ":       0x800000b1,
	"BTDX":       0x800000da,
	"BTF":        0x800026a0,
	"BTG":        0x8000009c,
	"BTM":        0x80000099,
	"BTN":        0x800003e8,
	"BTP":        0x80002327,
	"BTQ":        0x80002093,
	"BTR":        0x80002833,
	"BTT":        0x80008888,
	"BTV":        0x80001e61,
	"BTW":        0x80000309,
	"BTX":        0x800000a0,
	"BTY":        0x80003333,
	"BU":         0x8000020e,
	"BURST":      0x8000001e,
	"BUZZ":       0x800000a9,
	"CCC":        0x80000ccc,
	"CCN":        0x80000013,
	"CDN":        0x80000022,
	"CDY":        0x80000479,
	"CIVX":       0x800000f8,
	"CLAM":       0x80000017,
	"CLC":        0x8000076d,
	"CLO":        0x80000334,
	"CLUB":       0x8000004f,
	"CMP":        0x80000047,
	"CMT":        0x80000462,
	"CNMC":       0x800000cb,
	"COLX":       0x800007cf,
	"CPC":        0x80000121,
	"CRAVE":      0x800000ba,
	"CRW":        0x80000048,
	"CRX":        0x80000060,
	"DASH":       0x80000005,
	"DBIC":       0x80000068,
	"DCR":        0x8000002a,
	"DEO":        0x80000de0,
	"DFC":        0x80000539,
	"DGB":        0x80000014,
	"DGC":        0x80000012,
	"DIVI":       0x8000012d,
	"DLC":        0x80000066,
	"DMD":        0x80000098,
	"DNR":        0x80000074,
	"DOGE":       0x80000003,
	"DOPE":       0x80000035,
	"DST":        0x80000dec,
	"EAST":       0x80000271,
	"EC":         0x80000084,
	"ECN":        0x80000073,
	"EDRC":       0x80000038,
	"EFL":        0x8000004e,
	"EGEM":       0x800007c3,
	"ELA":        0x80000901,
	"ELLA":       0x800000a3,
	"EMC2":       0x80000029,
	"EOS":        0x800000c2,
	"EOSC":       0x800007e2,
	"ERC":        0x80000097,
	"ESN":        0x8000797e,
	"ETC":        0x8000003d,
	"ETF":        0x800000c7,
	"ETH":        0x8000003c,
	"ETHO":       0x8014095a,
	"ETP":        0x800008fe,
	"ETSC":       0x80000468,
	"EVO":        0x80000062,
	"EVT":        0x800000cf,
	"EXCL":       0x800000be,
	"EXP":        0x80000028,
	"FCT":        0x80000083,
	"FIC":        0x00001480,
	"FJC":        0x8000004b,
	"FLASH":      0x80000078,
	"FLO":        0x800000d8,
	"FO":         0x80002710,
	"FRST":       0x800000a7,
	"FTC":        0x80000008,
	"GAME":       0x80000065,
	"GB":         0x8000005e,
	"GBX":        0x800000b0,
	"GCR":        0x80000031,
	"GNX":        0x8000012c,
	"GO":         0x800017ac,
	"GOD":        0x8000270f,
	"GRC":        0x80000054,
	"GRS":        0x80000011,
	"GXC":        0x800008ff,
	"HLM":        0x800000e2,
	"HNC":        0x800000a8,
	"HNS":        0x000014e9,
	"HODL":       0x800007c5,
	"HSR":        0x800000ab,
	"HTML":       0x800000ac,
	"HUSH":       0x800000c5,
	"HYC":        0x80000575,
	"ICX":        0x8000004a,
	"ILT":        0x8011df89,
	"INSN":       0x80000044,
	"IOP":        0x80000042,
	"IOTA":       0x8000107a,
	"IXC":        0x80000056,
	"JBS":        0x8000001a,
	"KETH":       0x80010000,
	"KMD":        0x8000008d,
	"KOBO":       0x800000c4,
	"KOTO":       0x800001fe,
	"LAX":        0x801a2010,
	"LBC":        0x8000008c,
	"LBTC":       0x800003e6,
	"LCC":        0x800000c0,
	"LCH":        0x800000bd,
	"LDCN":       0x8000003f,
	"LET":        0x80000206,
	"LINX":       0x80000072,
	"LKR":        0x8000022d,
	"LSK":        0x80000086,
	"LTC":        0x80000002,
	"LTZ":        0x800000dd,
	"MARS":       0x8000006b,
	"MBRS":       0x800000aa,
	"MEC":        0x800000d9,
	"MEM":        0x8000014d,
	"MIX":        0x8000004c,
	"MKF":        0x80000225,
	"MNX":        0x800000b6,
	"MOIN":       0x80000027,
	"MONA":       0x80000016,
	"MONK":       0x800000d5,
	"MTR":        0x8000005d,
	"MUE":        0x8000001f,
	"MUSIC":      0x800000b8,
	"MXT":        0x800000b4,
	"MZC":        0x8000000d,
	"NANO":       0x80000100,
	"NAS":        0x80000a9e,
	"NAV":        0x80000082,
	"NBT":        0x8000000c,
	"NDAU":       0x80004e44,
	"NEBL":       0x80000092,
	"NEET":       0x800000d2,
	"NEO":        0x80000378,
	"NEOS":       0x80000019,
	"NIM":        0x800000f2,
	"NLC2":       0x80000095,
	"NLG":        0x80000057,
	"NMC":        0x80000007,
	"NOS":        0x800000e5,
	"NPW":        0x800000fd,
	"NRG":        0x800000cc,
	"NRO":        0x8000006e,
	"NSR":        0x8000000b,
	"NUKO":       0x8000012b,
	"NVC":        0x80000032,
	"NXS":        0x80000043,
	"NXT":        0x8000001d,
	"NYC":        0x800000b3,
	"ODN":        0x800000ad,
	"OK":         0x80000045,
	"OMNI":       0x800000c8,
	"USDT(OMNI)": 0x800000c8, //same as omni
	"ONE":        0x8000010e,
	"ONT":        0x80000400,
	"ONX":        0x800000ae,
	"OOT":        0x800000d4,
	"PART":       0x8000002c,
	"PHL":        0x800007c6,
	"PHR":        0x800001bc,
	"PIGGY":      0x80000076,
	"PINK":       0x80000075,
	"PIRL":       0x800000a4,
	"PIVX":       0x80000077,
	"PKB":        0x80000024,
	"PND":        0x80000025,
	"POA":        0x800000b2,
	"POLIS":      0x800007cd,
	"POT":        0x80000051,
	"PPC":        0x80000006,
	"PRJ":        0x80000215,
	"PSB":        0x8000003e,
	"PTC":        0x8000006d,
	"PUT":        0x8000007a,
	"PWR":        0x800057e8,
	"QKC":        0x85f5e0ff,
	"QRK":        0x80000052,
	"QTUM":       0x800008fd,
	"QVT":        0x80000328,
	"RAP":        0x80000141,
	"RBY":        0x80000010,
	"RDD":        0x80000004,
	"RIC":        0x8000008f,
	"RICHX":      0x80000050,
	"RIN":        0x800000cd,
	"RISE":       0x80000460,
	"ROGER":      0x80001b39,
	"ROI":        0x80000d31,
	"RPT":        0x8000008b,
	"RVN":        0x800000af,
	"SAFE":       0x80001a20,
	"SBC":        0x8000010f,
	"SBTC":       0x800022b8,
	"SDC":        0x80000023,
	"SDGO":       0x80003de5,
	"SH":         0x8000006a,
	"SHM":        0x8000005f,
	"SHR":        0x80000030,
	"SLP":        0x800000f5,
	"SLR":        0x8000003a,
	"SMART":      0x800000e0,
	"SMLY":       0x8000003b,
	"SOOM":       0x800000fa,
	"SSC":        0x80000900,
	"SSN":        0x8000009e,
	"SSP":        0x8000026a,
	"STAK":       0x800000bb,
	"START":      0x80000026,
	"STEEM":      0x80000087,
	"STO":        0x80000063,
	"STRAT":      0x80000069,
	"SYS":        0x80000039,
	"TES":        0x80000743,
	"THC":        0x80000071,
	"TOA":        0x8000009f,
	"TOMO":       0x80000379,
	"TPC":        0x80000036,
	"TRC":        0x80000053,
	"TRTL":       0x800007c0,
	"TRX":        0x800000c3,
	"TT":         0x800003e9,
	"TZC":        0x800000e8,
	"UBQ":        0x8000006c,
	"UC":         0x800000f7,
	"UFO":        0x800000ca,
	"UNIFY":      0x8000007c,
	"UNO":        0x8000005c,
	"USC":        0x80000070,
	"VAR":        0x800000e9,
	"VASH":       0x80000021,
	"VC":         0x8000007f,
	"VET":        0x80000332,
	"VIA":        0x8000000e,
	"VIPS":       0x8000077f,
	"VITE":       0x800a2c2a,
	"VIVO":       0x800000a6,
	"VOX":        0x80000081,
	"VTC":        0x8000001c,
	"WAN":        0x8057414e,
	"WAVES":      0x80579bfc,
	"WBTC":       0x800000bc,
	"WC":         0x800000b5,
	"WEB":        0x800000e3,
	"WHL":        0x80000096,
	"WICC":       0x8001869f,
	"X42":        0x80067932,
	"XAX":        0x800000db,
	"XBC":        0x80000041,
	"XCH":        0x8000000f,
	"XCP":        0x80000009,
	"XEM":        0x8000002b,
	"XFE":        0x800000c1,
	"XLM":        0x80000094,
	"XMCC":       0x800007ce,
	"XMR":        0x80000080,
	"XMX":        0x800007b9,
	"XMY":        0x8000005a,
	"XPM":        0x80000018,
	"XRB":        0x800000a5,
	"XRD":        0x80000200,
	"XRP":        0x80000090,
	"XSEL":       0x8000037a,
	"XSN":        0x80000180,
	"XSPEC":      0x800000d5,
	"XST":        0x8000007d,
	"XTZ":        0x800006c1,
	"XUEZ":       0x800000e1,
	"XVG":        0x8000004d,
	"XWC":        0x8000009b,
	"XZC":        0x80000088,
	"YAP":        0x80000210,
	"YCC":        0x80003334,
	"ZCL":        0x80000093,
	"ZEC":        0x80000085,
	"ZEN":        0x80000079,
	"ZEST":       0x80000103,
	"ZNY":        0x8000007b,
	"ZOOM":       0x80000020,
	"ZRC":        0x8000001b,
	"ZYD":        0x80000067,
	"kUSD":       0x857ab1e1,
}
