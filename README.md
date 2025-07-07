# FileManager - A Go CLI Tool for File Management


## Be Careful with file operations!
```aiignore
Be care when using file operations, especially delete and move commands. Always double-check the paths and filenames to avoid accidental data loss.
```



**FileManager** is a cross-platform command-line tool written in Go for managing files and directories on your operating system. It provides a fast and reliable way to perform file operations such as listing, copying, moving, deleting, and viewing file information — all from your terminal.

---

## 🚀 Features

- List files and directories with detailed metadata
- Delete files and folders
- Create new files or folders
- Search for files by name or extension
- Cross-platform: works on Linux, macOS, and Windows
- Simple, fast, and dependency-free

---

## 🛠 Installation

### Option 1: Build from source

1. Install Go (version 1.18+ recommended): [https://golang.org/dl/](https://golang.org/dl/)
2. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/filemanager.git
   cd filemanager
   go build -o filemanager
   
   ```
### This project is based on the book "Powerful Command-Line Applications in Go" by Ricardo Gerardi.
