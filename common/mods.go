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
