// Globals
const defaultGear = [
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

// This must be kept in sync with the enum in agents.go
const AGENT_TYPES = {
	"FIXED_3LB_1CL": 0,
	"FIXED_4LB_1CL": 1,
	"FIXED_5LB_1CL": 2,
	"FIXED_6LB_1CL": 3,
	"FIXED_7LB_1CL": 4,
	"FIXED_8LB_1CL": 5,
	"FIXED_9LB_1CL": 6,
	"FIXED_10LB_1CL": 7,
	"FIXED_LB_ONLY": 8,
	"ADAPTIVE": 9
};

// This must be kept in sync with the const values in stats.go
const STAT_IDX = {
	INT:        0,
	STAM:       1,
	SPELL_CRIT: 2,
	SPELL_HIT:  3,
	SPELL_DMG:  4,
	HASTE:      5,
	MP5:        6,
	MANA:       7,
	SPELL_PEN:  8,
	SPIRIT:     9
};
const STATS_LEN = 10;

class SimWorker {
	constructor() {
		this.numTasksRunning = 0;
		this.taskIdsToPromiseFuncs = {};
		this.worker = new window.Worker('simworker.js');

		let resolveReady = null;
		this.onReady = new Promise((_resolve, _reject) => {
			resolveReady = _resolve;
		});

		this.worker.onmessage = event => {
			if (event.data.msg == "ready") {
					this.worker.postMessage({ msg: "setID", payload: "1" });
					resolveReady();
			} else if (event.data.msg == "idconfirm") {
				// Do nothing
			} else {
					const id = event.data.id;
					if (!this.taskIdsToPromiseFuncs[id]) {
						console.warn("Unrecognized result id: " + id);
						return;
					}

					const promiseFuncs = this.taskIdsToPromiseFuncs[id];
					delete this.taskIdsToPromiseFuncs[id];
					this.numTasksRunnning--;

					promiseFuncs[0](event.data.payload);
			}
		};
	}

	makeTaskId() {
			let id = '';
			const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
			for (let i = 0; i < 16; i++) {
					id += characters.charAt(Math.floor(Math.random() * characters.length));
			}
			return id;
	}

	async doApiCall(request) {
		await this.onReady;

		const taskPromise = new Promise((resolve, reject) => {
			const id = this.makeTaskId();
			this.taskIdsToPromiseFuncs[id] = [resolve, reject];
			this.numTasksRunning++;
			
			this.worker.postMessage({
				msg: "apiCall",
				id: id,
				payload: request
			});
		});
		return await taskPromise;
	}
}

class WorkerPool {
	constructor(numWorkers) {
		this.workers = [];
		for (let i = 0; i < numWorkers; i++) {
			this.workers.push(new SimWorker());
		}
	}

	getLeastBusyWorker() {
		return this.workers.reduce(
				(curMinWorker, nextWorker) => curMinWorker.numTasksRunning < nextWorker.numTasksRunning ?
						curMinWorker : nextWorker);
	}

	async makeApiCall(request) {
		return await this.getLeastBusyWorker().doApiCall(request);
	}

	/**
   * The following methods are convenience wrappers for calling each API function.
	 */
	async getGearList() {
		const result = await this.makeApiCall({
			"RequestType": "GearList",
		});

		return result["GearList"];
	}

	async computeStats(computeStatsRequest) {
		const result = await this.makeApiCall({
			"RequestType": "ComputeStats",
			"ComputeStats": computeStatsRequest,
		});

		return result["ComputeStats"];
	}

	async statWeights(statWeightsRequest) {
		const result = await this.makeApiCall({
			"RequestType": "StatWeights",
			"StatWeights": statWeightsRequest,
		});

		return result["StatWeights"];
	}

	async runSimulation(simRequest) {
		const result = await this.makeApiCall({
			"RequestType": "Sim",
			"Sim": simRequest,
		});

		return result["Sim"];
	}

	async runBatchSimulation(batchSimRequest) {
		const result = await this.makeApiCall({
			"RequestType": "BatchSim",
			"BatchSim": batchSimRequest,
		});

		return result["BatchSim"];
	}

	async packOptions(packOptionsRequest) {
		const result = await this.makeApiCall({
			"RequestType": "PackOptions",
			"PackOptions": packOptionsRequest,
		});
		const packOptionsResult = result["PackOptions"];

		return new Uint8Array(packOptionsResult.Data, 0, packOptionsResult.Length);
	}
}

const workerPool = new WorkerPool(2);
workerPool.getGearList().then(gearListResult => popgear(gearListResult));

// Pulls options from the input 'options' pane for use in sims.
function getOptions() {
		// All of the names here need to match the corresponding properties in the Options struct in buffs.go
    var options = {};
		options.AgentType = AGENT_TYPES.ADAPTIVE;
    options.NumBloodlust = parseInt(document.getElementById("buffbl").value) || 0;
    options.NumDrums = parseInt(document.getElementById("buffdrums").value) || 0;
    options.DPSReportTime = 0;
    options.GCDMin = parseFloat(document.getElementById("custgcd").value) || 0;

		options.Buffs = {};
    options.Buffs.ArcaneInt = document.getElementById("buffai").checked;
    options.Buffs.GiftOftheWild = document.getElementById("buffgotw").checked;
    options.Buffs.BlessingOfKings = document.getElementById("buffbk").checked;
    options.Buffs.ImprovedBlessingOfWisdom = document.getElementById("buffibow").checked;
    options.Buffs.ImprovedDivineSpirit = document.getElementById("buffids").checked;
    options.Buffs.Moonkin = document.getElementById("buffmoon").checked;
    options.Buffs.MoonkinRavenGoddess = document.getElementById("buffmoonrg").checked;
    options.Buffs.SpriestDPS = parseInt(document.getElementById("buffspriest").value) || 0;
    options.Buffs.EyeOfNight = document.getElementById("buffeyenight").checked;
    options.Buffs.TwilightOwl = document.getElementById("bufftwilightowl").checked;
    options.Buffs.WaterShield = document.getElementById("sbufws").checked;
    options.Buffs.Race = parseInt(document.getElementById("sbufrace").value) || 0;

		// Target debuffs
    options.Buffs.JudgementOfWisdom = document.getElementById("debuffjow").checked;
    options.Buffs.ImpSealOfCrusader = document.getElementById("debuffisoc").checked;
    options.Buffs.Misery = document.getElementById("debuffmis").checked;

		options.Consumes = {};
    options.Consumes.FlaskOfBlindingLight = document.getElementById("confbl").checked;
    options.Consumes.FlaskOfMightyRestoration = document.getElementById("confmr").checked;
    options.Consumes.BrilliantWizardOil = document.getElementById("conbwo").checked;
    options.Consumes.MajorMageblood = document.getElementById("conmm").checked;
    options.Consumes.BlackendBasilisk = document.getElementById("conbb").checked;
    options.Consumes.SuperManaPotion = document.getElementById("consmp").checked;
    options.Consumes.DestructionPotion = document.getElementById("condp").checked;
    options.Consumes.DarkRune = document.getElementById("condr").checked;

		options.Totems = {};
    options.Totems.TotemOfWrath = parseInt(document.getElementById("totwr").value) || 0;
    options.Totems.WrathOfAir = document.getElementById("totwoa").checked;
    options.Totems.Cyclone2PC = document.getElementById("totcycl2p").checked;
    options.Totems.ManaStream = document.getElementById("totms").checked;

    options.Buffs.Custom = new Array(STATS_LEN).fill(0);
    options.Buffs.Custom[STAT_IDX.INT]        = parseInt(document.getElementById("custint").value) || 0;
    options.Buffs.Custom[STAT_IDX.SPELL_DMG]  = parseInt(document.getElementById("custsp").value) || 0;
    options.Buffs.Custom[STAT_IDX.SPELL_CRIT] = parseInt(document.getElementById("custsc").value) || 0;
    options.Buffs.Custom[STAT_IDX.SPELL_HIT]  = parseInt(document.getElementById("custsh").value) || 0;
    options.Buffs.Custom[STAT_IDX.HASTE]      = parseInt(document.getElementById("custha").value) || 0;
    options.Buffs.Custom[STAT_IDX.MP5]        = parseInt(document.getElementById("custmp5").value) || 0;
    options.Buffs.Custom[STAT_IDX.MANA]       = parseInt(document.getElementById("custmana").value) || 0;

		options.Talents = {
			LightningOverload:  5,
			ElementalPrecision: 3,
			NaturesGuidance:    3,
			TidalMastery:       5,
			ElementalMastery:   true,
			UnrelentingStorm:   3,
			CallOfThunder:      5,
			Concussion:         5,
			Convection:         5,
		};

		options.Encounter = {
			Duration: 0,
			NumClTargets: 1,
		};

    return options;
}

// basically this is a parser for the compact serializer for options.
//  for some reason I wrote the writer in go and the parser here. 
//  maybe its time to re-evaluate my life choices.
function setOptions(data) {
		if (data.byteLength < 3) {
			console.log('Empty options data loaded');
			return;
		}

    document.getElementById("buffbl").selectedIndex = data[1];
    document.getElementById("buffdrums").selectedIndex = data[2];

    const dst = new ArrayBuffer(data.byteLength);
    new Uint8Array(dst).set(new Uint8Array(data));

    const buffView = new DataView(dst, 3);

    let idx = 0;

    const buffOpt1 = buffView.getUint8(idx, true); idx++;
    const buffOpt2 = buffView.getUint8(idx, true); idx++;

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

    const numCustom = buffView.getUint8(idx, true); idx++;
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

    const consumOpt = buffView.getUint8(idx, true); idx++;
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
    const totemOpt = buffView.getUint8(idx, true); idx++;
    document.getElementById("totwoa").checked = (totemOpt & 1) == 1;
    document.getElementById("totms").checked = (totemOpt & 1 << 1) == 1 << 1;
    document.getElementById("totcycl2p").checked = (totemOpt & 1 << 2) == 1 << 2;
}

const castIDToName = {
    1: "LB",
    2: "CL",
    3: "TLC LB",
    999: "LB Overload", // this is just 1000-ID of the spell cast.
    998: "CL Overload",
}

// Only works for POD-type objects
function deepCopy(obj) {
	return JSON.parse(JSON.stringify(obj));
}

function runSimWithLogs(gearSpec) {
		const options = getOptions();
		options.AgentType = AGENT_TYPES.ADAPTIVE;
		options.Encounter.Duration = parseInt(document.getElementById("logdur").value);
		options.Encounter.NumClTargets = parseInt(document.getElementById("lognumClTargets").value);

		const simRequest = {
			Options: options,
			Gear: gearSpec,
			Iterations: 1,
			IncludeLogs: true,
		};

		workerPool.runSimulation(simRequest).then(simResult => {
				const logdiv = document.getElementById("simlogs");
				logdiv.innerText = simResult.Logs;
		});
}

// Populates the 'Sim' tab in the results pane.
function runSim(gearSpec) {
		const sharedOptions = getOptions();
		sharedOptions.Encounter = {
				Duration: parseInt(document.getElementById("dur").value),
				NumClTargets: parseInt(document.getElementById("numClTargets").value),
		};

		const sharedSimRequest = {
			Options: sharedOptions,
			Gear: gearSpec,
			Iterations: parseInt(document.getElementById("iters").value),
		};

    const pendingMetricHTML = `<div id="runningsim" uk-spinner="ratio: 1.5" style="margin:26%"></div>`;

		{
			const resultsElem = document.getElementById("simrotlb");
			resultsElem.innerHTML = pendingMetricHTML;
			
			const simRequest = deepCopy(sharedSimRequest);
			simRequest.Options.AgentType = AGENT_TYPES.FIXED_LB_ONLY;
			simRequest.Options.Encounter.Duration = 600;
			simRequest.Options.ExitOnOOM = true;

			workerPool.runSimulation(simRequest).then(simResult => {
					const oomAtText = simResult.oomAtAvg ? Math.round(simResult.oomAtAvg) : ">600";
					resultsElem.innerHTML = `<div><h3>Mana</h3><text class="simnums">${oomAtText}</text> sec<br /><text style="font-size:0.7em">to oom casting LB only ${Math.round(simResult.DpsAvg)} DPS</text></div>`
			});
		}

		{
			const resultsElem = document.getElementById("simrotpri");
			resultsElem.innerHTML = pendingMetricHTML;

			const simRequest = deepCopy(sharedSimRequest);
			simRequest.Options.AgentType = AGENT_TYPES.FIXED_3LB_1CL;
			simRequest.Options.Encounter.Duration = 600;
			simRequest.Options.DPSReportTime = 30;

			workerPool.runSimulation(simRequest).then(simResult => {
					const oomAtText = simResult.oomAtAvg ? Math.round(simResult.oomAtAvg) : ">600";
					const dps = Math.max(simResult.DpsAvg, simResult.DpsAtOomAvg);
					resultsElem.innerHTML = `<div><h3>Peak</h3><text class="simnums">${Math.round(dps)}</text> dps<br /><text style="font-size:0.7em">${oomAtText}s to oom using CL on CD.</text></div>`
			});
		}

		const resultsElem = document.getElementById("simrotai");
		resultsElem.innerHTML = pendingMetricHTML;

		const simRequest = deepCopy(sharedSimRequest);
		simRequest.Options.AgentType = AGENT_TYPES.ADAPTIVE;

    const start = new Date();
		workerPool.runSimulation(simRequest).then(simResult => {
        const end = new Date();
        console.log(`The sim took ${end - start} ms`);
        console.log("AI Stats: ", simResult);
        resultsElem.innerHTML = `<div><h3>Average</h3><text class="simnums">${Math.round(simResult.DpsAvg)}</text> +/- ${Math.round(simResult.DpsStDev)} dps<br /></div>`

        const rotstats = document.getElementById("rotstats");
        rotstats.innerHTML = "";
        Object.entries(simResult.Casts).forEach(castEntry => {
						const castID = castEntry[0];
						const cast = castEntry[1];
            if (cast.Count == 0) {
                return;
            }
            rotstats.innerHTML += `<text style="cursor:pointer" title="Avg Dmg: ${Math.round(cast.Dmg / cast.Count)} Crit: ${Math.round(cast.Crits / cast.Count * 100)}%">${castIDToName[castID]}: ${Math.round(cast.Count / simRequest.Iterations)}</text>`;
        });
        const percentOom = simResult.NumOom / simRequest.Iterations;
        if (percentOom > 0.02) {
            var dangerStyle = "";
            if (percentOom > 0.05 && percentOom <= 0.25) {
                dangerStyle = "border-color: #FDFD96;"
            } else if (percentOom > 0.25) {
                dangerStyle = "border-color: #FF6961;"
            }
            rotstats.innerHTML += `<text title="Downranking is not currently implemented." style="${dangerStyle};cursor: pointer">${Math.round(percentOom * 100)}% of simulations went OOM.`
        }

        const rotout = document.getElementById("rotout");
        const bounds = rotout.getBoundingClientRect();
				const chartCanvas = createDpsChartFromSimResult(simResult, bounds);

        const rotchart = document.getElementById("rotchart");
        rotchart.innerHTML = "";
        rotchart.appendChild(chartCanvas);
    });
}

function createDpsChartFromSimResult(simResult, chartBounds) {
		const chartCanvas = document.createElement("canvas");
		const ctx = chartCanvas.getContext('2d');
		chartCanvas.height = chartBounds.height - 30;
		chartCanvas.width = chartBounds.width;

		const min = simResult.DpsAvg - simResult.DpsStDev;
		const max = simResult.DpsAvg + simResult.DpsStDev;
		const vals = [];
		const colors = [];

		const labels = Object.keys(simResult.DpsHist)
		labels.forEach((k, i) => {
				vals.push(simResult.DpsHist[k]);
				const val = parseInt(k);
				if (val > min && val < max) {
						colors.push('#1E87F0');
				} else {
						colors.push('#FF6961');
				}
		});

		const chart = new Chart(ctx, {
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
		return chartCanvas;
}

// Populates the 'Hasted Rotations' tab in results pane.
function hastedRotations(gearSpec) {
    const sharedOptions = getOptions();
    sharedOptions.NumBloodlust = 0;
    sharedOptions.NumDrums = 0;
		sharedOptions.Encounter.Duration = 40;
		sharedOptions.Buffs.Custom[STAT_IDX.HASTE] = 0;

		const sharedSimRequest = {
			Options: sharedOptions,
			Gear: gearSpec,
			Iterations: 1000,
		};

    const hastes = [0, 50, 100, 200, 300, 400, 500, 600, 700, 788];
    const agentTypes = [
				AGENT_TYPES.FIXED_4LB_1CL,
				AGENT_TYPES.FIXED_5LB_1CL,
				AGENT_TYPES.FIXED_6LB_1CL,
    ];

    // TODO: Fix this to match the new return values now that process is done in go WASM code.

		workerPool.computeStats({
				Options: sharedOptions,
				Gear: gearSpec,
		}).then(computeStatsResult => {
				const finalHasteValue = computeStatsResult.FinalStats[STAT_IDX.HASTE];
				const rows = document.getElementById("hasterots").firstElementChild.firstElementChild.children;
				hastes.forEach((haste, i) => {
						// Subtract haste from gear/etc
						const hasteValue = haste - finalHasteValue;

						const row = rows[i];
						row.children[1].innerHTML = "<div uk-spinner=\"ratio: 0.5\"></div>";
						row.children[2].innerText = "";

						const batchSimRequest = {
								Requests: agentTypes.map(agentType => {
										const simRequest = deepCopy(sharedSimRequest);
										simRequest.Options.AgentType = agentType;
										simRequest.Options.Buffs.Custom[STAT_IDX.HASTE] = hasteValue;
										return simRequest;
								}),
						};
						workerPool.runBatchSimulation(batchSimRequest).then(batchSimResult => {
								let maxIdx = 0;
								let bestResult = null;
								batchSimResult.Results.forEach((simResult, idx) => {
										if (!bestResult || simResult.DpsAvg > bestResult.DpsAvg) {
												maxIdx = idx;
												bestResult = simResult;
										}
								});

								const rotTitle = "CL / " + (maxIdx + 4).toString() + "xLB";
								row.children[0].innerText = haste;
								row.children[1].innerText = rotTitle;
								row.children[2].innerText = Math.round(bestResult.DpsAvg) + " +/- " + Math.round(bestResult.DpsStDev);
						});
				});
		});
}

function stDevToConf90(stDev, N) {
	return 1.645 * stDev / Math.sqrt(N);
}

// Populates the 'Gear & Stat Weights' tab in results pane.
function calcStatWeights(gearSpec) {
		const statWeightsRequest = {
			Options: getOptions(),
			Gear: gearSpec,
			Iterations: parseInt(document.getElementById("switer").value),
		};
		statWeightsRequest.Options.Encounter.Duration = parseInt(document.getElementById("swdur").value);
		statWeightsRequest.Options.Encounter.NumClTargets = parseInt(document.getElementById("swnumClTargets").value);

		const statsTested = [
				STAT_IDX.SPELL_DMG, 
				STAT_IDX.INT,
				STAT_IDX.SPELL_CRIT,
				STAT_IDX.SPELL_HIT,
				STAT_IDX.HASTE,
				STAT_IDX.MP5,
		];
		const weightElems = statsTested.map((stat, i) => document.getElementById("w" + i.toString()));
		const weightConfidenceElems = statsTested.map((stat, i) => document.getElementById("wc" + i.toString()));

		statsTested.forEach((stat, i) => {
			weightElems[i].innerHTML = "<div uk-spinner=\"ratio: 1\"></div>";
			weightConfidenceElems[i].innerHTML = "";
		});

		workerPool.statWeights(statWeightsRequest).then(statWeightsResult => {
				statsTested.forEach((stat, i) => {
						const ep = statWeightsResult.EpValues[stat];
						const epStDev = statWeightsResult.EpValuesStDev[stat];
						const epConf90 = stDevToConf90(epStDev, statWeightsRequest.Iterations);

						weightElems[i].innerText = ep.toFixed(2);
						weightConfidenceElems[i].innerText = (ep - epConf90).toFixed(2) + " - " + (ep + epConf90).toFixed(2);
				});

				showGearRecommendations(statWeightsResult.EpValues);
		});
}

function showGearRecommendations(weights) {
    const itemWeightsBySlot = {};
    const curSlotWeights = new Array(20).fill(0);
    const csdVal = (((currentFinalStats[STAT_IDX.SPELL_DMG] * 0.795) + 603) * 2 * (currentFinalStats[STAT_IDX.SPELL_CRIT] / 2208) * 0.045) / 0.795;
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
                if (item.Stats[i]) {
                    value += item.Stats[i] * w
                }
            });
        }
        if (itemWeightsBySlot[item.Slot] == null) {
            itemWeightsBySlot[item.Slot] = [];
        }
        if (item.GemSlots != null && item.GemSlots.length > 0) {
            let numGems = item.GemSlots.length;
            if (item.GemSlots[0] == 1) {
                numGems--;
                // how to value a CSD
                // ~ spellpower * crit chance * 0.09 = increased damage per cast.
                value += csdVal;
            }
            value += (numGems * 9) * weights[STAT_IDX.SPELL_DMG]; // just for measuring use 9 spell power gems in every slot.
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
								const simRequest = {
									Options: getOptions(),
									Gear: toGearSpec(newgear),
									Iterations: parseInt(document.getElementById("switer").value),
								};
								simRequest.Options.AgentType = AGENT_TYPES.ADAPTIVE;
								simRequest.Options.Encounter = {
									Duration: parseInt(document.getElementById("swdur").value),
									NumClTargets: parseInt(document.getElementById("swnumClTargets").value),
								};
								workerPool.runSimulation(simRequest).then(simResult => {
                    col4.innerText = Math.round(simResult.DpsAvg).toString() + " +/- " + Math.round(simResult.DpsStDev).toString();
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
        runSimWithLogs(toGearSpec(gearUI.currentGear));
    });

    var simrunbut = document.getElementById("simrunbut");
    simrunbut.addEventListener("click", (event) => {
        runSim(toGearSpec(gearUI.currentGear));
    });

    var hastebut = document.getElementById("hastebut");
    hastebut.addEventListener("click", (event) => {
        hastedRotations(toGearSpec(gearUI.currentGear));
    });

    var caclweights = document.getElementById("calcstatweight");
    caclweights.addEventListener("click", (event) => {
        calcStatWeights(toGearSpec(gearUI.currentGear));
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

// toGearSpec returns a 'gear specification' which is the minimal amount of info needed
// to specify a set of gear.
function toGearSpec(gear) {
		const gearSpec = [];
		Object.values(gear).forEach(item => {
				if (!item) {
						return;
				}
				const itemSpec = {
					NameOrId: item.Name,
					Enchant: {
						NameOrId: (item.Enchant || {}).Name
					},
					Gems: (item.Gems || []).map(gem => {
						NameOrId: gem.Name
					})
				};
				gearSpec.push(itemSpec);
		});
		return gearSpec;
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
    const gearSpec = toGearSpec(gearlist); // converts to array with minimal data for serialization.
    // TODO: Is this the best way?
    localStorage.setItem("cachedGear.v2", JSON.stringify(gearSpec));
    exportGear(true); // this will update the URL

    const options = getOptions();
		options.AgentType = AGENT_TYPES.ADAPTIVE;

		const computeStatsRequest = {
			Options: options,
			Gear: gearSpec,
		};

		workerPool.computeStats(computeStatsRequest).then(computeStatsResult => {
				currentFinalStats = computeStatsResult.FinalStats;

        for (const [key, value] of Object.entries(computeStatsResult.GearOnly)) {
            var lab = document.getElementById(statToName[key].toLowerCase());
            if (lab != null) {
                lab.innerText = value.toFixed(0);
            }
        }

        var setlist = document.getElementById("setlist");
        setlist.innerHTML = computeStatsResult.Sets.join("<br />");

        for (const [key, value] of Object.entries(computeStatsResult.FinalStats)) {
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
    var cleanedGear = toGearSpec(gearUI.currentGear); // converts to array with minimal data for serialization.
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
    let gearCache = inputVal;
    if (inputVal[0] != "[") { // that is opening brace for a gear list in JSON, but not valid base64
        if (window.pako === undefined) {
            loadPako(() => {
                importGear(inputVal); // try again
            });
            return;
        } else {
            const infdata = pakoInflate(inputVal);
            gearCache = infdata.gear;
            setOptions(infdata.buffs);
        }
    }
    if (gearCache && gearCache.length > 0) {
        const parsedGear = JSON.parse(gearCache);
        if (parsedGear.length > 0) {
            const currentGear = gearUI.updateEquipped(parsedGear);
            updateGearStats(currentGear);
        }
    }
}

function pakoInflate(v) {
    const binary = atob(v);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < bytes.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    // var bytes = base2048.decode(v);
    const dv = new DataView(bytes.buffer);
    const leng = dv.getInt32(0);

    return { gear: pako.inflate(bytes.subarray(4, leng + 4), { to: 'string' }), buffs: bytes.subarray(leng + 4, bytes.length) };
}


var currentHash = "";

function exportGear(compressed) {
    const gearSpec = toGearSpec(gearUI.currentGear); // converts to array with minimal data for serialization.\
    const box = document.getElementById("importexport");
    const gearSpecStr = JSON.stringify(gearSpec);
    if (compressed) {
        if (window.pako === undefined) {
            loadPako(() => {
                exportGear(compressed); // try again
            });
            return;
        } else {
            workerPool.packOptions({ Options: getOptions() }).then(packedOptions => {
                const val = pako.deflate(gearSpecStr, { to: 'string' });
                const mergedArray = new Uint8Array(val.length + packedOptions.length + 4);
                var dv = new DataView(mergedArray.buffer);
                dv.setInt32(0, val.length);
                mergedArray.set(val, 4);
                mergedArray.set(packedOptions, 4 + val.length);
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
    box.value = gearSpecStr;
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
