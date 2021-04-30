package tbc

import (
	"fmt"
)

var Gems = []Gem{
	// {Name: "Destructive Skyfire Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Enigmatic Skyfire Diamond", Color: GemColorMeta, Stats: Stats{}},
	{Name: "Chaotic Skyfire Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatSpellCrit: 12}, Activate: ActivateCSD},
	// {Name: "Swift Skyfire Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Potent Unstable Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Swift Windfire Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Powerful Earthstorm Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Thundering Skyfire Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Relentless Earthstorm Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Tenacious Earthstorm Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Eternal Earthstorm Diamond", Color: GemColorMeta, Stats: Stats{}},
	// {Name: "Brutal Earthstorm Diamond", Color: GemColorMeta, Stats: Stats{}},

	{Name: "Bracing Earthstorm Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatSpellDmg: 14}},
	{Name: "Imbued Unstable Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatSpellDmg: 14}},
	{Name: "Ember Skyfire Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatSpellDmg: 14}, Activate: ActivateESD},
	{Name: "Swift Starfire Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatSpellDmg: 12}},
	{Name: "Mystical Skyfire Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{}, Activate: ActivateMSD},
	{Name: "Insightful Earthstorm Diamond", Quality: ItemQualityRare, Phase: 1, Color: GemColorMeta, Stats: Stats{StatInt: 12}, Activate: ActivateIED},
	{Name: "Runed Blood Garnet", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorRed, Stats: Stats{StatSpellDmg: 7}},
	{Name: "Runed Living Ruby", Quality: ItemQualityRare, Phase: 1, Color: GemColorRed, Stats: Stats{StatSpellDmg: 9}},
	{Name: "Runed Crimson Spinel", Quality: ItemQualityEpic, Phase: 3, Color: GemColorRed, Stats: Stats{StatSpellDmg: 12}},
	{Name: "Runed Ornate Ruby", Quality: ItemQualityEpic, Phase: 1, Color: GemColorRed, Stats: Stats{StatSpellDmg: 12}},
	{Name: "Don Julio's Heart", Quality: ItemQualityEpic, Phase: 1, Color: GemColorRed, Stats: Stats{StatSpellDmg: 14}},
	{Name: "Lustrous Azure Moonstone", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorBlue, Stats: Stats{StatMP5: 2}},
	{Name: "Lustrous Star of Elune", Quality: ItemQualityRare, Phase: 1, Color: GemColorBlue, Stats: Stats{StatMP5: 3}},
	{Name: "Lustrous Empyrean Sapphire", Quality: ItemQualityEpic, Phase: 1, Color: GemColorBlue, Stats: Stats{StatMP5: 4}},
	{Name: "Brilliant Golden Draenite", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorYellow, Stats: Stats{StatInt: 6}},
	{Name: "Brilliant Dawnstone", Quality: ItemQualityRare, Phase: 1, Color: GemColorYellow, Stats: Stats{StatInt: 8}},
	{Name: "Brilliant Lionseye", Quality: ItemQualityEpic, Phase: 3, Color: GemColorYellow, Stats: Stats{StatInt: 10}},
	{Name: "Gleaming Golden Draenite", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorYellow, Stats: Stats{StatSpellCrit: 6}},
	{Name: "Gleaming Dawnstone", Quality: ItemQualityRare, Phase: 1, Color: GemColorYellow, Stats: Stats{StatSpellCrit: 8}},
	{Name: "Gleaming Lionseye", Quality: ItemQualityEpic, Phase: 3, Color: GemColorYellow, Stats: Stats{StatSpellCrit: 10}},
	{Name: "Infused Fire Opal", Quality: ItemQualityEpic, Phase: 1, Color: GemColorOrange, Stats: Stats{StatInt: 4, StatSpellDmg: 6}},
	{Name: "Potent Flame Spessarite", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellCrit: 3, StatSpellDmg: 4}},
	{Name: "Potent Noble Topaz", Quality: ItemQualityRare, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellCrit: 4, StatSpellDmg: 5}},
	{Name: "Potent Pyrestone", Quality: ItemQualityEpic, Phase: 3, Color: GemColorOrange, Stats: Stats{StatSpellCrit: 5, StatSpellDmg: 6}},
	{Name: "Potent Fire Opal", Quality: ItemQualityEpic, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellDmg: 6, StatSpellCrit: 4}},
	{Name: "Potent Ornate Topaz", Quality: ItemQualityEpic, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellDmg: 6, StatSpellCrit: 5}},
	{Name: "Veiled Flame Spessarite", Quality: ItemQualityUncommon, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellHit: 3, StatSpellDmg: 4}},
	{Name: "Veiled Noble Topaz", Quality: ItemQualityRare, Phase: 1, Color: GemColorOrange, Stats: Stats{StatSpellHit: 4, StatSpellDmg: 5}},
	{Name: "Veiled Pyrestone", Quality: ItemQualityEpic, Phase: 3, Color: GemColorOrange, Stats: Stats{StatSpellHit: 5, StatSpellDmg: 6}},
	{Name: "Rune Covered Chrysoprase", Quality: ItemQualityEpic, Phase: 1, Color: GemColorGreen, Stats: Stats{StatMP5: 2, StatSpellCrit: 5}},
	{Name: "Dazzling Talasite", Quality: ItemQualityRare, Phase: 1, Color: GemColorGreen, Stats: Stats{StatMP5: 2, StatInt: 4}},
	{Name: "Glowing Nightseye", Quality: ItemQualityRare, Phase: 1, Color: GemColorPurple, Stats: Stats{StatSpellDmg: 5, StatStm: 6}},
	{Name: "Glowing Tanzanite", Quality: ItemQualityEpic, Phase: 1, Color: GemColorPurple, Stats: Stats{StatSpellDmg: 6, StatStm: 6}},
	{Name: "Infused Amethyst", Quality: ItemQualityEpic, Phase: 1, Color: GemColorPurple, Stats: Stats{StatSpellDmg: 6, StatStm: 6}},
	{Name: "Fluorescent Tanzanite", Quality: ItemQualityEpic, Phase: 1, Color: GemColorPurple, Stats: Stats{StatSpellDmg: 6, StatSpirit: 4}},
}

var Enchants = []Enchant{
	{Name: "Glyph of Power", Bonus: Stats{StatSpellDmg: 22, StatSpellHit: 14}, Slot: EquipHead},
	{Name: "Greater Inscription of the Orb", Bonus: Stats{StatSpellDmg: 12, StatSpellCrit: 15}, Slot: EquipShoulder},
	{Name: "Greater Inscription of Discipline", Bonus: Stats{StatSpellDmg: 18, StatSpellCrit: 10}, Slot: EquipShoulder},
	{Name: "Power of the Scourge", Bonus: Stats{StatSpellDmg: 15, StatSpellCrit: 14}, Slot: EquipShoulder},
	{Name: "Chest - Exceptional Stats", Bonus: Stats{StatStm: 6, StatInt: 6, StatSpirit: 6}, Slot: EquipChest},
	{Name: "Bracer - Spellpower", Bonus: Stats{StatSpellDmg: 15}, Slot: EquipWrist},
	{Name: "Gloves - Major Spellpower", Bonus: Stats{StatSpellDmg: 20}, Slot: EquipHands},
	{Name: "Runic Spellthread", Bonus: Stats{StatSpellDmg: 20}, Slot: EquipLegs},
	{Name: "Weapon - Major Spellpower", Bonus: Stats{StatSpellDmg: 40}, Slot: EquipWeapon},
	{Name: "Ring - Spellpower", Bonus: Stats{StatSpellDmg: 12}, Slot: EquipFinger},
}

var ItemLookup = map[string]Item{}
var GemLookup = map[string]Gem{}
var EnchantLookup = map[string]Enchant{}

func init() {
	for _, v := range Enchants {
		EnchantLookup[v.Name] = v
	}
	for _, v := range Gems {
		GemLookup[v.Name] = v
	}
	for _, v := range items {
		if it, ok := ItemLookup[v.Name]; ok {
			fmt.Printf("Found dup item: %s\n", v.Name)
			statsMatch := it.Slot == v.Slot
			for i, v := range v.Stats {
				if len(it.Stats) <= i {
					break
				}
				if it.Stats[i] != v {
					statsMatch = false
				}
			}
			if !statsMatch {
				// log.Printf("Mismatched slot/stats: \n\tMoreItem: \n%#v\n\t FirstSet: \n%#v", it, v)
			}
		} else {
			cv := v
			ItemLookup[cv.Name] = cv
		}
	}
}

type Item struct {
	Slot       byte
	SubSlot    byte `json:"subSlot,omitempty"`
	Name       string
	SourceZone string
	SourceDrop string
	Stats      Stats // Stats applied to wearer
	Phase      byte
	Quality    ItemQuality

	GemSlots    []GemColor
	SocketBonus Stats

	// Modified for each instance of the item.
	Gems    []Gem
	Enchant Enchant

	// For simplicity all items that produce an aura are 'activatable'.
	// Since we activate all items on CD, this works fine for stuff like Quags Eye.
	// TODO: is this the best design for this?
	Activate   ItemActivation `json:"-"` // Activatable Ability, produces an aura
	ActivateCD int            `json:"-"` // cooldown on activation, -1 means perm effect.
	CoolID     int32          `json:"-"` // ID used for cooldown
}

type ItemQuality byte

const (
	ItemQualityJunk      ItemQuality = iota // anything less than green
	ItemQualityUncommon                     // green
	ItemQualityRare                         // blue
	ItemQualityEpic                         // purple
	ItemQualityLegendary                    // orange
)

type Enchant struct {
	Name  string
	Bonus Stats
	Slot  byte // which slot does the enchant go on.
}

type Gem struct {
	Name     string
	Stats    Stats          // flat stats gem adds
	Activate ItemActivation `json:"-"` // Meta gems activate an aura on player when socketed. Assumes all gems are 'always active'
	Color    GemColor
	Phase    byte
	Quality  ItemQuality
	// Requirements  // Validate the gem can be used... later
}

type GemColor byte

const (
	GemColorUnknown GemColor = iota
	GemColorMeta
	GemColorRed
	GemColorBlue
	GemColorYellow
	GemColorGreen
	GemColorOrange
	GemColorPurple
	GemColorPrismatic
)

func (gm GemColor) Intersects(o GemColor) bool {
	if gm == o {
		return true
	}
	if gm == GemColorPrismatic || o == GemColorPrismatic {
		return true
	}
	if gm == GemColorMeta {
		return false // meta gems o nothing.
	}
	if gm == GemColorRed {
		return o == GemColorOrange || o == GemColorPurple
	}
	if gm == GemColorBlue {
		return o == GemColorGreen || o == GemColorPurple
	}
	if gm == GemColorYellow {
		return o == GemColorGreen || o == GemColorOrange
	}
	if gm == GemColorOrange {
		return o == GemColorYellow || o == GemColorRed
	}
	if gm == GemColorGreen {
		return o == GemColorYellow || o == GemColorBlue
	}
	if gm == GemColorPurple {
		return o == GemColorBlue || o == GemColorRed
	}

	return false // dunno wtf this is.
}

type ItemActivation func(*Simulation) Aura

type Equipment []Item

func NewEquipmentSet(names ...string) Equipment {
	e := Equipment{EquipTotem: Item{}}
	for _, v := range names {
		item, ok := ItemLookup[v]
		if !ok {
			fmt.Printf("Unable to find item: '%s'\n", v)
			continue
		}
		item.Gems = make([]Gem, len(item.GemSlots))
		if item.Slot == EquipFinger {
			if e[EquipFinger1].Name == "" {
				e[EquipFinger1] = item
			} else {
				e[EquipFinger2] = item
			}
		} else if item.Slot == EquipTrinket {
			if e[EquipTrinket1].Name == "" {
				e[EquipTrinket1] = item
			} else {
				e[EquipTrinket2] = item
			}
		} else {
			e[item.Slot] = item
		}
	}
	return e
}

const (
	EquipUnknown byte = iota
	EquipHead
	EquipNeck
	EquipShoulder
	EquipBack
	EquipChest
	EquipWrist
	EquipHands
	EquipWaist
	EquipLegs
	EquipFeet
	EquipFinger  // generic finger item
	EquipFinger1 // specific slot in equipment array
	EquipFinger2
	EquipTrinket  // generic trinket
	EquipTrinket1 // specific trinket slot in equip array
	EquipTrinket2
	EquipWeapon // holds 1 or 2h
	EquipOffhand
	EquipTotem
)

func (e Equipment) Clone() Equipment {
	ne := make(Equipment, len(e))
	for i, v := range e {
		vc := v
		ne[i] = vc
	}
	return ne
}

func (e Equipment) Stats() Stats {
	s := Stats{StatLen: 0}
	for _, item := range e {
		for k, v := range item.Stats {
			s[k] += v
		}
		isMatched := len(item.Gems) == len(item.GemSlots) && len(item.GemSlots) > 0
		for gi, g := range item.Gems {
			for k, v := range g.Stats {
				s[k] += v
			}
			isMatched = isMatched && g.Color.Intersects(item.GemSlots[gi])
			if !isMatched {
			}
		}
		if len(item.GemSlots) > 0 {
		}
		if isMatched {
			for k, v := range item.SocketBonus {
				if v == 0 {
					continue
				}
				s[k] += v
			}
		}
		for k, v := range item.Enchant.Bonus {
			s[k] += v
		}
	}
	return s
}

// Hopefully get access to:
// https://docs.google.com/spreadsheets/d/1XkLW3o9VrYg8VT84tCoINq-KxP9EA876RdhIzO-PcQk/edit#gid=1056257705

var items = []Item{
	// source: https://docs.google.com/spreadsheets/d/1X-XO9N1_MPIq-UIpTN13LrhXRoho9fe26YEEM48QmPk/edit#gid=2035379487
	{Slot: EquipHead, Name: "Gladiator's Mail Helm", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{15, 54, 18, 0, 37, 0, 0}, GemSlots: []GemColor{0x1, 0x2}, SocketBonus: Stats{}},
	{Slot: EquipHead, Name: "Spellstrike Hood", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{12, 16, 24, 16, 46, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{StatStm: 6}},
	{Slot: EquipHead, Name: "Incanter's Cowl", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{27, 15, 19, 0, 29, 0, 0}, GemSlots: []GemColor{0x1, 0x4}, SocketBonus: Stats{}},
	{Slot: EquipHead, Name: "Lightning Crown", Phase: 1, Quality: ItemQualityEpic, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{0, 0, 43, 0, 66, 0, 0}},
	{Slot: EquipHead, Name: "Hood of Oblivion", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{32, 27, 0, 0, 40, 0, 0}, GemSlots: []GemColor{0x1, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Exorcist's Mail Helm", Phase: 1, Quality: ItemQualityRare, SourceZone: "18 Spirit Shards", SourceDrop: "", Stats: Stats{16, 30, 24, 0, 29, 0, 0}, GemSlots: []GemColor{0x1}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: EquipHead, Name: "Tidefury Helm", Phase: 1, Quality: ItemQualityRare, SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{26, 32, 0, 0, 32, 0, 6}, GemSlots: []GemColor{0x1, 0x4}, SocketBonus: Stats{StatInt: 4}},
	{Slot: EquipHead, Name: "Windscale Hood", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Leatherworking BoE", SourceDrop: "", Stats: Stats{18, 16, 37, 0, 44, 0, 10}},
	{Slot: EquipHead, Name: "Shamanistic Helmet of Second Sight", Phase: 1, Quality: ItemQualityRare, SourceZone: "Teron Gorfiend, I am... - SMV Quest", SourceDrop: "", Stats: Stats{15, 12, 24, 0, 35, 0, 4}, GemSlots: []GemColor{0x4, 0x3, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Mana-Etched Crown", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Aeonus", SourceDrop: "", Stats: Stats{20, 27, 0, 0, 34, 0, 0}, GemSlots: []GemColor{0x1, 0x2}, SocketBonus: Stats{}},
	{Slot: EquipHead, Name: "Mag'hari Ritualist's Horns", Phase: 1, Quality: ItemQualityRare, SourceZone: "Hero of the Mag'har - Nagrand quest (Horde)", SourceDrop: "", Stats: Stats{16, 18, 15, 12, 50, 0, 0}},
	{Slot: EquipHead, Name: "Mage-Collar of the Firestorm", Phase: 1, Quality: ItemQualityRare, SourceZone: "H BF - The Maker", SourceDrop: "", Stats: Stats{33, 32, 23, 0, 39, 0, 0}},
	{Slot: EquipHead, Name: "Circlet of the Starcaller", Phase: 1, Quality: ItemQualityRare, SourceZone: "Dimensius the All-Devouring - NS Quest", SourceDrop: "", Stats: Stats{18, 27, 18, 0, 47, 0, 0}},
	{Slot: EquipHead, Name: "Mask of Inner Fire", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Chrono Lord Deja", SourceDrop: "", Stats: Stats{33, 30, 22, 0, 37, 0, 0}},
	{Slot: EquipHead, Name: "Mooncrest Headdress", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "Blast the Infernals! - SMV Quest", SourceDrop: "", Stats: Stats{16, 0, 21, 0, 44, 0, 0}},
	{Slot: EquipNeck, Name: "Pendant of Dominance", Phase: 1, Quality: ItemQualityEpic, SourceZone: "15,300 Honor & 10 EotS Marks", SourceDrop: "", Stats: Stats{12, 31, 16, 0, 26, 0, 0}, GemSlots: []GemColor{0x4}, SocketBonus: Stats{StatSpellCrit: 2}},
	{Slot: EquipNeck, Name: "Brooch of Heightened Potential", Phase: 1, Quality: ItemQualityRare, SourceZone: "SLabs - Blackheart the Inciter", SourceDrop: "", Stats: Stats{14, 15, 14, 9, 22, 0, 0}},
	{Slot: EquipNeck, Name: "Torc of the Sethekk Prophet", Phase: 1, Quality: ItemQualityRare, SourceZone: "Brother Against Brother - Auchindoun ", SourceDrop: "", Stats: Stats{18, 0, 21, 0, 19, 0, 0}},
	{Slot: EquipNeck, Name: "Natasha's Ember Necklace", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Hound-Master - BEM Quest", SourceDrop: "", Stats: Stats{15, 0, 10, 0, 29, 0, 0}},
	{Slot: EquipNeck, Name: "Warp Engineer's Prismatic Chain", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Mechano Lord Capacitus", SourceDrop: "", Stats: Stats{18, 17, 16, 0, 19, 0, 0}},
	{Slot: EquipNeck, Name: "Hydra-fang Necklace", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H UB - Ghaz'an", SourceDrop: "", Stats: Stats{16, 17, 0, 16, 19, 0, 0}},
	{Slot: EquipNeck, Name: "Natasha's Arcane Filament", Phase: 1, Quality: ItemQualityEpic, SourceZone: "The Hound-Master - BEM Quest", SourceDrop: "", Stats: Stats{10, 22, 0, 0, 29, 0, 0}},
	{Slot: EquipNeck, Name: "Omor's Unyielding Will", Phase: 1, Quality: ItemQualityRare, SourceZone: "H Ramps - Omar the Unscarred", SourceDrop: "", Stats: Stats{19, 19, 0, 0, 25, 0, 0}},
	{Slot: EquipNeck, Name: "Charlotte's Ivy", Phase: 1, Quality: ItemQualityEpic, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{19, 18, 0, 0, 23, 0, 0}},
	{Slot: EquipShoulder, Name: "Gladiator's Mail Spaulders", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{17, 33, 20, 0, 22, 0, 6}, GemSlots: []GemColor{0x2, 0x4}, SocketBonus: Stats{}},
	{Slot: EquipShoulder, Name: "Pauldrons of Wild Magic", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{28, 21, 23, 0, 33, 0, 0}},
	{Slot: EquipShoulder, Name: "Mana-Etched Spaulders", Phase: 1, Quality: ItemQualityRare, SourceZone: "H UB - Quagmirran", SourceDrop: "", Stats: Stats{17, 25, 16, 0, 20, 0, 0}, GemSlots: []GemColor{0x2, 0x4}, SocketBonus: Stats{}},
	{Slot: EquipShoulder, Name: "Spaulders of the Torn-heart", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Cipher of Damnation - SMV Quest", SourceDrop: "", Stats: Stats{7, 10, 18, 0, 40, 0, 0}},
	{Slot: EquipShoulder, Name: "Elekk Hide Spaulders", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "The Fallen Exarch - Terokkar Forest Quest", SourceDrop: "", Stats: Stats{12, 0, 28, 0, 25, 0, 0}},
	{Slot: EquipShoulder, Name: "Spaulders of Oblivion", Phase: 1, Quality: ItemQualityRare, SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{17, 25, 0, 0, 29, 0, 0}, GemSlots: []GemColor{0x4, 0x3}, SocketBonus: Stats{StatSpellHit: 3}},
	{Slot: EquipShoulder, Name: "Tidefury Shoulderguards", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - O'mrogg", SourceDrop: "", Stats: Stats{23, 18, 0, 0, 19, 0, 6}, GemSlots: []GemColor{0x2, 0x3}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipShoulder, Name: "Mantle of Three Terrors", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Chrono Lord Deja", SourceDrop: "", Stats: Stats{25, 29, 0, 12, 29, 0, 0}},
	{Slot: EquipBack, Name: "Ogre Slayer's Cover", Phase: 1, Quality: ItemQualityRare, SourceZone: "Cho'war the Pillager - Nagrand Quest", SourceDrop: "", Stats: Stats{18, 0, 16, 0, 20, 0, 0}},
	{Slot: EquipBack, Name: "Baba's Cloak of Arcanistry", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{15, 15, 14, 0, 22, 0, 0}},
	{Slot: EquipBack, Name: "Cloak of Woven Energy", Phase: 1, Quality: ItemQualityRare, SourceZone: "Hitting the Motherlode - Netherstorm Quest", SourceDrop: "", Stats: Stats{13, 6, 6, 0, 29, 0, 0}},
	{Slot: EquipBack, Name: "Sethekk Oracle Cloak", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{18, 18, 0, 12, 22, 0, 0}},
	{Slot: EquipBack, Name: "Terokk's Wisdom", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Terokk - Skettis Summoned Boss", SourceDrop: "", Stats: Stats{16, 18, 0, 0, 33, 0, 0}},
	{Slot: EquipBack, Name: "Cloak of the Black Void", Phase: 1, Quality: ItemQualityRare, SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{11, 0, 0, 0, 35, 0, 0}},
	{Slot: EquipBack, Name: "Cloak of Entropy", Phase: 1, Quality: ItemQualityRare, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{11, 0, 0, 10, 25, 0, 0}},
	{Slot: EquipBack, Name: "Sergeant's Heavy Cape", Phase: 1, Quality: ItemQualityEpic, SourceZone: "9,435 Honor & 20 AB Marks", SourceDrop: "", Stats: Stats{12, 33, 0, 0, 26, 0, 0}},
	{Slot: EquipChest, Name: "Gladiator's Mail Armor", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{23, 42, 23, 0, 32, 0, 7}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatSpellCrit: 4}},
	{Slot: EquipChest, Name: "Will of Edward the Odd", Phase: 1, Quality: ItemQualityEpic, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{30, 0, 30, 0, 53, 0, 0}},
	{Slot: EquipChest, Name: "Anchorite's Robe", Phase: 1, Quality: ItemQualityEpic, SourceZone: "The Aldor - Honored", SourceDrop: "", Stats: Stats{38, 16, 0, 0, 29, 0, 18}, GemSlots: []GemColor{0x4, 0x4, 0x3}},
	{Slot: EquipChest, Name: "Tidefury Chestpiece", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{22, 28, 0, 10, 36, 0, 4}, GemSlots: []GemColor{0x4, 0x4, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipChest, Name: "Auchenai Anchorite's Robe", Phase: 1, Quality: ItemQualityRare, SourceZone: "Everything Will Be Alright - AC Quest", SourceDrop: "", Stats: Stats{24, 0, 0, 23, 28, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatSpellCrit: 4}},
	{Slot: EquipChest, Name: "Mana-Etched Vestments", Phase: 1, Quality: ItemQualityRare, SourceZone: "OHF - Epoch Hunter", SourceDrop: "", Stats: Stats{25, 25, 17, 0, 29, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipChest, Name: "Robe of the Crimson Order", Phase: 1, Quality: ItemQualityEpic, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{23, 0, 0, 30, 50, 0, 0}},
	{Slot: EquipChest, Name: "Warp Infused Drape", Phase: 1, Quality: ItemQualityRare, SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{28, 27, 0, 12, 30, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{}},
	{Slot: EquipChest, Name: "Robe of Oblivion", Phase: 1, Quality: ItemQualityRare, SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{20, 30, 0, 0, 40, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{}},
	{Slot: EquipChest, Name: "Incanter's Robe", Phase: 1, Quality: ItemQualityRare, SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{22, 24, 8, 0, 29, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{}},
	{Slot: EquipChest, Name: "Robe of the Great Dark Beyond", Phase: 1, Quality: ItemQualityRare, SourceZone: "MT - Tavarok", SourceDrop: "", Stats: Stats{30, 25, 23, 0, 39, 0, 0}},
	{Slot: EquipChest, Name: "Worldfire Chestguard", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Dalliah the Doomsayer", SourceDrop: "", Stats: Stats{32, 33, 22, 0, 40, 0, 0}},
	{Slot: 0x6, Name: "General's Mail Bracers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "7,548 Honor & 20 WSG Marks", SourceDrop: "", Stats: Stats{12, 22, 14, 0, 20, 0, 0}, GemSlots: []GemColor{0x4}, SocketBonus: Stats{}},
	{Slot: 0x6, Name: "World's End Bracers", Phase: 1, Quality: ItemQualityRare, SourceZone: "H BF - Keli'dan the Breaker", SourceDrop: "", Stats: Stats{19, 18, 17, 0, 22, 0, 0}},
	{Slot: 0x6, Name: "Bracers of Havok", Phase: 1, Quality: ItemQualityRare, SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{12, 0, 0, 0, 30, 0, 0}, GemSlots: []GemColor{0x4}, SocketBonus: Stats{StatSpellCrit: 2}},
	{Slot: 0x6, Name: "Crimson Bracers of Gloom", Phase: 1, Quality: ItemQualityRare, SourceZone: "H Ramps - Omor the Unscarred", SourceDrop: "", Stats: Stats{18, 18, 0, 12, 22, 0, 0}},
	{Slot: 0x6, Name: "Bands of Negation", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H MT - Nexus- Prince Shaffar", SourceDrop: "", Stats: Stats{22, 25, 0, 0, 29, 0, 0}},
	{Slot: 0x6, Name: "Arcanium Signet Bands", Phase: 1, Quality: ItemQualityRare, SourceZone: "H UB - Hungarfen", SourceDrop: "", Stats: Stats{15, 14, 0, 0, 30, 0, 0}},
	{Slot: 0x6, Name: "Wave-Fury Vambraces", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H SV - Warlod Kalithresh", SourceDrop: "", Stats: Stats{18, 19, 0, 0, 22, 0, 5}},
	{Slot: 0x6, Name: "Mana Infused Wristguards", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "A Fate Worse Than Death - Netherstorm Quest", SourceDrop: "", Stats: Stats{8, 12, 0, 0, 25, 0, 0}},
	{Slot: 0x7, Name: "Mana-Etched Gloves", Phase: 1, Quality: ItemQualityRare, SourceZone: "H Ramps - Omor the Unscarred", SourceDrop: "", Stats: Stats{17, 25, 16, 0, 20, 0, 0}, GemSlots: []GemColor{0x2, 0x4}, SocketBonus: Stats{}},
	{Slot: 0x7, Name: "Earth Mantle Handwraps", Phase: 1, Quality: ItemQualityRare, SourceZone: "SV - Mekgineer Steamrigger", SourceDrop: "", Stats: Stats{18, 21, 16, 0, 19, 0, 0}, GemSlots: []GemColor{0x2, 0x4}, SocketBonus: Stats{StatInt: 3}},
	{Slot: 0x7, Name: "Gloves of Pandemonium", Phase: 1, Quality: ItemQualityRare, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{15, 0, 22, 10, 25, 0, 0}},
	{Slot: 0x7, Name: "Gladiator's Mail Gauntlets", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{18, 36, 21, 0, 32, 0, 0}},
	{Slot: 0x7, Name: "Thundercaller's Gauntlets", Phase: 1, Quality: ItemQualityRare, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{16, 16, 18, 0, 35, 0, 0}},
	{Slot: 0x7, Name: "Gloves of the High Magus", Phase: 1, Quality: ItemQualityRare, SourceZone: "News of Victory - SMV Quest", SourceDrop: "", Stats: Stats{18, 13, 22, 0, 26, 0, 0}},
	{Slot: 0x7, Name: "Tempest's Touch", Phase: 1, Quality: ItemQualityRare, SourceZone: "Return to Andormu - CoT Quest", SourceDrop: "", Stats: Stats{20, 10, 0, 0, 27, 0, 0}, GemSlots: []GemColor{0x3, 0x3}, SocketBonus: Stats{}},
	{Slot: 0x7, Name: "Gloves of the Deadwatcher", Phase: 1, Quality: ItemQualityRare, SourceZone: "H AC - Shirrak the Dead Watcher", SourceDrop: "", Stats: Stats{24, 24, 0, 18, 29, 0, 0}},
	{Slot: 0x7, Name: "Incanter's Gloves", Phase: 1, Quality: ItemQualityRare, SourceZone: "SV - Thespia", SourceDrop: "", Stats: Stats{24, 21, 14, 0, 29, 0, 0}},
	{Slot: 0x7, Name: "Starlight Gauntlets", Phase: 1, Quality: ItemQualityRare, SourceZone: "N UB - Hungarfen", SourceDrop: "", Stats: Stats{21, 10, 0, 0, 25, 0, 0}, GemSlots: []GemColor{0x3, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: 0x7, Name: "Gloves of Oblivion", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Kargath", SourceDrop: "", Stats: Stats{21, 33, 0, 20, 26, 0, 0}},
	{Slot: 0x7, Name: "Harmony's Touch", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "Building a Perimeter - Netherstorm Quest", SourceDrop: "", Stats: Stats{0, 18, 16, 0, 33, 0, 0}},
	{Slot: 0x8, Name: "Girdle of Living Flame", Phase: 1, Quality: ItemQualityRare, SourceZone: "H UB - Hungarfen", SourceDrop: "", Stats: Stats{17, 15, 0, 16, 29, 0, 0}, GemSlots: []GemColor{0x4, 0x3}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: 0x8, Name: "Wave-Song Girdle", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H AC - Exarch Maladaar", SourceDrop: "", Stats: Stats{25, 25, 23, 0, 32, 0, 0}},
	{Slot: 0x8, Name: "A'dal's Gift", Phase: 1, Quality: ItemQualityRare, SourceZone: "How to Break Into the Arcatraz - Quest", SourceDrop: "", Stats: Stats{25, 0, 21, 0, 34, 0, 0}},
	{Slot: 0x8, Name: "Sash of Arcane Visions", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H AC - Exarch Maladaar", SourceDrop: "", Stats: Stats{23, 18, 22, 0, 28, 0, 0}},
	{Slot: 0x8, Name: "Belt of Depravity", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{27, 31, 0, 17, 34, 0, 0}},
	{Slot: 0x8, Name: "Moonrage Girdle", Phase: 1, Quality: ItemQualityRare, SourceZone: "SV - Hydromancer Thespia", SourceDrop: "", Stats: Stats{22, 0, 20, 0, 25, 0, 0}},
	{Slot: 0x8, Name: "Sash of Serpentra", Phase: 1, Quality: ItemQualityRare, SourceZone: "SV - Warlord Kalithresh", SourceDrop: "", Stats: Stats{21, 31, 0, 17, 25, 0, 0}},
	{Slot: 0x8, Name: "Blackwhelp Belt", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "Whelps of the Wyrmcult - BEM Quest", SourceDrop: "", Stats: Stats{11, 0, 10, 0, 32, 0, 0}},
	{Slot: 0x9, Name: "Spellstrike Pants", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{8, 12, 26, 22, 46, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{StatStm: 6}},
	{Slot: 0x9, Name: "Stormsong Kilt", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H UB - The Black Stalker", SourceDrop: "", Stats: Stats{30, 25, 26, 0, 35, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: 0x9, Name: "Tempest Leggings", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Mag'har - Revered (Horde)", SourceDrop: "", Stats: Stats{11, 0, 18, 0, 44, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatMP5: 2}},
	{Slot: 0x9, Name: "Kurenai Kilt", Phase: 1, Quality: ItemQualityRare, SourceZone: "Kurenai - Revered (Ally)", SourceDrop: "", Stats: Stats{11, 0, 18, 0, 44, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatMP5: 2}},
	{Slot: 0x9, Name: "Breeches of the Occultist", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H BM - Aeonus", SourceDrop: "", Stats: Stats{22, 37, 23, 0, 26, 0, 0}, GemSlots: []GemColor{0x4, 0x4, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: 0x9, Name: "Pantaloons of Flaming Wrath", Phase: 1, Quality: ItemQualityRare, SourceZone: "H SH - Blood Guard Porung", SourceDrop: "", Stats: Stats{28, 0, 42, 0, 33, 0, 0}},
	{Slot: 0x9, Name: "Moonchild Leggings", Phase: 1, Quality: ItemQualityRare, SourceZone: "H BF - Broggok", SourceDrop: "", Stats: Stats{20, 26, 21, 0, 23, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatMP5: 2}},
	{Slot: 0x9, Name: "Haramad's Leggings of the Third Coin", Phase: 1, Quality: ItemQualityRare, SourceZone: "Undercutting the Competition - MT Quest", SourceDrop: "", Stats: Stats{29, 0, 16, 0, 27, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x4}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: 0x9, Name: "Gladiator's Mail Leggings", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{25, 54, 22, 0, 42, 0, 6}},
	{Slot: 0x9, Name: "Kirin Tor Master's Trousers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H SLabs - Murmur", SourceDrop: "", Stats: Stats{29, 27, 0, 0, 36, 0, 0}, GemSlots: []GemColor{0x2, 0x4, 0x3}, SocketBonus: Stats{StatSpellHit: 4}},
	{Slot: 0x9, Name: "Khadgar's Kilt of Abjuration", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Temporus", SourceDrop: "", Stats: Stats{22, 20, 0, 0, 36, 0, 0}, GemSlots: []GemColor{0x4, 0x3, 0x3}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: 0x9, Name: "Incanter's Trousers", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{30, 25, 18, 0, 42, 0, 0}},
	{Slot: 0x9, Name: "Mana-Etched Pantaloons", Phase: 1, Quality: ItemQualityRare, SourceZone: "H UB - The Black Stalker", SourceDrop: "", Stats: Stats{32, 34, 21, 0, 33, 0, 0}},
	{Slot: 0x9, Name: "Tidefury Kilt", Phase: 1, Quality: ItemQualityRare, SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{31, 39, 19, 0, 35, 0, 0}},
	{Slot: 0x9, Name: "Molten Earth Kilt", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{32, 24, 0, 0, 40, 0, 10}},
	{Slot: 0x9, Name: "Trousers of Oblivion", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{33, 42, 0, 12, 39, 0, 0}},
	{Slot: 0x9, Name: "Leggings of the Third Coin", Phase: 1, Quality: ItemQualityRare, SourceZone: "Levixus the Soul Caller - Auchindoun Quest", SourceDrop: "", Stats: Stats{26, 34, 12, 0, 32, 0, 4}},
	{Slot: 0xa, Name: "Sigil-Laced Boots", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{18, 24, 17, 0, 20, 0, 0}, GemSlots: []GemColor{0x2, 0x4}, SocketBonus: Stats{StatInt: 3}},
	{Slot: 0xa, Name: "General's Mail Sabatons", Phase: 1, Quality: ItemQualityEpic, SourceZone: "11,424 Honor & 40 EotS Marks", SourceDrop: "", Stats: Stats{23, 34, 24, 0, 28, 0, 0}},
	{Slot: 0xa, Name: "Moonstrider Boots", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Darkweaver Syth", SourceDrop: "", Stats: Stats{22, 21, 20, 0, 25, 0, 6}},
	{Slot: 0xa, Name: "Shattarath Jumpers", Phase: 1, Quality: ItemQualityRare, SourceZone: "Into the Heart of the Labyrinth - Auch. Quest", SourceDrop: "", Stats: Stats{17, 25, 0, 0, 29, 0, 0}, GemSlots: []GemColor{0x4, 0x3}, SocketBonus: Stats{StatInt: 3}},
	{Slot: 0xa, Name: "Wave-Crest Striders", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H BF - Keli'dan the Breaker", SourceDrop: "", Stats: Stats{26, 28, 0, 0, 33, 0, 8}},
	{Slot: 0xa, Name: "Extravagant Boots of Malice", Phase: 1, Quality: ItemQualityRare, SourceZone: "H MT - Tavarok", SourceDrop: "", Stats: Stats{24, 27, 0, 14, 30, 0, 0}},
	{Slot: 0xa, Name: "Magma Plume Boots", Phase: 1, Quality: ItemQualityRare, SourceZone: "H AC - Shirrak the Dead Watcher", SourceDrop: "", Stats: Stats{26, 24, 0, 14, 29, 0, 0}},
	{Slot: 0xa, Name: "Shimmering Azure Boots", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "Securing the Celestial Ridge - NS Quest", SourceDrop: "", Stats: Stats{19, 0, 0, 16, 23, 0, 5}},
	{Slot: 0xa, Name: "Boots of Blashpemy", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{29, 36, 0, 0, 36, 0, 0}},
	{Slot: 0xa, Name: "Boots of Ethereal Manipulation", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H Bot - Warp Splinter", SourceDrop: "", Stats: Stats{27, 27, 0, 0, 33, 0, 0}},
	{Slot: 0xa, Name: "Earthbreaker's Greaves", Phase: 1, Quality: ItemQualityRare, SourceZone: "Levixus the Soul Caller - Auchindoun Quest", SourceDrop: "", Stats: Stats{20, 27, 8, 0, 25, 0, 3}},
	{Slot: 0xa, Name: "Boots of the Nexus Warden", Phase: 1, Quality: ItemQualityUncommon, SourceZone: "The Flesh Lies... - Netherstorm Quest", SourceDrop: "", Stats: Stats{17, 27, 0, 18, 21, 0, 0}},
	{Slot: 0xb, Name: "Sparking Arcanite Ring", Phase: 1, Quality: ItemQualityRare, SourceZone: "H OHF - Epoch Hunter", SourceDrop: "", Stats: Stats{14, 13, 14, 10, 22, 0, 0}},
	{Slot: 0xb, Name: "Seer's Signit", Phase: 1, Quality: ItemQualityEpic, SourceZone: "The Scryers - Exalted", SourceDrop: "", Stats: Stats{0, 24, 12, 0, 34, 0, 0}},
	{Slot: 0xb, Name: "Ring of Conflict Survival", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H MT - Yor (Summoned Boss)", SourceDrop: "", Stats: Stats{0, 28, 20, 0, 23, 0, 0}},
	{Slot: 0xb, Name: "Ryngo's Band of Ingenuity", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Wrath-Scryer Soccothrates", SourceDrop: "", Stats: Stats{14, 12, 14, 0, 25, 0, 0}},
	{Slot: 0xb, Name: "Band of the Guardian", Phase: 1, Quality: ItemQualityRare, SourceZone: "Hero of the Brood - CoT Quest", SourceDrop: "", Stats: Stats{11, 0, 17, 0, 23, 0, 0}},
	{Slot: 0xb, Name: "Scintillating Coral Band", Phase: 1, Quality: ItemQualityRare, SourceZone: "SV - Hydromancer Thespia", SourceDrop: "", Stats: Stats{15, 14, 17, 0, 21, 0, 0}},
	{Slot: 0xb, Name: "Manastorm Band", Phase: 1, Quality: ItemQualityRare, SourceZone: "Shutting Down Manaforge Ara - Quest", SourceDrop: "", Stats: Stats{15, 0, 10, 0, 29, 0, 0}},
	{Slot: 0xb, Name: "Ashyen's Gift", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Cenarion Expedition - Exalted", SourceDrop: "", Stats: Stats{0, 30, 0, 21, 23, 0, 0}},
	{Slot: 0xb, Name: "Cobalt Band of Tyrigosa", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H MT - Nexus-Prince Shaffar", SourceDrop: "", Stats: Stats{17, 19, 0, 0, 35, 0, 0}},
	{Slot: 0xb, Name: "Seal of the Exorcist", Phase: 1, Quality: ItemQualityEpic, SourceZone: "50 Spirit Shards ", SourceDrop: "", Stats: Stats{0, 24, 0, 12, 28, 0, 0}},
	{Slot: 0xb, Name: "Lola's Eve", Phase: 1, Quality: ItemQualityEpic, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{14, 15, 0, 0, 29, 0, 0}},
	{Slot: 0xb, Name: "Yor's Collapsing Band", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H MT - Yor (Summoned Boss)", SourceDrop: "", Stats: Stats{20, 0, 0, 0, 23, 0, 0}},
	{Slot: 0x13, Name: "Totem of the Void", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Cache of the Legion", SourceDrop: "", Stats: Stats{StatSpellDmg: 55}},
	// {Slot: 0x13, Name: "Totem of the Pulsing Earth", Phase: 1, Quality: ItemQualityRare, SourceZone: "15 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	// {Slot: 0x13, Name: "Totem of Impact", Phase: 1, Quality: ItemQualityRare, SourceZone: "15 Mark of Thrallmar/ Honor Hold", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	// {Slot: 0x13, Name: "Totem of Lightning", Phase: 1, Quality: ItemQualityRare, SourceZone: "Colossal Menace - HFP Quest", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 0x11, Name: "Gladiator's Spellblade / Gavel", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{18, 28, 0, 0, 199, 0, 0}},
	{Slot: 0x11, Name: "Starlight Dagger", Phase: 1, Quality: ItemQualityRare, SourceZone: "H SP - Mennu the Betrayer", SourceDrop: "", Stats: Stats{15, 15, 0, 16, 121, 0, 0}},
	{Slot: 0x11, Name: "Runesong Dagger", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Warbringer O'mrogg", SourceDrop: "", Stats: Stats{11, 12, 20, 0, 121, 0, 0}},
	{Slot: 0x11, Name: "Bleeding Hollow Warhammer", Phase: 1, Quality: ItemQualityRare, SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{17, 12, 16, 0, 121, 0, 0}},
	{Slot: 0x11, Name: "Sky Breaker", Phase: 1, Quality: ItemQualityRare, SourceZone: "H AC - Avatar of the Martyred", SourceDrop: "", Stats: Stats{20, 13, 0, 0, 132, 0, 0}},
	{Slot: 0x12, Name: "Lamp of Peaceful Raidiance", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{14, 13, 13, 12, 21, 0, 0}},
	{Slot: 0x12, Name: "Manual of the Nethermancer", Phase: 1, Quality: ItemQualityRare, SourceZone: "Mech - Nethermancer Sepethrea", SourceDrop: "", Stats: Stats{15, 12, 19, 0, 21, 0, 0}},
	{Slot: 0x12, Name: "Draenei Honor Guard Shield", Phase: 1, Quality: ItemQualityRare, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{16, 0, 21, 0, 19, 0, 0}},
	{Slot: 0x12, Name: "Star-Heart Lamp", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Temporus", SourceDrop: "", Stats: Stats{18, 17, 0, 12, 22, 0, 0}},
	{Slot: 0x12, Name: "The Saga of Terokk", Phase: 1, Quality: ItemQualityRare, SourceZone: "Terokk's Legacy - Auchindoun Quest", SourceDrop: "", Stats: Stats{23, 0, 0, 0, 28, 0, 0}},
	{Slot: 0x12, Name: "Silvermoon Crest Shield", Phase: 1, Quality: ItemQualityRare, SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{20, 0, 0, 0, 23, 0, 5}},
	{Slot: 0x12, Name: "Spellbreaker's Buckler", Phase: 1, Quality: ItemQualityRare, SourceZone: "Akama's Promise - SMV Quest", SourceDrop: "", Stats: Stats{10, 22, 0, 0, 29, 0, 0}},
	{Slot: 0x12, Name: "Hortus' Seal of Brilliance", Phase: 1, Quality: ItemQualityRare, SourceZone: "SH - Warchief Kargath Bladefist", SourceDrop: "", Stats: Stats{20, 18, 0, 0, 23, 0, 0}},
	{Slot: 0x12, Name: "Gladiator's Endgame", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{14, 21, 0, 0, 19, 0, 0}},
	{Slot: 0x11, Name: "Gladiator's War Staff", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{35, 48, 36, 21, 199, 0, 0}},
	{Slot: 0x11, Name: "Terokk's Shadowstaff", Phase: 1, Quality: ItemQualityEpic, SourceZone: "H SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{42, 40, 37, 0, 168, 0, 0}},
	{Slot: 0x11, Name: "Auchenai Staff", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Aldor - Revered", SourceDrop: "", Stats: Stats{46, 0, 26, 19, 121, 0, 0}},
	{Slot: 0x11, Name: "Warpstaff of Arcanum", Phase: 1, Quality: ItemQualityRare, SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{38, 37, 26, 16, 121, 0, 0}},
	{Slot: 0x11, Name: "The Bringer of Death", Phase: 1, Quality: ItemQualityRare, SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{31, 32, 42, 0, 121, 0, 0}},
	{Slot: 0x11, Name: "Bloodfire Greatstaff", Phase: 1, Quality: ItemQualityRare, SourceZone: "BM - Aeonus", SourceDrop: "", Stats: Stats{42, 42, 28, 0, 121, 0, 0}},
	{Slot: 0x11, Name: "Ameer's Impulse Taser", Phase: 1, Quality: ItemQualityRare, SourceZone: "Nexus-King Salhadaar - Netherstorm Quest", SourceDrop: "", Stats: Stats{27, 27, 27, 17, 103, 0, 0}},
	{Slot: 0x11, Name: "Grand Scepter of the Nexus-Kings", Phase: 1, Quality: ItemQualityRare, SourceZone: "H MT - Nexus-Prince Shaffar", SourceDrop: "", Stats: Stats{43, 45, 0, 19, 121, 0, 0}},

	// source: https://docs.google.com/spreadsheets/d/1T4DEuq0yroEPb-11okC3qjj7aYfCGu2e6nT9LeT30zg/edit#gid=0
	{Slot: EquipHead, Name: "Uni-Mind Headdress", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Netherspite", Stats: Stats{StatStm: 31, StatInt: 40, StatSpellDmg: 46, StatSpellCrit: 25, StatSpellHit: 19}},
	{Slot: EquipHead, Name: "Wicked Witch's Hat", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Opera", Stats: Stats{StatStm: 37, StatInt: 38, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 32, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipHead, Name: "Cyclone Faceguard (Tier 4)", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Prince", Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 8}, GemSlots: []GemColor{GemColorMeta, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Cataclysm Headpiece (Tier 5)", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Lady Vashj", Stats: Stats{StatStm: 35, StatInt: 28, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 18, StatMP5: 7}, GemSlots: []GemColor{GemColorMeta, GemColorYellow}, SocketBonus: Stats{StatSpellHit: 5}},
	{Slot: EquipHead, Name: "Cowl of the Grand Engineer", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Void Reaver", Stats: Stats{StatStm: 22, StatInt: 27, StatSpellDmg: 53, StatHaste: 0, StatSpellCrit: 35, StatSpellHit: 16, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Magnified Moon Specs", Phase: 2, Quality: ItemQualityEpic, SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Leather)", Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 41, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorMeta, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Gadgetstorm Goggles", Phase: 2, Quality: ItemQualityEpic, SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Mail)", Stats: Stats{StatStm: 28, StatInt: 0, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 40, StatSpellHit: 12, StatMP5: 0}, GemSlots: []GemColor{GemColorMeta, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Destruction Holo-gogs", Phase: 2, Quality: ItemQualityEpic, SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Cloth)", Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 64, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorMeta, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Skyshatter Headguard (Tier 6)", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Archimonde", Stats: Stats{StatStm: 42, StatInt: 37, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 36, StatSpellHit: 0, StatMP5: 8}, GemSlots: []GemColor{GemColorMeta, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipHead, Name: "Cowl of the Illidari High Lord", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Illidan", Stats: Stats{StatStm: 33, StatInt: 31, StatSpellDmg: 64, StatHaste: 0, StatSpellCrit: 47, StatSpellHit: 21, StatMP5: 0}, GemSlots: []GemColor{GemColorMeta, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipNeck, Name: "Brooch of Unquenchable Fury", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Moroes", Stats: Stats{StatStm: 24, StatInt: 21, StatSpellDmg: 26, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 15, StatMP5: 0}},
	{Slot: EquipNeck, Name: "Manasurge Pendant", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Badges", Stats: Stats{StatStm: 24, StatInt: 22, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipNeck, Name: "Pendant of the Lost Ages", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Tidewalker", Stats: Stats{StatStm: 27, StatInt: 17, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipNeck, Name: "Adornment of Stolen Souls", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Prince", Stats: Stats{StatStm: 18, StatInt: 20, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipNeck, Name: "The Sun King's Talisman", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Kael Reward", Stats: Stats{StatStm: 22, StatInt: 16, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipNeck, Name: "Translucent Spellthread Necklace", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "RoS", Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 15, StatMP5: 0}},
	{Slot: EquipShoulder, Name: "Mantle of the Mind Flayer", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Aran", Stats: Stats{StatStm: 33, StatInt: 29, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipShoulder, Name: "Mantle of the Elven Kings", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Trash", Stats: Stats{StatStm: 27, StatInt: 18, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipShoulder, Name: "Cyclone Shoulderguards (Tier 4)", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Gruul's Lair", SourceDrop: "Maulgar", Stats: Stats{StatStm: 28, StatInt: 26, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 12, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipShoulder, Name: "Illidari Shoulderpads", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Tidewalker", Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 16, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipShoulder, Name: "Blood-cursed Shoulderpads", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Bloodboil", Stats: Stats{StatStm: 25, StatInt: 19, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipShoulder, Name: "Cataclysm Shoulderpads (Tier 5)", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "VoidReaver", Stats: Stats{StatStm: 26, StatInt: 19, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 6}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: EquipShoulder, Name: "Mantle of Nimble Thought", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Tailoring", Stats: Stats{StatStm: 37, StatInt: 26, StatSpellDmg: 44, StatHaste: 38, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipShoulder, Name: "Skyshatter Mantle (Tier 6)", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Mother", Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 27, StatSpellHit: 11, StatMP5: 4}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipShoulder, Name: "Hatefury Mantle", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Anetheron", Stats: Stats{StatStm: 15, StatInt: 18, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: EquipBack, Name: "Ruby Drape of the Mysticant", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Prince", Stats: Stats{StatStm: 22, StatInt: 21, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipBack, Name: "Shadow-Cloak of Dalaran", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Moroes", Stats: Stats{StatStm: 19, StatInt: 18, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipBack, Name: "Shawl of Shifting Probabilities", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Badges", Stats: Stats{StatStm: 18, StatInt: 16, StatSpellDmg: 21, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipBack, Name: "Royal Cloak of the Sunstriders", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Kaelthas", Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipBack, Name: "Brute Cloak of the Ogre-Magi", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Gruul's Lair", SourceDrop: "Maulgar", Stats: Stats{StatStm: 18, StatInt: 20, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipBack, Name: "Ancient Spellcloak of the Highborne", Phase: 1, Quality: ItemQualityEpic, SourceZone: "WorldBoss", SourceDrop: "Kazzak", Stats: Stats{StatStm: 0, StatInt: 15, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipBack, Name: "Cloak of the Illidari Council", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "IllidariCouncil", Stats: Stats{StatStm: 24, StatInt: 16, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipChest, Name: "Cyclone Chestguard (Tier 4)", Phase: 1, Quality: ItemQualityEpic, SourceZone: "GruulsLair", SourceDrop: "Maulgar", Stats: Stats{StatStm: 33, StatInt: 32, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 8}, GemSlots: []GemColor{GemColorRed, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellHit: 4}},
	{Slot: EquipChest, Name: "Netherstrike Breastplate", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Crafted", SourceDrop: "Leatherworking", Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 32, StatSpellHit: 0, StatMP5: 8}, GemSlots: []GemColor{GemColorBlue, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipChest, Name: "Robe of Hateful Echoes", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Hydross", Stats: Stats{StatStm: 34, StatInt: 36, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorRed, GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatStm: 6}},
	{Slot: EquipChest, Name: "Robe of the Shadow Council", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Teron", Stats: Stats{StatStm: 37, StatInt: 36, StatSpellDmg: 73, StatHaste: 0, StatSpellCrit: 28, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipChest, Name: "Robes of Rhonin", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Archimonde", Stats: Stats{StatStm: 55, StatInt: 38, StatSpellDmg: 81, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 27, StatMP5: 0}},
	{Slot: EquipChest, Name: "Cataclysm Chestpiece (Tier 5)", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Kaelthas", Stats: Stats{StatStm: 37, StatInt: 28, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 10}, GemSlots: []GemColor{GemColorBlue, GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipChest, Name: "Vestments of the Sea-Witch", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "LadyVashj", Stats: Stats{StatStm: 28, StatInt: 28, StatSpellDmg: 57, StatHaste: 0, StatSpellCrit: 31, StatSpellHit: 27, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipChest, Name: "Chestguard of Relentless Storms", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Trash", Stats: Stats{StatStm: 36, StatInt: 30, StatSpellDmg: 74, StatHaste: 0, StatSpellCrit: 46, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipChest, Name: "Skyshatter Breastplate (Tier 6)", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Illidan", Stats: Stats{StatStm: 42, StatInt: 41, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 27, StatSpellHit: 17, StatMP5: 7}, GemSlots: []GemColor{GemColorBlue, GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipWrist, Name: "Bands of Nefarious Deeds", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Maiden", Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 32, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWrist, Name: "Elunite Empowered Bracers", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "RoS", Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 6}},
	{Slot: EquipWrist, Name: "Focused Mana Bindings", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Akama", Stats: Stats{StatStm: 27, StatInt: 20, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 0}},
	{Slot: EquipWrist, Name: "Netherstrike Bracers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Crafted", SourceDrop: "Leatherworking", Stats: Stats{StatStm: 13, StatInt: 13, StatSpellDmg: 20, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 6}, GemSlots: []GemColor{GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 2}},
	{Slot: EquipWrist, Name: "Bands of the Coming Storm", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Supremus", Stats: Stats{StatStm: 28, StatInt: 28, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWrist, Name: "Mindstorm Wristbands", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Alar", Stats: Stats{StatStm: 13, StatInt: 13, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWrist, Name: "Cuffs of Devastation", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Winterchill", Stats: Stats{StatStm: 22, StatInt: 20, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 14, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow}, SocketBonus: Stats{StatStm: 3}},
	{Slot: EquipHands, Name: "Cyclone Handguards (Tier 4)", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Curator", Stats: Stats{StatStm: 26, StatInt: 29, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 6}},
	{Slot: EquipHands, Name: "Handwraps of Flowing Thought", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Huntsman", Stats: Stats{StatStm: 24, StatInt: 22, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 14, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellHit: 3}},
	{Slot: EquipHands, Name: "Cataclysm Handgrips (Tier 5)", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "LeotherastheBlind", Stats: Stats{StatStm: 25, StatInt: 27, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 19, StatMP5: 7}},
	{Slot: EquipHands, Name: "Gauntlets of the Sun King", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Kaelthas", Stats: Stats{StatStm: 28, StatInt: 29, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 28, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipHands, Name: "Anger-Spark Gloves", Phase: 1, Quality: ItemQualityEpic, SourceZone: "World Boss", SourceDrop: "Doomwalker", Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 20, StatMP5: 0}, GemSlots: []GemColor{GemColorRed, GemColorRed}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: EquipHands, Name: "Soul-Eater's Handwraps", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Magtheridon's Lair", SourceDrop: "Magtheridon", Stats: Stats{StatStm: 31, StatInt: 24, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipHands, Name: "Skyshatter Guantlets (Tier 6)", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Azgalor", Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 19, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 2}},
	{Slot: EquipWaist, Name: "Nethershard Girdle", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Moroes", Stats: Stats{StatStm: 22, StatInt: 30, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWaist, Name: "General's Mail Girdle", Phase: 1, Quality: ItemQualityEpic, SourceZone: "PvP", SourceDrop: "PvP", Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWaist, Name: "Malefic Girdle", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Illhoof", Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWaist, Name: "Monsoon Belt", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC/TK", SourceDrop: "Leatherworking", Stats: Stats{StatStm: 23, StatInt: 24, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 21, StatMP5: 0}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipWaist, Name: "Netherstrike Belt", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Crafted", SourceDrop: "Leatherworking", Stats: Stats{StatStm: 10, StatInt: 17, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 16, StatSpellHit: 0, StatMP5: 9}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellCrit: 3}},
	{Slot: EquipWaist, Name: "Belt of Divine Inspiration", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Gruul's Lair", SourceDrop: "Maulgar", Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipWaist, Name: "Cord of Screaming Terrors", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Lurker", Stats: Stats{StatStm: 34, StatInt: 15, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 24, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatStm: 4}},
	{Slot: EquipWaist, Name: "Girdle of Ruination", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Crafted", SourceDrop: "Tailoring", Stats: Stats{StatStm: 18, StatInt: 13, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorRed, GemColorYellow}, SocketBonus: Stats{StatStm: 4}},
	{Slot: EquipWaist, Name: "Belt of the Crescent Moon", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Kazrogal", Stats: Stats{StatStm: 25, StatInt: 27, StatSpellDmg: 44, StatHaste: 36, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWaist, Name: "Waistwrap of Infinity", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Supremus", Stats: Stats{StatStm: 31, StatInt: 22, StatSpellDmg: 56, StatHaste: 32, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWaist, Name: "Belt of Blasting", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC/TK", SourceDrop: "Tailoring", Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 30, StatSpellHit: 23, StatMP5: 0}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipWaist, Name: "Anetheron's Noose", Phase: 2, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Anetheron", Stats: Stats{StatStm: 22, StatInt: 23, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipWaist, Name: "Flashfire Girdle", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Akama", Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 44, StatHaste: 37, StatSpellCrit: 18, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipLegs, Name: "Cyclone Legguards (Tier 4)", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Gruul's Lair", SourceDrop: "Gruul", Stats: Stats{StatStm: 40, StatInt: 40, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 20, StatMP5: 8}},
	{Slot: EquipLegs, Name: "Trial-Fire Trousers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Opera", Stats: Stats{StatStm: 42, StatInt: 40, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipLegs, Name: "Trousers of the Astromancer", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Solarian", Stats: Stats{StatStm: 33, StatInt: 36, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorBlue, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipLegs, Name: "Cataclysm Leggings (Tier 5)", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Karathress", Stats: Stats{StatStm: 48, StatInt: 46, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 14, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 2}},
	{Slot: EquipLegs, Name: "Leggings of Devastation", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Mother", Stats: Stats{StatStm: 40, StatInt: 42, StatSpellDmg: 60, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 26, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipLegs, Name: "Skyshatter Legguards (Tier 6)", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "IllidariCouncil", Stats: Stats{StatStm: 40, StatInt: 42, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 20, StatMP5: 11}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipLegs, Name: "Leggings of the Seventh Circle", Phase: 1, Quality: ItemQualityEpic, SourceZone: "World Boss", SourceDrop: "Kazzak", Stats: Stats{StatStm: 0, StatInt: 22, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}, GemSlots: []GemColor{GemColorRed, GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipLegs, Name: "Leggings of Channeled Elements", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Kazrogal", Stats: Stats{StatStm: 25, StatInt: 28, StatSpellDmg: 59, StatHaste: 0, StatSpellCrit: 34, StatSpellHit: 18, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipFeet, Name: "Boots of the Infernal Coven", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Aran", Stats: Stats{StatStm: 27, StatInt: 27, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Ruby Slippers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Opera", Stats: Stats{StatStm: 33, StatInt: 29, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 16, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Windshear Boots", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Gruul's Lair", SourceDrop: "Gruul", Stats: Stats{StatStm: 37, StatInt: 32, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Blue Suede Shoes", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Kazrogal", Stats: Stats{StatStm: 37, StatInt: 32, StatSpellDmg: 56, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Boots of Blasting", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC/TK", SourceDrop: "Tailoring", Stats: Stats{StatStm: 25, StatInt: 25, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Boots of Foretelling", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Maiden", Stats: Stats{StatStm: 27, StatInt: 23, StatSpellDmg: 26, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorRed, GemColorYellow}, SocketBonus: Stats{StatInt: 3}},
	{Slot: EquipFeet, Name: "Hurricane Boots", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC/TK", SourceDrop: "Leatherworking", Stats: Stats{StatStm: 25, StatInt: 26, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 6}},
	{Slot: EquipFeet, Name: "Velvet Boots of the Guardian", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Lurker", Stats: Stats{StatStm: 21, StatInt: 21, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Boots of Oceanic Fury", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Najentus", Stats: Stats{StatStm: 28, StatInt: 36, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFeet, Name: "Naturewarden's Treads", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "RoS", Stats: Stats{StatStm: 39, StatInt: 18, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 7}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipFeet, Name: "Slippers of the Seacaller", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Najentus", Stats: Stats{StatStm: 25, StatInt: 18, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 0, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipFinger, Name: "Band of Crimson Fury", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Magtheridon's Lair", SourceDrop: "MagtheridonQuest", Stats: Stats{StatStm: 22, StatInt: 22, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 16, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Spectral Band of Innervation", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Huntsman", Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 29, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Band of Alar", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Alar", Stats: Stats{StatStm: 24, StatInt: 23, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Cryptic Dreams", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Badges", Stats: Stats{StatStm: 16, StatInt: 17, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Violet Signet of the Archmage", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Exalted", Stats: Stats{StatStm: 24, StatInt: 23, StatSpellDmg: 29, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Recurrence", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Chess", Stats: Stats{StatStm: 15, StatInt: 15, StatSpellDmg: 32, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Band of the Eternal Sage", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Exalted", Stats: Stats{StatStm: 28, StatInt: 25, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Endless Coils", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "LadyVashj", Stats: Stats{StatStm: 31, StatInt: 0, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Unrelenting Storms", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Trash", Stats: Stats{StatStm: 0, StatInt: 15, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Captured Storms", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Najentus", Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 19, StatMP5: 0}},
	{Slot: EquipFinger, Name: "Ring of Ancient Knowledge", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Trash", Stats: Stats{StatStm: 30, StatInt: 20, StatSpellDmg: 39, StatHaste: 31, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Gavel of Unearthed Secrets", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Lower City - Exalted", Stats: Stats{StatStm: 24, StatInt: 16, StatSpellDmg: 159, StatHaste: 0, StatSpellCrit: 15, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Eternium Runed Blade", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Crafted", SourceDrop: "Blacksmithing", Stats: Stats{StatStm: 0, StatInt: 19, StatSpellDmg: 168, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Gladiator's Gavel / Gladiator's Spellblade", Phase: 1, Quality: ItemQualityEpic, SourceZone: "PvP", SourceDrop: "PvP", Stats: Stats{StatStm: 28, StatInt: 18, StatSpellDmg: 199, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Nathrezim Mindblade", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Prince", Stats: Stats{StatStm: 18, StatInt: 18, StatSpellDmg: 203, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Talon of the Tempest", Phase: 1, Quality: ItemQualityEpic, SourceZone: "World Boss", SourceDrop: "Doomwalker", Stats: Stats{StatStm: 0, StatInt: 10, StatSpellDmg: 194, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 9, StatMP5: 0}, GemSlots: []GemColor{GemColorYellow, GemColorYellow}, SocketBonus: Stats{StatInt: 3}},
	{Slot: EquipWeapon, Name: "Hammer of Judgement", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Trash", Stats: Stats{StatStm: 33, StatInt: 22, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 22, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "The Maelstrom's Fury", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Najentus", Stats: Stats{StatStm: 33, StatInt: 21, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Staff of Infinite Mysteries", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Curator", Stats: Stats{StatStm: 61, StatInt: 51, StatSpellDmg: 185, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 23, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "The Nexus Key", Phase: 2, Quality: ItemQualityEpic, SourceZone: "TK", SourceDrop: "Kaelthas", Stats: Stats{StatStm: 76, StatInt: 52, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 51, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipWeapon, Name: "Zhar'doom, Greatstaff of the Devourer", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Illidan", Stats: Stats{StatStm: 70, StatInt: 47, StatSpellDmg: 259, StatHaste: 55, StatSpellCrit: 36, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Jewel of Infinite Possibilities", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Netherspite", Stats: Stats{StatStm: 19, StatInt: 18, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 21, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Dragonheart Flameshield", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Nightbane", Stats: Stats{StatStm: 19, StatInt: 33, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 7}},
	{Slot: EquipOffhand, Name: "Illidari Runeshield", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Trash", Stats: Stats{StatStm: 45, StatInt: 39, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Karaborian Talisman", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Magtheridon's Lair", SourceDrop: "Magtheridon", Stats: Stats{StatStm: 23, StatInt: 23, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Mazthoril Honor Shield", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Badges", Stats: Stats{StatStm: 16, StatInt: 29, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Talisman of Nightbane", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Nightbane", Stats: Stats{StatStm: 19, StatInt: 19, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Blind-Seers Icon", Phase: 3, Quality: ItemQualityEpic, SourceZone: "BT", SourceDrop: "Akama", Stats: Stats{StatStm: 25, StatInt: 16, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 24, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Khadgar's Knapsack", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrah", SourceDrop: "Badges", Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "FathomStone", Phase: 2, Quality: ItemQualityEpic, SourceZone: "SSC", SourceDrop: "Lurker", Stats: Stats{StatStm: 16, StatInt: 12, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Antonidas's Aegis of Rapt Concentration", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Archimonde", Stats: Stats{StatStm: 28, StatInt: 32, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}},
	{Slot: EquipOffhand, Name: "Chronicle of Dark Secrets", Phase: 3, Quality: ItemQualityEpic, SourceZone: "Hyjal", SourceDrop: "Winterchill", Stats: Stats{StatStm: 16, StatInt: 12, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 17, StatMP5: 0}},

	// Hand Written
	{Slot: EquipTrinket, Name: "Quagmirran's Eye", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Slave Pens", SourceDrop: "Quagmirran", Stats: Stats{StatSpellDmg: 37}, Activate: ActivateQuagsEye, ActivateCD: -1}, // -1 will trigger an activation only once
	{Slot: EquipTrinket, Name: "Icon of the Silver Crescent", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Shattrath", SourceDrop: "G'eras - 41 Badges", Stats: Stats{StatSpellDmg: 44}, Activate: createSpellDmgActivate(MagicIDBlessingSilverCrescent, 155, 20), ActivateCD: 120, CoolID: MagicIDISCTrink},
	{Slot: EquipTrinket, Name: "Natural Alignment Crystal", Phase: 0, Quality: ItemQualityEpic, SourceZone: "BWL", SourceDrop: "", Stats: Stats{}, Activate: ActivateNAC, ActivateCD: 300, CoolID: MagicIDNACTrink},
	{Slot: EquipTrinket, Name: "Neltharion's Tear", Phase: 0, Quality: ItemQualityEpic, SourceZone: "BWL", SourceDrop: "Nefarian", Stats: Stats{StatSpellDmg: 44, StatSpellHit: 16}},
	{Slot: EquipTrinket, Name: "Mark of the Champion", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "KT", Stats: Stats{StatSpellDmg: 85}},
	{Slot: EquipTrinket, Name: "Scryer's Bloodgem", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Scryers - Revered", SourceDrop: "", Stats: Stats{0, 0, 0, 32, 0, 0, 0}, Activate: createSpellDmgActivate(MagicIDSpellPower, 150, 15), ActivateCD: 90, CoolID: MagicIDScryerTrink},
	{Slot: EquipTrinket, Name: "Figurine - Living Ruby Serpent", Phase: 1, Quality: ItemQualityRare, SourceZone: "Jewelcarfting BoP", SourceDrop: "", Stats: Stats{23, 33, 0, 0, 0, 0, 0}, Activate: createSpellDmgActivate(MagicIDRubySerpent, 150, 20), ActivateCD: 300, CoolID: MagicIDRubySerpentTrink},
	{Slot: EquipTrinket, Name: "Xi'ri's Gift", Phase: 1, Quality: ItemQualityRare, SourceZone: "The Sha'tar - Revered", SourceDrop: "", Stats: Stats{0, 0, 32, 0, 0, 0, 0}, Activate: createSpellDmgActivate(MagicIDSpellPower, 150, 15), ActivateCD: 90, CoolID: MagicIDXiriTrink},
	{Slot: EquipTrinket, Name: "Shiffar's Nexus-Horn", Phase: 1, Quality: ItemQualityRare, SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{0, 0, 30, 0, 0, 0, 0}, Activate: ActivateNexusHorn, ActivateCD: -1},
	{Slot: EquipTrinket, Name: "Darkmoon Card: Crusade", Phase: 2, Quality: ItemQualityEpic, SourceZone: "Blessings Deck", SourceDrop: "", Activate: ActivateDCC, ActivateCD: -1},
	{Slot: EquipTrinket, Name: "The Lightning Capacitor", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "", Activate: ActivateTLC, ActivateCD: -1},
	{Slot: EquipTrinket, Name: "Eye of Magtheridon", Phase: 1, Quality: ItemQualityRare, SourceZone: "", SourceDrop: "", Stats: Stats{StatSpellDmg: 54}, Activate: ActivateEyeOfMag, ActivateCD: -1},
	{Slot: EquipTrinket, Name: "Sextant of Unstable Currents", Phase: 2, Quality: ItemQualityRare, SourceZone: "SSC", SourceDrop: "", Stats: Stats{StatSpellCrit: 40}, Activate: ActivateSextant, ActivateCD: -1},
	{Slot: EquipTrinket, Name: "Shifting Naaru Silver", Phase: 5, Quality: ItemQualityRare, SourceZone: "Sunwell", SourceDrop: "", Stats: Stats{StatHaste: 54}, Activate: createSpellDmgActivate(MagicIDShiftingNaaru, 320, 15), ActivateCD: 90, CoolID: MagicIDShiftingNaaruTrink},
	{Slot: EquipTrinket, Name: "The Skull of Gul'dan", Phase: 3, Quality: ItemQualityRare, SourceZone: "Black Temple", SourceDrop: "", Stats: Stats{StatSpellHit: 25, StatSpellDmg: 55}, Activate: createHasteActivate(MagicIDSkullGuldan, 175, 20), ActivateCD: 120, CoolID: MagicIDSkullGuldanTrink},
	{Slot: EquipTrinket, Name: "Hex Shrunken Head", Phase: 4, Quality: ItemQualityRare, SourceZone: "ZA", SourceDrop: "", Stats: Stats{StatSpellDmg: 53}, Activate: createSpellDmgActivate(MagicIDHexShunkHead, 211, 20), ActivateCD: 120, CoolID: MagicIDHexTrink},
	{Slot: EquipNeck, Name: "Eye of the Night", Phase: 1, Quality: ItemQualityRare, SourceZone: "Jewelcrafting", SourceDrop: "", Stats: Stats{StatSpellCrit: 26, StatSpellHit: 16, StatSpellPen: 15}, Activate: func(sim *Simulation) Aura {
		if sim.Options.Buffs.EyeOfNight {
			return Aura{}
		}
		activate := createSpellDmgActivate(MagicIDEyeOfTheNight, 34, 30*60)
		return activate(sim)
	}, ActivateCD: 3600, CoolID: MagicIDEyeOfTheNightTrink},
	{Slot: EquipNeck, Name: "Chain of the Twilight Owl", Phase: 1, Quality: ItemQualityRare, SourceZone: "Jewelcrafting", SourceDrop: "", Stats: Stats{StatStm: 0, StatInt: 19, StatSpellDmg: 21, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}, Activate: ActivateChainTO, ActivateCD: 3600, CoolID: MagicIDChainTOTrink},
	{Slot: EquipFinger, Name: "Evoker's Mark of the Redemption", Phase: 1, Quality: ItemQualityRare, SourceZone: "Quest SMV", SourceDrop: "Dissension Amongst the Ranks...", Stats: Stats{StatInt: 15, StatSpellDmg: 29, StatSpellCrit: 10}},
	{Slot: EquipFinger, Name: "Dreamcrystal Band", Phase: 1, Quality: ItemQualityRare, SourceZone: "Blades Edge Moutains", SourceDrop: "50 Apexis Shards", Stats: Stats{StatInt: 10, StatSpellDmg: 38, StatSpellCrit: 15}},
	{Slot: EquipChest, Name: "Windhawk Hauberk", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Leatherworking", SourceDrop: "", Stats: Stats{StatStm: 28, StatInt: 29, StatSpirit: 29, StatSpellDmg: 46, StatSpellCrit: 19}, GemSlots: []GemColor{GemColorBlue, GemColorYellow, GemColorBlue}, SocketBonus: Stats{StatSpellDmg: 5}},
	{Slot: EquipWaist, Name: "Windhawk Belt", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Leatherworking", SourceDrop: "", Stats: Stats{StatStm: 17, StatInt: 19, StatSpirit: 20, StatSpellDmg: 37, StatSpellCrit: 12}, GemSlots: []GemColor{GemColorBlue, GemColorYellow}, SocketBonus: Stats{StatSpellDmg: 4}},
	{Slot: EquipWrist, Name: "Windhawk Bracers", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Leatherworking", SourceDrop: "", Stats: Stats{StatStm: 22, StatInt: 17, StatSpirit: 7, StatSpellDmg: 27, StatSpellCrit: 16}, GemSlots: []GemColor{GemColorYellow}, SocketBonus: Stats{StatInt: 2}},
	{Slot: EquipHands, Name: "Tidefury Gauntlets", Phase: 1, Quality: ItemQualityRare, SourceZone: "", SourceDrop: "", Stats: Stats{StatStm: 22, StatInt: 26, StatSpellDmg: 29, StatMP5: 7}},
	{Slot: EquipWaist, Name: "Eyestalk Waist Cord", Phase: 0, Quality: ItemQualityEpic, SourceZone: "AQ40", SourceDrop: "C'thun", Stats: Stats{StatStm: 10, StatInt: 9, StatSpellDmg: 41, StatSpellCrit: 14}},
	{Slot: EquipLegs, Name: "Leggings of Polarity", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "Thaddius", Stats: Stats{StatStm: 20, StatInt: 14, StatSpellDmg: 44, StatSpellCrit: 28}},
	{Slot: EquipFinger, Name: "Ring of the Fallen God", Phase: 0, Quality: ItemQualityEpic, SourceZone: "AQ40", SourceDrop: "C'thun", Stats: Stats{StatStm: 5, StatInt: 6, StatSpellDmg: 37, StatSpellCrit: 8}},
	{Slot: EquipFinger, Name: "Band of the Inevitable", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "Noth", Stats: Stats{StatSpellDmg: 36, StatSpellHit: 8}},
	{Slot: EquipFinger, Name: "Seal of the Damned", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "Noth", Stats: Stats{StatStm: 17, StatSpellDmg: 21, StatSpellCrit: 14, StatSpellHit: 8}},
	{Slot: EquipShoulder, Name: "Pauldrons of Elemental Fury", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "Trash", Stats: Stats{StatStm: 19, StatInt: 21, StatSpellDmg: 26, StatSpellCrit: 14, StatSpellHit: 8}},
	{Slot: EquipBack, Name: "Cloak of the Necropolis", Phase: 0, Quality: ItemQualityEpic, SourceZone: "Naxx", SourceDrop: "Sapp", Stats: Stats{StatStm: 12, StatInt: 11, StatSpellDmg: 26, StatSpellCrit: 14, StatSpellHit: 8}},

	{Slot: EquipFeet, Name: "Glider's Sabatons of Nature's Wrath", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 78}},
	{Slot: EquipWaist, Name: "Lurker's Belt of Nature's Wrath", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 78}},
	{Slot: EquipWrist, Name: "Ravager's Bands of Nature's Wrath", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 58}},
	{Slot: EquipFeet, Name: "Glider's Sabatons of the Invoker", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 33, StatSpellCrit: 28}},
	{Slot: EquipWaist, Name: "Lurker's Belt of the Invoker", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 33, StatSpellCrit: 28}},
	{Slot: EquipWrist, Name: "Ravager's Bands of the Invoker", Phase: 1, Quality: ItemQualityEpic, SourceZone: "Kara", SourceDrop: "Beast Trash?", Stats: Stats{StatSpellDmg: 25, StatSpellCrit: 21}},

	// {Slot:EquipTrinket, Name:"Arcanist's Stone", Phase: 1, Quality: ItemQualityEpic, SourceZone:"H OHF - Epoch Hunter", SourceDrop:"", Stats:Stats{0, 0, 0, 25, 0, 0, 0} }
	// {Slot:EquipTrinket, Name:"Vengeance of the Illidari", Phase: 1, Quality: ItemQualityEpic, SourceZone:"Cruel's Intentions/Overlord - HFP Quest", SourceDrop:"", Stats:Stats{0, 0, 26, 0, 0, 0, 0} }
	{Slot: EquipTotem, Name: "Totem of Ancestral Guidance", Phase: 3, Quality: ItemQualityRare, SourceZone: "BT", SourceDrop: "", Stats: Stats{StatSpellDmg: 85}},
	{Slot: EquipTotem, Name: "Skycall Totem", Phase: 4, Quality: ItemQualityEpic, SourceZone: "Geras", SourceDrop: "20 Badges", Stats: Stats{}, Activate: ActivateSkycall, ActivateCD: -1}, // -1 will trigger an activation only once
}

//  C'thun ring, belt and gloves, Sapph cloak and trinket, KT trinket, 4hm ring

type ItemSet struct {
	Name    string
	Items   map[string]bool
	Bonuses map[int]ItemActivation // maps item count to activations
}

var sets = []ItemSet{
	{
		Name:  "Netherstrike",
		Items: map[string]bool{"Netherstrike Breastplate": true, "Netherstrike Bracers": true, "Netherstrike Belt": true},
		Bonuses: map[int]ItemActivation{3: func(sim *Simulation) Aura {
			sim.Buffs[StatSpellDmg] += 23
			return Aura{ID: MagicIDNetherstrike, Expires: 0}
		}},
	},
	{
		Name:  "The Twin Stars",
		Items: map[string]bool{"Charlotte's Ivy": true, "Lola's Eve": true},
		Bonuses: map[int]ItemActivation{2: func(sim *Simulation) Aura {
			sim.Buffs[StatSpellDmg] += 15
			return Aura{ID: MagicIDNetherstrike, Expires: 0}
		}},
	},
	{
		Name:  "Tidefury",
		Items: map[string]bool{"Tidefury Helm": true, "Tidefury Shoulderguards": true, "Tidefury Chestpiece": true, "Tidefury Kilt": true, "Tidefury Gauntlets": true},
		Bonuses: map[int]ItemActivation{4: func(sim *Simulation) Aura {
			if sim.Options.Buffs.WaterShield {
				sim.Buffs[StatMP5] += 3
			}
			return Aura{ID: MagicIDNetherstrike, Expires: 0}
		}},
	},
	{
		Name:    "Spellstrike",
		Items:   map[string]bool{"Spellstrike Hood": true, "Spellstrike Pants": true},
		Bonuses: map[int]ItemActivation{2: ActivateSpellstrike},
	},
	{
		Name:  "Mana Etched",
		Items: map[string]bool{"Mana-Etched Crown": true, "Mana-Etched Spaulders": true, "Mana-Etched Vestments": true, "Mana-Etched Gloves": true, "Mana-Etched Pantaloons": true},
		Bonuses: map[int]ItemActivation{4: ActivateManaEtched, 2: func(sim *Simulation) Aura {
			sim.Buffs[StatSpellHit] += 35
			return Aura{ID: MagicIDManaEtchedHit, Expires: 0}
		}},
	},
	{
		Name:  "Cyclone Regalia",
		Items: map[string]bool{"Cyclone Faceguard (Tier 4)": true, "Cyclone Shoulderguards (Tier 4)": true, "Cyclone Chestguard (Tier 4)": true, "Cyclone Handguards (Tier 4)": true, "Cyclone Legguards (Tier 4)": true},
		Bonuses: map[int]ItemActivation{4: ActivateCycloneManaReduce, 2: func(sim *Simulation) Aura {
			if !sim.Options.Totems.Cyclone2PC && sim.Options.Totems.WrathOfAir {
				sim.Buffs[StatSpellDmg] += 20 // only activate if we don't already have it from party/
			}
			return Aura{ID: MagicIDCyclone2pc, Expires: 0}
		}},
	},
	{
		Name:  "Windhawk",
		Items: map[string]bool{"Windhawk Hauberk": true, "Windhawk Belt": true, "Windhawk Bracers": true},
		Bonuses: map[int]ItemActivation{3: func(sim *Simulation) Aura {
			if sim.Options.Buffs.WaterShield {
				sim.Buffs[StatMP5] += 8
			}
			return Aura{ID: MagicIDWindhawk, Expires: 0}
		}},
	},
}
