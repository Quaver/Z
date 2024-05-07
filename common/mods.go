package common

type Mods int64

const (
	ModNoSliderVelocities Mods = 1 << iota
	ModSpeed05X
	ModSpeed06X
	ModSpeed07X
	ModSpeed08X
	ModSpeed09X
	ModSpeed11X
	ModSpeed12X
	ModSpeed13X
	ModSpeed14X
	ModSpeed15X
	ModSpeed16X
	ModSpeed17X
	ModSpeed18X
	ModSpeed19X
	ModSpeed20X
	ModStrict
	ModChill
	ModNoPause
	ModAutoplay
	ModPaused
	ModNoFail
	ModNoLongNotes
	ModRandomize
	ModSpeed055X
	ModSpeed065X
	ModSpeed075X
	ModSpeed085X
	ModSpeed095X
	ModInverse
	ModFullLN
	ModMirror
	ModCoop
	ModSpeed105X
	ModSpeed115X
	ModSpeed125X
	ModSpeed135X
	ModSpeed145X
	ModSpeed155X
	ModSpeed165X
	ModSpeed175X
	ModSpeed185X
	ModSpeed195X
	ModHealthAdjust
	ModNoMiss
	ModEnumMaxValue // This is only in place for looping purposes (i < ModEnumMaxValue - 1; i++)
)

var SpeedMods = []Mods{
	ModSpeed05X,
	ModSpeed055X,
	ModSpeed06X,
	ModSpeed065X,
	ModSpeed07X,
	ModSpeed075X,
	ModSpeed08X,
	ModSpeed085X,
	ModSpeed09X,
	ModSpeed095X,
	0,
	ModSpeed105X,
	ModSpeed11X,
	ModSpeed115X,
	ModSpeed12X,
	ModSpeed125X,
	ModSpeed13X,
	ModSpeed135X,
	ModSpeed14X,
	ModSpeed145X,
	ModSpeed15X,
	ModSpeed155X,
	ModSpeed16X,
	ModSpeed165X,
	ModSpeed17X,
	ModSpeed175X,
	ModSpeed18X,
	ModSpeed185X,
	ModSpeed19X,
	ModSpeed195X,
	ModSpeed20X,
}

// GetSpeedModFromMods Returns the active speed mod from a group of modifiers
func GetSpeedModFromMods(mods Mods) Mods {
	for _, speedMod := range SpeedMods {
		if speedMod != 0 && mods&speedMod != 0 {
			return speedMod
		}
	}

	return 0
}

// GetModStrings Returns the active modifiers from string
func GetModStrings() map[string]Mods {
	modMap := map[string]Mods{
		"NSV":   ModNoSliderVelocities,
		"0.5x":  ModSpeed05X,
		"0.55x": ModSpeed055X,
		"0.6x":  ModSpeed06X,
		"0.65x": ModSpeed065X,
		"0.7x":  ModSpeed07X,
		"0.75x": ModSpeed075X,
		"0.8x":  ModSpeed08X,
		"0.85x": ModSpeed085X,
		"0.9x":  ModSpeed09X,
		"0.95x": ModSpeed095X,
		"1.05x": ModSpeed105X,
		"1.1x":  ModSpeed11X,
		"1.15x": ModSpeed115X,
		"1.2x":  ModSpeed12X,
		"1.25x": ModSpeed125X,
		"1.3x":  ModSpeed13X,
		"1.35x": ModSpeed135X,
		"1.4x":  ModSpeed14X,
		"1.45x": ModSpeed145X,
		"1.5x":  ModSpeed15X,
		"1.55x": ModSpeed155X,
		"1.6x":  ModSpeed16X,
		"1.65x": ModSpeed165X,
		"1.7x":  ModSpeed17X,
		"1.75x": ModSpeed175X,
		"1.8x":  ModSpeed18X,
		"1.85x": ModSpeed185X,
		"1.9x":  ModSpeed19X,
		"1.95x": ModSpeed195X,
		"2.0x":  ModSpeed20X,
		"NF":    ModNoFail,
		"MR":    ModMirror,
		"NLN":   ModNoLongNotes,
		"FLN":   ModFullLN,
		"INV":   ModInverse,
	}

	return modMap
}
