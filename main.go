package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		src := c.Query("src")
		if src == "" {
			return fiber.NewError(fiber.StatusBadRequest, "src is required")
		}
		if _, err := url.Parse(src); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid src")
		}

		format := c.Query("format")
		if format != "png" && format != "jpg" && format != "svg" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid format, supported formats: png, jpg, svg")
		}

		quality := c.Query("quality", "90")
		if quality != "90" && format != "jpg" {
			return fiber.NewError(fiber.StatusBadRequest, "quality only supported for jpg format")
		}
		embedDiagram := c.Query("embed-diagram", "")
		if embedDiagram != "" && format != "svg" && format != "png" {
			return fiber.NewError(fiber.StatusBadRequest, "embed-diagram only supported for svg and png format")
		}
		embedSvgImages := c.Query("embed-svg-images", "")
		if embedSvgImages != "" && format != "svg" {
			return fiber.NewError(fiber.StatusBadRequest, "embed-svg-images only supported for svg format")
		}
		border := c.Query("border", "0")
		scale := c.Query("scale", "")
		width := c.Query("width", "")
		height := c.Query("height", "")
		pageIndex := c.Query("page-index", "")
		layers := c.Query("layers", "")
		svgTheme := c.Query("svg-theme", "light")
		if svgTheme != "light" && format != "svg" {
			return fiber.NewError(fiber.StatusBadRequest, "svg-theme only supported for svg format")
		}
		transparent := c.Query("transparent", "")

		req, err := http.NewRequest(http.MethodGet, src, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "cannot get file: "+err.Error())
		}

		l := log.Logger.With().
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Str("src", src).
			Logger()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			l.Err(err).Msg("cannot get file")
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			l.Err(err).Int("code", resp.StatusCode).Msg("response status code is not 200")
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		f, err := os.CreateTemp("", "*.drawio")
		if err != nil {
			l.Err(err).Msg("failed to create temporary file")
			return fiber.NewError(fiber.StatusInternalServerError)
		}
		defer f.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			l.Err(err).Msg("failed to save input file")
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		outFile := time.Now().String() + ".png"

		args := []string{
			f.Name(),
			"--no-sandbox",
			"--export",
			"--format", format,
			"--output", outFile,
			"--border", border,
		}
		if transparent != "" {
			args = append(args, "--transparent")
		}
		if format == "jpg" {
			args = append(args, "--quality", quality)
		}
		if embedDiagram != "" {
			args = append(args, "--embed-diagram")
		}
		if embedSvgImages != "" {
			args = append(args, "--embed-svg-images")
		}
		if scale != "" {
			args = append(args, "--scale", scale)
		}
		if width != "" {
			args = append(args, "--width", width)
		}
		if height != "" {
			args = append(args, "--height", height)
		}
		if pageIndex != "" {
			args = append(args, "--page-index", pageIndex)
		}
		if layers != "" {
			args = append(args, "--layers", layers)
		}
		if svgTheme != "" {
			args = append(args, "--svg-theme", svgTheme)
		}

		var sb bytes.Buffer
		cmd := exec.CommandContext(c.Context(), "drawio", args...)
		cmd.Stderr = &sb
		if err := cmd.Run(); err != nil {
			l.Err(err).Str("stderr", sb.String()).Msg("failed to convert")
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		l.Info().Str("source_file", f.Name()).Msg("converted")

		return c.SendFile(outFile, true)
	})

	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()
		<-ctx.Done()
		app.ShutdownWithContext(ctx)
	}()

	addr := os.Args[1]
	log.Info().Str("addr", addr).Msg("starting server")
	log.Fatal().Err(app.Listen(addr)).Send()
}
