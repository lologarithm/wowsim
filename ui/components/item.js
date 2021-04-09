// Individual Items in Gear

class ItemComponent {
    slot;
    name;
    innerdiv;
    img;
    maindiv;
    selector;

    socketComp;

    constructor(slot, allgems) {
        this.slot = slot;
        this.selector = new SelectorComponent(slot, allgems);

        var name = document.createElement("p");
        name.addEventListener("click", (e) => {this.showItemSelector(e)});
        name.id = slot+"label";
        name.innerText = "None";
        
        this.socketComp = new SocketsComponent(slot);
        this.socketComp.addSelectListener((socket, e)=>{
            this.showGemSelector(e, socket); ///ya ya ya, reversed parameters, ill deal with it later.
        });
        this.socketComp.div.addEventListener("click", (e) => {});

        var innerdiv = document.createElement("div");
        innerdiv.classList.add("equiplabel");
        innerdiv.appendChild(name);
        innerdiv.appendChild(this.socketComp.div);
    
        var img = document.createElement("img");
        img.id = slot+"icon";
        img.src = "../icons/Items/Temp.png"
        img.addEventListener("click", (e) => {this.showItemSelector(e)});

        var maindiv = document.createElement("div");
        maindiv.id = slot;
        maindiv.classList.add("equipslot");
        maindiv.appendChild(img);
        maindiv.appendChild(innerdiv);
        maindiv.appendChild(this.selector.selectordiv);

        this.name = name;
        this.innerdiv = innerdiv;
        this.img = img;
        this.maindiv = maindiv;
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

    updateEquipped(newItem) {
        if (newItem != null && newItem.Name != "") {
            this.name.innerText = newItem.Name;
            // gearlist.push(newItem.Name);

            this.img.src = slotToIcon[this.slot];
            if (newItem.GemSlots != null ) {
                 // updates the selector UI with the current gems/enchants (later)
                this.selector.updateEquipped(newItem);
                this.socketComp.updateSockets(newItem.GemSlots, newItem.Gems);
            }
            // gemdiv.innerHTML += '<div class="enchslot" style="float: right;"></div>';
        } else {
            this.name.innerText = "None";
            this.img.src = "";
            this.socketComp.updateSockets([], []);
        }
    }
}

function moveSelector(box, x, y) {
    if (x < 0) {
        x = 0;
    }
    if (y < 0) {
        y = 0;
    }
    if (x > window.innerWidth-400) {
        x = window.innerWidth - 400;
    }
    box.style.left = x;
    box.style.top = y;
}

// Click handler for 'remove' button on each slot.
function removeGear(event) {
    var ele = event.target;
    var slotid = ele.parentElement.parentElement.id;
    currentGear[slotid] = {Name: "None"};
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

    constructor() {
        this.div = document.createElement("div");
        this.div.style.height = "2.2em"; // so the text doesnt go right of the socket icons...
        
        this.selectedSocket = 0;
        this.sockets = [];
        this.listeners = [];
    }

    updateSockets(sockets, gems) {
        this.gems = gems;
        this.div.innerHTML = "";
        this.sockets = sockets;
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
            socketDiv.style.backgroundColor = color+",0.3)";
            socketDiv.style.border = '1px solid ' + color+",0.8)";
            
            socketDiv.addEventListener("click", (event)=>{
                this.listeners.forEach((h)=>{ h(idx, event) });
            });

            // TODO: gotta be a cleaner way to do this... ill fix it later.
            if (gems != null && gems != undefined) {
                if (gems[idx].Name != null) {
                    var img = document.createElement("img")
                    img.src = gemToIcon[gems[idx].Color]
                    if (img.src != undefined) {
                        socketDiv.appendChild(img);
                    }
                }
            }
            this.div.appendChild(socketDiv);
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

// For now hardcode an icon.
var gemToIcon = {
    1: "../icons/Gems/Gem_Diamond_07.png",
    2: "../icons/Gems/Gem_BloodGem_02.png",
    3: "../icons/Gems/Gem_AzureDraenite_02.png", // blue
    4: "../icons/Gems/Gem_GoldenDraenite_02.png", // yellow
    5: "../icons/Gems/Gem_DeepPeridot_02.png", // green
    6: "../icons/Gems/Gem_FlameSpessarite_02.png", // orange
    7: "../icons/Gems/Gem_EbonDraenite_02.png", // purple
}

