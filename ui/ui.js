
// Globals
var defaultGear = [
    {Name:"Shamanistic Helmet of Second Sight"},
    {Name:"Brooch of Heightened Potential"},
    {Name:"Pauldrons of Wild Magic"},
    {Name:"Ogre Slayer's Cover"},
    {Name:"Tidefury Chestpiece"},
    {Name:"World's End Bracers"},
    {Name:"Earth Mantle Handwraps"},
    {Name:"Wave-Song Girdle"},
    {Name:"Stormsong Kilt"},
    {Name:"Magma Plume Boots"},
    {Name:"Cobalt Band of Tyrigosa"},
    {Name:"Scintillating Coral Band"},
    {Name:"Totem of the Void"},
    {Name:"Mazthoril Honor Shield"},
    {Name:"Bleeding Hollow Warhammer"},
    {Name:"Quagmirran's Eye"}, 
    {Name:"Icon of the Silver Crescent"}
];


// This code is all for interacting with the workers.
var simlib = new window.Worker(`simworker.js`);
var simlib2 = new window.Worker(`simworker.js`);

var simlibBusy = false;

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

var worker2 =  false;
function statweight(iters, dur, gearlist, opts, statToMod, modAmount, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    var worker = simlib;
    if (worker2) {
        worker = simlib2;
        worker2 = false;
    } else {
        worker2 = true;
    }
    worker.postMessage({msg: "statweight", id: id, payload: {
        iters: iters, dur: dur, gearlist: gearlist, opts: opts, stat: statToMod, modVal: modAmount
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

function computeStats(gear, opts, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    simlib.postMessage({msg: "computeStats", id: id, payload: {gear: gear, opts: opts}});
}

// Pulls options from the input 'options' pane for use in sims.
function getOptions() {
    var options = {};


    options.buffai =  document.getElementById("buffai").checked;
    options.buffgotw =  document.getElementById("buffgotw").checked;
    options.buffbk =  document.getElementById("buffbk").checked;
    options.buffibow =  document.getElementById("buffibow").checked;
    options.buffids =  document.getElementById("buffids").checked;
    options.buffmoon =  document.getElementById("buffmoon").checked;
    options.sbufws =  document.getElementById("sbufws").checked;
    options.debuffjow =  document.getElementById("debuffjow").checked;
    options.debuffmis =  document.getElementById("debuffmis").checked;
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
    options.buffdrums = parseInt(document.getElementById("buffdrums").value) || 0;

    return options;
}

var castIDToName = {
    1: "LB",
    2: "CL",
    3: "TLC LB",
    999: "LB Overload",
    998: "CL Overload",
}

// processSimResult will take the output from a numer of sims
// and process stuff like DPS and std dev.
function processSimResult(output) {
    var resultStats = {};
    var out = output[0];
    var maxDPS = 0;

    var dpsHist = {};
    var oomDPSHist = {};

    var total = out.TotalDmgs.reduce(function(sum, value){
        var dps = value/out.SimSeconds;
        var dpsRounded = Math.round(dps/10) * 10;
        if (dpsHist[dpsRounded] == null) {
            dpsHist[dpsRounded] = 0;
        }
        dpsHist[dpsRounded] += 1;
        if (dps > maxDPS) {
            maxDPS = dps;
        }
        return sum + value;
    }, 0);
    var dps = total / out.SimSeconds / out.TotalDmgs.length;        
    var values = out.TotalDmgs;
    var avg = average(values);
    var dev = standardDeviation(values, avg) / out.SimSeconds;
    
    var oomat = 0;
    var numOOM = out.OOMAt.reduce(function(sum, value){
        if (value > 0) {
            oomat += value;
            return sum + 1;
        }
        return sum;
    }, 0);
    oomat /= (numOOM);

    var dpsAtOOM = 0;
    if (numOOM > 0) {
        out.DmgAtOOMs.forEach((v, i) => {
            dpsAtOOM += v / out.OOMAt[i];
        });
        dpsAtOOM /= numOOM;
    }

    var castStats = {
        1: 0, // TODO: expose these constants from the wasm somehow.
        2: 0,
        3: 0,
        999: 0,
        998: 0,
    };
    out.Casts.forEach((casts)=>{
        casts.forEach((cast)=>{
            if (!cast.IsLO)  {
                castStats[cast.ID] += 1
            } else {
                castStats[1000-cast.ID] += 1
            }
            
        });
    });
  
    resultStats.total = total;
    resultStats.maxDPS = maxDPS;
    resultStats.dps = dps;
    resultStats.dev = dev;
    resultStats.oomat = oomat;
    resultStats.numOOM = numOOM;
    resultStats.dpsAtOOM = dpsAtOOM;
    resultStats.casts = castStats;
    resultStats.dpsHist = dpsHist;

    return resultStats;
}

// Populates the 'Sim' tab in the results pane.
function runsim(currentGear) {

    var iters = parseInt(document.getElementById("iters").value);
    var dur = parseInt(document.getElementById("dur").value);
    console.log("Options: ", getOptions());

    var lbout = document.getElementById("simrotlb");
    var priout = document.getElementById("simrotpri");
    var aiout = document.getElementById("simrotai");

    var metricHTML = `<div id="runningsim" uk-spinner="ratio: 1.5" style="margin:26%"></div>`;

    lbout.innerHTML = metricHTML;
    priout.innerHTML = metricHTML;
    aiout.innerHTML = metricHTML;


    var veryMax = 0.0;

    var firstOpts = getOptions();
    firstOpts.exitoom = true;

    simulate(iters, 600, currentGear, firstOpts, [["pri", "CL6","LB12"]], 0, (out) => { 
        var stats = processSimResult(out);
        var max = stats.dps;
        if (stats.dpsAtOOM > max) {
            console.log(`DPS: ${stats.dps}, OOM DPS: ${stats.dpsAtOOM}`);
            max = stats.dpsAtOOM;
        }
        if (max > veryMax) {
            veryMax = max;
        } else {
            max = veryMax;
        }
        priout.innerHTML = `<div><h3>Peak</h3><text class="simnums">${Math.round(max)}</text> dps<br /><text style="font-size:0.7em">${Math.round(stats.oomat)}s to oom at peak dps.</text></div>`
    });
    simulate(iters, 600, currentGear, firstOpts, [["LB12"]], 0, (out) => {
        var stats = processSimResult(out);
        var ttoom = stats.oomat;
        if (ttoom == 0) {
            ttoom = "Never";
        }
        lbout.innerHTML = `<div><h3>Mana</h3><text class="simnums">${Math.round(ttoom)}</text> sec<br /><text style="font-size:0.7em">to oom casting LB only</text></div>`
    });

    var secondOpts = getOptions();
    secondOpts.useai = true;
    secondOpts.doopt = true;
    simulate(iters, dur, currentGear, secondOpts, null, 0, (out) => { 
        var stats = processSimResult(out);
        console.log("AI Casts: ", stats.casts);
        aiout.innerHTML = `<div><h3>Average</h3><text class="simnums">${Math.round(stats.dps)}</text> +/- ${Math.round(stats.dev)} dps<br /></div>`
        
        var rotstats = document.getElementById("rotstats");
        rotstats.innerHTML = "";
        Object.entries(stats.casts).forEach((entry) => {
            if (entry[1] == 0) {
                return;
            }
            rotstats.innerHTML += `<text>${castIDToName[entry[0]]}: ${Math.round(entry[1]/iters)}</text>`;
        });

        var chartcanvas = document.createElement("canvas"); // `<canvas id="myChart" width="600" height="400"></canvas>`;
        var rotout = document.getElementById("rotout");
        var bounds = rotout.getBoundingClientRect();

        // Dirty hack in case the prio casting runs out of mana after BL is done but the average case has enough mana to burn priority casting.
        if (stats.dps > veryMax) {
            if (priout.childNodes.length > 0) {
                priout.childNodes[0].childNodes[1].innerText = Math.round(stats.dps);
            }
            veryMax = stats.dps;
        }

        var rotchart = document.getElementById("rotchart");
        rotchart.innerHTML = "";
        chartcanvas.height = bounds.height - 30;
        chartcanvas.width = bounds.width;
        var ctx = chartcanvas.getContext('2d');
        rotchart.appendChild(chartcanvas);

        var min = stats.dps - stats.dev;
        var max = stats.dps + stats.dev;
        var labels = Object.keys(stats.dpsHist)
        var vals = [];
        var devvals = [];

        labels.forEach((k, i) => {
            var val = parseInt(k);
            if (val > min && val < max) {
                devvals.push(stats.dpsHist[k]);
                vals.push(0);
            } else {
                vals.push(stats.dpsHist[k]);
                devvals.push(0);
            }
            labels[i] += " DPS";
        });
        var myChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: "DPS",
                    data: vals,
                    backgroundColor: [
                        '#1e87f0',
                    ],
                },
                {
                    label: "Expected DPS",
                    data: devvals,
                    backgroundColor: [
                        '#FF6961',
                    ],
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

    });
}

// Populates the 'Hasted Rotations' tab in results pane.
function hastedRotations(currentGear) {
    console.log("Starting hasted rotations...");
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
    var rows = document.getElementById("hasterots").firstElementChild.firstElementChild.children;
    hastes.forEach( haste => {
        hasteCounter++;
        var myCounter = hasteCounter;
        var row = rows[myCounter];
        row.children[1].innerHTML = "<div uk-spinner=\"ratio: 0.5\"></div>";
        row.children[2].innerText = "";

        simulate(800, 40, currentGear, opts, rots, haste, (output) => {
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

// Populates the 'Gear & Stat Weights' tab in results pane.
function calcStatWeights(gear) {
    var iters = parseInt(document.getElementById("switer").value);
    var dur = parseInt(document.getElementById("swdur").value);
    var opts = getOptions();
    
    var baseDPS = 0.0;
    var sp_hitModDPS = 0.0
    var modDPS = [0, 0, 0, 0, 0, 0];
    var weights = [0, 0, 0, 0, 0, 0, 0]; // Int, X, Crit, Hit, Dmg, Haste, MP5
    modDPS.forEach((v, i)=>{
        var cell = document.getElementById("w"+i.toString());
        cell.innerHTML = "<div uk-spinner=\"ratio: 1\"></div>";
    });

    // A base DPS without any modified stats.
    statweight(iters, dur, gear, opts, 0, 0, (res) => {
        baseDPS = res;
    }); // base


    var done = [];
    var onfinish = () => {
        done.push(true);
        if (baseDPS == 0) {
            return;
        }
        if (modDPS[0] == 0) {
            return;
        }
        var ddps = modDPS[0] - baseDPS;

        modDPS.forEach((v, i)=>{
            if (v == 0) {
                return;
            }
            var cell = document.getElementById("w"+i.toString());
            // sphit uses different value;
            if (i == 3 && sp_hitModDPS != 0.0) {
                var sphitdiff = sp_hitModDPS - baseDPS;
                var weight = (v-baseDPS) / sphitdiff;
                if (weight < 0.01) {
                    weight = 0.0;
                }
                weights[3] = weight;
                cell.innerText = weight.toFixed(2);
                return; 
            }
            var weight = (v-baseDPS) / ddps;
            if (weight < 0.01) {
                weight = 0.0;
            }
            
            if (i == 0) {
                weights[4] = weight;
            } else if (i == 1) {
                weights[0] = weight;
            } else if (i == 2) {
                weights[2] = weight;
            } else if (i == 4) {
                weights[5] = weight;
            } else if (i == 5) {
                weights[6] = weight;
            }
            cell.innerText = weight.toFixed(2);
        });

        if (done.length == 7) {
            showGearRecommendations(weights);
        }
    };


    statweight(iters, dur, gear, opts, 4, 25, (res) => {sp_hitModDPS = res;onfinish();}); // sp
    statweight(iters, dur, gear, opts, 3, 25, (res) => {modDPS[3] = res;onfinish();}); // hit

    statweight(iters, dur, gear, opts, 4, 100, (res) => {modDPS[0] = res;onfinish();}); // sp
    statweight(iters, dur, gear, opts, 0, 100, (res) => {modDPS[1] = res;onfinish();}); // int
    statweight(iters, dur, gear, opts, 2, 100, (res) => {modDPS[2] = res;onfinish();}); // crit
    statweight(iters, dur, gear, opts, 5, 100, (res) => {modDPS[4] = res;onfinish();}); // haste
    statweight(iters, dur, gear, opts, 6, 100, (res) => {modDPS[5] = res;onfinish();}); // mp5


}

function showGearRecommendations(weights) {
    var itemWeightsBySlot = {};
    var curSlotWeights = {};
    var csdVal = (((currentFinalStats["StatSpellDmg"]*0.795)+603)*2 * (currentFinalStats["StatSpellCrit"]/2208) * 0.045) / 0.795;
    // process all items to find the weighted value.
    // find the value of each slots currently equipped item.
    Object.entries(gearUI.allitems).forEach((entry) => {
        var name = entry[0];
        var item = entry[1];

        var value = 0.0;
        if (item.Stats != null) {
            weights.forEach((w, i)=>{
                if (item.Stats[i] != null) {
                    value += item.Stats[i]*w
                }
            });
        }
        if (itemWeightsBySlot[item.Slot] == null) {
            itemWeightsBySlot[item.Slot] = [];
        }
        if (item.GemSlots != null && item.GemSlots.length > 0) {
            var numGems = item.GemSlots.length;
            if (item.GemSlots[0] == 1) {
                numGems--;
                // how to value a CSD
                // ~ spellpower * crit chance * 0.09 = increased damage per cast.
                value += csdVal;
            }
            value += (numGems * 9) * weights[4]; // just for measuring use 9 spell power gems in every slot.
        }
        var curEquip = gearUI.currentGear[slotToID[item.Slot]];
        if (curEquip != null && curEquip.Name == item.Name) {
            curSlotWeights[item.Slot] = value;
        }
        itemWeightsBySlot[item.Slot].push({Name: item.Name, Weight: value});
        itemWeightsBySlot[item.Slot] = itemWeightsBySlot[item.Slot].sort((a,b)=> b.Weight - a.Weight);
    });
    var uptab = document.getElementById("upgrades");
    uptab.innerHTML = "";
    
    var curSlot = -1;
    Object.entries(itemWeightsBySlot).forEach((entry)=>{
        if (entry[0] == 11 || entry[0] == 14) {
            // Skip rings/trinkets for now. Trinkets will be separate and rings need 
            // to check a finger1/2 instead of finger generically.
            return;
        }
        if (curSlot != entry[0]) {
            var row = document.createElement("tr");
            var col1 = document.createElement("td");
            slotToID[entry[0]].replace("equip", "");
            var title = slotToID[entry[0]].replace("equip", "");
            col1.innerHTML = "<h3>" +title.charAt(0).toUpperCase() + title.substr(1)+"</h3>";

            var col2 = document.createElement("td");
            var col3 = document.createElement("td");
            row.appendChild(col1);
            row.appendChild(col2);
            row.appendChild(col3);
            uptab.appendChild(row);
            curSlot = entry[0];
        }
        // get current item slot.
        var alt = 0;
        entry[1].forEach((v) => {
            alt++;
            var row = document.createElement("tr");
            if (alt%2 == 0) {
                row.style.backgroundColor = "#808080";
            }
            var col1 = document.createElement("td");
            col1.innerText = v.Name;
            var col2 = document.createElement("td");
            col2.innerText = Math.round(v.Weight - curSlotWeights[curSlot]);
            var col3 = document.createElement("td");
            var col4 = document.createElement("td");
            var simbut = document.createElement("button");
            simbut.innerText = "Sim";

            var item = Object.assign({Name: ""}, gearUI.allitems[v.Name]);
            simbut.addEventListener("click", (e)=>{
                col4.innerHTML = "<div uk-spinner=\"ratio: 0.5\"></div>";
                var newgear = {};
                Object.entries(gearUI.currentGear).forEach((entry)=>{
                    if (entry[0] == slotToID[item.Slot]) {
                        // replace
                        newgear[entry[0]] = item;
                        if (item.GemSlots != null) {
                            item.Gems = [];
                            item.GemSlots.forEach((color, i) => {             
                                if (color == 1) {
                                    item.Gems[i] = gearUI.allgems["Chaotic Skyfire Diamond"];
                                } else {
                                    item.Gems[i] = gearUI.allgems["Runed Living Ruby"];
                                }                                    
                            });    
                        }
                    } else {
                        newgear[entry[0]] = entry[1];
                    }
                });
                var iters = parseInt(document.getElementById("switer").value);
                var dur = parseInt(document.getElementById("swdur").value);
                var opts = getOptions();
                simulate(iters, dur, cleanGear(newgear), opts, null, null, (res)=>{
                    var statistics = processSimResult(res);
                    col4.innerText = Math.round(statistics.dps).toString() + " +/- " + Math.round(statistics.dev).toString();
                });
            });
            col3.appendChild(simbut);
            row.appendChild(col1);
            row.appendChild(col2);
            row.appendChild(col3);
            row.appendChild(col4);
            uptab.appendChild(row);
        })
    });
}

window.addEventListener('click', function(e){   
    if (gearUI == null) {
        return;
    }
    gearUI.hideSelectors(e);
});

window.addEventListener("keyup", (event) => {
    if (event.code == "Escape") {
        gearUI.hideSelectors();
    }
});

var gearUI;

// popgear will populate the allgear map from sim.
// Additionally it creates all the DOM elements for selecting gear.
function popgear(gearList) {
    gearUI = new GearUI(document.getElementById('gear'), gearList);

    var glist = defaultGear;
    
    var gearCache = localStorage.getItem('cachedGear.v2');
    if (gearCache && gearCache.length > 0) {
        var parsedGear = JSON.parse(gearCache);
        if (parsedGear.length > 0) {
            glist = parsedGear;
        }
    }
    console.log("Gear: ", glist);
    var currentGear = gearUI.updateEquipped(glist);

    gearUI.addChangeListener((item, slot)=>{
        updateGearStats(gearUI.currentGear);
    });
    updateGearStats(currentGear)

    var simrunbut = document.getElementById("simrunbut");
    simrunbut.addEventListener("click", (event)=>{
        runsim(cleanGear(gearUI.currentGear));
    });

    var hastebut = document.getElementById("hastebut");
    hastebut.addEventListener("click", (event)=>{
        hastedRotations(cleanGear(gearUI.currentGear));
    });

    var caclweights = document.getElementById("calcstatweight");
    caclweights.addEventListener("click", (event)=>{
        calcStatWeights(cleanGear(gearUI.currentGear));
    });
}

// clearGear strips off all parts of gear that is non-changing. This lets us pass minimal data to sim and store in local storage.
function cleanGear(gear) {
    var cleanedGear = [];
    Object.entries(gear).forEach((entry)=>{
        if (entry == null || entry == undefined) {
            return;
        }
        if (entry[1] == null || entry[1] == undefined) {
            return;
        }
        var it = {
            Name: entry[1].Name,
            Gems: [],
            // TODO: enchants
        }
        if (entry[1].Gems != null) {
            entry[1].Gems.forEach((g)=>{
                if (g == null) {
                    it.Gems.push("");
                    return;
                }
                it.Gems.push(g.Name);
            });    
        }
        cleanedGear.push(it);
    });
    return cleanedGear
}

var currentFinalStats = {};

// Updates the 'stats' pane in the viewport.
function updateGearStats(gearlist) {
    
    var cleanedGear = cleanGear(gearlist); // converts to array with minimal data for serialization.

    // TODO: Is this the best way?
    localStorage.setItem("cachedGear.v2", JSON.stringify(cleanedGear));

    computeStats(cleanedGear, null, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById(key.toLowerCase());
            if (lab != null) {
                lab.innerText = value;
            }
        }
    });

    var opts = getOptions();
    computeStats(cleanedGear, opts, (result) => {
        currentFinalStats = result;
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById("f"+key.toLowerCase());
            if (key.toLowerCase() == "statspellcrit") {
                lab.innerText = value.toString() + " ("  + (value/22.08).toFixed(1) + "%)";
            } else if (key.toLowerCase() == "statspellhit") {
                lab.innerText = value.toString() + " ("  + (value/12.6).toFixed(1) + "%)";
            } else if (key.toLowerCase() == "statspellhaste") {
                lab.innerText = value.toString() + " ("  + (value/15.76).toFixed(1) + "%)";
            } else if (lab == null) {
                // do nothing...
            } else {
                lab.innerText = value;
            }
        }    
    });
}


/// I hate html and javascript so much sometimes.

var panedrag = false;
var panediv = document.getElementById("panediv");
var calcpane = document.getElementById("calctabs");
var inpanel = document.getElementById("inputdata");
var outpanel = document.getElementById("calcdiv");
var h = window.innerHeight;

window.addEventListener("touchstart", (e)=>{
    if (panediv == null) {
        panediv = document.getElementById("panediv");
        calcpane = document.getElementById("calctabs");
        inpanel = document.getElementById("inputdata");
        outpanel = document.getElementById("calcdiv");
    }
    if (panediv.contains(e.target) || calcpane.contains(e.target)) {
        console.log("now dragging...");
        panedrag = true;
        h = window.innerHeight;
        e.preventDefault();
    }
    if( e.changedTouches.length > 1) {
        panedrag = false;
    }
});
window.addEventListener("mousedown", (e)=>{
    if (panediv == null) {
        panediv = document.getElementById("panediv");
        inpanel = document.getElementById("inputdata");
        outpanel = document.getElementById("calcdiv");
        calcpane = document.getElementById("calctabs");
    }
    if (panediv.contains(e.target) || calcpane.contains(e.target)) {
        console.log("now mouse dragging...");
        panedrag = true;
        h = window.innerHeight;
        e.preventDefault();
    }
});
window.addEventListener("touchend", (e)=>{
    panedrag = false;
});
window.addEventListener("mouseup", (e)=>{
    panedrag = false;
});

window.addEventListener("touchmove", (e)=>{
    if (panedrag) {
        console.log("dragging...", e)
        var percent = e.changedTouches[0].pageY / (h+60);
        inpanel.style.height = "calc(" + (percent*100).toString() + "% - 100px)";
        outpanel.style.height = "calc(" + (100-(percent*100)).toString() + "% - 100px)";
        e.preventDefault();
    }
});

window.addEventListener("mousemove", (e)=>{
    if (panedrag) {
        var percent = e.pageY / (h+60);
        inpanel.style.height = "calc(" + (percent*100).toString() + "% - 100px)";
        outpanel.style.height = "calc(" + (100-(percent*100)).toString() + "% - 100px)";
        e.preventDefault();
    }
});


var theme = "dark";
function toggletheme() {
    if (theme == "light") {
        document.getElementById("themebulb").src = "../icons/light-bulb.svg";
        document.children[0].children[1].classList.remove("lighttheme")
        document.children[0].classList.remove("lighttheme")
        
        document.children[0].classList.add("darktheme")
        document.children[0].children[1].classList.add("darktheme")
        document.children[0].children[1].classList.add("uk-light")

        toggleThemeClass("ltl", "dtl");
        toggleThemeClass("ltm", "dtm");
        toggleThemeClass("ltd", "dtd");

        theme = "dark";
    } else {
        document.getElementById("themebulb").src = "../icons/lightbulb.svg";

        document.children[0].children[1].classList.remove("uk-light") 
        document.children[0].classList.remove("darktheme")
        document.children[0].children[1].classList.remove("darktheme")

        document.children[0].children[1].classList.add("lighttheme")
        document.children[0].classList.add("lighttheme")

        toggleThemeClass("dtl", "ltl");
        toggleThemeClass("dtm", "ltm");
        toggleThemeClass("dtd", "ltd");

        theme = "light";
    }
}

function toggleThemeClass(rm, rp) {
    var elements = document.getElementsByClassName(rm);
    for (var i = elements.length-1; i >=0 ; i--) {
        var e = elements[i];
        e.classList.remove(rm);
        e.classList.add(rp);
    }
}