package main

import (
	"fmt"
	"yangjiafei/langdetect"
)

func main() {
	// Test cases for various languages and Chinese variants
	texts := map[string]string{
		"English": "Hello, world!",
		"Spanish": "¡Hola, mundo!",
		"French": "Bonjour le monde!",
		"Arabic": "مرحبا بالعالم!",
		"Russian": "Привет, мир!",
		"Japanese": "こんにちは、世界！",
		"Korean": "안녕하세요, 세계!",
		"German": "Hallo Welt!",
		"Simplified Chinese": "语言",
		"Traditional Chinese": "語言",
		"Dutch": "Hallo wereld!",		"Thai": "สวัสดีชาวโลก!",
		"Hindi": "नमस्ते दुनिया!",
		"Indonesian": "Halo dunia!",
		"Turkish": "Merhaba dünya!",
		"Slovak": "Ahoj svet!",
	}

	for langName, text := range texts {
		info := langdetect.Detect(text)
		fmt.Printf("Text (%s): '%s' -> Detected: %s, Code: %s\n", langName, text, info.Name, info.Code)
	}
}
