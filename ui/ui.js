
// Globals
var allgear = {};
var defaultGear = ["Shamanistic Helmet of Second Sight","Brooch of Heightened Potential","Pauldrons of Wild Magic","Ogre Slayer's Cover","Tidefury Chestpiece","World's End Bracers","Earth Mantle Handwraps","Wave-Song Girdle","Stormsong Kilt","Magma Plume Boots","Cobalt Band of Tyrigosa","Scintillating Coral Band","Totem of the Void","Mazthoril Honor Shield","Bleeding Hollow Warhammer", "Quagmirran's Eye", "Icon of the Silver Crescent"];
var currentGear = {};

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

var simlib = new window.Worker(`simworker.js`);
var simlib2 = new window.Worker(`simworker.js`);

var simlibBusy = false;
var simlib2Busy = false;

simlib.onmessage = (event) => {
    var m = event.data.msg;
    if (m == "ready") {
        simlib.postMessage({msg: "setID", payload: "1"});
        simlib.postMessage({msg: "getGearList"});
    } else if (m == "getGearList") {
        // do something
        popgear(event.data.payload);
    } else {
        var onComp = simrequests[event.data.id];
        if (onComp != null) {
            onComp(event.data.payload);
        }
        simlibBusy = false;
    }
}

simlib2.onmessage = (event) => {
    var m = event.data.msg;
    if (m == "ready") {
        simlib2.postMessage({msg: "setID", payload: "2"});
        return;
    }
    var onComp = simrequests[event.data.id];
    if (onComp != null) {
        onComp(event.data.payload);
    }
}

var simrequests = {};
function simulate(iters, dur, gearlist, opts, rots, haste, onComplete) {
    console.log("Called Simulate... rots:", rots);
    var id = makeid();
    simrequests[id] = onComplete
    var worker = simlib;
    if (simlibBusy) {
        worker = simlib2;
    } else {
        simlibBusy = true;
    }
    worker.postMessage({msg: "simulate", id: id, payload: {
        iters: iters, dur: dur, gearlist: gearlist, opts: opts, rots: rots, haste: haste,
    }});
}

function makeid() {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < 16; i++ ) {
       result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

function computestats(gear, opts, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    simlib.postMessage({msg: "computeStats", id: id, payload: {gear: gear, opts: opts}});
}

function getOptions() {
    var options = {};


    options.buffai =  document.getElementById("buffai").checked;
    options.buffgotw =  document.getElementById("buffgotw").checked;
    options.buffbk =  document.getElementById("buffbk").checked;
    options.buffibow =  document.getElementById("buffibow").checked;
    options.buffmoon =  document.getElementById("buffmoon").checked;
    options.sbufws =  document.getElementById("sbufws").checked;
    options.debuffjow =  document.getElementById("debuffjow").checked;
    options.confbl =  document.getElementById("confbl").checked;
    options.confmr =  document.getElementById("confmr").checked;
    options.conbwo =  document.getElementById("conbwo").checked;
    options.conmm =  document.getElementById("conmm").checked;
    options.conbb =  document.getElementById("conbb").checked;
    options.consmp =  document.getElementById("consmp").checked;
    options.condr =  document.getElementById("condr").checked;
    options.totms =  document.getElementById("totms").checked;
    options.totwoa =  document.getElementById("totwoa").checked;

    options.buffbl =  parseInt(document.getElementById("buffbl").value) || 0;
    options.buffspriest = parseInt(document.getElementById("buffspriest").value) || 0;
    options.totwr =  parseInt(document.getElementById("totwr").value) || 0;
    options.buffdrum = 0; // todo, drums

    options.doopt = document.getElementById("doopt").checked;
    return options;
}

// Actually runs the sim. Uses the 'currentGear' global to populate the call.
function runsim() {
    var outele1 = document.getElementById("output1");
    var outele2 = document.getElementById("output2");

    var iters = parseInt(document.getElementById("iters").value);
    var dur = parseInt(document.getElementById("dur").value);

    var metricHTML = "<br /><div id=\"runningsim\" uk-spinner=\"ratio: 1.5\"></div><hr />";
    outele1.innerHTML = metricHTML;
    outele2.innerText = "";

    var gearlist = [];
    slotToID.forEach(k => {
        var item = currentGear[k];
        if (item != null && item.Name != "") {
            gearlist.push(item.Name);
        }
    });
    console.log("Options: ", getOptions());
    
    // #1 simulate LB
    // #2 simulate CL->LB priority cast
    // #3 if dur > LB.OOM run pure LB sim
    //    if du < CL->LB.OOM run priority sim.
    //    else, optimize sim

    var lbmetrics = null;
    var primetrics = null;
    var doExit = false;
    var includeFullDPS = true;

    var processSimResult = function(output) {
        console.log("Processing Results:", output)
        var optimal = {};
        var maxdps = 0.0;
        output.forEach(out => {
            var fulloutput = "";
            var total = out.TotalDmgs.reduce(function(sum, value){
                return sum + value;
            }, 0);
            var dps = total / out.SimSeconds;
            if (total/out.SimSeconds > maxdps) {
                maxdps = total/out.SimSeconds;
                optimal = out;
            }
    
            var values = out.TotalDmgs;
            var avg = average(values);
            var dev = standardDeviation(values, avg);
            var simdur = out.SimSeconds;
            if (out.Rotation[0] == "pri") {
                fulloutput += "Priority: " + out.Rotation.slice(1).join(", ") + "<br />";
                primetrics = out;
            } else if (out.Rotation.length > 6 && out.Rotation[0] == "CL6" && out.Rotation[1] == "LB12") {
                fulloutput += "Rotation: 1CL : " + (out.Rotation.length-1).toString() + "LB<br />";
                // fulloutput += "Rotation: " + out.Rotation.join(", ") + "<br />";
                if (out.Rotation.length == 1 && out.Rotation[0] == "LB12") {
                    lbmetrics = out;
                }
            } else {
                fulloutput += "Rotation: " + out.Rotation.join(", ") + "<br />";
                if (out.Rotation.length == 1 && out.Rotation[0] == "LB12") {
                    lbmetrics = out;
                }
            }
            var oomat = 0;
            var numOOM = out.OOMAt.reduce(function(sum, value){
                if (value > 0) {
                    oomat += value;
                    return sum + 1;
                }
                return sum;
            }, 0);

            if (includeFullDPS) {
                fulloutput += "Duration: " + out.SimSeconds + " seconds.<br />"
                fulloutput += "DPS: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";
            }
            if (numOOM > 0) {
                var values = out.DmgAtOOMs;
                var avg = average(values);
                var dev = standardDeviation(values, avg);
                var simdur = Math.round(oomat/numOOM);
                out.averageoom = Math.round(oomat/numOOM);

                if (includeFullDPS) {
                    fulloutput += "Went OOM: " + numOOM + " / " + iters + " simulations.<br />" 
                } else {
                    fulloutput += "Average time to OOM: " + Math.round(oomat/numOOM) + " seconds.<br />";
                    fulloutput += "DPS at time of OOM: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";        
                }
            } else {
                out.averageoom = 100000; // a big number
            }
            if (out.Rotation[0].startsWith("AI")) {
                var castStats = {
                    19: 0,
                    18: 0
                };
                out.Casts.forEach((casts)=>{
                    casts.forEach((cast)=>{
                        castStats[cast.ID] += 1
                    });
                });
                fulloutput += "Average Casts:<br />"
                fulloutput += " LB12: " + Math.round(castStats[18] / iters) + "<br />";
                fulloutput += " CL6: " + Math.round(castStats[19] / iters) + "<br />";
            }
            fulloutput += "<br />";
            console.log("Appending Full Output: ", fulloutput)
            outele2.innerHTML += fulloutput;
        });

        if (lbmetrics == null || primetrics == null) {
            return;
        }
        if (doExit) {
            outele1.innerHTML = "<hr />";
            return;
        }
        doExit = true;
        var realOpts = getOptions();
        if (!realOpts.doopt) {
            outele1.innerHTML = "<hr />";
            return;
        }
        outele2.innerHTML += "<hr /><p>Optimization Result:<br />";
        includeFullDPS = true;
        realOpts.useai = true;
        if (lbmetrics.averageoom < dur) {
            // set LB wins
            outele2.innerHTML += "-- You probably will need to downrank. -- <br />"
            simulate(iters, dur, gearlist, realOpts, [["LB12"]], 0, processSimResult);
        } else if (primetrics.averageoom > dur) {
            // set pri wins
            simulate(iters, dur, gearlist, realOpts, [["pri", "CL6","LB12"]], 0, processSimResult);
        } else {
            simulate(iters, dur, gearlist, realOpts, null, 0, processSimResult);
        }
    };

    var firstOpts = getOptions();
    if (firstOpts.doopt) {
        includeFullDPS = false;
        firstOpts.exitoom = true;
        firstOpts.doopt = false;
        simulate(iters, 600, gearlist, firstOpts, [["LB12"]], 0, processSimResult);
        simulate(iters, 600, gearlist, firstOpts, [["pri", "CL6","LB12"]], 0, processSimResult);    
    } else {
        simulate(iters, dur, gearlist, firstOpts, [["LB12"]], 0, processSimResult);
        simulate(iters, dur, gearlist, firstOpts, [["pri", "CL6","LB12"]], 0, processSimResult);    
    }
}

function hastedRotations() {
    console.log("Starting hasted rotations...");
    var gearlist = [];
    slotToID.forEach(k => {
        var item = currentGear[k];
        if (item != null && item.Name != "") {
            gearlist.push(item.Name);
        }
    });
    var opts = getOptions();
    opts.buffbl = 0;
    opts.buffdrum = 0;

    var hastes = [100, 200, 300, 400, 500, 600, 700, 788];
    var rots = [
        ["CL6", "LB12", "LB12", "LB12", "LB12"],
        ["CL6", "LB12", "LB12", "LB12", "LB12", "LB12"],
        ["CL6", "LB12", "LB12", "LB12", "LB12", "LB12", "LB12"]
    ];


    var hasteCounter = 0;
    hastes.forEach( haste => {
        hasteCounter++;
        var myCounter = hasteCounter;
        simulate(300, 60, gearlist, opts, rots, haste, (output) => {
            var maxdmg = 0.0;
            var maxrot = {};
    
            output.forEach(out => {
                var total = out.TotalDmgs.reduce(function(sum, value){
                    return sum + value;
                }, 0);
                if (total > maxdmg) {
                    maxrot = out;
                    maxdmg = total;
                }
            });
            
            var values = maxrot.TotalDmgs;
            var avg = average(values);
            var dev = standardDeviation(values, avg);
            var simdur = maxrot.SimSeconds;
            var rotTitle = "CL / " + (maxrot.Rotation.length-1).toString() + "xLB";
            var rows = document.getElementById("hasterots").firstElementChild.firstElementChild.children;
            var row = rows[myCounter];
            row.children[0].innerText = haste;
            row.children[1].innerText = rotTitle;
            row.children[2].innerText = "" + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur);
        });
    });
}

function standardDeviation(values, avg){
    var squareDiffs = values.map(function(value){
        var diff = value - avg;
        var sqrDiff = diff * diff;
        return sqrDiff;
    });

    var avgSquareDiff = average(squareDiffs);
    var stdDev = Math.sqrt(avgSquareDiff);
    return stdDev;
}
  
function average(data){
    var sum = data.reduce(function(sum, value){
        return sum + value;
    }, 0);

    var avg = sum / data.length;
    return avg;
}


window.addEventListener('click', function(e){   
    var allSelectors = document.getElementsByClassName("equipselector");
    for (var i = 0; i < allSelectors.length; i++) {
        var sel = allSelectors[i];
        if (sel.contains(e.target)){
            // Clicked in box
        } else{
            clearSearchEle(sel);
        }    
    }

});

window.addEventListener("keyup", (event) => {
    if (event.code == "Escape") {
        var allSelectors = document.getElementsByClassName("equipselector");
        for (var i = 0; i < allSelectors.length; i++) {
            clearSearchEle(allSelectors[i]);
        }
    }
});

// popgear will populate the allgear map from sim.
// Additionally it creates all the DOM elements for selecting gear.
function popgear(gearList) {
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
        var search = document.createElement("input");
        search.id = slot+"search";
        search.addEventListener("keyup", searchHandler);
        
        var clearbutton = document.createElement("button");
        clearbutton.innerText = "Clear";
        clearbutton.addEventListener("click", removeGear);
    
        var selectordiv = document.createElement("div");
        selectordiv.id = slot+"selector";
        selectordiv.classList.add("equipselector");
        selectordiv.style.display = "none";
        selectordiv.appendChild(search);
        selectordiv.appendChild(clearbutton);
        var selectorlist = document.createElement("div");
        selectorlist.id = slot+"selectorlist";
        selectordiv.appendChild(selectorlist);
    
        var name = document.createElement("p");
        name.addEventListener("click", focusSearch);
        name.id = slot+"label";
        name.innerText = "None";
        var gemdiv = document.createElement("div");
        gemdiv.id = slot+"enchants";
    
        var innerdiv = document.createElement("div");
        innerdiv.classList.add("equiplabel");
        innerdiv.appendChild(name);
        innerdiv.appendChild(gemdiv);
    
        var img = document.createElement("img");
        img.id = slot+"icon";
        img.addEventListener("click", focusSearch);
        img.src = "/icons/Items/Temp.png"
        var maindiv = document.createElement("div");
        maindiv.id = slot;
        maindiv.classList.add("equipslot");
        maindiv.appendChild(img);
        maindiv.appendChild(innerdiv);
        maindiv.appendChild(selectordiv);
        maindiv.addEventListener("blur", ()=>{maindiv.style.display = "none";});
        holder.appendChild(maindiv);
    });

    console.log("Items: ", gearList);
    gearList.Items.forEach(g => {
        if (g.GemSlots != null) {
            var gemslots = [];
            var gbytes = atob(g.GemSlots);
            for (var i = 0; i < gbytes.length; i++) {
                gemslots.push(gbytes.codePointAt(i));
            }
            g.GemSlots = gemslots;
        }
        allgear[g.Name] = g;


        try {
            var listItem = document.createElement("div");
            listItem.classList.add("equipselitem");
            listItem.innerText = g.Name;
            listItem.addEventListener("click", gearClickHandler);

            slotid = slotToID[g.Slot];
            if (slotid == "equipfinger" || slotid == "equiptrinket") {
                var itemlist = document.getElementById(slotid+"1selectorlist");
                itemlist.appendChild(listItem);

                var listItem2 = document.createElement("li");
                listItem2.classList.add("equipselitem");
                listItem2.innerText = g.Name;
                listItem2.addEventListener("click", gearClickHandler);
                var itemlist2 = document.getElementById(slotid+"2selectorlist");
                itemlist2.appendChild(listItem2);
            } else {
                var itemlist = document.getElementById(slotid+"selectorlist");
                itemlist.appendChild(listItem);
            }
        } catch (e) {
            console.log("Failed to intialize lootz: ", e);
        }
    });

    var finger1done = false;
    var trink1done = false;
    var glist = defaultGear;
    // TODO: make this store in like local storage or something so people cache gear choices.
    var gearCache = localStorage.getItem('cachedGear');
    if (gearCache && gearCache.length > 0) {
        var parsedGear = JSON.parse(gearCache);
        if (parsedGear.length > 0) {
            glist = parsedGear;
        }
    }

    glist.forEach(inm => {
        if (inm == "" || inm == "None") {
            return;
        }
        var item = allgear[inm];
        if (item == null) {
            return;
        }
        var slotid = slotToID[item.Slot];

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
        currentGear[slotid] = item;
    })

    updateGear(currentGear);
}

// Click handler for each item in slot list dropdown.
function gearClickHandler(event) {
    console.log("Gear Clicked: ", event);

    var slotid = event.target.parentElement.parentElement.parentElement.id;
    currentGear[slotid] = allgear[event.target.innerText];
    updateGear(currentGear);

    clearSearchEle(event.target.parentElement.parentElement);
}

// Click handler for 'remove' button on each slot.
function removeGear(ele) {
    var slotid = ele.parentElement.parentElement.id;
    console.log("Remove Slot: ", slotid);
    currentGear[slotid] = {name: "None"};
    updateGear(currentGear);

    var $dropdown = UIkit.dropdown(ele.parentElement);
    $dropdown.hide(0);
}

// For now hardcode an icon.
var slotToIcon = {
    "equiphead": "/icons/Armor/INV_Helmet_06.png",
    "equipneck": "/icons/Armor/INV_Jewelry_Necklace_07.png",
    "equipshoulder": "/icons/Armor/INV_Shoulder_14.png",
    "equipback": "/icons/Armor/INV_Misc_Cape_16.png",
    "equipchest": "/icons/Armor/INV_Chest_Chain_04.png",
    "equipwrist": "/icons/Armor/INV_Bracer_09.png",
    "equiphands": "/icons/Armor/INV_Gauntlets_26.png",
    "equipwaist": "/icons/Armor/INV_Belt_19.png",
    "equiplegs": "/icons/Armor/INV_Pants_03.png",
    "equipfeet": "/icons/Armor/INV_Boots_Wolf.png",
    "equipfinger1": "/icons/Armor/INV_Jewelry_Ring_04.png",
    "equipfinger2": "/icons/Armor/INV_Jewelry_Ring_05.png",
    "equiptrinket1": "/icons/Armor/INV_Jewelry_Talisman_09.png",
    "equiptrinket2": "/icons/Armor/INV_Jewelry_Talisman_10.png",
    "equipweapon": "/icons/Weapons/INV_Sword_39.png",
    "equipoffhand": "/icons/Armor/INV_Shield_20.png",
    "equiptotem": "/icons/Spells/Spell_Nature_InvisibilityTotem.png"
}
// updateGear will update the gear UI elements (to redraw when new gear is selected)
function updateGear(newGear) {
    var gearlist = [];
    slotToID.forEach(k => {
        var item = newGear[k];
        if (item != null && item.Name != "") {
            var nameele = document.getElementById(k+"label");
            nameele.innerText = item.Name;
            gearlist.push(item.Name);

            var iconImg = document.getElementById(k+"icon");
            iconImg.src = slotToIcon[k];

            var gemdiv = document.getElementById(k+"enchants");
            gemdiv.innerHTML = "";
            if (item.GemSlots != null) {
                item.GemSlots.forEach((gem) => {
                    var color = "rgba(30, 30, 30";
                    if (gem == 2) {
                        color = "rgba(250, 30, 30";
                    } else if (gem == 3) {
                        color = "rgba(30, 30, 250";
                    } else if (gem == 4) {
                        color = "rgba(250, 250, 30";
                    }
                    gemdiv.innerHTML += '<div class="gemslot" style="background-color: ' + color + ', 0.3);border:1px solid ' + color + ', 0.8)"></div>';
                });    
            }
            gemdiv.innerHTML += '<div class="enchslot" style="float: right;"></div>';
        }
    });

    currentGear = newGear;
    localStorage.setItem("cachedGear", JSON.stringify(gearlist));
    
    computestats(gearlist, null, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById(key.toLowerCase());
            lab.innerText = value;
        }
    });

    var opts = getOptions();
    computestats(gearlist, opts, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById("f"+key.toLowerCase());
            lab.innerText = value;
        }    
    });
}

// turns a search string into a list of lowercase terms and if there is punctuation or not.
function getSearchTerms(text) {
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
function find(search, value) {
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

// completeSearch finds the highlighted element and assigns it to the current gear set.
// Then closes the search area.
function completeSearch(ele) {
    var numChild = ele.parentElement.children[2].childElementCount;
    for (var i = 0; i < numChild; i++) {
        var le = ele.parentElement.children[2].children[i];
        if (le.classList.contains("lisearch")) {
            var slotid = ele.parentElement.parentElement.id;
            currentGear[slotid] = allgear[le.innerText];
            updateGear(currentGear);
            break;
        }
    }

    clearSearchEle(ele.parentElement);
}

function arrow(ele, search, numChild) {
    var mark = false;
    var last;
    var found = false;

    return function(i) {
        var le = ele.parentElement.children[2].children[i];
        if (mark) { // find next correct search result.
            found = find(search, le.innerText);
            if (found) {
                le.classList.add("lisearch");
                last.classList.remove("lisearch");
                return true; // exit now.
            }
        } else {
            if (le.classList.contains("lisearch")) {
                mark = true;
                last = le;
            }
        }
        return false;
    }
}
function handleSearchDown(ele, search, numChild) {
    var handledir = arrow(ele, search, numChild);
    for (var i = 0; i < numChild; i++) {
        if (handledir(i)) {
            return;
        }
    }
}
function handleSearchUp(ele, search, numChild) {
    var handledir = arrow(ele, search, numChild);
    for (var i = numChild-1; i>=0; i--) {
        if (handledir(i)) {
            return;
        }
    }
}
// Uses text from element to find item slot list.
// Ignores case and punctuation unless punctuation is included in the search.
// Spaces in search are implicit 'and'
function searchHandler(event) {
    var ele = event.target;
    if (event.code == "Enter" && ele.value != "") { // Enter
        completeSearch(ele)
        return;
    }

    var search = getSearchTerms(ele.value);
    // Number of things to search
    var numChild = ele.parentElement.children[2].childElementCount;

    if (event.code == "ArrowUp") {
        handleSearchUp(ele, search, numChild);
        return
    } else if (event.code == "ArrowDown") {
        handleSearchDown(ele, search, numChild);
        return
    }
    var firstFound = false;
    for (var i = 0; i < numChild; i++) {
        var le = ele.parentElement.children[2].children[i];
        var found = find(search, le.innerText)
        // Show / Hide item
        if (found) {
            if (!firstFound) {
                le.classList.add("lisearch");
                firstFound = true;
            } else {
                le.classList.remove("lisearch");
            }
            le.style.removeProperty("color")
        } else {
            le.style.color = "#888888";
            le.classList.remove("lisearch");
        }
    }
}

function clearSearch(event) {
    clearSearchEle(event.target.parentElement);
}

function clearSearchEle(ele) {
    ele.children[0].value = "";
    var numChild = ele.children[2].childNodes.length;
    for (var i = 0; i < numChild; i++) {
        var le = ele.children[2].children[i];
        le.style.removeProperty("color")
        le.classList.remove("lisearch");
    }
    ele.style.display = "none";
}

// focuses the search box (useful for making UI better)
function focusSearch(event) {
    var ele = event.target;

    // Get the 
    var tb;
    if (ele.id.includes("label")) {
        tb = document.getElementById(ele.id.replace("label", "selector"));
    } else {
        tb = document.getElementById(ele.parentElement.id+"selector");;
    }
    tb.style.display = "block";
    tb.firstElementChild.focus();

    var allSelectors = document.getElementsByClassName("equipselector");
    for (var i = 0; i < allSelectors.length; i++) {
        if (allSelectors[i].id != tb.id) {
            clearSearchEle(allSelectors[i]);
        }
    }

    event.preventDefault();
    event.stopPropagation();
}

// function 