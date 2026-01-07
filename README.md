<table align="center">
  <tr>
    <td>
      <img src="Icon.png" width="64" alt="RegistryMate Icon">
    </td>
    <td>
      <h1>RegistryMate</h1>
    </td>
  </tr>
</table>

<br>

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)
[![Fyne](https://img.shields.io/badge/Fyne-2.6+-orange?logo=go)](https://fyne.io/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/javaLux/registrymate)

A cross-platform GUI tool built with **Go** and **Fyne** to **easily and correctly create Kubernetes ImagePullSecrets**, ensuring accuracy and reducing common errors.<br>
The app focuses on **usability, transparency, and reliability** when working with private container registries and Kubernetes YAML manifests.

---

## Table of Contents

- [Features](#features)
- [Screenshots](#screenshots)
- [Installation](#installation)
  - [Linux (Debian/Ubuntu)](#linux-debianubuntu)
  - [macOS](#macos)
  - [Windows](#windows)
- [Running from Source](#running-from-source)
- [Building from Source](#building-from-source)
- [Usage](#usage-instructions)
- [Contribution](#contributing)
- [License](#license)

## Features

- Create Kubernetes **ImagePullSecrets** as valid YAML
- Copy the generated secret to clipboard or save it to a file
- Built-in **history** for registries and secret metadata
  - Stores up to **100 entries**
  - Oldest entries are automatically removed when the limit is reached
- **Base64 Encode / Decode** utility for Docker-Config JSON string

## Screenshots

<p align="center">
  <img src="../assets/registrymate-light.png" alt="RegistryMate Light-Theme" width="700" />
</p>

<br>

<p align="center">
  <img src="../assets/registrymate-dark.png" alt="RegistryMate Dark-Theme" width="700" />
</p>

## Installation
The easiest way to install is to use the precompiled binaries for your platform.

### ToDo

* [ ] Setup CD-Pipeline
* [ ] Release v1.0
  * [ ] Windows build
  * [ ] macOS build
  * [ ] Linux build


### Linux (Debian/Ubuntu)
- work in progress

### macOS
- work in progress

### Windows
- work in progress

## Running from Source

### Prerequisites

- **Go >= 1.22**
- **Fyne dependencies** see the [Quick-Start-Guide](https://docs.fyne.io/started/quick/)

### Run from Source

```bash
git clone https://github.com/javaLux/registrymate.git
cd registrymate
go mod tidy
go run .
```

The GUI window will open, allowing you to create easily Kubernetes ImagePullSecrets

## Building from Source

### Simple Build

#### Linux

```bash
GOOS=linux GOARCH=amd64 go build
```

#### macOS

```bash
GOOS=darwin GOARCH=amd64 go build
```

#### Windows

```bash
GOOS=windows GOARCH=amd64 go build
```

### Production Build with Packaging

For production-ready packages with icons, metadata, and installer formats:

#### Install Build Tools

```bash
# Install Fyne CLI
go install fyne.io/tools/cmd/fyne@latest

# For Linux packages, also install:
sudo apt-get install dpkg-dev imagemagick wget fuse libfuse2
```

#### Automated Build Script (Linux)

Use the provided build script that creates both `.deb` and `.AppImage` packages:

```bash
# Make the script executable
chmod +x build.sh

# Run the build
./build.sh
```

This will create:
- `dist/registrymate_<Version>_amd64.deb` - Debian package
- `dist/registrymate_<Version>-x86_64.AppImage` - Universal Linux package
- `dist/SHA256SUMS` - Checksums for verification
- `dist/MD5SUMS` - MD5 checksums

#### Manual Packaging with Fyne

##### Linux

```bash
fyne package --os linux --release
```

##### macOS

```bash
fyne package --os darwin --release
```

##### Windows

```bash
fyne package --os windows --release
```

#### Installing Locally (Development)

To install your development build system-wide:

```bash
fyne install --icon Icon.png
```

## Usage

1. Start the application
2. Enter required information:
    - Registry:
      - **URL**
      - **Username**
      - **Password or token**

3. Optional - Secret-Metadata
    - These values must comply with [Kubernetes-Naming-Rules](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/).
      - **Name** -> If not set or invalid, a random name will be generated
      - **Namespace** -> If invalid, it is omitted.

4. Generate the ImagePullSecret by pressing the ► Button or hit ENTER
5. Choose one of the following actions:
   - Copy the YAML to clipboard
   - Save the YAML to a file
6. Use the Base64 Encode / Decode section if you need to inspect or verify the Docker-Config JSON string

Previously used registries and metadata are stored in the history and can be reused quickly.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the **MIT License**.

See the [LICENSE](LICENSE) file for details.
