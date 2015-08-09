function Erase(m, N){
	Erode(m, N);
	Blur(m, N);
	return;
	for(var i = 0; i < N; i++){
		Erode1H(m);
		Erode1V(m);

		Blur1H(m);
		Blur1V(m);
	}
}

function Blur(m, N){
	for(var i = 0; i < N; i++){
		Blur1H(m);
		Blur1V(m);
	}
}

function Blur1H(m){
	var w = m.width;
	var h = m.height;
	var data = m.data;
	for(var y = 0; y < h; y++){
		var i = y * w * 4;
		var p = data[i];
		var z = p;

		for(var x = 0; x < w-1; x++){
			var n = data[i+4];
			data[i] = (p + z + n) / 3;
			p = z; z = n;
			i += 4;
		}
		var n = data[i];
		data[i] = (p + n + n) / 3;
	}
}

function Blur1V(m){
	var w = m.width;
	var stride = w * 4;
	var h = m.height;
	var data = m.data;
	for(var x = 0; x < w; x++){
		i = x * 4;
		var p = data[i];
		var z = p;
		for(var y = 0; y < h-1; y++){
			var n = data[i+stride];
			data[i] = (p + z + n) / 3;
			p = z; z = n;
			i += stride;
		}
		var n = data[i];
		data[i] = (p + n + n) / 3;
	}
}


function Erode(m, N){
	for(var i = 0; i < N; i++){
		Erode1H(m);
		Erode1V(m);
	}
}

function Erode1H(m){
	var w = m.width;
	var h = m.height;
	var data = m.data;
	for(var y = 0; y < h; y++){
		var i = y * w * 4;
		var p = data[i];
		var z = p;

		for(var x = 0; x < w-1; x++){
			var n = data[i+4];
			data[i] = Math.max(p, z, n);
			p = z; z = n;
			i += 4;
		}
		var n = data[i];
		data[i] = Math.max(p, n);
	}
}

function Erode1V(m){
	var w = m.width;
	var stride = w * 4;
	var h = m.height;
	var data = m.data;
	for(var x = 0; x < w; x++){
		i = x * 4;
		var p = data[i];
		var z = p;
		for(var y = 0; y < h-1; y++){
			var n = data[i+stride];
			data[i] = Math.max(p, z, n);
			p = z; z = n;
			i += stride;
		}
		var n = data[i];
		data[i] = Math.max(p, n);
	}
}
