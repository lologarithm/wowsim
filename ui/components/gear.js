// Gear tab

// converts a slot number to its string id. Useful for finding the slot in the UI for an item.
const slotToID = [
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
    allenchants;

    itemsByID;
    enchantsByID;
    gemsByID;

    compdiv;
    
    itemCompSlots; // index of itemComps by slot
    itemComps;
    currentGear;

    changeHandlers;

    constructor(node, gearList) {
        this.changeHandlers = [];
        this.allitems = {};
        this.itemsByID = {};
        this.allgems = {};
        this.gemsByID = {};
        this.allenchants = {};
        this.enchantsByID = {};
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
                    // remove any pointers to the root item after cloning properties.
                    item.Enchant = null; 
                    item.Gems = null;
                    this.currentGear[slot] = item;

                    // If the item is a 2h, remove any offhand that is currently equipped.
                    if (item.subSlot == 2  && this.currentGear["equipoffhand"] != null) {
                        this.currentGear["equipoffhand"] = null;
                        this.itemCompSlots["equipoffhand"].updateEquipped(null);
                    }
                    
                    itemComp.updateEquipped(item);
                } else if (change.gem != null) {
                    if (this.currentGear[slot].Gems == null) {
                        this.currentGear[slot].Gems = [];
                        this.currentGear[slot].GemSlots.forEach(()=>{
                            this.currentGear[slot].Gems.push({})
                        });
                    }
                    if (change.gem.socket == -1) {
                        // all sockets
                        if (change.gem.name == "none") {
                            this.currentGear[slot].Gems = null;
                        }
                        // TODO: set all slots to a real gem
                    } else {
                        this.currentGear[slot].Gems[change.gem.socket] = this.allgems[change.gem.name];
                    }
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
            this.itemsByID[g.ID] = g;

            var slot = slotToID[g.Slot];
            if (slot == "equipfinger" || slot == "equiptrinket") {
                this.itemCompSlots[slot+"1"].addItem(g);
                this.itemCompSlots[slot+"2"].addItem(g);
            } else {
                this.itemCompSlots[slot].addItem(g);
            }
        });
        this.setWeights(defaultStatWeights);

        gearList.Gems.forEach((gem)=>{
            this.allgems[gem.Name] = gem;
            this.gemsByID[gem.ID] = gem;
        });
        gearList.Enchants.forEach((enchant)=>{
            this.allenchants[enchant.Name] = enchant;
            this.enchantsByID[enchant.ID] = enchant;
        });
    
        this.compdiv = node;
    }

    // triggered on a click of an item in selector
    // or on pressing 'enter' in a search box
    addChangeListener(handler) {
        this.changeHandlers.push(handler);
    }

    // Sets the EP weights for all items.
    setWeights(weights) {
        Object.values(this.itemCompSlots).forEach(itemComp => itemComp.setWeights(weights));
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
        newGear.forEach(item => {
            const nameOrId = item.ID || item.Name;
            const realItem = this.allitems[nameOrId] || this.itemsByID[nameOrId];
            if (!realItem) {
              return;
            }

            let slotid = slotToID[realItem.Slot];
    
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

            if (item.Gems) {
              realItem.Gems = item.Gems.map(gem => {
                if (!gem) {
                  return {};
                }

                if (typeof gem === 'string') {
                  return this.allgems[gem] || {};
                }
                return this.allgems[gem.Name] || this.gemsByID[gem.ID] || {};
              });
            }
            if (item.g) {
                realItem.Gems = item.g.map(gem => this.gemsByID[gem] || {});
            }
            if (item.Enchant) {
                realItem.Enchant = this.allenchants[item.Enchant] 
                    || this.allenchants[item.Enchant.Name] 
                    || this.enchantsByID[item.Enchant.ID];
            }
            if (item.e && item.e > 0) {
                realItem.Enchant = this.enchantsByID[item.e];
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
