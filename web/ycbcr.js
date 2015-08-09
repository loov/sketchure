function Linearize(data){
	for(var i = 0; i < data.length; i += 4){
		data[i+0] = Linearize.LUT[data[i+0]];
		data[i+1] = Linearize.LUT[data[i+1]];
		data[i+2] = Linearize.LUT[data[i+2]];
	}
}
Linearize.value = function(v){
	return v < 0.04045 ? v / 12.92 : Math.pow((v+0.055)/1.055, 2.4);
};

Linearize.LUT = new Uint8ClampedArray(256);
for(var i = 0; i < Linearize.LUT.length; i++){
	Linearize.LUT[i] = Linearize.value(i / 0xFF) * 0xFF;
}

function Delinearize(data){
	for(var i = 0; i < data.length; i += 4){
		data[i+0] = Delinearize.LUT[data[i+0]];
		data[i+1] = Delinearize.LUT[data[i+1]];
		data[i+2] = Delinearize.LUT[data[i+2]];
	}
}
Delinearize.value = function(v){
	return v < 0.0031308 ? v * 12.92 : 1.055*Math.pow(v, 1.0/2.4) - 0.055;
};

Delinearize.LUT = new Uint8ClampedArray(256);
for(var i = 0; i < Delinearize.LUT.length; i++){
	Delinearize.LUT[i] = Delinearize.value(i / 0xFF) * 0xFF;
}

function RGBtoYCBCR(data){
	var Y, Cb, Cr, R, G, B;
	for(var i = 0; i < data.length; i += 4){
		R = +data[i+0]; G = +data[i+1]; B = +data[i+2];

		Y  = ( 0.2990*R + 0.5870*G + 0.1140*B)|0;
		Cb = (-0.1687*R - 0.3313*G + 0.5000*B + 128)|0;
		Cr = ( 0.5000*R - 0.4187*G - 0.0813*B + 128)|0;

		data[i+0] = Y; data[i+1] = Cb; data[i+2] = Cr;
	}
}

function YCBCRtoRGB(data){
	var Y, Cb, Cr, R, G, B;
	for(var i = 0; i < data.length; i += 4){
		Y = +data[i+0]; Cb = +data[i+1]; Cr = +data[i+2];

		R = (Y + 1.40200*(Cr-128))|0;
    	G = (Y - 0.34414*(Cb-128) - 0.71414*(Cr-128))|0;
    	B = (Y + 1.77200*(Cb-128))|0;

		data[i+0] = R; data[i+1] = G; data[i+2] = B;
	}
}
