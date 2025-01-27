package routes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func SaveFile(file io.Reader, filename string) (string, error) {
	// Créer le dossier files s'il n'existe pas
	if err := os.MkdirAll("files", 0755); err != nil {
		return "", fmt.Errorf("erreur lors de la création du dossier: %v", err)
	}

	// Créer le chemin complet du fichier
	filePath := filepath.Join("files", filename)

	// Créer le fichier
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier: %v", err)
	}
	defer dst.Close()

	// Copier le contenu
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("erreur lors de la copie du fichier: %v", err)
	}

	return filePath, nil
}
