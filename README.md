# Unique IP Address Counter

This Go project efficiently counts the number of unique IPv4 addresses in a potentially very large text file. The file, `large_ip_file.txt`, can contain millions of IP addresses, and this project is designed to handle such large datasets while optimizing for memory and time efficiency.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Algorithm Explanation](#algorithm-explanation)
- [Performance Considerations](#performance-considerations)
- [Contributing](#contributing)
- [License](#license)

## Overview

Given a large file containing IPv4 addresses, where each line contains one address, this program counts the number of unique IP addresses using a combination of external sorting and merging techniques. The algorithm is designed to handle files that are too large to fit into memory.

## Features

- **Memory Efficient:** Processes large files in chunks to avoid memory overflows.
- **Parallel Processing:** Utilizes Go's goroutines to sort chunks of data concurrently.
- **Scalable:** Can handle files of arbitrary size by managing disk storage effectively.
- **Duplicate Detection:** Efficiently counts unique IP addresses by sorting and merging data.

## Installation

1. **Install Go**: Make sure you have Go installed. You can download it from the [official website](https://golang.org/dl/).

2. **Clone the Repository**:
    ```bash
    git clone https://github.com/sergeyHudzenko/go-unique-ip-counter.git
    cd go-unique-ip-counter
    ```

3. **Initialize the Go Module**:
    ```bash
    go mod init unique-ip-counter
    ```

## Usage

1. **Prepare Your Input File**: Place your large IP file in the project directory. Make sure it's named `large_ip_file.txt` or update the filename in `main.go` accordingly.

2. **Run the Program**:
    ```bash
    go run main.go
    ```

3. **Build the Executable (Optional)**:
    If you prefer to build an executable, run:
    ```bash
    go build -o unique-ip-counter
    ```
    Then execute the program:
    ```bash
    ./unique-ip-counter
    ```

## Algorithm Explanation

This project uses an optimized algorithm to count unique IP addresses:

1. **Chunk Processing**:
    - The file is read in chunks (adjustable size).
    - Each chunk is sorted in memory and then written to a temporary file.

2. **Concurrent Sorting**:
    - Sorting of chunks is done concurrently using Go's goroutines, taking advantage of multi-core processors.

3. **External Merge Sort**:
    - The sorted chunks are merged using a min-heap to efficiently combine the sorted files while counting unique IP addresses.

4. **Memory Efficiency**:
    - By processing the file in manageable chunks and using disk storage for intermediate results, the algorithm minimizes memory usage.

## Performance Considerations

- **Memory Usage**: The program is designed to use as little memory as possible by processing data in chunks and writing intermediate results to disk.
- **Concurrency**: By using goroutines, the program can significantly reduce the time spent on sorting large datasets.
- **Disk I/O**: While the program optimizes memory usage, its performance can still be influenced by disk I/O speeds, especially when handling very large files.

## Contributing

Contributions are welcome! If you find a bug or have an idea for an improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
