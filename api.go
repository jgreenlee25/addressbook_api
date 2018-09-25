package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TODO: import/export CSV
// TODO: *optional* rewrite Echo API to HTTP API w/ native testing
// TODO: check REST API guidelines to make sure we're returning the correct status codes

// CLASSES/INTERFACES AND STRUCTS
type Controller interface {
	Get(echo.Context) error
	Post(echo.Context) error
	Put(echo.Context) error
	Delete(echo.Context) error
	List(echo.Context) error
}

type Address struct {
	ID    int    `json:"id"`
	First string `json:"first"`
	Last  string `json:"last"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// ADDRESS CONTROLLER
type AddressController struct {
	db  map[int]*Address
	seq int
	Controller
}

func NewAddressController(data map[int]*Address) *AddressController {
	return &AddressController{db: data, seq: 1}
}

func (this *AddressController) Get(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	if val, ok := this.db[id]; ok {
		return c.JSON(http.StatusOK, val)
	}
	return c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (this *AddressController) Post(c echo.Context) (err error) {
	payload := &Address{
		ID: this.seq,
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, http.StatusText(400))
	}

	this.db[this.seq] = payload
	this.seq++

	return c.JSON(http.StatusCreated, payload)
}

func (this *AddressController) Put(c echo.Context) (err error) {
	payload := new(Address)
	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, http.StatusText(400))
	}

	id, _ := strconv.Atoi(c.Param("id"))

	if val, ok := this.db[id]; ok {
		val.First = payload.First
		val.Last = payload.Last
		val.Email = payload.Email
		val.Phone = payload.Phone

	} else {
		return c.JSON(404, http.StatusText(404))
	}

	return c.JSON(http.StatusOK, this.db[id])
}

func (this *AddressController) Delete(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(this.db, id)
	return c.NoContent(http.StatusNoContent)
}

func (this *AddressController) List(c echo.Context) (err error) {
	m := make([]*Address, 0, len(this.db))
	for _, val := range this.db {
		m = append(m, val)
	}
	sort.Slice(m[:], func(i, j int) bool {
		return m[i].ID < m[j].ID
	})
	return c.JSON(http.StatusOK, m)
}

func (this *AddressController) ExportCSV(c echo.Context) (err error) { // w http.ResponseWriter, r *http.Request
	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.

	wr.Write([]string{"First", "Last", "Email", "Phone"}) // write CSV header

	for _, record := range this.db {
		strs := make([]string, 0)
		strs = append(strs, record.First)
		strs = append(strs, record.Last)
		strs = append(strs, record.Email)
		strs = append(strs, record.Phone)
		wr.Write(strs)
	}

	wr.Flush() // writes the csv writer data to the buffered data io writer(b(bytes.buffer))

	c.Response().Header().Set("Content-Type", "text/csv")
	// c.Response().Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
	c.Response().WriteHeader(http.StatusOK)

	c.Response().Write(b.Bytes())
	c.Response().Flush()
	return nil
}

func (this *AddressController) ImportCSV(c echo.Context) (err error) {
	if payload, err := ReadCSVFromHttpRequest(c.Request()); err == nil {
		debug("payload", payload)

	} else {
		return c.JSON(400, http.StatusText(400))
	}

	return c.JSON(http.StatusOK, "Successfully imported CSV data into Address Book.")
}

func ReadCSVFromHttpRequest(req *http.Request) ([][]string, error) {
	reader := csv.NewReader(req.Body)
	var results [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		results = append(results, record)
	}
	return results, nil
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	data := make(map[int]*Address)

	// Routing
	ac := *NewAddressController(data)
	e.GET("/address/:id", ac.Get)
	e.GET("/address", ac.List)
	e.GET("/address/export", ac.ExportCSV)
	e.POST("/address/import", ac.ImportCSV)
	e.POST("/address", ac.Post)
	e.PUT("/address/:id", ac.Put)
	e.DELETE("/address/:id", ac.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}

// HELPERS
func debug(s string, o interface{}) {
	if os, ok := o.(string); ok {
		fmt.Println(s + " = " + os)

	} else {
		fmt.Printf(s+" = %#v\n", o)

		response, _ := json.MarshalIndent(&o, "", "  ")
		fmt.Println(string(response))
	}
}
