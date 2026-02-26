package main

/*
package main

import (
	"backend/internal/api/dto"
	"backend/internal/domain/service"
	"backend/internal/infra"
	"backend/internal/infra/repo"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Download(videoID string) (io.ReadCloser, error) {
	// Create a pipe
	pr, pw := io.Pipe()

	// token := "MniIuwfQi9boiIiHAWCMVIOgxWjEpX6_hqZ5p3rEO3IVGME6IvmCTlJCrN67VRwKsUscHgB-wBaEAGYVfoGHaMoVqTxuRlBiLgVRKy0QUbJVNWcUI3M1RzQFkSdysAJh4pnw4b15_Aq6rDRHaqqUgvi02N8K4dkLdUM="

	// Prepare command
	cmd := exec.Command(
		"yt-dlp",
		"--cache-dir", "./cache", // cache
		"-N", strconv.Itoa(runtime.NumCPU()), // multi-threading
		"--socket-timeout", "3",
		"--retries", "3",
		// TODO: optimization shit, add youtube potoken and download a specific custom yt-dlp binary, that supports sabr, use yt music host
		// "--extractor-args", "youtube:formats=duplicate;webpage_skip=1;player_skip=initial_data;innertube_host=www.youtube.com;player-client=web;skip=translated_subs,hls,dash;po_token=web.gvs+"+token,
		"--extractor-args", "youtube:formats=duplicate;webpage_skip=1;player_skip=initial_data;skip=translated_subs,hls,dash", // player-client=tv_simply;
		"-f", "ba[ext=m4a]",
		"-o", "-",
		"https://www.youtube.com/watch?v="+videoID,
	)

	println(cmd.String())

	// Set stdout to our pipe writer
	cmd.Stdout = pw

	// Capture stderr for logging
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	// Start the command (don't wait for it to finish)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("download failed to start: %w\nError: %s", err, stderrBuf.String())
	}

	// In a goroutine, wait for the command to finish and close the writer
	go func() {
		defer pw.Close()
		if err := cmd.Wait(); err != nil {
			pw.CloseWithError(fmt.Errorf("download failed: %w\nError: %s", err, stderrBuf.String()))
		}
	}()

	return pr, nil
}

func WriteMetadata(song dto.Track, dir string, input io.Reader) error {
	// Build the ffmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-y",           // Overwrite output without asking
		"-i", "pipe:0", // Input from stdin
		"-i", song.Thumbnail, // Input thumbnail
		"-map", "0:a", // Use audio from first input
		"-map", "1:v", // Use video (thumbnail) from second input
		"-disposition:v:0", "attached_pic", // Set thumbnail as cover
		"-c:a", "copy", // TODO: BIG WARNING, IF AUDIO IS NOT OPUS, ADD CODEC CHECKING
		"-vn",
		// "-c:a", "aac", // Use AAC audio codec
		// "-c:v", "copy", // Copy thumbnail without re-encoding
		"-metadata", fmt.Sprintf("title=%s", escapeMetadata(song.Title)),
		"-metadata", fmt.Sprintf("artist=%s", escapeMetadata(song.Authors)),
		"./dl/"+dir+"/"+song.Id+".m4a", // Output file
	)

	// Set stdin to our input reader
	cmd.Stdin = input

	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}

	return nil
}

// escapeMetadata properly escapes special characters in metadata
func escapeMetadata(s string) string {
	// FFmpeg requires special characters to be escaped with backslash or single quotes
	// We'll use single quotes which are safer for most cases
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}

func main() {
	cfg, err := infra.NewConfig()
	if err != nil {
		panic(err)
	}

	pgx, err := pgxpool.New(context.Background(), cfg.DbUrl)
	if err != nil {
		panic(err)
	}

	tr := repo.NewTrackRepo(pgx)
	plService := service.NewPlaylistService(pgx, tr)

	list, err := plService.GetAll(context.Background(), 687627953)
	if err != nil {
		panic(err)
	}

	for _, pl := range list {
		fullPl, err := plService.GetById(context.Background(), pl.Id, 687627953)
		if err != nil {
			panic(err)
		}

		err = downloadPlaylist(context.Background(), fullPl)
		if err != nil {
			panic(err)
		}
	}
}

func downloadPlaylist(ctx context.Context, pl dto.Playlist) error {
	err := os.Mkdir("./dl/"+pl.Title, os.ModePerm)
	if err != nil {
		return err
	}

	var allowedTracks []dto.Track
	for _, t := range pl.Tracks {
		if slices.Contains(pl.AllowedIds, t.Id) {
			allowedTracks = append(allowedTracks, t)
		}
	}

	for _, t := range allowedTracks {
		dl, err := Download(t.Id)
		if err != nil {
			return err
		}
		defer dl.Close()

		err = WriteMetadata(t, pl.Title, dl)
		if err != nil {
			return err
		}
	}

	return nil
}

*/
