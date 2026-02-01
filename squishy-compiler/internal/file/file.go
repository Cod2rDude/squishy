package file

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    "github.com/Cod2rDude/squishy/squishy-compiler/internal/cli/ui"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/config"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
    "github.com/Cod2rDude/squishy/squishy-compiler/internal/util"
)

// Public Functions
func IsADirectory(path string) (bool, *errors.StackError) {
    fileInfo, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, errors.New(errors.EmptyError, err.Error())
    }

    return fileInfo.IsDir(), nil
}

func IsAValidFile(path string) (bool, *errors.StackError) {
    fileInfo, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, errors.New(errors.PathDoesntExist, path)
        }
        return false, errors.New(errors.EmptyError, err.Error())
    }

    if fileInfo.IsDir() {
        return false, errors.New(errors.PathIsADirectoryNotAFile, path)
    }

    return true, nil
}

func HasAValidExtension(path string, extension string) (bool, *errors.StackError) {
    doesExist, err := IsAValidFile(path)

    if !doesExist {
        return false, err
    }

    doesHaveValidExtension := strings.EqualFold(filepath.Ext(path), extension)

    if doesHaveValidExtension {
        return true, nil
    }

    return false, errors.New(errors.InvalidExtension, filepath.Ext(path), extension)
}

func HasAnyValidExtension(path string, extensions map[string]bool) (bool, *errors.StackError) {
    doesExist, err := IsAValidFile(path)

    if !doesExist {
        return false, err
    }

    ext := strings.ToLower(filepath.Ext(path))

    if extensions[ext] {
        return true, nil
    }

    return false, errors.New(errors.InvalidExtension, filepath.Ext(path), util.ConcatStringIndexedMapToIndexOnlyString(extensions, ", "))
}

func FileToString(path string) (string, *errors.StackError) {
    fileContent, err := os.ReadFile(path)
    if err != nil {
        return "", errors.New(errors.EmptyError, err.Error())
    }

    return string(fileContent), nil
}

func CopyFile(src string, dest string) *errors.StackError {
    if isValidFile, err := IsAValidFile(src); !isValidFile || err != nil {
        return errors.New(errors.SourceFileIsntValid, src)
    }

    if isDirectory, err := IsADirectory(dest); !isDirectory || err != nil {
        return errors.New(errors.DestinationDirectoryIsntValid, dest)
    }

    sourceFile, err := os.Open(src)
    if err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }
    defer sourceFile.Close()

    finalDestPath := filepath.Join(dest, filepath.Base(src))
    destFile, err := os.Create(finalDestPath)
    if err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }
    defer destFile.Close()

    if _, err := io.Copy(destFile, sourceFile); err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    if err := destFile.Sync(); err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    return nil
}

func RenameFile(path string, name string) *errors.StackError {
    if isValidFile, err := IsAValidFile(path); !isValidFile || err != nil {
        return errors.New(errors.SourceFileIsntValid, path)
    }

    newPath := filepath.Join(filepath.Dir(path), name)

    if err := os.Rename(path, newPath); err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    return nil
}

func ReadLine(path string, lineNumber int) (string, *errors.StackError) {
    if isValidFile, err := IsAValidFile(path); !isValidFile || err != nil {
        return "", errors.New(errors.SourceFileIsntValid, path)
    }

    file, err := os.Open(path)
    if err != nil {
        return "", errors.New(errors.EmptyError, err.Error())
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    currentLine := 1

    for scanner.Scan() {
        if currentLine == lineNumber {
            return scanner.Text(), nil
        }
        currentLine++
    }

    if err := scanner.Err(); err != nil {
        return "", errors.New(errors.EmptyError, err.Error())
    }

    return "", errors.New(errors.EmptyError, fmt.Sprintf("Line %d not found in file", lineNumber))
}

func WriteToLine(path string, lineNumber int, newText string) *errors.StackError {
    if isValidFile, err := IsAValidFile(path); !isValidFile || err != nil {
        return errors.New(errors.SourceFileIsntValid, path)
    }

    content, err := os.ReadFile(path)
    if err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    lines := strings.Split(string(content), "\n")

    if lineNumber < 1 || lineNumber > len(lines) {
        return errors.New(errors.EmptyError, fmt.Sprintf("Line %d is out of range", lineNumber))
    }

    lines[lineNumber-1] = newText

    output := strings.Join(lines, "\n")
    if err := os.WriteFile(path, []byte(output), 0644); err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    return nil
}

func CreateAndWriteFile(dir string, name string, content string) *errors.StackError {
    if isDir, err := IsADirectory(dir); !isDir || err != nil {
        return errors.New(errors.DestinationDirectoryIsntValid)
    }

    path := filepath.Join(dir, name)

    if exists , _ := IsAValidFile(path); exists {
        ui.Log(config.GRANDMASTER, "warning", fmt.Sprintf("File '%s' already exists. Deleting...", path))
        if err := os.Remove(path); err != nil {
            return errors.New(errors.EmptyError, err.Error())
        }
    }

    file, err := os.OpenFile(path, os.O_CREATE, 0644)
    if err != nil {
        if os.IsExist(err) {
            ui.Log(config.GRANDMASTER, "warning", fmt.Sprintf("File '%s' already exists. Deleting...", path))
            if err2 := os.Remove(path); err2 != nil {
                return errors.New(errors.EmptyError, err2.Error())
            }
        } else {
            return errors.New(errors.EmptyError, err.Error())
        }
    }
    defer file.Close()

    if _, err := file.WriteString(content); err != nil {
        return errors.New(errors.EmptyError, err.Error())
    }

    return nil
}