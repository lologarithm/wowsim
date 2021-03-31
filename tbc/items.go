package tbc

import "fmt"

// TODO: is the item separated structure better?
//   or should it just be a giant list?
var items = struct {
	Head     []Item
	Neck     []Item
	Shoulder []Item
	Back     []Item
	Chest    []Item
	Wrist    []Item
	Hands    []Item
	Waist    []Item
	Legs     []Item
	Feet     []Item
	Finger   []Item
	Trinket  []Item
	MainHand []Item
	OffHand  []Item
}{
	Head: []Item{
		{Slot: EquipHead, Name: "Uni-Mind Headdress", SourceZone: "Kara", SourceDrop: "Netherspite",
			Stats: Stats{StatStm: 31, StatInt: 40, StatSpellDmg: 46, StatSpellCrit: 25, StatSpellHit: 19}},
		{Slot: EquipHead, Name: "Wicked Witch's Hat", SourceZone: "Kara", SourceDrop: "Opera",
			Stats: Stats{StatStm: 37, StatInt: 38, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 32, StatSpellHit: 0, StatMP5: 0}},
		{Slot: EquipHead, Name: "Cyclone Faceguard (Tier 4)", SourceZone: "Kara", SourceDrop: "Prince",
			Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 8}},
		{Slot: EquipHead, Name: "Spellstrike Hood", SourceZone: "Crafted", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 16, StatInt: 12, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 16, StatMP5: 0}},
		{Slot: EquipHead, Name: "Cataclysm Headpiece (Tier 5)", SourceZone: "SSC", SourceDrop: "Lady Vashj",
			Stats: Stats{StatStm: 35, StatInt: 28, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 18, StatMP5: 7}},
		{Slot: EquipHead, Name: "Cowl of the Grand Engineer", SourceZone: "TK", SourceDrop: "Void Reaver",
			Stats: Stats{StatStm: 22, StatInt: 27, StatSpellDmg: 53, StatHaste: 0, StatSpellCrit: 35, StatSpellHit: 16, StatMP5: 0}},
		{Slot: EquipHead, Name: "Magnified Moon Specs", SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Leather)",
			Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 41, StatSpellHit: 0, StatMP5: 0}},
		{Slot: EquipHead, Name: "Gadgetstorm Goggles", SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Mail)",
			Stats: Stats{StatStm: 28, StatInt: 0, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 40, StatSpellHit: 12, StatMP5: 0}},
		{Slot: EquipHead, Name: "Destruction Holo-gogs", SourceZone: "Crafted (Patch 2.1)", SourceDrop: "Engineering (Cloth)",
			Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 64, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 0, StatMP5: 0}},
		{Slot: EquipHead, Name: "Skyshatter Headguard (Tier 6)", SourceZone: "Hyjal", SourceDrop: "Archimonde",
			Stats: Stats{StatStm: 42, StatInt: 37, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 36, StatSpellHit: 0, StatMP5: 8}},
		{Slot: EquipHead, Name: "Cowl of the Illidari High Lord", SourceZone: "BT", SourceDrop: "Illidan",
			Stats: Stats{StatStm: 33, StatInt: 31, StatSpellDmg: 64, StatHaste: 0, StatSpellCrit: 47, StatSpellHit: 21, StatMP5: 0}},
	},
	// Meta Gem, Red Gems, Orange Gems, Purple Gems, Socket Bonus
	Neck: []Item{
		{Slot: EquipNeck, Name: "Brooch of Unquenchable Fury", SourceZone: "Kara", SourceDrop: "Moroes",
			Stats: Stats{StatStm: 24, StatInt: 21, StatSpellDmg: 26, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 15, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipNeck, Name: "Manasurge Pendant", SourceZone: "Shattrah", SourceDrop: "Badges",
			Stats: Stats{StatStm: 24, StatInt: 22, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, // , 0, 0, 0, 0, 0},
		{Slot: EquipNeck, Name: "Pendant of the Lost Ages", SourceZone: "SSC", SourceDrop: "Tidewalker",
			Stats: Stats{StatStm: 27, StatInt: 17, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipNeck, Name: "Adornment of Stolen Souls", SourceZone: "Kara", SourceDrop: "Prince",
			Stats: Stats{StatStm: 18, StatInt: 20, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipNeck, Name: "The Sun King's Talisman", SourceZone: "TK", SourceDrop: "Kael Reward",
			Stats: Stats{StatStm: 22, StatInt: 16, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipNeck, Name: "Translucent Spellthread Necklace", SourceZone: "BT", SourceDrop: "RoS",
			Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 15, StatMP5: 0}}, // 0, 0, 0, 0, 0},
	},
	Shoulder: []Item{
		{Slot: EquipShoulder, Name: "Mantle of the Mind Flayer", SourceZone: "Kara", SourceDrop: "Aran",
			Stats: Stats{StatStm: 33, StatInt: 29, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Mantle of the Elven Kings", SourceZone: "TK", SourceDrop: "Trash",
			Stats: Stats{StatStm: 27, StatInt: 18, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Cyclone Shoulderguards (Tier 4)", SourceZone: "Gruul's Lair", SourceDrop: "Maulgar",
			Stats: Stats{StatStm: 28, StatInt: 26, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 12, StatSpellHit: 0, StatMP5: 0}}, //, 0, 18, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Illidari Shoulderpads", SourceZone: "SSC", SourceDrop: "Tidewalker",
			Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 16, StatSpellHit: 0, StatMP5: 0}}, //, 0, 18, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Blood-cursed Shoulderpads", SourceZone: "BT", SourceDrop: "Bloodboil",
			Stats: Stats{StatStm: 25, StatInt: 19, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Cataclysm Shoulderpads (Tier 5)", SourceZone: "TK", SourceDrop: "VoidReaver",
			Stats: Stats{StatStm: 26, StatInt: 19, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 6}}, //, 0, 18, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Mantle of Nimble Thought", SourceZone: "BT", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 37, StatInt: 26, StatSpellDmg: 44, StatHaste: 38, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Skyshatter Mantle (Tier 6)", SourceZone: "BT", SourceDrop: "Mother",
			Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 27, StatSpellHit: 11, StatMP5: 4}}, //, 0, 24, 0, 0, 0},
		{Slot: EquipShoulder, Name: "Hatefury Mantle", SourceZone: "Hyjal", SourceDrop: "Anetheron",
			Stats: Stats{StatStm: 15, StatInt: 18, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}}, //, 0, 24, 0, 0, 0},
	},
	Back: []Item{
		{Slot: EquipBack, Name: "Ruby Drape of the Mysticant", SourceZone: "Kara", SourceDrop: "Prince",
			Stats: Stats{StatStm: 22, StatInt: 21, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Shadow-Cloak of Dalaran", SourceZone: "Kara", SourceDrop: "Moroes",
			Stats: Stats{StatStm: 19, StatInt: 18, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Shawl of Shifting Probabilities", SourceZone: "Shattrah", SourceDrop: "Badges",
			Stats: Stats{StatStm: 18, StatInt: 16, StatSpellDmg: 21, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Royal Cloak of the Sunstriders", SourceZone: "TK", SourceDrop: "Kaelthas",
			Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Brute Cloak of the Ogre-Magi", SourceZone: "Gruul'sLair", SourceDrop: "Maulgar",
			Stats: Stats{StatStm: 18, StatInt: 20, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Ancient Spellcloak of the Highborne", SourceZone: "WorldBoss", SourceDrop: "Kazzak",
			Stats: Stats{StatStm: 0, StatInt: 15, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipBack, Name: "Cloak of the Illidari Council", SourceZone: "BT", SourceDrop: "IllidariCouncil",
			Stats: Stats{StatStm: 24, StatInt: 16, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
	},
	Chest: []Item{
		{Slot: EquipChest, Name: "Cyclone Chestguard (Tier 4)", SourceZone: "GruulsLair", SourceDrop: "Maulgar",
			Stats: Stats{StatStm: 33, StatInt: 32, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 8}}, // 0, 27, 0, 0, 0},
		{Slot: EquipChest, Name: "Netherstrike Breastplate", SourceZone: "Crafted", SourceDrop: "Leatherworking",
			Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 32, StatSpellHit: 0, StatMP5: 8}}, // 0, 27, 0, 0, 0},
		{Slot: EquipChest, Name: "Robe of Hateful Echoes", SourceZone: "SSC", SourceDrop: "Hydross",
			Stats: Stats{StatStm: 34, StatInt: 36, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 0, StatMP5: 0}}, // 0, 27, 0, 0, 0},
		{Slot: EquipChest, Name: "Robe of the Shadow Council", SourceZone: "BT", SourceDrop: "Teron",
			Stats: Stats{StatStm: 37, StatInt: 36, StatSpellDmg: 73, StatHaste: 0, StatSpellCrit: 28, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipChest, Name: "Robes of Rhonin", SourceZone: "Hyjal", SourceDrop: "Archimonde",
			Stats: Stats{StatStm: 55, StatInt: 38, StatSpellDmg: 81, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 27, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipChest, Name: "Cataclysm Chestpiece (Tier 5)", SourceZone: "TK", SourceDrop: "Kaelthas",
			Stats: Stats{StatStm: 37, StatInt: 28, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 10}}, // 0, 27, 0, 0, 0},
		{Slot: EquipChest, Name: "Vestments of the Sea-Witch", SourceZone: "SSC", SourceDrop: "LadyVashj",
			Stats: Stats{StatStm: 28, StatInt: 28, StatSpellDmg: 57, StatHaste: 0, StatSpellCrit: 31, StatSpellHit: 27, StatMP5: 0}}, // 0, 27, 0, 0, 0},
		{Slot: EquipChest, Name: "Chestguard of Relentless Storms", SourceZone: "Hyjal", SourceDrop: "Trash",
			Stats: Stats{StatStm: 36, StatInt: 30, StatSpellDmg: 74, StatHaste: 0, StatSpellCrit: 46, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipChest, Name: "Skyshatter Breastplate (Tier 6)", SourceZone: "BT", SourceDrop: "Illidan",
			Stats: Stats{StatStm: 42, StatInt: 41, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 27, StatSpellHit: 17, StatMP5: 7}}, // 0, 36, 0, 0, 0},
	},
	Wrist: []Item{
		{Slot: EquipWrist, Name: "Bands of Nefarious Deeds", SourceZone: "Kara", SourceDrop: "Maiden",
			Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 32, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWrist, Name: "Elunite Empowered Bracers", SourceZone: "BT", SourceDrop: "RoS",
			Stats: Stats{StatStm: 27, StatInt: 22, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 6}}, //0, 0, 0, 0, 0,
		{Slot: EquipWrist, Name: "Focused Mana Bindings", SourceZone: "BT", SourceDrop: "Akama",
			Stats: Stats{StatStm: 27, StatInt: 20, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWrist, Name: "Netherstrike Bracers", SourceZone: "Crafted", SourceDrop: "Leatherworking",
			Stats: Stats{StatStm: 13, StatInt: 13, StatSpellDmg: 20, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 6}}, //0, 0, 8, 0, 2,
		{Slot: EquipWrist, Name: "Bands of the Coming Storm", SourceZone: "BT", SourceDrop: "Supremus",
			Stats: Stats{StatStm: 28, StatInt: 28, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWrist, Name: "Mindstorm Wristbands", SourceZone: "TK", SourceDrop: "Alar",
			Stats: Stats{StatStm: 13, StatInt: 13, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWrist, Name: "Cuffs of Devastation", SourceZone: "Hyjal", SourceDrop: "Winterchill",
			Stats: Stats{StatStm: 22, StatInt: 20, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 14, StatSpellHit: 0, StatMP5: 0}}, //0, 12, 0, 0, 0,
	},
	Hands: []Item{
		{Slot: EquipHands, Name: "Cyclone Handguards (Tier 4)", SourceZone: "Kara", SourceDrop: "Curator",
			Stats: Stats{StatStm: 26, StatInt: 29, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 19, StatMP5: 6}}, //0, 0, 0, 0, 0,
		{Slot: EquipHands, Name: "Handwraps of Flowing Thought", SourceZone: "Kara", SourceDrop: "Huntsman",
			Stats: Stats{StatStm: 24, StatInt: 22, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 14, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipHands, Name: "Cataclysm Handgrips (Tier 5)", SourceZone: "TK", SourceDrop: "LeotherastheBlind",
			Stats: Stats{StatStm: 25, StatInt: 27, StatSpellDmg: 41, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 19, StatMP5: 7}}, //0, 0, 0, 0, 0,
		{Slot: EquipHands, Name: "Gauntlets of the Sun King", SourceZone: "TK", SourceDrop: "Kaelthas",
			Stats: Stats{StatStm: 28, StatInt: 29, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 28, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipHands, Name: "Anger-Spark Gloves", SourceZone: "World Boss", SourceDrop: "Doomwalker",
			Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 20, StatMP5: 0}}, //0, 18, 0, 0, 2,
		{Slot: EquipHands, Name: "Soul-Eater's Handwraps", SourceZone: "Magtheridon's Lair", SourceDrop: "Magtheridon",
			Stats: Stats{StatStm: 31, StatInt: 24, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipHands, Name: "Skyshatter Guantlets (Tier 6)", SourceZone: "Hyjal", SourceDrop: "Azgalor",
			Stats: Stats{StatStm: 30, StatInt: 31, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 19, StatMP5: 0}}, //0, 9, 0, 0, 2,
	},
	Waist: []Item{
		{Slot: EquipWaist, Name: "Nethershard Girdle", SourceZone: "Kara", SourceDrop: "Moroes",
			Stats: Stats{StatStm: 22, StatInt: 30, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWaist, Name: "General's Mail Girdle", SourceZone: "PvP", SourceDrop: "PvP",
			Stats: Stats{StatStm: 34, StatInt: 23, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWaist, Name: "Malefic Girdle", SourceZone: "Kara", SourceDrop: "Illhoof",
			Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWaist, Name: "Monsoon Belt", SourceZone: "SSC/TK", SourceDrop: "Leatherworking",
			Stats: Stats{StatStm: 23, StatInt: 24, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 21, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Netherstrike Belt", SourceZone: "Crafted", SourceDrop: "Leatherworking",
			Stats: Stats{StatStm: 10, StatInt: 17, StatSpellDmg: 30, StatHaste: 0, StatSpellCrit: 16, StatSpellHit: 0, StatMP5: 9}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Belt of Divine Inspiration", SourceZone: "Gruul's Lair", SourceDrop: "Maulgar",
			Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Cord of Screaming Terrors", SourceZone: "SSC", SourceDrop: "Lurker",
			Stats: Stats{StatStm: 34, StatInt: 15, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 24, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Girdle of Ruination", SourceZone: "Crafted", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 18, StatInt: 13, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Belt of the Crescent Moon", SourceZone: "Hyjal", SourceDrop: "Kazrogal",
			Stats: Stats{StatStm: 25, StatInt: 27, StatSpellDmg: 44, StatHaste: 36, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWaist, Name: "Waistwrap of Infinity", SourceZone: "BT", SourceDrop: "Supremus",
			Stats: Stats{StatStm: 31, StatInt: 22, StatSpellDmg: 56, StatHaste: 32, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
		{Slot: EquipWaist, Name: "Belt of Blasting", SourceZone: "SSC/TK", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 30, StatSpellHit: 23, StatMP5: 0}}, //0, 18, 0, 0, 0,
		{Slot: EquipWaist, Name: "Anetheron's Noose", SourceZone: "Hyjal", SourceDrop: "Anetheron",
			Stats: Stats{StatStm: 22, StatInt: 23, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}}, //0, 24, 0, 0, 0,
		{Slot: EquipWaist, Name: "Flashfire Girdle", SourceZone: "BT", SourceDrop: "Akama",
			Stats: Stats{StatStm: 27, StatInt: 26, StatSpellDmg: 44, StatHaste: 37, StatSpellCrit: 18, StatSpellHit: 0, StatMP5: 0}}, //0, 0, 0, 0, 0,
	},
	Legs: []Item{
		{Slot: EquipLegs, Name: "Cyclone Legguards (Tier 4)", SourceZone: "Gruul's Lair", SourceDrop: "Gruul",
			Stats: Stats{StatStm: 40, StatInt: 40, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 20, StatMP5: 8}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipLegs, Name: "Trial-Fire Trousers", SourceZone: "Kara", SourceDrop: "Opera",
			Stats: Stats{StatStm: 42, StatInt: 40, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 25, 0, 5},
		{Slot: EquipLegs, Name: "Trousers of the Astromancer", SourceZone: "TK", SourceDrop: "Solarian",
			Stats: Stats{StatStm: 33, StatInt: 36, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 27, 0, 0, 0},
		{Slot: EquipLegs, Name: "Cataclysm Leggings (Tier 5)", SourceZone: "TK", SourceDrop: "Karathress",
			Stats: Stats{StatStm: 48, StatInt: 46, StatSpellDmg: 54, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 14, StatMP5: 0}}, //, 0, 0, 8, 0, 2},
		{Slot: EquipLegs, Name: "Spellstrike Pants", SourceZone: "Crafted", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 12, StatInt: 8, StatSpellDmg: 46, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 22, StatMP5: 0}}, //, 0, 27, 0, 0, 0},
		{Slot: EquipLegs, Name: "Leggings of Devastation", SourceZone: "BT", SourceDrop: "Mother",
			Stats: Stats{StatStm: 40, StatInt: 42, StatSpellDmg: 60, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 26, StatMP5: 0}}, //, 0, 36, 0, 0, 0},
		{Slot: EquipLegs, Name: "Skyshatter Pants (Tier 6)", SourceZone: "BT", SourceDrop: "IllidariCouncil",
			Stats: Stats{StatStm: 40, StatInt: 42, StatSpellDmg: 62, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 20, StatMP5: 11}}, //, 0, 0, 11, 0, 2},
		{Slot: EquipLegs, Name: "Leggings of the Seventh Circle", SourceZone: "World Boss", SourceDrop: "Kazzak",
			Stats: Stats{StatStm: 0, StatInt: 22, StatSpellDmg: 50, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}}, //, 0, 9, 16, 0, 5},
		{Slot: EquipLegs, Name: "Leggings of Channeled Elements", SourceZone: "Hyjal", SourceDrop: "Kazrogal",
			Stats: Stats{StatStm: 25, StatInt: 28, StatSpellDmg: 59, StatHaste: 0, StatSpellCrit: 34, StatSpellHit: 18, StatMP5: 0}}, //, 0, 36, 0, 0, 0},
	},
	Feet: []Item{
		{Slot: EquipFeet, Name: "Boots of the Infernal Coven", SourceZone: "Kara", SourceDrop: "Aran",
			Stats: Stats{StatStm: 27, StatInt: 27, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Ruby Slippers", SourceZone: "Kara", SourceDrop: "Opera",
			Stats: Stats{StatStm: 33, StatInt: 29, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 16, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Windshear Boots", SourceZone: "Gruul's Lair", SourceDrop: "Gruul",
			Stats: Stats{StatStm: 37, StatInt: 32, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Blue Suede Shoes", SourceZone: "Hyjal", SourceDrop: "Kazrogal",
			Stats: Stats{StatStm: 37, StatInt: 32, StatSpellDmg: 56, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 18, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Boots of Blasting", SourceZone: "SSC/TK", SourceDrop: "Tailoring",
			Stats: Stats{StatStm: 25, StatInt: 25, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 25, StatSpellHit: 18, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Boots of Foretelling", SourceZone: "Kara", SourceDrop: "Maiden",
			Stats: Stats{StatStm: 27, StatInt: 23, StatSpellDmg: 26, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}}, //, 0, 18, 0, 0, 0},
		{Slot: EquipFeet, Name: "Hurricane Boots", SourceZone: "SSC/TK", SourceDrop: "Leatherworking",
			Stats: Stats{StatStm: 25, StatInt: 26, StatSpellDmg: 39, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 6}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Velvet Boots of the Guardian", SourceZone: "SSC", SourceDrop: "Lurker",
			Stats: Stats{StatStm: 21, StatInt: 21, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Boots of Oceanic Fury", SourceZone: "BT", SourceDrop: "Najentus",
			Stats: Stats{StatStm: 28, StatInt: 36, StatSpellDmg: 55, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFeet, Name: "Naturewarden's Treads", SourceZone: "BT", SourceDrop: "RoS",
			Stats: Stats{StatStm: 39, StatInt: 18, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 26, StatSpellHit: 0, StatMP5: 7}}, //, 0, 24, 0, 0, 0},
		{Slot: EquipFeet, Name: "Slippers of the Seacaller", SourceZone: "BT", SourceDrop: "Najentus",
			Stats: Stats{StatStm: 25, StatInt: 18, StatSpellDmg: 44, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 0, StatMP5: 0}}, //, 0, 24, 0, 0, 0},
	},
	Finger: []Item{
		{Slot: EquipFinger, Name: "Band of Crimson Fury", SourceZone: "Magtheridon's Lair", SourceDrop: "MagtheridonQuest",
			Stats: Stats{StatStm: 22, StatInt: 22, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 16, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Spectral Band of Innervation", SourceZone: "Kara", SourceDrop: "Huntsman",
			Stats: Stats{StatStm: 22, StatInt: 24, StatSpellDmg: 29, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Band of Alar", SourceZone: "TK", SourceDrop: "Alar",
			Stats: Stats{StatStm: 24, StatInt: 23, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Cryptic Dreams", SourceZone: "Shattrah", SourceDrop: "Badges",
			Stats: Stats{StatStm: 16, StatInt: 17, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Violet Signet of the Archmage", SourceZone: "Kara", SourceDrop: "Exalted",
			Stats: Stats{StatStm: 24, StatInt: 23, StatSpellDmg: 29, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Recurrence", SourceZone: "Kara", SourceDrop: "Chess",
			Stats: Stats{StatStm: 15, StatInt: 15, StatSpellDmg: 32, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Band of the Eternal Sage", SourceZone: "Hyjal", SourceDrop: "Exalted",
			Stats: Stats{StatStm: 28, StatInt: 25, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 24, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Endless Coils", SourceZone: "SSC", SourceDrop: "LadyVashj",
			Stats: Stats{StatStm: 31, StatInt: 0, StatSpellDmg: 37, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Unrelenting Storms", SourceZone: "Kara", SourceDrop: "Trash",
			Stats: Stats{StatStm: 0, StatInt: 15, StatSpellDmg: 43, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Captured Storms", SourceZone: "BT", SourceDrop: "Najentus",
			Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 29, StatSpellHit: 19, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipFinger, Name: "Ring of Ancient Knowledge", SourceZone: "BT", SourceDrop: "Trash",
			Stats: Stats{StatStm: 30, StatInt: 20, StatSpellDmg: 39, StatHaste: 31, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
	},
	MainHand: []Item{
		{Slot: EquipWeapon, Name: "Gavel of Unearthed Secrets", SourceZone: "Shattrah", SourceDrop: "LowerCityExalted",
			Stats: Stats{StatStm: 24, StatInt: 16, StatSpellDmg: 159, StatHaste: 0, StatSpellCrit: 15, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Eternium Runed Blade", SourceZone: "Crafted", SourceDrop: "Blacksmithing",
			Stats: Stats{StatStm: 0, StatInt: 19, StatSpellDmg: 168, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Gladiator's Gavel / Gladiator's Spellblade", SourceZone: "PvP", SourceDrop: "PvP",
			Stats: Stats{StatStm: 28, StatInt: 18, StatSpellDmg: 199, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Nathrezim Mindblade", SourceZone: "Kara", SourceDrop: "Prince",
			Stats: Stats{StatStm: 18, StatInt: 18, StatSpellDmg: 203, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Talon of the Tempest", SourceZone: "World Boss", SourceDrop: "Doomwalker",
			Stats: Stats{StatStm: 0, StatInt: 10, StatSpellDmg: 194, StatHaste: 0, StatSpellCrit: 19, StatSpellHit: 9, StatMP5: 0}}, //, 0, 18, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Hammer of Judgement", SourceZone: "Hyjal", SourceDrop: "Trash",
			Stats: Stats{StatStm: 33, StatInt: 22, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 22, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "The Maelstrom's Fury", SourceZone: "BT", SourceDrop: "Najentus",
			Stats: Stats{StatStm: 33, StatInt: 21, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 22, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Staff of Infinite Mysteries", SourceZone: "Kara", SourceDrop: "Curator",
			Stats: Stats{StatStm: 61, StatInt: 51, StatSpellDmg: 185, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 23, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "The Nexus Key", SourceZone: "TK", SourceDrop: "Kaelthas",
			Stats: Stats{StatStm: 76, StatInt: 52, StatSpellDmg: 236, StatHaste: 0, StatSpellCrit: 51, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipWeapon, Name: "Zhar'doom, Greatstaff of the Devourer", SourceZone: "BT", SourceDrop: "Illidan",
			Stats: Stats{StatStm: 70, StatInt: 47, StatSpellDmg: 259, StatHaste: 55, StatSpellCrit: 36, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
	},
	OffHand: []Item{
		{Slot: EquipOffhand, Name: "Jewel of Infinite Possibilities", SourceZone: "Kara", SourceDrop: "Netherspite",
			Stats: Stats{StatStm: 19, StatInt: 18, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 21, StatMP5: 0}}, // , 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Dragonheart Flameshield", SourceZone: "Kara", SourceDrop: "Nightbane",
			Stats: Stats{StatStm: 19, StatInt: 33, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 7}}, // , 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Illidari Runeshield", SourceZone: "BT", SourceDrop: "Trash",
			Stats: Stats{StatStm: 45, StatInt: 39, StatSpellDmg: 34, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Karaborian Talisman", SourceZone: "Magtheridon's Lair", SourceDrop: "Magtheridon",
			Stats: Stats{StatStm: 23, StatInt: 23, StatSpellDmg: 35, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Mazthoril Honor Shield", SourceZone: "Shattrah", SourceDrop: "Badges",
			Stats: Stats{StatStm: 16, StatInt: 29, StatSpellDmg: 23, StatHaste: 0, StatSpellCrit: 21, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Talisman of Nightbane", SourceZone: "Kara", SourceDrop: "Nightbane",
			Stats: Stats{StatStm: 19, StatInt: 19, StatSpellDmg: 28, StatHaste: 0, StatSpellCrit: 17, StatSpellHit: 0, StatMP5: 0}}, // 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Blind-Seers Icon", SourceZone: "BT", SourceDrop: "Akama",
			Stats: Stats{StatStm: 25, StatInt: 16, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 24, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Khadgar's Knapsack", SourceZone: "Shattrah", SourceDrop: "Badges",
			Stats: Stats{StatStm: 0, StatInt: 0, StatSpellDmg: 49, StatHaste: 0, StatSpellCrit: 0, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "FathomStone", SourceZone: "SSC", SourceDrop: "Lurker",
			Stats: Stats{StatStm: 16, StatInt: 12, StatSpellDmg: 36, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Antonidas's Aegis of Rapt Concentration", SourceZone: "Hyjal", SourceDrop: "Archimonde",
			Stats: Stats{StatStm: 28, StatInt: 32, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 20, StatSpellHit: 0, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
		{Slot: EquipOffhand, Name: "Chronicle of Dark Secrets", SourceZone: "Hyjal", SourceDrop: "Winterchill",
			Stats: Stats{StatStm: 16, StatInt: 12, StatSpellDmg: 42, StatHaste: 0, StatSpellCrit: 23, StatSpellHit: 17, StatMP5: 0}}, //, 0, 0, 0, 0, 0},
	},
}

var ItemLookup = map[string]*Item{}

func IL(name string) *Item {
	return ItemLookup[name]
}

func init() {
	for _, v := range items.Head {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Neck {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Shoulder {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Back {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Chest {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Wrist {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Hands {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Waist {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Legs {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Feet {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Finger {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.Trinket {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.MainHand {
		cv := v
		ItemLookup[cv.Name] = &cv
	}
	for _, v := range items.OffHand {
		cv := v
		ItemLookup[cv.Name] = &cv
	}

	for _, v := range moreItems {
		if it, ok := ItemLookup[v.Name]; ok {
			// log.Printf("Found dup item: %s", v.Name)
			statsMatch := it.Slot == v.Slot
			for i, v := range v.Stats {
				if it.Stats[i] != v {
					statsMatch = false
				}
			}
			if !statsMatch {
				// log.Printf("Mismatched slot/stats: \n\tMoreItem: \n%#v\n\t FirstSet: \n%#v", it, v)
			}
		} else {
			cv := v
			ItemLookup[cv.Name] = &cv
		}
	}
}

// Darkmoon Card: Crusade	Blessings Deck	78
// Scryer's Bloodgem	The Scryers - Revered	74.6 / 21						32
// Quagmirran's Eye	H SP - Quagmirran	68.6			37
// Arcanist's Stone	H OHF - Epoch Hunter	66.7 / 24.2						25
// Icon of the Silver Crescent	41 Badge of Justice - G'eras	64.8			43
// Shiffar's Nexus-Horn	Arc - Harbinger Skyriss	45.4				30
// Xi'ri's Gift	The Sha'tar - Revered	40				32
// Vengeance of the Illidari	Cruel's Intentions/Overlord - HFP Quest	29.3				26
// Figurine - Living Ruby Serpent	Jewelcarfting BoP		33	23

type Item struct {
	Slot       int
	Name       string
	SourceZone string
	SourceDrop string
	Stats      Stats // Stats applied to wearer
	Aura       Aura  `json:"-"` // Aura item applies when worn
}

type Equipment []Item

func NewEquipmentSet(names ...string) Equipment {
	e := Equipment{EquipTotem: Item{}}
	for _, v := range names {
		item, ok := ItemLookup[v]
		if !ok {
			fmt.Printf("Unable to find item: '%s'", v)
			return e
		}
		if item.Slot == EquipFinger {
			if e[EquipFinger1].Name == "" {
				e[EquipFinger1] = *item
			} else {
				e[EquipFinger2] = *item
			}
		} else if item.Slot == EquipTrinket {
			if e[EquipTrinket1].Name == "" {
				e[EquipTrinket1] = *item
			} else {
				e[EquipTrinket2] = *item
			}
		} else {
			e[item.Slot] = *item
		}
	}
	return e
}

const (
	EquipUnknown int = iota
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

func (e Equipment) Stats() Stats {
	s := Stats{StatLen: 0}
	for _, item := range e {
		for k, v := range item.Stats {
			s[k] += v
		}
	}
	return s
}

// "Slot","Name","Location","Boss","Sta","Int","Sp","Haste","Crit","Hit","MP5","Meta Gem","Red Gems","Orange Gems","Purple Gems","Socket Bonus"

// "Color","Gem","Value",,,,,,,,,,,,,
// "Meta","Chaotic Skyfire Diamond (ZA Patch)","47.8",,,,,,,,,,,,,
// "Meta","Mystical Skyfire Diamond","10.1",,,,,,,,,,,,,
// "Red","Runed Living Ruby","9",,,,,,,,,,,,,
// "Orange","Potent Noble Topaz","8.2",,,,,,,,,,,,,
// "Purple","Glowing Nightseye","5",,,,,,,,,,,,,
// "Red (Epic)","Runed Crimson Spinel","12",,,,,,,,,,,,,
// "Orange (Epic)","Potent Pyrestone","10.5",,,,,,,,,,,,,
// "Purple (Epic)","Glowing Shadowsong","6",,,,,,,,,,,,,
// "Meta","Chaotic Skyfire Diamond","47.8",,,,,,,,,,,,,
// "Meta","Mystical Skyfire Diamond","10.1",,,,,,,,,,,,,
// "Red","Runed Living Ruby","9",,,,,,,,,,,,,

var moreItems = []Item{
	{Slot: 1, Name: "Gadgetstorm Goggles", SourceZone: "Engineering BoP", SourceDrop: "", Stats: Stats{0, 28, 40, 12, 55, 0, 0}},
	{Slot: 1, Name: "Gladiator's Mail Helm", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{15, 54, 18, 0, 37, 0, 0}},
	{Slot: 1, Name: "Spellstrike Hood", SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{12, 16, 24, 16, 46, 0, 0}},
	{Slot: 1, Name: "Incanter's Cowl", SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{27, 15, 19, 0, 29, 0, 0}},
	{Slot: 1, Name: "Lightning Crown", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{0, 0, 43, 0, 66, 0, 0}},
	{Slot: 1, Name: "Hood of Oblivion", SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{32, 27, 0, 0, 40, 0, 0}},
	{Slot: 1, Name: "Exorcist's Mail Helm", SourceZone: "18 Spirit Shards", SourceDrop: "", Stats: Stats{16, 30, 24, 0, 29, 0, 0}},
	{Slot: 1, Name: "Tidefury Helm", SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{26, 32, 0, 0, 32, 0, 6}},
	{Slot: 1, Name: "Windscale Hood", SourceZone: "Leatherworking BoE", SourceDrop: "", Stats: Stats{18, 16, 37, 0, 44, 0, 10}},
	{Slot: 1, Name: "Shamanistic Helmet of Second Sight", SourceZone: "Teron Gorfiend, I am... - SMV Quest", SourceDrop: "", Stats: Stats{15, 12, 24, 0, 35, 0, 4}},
	{Slot: 1, Name: "Mana-Etched Crown", SourceZone: "BM - Aeonus", SourceDrop: "", Stats: Stats{20, 27, 0, 0, 34, 0, 0}},
	{Slot: 1, Name: "Mag'hari Ritualist's Horns", SourceZone: "Hero of the Mag'har - Nagrand quest (Horde)", SourceDrop: "", Stats: Stats{16, 18, 15, 12, 50, 0, 0}},
	{Slot: 1, Name: "Mage-Collar of the Firestorm", SourceZone: "H BF - The Maker", SourceDrop: "", Stats: Stats{33, 32, 23, 0, 39, 0, 0}},
	{Slot: 1, Name: "Circlet of the Starcaller", SourceZone: "Dimensius the All-Devouring - NS Quest", SourceDrop: "", Stats: Stats{18, 27, 18, 0, 47, 0, 0}},
	{Slot: 1, Name: "Mask of Inner Fire", SourceZone: "BM - Chrono Lord Deja", SourceDrop: "", Stats: Stats{33, 30, 22, 0, 37, 0, 0}},
	{Slot: 1, Name: "Mooncrest Headdress", SourceZone: "Blast the Infernals! - SMV Quest", SourceDrop: "", Stats: Stats{16, 0, 21, 0, 44, 0, 0}},
	{Slot: 2, Name: "Pendant of Dominance", SourceZone: "15,300 Honor & 10 EotS Marks", SourceDrop: "", Stats: Stats{12, 31, 16, 0, 26, 0, 0}},
	{Slot: 2, Name: "Brooch of Heightened Potential", SourceZone: "SLabs - Blackheart the Inciter", SourceDrop: "", Stats: Stats{12, 15, 14, 9, 22, 0, 0}},
	{Slot: 2, Name: "Torc of the Sethekk Prophet", SourceZone: "Brother Against Brother - Auchindoun ", SourceDrop: "", Stats: Stats{18, 0, 21, 0, 19, 0, 0}},
	{Slot: 2, Name: "Natasha's Ember Necklace", SourceZone: "The Hound-Master - BEM Quest", SourceDrop: "", Stats: Stats{15, 0, 10, 0, 29, 0, 0}},
	{Slot: 2, Name: "Warp Engineer's Prismatic Chain", SourceZone: "Mech - Mechano Lord Capacitus", SourceDrop: "", Stats: Stats{18, 17, 16, 0, 19, 0, 0}},
	{Slot: 2, Name: "Eye of the Night", SourceZone: "Jewelcrafting BoE", SourceDrop: "", Stats: Stats{0, 0, 26, 16, 0, 0, 0}},
	{Slot: 2, Name: "Hydra-fang Necklace", SourceZone: "H UB - Ghaz'an", SourceDrop: "", Stats: Stats{16, 17, 0, 16, 19, 0, 0}},
	{Slot: 2, Name: "Manasurge Pendant", SourceZone: "25 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{22, 24, 0, 0, 28, 0, 0}},
	{Slot: 2, Name: "Natasha's Arcane Filament", SourceZone: "The Hound-Master - BEM Quest", SourceDrop: "", Stats: Stats{10, 22, 0, 0, 29, 0, 0}},
	{Slot: 2, Name: "Omor's Unyielding Will", SourceZone: "H Ramps - Omar the Unscarred", SourceDrop: "", Stats: Stats{19, 19, 0, 0, 25, 0, 0}},
	{Slot: 2, Name: "Charlotte's Ivy", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{19, 18, 0, 0, 23, 0, 0}},
	{Slot: 3, Name: "Gladiator's Mail Spaulders", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{17, 33, 20, 0, 22, 0, 6}},
	{Slot: 3, Name: "Pauldrons of Wild Magic", SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{28, 21, 23, 0, 33, 0, 0}},
	{Slot: 3, Name: "Mana-Etched Spaulders", SourceZone: "H UB - Quagmirran", SourceDrop: "", Stats: Stats{17, 25, 16, 0, 20, 0, 0}},
	{Slot: 3, Name: "Spaulders of the Torn-heart", SourceZone: "The Cipher of Damnation - SMV Quest", SourceDrop: "", Stats: Stats{7, 10, 18, 0, 40, 0, 0}},
	{Slot: 3, Name: "Elekk Hide Spaulders", SourceZone: "The Fallen Exarch - Terokkar Forest Quest", SourceDrop: "", Stats: Stats{12, 0, 28, 0, 25, 0, 0}},
	{Slot: 3, Name: "Spaulders of Oblivion", SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{17, 25, 0, 0, 29, 0, 0}},
	{Slot: 3, Name: "Tidefury Shoulderguards", SourceZone: "SH - O'mrogg", SourceDrop: "", Stats: Stats{23, 18, 0, 0, 19, 0, 6}},
	{Slot: 3, Name: "Mantle of Three Terrors", SourceZone: "BM - Chrono Lord Deja", SourceDrop: "", Stats: Stats{25, 29, 0, 12, 29, 0, 0}},
	{Slot: 4, Name: "Shawl of Shifting Probabilities", SourceZone: "25 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{16, 18, 22, 0, 21, 0, 0}},
	{Slot: 4, Name: "Ogre Slayer's Cover", SourceZone: "Cho'war the Pillager - Nagrand Quest", SourceDrop: "", Stats: Stats{18, 0, 16, 0, 20, 0, 0}},
	{Slot: 4, Name: "Baba's Cloak of Arcanistry", SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{15, 15, 14, 0, 22, 0, 0}},
	{Slot: 4, Name: "Cloak of Woven Energy", SourceZone: "Hitting the Motherlode - Netherstorm Quest", SourceDrop: "", Stats: Stats{13, 6, 6, 0, 29, 0, 0}},
	{Slot: 4, Name: "Sethekk Oracle Cloak", SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{18, 18, 0, 12, 22, 0, 0}},
	{Slot: 4, Name: "Terokk's Wisdom", SourceZone: "Terokk - Skettis Summoned Boss", SourceDrop: "", Stats: Stats{16, 18, 0, 0, 33, 0, 0}},
	{Slot: 4, Name: "Cloak of the Black Void", SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{11, 0, 0, 0, 35, 0, 0}},
	{Slot: 4, Name: "Cloak of Entropy", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{11, 0, 0, 10, 25, 0, 0}},
	{Slot: 4, Name: "Sergeant's Heavy Cape", SourceZone: "9,435 Honor & 20 AB Marks", SourceDrop: "", Stats: Stats{12, 33, 0, 0, 26, 0, 0}},
	{Slot: 5, Name: "Netherstrike Breastplate", SourceZone: "Leatherworking BoP - Req. Dragonscale LW", SourceDrop: "", Stats: Stats{23, 34, 32, 0, 37, 0, 8}},
	{Slot: 5, Name: "Gladiator's Mail Armor", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{23, 42, 23, 0, 32, 0, 7}},
	{Slot: 5, Name: "Will of Edward the Odd", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{30, 0, 30, 0, 53, 0, 0}},
	{Slot: 5, Name: "Anchorite's Robe", SourceZone: "The Aldor - Honored", SourceDrop: "", Stats: Stats{38, 16, 0, 0, 29, 0, 18}},
	{Slot: 5, Name: "Tidefury Chestpiece", SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{22, 28, 0, 10, 36, 0, 4}},
	{Slot: 5, Name: "Auchenai Anchorite's Robe", SourceZone: "Everything Will Be Alright - AC Quest", SourceDrop: "", Stats: Stats{24, 0, 0, 23, 28, 0, 0}},
	{Slot: 5, Name: "Mana-Etched Vestments", SourceZone: "OHF - Epoch Hunter", SourceDrop: "", Stats: Stats{25, 25, 17, 0, 29, 0, 0}},
	{Slot: 5, Name: "Robe of the Crimson Order", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{23, 0, 0, 30, 50, 0, 0}},
	{Slot: 5, Name: "Warp Infused Drape", SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{28, 27, 0, 12, 30, 0, 0}},
	{Slot: 5, Name: "Robe of Oblivion", SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{20, 30, 0, 0, 40, 0, 0}},
	{Slot: 5, Name: "Incanter's Robe", SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{22, 24, 8, 0, 29, 0, 0}},
	{Slot: 5, Name: "Robe of the Great Dark Beyond", SourceZone: "MT - Tavarok", SourceDrop: "", Stats: Stats{30, 25, 23, 0, 39, 0, 0}},
	{Slot: 5, Name: "Worldfire Chestguard", SourceZone: "Arc - Dalliah the Doomsayer", SourceDrop: "", Stats: Stats{32, 33, 22, 0, 40, 0, 0}},
	{Slot: 6, Name: "Netherstrike Bracers", SourceZone: "Leatherworking BoP - Req. Dragonscale LW", SourceDrop: "", Stats: Stats{13, 13, 17, 0, 20, 0, 6}},
	{Slot: 6, Name: "General's Mail Bracers", SourceZone: "7,548 Honor & 20 WSG Marks", SourceDrop: "", Stats: Stats{12, 22, 14, 0, 20, 0, 0}},
	{Slot: 6, Name: "World's End Bracers", SourceZone: "H BF - Keli'dan the Breaker", SourceDrop: "", Stats: Stats{19, 18, 17, 0, 22, 0, 0}},
	{Slot: 6, Name: "Bracers of Havok", SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{12, 0, 0, 0, 30, 0, 0}},
	{Slot: 6, Name: "Crimson Bracers of Gloom", SourceZone: "H Ramps - Omor the Unscarred", SourceDrop: "", Stats: Stats{18, 18, 0, 12, 22, 0, 0}},
	{Slot: 6, Name: "Bands of Negation", SourceZone: "H MT - Nexus- Prince Shaffar", SourceDrop: "", Stats: Stats{22, 25, 0, 0, 29, 0, 0}},
	{Slot: 6, Name: "Arcanium Signet Bands", SourceZone: "H UB - Hungarfen", SourceDrop: "", Stats: Stats{15, 14, 0, 0, 30, 0, 0}},
	{Slot: 6, Name: "Wave-Fury Vambraces", SourceZone: "H SV - Warlod Kalithresh", SourceDrop: "", Stats: Stats{18, 19, 0, 0, 22, 0, 5}},
	{Slot: 6, Name: "Mana Infused Wristguards", SourceZone: "A Fate Worse Than Death - Netherstorm Quest", SourceDrop: "", Stats: Stats{8, 12, 0, 0, 25, 0, 0}},
	{Slot: 7, Name: "Mana-Etched Gloves", SourceZone: "H Ramps - Omor the Unscarred", SourceDrop: "", Stats: Stats{17, 25, 16, 0, 20, 0, 0}},
	{Slot: 7, Name: "Earth Mantle Handwraps", SourceZone: "SV - Mekgineer Steamrigger", SourceDrop: "", Stats: Stats{18, 21, 16, 0, 19, 0, 0}},
	{Slot: 7, Name: "Gloves of Pandemonium", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{15, 0, 22, 10, 25, 0, 0}},
	{Slot: 7, Name: "Gladiator's Mail Gauntlets", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{18, 36, 21, 0, 32, 0, 0}},
	{Slot: 7, Name: "Thundercaller's Gauntlets", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{16, 16, 18, 0, 35, 0, 0}},
	{Slot: 7, Name: "Gloves of the High Magus", SourceZone: "News of Victory - SMV Quest", SourceDrop: "", Stats: Stats{18, 13, 22, 0, 26, 0, 0}},
	{Slot: 7, Name: "Tempest's Touch", SourceZone: "Return to Andormu - CoT Quest", SourceDrop: "", Stats: Stats{20, 10, 0, 0, 27, 0, 0}},
	{Slot: 7, Name: "Gloves of the Deadwatcher", SourceZone: "H AC - Shirrak the Dead Watcher", SourceDrop: "", Stats: Stats{24, 24, 0, 18, 29, 0, 0}},
	{Slot: 7, Name: "Incanter's Gloves", SourceZone: "SV - Thespia", SourceDrop: "", Stats: Stats{24, 21, 14, 0, 29, 0, 0}},
	{Slot: 7, Name: "Starlight Gauntlets", SourceZone: "N UB - Hungarfen", SourceDrop: "", Stats: Stats{21, 10, 0, 0, 25, 0, 0}},
	{Slot: 7, Name: "Gloves of Oblivion", SourceZone: "SH - Kargath", SourceDrop: "", Stats: Stats{21, 33, 0, 20, 26, 0, 0}},
	{Slot: 7, Name: "Harmony's Touch", SourceZone: "Building a Perimeter - Netherstorm Quest", SourceDrop: "", Stats: Stats{0, 18, 16, 0, 33, 0, 0}},
	{Slot: 8, Name: "Girdle of Ruination", SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{13, 18, 20, 0, 39, 0, 0}},
	{Slot: 8, Name: "Girdle of Living Flame", SourceZone: "H UB - Hungarfen", SourceDrop: "", Stats: Stats{17, 15, 0, 16, 29, 0, 0}},
	{Slot: 8, Name: "Wave-Song Girdle", SourceZone: "H AC - Exarch Maladaar", SourceDrop: "", Stats: Stats{25, 25, 23, 0, 32, 0, 0}},
	{Slot: 8, Name: "A'dal's Gift", SourceZone: "How to Break Into the Arcatraz - Quest", SourceDrop: "", Stats: Stats{25, 0, 21, 0, 34, 0, 0}},
	{Slot: 8, Name: "Netherstrike Belt", SourceZone: "Leatherworking BoP - Req. Dragonscale LW", SourceDrop: "", Stats: Stats{17, 10, 16, 0, 30, 0, 9}},
	{Slot: 8, Name: "General's Mail Girdle", SourceZone: "14,280 Honor & 40 AB Marks", SourceDrop: "", Stats: Stats{23, 34, 24, 0, 28, 0, 0}},
	{Slot: 8, Name: "Sash of Arcane Visions", SourceZone: "H AC - Exarch Maladaar", SourceDrop: "", Stats: Stats{23, 18, 22, 0, 28, 0, 0}},
	{Slot: 8, Name: "Belt of Depravity", SourceZone: "H Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{27, 31, 0, 17, 34, 0, 0}},
	{Slot: 8, Name: "Moonrage Girdle", SourceZone: "SV - Hydromancer Thespia", SourceDrop: "", Stats: Stats{22, 0, 20, 0, 25, 0, 0}},
	{Slot: 8, Name: "Sash of Serpentra", SourceZone: "SV - Warlord Kalithresh", SourceDrop: "", Stats: Stats{21, 31, 0, 17, 25, 0, 0}},
	{Slot: 8, Name: "Blackwhelp Belt", SourceZone: "Whelps of the Wyrmcult - BEM Quest", SourceDrop: "", Stats: Stats{11, 0, 10, 0, 32, 0, 0}},
	{Slot: 9, Name: "Spellstrike Pants", SourceZone: "Tailoring BoE", SourceDrop: "", Stats: Stats{8, 12, 26, 22, 46, 0, 0}},
	{Slot: 9, Name: "Stormsong Kilt", SourceZone: "H UB - The Black Stalker", SourceDrop: "", Stats: Stats{30, 25, 26, 0, 35, 0, 0}},
	{Slot: 9, Name: "Tempest Leggings", SourceZone: "The Mag'har - Revered (Horde)", SourceDrop: "", Stats: Stats{11, 0, 18, 0, 44, 0, 0}},
	{Slot: 9, Name: "Kurenai Kilt", SourceZone: "Kurenai - Revered (Ally)", SourceDrop: "", Stats: Stats{11, 0, 18, 0, 44, 0, 0}},
	{Slot: 9, Name: "Breeches of the Occultist", SourceZone: "H BM - Aeonus", SourceDrop: "", Stats: Stats{22, 37, 23, 0, 26, 0, 0}},
	{Slot: 9, Name: "Pantaloons of Flaming Wrath", SourceZone: "H SH - Blood Guard Porung", SourceDrop: "", Stats: Stats{28, 0, 42, 0, 33, 0, 0}},
	{Slot: 9, Name: "Moonchild Leggings", SourceZone: "H BF - Broggok", SourceDrop: "", Stats: Stats{20, 26, 21, 0, 23, 0, 0}},
	{Slot: 9, Name: "Haramad's Leggings of the Third Coin", SourceZone: "Undercutting the Competition - MT Quest", SourceDrop: "", Stats: Stats{29, 0, 16, 0, 27, 0, 0}},
	{Slot: 9, Name: "Gladiator's Mail Leggins", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{25, 54, 22, 0, 42, 0, 6}},
	{Slot: 9, Name: "Kirin Tor Master's Trousers", SourceZone: "H SLabs - Murmur", SourceDrop: "", Stats: Stats{29, 27, 0, 0, 36, 0, 0}},
	{Slot: 9, Name: "Khadgar's Kilt of Abjuration", SourceZone: "BM - Temporus", SourceDrop: "", Stats: Stats{22, 20, 0, 0, 36, 0, 0}},
	{Slot: 9, Name: "Incanter's Trousers", SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{30, 25, 18, 0, 42, 0, 0}},
	{Slot: 9, Name: "Mana-Etched Pantaloons", SourceZone: "H UB - The Black Stalker", SourceDrop: "", Stats: Stats{32, 34, 21, 0, 33, 0, 0}},
	{Slot: 9, Name: "Tidefury Kilt", SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{31, 39, 19, 0, 35, 0, 0}},
	{Slot: 9, Name: "Molten Earth Kilt", SourceZone: "Mech - Pathaleon the Calculator", SourceDrop: "", Stats: Stats{32, 24, 0, 0, 40, 0, 10}},
	{Slot: 9, Name: "Trousers of Oblivion", SourceZone: "SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{33, 42, 0, 12, 39, 0, 0}},
	{Slot: 9, Name: "Leggings of the Third Coin", SourceZone: "Levixus the Soul Caller - Auchindoun Quest", SourceDrop: "", Stats: Stats{26, 34, 12, 0, 32, 0, 4}},
	{Slot: 10, Name: "Sigil-Laced Boots", SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{18, 24, 17, 0, 20, 0, 0}},
	{Slot: 10, Name: "General's Mail Sabatons", SourceZone: "11,424 Honor & 40 EotS Marks", SourceDrop: "", Stats: Stats{23, 34, 24, 0, 28, 0, 0}},
	{Slot: 10, Name: "Moonstrider Boots", SourceZone: "SH - Darkweaver Syth", SourceDrop: "", Stats: Stats{22, 21, 20, 0, 25, 0, 6}},
	{Slot: 10, Name: "Shattarath Jumpers", SourceZone: "Into the Heart of the Labyrinth - Auch. Quest", SourceDrop: "", Stats: Stats{17, 25, 0, 0, 29, 0, 0}},
	{Slot: 10, Name: "Wave-Crest Striders", SourceZone: "H BF - Keli'dan the Breaker", SourceDrop: "", Stats: Stats{26, 28, 0, 0, 33, 0, 8}},
	{Slot: 10, Name: "Extravagant Boots of Malice", SourceZone: "H MT - Tavarok", SourceDrop: "", Stats: Stats{24, 27, 0, 14, 30, 0, 0}},
	{Slot: 10, Name: "Magma Plume Boots", SourceZone: "H AC - Shirrak the Dead Watcher", SourceDrop: "", Stats: Stats{26, 24, 0, 14, 29, 0, 0}},
	{Slot: 10, Name: "Shimmering Azure Boots", SourceZone: "Securing the Celestial Ridge - NS Quest", SourceDrop: "", Stats: Stats{19, 0, 0, 16, 23, 0, 5}},
	{Slot: 10, Name: "Boots of Blashpemy", SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{29, 36, 0, 0, 36, 0, 0}},
	{Slot: 10, Name: "Boots of Ethereal Manipulation", SourceZone: "H Bot - Warp Splinter", SourceDrop: "", Stats: Stats{27, 27, 0, 0, 33, 0, 0}},
	{Slot: 10, Name: "Earthbreaker's Greaves", SourceZone: "Levixus the Soul Caller - Auchindoun Quest", SourceDrop: "", Stats: Stats{20, 27, 8, 0, 25, 0, 3}},
	{Slot: 10, Name: "Boots of the Nexus Warden", SourceZone: "The Flesh Lies... - Netherstorm Quest", SourceDrop: "", Stats: Stats{17, 27, 0, 18, 21, 0, 0}},
	{Slot: 11, Name: "Sparking Arcanite Ring", SourceZone: "H OHF - Epoch Hunter", SourceDrop: "", Stats: Stats{14, 13, 14, 10, 22, 0, 0}},
	{Slot: 11, Name: "Ring of Cryptic Dreams", SourceZone: "25 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{17, 16, 20, 0, 23, 0, 0}},
	{Slot: 11, Name: "Seer's Signit", SourceZone: "The Scryers - Exalted", SourceDrop: "", Stats: Stats{0, 24, 12, 0, 34, 0, 0}},
	{Slot: 11, Name: "Ring of Conflict Survival", SourceZone: "H MT - Yor (Summoned Boss)", SourceDrop: "", Stats: Stats{0, 28, 20, 0, 23, 0, 0}},
	{Slot: 11, Name: "Ryngo's Band of Ingenuity", SourceZone: "Arc - Wrath-Scryer Soccothrates", SourceDrop: "", Stats: Stats{14, 12, 14, 0, 25, 0, 0}},
	{Slot: 11, Name: "Band of the Guardian", SourceZone: "Hero of the Brood - CoT Quest", SourceDrop: "", Stats: Stats{11, 0, 17, 0, 23, 0, 0}},
	{Slot: 11, Name: "Scintillating Coral Band", SourceZone: "SV - Hydromancer Thespia", SourceDrop: "", Stats: Stats{15, 14, 17, 0, 21, 0, 0}},
	{Slot: 11, Name: "Manastorm Band", SourceZone: "Shutting Down Manaforge Ara - Quest", SourceDrop: "", Stats: Stats{15, 0, 10, 0, 29, 0, 0}},
	{Slot: 11, Name: "Ashyen's Gift", SourceZone: "Cenarion Expedition - Exalted", SourceDrop: "", Stats: Stats{0, 30, 0, 21, 23, 0, 0}},
	{Slot: 11, Name: "Cobalt Band of Tyrigosa", SourceZone: "H MT - Nexus-Prince Shaffar", SourceDrop: "", Stats: Stats{17, 19, 0, 0, 35, 0, 0}},
	{Slot: 11, Name: "Seal of the Exorcist", SourceZone: "50 Spirit Shards ", SourceDrop: "", Stats: Stats{0, 24, 0, 12, 28, 0, 0}},
	{Slot: 11, Name: "Lola's Eve", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{14, 15, 0, 0, 29, 0, 0}},
	{Slot: 11, Name: "Yor's Collapsing Band", SourceZone: "H MT - Yor (Summoned Boss)", SourceDrop: "", Stats: Stats{20, 0, 0, 0, 23, 0, 0}},
	{Slot: 11, Name: "Click Here for Trinket/Set Bonus Sims", SourceZone: "", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 11, Name: "Darkmoon Card: Crusade", SourceZone: "Blessings Deck", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 11, Name: "Scryer's Bloodgem", SourceZone: "The Scryers - Revered", SourceDrop: "", Stats: Stats{0, 0, 0, 32, 0, 0, 0}},
	{Slot: 11, Name: "Quagmirran's Eye", SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 37, 0, 0}},
	{Slot: 11, Name: "Arcanist's Stone", SourceZone: "H OHF - Epoch Hunter", SourceDrop: "", Stats: Stats{0, 0, 0, 25, 0, 0, 0}},
	{Slot: 11, Name: "Icon of the Silver Crescent", SourceZone: "41 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 43, 0, 0}},
	{Slot: 11, Name: "Shiffar's Nexus-Horn", SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{0, 0, 30, 0, 0, 0, 0}},
	{Slot: 11, Name: "Xi'ri's Gift", SourceZone: "The Sha'tar - Revered", SourceDrop: "", Stats: Stats{0, 0, 32, 0, 0, 0, 0}},
	{Slot: 11, Name: "Vengeance of the Illidari", SourceZone: "Cruel's Intentions/Overlord - HFP Quest", SourceDrop: "", Stats: Stats{0, 0, 26, 0, 0, 0, 0}},
	{Slot: 11, Name: "Figurine - Living Ruby Serpent", SourceZone: "Jewelcarfting BoP", SourceDrop: "", Stats: Stats{23, 33, 0, 0, 0, 0, 0}},
	{Slot: 19, Name: "Totem of the Void", SourceZone: "Mech - Cache of the Legion", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 19, Name: "Totem of the Pulsing Earth", SourceZone: "15 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 19, Name: "Totem of Impact", SourceZone: "15 Mark of Thrallmar/ Honor Hold", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 19, Name: "Totem of Lightning", SourceZone: "Colossal Menace - HFP Quest", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 0, 0, 0}},
	{Slot: 17, Name: "Gladiator's Spellblade / Gavel", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{18, 28, 0, 0, 199, 0, 0}},
	{Slot: 17, Name: "Eternium Runed Blade", SourceZone: "Blacksmithing BoE", SourceDrop: "", Stats: Stats{19, 0, 21, 0, 168, 0, 0}},
	{Slot: 17, Name: "Gavel of Unearthed Secrets", SourceZone: "Lower City - Exalted", SourceDrop: "", Stats: Stats{16, 24, 15, 0, 159, 0, 0}},
	{Slot: 17, Name: "Starlight Dagger", SourceZone: "H SP - Mennu the Betrayer", SourceDrop: "", Stats: Stats{15, 15, 0, 16, 121, 0, 0}},
	{Slot: 17, Name: "Runesong Dagger", SourceZone: "SH - Warbringer O'mrogg", SourceDrop: "", Stats: Stats{11, 12, 20, 0, 121, 0, 0}},
	{Slot: 17, Name: "Bleeding Hollow Warhammer", SourceZone: "H SP - Quagmirran", SourceDrop: "", Stats: Stats{17, 12, 16, 0, 121, 0, 0}},
	{Slot: 17, Name: "Sky Breaker", SourceZone: "H AC - Avatar of the Martyred", SourceDrop: "", Stats: Stats{20, 13, 0, 0, 132, 0, 0}},
	{Slot: 18, Name: "Mazthoril Honor Shield", SourceZone: "33 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{17, 16, 21, 0, 23, 0, 0}},
	{Slot: 18, Name: "Lamp of Peaceful Raidiance", SourceZone: "Arc - Harbinger Skyriss", SourceDrop: "", Stats: Stats{14, 13, 13, 12, 21, 0, 0}},
	{Slot: 18, Name: "Khadgar's Knapsack", SourceZone: "25 Badge of Justice - G'eras", SourceDrop: "", Stats: Stats{0, 0, 0, 0, 49, 0, 0}},
	{Slot: 18, Name: "Manual of the Nethermancer", SourceZone: "Mech - Nethermancer Sepethrea", SourceDrop: "", Stats: Stats{15, 12, 19, 0, 21, 0, 0}},
	{Slot: 18, Name: "Draenei Honor Guard Shield", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{16, 0, 21, 0, 19, 0, 0}},
	{Slot: 18, Name: "Star-Heart Lamp", SourceZone: "BM - Temporus", SourceDrop: "", Stats: Stats{18, 17, 0, 12, 22, 0, 0}},
	{Slot: 18, Name: "The Saga of Terokk", SourceZone: "Terokk's Legacy - Auchindoun Quest", SourceDrop: "", Stats: Stats{23, 0, 0, 0, 28, 0, 0}},
	{Slot: 18, Name: "Silvermoon Crest Shield", SourceZone: "SLabs - Murmur", SourceDrop: "", Stats: Stats{20, 0, 0, 0, 23, 0, 5}},
	{Slot: 18, Name: "Spellbreaker's Buckler", SourceZone: "Akama's Promise - SMV Quest", SourceDrop: "", Stats: Stats{10, 22, 0, 0, 29, 0, 0}},
	{Slot: 18, Name: "Hortus' Seal of Brilliance", SourceZone: "SH - Warchief Kargath Bladefist", SourceDrop: "", Stats: Stats{20, 18, 0, 0, 23, 0, 0}},
	{Slot: 18, Name: "Gladiator's Endgame", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{14, 21, 0, 0, 19, 0, 0}},
	{Slot: 17, Name: "Gladiator's War Staff", SourceZone: "Arena Season 1 Reward", SourceDrop: "", Stats: Stats{35, 48, 36, 21, 199, 0, 0}},
	{Slot: 17, Name: "Terokk's Shadowstaff", SourceZone: "H SH - Talon King Ikiss", SourceDrop: "", Stats: Stats{42, 40, 37, 0, 168, 0, 0}},
	{Slot: 17, Name: "Auchenai Staff", SourceZone: "The Aldor - Revered", SourceDrop: "", Stats: Stats{46, 0, 26, 19, 121, 0, 0}},
	{Slot: 17, Name: "Warpstaff of Arcanum", SourceZone: "Bot - Warp Splinter", SourceDrop: "", Stats: Stats{38, 37, 26, 16, 121, 0, 0}},
	{Slot: 17, Name: "The Bringer of Death", SourceZone: "BoE World Drop", SourceDrop: "", Stats: Stats{31, 32, 42, 0, 121, 0, 0}},
	{Slot: 17, Name: "Bloodfire Greatstaff", SourceZone: "BM - Aeonus", SourceDrop: "", Stats: Stats{42, 42, 28, 0, 121, 0, 0}},
	{Slot: 17, Name: "Ameer's Impulse Taser", SourceZone: "Nexus-King Salhadaar - Netherstorm Quest", SourceDrop: "", Stats: Stats{27, 27, 27, 17, 103, 0, 0}},
	{Slot: 17, Name: "Grand Scepter of the Nexus-Kings", SourceZone: "H MT - Nexus-Prince Shaffar", SourceDrop: "", Stats: Stats{43, 45, 0, 19, 121, 0, 0}},
}
