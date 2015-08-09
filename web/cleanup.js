function CloneImageData(m){
	return {
		data: new Float32Array(m.data),
		width: m.width,
		height: m.height
	};
}

function CleanupByBase(m, opts){
	function average(m){
		var t = 0;
		var data = m.data;
		for(var i = 0; i < data.length; i += 4){
			t += data[i];
		}
		return t / data.length;
	}

	var white = opts.whiteness;
	var lineWidth = opts.lineWidth;
	if(lineWidth === undefined){
		lineWidth = Math.min(m.width, m.height)*0.01|0;
	}

	Median(m, 5);

	var base = CloneImageData(m);
	Erase(base, lineWidth);

	var avg = average(m);
	var invspan = 1 / (avg / white);

	for(var i = 0; i < m.data.length; i += 4){
		m.data[i] = white + (m.data[i] - base.data[i]) * invspan;
	}
}