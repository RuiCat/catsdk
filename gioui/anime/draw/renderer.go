package draw

import (
	"image"
	"image/color"
	"image/gif"
	"mat/asm/f32"
	"math"
	"math/rand"
	"runtime"
	"sdk/x/image/draw"
	"sync"
	"time"
)

// Renderer 表示一个渲染器。
type Renderer struct {
	Scene              *Scene  // 是当前渲染的场景。
	Camera             *Camera // 是摄像机对象。
	Sampler            Sampler // 是采样器，用于生成随机点。
	Buffer             *Buffer // 是缓存区，用于存储渲染结果。
	SamplesPerPixel    int     // 表示每个像素的采样次数。
	StratifiedSampling bool    // 是否使用有序采样。
	AdaptiveSamples    int     // 使用自适应采样时，最大采样次数。
	AdaptiveThreshold  float32 // 自适应采样阈值。
	AdaptiveExponent   float32 // 自适应采样指数。
	FireflySamples     int     // 火花采样时的采样次数。
	FireflyThreshold   float32 // 火花采样阈值。
	NumCPU             int     // 是用于渲染的 CPU core 数量。
	Verbose            bool    // 是否开启调试模式。
}

func NewRenderer(scene *Scene, camera *Camera, sampler Sampler, w, h int) *Renderer {
	r := Renderer{}
	r.Scene = scene
	r.Camera = camera
	r.Sampler = sampler
	r.Buffer = NewBuffer(w, h)
	r.SamplesPerPixel = 20
	r.StratifiedSampling = true
	r.AdaptiveSamples = 0
	r.AdaptiveThreshold = 1
	r.AdaptiveExponent = 1
	r.FireflySamples = 0
	r.FireflyThreshold = 1
	r.NumCPU = runtime.NumCPU()
	r.Verbose = true
	return &r
}

func (r *Renderer) run() {
	scene := r.Scene
	camera := r.Camera
	sampler := r.Sampler
	buf := r.Buffer
	w, h := buf.W, buf.H
	spp := r.SamplesPerPixel
	sppRoot := int(math.Sqrt(float64(r.SamplesPerPixel)))
	ncpu := r.NumCPU

	runtime.GOMAXPROCS(ncpu)
	scene.Compile()
	ch := make(chan int, h)

	scene.rays = 0
	for i := 0; i < ncpu; i++ {
		go func(i int) {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			for y := i; y < h; y += ncpu {
				for x := 0; x < w; x++ {
					if r.StratifiedSampling {
						// stratified subsampling
						for u := 0; u < sppRoot; u++ {
							for v := 0; v < sppRoot; v++ {
								fu := (float64(u) + 0.5) / float64(sppRoot)
								fv := (float64(v) + 0.5) / float64(sppRoot)
								ray := camera.CastRay(x, y, w, h, fu, fv, rnd)
								sample := sampler.Sample(scene, ray, rnd)
								buf.AddSample(x, y, sample)
							}
						}
					} else {
						// random subsampling
						for i := 0; i < spp; i++ {
							fu := rnd.Float64()
							fv := rnd.Float64()
							ray := camera.CastRay(x, y, w, h, fu, fv, rnd)
							sample := sampler.Sample(scene, ray, rnd)
							buf.AddSample(x, y, sample)
						}
					}
					// adaptive sampling
					if r.AdaptiveSamples > 0 {
						v := buf.StandardDeviation(x, y).MaxComponent()
						v = Clamp(v/r.AdaptiveThreshold, 0, 1)
						v = f32.Pow(v, r.AdaptiveExponent)
						samples := int(v * float32(r.AdaptiveSamples))
						for i := 0; i < samples; i++ {
							fu := rnd.Float64()
							fv := rnd.Float64()
							ray := camera.CastRay(x, y, w, h, fu, fv, rnd)
							sample := sampler.Sample(scene, ray, rnd)
							buf.AddSample(x, y, sample)
						}
					}
					// firefly reduction
					if r.FireflySamples > 0 {
						if buf.StandardDeviation(x, y).MaxComponent() > r.FireflyThreshold {
							for i := 0; i < r.FireflySamples; i++ {
								fu := rnd.Float64()
								fv := rnd.Float64()
								ray := camera.CastRay(x, y, w, h, fu, fv, rnd)
								sample := sampler.Sample(scene, ray, rnd)
								buf.AddSample(x, y, sample)
							}
						}
					}
				}
				ch <- 1
			}
		}(i)
	}
	for i := 0; i < h; i++ {
		<-ch
	}
}

func (r *Renderer) writeImage(path string, buf *Buffer, channel Channel, wg *sync.WaitGroup) {
	defer wg.Done()
	im := buf.Image(channel)
	if err := SavePNG(path, im); err != nil {
		panic(err)
	}
}

func (r *Renderer) Render() image.Image {
	r.run()
	return r.Buffer.Image(ColorChannel)
}

func (r *Renderer) IterativeRender(iterations int) *gif.GIF {
	g := &gif.GIF{
		Delay:    make([]int, iterations),
		Image:    make([]*image.Paletted, iterations),
		Disposal: make([]byte, iterations),
	}
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		r.run()
		buf := r.Buffer.Copy()
		wg.Add(1)
		go func() {
			defer wg.Done()
			im := buf.Image(ColorChannel)
			bounds := im.Bounds()
			palettedImage := image.NewPaletted(bounds, quantizeMedianCutQuantizer.Quantize(make([]color.Color, 0, 256), im))
			draw.Src.Draw(palettedImage, bounds, im, bounds.Min)
			g.Image[i] = palettedImage
			g.Delay[i] = 30
			g.Disposal[i] = 3
		}()
	}
	wg.Wait()
	return g
}

func (r *Renderer) ChannelRender() <-chan image.Image {
	ch := make(chan image.Image)
	go func() {
		for i := 1; ; i++ {
			r.run()
			ch <- r.Buffer.Image(ColorChannel)
		}
	}()
	return ch
}

func (r *Renderer) FrameRender(path string, iterations int, wg *sync.WaitGroup) {
	for i := 1; i <= iterations; i++ {
		r.run()
	}
	buf := r.Buffer.Copy()
	wg.Add(1)
	go r.writeImage(path, buf, ColorChannel, wg)
}

func (r *Renderer) TimedRender(duration time.Duration) image.Image {
	start := time.Now()
	for {
		r.run()
		if time.Since(start) > duration {
			break
		}
	}
	return r.Buffer.Image(ColorChannel)
}
