# Digital Organism

### Overview
Digital Organism is designed to illustrate Darwin's evolutionary theory using digital organism.


### Installation
This program is based on go program. Please make sure go is installed before using.

Run ```go build``` to compile.

### Command Line
Example:
```
./DigitalOrganism 200 5000 50 0.001 0.002 0.001 0.002 3 pct.txt length.txt genome.txt
```
* Arg[1]: board size. The size of the board would be Arg[1]* Arg[1]
* Arg[2]: numGen. The program would run Arg[2] generation
* Arg[3]: interval. The gif and report would generate every Arg[3] generations.
* Arg[4]: mutation rate. Indel rate for eukaryotic organisms would be Arg[4]
* Arg[5]: Promutation rate. Indel rate for prokaryotic organisms would be Arg[5]
* Arg[6]: mismatch rate. Mismatch rate for eukaryotic would be Arg[6]
* Arg[7]: Promismatch rate. Mismatch rate for prokaryotic organisms would be Arg[7]
* Arg[8]: Genelength. The initial length of gene of both eukaryotic and prokaryotic organism is Arg[7] times the base gene[E,R,N]. e.g. An input of 2 would make prokaryotic’s initial gene to be [E,R,N, E,R,N], and eukaryotic a pair of [E,R,N, E,R,N].
* Arg[9]: pct file name. The file in which we store stats of percentage of both species’ representation on board.
* Arg[10]: length file name. The file in which we store stats of gene length of both species on board.
* Arg[11]: genome file name. The file in which we store stats of genome content of both species on board.
* Arg[12]: file for reporting the organisms' info with top 10 energy storage.
