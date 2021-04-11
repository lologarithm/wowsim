
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

// Actually runs the sim.
function runsim(currentGear) {
    var outele1 = document.getElementById("output1");
    var outele2 = document.getElementById("output2");

    var iters = parseInt(document.getElementById("iters").value);
    var dur = parseInt(document.getElementById("dur").value);

    var metricHTML = "<br /><div id=\"runningsim\" uk-spinner=\"ratio: 1.5\"></div><hr />";
    outele1.innerHTML = metricHTML;
    outele2.innerText = "";

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
            } else if (numOOM == 0) {
                fulloutput += "Average time to OOM: Never<br />";
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
                    fulloutput += "Average time to OOM: " + Math.round(oomat/numOOM) + " seconds.<br />";
                } else {
                    fulloutput += "Average time to OOM: " + Math.round(oomat/numOOM) + " seconds.<br />";
                    fulloutput += "DPS at time of OOM: " + Math.round(avg/simdur) + " +/- " + Math.round(dev/simdur) + "<br />";        
                }
            } else {
                out.averageoom = 100000; // a big number
            }
            if (out.Rotation[0].startsWith("AI")) {
                var castStats = {
                    1: 0, // TODO: expose these constants from the wasm somehow.
                    2: 0
                };
                out.Casts.forEach((casts)=>{
                    casts.forEach((cast)=>{
                        castStats[cast.ID] += 1
                    });
                });
                fulloutput += "Average Casts:<br />"
                fulloutput += " LB12: " + Math.round(castStats[1] / iters) + "<br />";
                fulloutput += " CL6: " + Math.round(castStats[2] / iters) + "<br />";
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
            simulate(iters, dur, currentGear, realOpts, [["LB12"]], 0, processSimResult);
        } else if (primetrics.averageoom > dur) {
            // set pri wins
            simulate(iters, dur, currentGear, realOpts, [["pri", "CL6","LB12"]], 0, processSimResult);
        } else {
            simulate(iters, dur, currentGear, realOpts, null, 0, processSimResult);
        }
    };

    var firstOpts = getOptions();
    if (firstOpts.doopt) {
        includeFullDPS = false;
        firstOpts.exitoom = true;
        firstOpts.doopt = false;
        simulate(iters, 600, currentGear, firstOpts, [["LB12"]], 0, processSimResult);
        simulate(iters, 600, currentGear, firstOpts, [["pri", "CL6","LB12"]], 0, processSimResult);    
    } else {
        simulate(iters, dur, currentGear, firstOpts, [["LB12"]], 0, processSimResult);
        simulate(iters, dur, currentGear, firstOpts, [["pri", "CL6","LB12"]], 0, processSimResult);    
    }
}

function hastedRotations(currentGear) {
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
        simulate(400, 30, gearlist, opts, rots, haste, (output) => {
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

function calcStatWeights(gear) {
    var iters = parseInt(document.getElementById("switer").value);
    var dur = parseInt(document.getElementById("swdur").value);
    var opts = getOptions();
    
    var baseDPS = 0.0;
    var sp_hitModDPS = 0.0
    var modDPS = [0, 0, 0, 0, 0, 0];

    modDPS.forEach((v, i)=>{
        var cell = document.getElementById("w"+i.toString());
        cell.innerHTML = "<div uk-spinner=\"ratio: 1\"></div>";
    });

    // A base DPS without any modified stats.
    statweight(iters, dur, gear, opts, 0, 0, (res) => {
        baseDPS = res;
    }); // base


    var onfinish = () => {
        if (baseDPS == 0) {
            return;
        }
        if (modDPS[0] == 0) {
            return;
        }
        var ddps = modDPS[0] - baseDPS;

        var complete = true;
        modDPS.forEach((v, i)=>{
            if (v == 0) {
                complete = false;
                return;
            }
            var cell = document.getElementById("w"+i.toString());
            // sphit uses different value;
            if (i == 3 && sp_hitModDPS != 0.0) {
                var sphitdiff = sp_hitModDPS - baseDPS;
                var weight = (v-baseDPS) / sphitdiff;
                cell.innerText = weight.toFixed(2);
                return; 
            }
            var weight = (v-baseDPS) / ddps;
            cell.innerText = weight.toFixed(2);
        });

        if (complete) {
            
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

function updateGearStats(gearlist) {
    
    var cleanedGear = cleanGear(gearlist); // converts to array with minimal data for serialization.

    // TODO: Is this the best way?
    localStorage.setItem("cachedGear.v2", JSON.stringify(cleanedGear));

    computeStats(cleanedGear, null, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById(key.toLowerCase());
            lab.innerText = value;
        }
    });

    var opts = getOptions();
    computeStats(cleanedGear, opts, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById("f"+key.toLowerCase());
            if (key.toLowerCase() == "statspellcrit") {
                lab.innerText = value.toString() + " ("  + (value/22.08).toFixed(1) + "%)";
            } else if (key.toLowerCase() == "statspellhit") {
                lab.innerText = value.toString() + " ("  + (value/12.6).toFixed(1) + "%)";
            } else if (key.toLowerCase() == "statspellhaste") {
                lab.innerText = value.toString() + " ("  + (value/15.76).toFixed(1) + "%)";
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
