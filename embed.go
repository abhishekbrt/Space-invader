package main

import (
	"embed"
	rl "github.com/gen2brain/raylib-go/raylib"
	"path/filepath"
)

//go:embed assets/*
var assetsFS embed.FS

// Helper function to get file extension with dot
func getFileExtension(path string) string {
	return filepath.Ext(path) // Returns extension with dot, e.g., ".png"
}

// LoadImageFromEmbedded loads an image from embedded files and returns a pointer to the image
func LoadImageFromEmbedded(path string) *rl.Image {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	ext := getFileExtension(path)
	return rl.LoadImageFromMemory(ext, data, int32(len(data)))
}

// LoadTextureFromEmbedded loads a texture from embedded files
func LoadTextureFromEmbedded(path string) rl.Texture2D {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Get file extension including the dot
	ext := getFileExtension(path)

	// Create image from memory
	image := rl.LoadImageFromMemory(ext, data, int32(len(data)))

	// Load texture from image
	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)

	return texture
}

// LoadSoundFromEmbedded loads a sound from embedded files
func LoadSoundFromEmbedded(path string) rl.Sound {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Get file extension including the dot
	ext := getFileExtension(path)

	var wave rl.Wave
	wave = rl.LoadWaveFromMemory(ext, data, int32(len(data)))

	sound := rl.LoadSoundFromWave(wave)
	rl.UnloadWave(wave)

	return sound
}
