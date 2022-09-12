# vram

Tool for retrieving the average amount of VRAM across all GPUs, or list all available GPUs and the VRAM for each of them.

## Building vram

    go build

## Usage

Retreive the average VRAM as a number, followed by ` MiB`:

    vram

Listing the VRAM for all available GPUs:

    vram -l

Version information:

    vram --version
