package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Logger gestiona el archivo de persistencia
type Logger struct {
	file *os.File
}

// NewLogger abre (o crea) el archivo de logs
func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// Write guarda una operación en disco: "SET key value"
func (l *Logger) Write(op, key, value string) error {
	// Formato simple: OP,KEY,VALUE (separado por comas o espacios)
	// Para simplificar, usaremos un formato de texto línea por línea
	line := fmt.Sprintf("%s,%s,%s\n", op, key, value)
	_, err := l.file.WriteString(line)
	return err
}

// Recover lee el archivo y restaura los datos en el Store
func (l *Logger) Recover(s *Store) error {
	// Volvemos al inicio del archivo
	l.file.Seek(0, 0)
	
	scanner := bufio.NewScanner(l.file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 3) // Dividir en 3 partes máx
		if len(parts) < 2 {
			continue
		}
		
		op := parts[0]
		key := parts[1]
		
		if op == "SET" && len(parts) == 3 {
			s.data[key] = parts[2] // Restauramos directamente (sin mutex porque es el inicio)
		} else if op == "DEL" {
			delete(s.data, key)
		}
	}
	return scanner.Err()
}

func (l *Logger) Close() error {
	return l.file.Close()
}