
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
    var onComp = simrequests[event.data.id];
    if (onComp != null) {
        onComp(event.data.payload);
    }
}

var simrequests = {};
function simulate(iters, dur, gearlist, opts, rots, haste, onComplete) {
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
// TODO: unique ID each request...

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
        if (item != null && item.name != "") {
            gearlist.push(item.name);
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
        var optimal = {};
        var maxdps = 0.0;
        var fulloutput = "";
        output.forEach(out => {
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
                    fulloutput += "Went OOM: " + numOOM + " / " + iters + " simulations. " 
                }
                fulloutput += "Average time when OOM: " + Math.round(oomat/numOOM) + " seconds.<br />";
                fulloutput += "DPS at time of OOM: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";    
            } else {
                out.averageoom = 100000; // a big number
            }
            fulloutput += "<br />";
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
        outele2.innerHTML += "<hr /><p>Optimized Rotation:<br />";
        includeFullDPS = true;
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
        if (item != null && item.name != "") {
            gearlist.push(item.name);
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

// popgear will populate the allgear map from sim.
function popgear(geararray) {
    geararray.forEach(g => {
        allgear[g.name] = g;

        try {
            var listItem = document.createElement("li");
            listItem.innerText = g.name;
            
            slotid = slotToID[g.slot];
            if (slotid == "equipfinger") {
                var nav = UIkit.nav("#equipfinger1 div.uk-dropdown ul.uk-nav.uk-dropdown-nav");
                nav.$el.addEventListener("click", gearClickHandler);
                nav.$el.appendChild(listItem);

                var listItem2 = document.createElement("li");
                listItem2.innerText = g.name;
                var nav = UIkit.nav("#equipfinger2 div.uk-dropdown ul.uk-nav.uk-dropdown-nav");
                nav.$el.addEventListener("click", gearClickHandler);
                nav.$el.appendChild(listItem2);
            } else if (slotid == "equiptrinket") {
                var nav = UIkit.nav("#equiptrinket1 div.uk-dropdown ul.uk-nav.uk-dropdown-nav");
                nav.$el.addEventListener("click", gearClickHandler);
                nav.$el.appendChild(listItem);
                
                var listItem2 = document.createElement("li");
                listItem2.innerText = g.name;
                var nav = UIkit.nav("#equiptrinket2 div.uk-dropdown ul.uk-nav.uk-dropdown-nav");
                nav.$el.addEventListener("click", gearClickHandler);
                nav.$el.appendChild(listItem2);
            } else {
                var nav = UIkit.nav("#" + slotid + " div.uk-dropdown ul.uk-nav.uk-dropdown-nav");
                nav.$el.addEventListener("click", gearClickHandler);
                nav.$el.appendChild(listItem);
            }
        } catch (e) {
            console.log("Failed to intialize lootz: ", e);
        }
    });

    var finger1done = false;
    var trink1done = false;
    var glist = [];
    // TODO: make this store in like local storage or something so people cache gear choices.
    var gearCache = localStorage.getItem('cachedGear');
    if (gearCache) {
        glist = JSON.parse(gearCache);
    } else {
        glist = defaultGear;
    }

    glist.forEach(inm => {
        var item = allgear[inm];
        var slotid = slotToID[item.slot];

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

    UIkit.update(element = document.body, type = 'update');
}

// Click handler for each item in slot list dropdown.
function gearClickHandler(event) {
    console.log("Gear Clicked: ", event);

    var slotid = event.target.parentElement.parentElement.parentElement.id;
    currentGear[slotid] = allgear[event.target.innerText];
    updateGear(currentGear);

    var $dropdown = UIkit.dropdown(event.target.parentElement.parentElement);
    $dropdown.hide(0);
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

// updateGear will update the gear UI elements (to redraw when new gear is selected)
function updateGear(newGear) {
    var gearlist = [];
    slotToID.forEach(k => {
        var item = newGear[k];
        if (item != null && item.name != "") {
            var button = document.getElementById(k).firstElementChild;
            button.innerText = item.name;
            gearlist.push(item.name);
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
    var $dropdown = UIkit.dropdown(ele.parentElement);
    $dropdown.hide(0);
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
function searchHandler(ele, event) {
    if (event.keyCode == 13 && ele.value != "") { // Enter
        completeSearch(ele)
        return;
    }

    var search = getSearchTerms(ele.value);
    // Number of things to search
    var numChild = ele.parentElement.children[2].childElementCount;

    if (event.keyCode == 38) {
        handleSearchUp(ele, search, numChild);
        return
    } else if (event.keyCode == 40) {
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

function clearSearch(e) {
    e.value = "";
    var numChild = e.parentElement.children[2].childElementCount;
    for (var i = 0; i < numChild; i++) {
        var le = e.parentElement.children[2].children[i];
        le.style.removeProperty("color")
        le.classList.remove("lisearch");
    }
}

// focuses the search box (useful for making UI better)
function focusSearch(ele, event) {
    var tb = ele.parentElement.children[1].children[0];
    setTimeout( () => {
        tb.focus();
    }, 10);
}

// function 