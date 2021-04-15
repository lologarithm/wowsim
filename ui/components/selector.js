// Item, Gem, and Enchant selector

class SelectorComponent {
    // UI Elements
    search; // search box for finding items.
    clearbutton; // Button to clear current item.
    itemselector; // UI to search & choose a different item
    selectorlist; // Container of items that can be selected
    seltabs;
    tab1;
    tab2;
    tab3;
    gemseldiv;
    gemsList; // List of gems to select when socket is selected.

    // Main Div
    selectordiv;
    
    // State
    slot;
    
    highlighted; // currently highlighted item
    highText; // text of highlighted item
    items; // searchable items, mirrors selectorlist.children
    
    // SubComp
    sockComp;

    allgear;
    constructor(slot, allgear) {
        this.allgear = allgear;
        this.changeHandlers = [];
        this.items = [];
        this.slot = slot;
        this.sockComp = new SocketsComponent(slot);
        this.gemsList = document.createElement("div");
        this.sockComp.addSelectListener((socket)=>{
            this.gemSelected(socket);
        });

        var gemseldiv = document.createElement("div");
        gemseldiv.style.display = "none";
        gemseldiv.appendChild(this.sockComp.div);
        gemseldiv.appendChild(this.gemsList);

        var search = document.createElement("input");
        search.addEventListener("keyup", (e) => {this.searchHandler(e)});

        var closebut = document.createElement("button");
        closebut.innerText = "X";
        closebut.style.backgroundColor = "red";
        closebut.style.marginLeft = "5px";
        closebut.addEventListener("click", (e) => {this.hide()});

        var clearbutton = document.createElement("button");
        clearbutton.innerText = "Remove";
        clearbutton.addEventListener("click", (e) => {this.notifyItemChange("None")});
    
        var itemselector = document.createElement("div");
        var selectorlist = document.createElement("div");
        itemselector.appendChild(search);
        itemselector.appendChild(selectorlist);

        var seltabs =  document.createElement("div");
        seltabs.classList.add("selectortabs");
        var tab1 = document.createElement("div");
        tab1.innerText = "Item";
        tab1.classList.add("selectortab");
        tab1.addEventListener("click", ()=> { this.focus("item") });

        var tab2 = document.createElement("div");
        tab2.innerText = "Gems";
        tab2.classList.add("selectortab");
        tab2.addEventListener("click", ()=> { this.focus("gem") });
        
        var tab3 = document.createElement("div");
        tab3.innerText = "Enchants";
        tab3.classList.add("selectortab");
        seltabs.appendChild(tab1);
        seltabs.appendChild(tab2);
        seltabs.appendChild(tab3);

        var selectordiv = document.createElement("div");
        selectordiv.classList.add("equipselector");
        if (theme == "dark") {
            selectordiv.classList.add("dtl");
        } else {
            selectordiv.classList.add("ltl");
        }
        selectordiv.style.display = "none";
        selectordiv.appendChild(closebut);
        selectordiv.appendChild(clearbutton);
        selectordiv.appendChild(seltabs);
        selectordiv.appendChild(gemseldiv);
        selectordiv.appendChild(itemselector);


        this.gemseldiv = gemseldiv;
        this.search = search;
        this.clearbutton = clearbutton;
        this.itemselector = itemselector;
        this.selectorlist = selectorlist;
        this.seltabs = seltabs;
        this.tab1 = tab1;
        this.tab2 = tab2;
        this.tab3 = tab3;
        this.selectordiv = selectordiv;

    }

    changeHandlers;

    // triggered on a click of an item in selector
    // or on pressing 'enter' in search box
    addSelectedListener(handler) {
        this.changeHandlers.push(handler);
    }

    // Updates the socket selector UI
    updateEquipped(item) {
        if (item.GemSlots == null || item.GemSlots.length == 0) {
            return;
        }

        this.sockComp.updateSockets(item.GemSlots, item.Gems);
    }

    // Adds a new item that can be searched.
    addItem(item) {
        var listItem = document.createElement("div");
        listItem.classList.add("equipselitem");
        listItem.innerText = item.Name;
        listItem.addEventListener("click", (e)=>{
            this.gearClickHandler(item.Name);
        });
            
        this.selectorlist.appendChild(listItem);
        this.items.push(item.Name);
    }

    show() {
        this.selectordiv.style.display = "block";
    }
    focus(tab, subitem) {
        if (tab == "item") {
            this.tab1.classList.add("selactive");
            this.tab2.classList.remove("selactive");
            this.gemseldiv.style.display = "none";
            this.itemselector.style.display = "block";
            this.search.focus();
        } else if (tab == "gem") {
            this.tab2.classList.add("selactive");
            this.tab1.classList.remove("selactive");
            this.gemseldiv.style.display = "block";
            this.itemselector.style.display = "none";
            if (subitem != null) {
                // activate a specific gem
                this.gemSelected(subitem);
            }
            // this.sockComp.
        }
    }
    hide(e) {
        if (e == null) {
            // by default just hide
            this.selectordiv.style.display = "none";
            return;
        }

        // If we are provided a click location, see if its clicking us or not.
        if (!this.selectordiv.contains(e.target)) {
            this.selectordiv.style.display = "none";
        }
    }
    gemSelected(socket) {
        var newList = this.gemSelector(this.sockComp.sockets[socket]); // get color of socket to filter new gems list.
        this.gemseldiv.replaceChild(newList, this.gemsList); // replace the list inside the gem selector
        this.gemsList = newList;
        this.sockComp.selectSocket(socket); // now activate in the gem socket selector UI
    }

    gemSelector(color) {
        var div = document.createElement("div");
        div.classList.add("gemselectorlist")
        Object.entries(this.allgear.Gems).filter((v) => {
            return true; //colorIntersects(color, v[1].Color);
        }).forEach((gem) => {
            var name = gem[1].Name;
            var itemdiv = document.createElement("div");
            itemdiv.innerText = name
            itemdiv.addEventListener("click", (e)=>{
                // using global itemselector here feels weird...
                this.notifyGemChange(name, this.sockComp.selectedSocket);
            });
            div.appendChild(itemdiv);
        });
        return div;
    }

    // Item Search
    
    // turns a search string into a list of lowercase terms and if there is punctuation or not.
    getSearchTerms(text) {
        // Use lower case in search
        var search = text.toLowerCase();

        // Check for punctuation. ya ya ya, I should just use regex check instead of replace.
        var snp = search.replace(/[^\w\s]|_/g, "");
        var doPunc = false;
        if (snp != search) {
            doPunc = true;
        }
        
        // Now split search into terms
        var sterms = search.split(' ')

        return {"terms": sterms, "punc": doPunc}
    }

    // takes search result from 'getSearchTerms' and searches the value string.
    find(search, value) {
        value = value.toLowerCase();
        if (!search.punc) {
            value = value.replace(/[^\w\s]|_/g, "");
        }
        
        var found = true;
        search.terms.forEach(st => {
            // If value is missing any search term, break search.
            if (!value.includes(st)) {
                found = false;
                return false;
            }
        });
        return found;
    }

    notifyGemChange(gem, socket) {
        this.changeHandlers.forEach((h)=>{
            h({gem: {name: gem, socket: socket}});
        });
    }

    notifyItemChange(value) {
        this.changeHandlers.forEach((h)=>{
            h({item: value});
        });
    }

    // completeSearch clears the search area and notifys of the searched item.
    completeSearch() {
        var itemText = this.highText;
        this.clearSearch(); // clears current search state.
        return itemText;
    }

    arrow(search, last, F) {
        return function(i, node, item) {
            if (F(search, item)) {
                node.childNodes[i].classList.add("lisearch");
                node.childNodes[last].classList.remove("lisearch");
                return true; // exit now.
            }
            return false;
        }
    }
    handleSearchDown(hndlr) {
        for (var i = this.highlighted+1; i < this.items.length; i++) {
            if (hndlr(i, this.selectorlist, this.items[i])) {
                this.highText = this.items[i];
                this.highlighted = i;
                return;
            }
        }
    }
    handleSearchUp(hndlr) {
        for (var i = this.highlighted-1; i>=0; i--) {
            if (hndlr(i, this.selectorlist, this.items[i])) {
                this.highText = this.items[i];
                this.highlighted = i;
                return;
            }
        }
    }
    // Uses text from element to find item slot list.
    // Ignores case and punctuation unless punctuation is included in the search.
    // Spaces in search are implicit 'and'
    searchHandler(event) {
        if (event.code == "Enter" && this.highText != "") { // Enter
            this.notifyItemChange(this.completeSearch());
            return;
        }

        var search = this.getSearchTerms(this.search.value);
        
        // Handle direction function.
        var hndlr = this.arrow(search, this.highlighted, this.find);
        if (event.code == "ArrowUp") {
            this.handleSearchUp(hndlr);
            event.preventDefault();
            return
        } else if (event.code == "ArrowDown") {
            this.handleSearchDown(hndlr);
            event.preventDefault();
            return
        }

        var firstFound = false;
        for (var i = 0; i < this.items.length; i++) {
            var found = this.find(search, this.items[i]);

            // Show / Hide item
            var le = this.selectorlist.childNodes[i];
            if (found) {
                le.style.removeProperty("display");
                if (!firstFound) {
                    this.highText = this.items[i];
                    this.highlighted = i;
                    le.classList.add("lisearch");
                    firstFound = true;
                } else {
                    le.classList.remove("lisearch");
                }
            } else {
                le.style.display = "none";
                le.classList.remove("lisearch");
            }
        }
    }

    clearSearch() {
        this.search.value = "";

        // now iterate all the children, clearing temp search state.
        var numChild = this.selectorlist.childNodes.length;
        for (var i = 0; i < numChild; i++) {
            var le = this.selectorlist.children[i];
            le.style.removeProperty("display");
            le.classList.remove("lisearch");
        }

        this.highlighted = -1;
        this.highText = "";
    }

    // Click handler for each item in slot list dropdown.
    gearClickHandler(name) {
        this.notifyItemChange(name);
        this.clearSearch();
    }

}


function colorIntersects(color, intersect) {
    if (color == intersect) {
        return true;
    }
    if (color == 8 || intersect == 8) { // prismatic intersects everything.
        return true;
    }

    if (color == 1) {
        return false; // meta gems intersect nothing.
    }
    if (color == 2) { // red
        return intersect == 6 || intersect == 7
    }
    if (color == 3) { // blue
        return intersect == 5 || intersect == 7
    }
    if (color == 4) { // yellow
        return intersect == 5 || intersect == 6
    }

    return false; // dunno wtf this is.
}
