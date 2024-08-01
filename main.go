package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/pojntfx/go-nbd/pkg/server"
)

func main() {
	basedir := flag.String("basedir", "./SPARSEBD/", "directory to store the block device chunk files")
	filesuffix := flag.String("filesuffix", ".BDC", "suffix to add to names of tree leaf files")
	treespecStr := flag.String("treespec", "64,64,32", "number of nodes on each level of the tree (the size of the block device = total number of nodes * file size)")
	fileSize := flag.Int64("filesize", 32768, "leaf file size in bytes (consider setting it to a multiple of 512) (the size of the block device = total number of nodes * file size)")
	listenaddr := flag.String("listenaddr", "127.0.0.1:11114", "TCP address to listen on")
	exportName := flag.String("exportname", "sparsebd", "NBD export name")
	exportDescription := flag.String("exportdesc", "sparse storage backed NBD server", "NBD export description")
	loglevelStr := flag.String("loglevel", "info", "logging level (setting to debug might cause slowdowns)")
	flag.Parse()

	loglevel, err := parseLogLevel(*loglevelStr)
	if err != nil {
		panic(fmt.Sprintf("unsupported log level: %v", *loglevelStr))
	}

	loggingSetup(loglevel)

	treespec, err := parseTreeSpec(*treespecStr)
	if err != nil {
		globalLogger.Error("error parsing tree spec", "error", err)
		os.Exit(127)
	}

	l, err := net.Listen("tcp", *listenaddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	globalLogger.Info("server listening", "addr", *listenaddr)

	back := &fileChunkBackend{
		*basedir,
		*filesuffix,
		treespec,
		*fileSize,
		sync.RWMutex{},
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			globalLogger.Warn("could not accept connection", "error", err)
			continue
		}
		globalLogger.Info("client connected", "addr", conn.RemoteAddr())

		go func() {
			defer conn.Close()
			err := server.Handle(
				conn,
				[]*server.Export{
					{
						Name:        *exportName,
						Description: *exportDescription,
						Backend:     back,
					},
				},
				&server.Options{
					ReadOnly:           false,
					MinimumBlockSize:   512,
					PreferredBlockSize: 512,
					MaximumBlockSize:   512,
					SupportsMultiConn:  true,
				})
			if err != nil {
				globalLogger.Error("error encountered in server", "error", err)
				if err == io.EOF {
					globalLogger.Info("error is EOF, did you specify the export name correctly?")
				}
			}
		}()
	}
}
