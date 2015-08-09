function tof(x){
	x = +x;
	return x > 0.008856 ? Math.cbrt(x) : (903.3*x + 16) / 116;
}

function ConvertRGBToFloats(uints, floats){
	for(var i = 0; i < uints.length; i++){
		floats[i] = (uints[i]|0)/0xFF;
	}
}

function ConvertFloatsToRGB(floats, uints){
	for(var i = 0; i < floats.length; i++){
		uints[i] = (+floats[i])*0xFF;
	}
}

function Linearize(data){
	function linearize(v){
		return v < 0.04045 ? v / 12.92 : Math.pow((v+0.055)/1.055, 2.4);
	}

	for(var i = 0; i < data.length; i += 4){
		data[i+0] = linearize(+data[i+0]);
		data[i+1] = linearize(+data[i+1]);
		data[i+2] = linearize(+data[i+2]);
	}
}

function Delinearize(data){
	function delinearize(v){
		return v < 0.0031308 ? v * 12.92 : 1.055*Math.pow(v, 1.0/2.4) - 0.055;
	}

	for(var i = 0; i < data.length; i += 4){
		data[i+0] = delinearize(+data[i+0]);
		data[i+1] = delinearize(+data[i+1]);
		data[i+2] = delinearize(+data[i+2]);
	}
}

function RGBtoXYZ(data){
	var x, y, z, r, g, b;
	for(var i = 0; i < data.length; i += 4){
		r = +data[i+0]; g = +data[i+1]; b = +data[i+2];

		x = 0.4124564*r + 0.3575761*g + 0.1804375*b;
		y = 0.2126729*r + 0.7151522*g + 0.0721750*b;
		z = 0.0193339*r + 0.1191920*g + 0.9503041*b;

		data[i+0] = x; data[i+1] = y; data[i+2] = z;
	}
}

function XYZtoRGB(data){
	var x, y, z, r, g, b;
	for(var i = 0; i < data.length; i += 4){
		x = +data[i+0]; y = +data[i+1]; z = +data[i+2];

		r = +3.2404542*x - 1.5371385*y - 0.4985314*z;
		g = -0.9692660*x + 1.8760108*y + 0.0415560*z;
		b = +0.0556434*x - 0.2040259*y + 1.0572252*z;

		data[i+0] = r; data[i+1] = g; data[i+2] = b;
	}
}

function XYZtoLAB(data){
	var x, y, z, L, a, b;
	var fx, fy, fz;

	for(var i = 0; i < data.length; i += 4){
		x = +data[i+0]; y = +data[i+1]; z = +data[i+2];

		fx = tof(x / 0.95047);
		fy = tof(y / 1.00000);
		fz = tof(z / 1.08883);

		L = 116*fy - 16;
		a = 500 * (fx - fy);
		b = 200 * (fy - fz);

		data[i+0] = L; data[i+1] = a; data[i+2] = b;
	}
}

function LABtoXYZ(data){
	var x, y, z, L, a, b;
	var fx, fy, fz, xr, yr, zr;

	for(var i = 0; i < data.length; i += 4){
		L = +data[i+0]; a = +data[i+1]; b = +data[i+2];

		fy = (L + 16) / 116;
		fx = a/500 + fy;
		fz = fy - b/200;

		xr = fx > 0.20689303442 ? Math.pow(fx,3) : (116*fx - 16) / 903.3;
		yr = L > 7.9996248 ? Math.pow(fy,3) : L / 903.3;
		zr = fz > 0.20689303442 ? Math.pow(fz,3) : (116*fz - 16) / 903.3;

		data[i+0] = xr * 0.95047;
		data[i+1] = yr * 1.00000;
		data[i+2] = zr * 1.08883;
	}
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

function SRGBtoLAB(rgb, lab){
	ConvertRGBToFloats(rgb, lab);
	Linearize(lab);
	RGBtoXYZ(lab);
	XYZtoLAB(lab);
}

function LABtoSRGB(lab, rgb){
	LABtoXYZ(lab);
	XYZtoRGB(lab);
	Delinearize(lab);
	ConvertFloatsToRGB(lab, rgb);
}