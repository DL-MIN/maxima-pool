/*******************************************************************************
 * maxima pool
 *
 * @author     Lars Thoms
 * @date       2022-08-17
 ******************************************************************************/

package main

import (
    "archive/zip"
    "context"
    "errors"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

var hostPtr *string
var portPtr *uint
var timeoutPtr *time.Duration
var tempDirPtr *string
var maximaExecPtr *string


/**
 * Handles maxima server request
 *
 * @param      writer   HTTP response
 * @param      request  HTTP request
 *
 * @return     void
 */
func maximaHandler(writer http.ResponseWriter, request *http.Request) {

    // Check method type
    if request.Method != "POST" {
        http.Error(writer, "Method not supported", http.StatusRequestedRangeNotSatisfiable)
        return
    }

    // Check for valid request
    if err := request.ParseForm(); err != nil {
        http.Error(writer, "Request is invalid", http.StatusRequestedRangeNotSatisfiable)
        return
    }

    // Parse 'input' as maxima code, timeout as milliseconds, and ploturlbase
    code := request.FormValue("input")
    timeout, _ := time.ParseDuration(request.FormValue("timeout") + "ms")
    plotURL := request.FormValue("ploturlbase")

    // Check for empty maxima code
    if len(code) == 0 {
        http.Error(writer, "Input not specified", http.StatusRequestedRangeNotSatisfiable)
        return
    }

    // Set maximum execution time
    if timeout == 0 || timeout > 10 {
        timeout = *timeoutPtr
    }

    // Set plot URL
    if len(plotURL) == 0 {
        plotURL = "!ploturl!"
    }

    // Run maxima
    maximaOut, tempDir, err := runMaxima(code, timeout, plotURL)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusRequestedRangeNotSatisfiable)
        return
    }

    // Package and send result
    sendResult(writer, maximaOut, tempDir)
}


/**
 * Run maxima with given code
 *
 * @param      code     The maxima code
 * @param      timeout  Maximum execution time
 * @param      plotURL  Replacement of a plot url in output
 *
 * @return     Maxima output, path to temporary directory, and an error message
 */
func runMaxima(code string, timeout time.Duration, plotURL string) (string, string, error) {

    // Create temporary directory
    tempDir, err := os.MkdirTemp(*tempDirPtr, "maxima-")
    if err != nil {
        return "", "", fmt.Errorf("I/O error, unable to create temp dir")
    }

    // Create plot directory
    plotDir := filepath.Join(tempDir, "plots")
    if err := os.Mkdir(plotDir, 0750); err != nil {
        return "", tempDir, fmt.Errorf("I/O error, unable to create temp plot dir")
    }

    // Create context with timeout to cancel subprocess
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
 
    // Create a subprocess
    maximaCmd := exec.CommandContext(ctx, *maximaExecPtr, "--quiet")

    // Access stdin of subprocess
    maximaIn, err := maximaCmd.StdinPipe()
    if err != nil {
        return "", tempDir, errors.New("I/O error, stdin expected")
    }
    
    // Access stdout of subprocess
    maximaOut, err := maximaCmd.StdoutPipe()
    if err != nil {
        return "", tempDir, errors.New("I/O error, stdout expected")
    }

    // Start subprocess
    if err := maximaCmd.Start(); err != nil {
        return "", tempDir, errors.New("Exec error, could not start")
    }

    // Pipe maxima code into subprocess
    maximaIn.Write([]byte(fmt.Sprintf("maxima_tempdir:\"%s/\"$\nIMAGE_DIR:\"%s/\"$\nURL_BASE:\"%s\"$\n%s", tempDir, plotDir, plotURL, code)))
    maximaIn.Close()

    // Read output of subprocess
    maximaOutBytes, err := io.ReadAll(maximaOut)
    if err != nil {
        return "", tempDir, errors.New("I/O error, stdout expected")
    }

    // Wait until subprocess exits
    if err := maximaCmd.Wait(); err != nil {
        return "", tempDir, errors.New("Exec error, ungraceful exit")
    }

    // Check for timeout
    if (ctx.Err() == context.DeadlineExceeded) {
        return "", tempDir, errors.New("Exec error, max execution time exeeded")
    }

    // Check for unknown errors
    if (ctx.Err() != nil) {
        return "", tempDir, errors.New("Exec error, unknown error")
    }

    return string(maximaOutBytes), tempDir, nil
}


func sendResult(writer http.ResponseWriter, output string, tempDir string) {

    // Remove temporary directory
    defer os.RemoveAll(tempDir)

    // List all generated plots
    plotDir := filepath.Join(tempDir, "plots")
    plotDirList, err := os.ReadDir(plotDir)
    if err != nil {
        http.Error(writer, "I/O error, could not list plots", http.StatusRequestedRangeNotSatisfiable)
    }

    // No plots, send plaintext only
    if len(plotDirList) == 0 {
        writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")
        writer.Write([]byte(output))

    // Plots exist, send zip archive
    } else {
        writer.Header().Set("Content-Type", "application/zip;charset=UTF-8")

        // Create zip archive
        zipFile := zip.NewWriter(writer)
        defer zipFile.Close()

        // Add OUTPUT file
        zipFileInput, err := zipFile.Create("OUTPUT")
        if err != nil {
            http.Error(writer, "ZIP error, could not add OUTPUT file", http.StatusRequestedRangeNotSatisfiable)
            return
        }
        zipFileInput.Write([]byte(output))

        // Add all plots
        for _, plotFile := range plotDirList {
            plotContent, err := os.Open(filepath.Join(plotDir, plotFile.Name()))
            defer plotContent.Close()

            if err != nil {
                http.Error(writer, "I/O error, could not open plot", http.StatusRequestedRangeNotSatisfiable)
                return
            }

            zipFileInput, err := zipFile.Create("/" + plotFile.Name())
            if err != nil {
                http.Error(writer, "ZIP error, could not add plot", http.StatusRequestedRangeNotSatisfiable)
                return
            }
            io.Copy(zipFileInput, plotContent)
        }
    }
}


/**
 * Main
 *
 * @return     exit code
 */
func main() {

    // CLI options
    hostPtr = flag.String("host", "127.0.0.1", "Bind to specific host or IP")
    portPtr = flag.Uint("port", 8000, "Bind to specific port")
    timeoutPtr = flag.Duration("timeout", 10 * time.Second, "Max execution time of maxima subprocesses")
    tempDirPtr = flag.String("tmp", "/tmp", "Temporary directory for maxima and gnuplot")
    maximaExecPtr = flag.String("maxima", "maxima", "Path to maxima or an optimized maxima snapshot")
    flag.Parse()

    // Handle default route
    http.HandleFunc("/", maximaHandler)

    // Start application server
    fmt.Printf("Starting server at %s:%d\n", *hostPtr, *portPtr)
    fmt.Printf("Max execution time is %s\n", (*timeoutPtr).String())
    if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *hostPtr, *portPtr), nil); err != nil {
        log.Fatal(err)
    }
}
