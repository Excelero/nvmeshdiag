# Description
NVMesh Diag is a tool which can be used to capture hardware and system configuration details required to perform health checks, check for crucial NVMesh software dependencies, verify best practices and to gather vital information for installation and deployment planning.
# Synopsis/Execute the program
./nvmeshdiag
# Operating System Support
All major linux distributions.
# Installation and Execution
Two options are available to easily install and execute nvmeshdiag.
Download and execute the compiled binary, or
Download the source code and compile the executable

## Download the binary:
Go to https://github.com/Excelero/nvmeshdiag/bin and download the binary or clone the git repo with `git clone https://github.com/Excelero/nvmeshdiag` and navigate to the `bin` folder. Verify if the binary is executable and run it by simply typing `./nvmeshdiag`
## Download and compile the source:
First, as NVMesh drag is written in Go, if you want to compile the source code by yourself you need to install Go on your host. Please see the details on how to download and install Go here: https://golang.org
Second, clone the git repository or download the latest source code and compile from the source by changing into the `src` directory and compile it by running 'go build -o nvmeshdiag'. The binary will be located in same folder.
