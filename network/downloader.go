package network

import (
	"bytes"
	"errors"
	"github.com/valyala/fasthttp"
	"net/url"
	"orchestra/utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Helper function
func isSuccessfulConnection(code int) bool {
	return code > 199 && code < 300
}

// Helper function
func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

type (
	task struct {
		start, end uint64
		id         uint
	}

	TaskStatus       uint
	DownloaderStatus uint

	Downloader struct {
		// Public SECTION

		/*
			How many threads (goroutines) should be employed in the procedure.
			0 represents using number of the logical processors.
		*/
		ThreadCount uint
		/*
			The URL which the downloader sends attempts to.
		*/
		URL *url.URL
		/*
			Target file path where the downloader saves the file.
			Leave empty "" if you want to let the program decide automatically.
		*/
		TargetFilePath string

		/*
			Rename the downloaded file if the header of GET request returns "Content-Disposition".
			Set true to rename automatically, otherwise do not
		*/
		AutoRename bool

		/*
			Whether to connect by proxy or not
		*/
		EnableProxy bool

		/*
			Leave empty string "" if you do not want to use proxy
		*/
		ProxyURL string

		// Private SECTION

		/*
			Bytes a thread (goroutine) will download
		*/
		sizePerThread uint64
		/*
			The total length of the target file
		*/
		length uint64
		/*
			The type of the file
		*/
		contentType string
		/*
			true - Support multi-thread downloading, otherwise do not
		*/
		supportPartition bool
		/*
			Bytes the downloader has downloaded so far
		*/
		bytesDownloaded uint64

		requestURI     *fasthttp.URI
		group          sync.WaitGroup
		err            error
		threadUsed     uint
		target         *os.File
		hasDisposition bool
		streamFinished bool
		ready          bool
		Suggested      string // The name of target file

		// Callbacks
		TaskResponses       chan TaskResponse // Channel passing responses
		DownloaderResponses chan DownloaderResponse
	}

	TaskResponse struct {
		Id     uint       // The id of the task
		Status TaskStatus // Task status
	}

	DownloaderResponse struct {
		Status    DownloaderStatus
		Exception error
	}
)

// TaskStatus
const (
	TaskFinished TaskStatus = iota
	TaskProtocolConnected
	TaskStarted
	TaskUpdated
)

// DownloaderStatus
const (
	DownloaderStarted DownloaderStatus = iota
	DownloaderFinished
	DownloaderWaiting
	DownloaderTaskDispatched
	DownloaderTaskCompleted
	DownloaderFileWritten
	DownloaderFileExists
	DownloaderUnableToCreateFile
	DownloaderFailedToTryConnection
)

func (d *Downloader) try() error {
	var err error = nil
	var req = fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodHead)
	req.SetURI(d.requestURI)
	defer fasthttp.ReleaseRequest(req)
	var resp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err = fasthttp.Do(req, resp); err != nil {
		return err
	} else {
		if isSuccessfulConnection(resp.StatusCode()) {
			// Successfully connected
			d.length = uint64(resp.Header.ContentLength())
			d.supportPartition = bytes.
				Equal(utils.ConvertStringToByteSlice("bytes"),
					resp.Header.Peek("Accept-Ranges"))
			d.contentType = string(resp.Header.ContentType())
			dis := resp.Header.Peek("Content-Disposition")
			d.hasDisposition = dis != nil
			if d.hasDisposition {
				vs := utils.Split(dis, ';')
				for _, v := range vs {
					if i := utils.Index(v, '"'); utils.HasPrefix(v, "filename") &&
						i != -1 {
						d.Suggested = utils.ConvertByteSliceToString(v[i+1 : utils.LastIndex(v, '"')])
					}
				}
			} else {
				ur := utils.ConvertStringToByteSlice(d.URL.Path)
				ur = ur[utils.LastIndex(ur, '/')+1:]
				if i := utils.LastIndex(ur, '.'); i == -1 {
					// Cannot find any dots in the filename
					d.Suggested = utils.
						ConvertByteSliceToString(append(ur,
							utils.SuffixOf(d.contentType)...)) // Guess the possible suffix of the file
				} else {
					d.Suggested = utils.ConvertByteSliceToString(ur)
				}
			}
		} else {
			err = errors.New("failed to connect, because the host returned: " +
				strconv.FormatInt(int64(resp.StatusCode()), 10))
			return err
		}
		return nil
	}
}

func (d *Downloader) dispatch(t task) {
	// Start downloading a single chunk
	d.TaskResponses <- TaskResponse{
		Id:     t.id,
		Status: TaskStarted,
	}
	d.DownloaderResponses <- DownloaderResponse{DownloaderTaskDispatched, nil}
	defer func() {
		d.group.Done()
		d.TaskResponses <- TaskResponse{Id: t.id, Status: TaskFinished}
		d.DownloaderResponses <- DownloaderResponse{DownloaderTaskCompleted, nil}
	}()
	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetURI(d.requestURI)
	if d.supportPartition {
		req.Header.SetByteRange(int(t.start), int(t.end))
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if d.err = fasthttp.Do(req, resp); d.err == nil &&
		isSuccessfulConnection(resp.StatusCode()) {
		d.TaskResponses <- TaskResponse{
			Id:     t.id,
			Status: TaskProtocolConnected,
		}
		var cursor = t.start
		n, _ := d.target.WriteAt(resp.Body(), int64(cursor))
		d.TaskResponses <- TaskResponse{Id: t.id, Status: TaskUpdated}
		atomic.AddUint64(&d.bytesDownloaded, uint64(n))
		cursor += uint64(n)
	}
}

// Launch simply call the Start() function, but returns no error.
// To get the error, please call Error() function.
func (d *Downloader) Launch() {
	d.err = d.Init()
	d.err = d.Start()
}

func (d *Downloader) Init() error {
	d.DownloaderResponses = make(chan DownloaderResponse, 8)
	d.requestURI = fasthttp.AcquireURI()
	d.requestURI.SetScheme(d.URL.Scheme)
	d.requestURI.SetHost(fasthttp.AddMissingPort(d.URL.Host, d.URL.Scheme == "https"))
	d.requestURI.SetPath(d.URL.Path)
	if d.ThreadCount == 0 {
		d.ThreadCount = uint(d.length / (16 * 1024)) // 16 KB per goroutine
		if n := uint(runtime.NumCPU()); d.ThreadCount < n {
			d.ThreadCount = n
		}
	}
	if d.EnableProxy {
		// TODO
	}
	if err := d.try(); err != nil {
		d.DownloaderResponses <- DownloaderResponse{DownloaderFailedToTryConnection, err}
		return err
	}
	if fi, err := os.Stat(d.TargetFilePath); err == nil {
		if fi.IsDir() {
			if d.AutoRename {
				d.TargetFilePath = filepath.Join(d.TargetFilePath, d.Suggested)
			} else {
				d.TargetFilePath = filepath.Join(d.TargetFilePath,
					strconv.FormatInt(time.Now().Unix(), 16)+
						".orchestra_temporary")
			}
		} else {
			d.DownloaderResponses <- DownloaderResponse{DownloaderFileExists, err}
			return os.ErrExist
		}
	}
	d.TargetFilePath, _ = filepath.Abs(d.TargetFilePath)
	d.sizePerThread = d.length / uint64(d.ThreadCount)
	if d.supportPartition {
		d.threadUsed = d.ThreadCount
	} else {
		d.threadUsed = 1
	}
	d.TaskResponses = make(chan TaskResponse, d.threadUsed)
	d.ready = true
	return nil
}

func (d *Downloader) Start() error {
	if !d.ready {
		if err := d.Init(); err != nil {
			return err
		}
	}
	if f, err := createFile(d.TargetFilePath); err == nil {
		d.target = f
	} else {
		d.DownloaderResponses <- DownloaderResponse{DownloaderUnableToCreateFile, err}
		return err
	}
	for i := uint(0); i < d.threadUsed; i++ {
		var s = uint64(i) * d.sizePerThread
		var e = s + d.sizePerThread
		if e > d.length {
			e = d.length
		}
		d.group.Add(1)
		go d.dispatch(task{s, e, i})
	}
	d.DownloaderResponses <- DownloaderResponse{DownloaderStarted, nil}
	return nil
}

func (d *Downloader) Wait() {
	d.DownloaderResponses <- DownloaderResponse{DownloaderWaiting, nil}
	d.group.Wait()
	d.Close()
	d.DownloaderResponses <- DownloaderResponse{DownloaderFinished, nil}
}

func (d *Downloader) Error() error {
	return d.err
}

func (d *Downloader) Catch(err error) {
	d.err = err
}

func (d *Downloader) IsCompleted() bool {
	return d.bytesDownloaded == d.length
}

func (d *Downloader) Close() {
	if !d.streamFinished { // Make sure the Close() function can only be called once
		return
	}
	d.streamFinished = true  // Mark the flag
	d.err = d.target.Close() // Close the file
	fasthttp.ReleaseURI(d.requestURI)
	d.DownloaderResponses <- DownloaderResponse{DownloaderFileWritten, d.err}
}

func (d *Downloader) BytesDownloaded() uint64 {
	return d.bytesDownloaded
}

func (d *Downloader) TotalLength() uint64 {
	return d.length
}
