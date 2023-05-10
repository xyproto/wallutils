# vram

Tool for retrieving the minimum amount of available VRAM for all non-integrated GPUs, or list all available GPUs.

If only integrated GPUs are available, the minimum amount of VRAM for those will be outputted instead.

## Building vram

    go build

## Usage

Retrieve the minimum amount of available VRAM as a number, followed by ` MiB`:

    vram

Listing the bus ID, a description and the available VRAM for all GPUs:

    vram -l

Version information:

    vram --version
