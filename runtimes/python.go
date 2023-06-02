package runtimes

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pterm/pterm"
	"github.com/valyala/fasthttp"
	"orchestra/utils"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func FetchPythonVersions() ([]ReleaseVersion, error) {
	b2s := func(b []byte) string {
		return *(*string)(unsafe.Pointer(&b))
	}
	req, uri, resp := fasthttp.AcquireRequest(), fasthttp.AcquireURI(), fasthttp.AcquireResponse()
	_ = uri.Parse(nil, utils.ConvertStringToByteSlice("https://www.python.org/downloads/"))
	defer fasthttp.ReleaseURI(uri)
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)
	if err := fasthttp.Do(req, resp); err == nil {
		if doc, er := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body())); er == nil {
			var v []ReleaseVersion
			doc.Find("#content > div > section > div.row.download-list-widget > ol > li").
				Each(func(i int, selection *goquery.Selection) {
					l, _ := selection.Find("span.release-download > a").Attr("href")
					v = append(v, ReleaseVersion{
						selection.Find("span.release-number > a").Text(),
						selection.Find("span.release-date").Text(),
						nil,
						b2s(append(append(append(uri.Scheme(), "://"...), uri.Host()...), l...)),
						true,
						"md5",
					})
				})
			return v, nil
		} else {
			return nil, er
		}
	} else {
		return nil, err
	}
}

func FetchPythonReleaseFiles(v *ReleaseVersion, filterFunc func(*ReleaseFileVersion) bool) ([]ReleaseFileVersion, error) {
	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)
	if err := uri.Parse(nil, utils.ConvertStringToByteSlice(v.FilesLink)); err == nil {
		req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		req.SetURI(uri)
		req.Header.SetMethod(fasthttp.MethodGet)
		if er := fasthttp.Do(req, resp); er == nil {
			if doc, e := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body())); e == nil {
				var vs []ReleaseFileVersion
				doc.Find("#content > div > section > article > table > tbody > tr").
					Each(func(i int, selection *goquery.Selection) {
						ver := selection.Find("td:nth-child(1) > a")
						n := ver.Text()
						l, _ := ver.Attr("href")
						os := selection.Find("td:nth-child(2)").Text()
						nv := ReleaseFileVersion{
							Name: n,
							Link: l,
							FileType: func() ReleaseFileType {
								if os == "Source release" {
									return TypeSourceCode
								}
								if strings.Contains(n, "installer") {
									return TypeInstaller
								} else {
									return TypeCompressed
								}
							}(),
							SupportedOs: StringToOsType(n),
							Digest:      selection.Find("td:nth-child(4)").Text(),
						}
						if filterFunc != nil && filterFunc(&nv) {
							vs = append(vs, nv)
						}
					})
				return vs, nil
			} else {
				return nil, e
			}
		} else {
			return nil, er
		}
	} else {
		return nil, err
	}
}

func DeployPython() error {
	if vs, err := FetchPythonVersions(); err == nil {
		var data [][]string
		data = append(data, []string{"", "Version", "Release Date"})
		for i, v := range vs {
			data = append(data, []string{strconv.FormatInt(int64(i), 10), v.Version, v.Date})
		}
		_ = pterm.DefaultTable.
			WithBoxed(true).
			WithHasHeader(true).
			WithLeftAlignment(true).
			WithData(data).Render()
		fmt.Print("Please choose the index of Python version you want to download: ")
		var i uint
		_, _ = fmt.Scanln(&i)
		fmt.Println("Searching for", data[i+1][1])
		if fv, er := FetchPythonReleaseFiles(&vs[i+1], func(version *ReleaseFileVersion) bool {
			if version.SupportedOs != StringToOsType(runtime.GOOS) {
				return false
			}
			if !strings.Contains(version.Name, strconv.FormatInt(32<<(^uint(0)>>63), 10)) {
				return false
			}
			return true
		}); er == nil {
			var vd = make([][]string, len(fv)+1)
			vd[0] = []string{"", "Name"}
			for index, value := range fv {
				vd[index+1] = []string{
					strconv.FormatInt(int64(index), 10),
					value.Name,
				}
			}
			_ = pterm.DefaultTable.
				WithBoxed(true).
				WithHasHeader(true).
				WithLeftAlignment(true).
				WithData(vd).Render()
			return nil
		} else {
			return er
		}
	} else {
		return err
	}
}
