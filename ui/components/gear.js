// Gear tab

// converts a slot number to its string id. Useful for finding the slot in the UI for an item.
var slotToID = [
    "equipnone",
    "equiphead",
    "equipneck",
    "equipshoulder",
    "equipback",
    "equipchest",
    "equipwrist",
    "equiphands",
    "equipwaist",
    "equiplegs",
    "equipfeet",
    "equipfinger",
    "equipfinger1",
    "equipfinger2",
    "equiptrinket",
    "equiptrinket1",
    "equiptrinket2",
    "equipweapon",
    "equipoffhand",
    "equiptotem"
]

class GearUI {
    allitems;
    allgems;
    
    compdiv;
    
    itemCompSlots; // index of itemComps by slot
    itemComps;
    currentGear;

    changeHandlers;

    constructor(node, gearList) {
        this.changeHandlers = [];
        this.allitems = {};
        this.allgems = {};
        this.allenchants = {};
        this.itemComps = [];
        this.itemCompSlots = {};
        this.currentGear = {};

        var holderl = document.getElementById("gearleft");
        var holderr = document.getElementById("gearright");

        var holder = holderl;
        slotToID.forEach((slot) => {
            // ya ya ya, this is terrible code. I dont have time right now.
            if (slot == "equipnone" || slot == "equipfinger" || slot == "equiptrinket") {
                return;
            }
            if (slot == "equiphands") {
                holder = holderr;
            }
            if (slot == "equipweapon") {
                holder = holderl;
            }
            if (slot == "equiptotem") {
                holder = holderr;
            }
            
            var itemComp = new ItemComponent(slot, gearList);
            itemComp.selector.addSelectedListener((change)=>{
                if (change.item != null) {
                    var item = Object.assign({Name: ""},this.allitems[change.item]);
                    this.currentGear[slot] = item;
                    itemComp.updateEquipped(item);
                } else if (change.gem != null) {
                    if (this.currentGear[slot].Gems == null) {
                        this.currentGear[slot].Gems = [];
                        this.currentGear[slot].GemSlots.forEach(()=>{
                            this.currentGear[slot].Gems.push({})
                        });
                    }
                    this.currentGear[slot].Gems[change.gem.socket] = this.allgems[change.gem.name];
                    itemComp.updateEquipped(this.currentGear[slot]);
                } else if (change.enchant != null) {
                    this.currentGear[slot].Enchant = this.allenchants[change.enchant.name];
                    itemComp.updateEquipped(this.currentGear[slot]);
                }

                this.changeHandlers.forEach((h)=>{
                    h(item, slot);
                });
                itemComp.selector.hide();
            });
            this.itemCompSlots[slot] = itemComp;
            this.itemComps.push(itemComp);

            holder.appendChild(itemComp.maindiv);
        });
        
        gearList.Items.forEach(g => {
            if (g.GemSlots != null) {
                var gemslots = [];
                var gbytes = atob(g.GemSlots);
                for (var i = 0; i < gbytes.length; i++) {
                    gemslots.push(gbytes.codePointAt(i));
                }
                g.GemSlots = gemslots;
            }
            this.allitems[g.Name] = g;

            var slot = slotToID[g.Slot];
            if (slot == "equipfinger" || slot == "equiptrinket") {
                this.itemCompSlots[slot+"1"].addItem(g);
                this.itemCompSlots[slot+"2"].addItem(g);
            } else {
                this.itemCompSlots[slot].addItem(g);
            }
        });
        gearList.Gems.forEach((gem)=>{
            this.allgems[gem.Name] = gem;
        });
        gearList.Enchants.forEach((enchant)=>{
            this.allenchants[enchant.Name] = enchant;
        });
    
        this.compdiv = node;
    }

    // triggered on a click of an item in selector
    // or on pressing 'enter' in a search box
    addChangeListener(handler) {
        this.changeHandlers.push(handler);
    }

    updateItemSlot(newItem, slotid) {
        this.currentGear[slotid] = newItem;
        this.itemCompSlots[slotid].updateEquipped(newItem);
    }

    removeEquipped() {
        var emptyItem = {Name: "None"};
        Object.keys(this.currentGear).forEach((key) => {
            // remove item.
            this.currentGear[key] = emptyItem;
            this.itemCompSlots[key].updateEquipped(null);
        });
    }

    setPhase(filter) {
        this.itemComps.forEach((comp)=>{
            comp.selector.setPhase(filter);
        });
    }
    setFilter(filter) {
        this.itemComps.forEach((comp)=>{
            comp.selector.setFilter(filter);
        });
    }

    // updateEquipped will update the gear UI elements (to redraw when new gear is selected)
    updateEquipped(newGear) {
        this.removeEquipped();

        var finger1done = false;
        var trink1done = false;

        // Take each item, find its slot
        newGear.forEach( (item) => {
            var inm = item.Name;
            if (inm == "" || inm == "None") {
                return;
            }
            var realItem = Object.assign({}, this.allitems[inm]);
            if (realItem == null || realItem.Name == null) {
                return;
            }
            var slotid = slotToID[realItem.Slot];
    
            if (slotid == "equipfinger") {
                if (!finger1done) {
                    slotid = "equipfinger1";
                    finger1done = true;
                } else {
                    slotid = "equipfinger2";
                }
            } else if (slotid == "equiptrinket") {
                if (!trink1done) {
                    slotid = "equiptrinket1";
                    trink1done   = true;
                } else {
                    slotid = "equiptrinket2";
                }
            }
            if (item.Gems != null && item.Gems.length > 0) {
                realItem.Gems = [];
                item.Gems.forEach((g, idx) => {
                    var gem = this.allgems[g];
                    if (gem == null) {
                        gem = {}; // empty object for gem sentinal?
                    }
                    realItem.Gems.push(gem);
                });
            }
            if (item.Enchant != null && item.Enchant != "") {
                realItem.Enchant = this.allenchants[item.Enchant];
            }
            this.updateItemSlot(realItem, slotid)
        });
    
        // this.currentGear = newGear;
        return this.currentGear;
    }

    hideSelectors(e) {
        this.itemComps.forEach((ic) => {
            if (e != null && ic.maindiv.contains(e.target)) {
                return; // dont close selectors for the item we are clicking on.
            }
            ic.selector.hide(e);
        });
    }

    gearList() {
        var list = [];
        Object.entries(this.currentGear).forEach( (entry)=>{
            list.push(entry[1].Name);
        } );

        return list;
    }
}

// For now hardcode an icon.
var slotToIcon = {
    "equiphead": "../icons/Armor/INV_Helmet_06.png",
    "equipneck": "../icons/Armor/INV_Jewelry_Necklace_07.png",
    "equipshoulder": "../icons/Armor/INV_Shoulder_14.png",
    "equipback": "../icons/Armor/INV_Misc_Cape_16.png",
    "equipchest": "../icons/Armor/INV_Chest_Chain_04.png",
    "equipwrist": "../icons/Armor/INV_Bracer_09.png",
    "equiphands": "../icons/Armor/INV_Gauntlets_26.png",
    "equipwaist": "../icons/Armor/INV_Belt_19.png",
    "equiplegs": "../icons/Armor/INV_Pants_03.png",
    "equipfeet": "../icons/Armor/INV_Boots_Wolf.png",
    "equipfinger1": "../icons/Armor/INV_Jewelry_Ring_04.png",
    "equipfinger2": "../icons/Armor/INV_Jewelry_Ring_05.png",
    "equiptrinket1": "../icons/Armor/INV_Jewelry_Talisman_09.png",
    "equiptrinket2": "../icons/Armor/INV_Jewelry_Talisman_10.png",
    "equipweapon": "../icons/Weapons/INV_Sword_39.png",
    "equipoffhand": "../icons/Armor/INV_Shield_20.png",
    "equiptotem": "../icons/Spells/Spell_Nature_InvisibilityTotem.png"
}

