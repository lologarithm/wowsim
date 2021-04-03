
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

function getOptions() {
    var options = {};


    options.buffai =  document.getElementById("buffai").value == "on";
    options.buffgotw =  document.getElementById("buffgotw").value == "on";
    options.buffbk =  document.getElementById("buffbk").value == "on";
    options.buffibow =  document.getElementById("buffibow").value == "on";
    options.buffmoon =  document.getElementById("buffmoon").value == "on";
    options.sbufws =  document.getElementById("sbufws").value == "on";
    options.debuffjow =  document.getElementById("debuffjow").value == "on";
    options.confbl =  document.getElementById("confbl").value == "on";
    options.confmr =  document.getElementById("confmr").value == "on";
    options.conbwo =  document.getElementById("conbwo").value == "on";
    options.conmm =  document.getElementById("conmm").value == "on";
    options.conbb =  document.getElementById("conbb").value == "on";
    options.consmp =  document.getElementById("consmp").value == "on";
    options.condr =  document.getElementById("condr").value == "on";
    options.totms =  document.getElementById("totms").value == "on";
    options.totwoa =  document.getElementById("totwoa").value == "on";

    options.buffbl =  parseInt(document.getElementById("buffbl").value);
    options.buffspriest =  parseInt(document.getElementById("buffspriest").value);
    options.totwr =  parseInt(document.getElementById("totwr").value);
    options.buffdrum = 0; // todo, drums

    return options;
}

// Actually runs the sim. Uses the 'currentGear' global to populate the call.
function runsim() {
    var iters = document.getElementById("iters").value;
    var dur = document.getElementById("dur").value;
    var gearlist = [];
    slotToID.forEach(k => {
        var item = currentGear[k];
        if (item != null && item.name != "") {
            gearlist.push(item.name);
        }
    });
    console.log("Options: ", getOptions());
    var resStr = simulate(parseInt(iters), parseInt(dur), gearlist, getOptions());
    var output = JSON.parse(resStr);
    console.log("Result:", output);
    
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
        } else {
            fulloutput += "Rotation: " + out.Rotation.join(", ") + "<br />";
        }
        fulloutput += "Duration: " + out.SimSeconds + " seconds.<br />"
        fulloutput += "DPS: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";

        var oomat = 0;
        var numOOM = out.OOMAt.reduce(function(sum, value){
            if (value > 0) {
                oomat += value;
                return sum + 1;
            }
            return sum;
        }, 0);
        if (numOOM > 0) {
            var values = out.DmgAtOOMs;
            var avg = average(values);
            var dev = standardDeviation(values, avg);
            var simdur = Math.round(oomat/numOOM);
    
            fulloutput += "Went OOM: " + numOOM + " / " + iters + " simulations. Average time when OOM: " + Math.round(oomat/numOOM) + " seconds.<br />";
            fulloutput += "DPS at time of OOM: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";    
        }
        fulloutput += "<br />";
    });

    var outele = document.getElementById("output");
    var statele = document.getElementById("fstats");
    
    var simdur = optimal.SimSeconds;
    var metricHTML = "<br /><hr />Optimal Casting Rotation for a " + simdur + " second fight -> ";
    if (optimal.Rotation.length == 1) {
        metricHTML += "<b>LB12 Spam</b>";
    } else if (optimal.Rotation[0] == "pri") {
        metricHTML += "<b>CL on CD</b>";
    } else {
        var numLB = optimal.Rotation.length - 1
        metricHTML += " <b>" + numLB + "LB : 1CL</b>";
    }
    var values = optimal.TotalDmgs;
    var avg = average(values);
    var dev = standardDeviation(values, avg);
    metricHTML += " @ " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + " DPS<hr /><br /><br />";

    metricHTML += fulloutput;
    outele.innerHTML = metricHTML;
    // statele.innerHTML = JSON.stringify(output.stats).replaceAll(",", "\n", );
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

    var fulloutput = "<table class=\"hastetable\"><tr><th>Haste</th><th>Rotation</th><th>DPS</th></tr>";

    var hastes = [100, 200, 300, 400, 500, 600, 700, 788];
    var rots = [
        ["CL6", "LB12", "LB12", "LB12", "LB12"],
        ["CL6", "LB12", "LB12", "LB12", "LB12", "LB12"],
        ["CL6", "LB12", "LB12", "LB12", "LB12", "LB12", "LB12"]
    ];

    hastes.forEach( haste => {
        console.log("Running haste: ", haste);
        var opts = getOptions();
        opts.buffbl = 0;
        opts.buffdrum = 0;

        var resStr = simulate(200, 60, gearlist, opts, rots, haste);
        var output = JSON.parse(resStr);
        
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
        fulloutput += "<tr><td>" + haste + "</td>"
        fulloutput += "<td>"
        fulloutput += rotTitle
        fulloutput += "</td>"
        fulloutput += "<td>"
        fulloutput += "" + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";
        fulloutput += "</td>"
    });

    fulloutput += "</table>"
    var outele = document.getElementById("hasterots");
    outele.innerHTML = fulloutput;
    // statele.innerHTML = JSON.stringify(output.stats).replaceAll(",", "\n", );
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
function popgear() {
    var geararray = JSON.parse(gearlist());
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

    // TODO: make this store in like local storage or something so people cache gear choices.
    var finger1done = false;
    var trink1done = false;
    defaultGear.forEach(inm => {
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
    slotToID.forEach(k => {
        var item = newGear[k];
        if (item != null && item.name != "") {
            var button = document.getElementById(k).firstElementChild;
            button.innerText = item.name;
        }
    });

    currentGear = newGear;
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
            le.style.color = "#F0F0F0";
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
        // le.style.display = "block";
        le.style.color = "#F0F0F0";
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