// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// file2byteslice is a dead simple tool to embed a file to Go.
package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hajimehoshi/file2byteslice"
)

func genFile(pngFileName string, genFileName string, varName string, spriteSize []int) error {
	var out io.Writer
	var f *os.File
	var err error

	if genFileName != "" {
		f, err = os.Create(genFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var in io.Reader
	if pngFileName != "" {
		f, err := os.Open(pngFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	if err := file2byteslice.Write(out, in, false, "", "asset", varName); err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("var %v_Size = []int{%d,%d} \n", varName, spriteSize[0], spriteSize[1]))
	if err != nil {
		return err
	}

	return nil
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		ext := filepath.Ext(s)
		if ext == ".go" {
			return os.Remove(s)
		} else if ext == ".png" {
			generatedFileName := strings.Replace(s, ".png", ".go", -1)
			varName, spriteSize := genVarNameAndSpriteSize(s)
			fmt.Printf("Found: %v, gen: %v, var: %v\n", s, generatedFileName, varName)
			return genFile(s, generatedFileName, varName, spriteSize)
		}
	}
	return nil
}

// my_image.png -> MyImage_png
func genVarNameAndSpriteSize(pngFileName string) (string, []int) {
	resultStrings := []string{}
	sizePerSprite := []int{}

	{
		spriteSizeComponentRegex, err := regexp.Compile("-\\d+x\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		sizePerSpriteComponent := spriteSizeComponentRegex.FindString(pngFileName)
		spriteSizeRegex, err := regexp.Compile("\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		spriteSizesString := spriteSizeRegex.FindAllString(sizePerSpriteComponent, 2)
		for _, sizeString := range spriteSizesString {
			size, err := strconv.Atoi(sizeString)
			if err != nil {
				log.Fatalln(err)
			}
			sizePerSprite = append(sizePerSprite, int(size))
		}
	}

	extRemoved := strings.Replace(filepath.Base(pngFileName), ".png", "", -1)
	for _, part := range strings.Split(extRemoved, "_") {
		parts := strings.SplitN(part, "", 2)
		resultStrings = append(resultStrings, strings.ToUpper(parts[0]), parts[1])
	}

	return strings.ReplaceAll(strings.Join(resultStrings, ""), "-", "_"), sizePerSprite
}

func main() {
	filepath.WalkDir("./asset", walk)
}
