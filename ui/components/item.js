// Individual Items in Gear

class ItemComponent {
    slot;
    filterLevel;

    // UI Elements
    name;
    innerdiv;
    img;
    maindiv;
    // statpop;

    // Sub Components
    selector;
    socketComp;

    constructor(slot, allgear) {
        this.filterLevel = 0;
        this.slot = slot;
        this.selector = new SelectorComponent(slot, allgear);

        var name = document.createElement("p");
        name.addEventListener("click", (e) => { this.showItemSelector(e) });
        name.innerText = "None";

        this.socketComp = new SocketsComponent(slot, allgear.Enchants);
        this.socketComp.addSelectListener((socket, e) => {
            if (socket == -1) {
                this.showEnchantSelector(e);
                return;
            }
            this.showGemSelector(e, socket); ///ya ya ya, reversed parameters, ill deal with it later.
        });
        this.socketComp.div.addEventListener("click", (e) => { });

        var innerdiv = document.createElement("div");
        innerdiv.classList.add("equiplabel");
        innerdiv.appendChild(name);
        innerdiv.appendChild(this.socketComp.div);

        var img = document.createElement("img");
        img.id = slot + "icon";
        img.src = "";
        
        var maindiv = document.createElement("div");
        maindiv.id = slot;
        maindiv.classList.add("equipslot");

        var imgLink = document.createElement("a");
        imgLink.target = "_blank";
        imgLink.appendChild(img);
        maindiv.appendChild(imgLink);
        
        maindiv.appendChild(innerdiv);
        maindiv.appendChild(this.selector.selectordiv);

        this.name = name;
        this.innerdiv = innerdiv;
        this.img = img;
        this.maindiv = maindiv;
    }

    showEnchantSelector(event, socket) {
        this.selector.show(event);
        this.selector.focus("enchant");
        moveSelector(this.selector.selectordiv, event.clientX, event.clientY);
    }

    showGemSelector(event, socket) {
        this.selector.show(event);
        this.selector.focus("gem", socket);
        moveSelector(this.selector.selectordiv, event.clientX, event.clientY);
    }

    showItemSelector(event) {
        this.selector.show(event);
        this.selector.focus("item");
        moveSelector(this.selector.selectordiv, event.clientX, event.clientY);
    }

    addItem(newItem) {
        this.selector.addItem(newItem);
    }

    setWeights(weights) {
        this.selector.setWeights(weights);
    }

    updateEquipped(newItem) {
        console.log("New Item Equipped: ", newItem);
        if (newItem != null && newItem.Name != "") {
            this.name.innerText = newItem.Name;
            // remove old quality classnames
            this.name.classList.forEach(classname => {
                if (classname.includes("quality")) {
                    this.name.classList.remove(classname);
                }
            });
            this.name.classList.add(itemQualityCssClass(newItem.Quality));

            let wowheadData = ""
            if (newItem.Gems && newItem.Gems.length > 0) {
                wowheadData += "gems=";
                newItem.Gems.forEach((g, i) => {
                    if (i != 0) {
                        wowheadData += ":";
                    }
                    wowheadData += g.ID;
                });
            }
            if (newItem.Enchant) {
                if (wowheadData.length > 0) {
                    wowheadData += "&"
                }
                wowheadData += "ench=" + newItem.Enchant.EffectID;
            }

            setItemIcon(newItem.ID, this.img);
          
            this.img.parentNode.setAttribute("href", "https://tbc.wowhead.com/item=" + newItem.ID);
            this.img.parentNode.setAttribute("data-wowhead", wowheadData);
            // updates the selector UI with the current gems/enchants (later)
            this.selector.updateEquipped(newItem);
            this.socketComp.updateSockets(newItem.GemSlots, newItem.Gems);
        } else {
            this.name.innerText = "None";
            this.img.src = "";
            this.socketComp.updateSockets([], []);
        }
    }
}

var statnames = ["Int", "Stm", "SpellCrit", "SpellHit", "SpellDmg", "Haste", "MP5", "Mana", "SpellPen"];

function moveSelector(box, x, y) {
    if (x < 0) {
        x = 0;
    }
    if (y < 0) {
        y = 0;
    }
    if (x > window.innerWidth - 400) {
        x = window.innerWidth - 400;
    }
    box.style.left = x;
    box.style.top = y;
}

// Click handler for 'remove' button on each slot.
function removeGear(event) {
    var ele = event.target;
    var slotid = ele.parentElement.parentElement.id;
    currentGear[slotid] = { Name: "None" };
    updateGear(currentGear);
}


class SocketsComponent {
    // UI
    div;

    // State
    sockets;
    gems; // currently socketed gems.
    selectedSocket;
    listeners;
    slot; // Slot this selector is for
    enchant; // enchant in the enchant socket
    enchants; // list of all enchants

    constructor(slot, enchants) {
        this.enchants = enchants;
        this.div = document.createElement("div");
        this.div.style.height = "2.2em"; // so the text doesnt go right of the socket icons...

        this.slot = slot;
        this.selectedSocket = 0;
        this.sockets = [];
        this.listeners = [];
    }

    updateSockets(sockets, gems) {
        this.gems = gems;
        this.div.innerHTML = "";
        this.sockets = sockets;
        if (sockets != null) {
            sockets.forEach((socket, idx) => {
                var color = "rgba(30, 30, 30";
                if (socket == 2) {
                    color = "rgba(250, 30, 30";
                } else if (socket == 3) {
                    color = "rgba(30, 30, 250";
                } else if (socket == 4) {
                    color = "rgba(250, 250, 30";
                }
                var socketDiv = document.createElement("div");
                socketDiv.classList.add("gemslot");
                socketDiv.style.backgroundColor = color + ",0.3)";
                socketDiv.style.border = '1px solid ' + color + ",0.8)";

                socketDiv.addEventListener("click", (event) => {
                    this.listeners.forEach((h) => { h(idx, event) });
                });

                // TODO: gotta be a cleaner way to do this... ill fix it later.
                if (gems && gems[idx]) {
                    if (gems[idx].Name != null) {
                        var img = document.createElement("img")
                        setItemIcon(gems[idx].ID, img);
                        var text = gems[idx].Name;
                        gems[idx].Stats.forEach((v, i) => {
                            if (v > 0) {
                                text += `\n${statnames[i]}: ${v.toString()}`;
                            }
                        });
                        img.title = text;

                        if (img.src != undefined) {
                            socketDiv.appendChild(img);
                        }
                    }
                }
                this.div.appendChild(socketDiv);
            });
        }

        var addedEnc = false;
        this.enchants.forEach((e) => {
            if (addedEnc) { return; }
            var slot = this.slot;
            if (slot == "equipfinger1" || slot == "equipfinger2") {
                slot = "equipfinger";
            }
            if (slotToID[e.Slot] == slot) {
                var enchdiv = document.createElement("div");
                enchdiv.classList.add("enchslot");
                enchdiv.addEventListener("click", (event) => {
                    this.listeners.forEach((h) => { h(-1, event) });
                });
                this.div.appendChild(enchdiv);
                addedEnc = true;
            }
        });
    }

    selectSocket(idx, event) {
        // activate socket and then 
        this.selectedSocket = idx;
        for (var i = 0; i < this.div.childNodes.length; i++) {
            if (i == idx) {
                this.div.childNodes[i].style.border = "1px solid white";
            } else {
                this.div.childNodes[i].style.border = "1px solid black";
            }
        }
    }

    addSelectListener(hnd) {
        this.listeners.push(hnd);
    }
}

function setItemIcon(id, imgElement) {
    if (itemToIcon.has(id)) {
        imgElement.src = itemToIcon.get(id);
    } else {
        fetch('https://tbc.wowhead.com/tooltip/item/'+id)
        .then(response => response.json())
        .then(itemInfo => {
            itemToIcon.set(id, "https://wow.zamimg.com/images/wow/icons/large/" + itemInfo.icon + ".jpg");
            imgElement.src = itemToIcon.get(id);
        });
    }
}

var itemToIcon = new Map();