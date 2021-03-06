(function(cleanup) {
	"use strict";

	function YCbCr(width, height) {
		var sz = width * height;
		this.Y = new Uint8ClampedArray(sz);
		this.Cb = new Uint8ClampedArray(sz);
		this.Cr = new Uint8ClampedArray(sz);
		this.width = width;
		this.height = height;
	}

	YCbCr.prototype = {
		desaturate: function() {
			var cb = this.Cb;
			var cr = this.Cr;
			var i;
			for (i = 0; i < cb.length; i++) {
				cb[i] = 128;
			}
			for (i = 0; i < cr.length; i++) {
				cr[i] = 128;
			}
		},
		assignImageData: function(imagedata) {
			var src = imagedata.data;
			var dst = this;

			for (var i = 0; i < src.length; i += 4) {
				var r = src[i + 0];
				var g = src[i + 1];
				var b = src[i + 2];


				var yy = (0.2990 * r + 0.5870 * g + 0.1140 * b) | 0;
				var cb = (-0.1687 * r - 0.3313 * g + 0.5000 * b + 128.0) | 0;
				var cr = (0.5000 * r - 0.4187 * g - 0.0813 * b + 128.0) | 0;


				var k = i >> 2;
				dst.Y[k] = yy;
				dst.Cb[k] = cb;
				dst.Cr[k] = cr;
			}
		},
		assignToImageData: function(imagedata) {
			var src = this;
			var dst = imagedata.data;

			for (var i = 0; i < dst.length; i += 4) {
				var k = i >> 2;
				var yy = src.Y[k];
				var cb = src.Cb[k] - 128.0;
				var cr = src.Cr[k] - 128.0;

				var r = (yy + 1.40200 * cr) | 0;
				var g = (yy - 0.34414 * cb - 0.71414 * cr) | 0;
				var b = (yy + 1.77200 * cb) | 0;

				dst[i + 0] = r;
				dst[i + 1] = g;
				dst[i + 2] = b;
			}
		}
	};

	function Channel(data, width, height) {
		this.data = data;
		this.width = width;
		this.height = height;
	}
	Channel.prototype = {
		clone: function() {
			return new Channel(
				new Uint8ClampedArray(this.data),
				this.width,
				this.height
			);
		},
		average: function() {
			var data = this.data;
			var w = this.width;
			var h = this.height;

			var t = 0.0;
			for (var y = 0; y < h; y++) {
				var i = y * w;
				var e = i + w;
				for (; i < e; i++) {
					t += data[i];
				}
			}

			return t / (w * h);
		},

		blur: function(steps) {
			// average
			function op(a, b, c) {
				return (a + b + c) / 3;
			}

			function h3(data, w, h) {
				for (var y = 0; y < h; y++) {
					var i = y * w;
					var e = (y + 1) * w - 1;
					var p = data[i];
					var z = data[i];
					for (; i < e; i++) {
						var n = data[i + 1];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			function v3(data, w, h) {
				for (var x = 0; x < w; x++) {
					var i = x;
					var e = (h - 1) * w + x;
					var p = data[i];
					var z = data[i];

					for (; i < e; i += w) {
						var n = data[i + w];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			var data = this.data;
			var width = this.width;
			var height = this.height;
			for (var i = 0; i < steps; i++) {
				h3(data, width, height);
				v3(data, width, height);
			}
		},
		erode: function(steps) {
			// max of a, b, c
			function op(a, b, c) {
				if (a >= b) {
					if (a >= c) {
						return a;
					}
				} else if (b > c) {
					return b;
				}
				return c;
			}

			function h3(data, w, h) {
				for (var y = 0; y < h; y++) {
					var i = y * w;
					var e = (y + 1) * w - 1;
					var p = data[i];
					var z = data[i];
					for (; i < e; i++) {
						var n = data[i + 1];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			function v3(data, w, h) {
				for (var x = 0; x < w; x++) {
					var i = x;
					var e = (h - 1) * w + x;
					var p = data[i];
					var z = data[i];

					for (; i < e; i += w) {
						var n = data[i + w];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			var data = this.data;
			var width = this.width;
			var height = this.height;
			for (var i = 0; i < steps; i++) {
				h3(data, width, height);
				v3(data, width, height);
			}
		},
		median: function(steps) {
			// median of a, b, c
			function op(a, b, c) {
				if (a > b) {
					if (b > c) {
						return b;
					} else if (a < c) {
						return a;
					}
				} else {
					if (a > c) {
						return a;
					} else if (b < c) {
						return b;
					}
				}
				return c;
			}

			function h3(data, w, h) {
				for (var y = 0; y < h; y++) {
					var i = y * w;
					var e = (y + 1) * w - 1;
					var p = data[i];
					var z = data[i];
					for (; i < e; i++) {
						var n = data[i + 1];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			function v3(data, w, h) {
				for (var x = 0; x < w; x++) {
					var i = x;
					var e = (h - 1) * w + x;
					var p = data[i];
					var z = data[i];

					for (; i < e; i += w) {
						var n = data[i + w];
						data[i] = op(p, z, n);
						p = z;
						z = n;
					}
					data[i] = op(p, data[i], data[i]);
				}
			}

			var data = this.data;
			var width = this.width;
			var height = this.height;
			for (var i = 0; i < steps; i++) {
				h3(data, width, height);
				v3(data, width, height);
			}
		}
	};


	function pipeline() {
		return {
			done: 0,
			total: 0,
			steps: [],
			step: function(name, weight, fn) {
				this.steps.push({
					name: name,
					weight: weight,
					fn: fn
				});
				this.total += weight;
				return this;
			},
			run: function(onprogress, onprogressnext) {
				for (var i = 0; i < this.steps.length; i++) {
					var step = this.steps[i];
					onprogressnext && onprogressnext(step.name);
					step.fn();
					this.done += step.weight;
					onprogress && onprogress(this.done / this.total);
				}

				onprogressnext && onprogressnext("done");
			}
		}
	}


	cleanup.ImageData = function(imagedata, opts) {
		var image, L, base;

		pipeline().
		step("convert", 1, function() {
			image = new YCbCr(imagedata.width, imagedata.height);
			image.assignImageData(imagedata);
		}).
		step("desaturate", 1, function() {
			if ((typeof opts.desaturate === "undefined") || opts.desaturate) {
				image.desaturate();
			}
		}).
		step("prepare", 1, function() {
			L = new Channel(image.Y, image.width, image.height);
		}).
		step("remove hot-pixels", 1, function() {
			L.median(1);
		}).
		step("creating background", 1, function() {
			base = L.clone();
		}).
		step("erasing lines from background", opts.lineWidth, function() {
			base.erode(opts.lineWidth);
		}).
		step("blurring artifacts", opts.lineWidth, function() {
			base.blur(opts.lineWidth);
		}).
		step("renormalizing image", 2, function() {
			var average = base.average();
			var white = opts.whiteness * 255.0 / 100.0;
			var invspan = 1.0 / (average / white);
			for (var y = 0; y < L.height; y++) {
				var i = y * L.width;
				var e = i + L.width;
				for (; i < e; i++) {
					L.data[i] = white + (L.data[i] - base.data[i]) * invspan;
				}
			}
		}).
		step("unconverting image", 1, function() {
			image.assignToImageData(imagedata);
		}).run(opts.onprogress, opts.onprogressnext);
	};
})(this.cleanup = {});