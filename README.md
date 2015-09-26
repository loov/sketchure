Sketchure is an algorithm that takes an image of a sketch taken in poor light
conditions and removes light variances on the paper.

On the left are input images taken with a phone and on the right are the
images run through the algorithm with default parameters.

![Comparison](https://raw.github.com/loov/sketchure/master/comparison.jpg)

## Algorithm

The algorithm itself is pretty trivial.

1. Convert image into a color-space that contains luminance values.
   YUV, YCbCr, YPbYPr, CIELab are all good options here. After comparing
   different approaches, I got the best results with YCbCr and 8bit precision.
2. Remove any noise created by the camera; we are assuming that we are in
   poor light conditions and the cameras are not perfects. For this step
   I used `median` filter with 3x3 kernel. Of course something more advanced
   can be used.
3. To extract the background we create a copy of the luminance channel and
   call it `base`.
4. Filter `base` through `erode` with a kernel size equal to the maximum line widht
   in pixels. This mostly erases the lines, while preserving the luminance of the
   paper.
5. Filter `base` through `blur`. During `erode` it created boxing artifacts,
   this makes them less visible. It will also smoothen the background tone,
   since we assume that the paper is not heavily textured.
6. Calculate the new luminance value with `whiteness + (L - B) / (average / whiteness)`.
   Where `L` is the old luminance value, `B` is the `base` luminance, `average` is the
   average luminance of the `base`, `whiteness` is a predefined constant to define the
   maximum `white` value.
7. Optionally desaturate the image.
8. Convert image back to the original color-space.

There are two reference implementations of the algorithm: `cleanup/normalize.go` and `js/cleanup.js`.