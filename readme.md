# Usage

## Run server
```bash
docker run --rm -p 8080:8080 ghcr.io/rprtr258/drawishki:latest
```

## Render your images
```url
http://localhost:8080/?src=DRAWIO_FILE_URL&format=png&transparent=1
```

where `DRAWIO_FILE_URL` is address to your draw.io file reachable from server.

Can be used in markdown:
```md
![](http://localhost:8080/?src=DRAWIO_FILE_URL&format=png&transparent=1)
```

# Options
## Required
- `src` - url to draw.io file
- `format` - image format: `png`, `svg`, `jpg`

## Optional
- `transparent` - transparent background
- `quality` - jpg image quality, available only for `format=jpg`
- `embed-diagram` - embed draw.io diagram, available only for `format=svg` and `format=png`
- `embed-svg-images` - embed images into svg, available only for `format=svg`
- `border` - border size, default is `0`
- `scale` - image scale
- `width` - image width
- `height` - image height
- `page-index` - page index, default is `0`
- `layers` - layers to render
- `svg-theme` - svg theme: `light`, `dark`, default is `light`, available only for `format=svg`
