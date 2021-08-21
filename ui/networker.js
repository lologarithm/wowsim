postMessage({
	msg: "ready"
});

var workerID = "";

addEventListener('message', async (e) => {
	const msg = e.data.msg;
	const payload = e.data.payload;

	if (msg == "apiCall") {
		const id = e.data.id;
		let url = '/api'

		let response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json;charset=utf-8'
			},
			body: JSON.stringify(payload)
		});

		let result = await response.json();
		postMessage({
			msg: "apiCall",
			payload: result,
			id: id,
		});
	} else if (msg == "setID") {
		workerID = payload;
		postMessage({ msg: "idconfirm" })
	}
}, false);
