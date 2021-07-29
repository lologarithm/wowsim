// Globals
var defaultGear = [
    { Name: "Shamanistic Helmet of Second Sight" },
    { Name: "Brooch of Heightened Potential" },
    { Name: "Pauldrons of Wild Magic" },
    { Name: "Ogre Slayer's Cover" },
    { Name: "Tidefury Chestpiece" },
    { Name: "World's End Bracers" },
    { Name: "Earth Mantle Handwraps" },
    { Name: "Wave-Song Girdle" },
    { Name: "Stormsong Kilt" },
    { Name: "Magma Plume Boots" },
    { Name: "Cobalt Band of Tyrigosa" },
    { Name: "Scintillating Coral Band" },
    { Name: "Totem of the Void" },
    { Name: "Mazthoril Honor Shield" },
    { Name: "Bleeding Hollow Warhammer" },
    { Name: "Quagmirran's Eye" },
    { Name: "Icon of the Silver Crescent" }
];


// This code is all for interacting with the workers.
var simlib = new window.Worker(`simworker.js`);
var simlib2 = new window.Worker(`simworker.js`);

var simlibBusy = false;

simlib.onmessage = (event) => {
    var m = event.data.msg;
    if (m == "ready") {
        simlib.postMessage({ msg: "setID", payload: "1" });
        simlib.postMessage({ msg: "getGearList" });
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
        simlib2.postMessage({ msg: "setID", payload: "2" });
        return;
    }
    var onComp = simrequests[event.data.id];
    if (onComp != null) {
        onComp(event.data.payload);
    }
}

var simrequests = {};
function simulate(iters, dur, numClTargets, gearlist, opts, agentTypes, haste, fullLogs, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    var worker = simlib;
    if (simlibBusy) {
        worker = simlib2;
    } else {
        simlibBusy = true;
    }
    worker.postMessage({
        msg: "simulate", id: id, payload: {
            iters: iters,
            dur: dur,
            numClTargets: numClTargets,
            gearlist: gearlist,
            opts: opts,
            agentTypes: agentTypes,
            haste: haste,
            fullLogs: fullLogs
        }
    });
}

var worker2 = false;
function statweight(iters, dur, numClTargets, gearlist, opts, statToMod, modAmount, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    var worker = simlib;
    if (worker2) {
        worker = simlib2;
        worker2 = false;
    } else {
        worker2 = true;
    }
    worker.postMessage({
        msg: "statweight", id: id, payload: {
            iters: iters,
            dur: dur,
            numClTargets: numClTargets,
            gearlist: gearlist,
            opts: opts,
            stat: statToMod,
            modVal: modAmount
        }
    });
}

function packOptions(opt, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    simlib.postMessage({
        msg: "packopt", id: id, payload: {
            opt: opt,
        }
    });
}

function makeid() {
    var result = '';
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < 16; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

function computeStats(gear, opts, onComplete) {
    var id = makeid();
    simrequests[id] = onComplete
    simlib.postMessage({ msg: "computeStats", id: id, payload: { gear: gear, opts: opts } });
}

// Pulls options from the input 'options' pane for use in sims.
function getOptions() {
    var options = {};


    options.buffai = document.getElementById("buffai").checked;
    options.buffgotw = document.getElementById("buffgotw").checked;
    options.buffbk = document.getElementById("buffbk").checked;
    options.buffibow = document.getElementById("buffibow").checked;
    options.buffids = document.getElementById("buffids").checked;
    options.buffmoon = document.getElementById("buffmoon").checked;
    options.buffmoonrg = document.getElementById("buffmoonrg").checked;
    options.sbufws = document.getElementById("sbufws").checked;
    options.debuffjow = document.getElementById("debuffjow").checked;
    options.debuffisoc = document.getElementById("debuffisoc").checked;
    options.debuffmis = document.getElementById("debuffmis").checked;
    options.confbl = document.getElementById("confbl").checked;
    options.confmr = document.getElementById("confmr").checked;
    options.conbwo = document.getElementById("conbwo").checked;
    options.conmm = document.getElementById("conmm").checked;
    options.conbb = document.getElementById("conbb").checked;
    options.consmp = document.getElementById("consmp").checked;
    options.condp = document.getElementById("condp").checked;
    options.condr = document.getElementById("condr").checked;
    options.totms = document.getElementById("totms").checked;
    options.totwoa = document.getElementById("totwoa").checked;
    options.totcycl2p = document.getElementById("totcycl2p").checked;
    options.buffeyenight = document.getElementById("buffeyenight").checked;
    options.bufftwilightowl = document.getElementById("bufftwilightowl").checked;

    options.buffbl = parseInt(document.getElementById("buffbl").value) || 0;
    options.buffspriest = parseInt(document.getElementById("buffspriest").value) || 0;
    options.totwr = parseInt(document.getElementById("totwr").value) || 0;
    options.buffdrums = parseInt(document.getElementById("buffdrums").value) || 0;
    options.sbufrace = parseInt(document.getElementById("sbufrace").value) || 0;

    options.custom = {};
    options.custom.custint = parseInt(document.getElementById("custint").value) || 0;
    options.custom.custsp = parseInt(document.getElementById("custsp").value) || 0;
    options.custom.custsc = parseInt(document.getElementById("custsc").value) || 0;
    options.custom.custsh = parseInt(document.getElementById("custsh").value) || 0;
    options.custom.custha = parseInt(document.getElementById("custha").value) || 0;
    options.custom.custmp5 = parseInt(document.getElementById("custmp5").value) || 0;
    options.custom.custmana = parseInt(document.getElementById("custmana").value) || 0;

    options.dpsReportTime = 0;
    options.gcd = parseFloat(document.getElementById("custgcd").value) || 0;

    return options;
}

// basically this is a parser for the compact serializer for options.
//  for some reason I wrote the writer in go and the parser here. 
//  maybe its time to re-evaluate my life choices.
function setOptions(data) {

    document.getElementById("buffbl").selectedIndex = data[1];
    document.getElementById("buffdrums").selectedIndex = data[2];

    var dst = new ArrayBuffer(data.byteLength);
    new Uint8Array(dst).set(new Uint8Array(data));

    var buffView = new DataView(dst, 3);

    var idx = 0;

    var buffOpt1 = buffView.getUint8(idx, true); idx++;
    var buffOpt2 = buffView.getUint8(idx, true); idx++;

    document.getElementById("buffai").checked = (buffOpt1 & 1) == 1;
    document.getElementById("buffgotw").checked = (buffOpt1 & 1 << 1) == 1 << 1;
    document.getElementById("buffbk").checked = (buffOpt1 & 1 << 2) == 1 << 2;
    document.getElementById("buffibow").checked = (buffOpt1 & 1 << 3) == 1 << 3;
    document.getElementById("buffids").checked = (buffOpt1 & 1 << 4) == 1 << 4;
    document.getElementById("buffmoon").checked = (buffOpt1 & 1 << 5) == 1 << 5;
    document.getElementById("buffmoonrg").checked = (buffOpt1 & 1 << 6) == 1 << 6;
    document.getElementById("buffeyenight").checked = (buffOpt1 & 1 << 7) == 1 << 7;

    document.getElementById("bufftwilightowl").checked = (buffOpt2 & 1) == 1;
    document.getElementById("sbufws").checked = (buffOpt2 & 1 << 1) == 1 << 1;
    document.getElementById("debuffjow").checked = (buffOpt2 & 1 << 2) == 1 << 2;
    document.getElementById("debuffisoc").checked = (buffOpt2 & 1 << 3) == 1 << 3;
    document.getElementById("debuffmis").checked = (buffOpt2 & 1 << 4) == 1 << 4;

    idx++; // water shield procs not implemented
    document.getElementById("buffspriest").value = buffView.getUint16(idx, true); idx += 2;
    document.getElementById("sbufrace").selectedIndex = buffView.getUint8(idx, true); idx++;

    var numCustom = buffView.getUint8(idx, true); idx++;
    if (numCustom > 0) {
        document.getElementById("custint").value = buffView.getFloat64(7, true);
        // document.getElementById("custstm").value = buffView.getFloat64(7+8*1);
        document.getElementById("custsc").value = buffView.getFloat64(7 + 8 * 2, true);
        document.getElementById("custsh").value = buffView.getFloat64(7 + 8 * 3, true);
        document.getElementById("custsp").value = buffView.getFloat64(7 + 8 * 4, true);
        document.getElementById("custha").value = buffView.getFloat64(7 + 8 * 5, true);
        document.getElementById("custmp5").value = buffView.getFloat64(7 + 8 * 6, true);
        document.getElementById("custmana").value = buffView.getFloat64(7 + 8 * 7, true);
        idx += numCustom * 8;
    } else {
        document.getElementById("custint").value = 0;
        // document.getElementById("custstm").value = 0;
        document.getElementById("custsc").value = 0;
        document.getElementById("custsh").value = 0;
        document.getElementById("custsp").value = 0;
        document.getElementById("custha").value = 0;
        document.getElementById("custmp5").value = 0;
        document.getElementById("custmana").value = 0;
    }

    var consumOpt = buffView.getUint8(idx, true); idx++;
    document.getElementById("conbwo").checked = (consumOpt & 1) == 1;
    document.getElementById("conmm").checked = (consumOpt & 1 << 1) == 1 << 1;
    document.getElementById("confbl").checked = (consumOpt & 1 << 2) == 1 << 2;
    document.getElementById("confmr").checked = (consumOpt & 1 << 3) == 1 << 3;
    document.getElementById("conbb").checked = (consumOpt & 1 << 4) == 1 << 4;
    document.getElementById("condp").checked = (consumOpt & 1 << 5) == 1 << 5;
    document.getElementById("consmp").checked = (consumOpt & 1 << 6) == 1 << 6;
    document.getElementById("condr").checked = (consumOpt & 1 << 7) == 1 << 7;

    // talents
    idx += 9;

    document.getElementById("totwr").selectedIndex = buffView.getUint8(idx, true); idx++;
    var totemOpt = buffView.getUint8(idx, true); idx++;
    document.getElementById("totwoa").checked = (totemOpt & 1) == 1;
    document.getElementById("totms").checked = (totemOpt & 1 << 1) == 1 << 1;
    document.getElementById("totcycl2p").checked = (totemOpt & 1 << 2) == 1 << 2;
}

var castIDToName = {
    1: "LB",
    2: "CL",
    3: "TLC LB",
    999: "LB Overload", // this is just 1000-ID of the spell cast.
    998: "CL Overload",
}

// Populates the 'Sim' tab in the results pane.
function runsim(currentGear, fullLogs) {
    if (fullLogs) {
        var dur = parseInt(document.getElementById("logdur").value);
        var numClTargets = parseInt(document.getElementById("lognumClTargets").value);
        var firstOpts = getOptions();
        simulate(1, dur, numClTargets, currentGear, firstOpts, ["Adaptive"], null, true, (out) => {
            var logdiv = document.getElementById("simlogs");
            logdiv.innerText = out[0].Logs;
        });
        return;
    }

    var iters = parseInt(document.getElementById("iters").value);
    var dur = parseInt(document.getElementById("dur").value);
    var numClTargets = parseInt(document.getElementById("numClTargets").value);

    var lbout = document.getElementById("simrotlb");
    var priout = document.getElementById("simrotpri");
    var aiout = document.getElementById("simrotai");

    var metricHTML = `<div id="runningsim" uk-spinner="ratio: 1.5" style="margin:26%"></div>`;

    lbout.innerHTML = metricHTML;
    priout.innerHTML = metricHTML;
    aiout.innerHTML = metricHTML;

    var firstOpts = getOptions();
    firstOpts.exitoom = true;
    simulate(iters, 600, numClTargets, currentGear, firstOpts, ["FixedLBOnly"], 0, false, (out) => {
        var stats = out[0];
        var ttoom = stats.oomat;
        if (ttoom == 0) {
            ttoom = ">600";
        } else {
            ttoom = Math.round(ttoom);
        }
        lbout.innerHTML = `<div><h3>Mana</h3><text class="simnums">${ttoom}</text> sec<br /><text style="font-size:0.7em">to oom casting LB only ${Math.round(stats.dps)} DPS</text></div>`
    });
    firstOpts.dpsReportTime = 30; // report dps for 30 seconds only.
    simulate(iters, 600, numClTargets, currentGear, firstOpts, ["Fixed3LB1CL"], 0, false, (out) => {
        var stats = out[0];
        var dps = Math.max(stats.dps, stats.dpsAtOOM);
        var oomat = stats.oomat;
        if (oomat == 0) {
            oomat = ">600"
        } else {
            oomat = Math.round(oomat);
        }
        priout.innerHTML = `<div><h3>Peak</h3><text class="simnums">${Math.round(dps)}</text> dps<br /><text style="font-size:0.7em">${oomat}s to oom using CL on CD.</text></div>`
    });


    var secondOpts = getOptions();
    var start = new Date();
    simulate(iters, dur, numClTargets, currentGear, secondOpts, ["Adaptive"], 0, false, (out) => {
        var end = new Date();
        console.log(`The sim took ${end - start} ms`);
        var stats = out[0];
        console.log("AI Stats: ", stats);
        aiout.innerHTML = `<div><h3>Average</h3><text class="simnums">${Math.round(stats.dps)}</text> +/- ${Math.round(stats.dev)} dps<br /></div>`

        var rotstats = document.getElementById("rotstats");
        rotstats.innerHTML = "";
        Object.entries(stats.casts).forEach((entry) => {
            if (entry[1].count == 0) {
                return;
            }
            var cstat = entry[1];
            rotstats.innerHTML += `<text style="cursor:pointer" title="Avg Dmg: ${Math.round(cstat.dmg / cstat.count)} Crit: ${Math.round(cstat.crits / cstat.count * 100)}%">${castIDToName[entry[0]]}: ${Math.round(cstat.count / iters)}</text>`;
        });
        var percentoom = stats.numOOM / iters;
        if (percentoom > 0.02) {
            var dangerStyle = "";
            if (percentoom > 0.05 && percentoom <= 0.25) {
                dangerStyle = "border-color: #FDFD96;"
            } else if (percentoom > 0.25) {
                dangerStyle = "border-color: #FF6961;"
            }
            rotstats.innerHTML += `<text title="Downranking is not currently implemented." style="${dangerStyle};cursor: pointer">${Math.round(stats.numOOM / iters * 100)}% of simulations went OOM.`
        }

        var chartcanvas = document.createElement("canvas");
        var rotout = document.getElementById("rotout");
        var bounds = rotout.getBoundingClientRect();

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
        var colors = [];

        labels.forEach((k, i) => {
            vals.push(stats.dpsHist[k]);
            var val = parseInt(k);
            if (val > min && val < max) {
                colors.push('#1E87F0');
            } else {
                colors.push('#FF6961');
            }
        });

        var myChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    data: vals,
                    backgroundColor: colors,
                }]
            },
            options: {
                plugins: {
                    legend: {
                        display: false,
                        labels: {}
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            display: false
                        }
                    }
                }
            }
        });
    });
}

// Populates the 'Hasted Rotations' tab in results pane.
function hastedRotations(currentGear) {
    var opts = getOptions();
    opts.buffbl = 0;
    opts.buffdrum = 0;

    var hastes = [0, 50, 100, 200, 300, 400, 500, 600, 700, 788];
    var agentTypes = [
				"Fixed4LB1CL",
				"Fixed5LB1CL",
				"Fixed6LB1CL",
    ];

    // TODO: Fix this to match the new return values now that process is done in go WASM code.

    var hasteCounter = 0;
    var rows = document.getElementById("hasterots").firstElementChild.firstElementChild.children;
    hastes.forEach(haste => {
        hasteCounter++;
        var myCounter = hasteCounter;
        var row = rows[myCounter];
        row.children[1].innerHTML = "<div uk-spinner=\"ratio: 0.5\"></div>";
        row.children[2].innerText = "";

        simulate(1000, 40, 1, currentGear, opts, agentTypes, haste, false, (output) => {
            var maxdmg = 0.0;
            var maxrot = {};
						var maxIdx = 0;

            output.forEach((out, i) => {
                if (maxdmg < out.dps) {
                    maxdmg = out.dps
                    maxrot = out;
										maxIdx = i;
                }
            });

            var simdur = maxrot.SimSeconds;
            var rotTitle = "CL / " + (maxIdx + 4).toString() + "xLB";
            row.children[0].innerText = haste;
            row.children[1].innerText = rotTitle;
            row.children[2].innerText = "" + Math.round(maxrot.dps) + " +/- " + Math.round(maxrot.dev);
        });
    });
}

// Populates the 'Gear & Stat Weights' tab in results pane.
function calcStatWeights(gear) {
    var iters = parseInt(document.getElementById("switer").value);
    var dur = parseInt(document.getElementById("swdur").value);
    var numClTargets = parseInt(document.getElementById("swnumClTargets").value);
    var opts = getOptions();
		opts.agenttype = "Adaptive";

    var baseDPS = 0.0;
    var baseConf = 0.0;

    var sp_hitModDPS = 0.0;
    var sp_hitModConf = 0.0;

    var modDPS = [0, 0, 0, 0, 0, 0]; // SP, Int, Crit, Hit, Haste, MP5
    var modConf = [0, 0, 0, 0, 0, 0]
    var weights = [0, 0, 0, 0, 0, 0, 0]; // Int, X, Crit, Hit, Dmg, Haste, MP5
    modDPS.forEach((v, i) => {
        var cell = document.getElementById("w" + i.toString());
        cell.innerHTML = "<div uk-spinner=\"ratio: 1\"></div>";
    });

    // A base DPS without any modified stats.
    statweight(iters, dur, numClTargets, gear, opts, 0, 0, (res) => {
        var resVals = res.split(",")
        baseDPS = parseFloat(resVals[0]);
        baseConf = parseFloat(resVals[2]);
        if (baseDPS < 1) {
            // we failed.
            modDPS.forEach((v, i) => {
                var cell = document.getElementById("w" + i.toString());
                cell.innerHTML = `<text style="color:#FF6961">OOM</text>`;
            });
            var uptab = document.getElementById("upgrades");
            var nr = document.createElement("text");
            nr.innerText = `Simulations went OOM and so weights will be incorrect as downranking is not yet implemented.`;
            uptab.appendChild(nr);
        }
        console.log(`Base DPS: ${baseDPS} +/- ${baseConf}`);
    }); // base


    var done = [];
    var onfinish = () => {
        done.push(true);
        if (baseDPS < 1) {
            return;
        }
        if (modDPS[0] == 0) {
            return;
        }
        var baseMax = baseDPS + baseConf
        var baseMin = baseDPS - baseConf

        var ddpsMax = (modDPS[0] + modConf[0]) - baseMin;
        var ddpsMin = (modDPS[0] - modConf[0]) - baseMax;

        modDPS.forEach((v, i) => {
            if (v == 0) {
                return;
            }
            var cell = document.getElementById("w" + i.toString());
            var cellConf = document.getElementById("wc" + i.toString());
            if (v == -1) {
                cell.innerHTML = `<text style="color:#FF6961">OOM</text>`;
                return;
            }
            // sphit uses different value;
            if (i == 3 && sp_hitModDPS != 0.0) {
                var sphitMax = sp_hitModDPS + sp_hitModConf - baseMin;
                var sphitMin = sp_hitModDPS - sp_hitModConf - baseMax;

                var wmax = (v + modConf[i] - baseMin) / sphitMax;
                var wmin = (v - modConf[i] - baseMax) / sphitMin;
                var weight = (wmax + wmin) / 2;
                if (weight < 0.01) {
                    weight = 0.0;
                }
                if (currentFinalStats[3] > 202) {
                    weight = 0.0; // Just going to force 0 weight if you are hit capped.
                }

                weights[3] = weight;
                cell.innerText = weight.toFixed(2);
                cellConf.innerText = wmin.toFixed(2) + " - " + wmax.toFixed(2);
                return;
            }

            var wmax = (v + modConf[i] - baseMin) / ddpsMax;
            var wmin = (v - modConf[i] - baseMax) / ddpsMin;

            var weight = (wmax + wmin) / 2;
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
            if (wmax != wmin) {
                cellConf.innerText = wmin.toFixed(2) + " - " + wmax.toFixed(2);
            }
        });

        if (done.length == 7) {
            var oomed = true;
            modDPS.forEach((v) => {
                if (v > 0) {
                    oomed = false;
                }
            });
            if (oomed) {
                return;
            }
            showGearRecommendations(weights);
        }
    };

    statweight(iters, dur, numClTargets, gear, opts, 4, 20, (res) => {
        var resVals = res.split(",");
        sp_hitModDPS = parseFloat(resVals[0]);
        sp_hitModConf = parseFloat(resVals[2]);
        onfinish();
    }); // sp
    statweight(iters, dur, numClTargets, gear, opts, 3, 20, (res) => {
        var resVals = res.split(",");
        modConf[3] = parseFloat(resVals[2])
        modDPS[3] = parseFloat(resVals[0]);
        onfinish();
    }); // hit
    statweight(iters, dur, numClTargets, gear, opts, 4, 50, (res) => {
        var resVals = res.split(",");
        modConf[0] = parseFloat(resVals[2])
        modDPS[0] = parseFloat(resVals[0]);
        onfinish();
    }); // sp
    statweight(iters, dur, numClTargets, gear, opts, 0, 50, (res) => {
        var resVals = res.split(",");
        modConf[1] = parseFloat(resVals[2])
        modDPS[1] = parseFloat(resVals[0]);
        onfinish();
    }); // int
    statweight(iters, dur, numClTargets, gear, opts, 2, 50, (res) => {
        var resVals = res.split(",");
        modConf[2] = parseFloat(resVals[2])
        modDPS[2] = parseFloat(resVals[0]);
        onfinish();
    }); // crit
    statweight(iters, dur, numClTargets, gear, opts, 5, 50, (res) => {
        var resVals = res.split(",");
        modConf[4] = parseFloat(resVals[2])
        modDPS[4] = parseFloat(resVals[0]);
        onfinish();
    }); // haste
    statweight(iters, dur, numClTargets, gear, opts, 6, 50, (res) => {
        var resVals = res.split(",");
        modConf[5] = parseFloat(resVals[2])
        modDPS[5] = parseFloat(resVals[0]);
        onfinish();
    }); // mp5
}

function showGearRecommendations(weights) {
    var itemWeightsBySlot = {};
    var curSlotWeights = {};
    // 4 == dmg, 2 == crit
    var csdVal = (((currentFinalStats[4] * 0.795) + 603) * 2 * (currentFinalStats[2] / 2208) * 0.045) / 0.795;
    // process all items to find the weighted value.
    // find the value of each slots currently equipped item.
    Object.entries(gearUI.allitems).forEach((entry) => {
        var name = entry[0];
        var item = entry[1];
        if (item.Quality <= qualityFilter || item.Phase > phaseFilter) {
            return;
        }
        var value = 0.0;
        if (item.Stats != null) {
            weights.forEach((w, i) => {
                if (item.Stats[i] != null) {
                    value += item.Stats[i] * w
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
        var slotid = slotToID[item.Slot];
        if (slotid == "equipfinger") {
            slotid = "equipfinger1"
        }
        var curEquip = gearUI.currentGear[slotid];
        if (curEquip != null && curEquip.Name == item.Name) {
            curSlotWeights[item.Slot] = value;
        }
        itemWeightsBySlot[item.Slot].push({ Name: item.Name, Weight: value });
        itemWeightsBySlot[item.Slot] = itemWeightsBySlot[item.Slot].sort((a, b) => b.Weight - a.Weight);
    });
    var uptab = document.getElementById("upgrades");
    uptab.innerHTML = "";

    var curSlot = -1;
    Object.entries(itemWeightsBySlot).forEach((entry) => {
        if (entry[0] == 14) {
            // Skip trinkets for now. Trinkets will be separate
            return;
        }
        if (curSlot != entry[0]) {
            var row = document.createElement("tr");
            var col1 = document.createElement("td");
            slotToID[entry[0]].replace("equip", "");
            var title = slotToID[entry[0]].replace("equip", "");
            col1.innerHTML = "<h3>" + title.charAt(0).toUpperCase() + title.substr(1) + "</h3>";

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
            if (alt % 2 == 0) {
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

            var item = Object.assign({ Name: "" }, gearUI.allitems[v.Name]);
            simbut.addEventListener("click", (e) => {
                col4.innerHTML = "<div uk-spinner=\"ratio: 0.5\"></div>";
                var newgear = {};
                var slotID = slotToID[item.Slot];
                if (slotID == "equipfinger") {
                    slotID = "equipfinger1"; // hardcode finger 1 replacement.
                }
                Object.entries(gearUI.currentGear).forEach((entry) => {
                    if (entry[0] == slotID) {
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
                        item.Enchant = entry[1].Enchant;
                    } else {
                        newgear[entry[0]] = entry[1];
                    }
                });
                var iters = parseInt(document.getElementById("switer").value);
                var dur = parseInt(document.getElementById("swdur").value);
                var numClTargets = parseInt(document.getElementById("swnumClTargets").value);
                var opts = getOptions();
                simulate(iters, dur, numClTargets, cleanGear(newgear), opts, ["Adaptive"], null, false, (res) => {
                    var statistics = res[0];
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

window.addEventListener('click', function (e) {
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

    var importVal = location.hash.substr(1);

    if (importVal.length > 0) {
        importGear(importVal);
    } else {
        var glist = defaultGear;
        var gearCache = localStorage.getItem('cachedGear.v2');
        if (gearCache && gearCache.length > 0) {
            var parsedGear = JSON.parse(gearCache);
            if (parsedGear.length > 0) {
                glist = parsedGear;
            }
        }
        gearUI.updateEquipped(glist);
    }

    var currentGear = gearUI.currentGear;

    gearUI.addChangeListener((item, slot) => {
        updateGearStats(gearUI.currentGear);
    });
    updateGearStats(currentGear)

    var simlogrun = document.getElementById("simlogrun");
    simlogrun.addEventListener("click", (event) => {
        runsim(cleanGear(gearUI.currentGear), true);
    });

    var simrunbut = document.getElementById("simrunbut");
    simrunbut.addEventListener("click", (event) => {
        runsim(cleanGear(gearUI.currentGear));
    });

    var hastebut = document.getElementById("hastebut");
    hastebut.addEventListener("click", (event) => {
        hastedRotations(cleanGear(gearUI.currentGear));
    });

    var caclweights = document.getElementById("calcstatweight");
    caclweights.addEventListener("click", (event) => {
        calcStatWeights(cleanGear(gearUI.currentGear));
    });

    var inputs = document.querySelectorAll("#buffs input");
    for (var i = 0; i < inputs.length; i++) {
        var inp = inputs[i];
        inp.addEventListener("input", (e) => {
            updateGearStats(gearUI.currentGear);
        });
    }
    var selects = document.querySelectorAll("#buffs select");
    for (var i = 0; i < selects.length; i++) {
        var sel = selects[i];
        sel.addEventListener("change", (e) => {
            updateGearStats(gearUI.currentGear);
        })
    }

    updateGearSetList();

    changePhaseFilter({ target: document.getElementById("phasesel") });
    changeQualityFilter({ target: document.getElementById("qualsel") })


    window.addEventListener('hashchange', () => {
        var importVal = location.hash;
        if (importVal.length > 0) {
            importVal = importVal.substr(1);
            if (importVal != currentHash) {
                currentHash = importVal;
                importGear(importVal);
            }
        }
    });
}

// clearGear strips off all parts of gear that is non-changing. This lets us pass minimal data to sim and store in local storage.
function cleanGear(gear) {
    var cleanedGear = [];
    Object.entries(gear).forEach((entry) => {
        if (entry == null || entry == undefined) {
            return;
        }
        if (entry[1] == null || entry[1] == undefined) {
            return;
        }
        var it = {
            ID: entry[1].ID
        }
        if (entry[1].Gems != null) {
            it.g = [];
            entry[1].Gems.forEach((g) => {
                if (g == null) {
                    it.g.push(0);
                    return;
                }
                it.g.push(g.ID);
            });
        }
        if (entry[1].Enchant != null && entry[1].Enchant.ID > 0) {
            it.e = entry[1].Enchant.ID;
        }
        cleanedGear.push(it);
    });
    return cleanedGear
}

var currentFinalStats = {};

var statToName = [
    "StatInt",
    "StatStm",
    "StatSpellCrit",
    "StatSpellHit",
    "StatSpellDmg",
    "StatHaste",
    "StatMP5",
    "StatMana",
    "StatSpellPen",
    "StatSpirit",
    "StatLen",
]
// Updates the 'stats' pane in the viewport.
function updateGearStats(gearlist) {

    var cleanedGear = cleanGear(gearlist); // converts to array with minimal data for serialization.

    // TODO: Is this the best way?
    localStorage.setItem("cachedGear.v2", JSON.stringify(cleanedGear));

    computeStats(cleanedGear, null, (result) => {
        for (const [key, value] of Object.entries(result)) {
            var lab = document.getElementById(statToName[key].toLowerCase());
            if (lab != null) {
                lab.innerText = value.toFixed(0);
            }
        }
    });

    var opts = getOptions();
		opts.agenttype = "Adaptive";
    computeStats(cleanedGear, opts, (result) => {
        currentFinalStats = result.Stats;
        var sets = result.Sets;
        var setlist = document.getElementById("setlist");
        setlist.innerHTML = sets.join("<br />");
        for (const [key, value] of Object.entries(currentFinalStats)) {
            var lab = document.getElementById("f" + statToName[key].toLowerCase());
            if (key == 2) {
                lab.innerText = value.toFixed(0).toString() + " (" + (value / 22.08).toFixed(1) + "%)";
            } else if (key == 3) {
                lab.innerText = value.toFixed(0).toString() + " (" + (value / 12.6).toFixed(1) + "%)";
            } else if (key == 5) {
                lab.innerText = value.toFixed(0).toString() + " (" + (value / 15.76).toFixed(1) + "%)";
            } else if (lab == null) {
                // do nothing...
            } else {
                lab.innerText = value.toFixed(0);
            }
        }
    });

    exportGear(true); // this will update the URL
}


/// I hate html and javascript so much sometimes.

var panedrag = false;
var panediv = document.getElementById("panediv");
var calcpane = document.getElementById("calctabs");
var inpanel = document.getElementById("inputdata");
var outpanel = document.getElementById("calcdiv");
var h = window.innerHeight;

window.addEventListener("touchstart", (e) => {
    if (panediv == null) {
        panediv = document.getElementById("panediv");
        calcpane = document.getElementById("calctabs");
        inpanel = document.getElementById("inputdata");
        outpanel = document.getElementById("calcdiv");
    }
    if (panediv.contains(e.target) || calcpane.contains(e.target)) {
        panedrag = true;
        h = window.innerHeight;
        e.preventDefault();
    }
    if (e.changedTouches.length > 1) {
        panedrag = false;
    }
});
window.addEventListener("mousedown", (e) => {
    if (panediv == null) {
        panediv = document.getElementById("panediv");
        inpanel = document.getElementById("inputdata");
        outpanel = document.getElementById("calcdiv");
        calcpane = document.getElementById("calctabs");
    }
    if (panediv.contains(e.target) || calcpane.contains(e.target)) {
        panedrag = true;
        h = window.innerHeight;
        e.preventDefault();
    }
});
window.addEventListener("touchend", (e) => {
    panedrag = false;
});
window.addEventListener("mouseup", (e) => {
    panedrag = false;
});

window.addEventListener("touchmove", (e) => {
    if (panedrag) {
        var percent = e.changedTouches[0].pageY / (h + 60);
        inpanel.style.height = "calc(" + (percent * 100).toString() + "% - 100px)";
        outpanel.style.height = "calc(" + (100 - (percent * 100)).toString() + "% - 100px)";
        e.preventDefault();
    }
});

window.addEventListener("mousemove", (e) => {
    if (panedrag) {
        var percent = e.pageY / (h + 60);
        inpanel.style.height = "calc(" + (percent * 100).toString() + "% - 100px)";
        outpanel.style.height = "calc(" + (100 - (percent * 100)).toString() + "% - 100px)";
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
    for (var i = elements.length - 1; i >= 0; i--) {
        var e = elements[i];
        e.classList.remove(rm);
        e.classList.add(rp);
    }
}

var pulloutRight = -200;
function pulloutToggle() {
    var root = document.getElementById('root');
    var po = document.getElementById('pullout');

    if (pulloutRight < 0) {
        pulloutRight = 0;
        root.style.width = "calc(100% - 200px)";
    } else {
        pulloutRight = -200;
        root.style.width = "100%";
    }

    po.style.right = pulloutRight.toString() + "px";
}

function removegear() {
    gearUI.removeEquipped();
}

var phaseFilter = 5;
function changePhaseFilter(e) {
    var filter = e.target.value;
    phaseFilter = parseInt(filter);
    gearUI.setPhase(phaseFilter);
}

var qualityFilter = 0;
function changeQualityFilter(e) {
    var filter = e.target.value;
    qualityFilter = parseInt(filter);
    gearUI.setFilter(qualityFilter);
}

function saveGearSet() {
    var name = document.getElementById("gearname").value;
    var cleanedGear = cleanGear(gearUI.currentGear); // converts to array with minimal data for serialization.
    localStorage.setItem("stored." + name, JSON.stringify(cleanedGear));

    updateGearSetList();
}

function updateGearSetList() {
    var loadlist = document.getElementById("gearloader");
    loadlist.innerHTML = "";
    Object.keys(localStorage).forEach(k => {
        if (k.indexOf('stored.') == 0) {
            var name = k.split("stored.")[1];
            var item = document.createElement("option")
            item.innerText = name;
            loadlist.appendChild(item);
        }
        return false;
    })
}

function loadGearSet(event) {
    var name = document.getElementById("gearloader").value;
    var gearCache = localStorage.getItem("stored." + name);
    if (gearCache && gearCache.length > 0) {
        var parsedGear = JSON.parse(gearCache);
        if (parsedGear.length > 0) {
            var currentGear = gearUI.updateEquipped(parsedGear);
            updateGearStats(currentGear);
        }
    }
}

function deleteGearSet() {
    var name = document.getElementById("gearloader").value;
    localStorage.removeItem("stored." + name);
    updateGearSetList();
}

function importGearHandler() {
    var inputVal = document.getElementById("importexport").value;
    importGear(inputVal);
}

function importGear(inputVal) {
    var gearCache = inputVal;
    if (inputVal[0] != "[") { // that is opening brace for a gear list in JSON, but not valid base64
        if (window.pako === undefined) {
            loadPako(() => {
                importGear(inputVal); // try again
            });
            return;
        } else {
            var infdata = pakoInflate(inputVal);
            gearCache = infdata.gear;
            setOptions(infdata.buffs);
        }
    }
    if (gearCache && gearCache.length > 0) {
        var parsedGear = JSON.parse(gearCache);
        if (parsedGear.length > 0) {
            var currentGear = gearUI.updateEquipped(parsedGear);
            updateGearStats(currentGear);
        }
    }
}

function pakoInflate(v) {
    var binary = atob(v);
    var bytes = new Uint8Array(binary.length);
    for (let i = 0; i < bytes.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    // var bytes = base2048.decode(v);
    var dv = new DataView(bytes.buffer);
    var leng = dv.getInt32(0);

    return { gear: pako.inflate(bytes.subarray(4, leng + 4), { to: 'string' }), buffs: bytes.subarray(leng + 4, bytes.length) };
}


var currentHash = "";

function exportGear(compressed) {
    var cleanedGear = cleanGear(gearUI.currentGear); // converts to array with minimal data for serialization.\
    var box = document.getElementById("importexport");
    var enc = JSON.stringify(cleanedGear);
    if (compressed) {
        if (window.pako === undefined) {
            loadPako(() => {
                exportGear(compressed); // try again
            });
            return;
        } else {
            packOptions(getOptions(), (r) => {
                var val = pako.deflate(enc, { to: 'string' });
                var mergedArray = new Uint8Array(val.length + r.length + 4);
                var dv = new DataView(mergedArray.buffer);
                dv.setInt32(0, val.length);
                mergedArray.set(val, 4);
                mergedArray.set(r, 4 + val.length);
                // var output = base2048.encode(mergedArray);
                var output = btoa(String.fromCharCode(...mergedArray));
                // console.log(`JSON Size: ${enc.length}, Bin Size: ${mergedArray.length}, b2048: ${output.length}`);    
                // box.value = output;
                if (currentHash != output) {
                    currentHash = output;
                    location.hash = output;
                }
            });
            return;
        }
    }
    box.value = enc;
}

function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split('&');
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) == variable) {
            return decodeURIComponent(pair[1]);
        }
    }

    return "";
}

function loadPako(onLoad) {
    var script = document.createElement('script');
    script.onload = onLoad
    script.src = "pako.js";
    document.head.appendChild(script);
}
