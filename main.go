package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sort"
	"time"
	"flag"
	"errors"
	b64 "encoding/base64"
	"strings"
	"text/template"
	"bytes"
)

var filepath102 string = "./testdata/sewobe102.json"

var filepath70 string = "./testdata/sewobe70.json"

var filepath158 string = "./testdata/sewobe158.json"

var aemter_input_simple_example InputAemterKonfig = InputAemterKonfig{
	ExcludeAemter: []string{
		"04 Schatzmeister",
		"23 Redakteur (B)",
		"52 Kantinier (KA)",
		"73 Hüttenrat F",
		"74 Hüttenrat HH",
		"64 Gauobmann HH",
		"71 Hüttenrat Mittelrhein",
		"30 Kantinier (B)",
		"21 Aktiventutor (B)",
		"29 IT-Manager",
		"99 Webteam",
		"31 Assistent der Haus- und Grundstücksverwaltung (B)",
		"76 Hüttenrat KA",
		"77 Hüttenrat Kurpfalz",
		"78 Hüttenrat M",
		"68 Gauobmann M",
		"99 Webteam",
		"09 Vorsitzender der Segelabteilung",
		"99 Webteam",
		"99 Webteam",
		"46 Redakteur (KA)",
		"99 Webteam",
		"95 Archivteam",
		"13 Haus- und Grundstücksverwalter (Karlsruhe)",
		"07 Wissenschaftlicher Koordinator",
		"26 Technik- und Hauswart (HH)",
		"40 Fürsteher",
		"43 Statist (KA)",
		"75 Hüttenrat H",
		"65 Gauobmann H",
		"99 Webteam",
		"14 Vorsitzender der Kassenprüfungskommission",
		"42 Veranstaltungswart (KA)",
		"54 Opi",
		"56 Hauskommission",
		"12 Archivar",
		"27 Segel- und Hauswart (PH)",
		"51 Skihüttenwart",
		"58 Multimediawart",
		"44 Kassierer (KA)",
		"99 Webteam",
		"10 AH Kassierer",
		"99 Webteam",
		"98 SEWOBE Administrator",
		"11 Haus- und Grundstücksverwalter (Berlin)",
		"20 Aktivensprecher",
		"22 Veranstaltungswart (B)",
		"45 Rechnerwart",
		"47 www-Wart",
		"53 Aktivenmitglied im Hüttenrat (KA)",
		"41 Aktiventutor (KA)",
		"70 Hüttenrat B",
		"59 Vertreter der Hüttenakademie (KA)",
		"66 Gauobmann KA",
		"99 Webteam",
	},
	Gruppen: []InputAemterGruppe{
		InputAemterGruppe{
			Titel:        "Gruppe A",
			ZeigeNeuwahl: true,
			Eintraege: []InputAemterEintrag{
				InputAemterEintrag{
					Titel: "Vorstand",
					AemterIDs: []string{
						"15 Aktivenmitglied im Vorstand",
						"01 Vorsitzender des Vorstands",
						"02 Stellvertretender Vorsitzender",
					},
				},

				InputAemterEintrag{
					Titel: "",
					AemterIDs: []string{
						"08 WR-Vorsitzender",
					},
				},
			},
		},
		InputAemterGruppe{
			Titel:        "Gruppe B",
			ZeigeNeuwahl: false,
			Eintraege: []InputAemterEintrag{
				InputAemterEintrag{
					Titel: "Gau-Obmann Frankfurt",
					AemterIDs: []string{
						"63 Gauobmann F",
					},
				},
				InputAemterEintrag{
					Titel: "Hüttenrat Karlrsruhe",
					AemterIDs: []string{
						"76 Hüfcttenrat KA",
					},
				},
				InputAemterEintrag{
					Titel: "Nicht vorhandenes Amt",
					AemterIDs: []string{
						"101 Nicht vorhanden",
					},
				},
			},
		},
	},
}

func configToB64(konfig InputAemterKonfig) string {
	jsonbytes, _ := json.Marshal(konfig)
	var b []byte = make([]byte, b64.RawURLEncoding.EncodedLen(len(jsonbytes)))
	b64.RawURLEncoding.Encode(b, jsonbytes)
	return string(b)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[WIP] API Gateway für Sewobe Datenbank\nAktuell: \n./geburtstage (enthält json, csv, und tex Version für beliebige Intervalle)\n./aemterliste (Ämterliste für HüMi etc.)\n./mitgliederstatistik (Mitgliederstatistik für HüMi etc.)")
}

func geburtstagshandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()
	var vj, vm, vt, bj, bm, bt int
	var format = "json"
	fromstr, fromstrexist := query["von"]
	if fromstrexist {
		von, e1 := time.Parse(time.RFC3339, fromstr[0]+"T00:00:00+00:00")
		if e1 != nil {
			fromstrexist = false
		} else {
			vj = von.Year()
			vm = int(von.Month())
			vt = von.Day()
		}
	}
	tostr, tostrexist := query["bis"]
	if tostrexist {
		bis, e1 := time.Parse(time.RFC3339, tostr[0]+"T00:00:00+00:00")
		if e1 != nil {
			tostrexist = false
		} else {
			bj = bis.Year()
			bm = int(bis.Month())
			bt = bis.Day()
		}
	}
	formatstr, formatstrexist := query["format"]
	if formatstrexist {
		format = formatstr[0]
	}
	if !fromstrexist {
		io.WriteString(w, "Es fehlt ein Query Parameter in der form 'von=YYYY-MM-DD', z.B. '...geburtstage?von=2011-03-16")
	} else if !tostrexist {
		io.WriteString(w, "Es fehlt ein Query Parameter in der form 'bis=YYYY-MM-DD', z.B. '...geburtstage?von=2011-03-16&bis=2012-01-14")
	} else {
		io.WriteString(w, geburtstagsliste(filepath102, format, vt, vm, vj, bt, bm, bj))
	}
}

func aemterlistehandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query()
	var format = "json"
	var konfigstring = ""
	konfigstr, konfigstrexist := query["konfiguration"]
	formatstr, formatstrexist := query["format"]

	if (formatstrexist) {
		format = formatstr[0]
	}
	if (konfigstrexist) {
		konfigstring = konfigstr[0]
	}
	if (len(konfigstring) == 0) {
		jsonbytes, _ := json.MarshalIndent(aemter_input_simple_example, "", "    ")
		var jsonstr = string(jsonbytes)
		io.WriteString(w, "Es fehlt ein Query Parameter 'konfiguration' der Form (als url encoded base64 string ohne padding ala '=' am Ende, z.B. per http://kjur.github.io/jsjws/tool_b64uenc.html):\n\n"+jsonstr+"\n\n\nAls Base64:\n"+"aemterliste?format=tex&konfiguration="+configToB64(aemter_input_simple_example))
	} else {
		inputjson, inputjsonerr := aemter_input_parser(konfigstring)
		if (inputjsonerr != nil) {
			io.WriteString(w, inputjsonerr.Error())
		} else {
			result, err := aemter_aggregator(inputjson, filepath70, format)
			if (err != nil) {
				io.WriteString(w, err.Error())
			} else {
				io.WriteString(w, result)
			}
		}
	}
}

var current_aemterliste_json_config string = `
{
  "ExcludeAemter": [
    "99 Webteam",
    "95 Archivteam",
    "12 Archivar",
    "97 Stiftungsfestkommission",
    "98 SEWOBE Administrator"
  ],
  "Gruppen": [
    {
      "Titel": "Altherrenschaft",
      "ZeigeNeuwahl": true,
      "Eintraege": [
        {
          "AemterIDs": [
            "01 Vorsitzender des Vorstands"
          ],
          "Titel": "Vorsitzender des Vorstands"
        },
        {
          "AemterIDs": [
            "02 Stellvertretender Vorsitzender",
            "08 WR-Vorsitzender",
            "15 Aktivenmitglied im Vorstand"
          ],
          "Titel": "Stellvertretende Vorsitzende"
        },
        {
          "AemterIDs": [
            "04 Schatzmeister"
          ],
          "Titel": "Schatzmeister"
        },
        {
          "AemterIDs": [
            "07 Wissenschaftlicher Koordinator"
          ],
          "Titel": "Wissenschaftlicher Koordinator"
        },
        {
          "AemterIDs": [
            "08 WR-Vorsitzender"
          ],
          "Titel": "WR-Vorsitzender"
        },
        {
          "AemterIDs": [
            "09 Vorsitzender der Segelabteilung"
          ],
          "Titel": "Vorsitzender der Segelabteilung"
        },
        {
          "AemterIDs": [
            "10 AH Kassierer"
          ],
          "Titel": "AH Kassierer"
        },
        {
          "AemterIDs": [
            "11 Haus- und Grundstücksverwalter (Berlin)"
          ],
          "Titel": "Haus- und Grundstücksverwalter (Berlin)"
        },
        {
          "AemterIDs": [
            "13 Haus- und Grundstücksverwalter (Karlsruhe)"
          ],
          "Titel": "Haus- und Grundstücksverwalter (Karlsruhe)"
        },
        {
          "AemterIDs": [
            "12 Archivar"
          ],
          "Titel": "Archivar"
        },
        {
          "AemterIDs": [
            "14 Vorsitzender der Kassenprüfungskommission"
          ],
          "Titel": "Vorsitzender der Kassenprüfungskommission"
        }
      ]
    },
    {
      "Titel": "Vertreter der Hüttengaue im Hüttenrat",
      "ZeigeNeuwahl": true,
      "Eintraege": [
        {
          "AemterIDs": [
            "70 Hüttenrat B"
          ],
          "Titel": "Berlin"
        },
        {
          "AemterIDs": [
            "71 Hüttenrat Mittelrhein"
          ],
          "Titel": "Mittelrhein"
        },
        {
          "AemterIDs": [
            "61 Hüttenrat Ruhr"
          ],
          "Titel": "Ruhr"
        },
        {
          "AemterIDs": [
            "73 Hüttenrat F"
          ],
          "Titel": "Frankfurt"
        },
        {
          "AemterIDs": [
            "74 Hüttenrat HH"
          ],
          "Titel": "Hamburg"
        },
        {
          "AemterIDs": [
            "75 Hüttenrat H"
          ],
          "Titel": "Hannover"
        },
        {
          "AemterIDs": [
            "76 Hüttenrat KA"
          ],
          "Titel": "Karlsruhe"
        },
        {
          "AemterIDs": [
            "77 Hüttenrat Kurpfalz"
          ],
          "Titel": "Kurpfalz"
        },
        {
          "AemterIDs": [
            "78 Hüttenrat M"
          ],
          "Titel": "München"
        },
        {
          "AemterIDs": [
            "79 Hüttenrat S"
          ],
          "Titel": "Stuttgart"
        }
      ]
    },
    {
      "Titel": "Liste der Gauobleute",
      "ZeigeNeuwahl": false,
      "Eintraege": [
        {
          "AemterIDs": [
            "60 Gauobmann B"
          ],
          "Titel": "Berlin"
        },
        {
          "AemterIDs": [
            "61 Gauobmann Mittelrhein"
          ],
          "Titel": "Mittelrhein"
        },
        {
          "AemterIDs": [
            "62 Gauobmann R"
          ],
          "Titel": "Ruhr"
        },
        {
          "AemterIDs": [
            "63 Gauobmann F"
          ],
          "Titel": "Frankfurt"
        },
        {
          "AemterIDs": [
            "64 Gauobmann HH"
          ],
          "Titel": "Hamburg"
        },
        {
          "AemterIDs": [
            "65 Gauobmann H"
          ],
          "Titel": "Hannover"
        },
        {
          "AemterIDs": [
            "66 Gauobmann KA"
          ],
          "Titel": "Karlsruhe"
        },
        {
          "AemterIDs": [
            "67 Gauobmann Kurpfalz"
          ],
          "Titel": "Kurpfalz"
        },
        {
          "AemterIDs": [
            "68 Gauobmann M"
          ],
          "Titel": "München"
        },
        {
          "AemterIDs": [
            "69 Gauobmann S"
          ],
          "Titel": "Stuttgart"
        }
      ]
    },






    {
      "Titel": "Aktivitas Karlsruhe",
      "ZeigeNeuwahl": true,
      "Eintraege": [
        {
          "AemterIDs": [
            "40 Fürsteher"
          ],
          "Titel": "Fürsteher"
        },
        {
          "AemterIDs": [
            "44 Kassierer (KA)"
          ],
          "Titel": "Kassierer"
        },
        {
          "AemterIDs": [
            "42 Veranstaltungswart (KA)"
          ],
          "Titel": "Veranstaltungswart"
        },
        {
          "AemterIDs": [
            "41 Aktiventutor (KA)"
          ],
          "Titel": "Aktiventutor"
        },
        {
          "AemterIDs": [
            "44 Kassierer (KA)",
            "56 Hauskommission"
          ],
          "Titel": "Hauskommission"
        },
        {
          "AemterIDs": [
            "50 Hauswart (KA)"
          ],
          "Titel": "Hauswart"
        },
        {
          "AemterIDs": [
            "51 Skihüttenwart"
          ],
          "Titel": "Skihüttenwart"
        },
        {
          "AemterIDs": [
            "52 Kantinier (KA)"
          ],
          "Titel": "Kantinier"
        },
        {
          "AemterIDs": [
            "46 Redakteur (KA)"
          ],
          "Titel": "Redakteur"
        },
        {
          "AemterIDs": [
            "43 Statist (KA)"
          ],
          "Titel": "Statist"
        },
        {
          "AemterIDs": [
            "49 Bibliothekswart (KA)"
          ],
          "Titel": "Bibliothekswart"
        },
        {
          "AemterIDs": [
            "54 Opi"
          ],
          "Titel": "Opi"
        },
        {
          "AemterIDs": [
            "47 www-Wart"
          ],
          "Titel": "www-Wart"
        },
        {
          "AemterIDs": [
            "58 Multimediawart"
          ],
          "Titel": "Multimediawart"
        },
        {
          "AemterIDs": [
            "45 Rechnerwart"
          ],
          "Titel": "Rechnerwart"
        },
        {
          "AemterIDs": [
            "59 Vertreter der Hüttenakademie (KA)"
          ],
          "Titel": "Hüttenakademie"
        },
        {
          "AemterIDs": [
            "53 Aktivenmitglied im Hüttenrat (KA)"
          ],
          "Titel": "Mitglied im Hüttenrat"
        }
      ]
    },
    {
      "Titel": "Aktivitas Berlin",
      "ZeigeNeuwahl": true,
      "Eintraege": [
        {
          "AemterIDs": [
            "20 Aktivensprecher"
          ],
          "Titel": "Aktivensprecher"
        },
        {
          "AemterIDs": [
            "21 Aktiventutor (B)"
          ],
          "Titel": "Aktiventutor"
        },
        {
          "AemterIDs": [
            "26 Technik- und Hauswart (HH)"
          ],
          "Titel": "Haus- und Hofwart HH"
        },
        {
          "AemterIDs": [
            "27 Segel- und Hauswart (PH)"
          ],
          "Titel": "Segel- und Pichelhüttenwart"
        },
        {
          "AemterIDs": [
            "24 Statist (B)"
          ],
          "Titel": "Statist"
        },
        {
          "AemterIDs": [
            "23 Redakteur (B)"
          ],
          "Titel": "Redakteur"
        },
        {
          "AemterIDs": [
            "25 Kassierer (B)"
          ],
          "Titel": "Kassierer"
        },
        {
          "AemterIDs": [
            "22 Veranstaltungswart (B)"
          ],
          "Titel": "Veranstaltungswart"
        },
        {
          "AemterIDs": [
            "34 Bibliothekswart (B)"
          ],
          "Titel": "Bibliothekswart"
        },
        {
          "AemterIDs": [
            "31 Assistent der Haus- und Grundstücksverwaltung (B)"
          ],
          "Titel": "Verantwortlicher der Vermietungen"
        },
        {
          "AemterIDs": [
            "29 IT-Manager"
          ],
          "Titel": "IT-Manager"
        },
        {
          "AemterIDs": [
            "28 Aktivenmitglied in der Kassenprüfungskommission"
          ],
          "Titel": "Mitglied in der Kassenprüfung"
        },
        {
          "AemterIDs": [
            "30 Kantinier (B)"
          ],
          "Titel": "Kantinier"
        },
        {
          "AemterIDs": [
            "32 Aktivenmitglied im Hüttenrat (B)"
          ],
          "Titel": "Mitglied im Hüttenrat"
        }
      ]
    }
  ]
}

`

var aemterliste_html = `
<html>
    <head>
    <title></title>
    </head>
    <body>
    	<h1>AVH Ämterliste Generator</h1>
    	<p>Du kannst hier das Format auswählen und die Konfiguration anpassen, bevor du das Ergebnis erhälst.</p>

    	<form action="" id="aemterform" method="post">
  			Format: <select name="format">
				<option value="tex">LaTeX</option>
				<option value="json">JSON</option>
				<option value="csv">CSV</option>
			  </select>
  			<input type="submit">
		</form>
		<br>
		<textarea rows="30" cols="50" name="json" form="aemterform">
		`+ current_aemterliste_json_config + `
		</textarea>
    </body>
</html>
`

func aemter_input_parser_withoutb64(decoded string) (InputAemterKonfig, error) {
	var k InputAemterKonfig
	if err := json.Unmarshal([]byte(decoded), &k); err != nil {
		return InputAemterKonfig{}, err
	} else {
		return k, nil
	}
}

func aemterliste_v2_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, aemterliste_html)
	} else {
		r.ParseForm()
		inputjson, inputjsonerr := aemter_input_parser_withoutb64(r.Form["json"][0])
		if (inputjsonerr != nil) {
			io.WriteString(w, inputjsonerr.Error())
		} else {
			result, err := aemter_aggregator(inputjson, filepath70, r.Form["format"][0])
			if (err != nil) {
				io.WriteString(w, err.Error())
			} else {
				io.WriteString(w, result)
			}
		}
	}
}

func mitgliederstatistikhandler(w http.ResponseWriter, r *http.Request) {

	var query = r.URL.Query()
	var format = "json"
	var stichtag Datum

	stichtagstr, stichtagstrexist := query["stichtag"]
	formatstr, formatstrexist := query["format"]
	if stichtagstrexist {
		stichtagtemp, stichtagerr := datum_parse(stichtagstr[0])
		if stichtagerr != nil {
			io.WriteString(w, stichtagerr.Error())
			return
		} else {
			stichtag = stichtagtemp
		}
	} else {
		stichtag = datum_jetzt()
	}
	if (formatstrexist) {
		format = formatstr[0]
	}
	result, err := mitgliederstatistik_aggregator(stichtag, filepath158, format)
	if (err != nil) {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, result)
	}
}

func main() {
	flag.StringVar(&filepath102, "filepath102", filepath102, "a filepath to sewobe102.json file")
	flag.StringVar(&filepath70, "filepath70", filepath70, "a filepath to sewobe70.json file")
	flag.StringVar(&filepath158, "filepath158", filepath158, "a filepath to sewobe158.json file")
	portPtr := flag.Int("port", 8080, "the port this http server is running under")
	flag.Parse()
	port := fmt.Sprint(":", *portPtr)

	http.HandleFunc("/", hello)
	http.HandleFunc("/geburtstage", geburtstagshandler)
	http.HandleFunc("/aemterliste", aemterlistehandler)
	http.HandleFunc("/aemterliste_v2", aemterliste_v2_handler)
	http.HandleFunc("/mitgliederstatistik", mitgliederstatistikhandler)

	http.ListenAndServe(port, nil)
}

func geburtstagsliste(filepath string, fmttype string, nachTag int, nachMonat int, nachJahr int, vorTag int, vorMonat int, vorJahr int) string {
	raw, err := readAndConvert(filepath)
	if (err != nil) {
		return err.Error()
	}
	//fmt.Println("Raw:\n", raw)
	var filtered []Datarow = filter(raw, nachTag, nachMonat, nachJahr, vorTag, vorMonat, vorJahr)
	//fmt.Println("Filtered:\n", filtered)
	var list []DatarowAlter = berechneAlter(filtered, nachTag, nachMonat, nachJahr)

	//fmt.Println("Altered:\n", list)
	if fmttype == "json" {
		jsonbytes, err := json.Marshal(list)
		if (err != nil) {
			abort("json encoding did not work")
			return ""
		} else {
			return string(jsonbytes)
		}
	} else if fmttype == "json-pretty" {
		jsonbytes, err := json.MarshalIndent(list, "", "    ")
		if (err != nil) {
			abort("json encoding did not work")
			return ""
		} else {
			return string(jsonbytes)
		}
	} else if fmttype == "tex" {
		return convertLatex(list)
	} else if fmttype == "csv" {
		return convertCSV(list)
	} else {
		return "format of output not recognized, choose json, tex, or csv"
	}
}

//INCLUDES sowohl Anfang als auch Ende
func filter(rows []Datarow, nachTag int, nachMonat int, nachJahr int, vorTag int, vorMonat int, vorJahr int) []Datarow {
	var p []Datarow
	//println(nachTag, ".", nachMonat, ".", nachJahr, " bis ", vorTag, ".", vorMonat, ".", vorJahr)

	for j := nachJahr; j <= vorJahr; j++ {

		//fmt.Println("Jahr: ", j)
		var beginmonth int = 1
		if j == nachJahr {
			beginmonth = nachMonat
		}
		var endmonth int = 12
		if j == vorJahr {
			endmonth = vorMonat
		}
		for m := beginmonth; m <= endmonth; m++ {
			//fmt.Println("Monat: ", m, "Jahr: ", j)
			var begintag int = 1;
			if j == nachJahr && m == nachMonat {
				begintag = nachTag
			}
			var endtag int = 31;
			if j == vorJahr && m == vorMonat {
				endtag = vorTag
			}
			for d := begintag; d <= endtag; d++ {
				//fmt.Println("Tag: ", d, "Monat: ", m, "Jahr: ", j)
				for i := 0; i < len(rows); i++ {
					var e Datarow = rows[i]
					if (e.Tag == d && e.Monat == m) {
						//fmt.Println("appended:", e)
						p = append(p, e)
					} else {
						//fmt.Println("did not append:", e)
					}
				}
			}
		}
	}
	return p
}

//TODO: bug: wenn jemand doppelt erscheint (interval >> 1 jahr), dann wird das Alter nur logischerweise einmal berechnet
func berechneAlter(rows []Datarow, nachTag int, nachMonat int, nachJahr int) []DatarowAlter {
	var p []DatarowAlter
	for i := 0; i < len(rows); i++ {
		var alter int = -1
		if rows[i].Monat > nachMonat {
			//danach
			alter = nachJahr - rows[i].Jahr + 0
		} else if rows[i].Monat == nachMonat && rows[i].Tag >= nachTag {
			//danach
			alter = nachJahr - rows[i].Jahr + 0
		} else {
			//davor
			alter = nachJahr - rows[i].Jahr + 1
		}
		p = append(p, DatarowAlter{
			Tag:      rows[i].Tag,
			Monat:    rows[i].Monat,
			Jahr:     rows[i].Jahr,
			Vorname:  rows[i].Vorname,
			Biername: rows[i].Biername,
			Nachname: rows[i].Nachname,
			Alter:    alter,
		})
	}
	return p
}

func convertLatex(rows []DatarowAlter) string {
	var pre string = `% Geburtstagliste hier einfügen
\newpage
\phantomsection
\addcontentsline{toc}{section}{Geburtstage}
\markright{Geburtstage}
{\Huge Geburtstagsliste}

\begin{longtable}[width=\textwidth]{p{0.2\textwidth} p {0.2\textwidth} p{0.2\textwidth}p{0.2\textwidth}p{0.2\textwidth}}
\textbf{Geburtsdatum} & \textbf{Vorname} & \textbf{Biername} & \textbf{Nachname} & \textbf{Alter} \\
`
	var post string = `




\end{longtable}
\newpage`

	var str string = pre
	for i := 0; i < len(rows); i++ {
		str += fmt.Sprintf("%02d.%02d.%04d & %s & %s & %s & %d \\\\\n", rows[i].Tag, rows[i].Monat, rows[i].Jahr, rows[i].Vorname, rows[i].Biername, rows[i].Nachname, rows[i].Alter)
	}
	str += post
	return str
}

func convertCSV(rows []DatarowAlter) string {
	var str string = "Geburtstag, Vorname, Biername, Nachname, neues Alter\n"
	for i := 0; i < len(rows); i++ {
		str += fmt.Sprintf("%02d.%02d.%04d,%s,%s,%s,%d\n", rows[i].Tag, rows[i].Monat, rows[i].Jahr, rows[i].Vorname, rows[i].Biername, rows[i].Nachname, rows[i].Alter)
	}
	return str
}

func aemter_input_parser(queryvar string) (InputAemterKonfig, error) {
	decoded, b64err := b64.RawURLEncoding.DecodeString(queryvar)
	if (b64err != nil) {
		return InputAemterKonfig{}, b64err
	} else {
		var k InputAemterKonfig
		if err := json.Unmarshal([]byte(decoded), &k); err != nil {
			return InputAemterKonfig{}, err
		} else {
			return k, nil
		}
	}
}

func aemter_aggregator(konfig InputAemterKonfig, sewobe_filepath string, format string) (string, error) {
	konstellation, err := aemter_data_picker(konfig, sewobe_filepath)
	if (format == "json") {
		if (err != nil) {
			return "", err
		}
		jsonbytes, _ := json.MarshalIndent(konstellation, "", "    ")
		return string(jsonbytes), nil
	} else if (format == "tex") {
		if (err != nil) {
			return "", err
		}
		return aemter_tex_formatter(konstellation), nil
	} else if (format == "csv") {
		if (err != nil) {
			return "", err
		}
		return aemter_csv_formatter(konstellation, ","), nil
	} else {
		return "", errors.New("Fehler: Query Parameter 'format' sollte 'json', 'tex', oder 'json' sein. Diesmal war es '" + format + "'.")
	}

}

func aemter_csv_formatter(konstellation AemterKonstellation, sep string) string {
	var lines []string
	if (len(konstellation.UnzugeordneteAemterIDs) > 0) {
		lines = append(lines, "UNZUGEORDNETE AEMTER!!! Kontrolliere deine Konfiguration:")
		lines = append(lines, strings.Join(konstellation.UnzugeordneteAemterIDs, ","))
	}
	lines = append(lines, "Gruppe,ZeigeNeuwahl,BlockTitel,AemterID,Vorname,Biername,Nachname,Wahlveranstaltung,Wahlsemester,Wahlsemester,Vakant")
	for _, gruppe := range konstellation.Gruppen {
		for _, block := range gruppe.Bloecke {
			for _, zeile := range block.Zeilen {
				lines = append(lines, strings.Join([]string{
					gruppe.Titel,
					strconv.FormatBool(gruppe.ZeigeNeuwahl),
					block.Titel,
					zeile.AemterID,
					zeile.Vorname,
					zeile.Biername,
					zeile.Nachname,
					zeile.Wahlveranstaltung,
					zeile.Wahlsemester,
					strconv.FormatBool(zeile.Vakant),
				}, sep))
			}
		}
	}
	return strings.Join(lines, "\n");
}

func aemter_tex_formatter(konstellation AemterKonstellation) string {
	var lines []string
	var prelogueKonstellation = `\newpage
\phantomsection
\addcontentsline{toc}{section}{Ämterliste}
\markright{Ämterlisten}
{\noindent\Huge Ämter des A.V. Hütte}\\`
	lines = append(lines, prelogueKonstellation)

	if (len(konstellation.UnzugeordneteAemterIDs) > 0) {
		lines = append(lines, "Unzugeordnete Ämter: "+strings.Join(konstellation.UnzugeordneteAemterIDs, " - - ")+" \\")
	}

	var prelogueGruppe1a = `\needspace{.25\textheight}
\noindent
\textbf{`
	var prelogueGruppe1b = `}\\
\begin{tabular}{p{0.30\textwidth} p{0.35\textwidth} p{0.23\textwidth}}
& & Neuwahl: \\ \hline`
	var prelogueGruppe2a = `\needspace{.25\textheight}
\noindent
\textbf{`
	var prelogueGruppe2b = `}\\
\begin{tabular}{p{0.30\textwidth} p{0.35\textwidth} p{0.23\textwidth}}
& & \\ \hline`

	var postlogueGruppe = `\end{tabular}\\
\newssep`

	for _, gruppe := range konstellation.Gruppen {
		if (gruppe.ZeigeNeuwahl) {
			lines = append(lines, prelogueGruppe1a+gruppe.Titel+prelogueGruppe1b)
		} else {
			lines = append(lines, prelogueGruppe2a+gruppe.Titel+prelogueGruppe2b)
		}
		for _, block := range gruppe.Bloecke {
			for index_zeile, zeile := range block.Zeilen {
				var title string = "TestTitle"
				if (index_zeile > 0) {
					title = ""
				} else {
					if (len(block.Titel) > 0) {
						title = block.Titel
					} else {
						title = block.Zeilen[0].AemterID //TODO: remove digits from auto id
					}
				}
				var neuwahl string = zeile.Wahlveranstaltung + " " + zeile.Wahlsemester
				if (!gruppe.ZeigeNeuwahl) {
					neuwahl = ""
				}
				var vorname string = zeile.Vorname
				if (zeile.Vakant) {
					vorname = "vakant"
				}

				var newstring = strings.Join([]string{title, " & "}, "")
				if (zeile.Vakant) {
					newstring += "vakant"
				} else {
					newstring += vorname + " (\\textbf{" + zeile.Biername + "}) " + zeile.Nachname
				}
				newstring += " & " + neuwahl + "\\\\"

				lines = append(lines, newstring)
			}
		}
		lines = append(lines, postlogueGruppe)
	}

	return strings.Join(lines, "\n")
}

func aemter_data_picker(konfig InputAemterKonfig, sewobe_filepath string) (AemterKonstellation, error) {

	//TODO bug: aemter werden direkt aus liste gelöscht nach erstem Eintrag. Manche Ämter kommen aber mehrmals vor
	//=> lösche nicht aus Liste, sondern führe eine zweite

	sewobe, sewobeerr := readAndParseSewobe(sewobe_filepath)
	if (sewobeerr != nil) {
		return AemterKonstellation{}, sewobeerr
	}
	var people []AemterInternalLine = aemter_splitup(sewobe)
	var peopleList []AemterInternalLine = aemter_splitup(sewobe)

	var konstellation = AemterKonstellation{}

	//for each in konfig, suche passende Kandidaten (und appende sie alle)
	for gid, gruppe := range konfig.Gruppen {
		konstellation.Gruppen = append(konstellation.Gruppen, AemterGruppe{
			Titel:        gruppe.Titel,
			ZeigeNeuwahl: gruppe.ZeigeNeuwahl,
		})
		for eid, eintrag := range gruppe.Eintraege {
			konstellation.Gruppen[gid].Bloecke = append(konstellation.Gruppen[gid].Bloecke, AemterBlock{
				Titel: eintrag.Titel,
			})
			for _, aemterid := range eintrag.AemterIDs {
				//falls keine InternalLine mit aemterid verfuegbar, setze vakant auf true
				var vak bool = true
				for i := 0; i < len(people); i++ {
					var person = people[i]
					if (person.AemterID == aemterid) {
						//append alle, die passen
						konstellation.Gruppen[gid].Bloecke[eid].Zeilen = append(konstellation.Gruppen[gid].Bloecke[eid].Zeilen, AemterZeile{
							aemterid,
							person.Vorname,
							person.Biername,
							person.Nachname,
							person.Wahlveranstaltung,
							person.Wahlsemester,
							false,
						})
						vak = false
						//alle genommenen Kandidaten werden aus dem Splice entfernt:

						for k := 0; k < len(peopleList); k++ {
							if peopleList[k].AemterID == aemterid {
								peopleList = append(peopleList[:k], peopleList[k+1:]...) //delete element
								//slice is now one element smaller
								break
							}
						}
					}
				}
				if (vak) {
					konstellation.Gruppen[gid].Bloecke[eid].Zeilen = append(konstellation.Gruppen[gid].Bloecke[eid].Zeilen, AemterZeile{
						aemterid,
						"",
						"",
						"",
						"",
						"",
						true,
					})
				}

			}
		}
	}

	//alle übrigen wandern in das "übrig" splice anhand ihrer AemterID
	for _, u := range peopleList {
		//filtern von excluded
		var found = false
		for _, s := range konfig.ExcludeAemter {
			if (strings.TrimSpace(s) == strings.TrimSpace(u.AemterID)) {
				found = true
			}
		}
		if (!found) {
			konstellation.UnzugeordneteAemterIDs = append(konstellation.UnzugeordneteAemterIDs, strings.TrimSpace(u.AemterID))
		}
	}
	sort.Strings(konstellation.UnzugeordneteAemterIDs)
	return konstellation, nil
}

func aemter_splitup(input map[string]Sewoberaw) []AemterInternalLine {
	list := make([]AemterInternalLine, 0)
	for _, value := range input {
		aemterstrings := strings.Split(value.DATENSATZ["AMT"], ",")
		for _, amt := range aemterstrings {
			list = append(list, AemterInternalLine{
				AemterID:          strings.TrimSpace(amt),
				Vorname:           value.DATENSATZ["VORNAME-PRIVATPERSON"],
				Biername:          value.DATENSATZ["BIERNAME"],
				Nachname:          value.DATENSATZ["NACHNAME-PRIVATPERSON"],
				Wahlveranstaltung: value.DATENSATZ["NEUWAHL"],
				Wahlsemester:      value.DATENSATZ["JAHR"],
			})
		}

	}
	return list
}

func calculate_age(geburtstag, stichtag Datum) int {
	//entweder ist der monat davor, oder der monat ist gleich und der tag ist davor
	var davor bool = (stichtag.Monat < geburtstag.Monat) || ((geburtstag.Monat == stichtag.Monat) && (stichtag.Tag < geburtstag.Tag))
	if (davor) {
		return stichtag.Jahr - geburtstag.Jahr - 1
	} else {
		return stichtag.Jahr - geburtstag.Jahr - 0
	}
}

func mitgliederstatistik_csv_formatter(statistik MitgliederstatistikOutput, sep string) string {
	var lines []string
	lines = append(lines, strings.Join([]string{"Alter", "Geschlecht", "Status", "Anzahl"}, sep))
	for alter, geschlechtermap := range statistik.Statistik {
		for geschlecht, statusmap := range geschlechtermap {
			for status, anzahl := range statusmap {
				lines = append(lines, strings.Join([]string{fmt.Sprint(alter), geschlecht, status, fmt.Sprint(anzahl)}, sep))
			}
		}
	}
	sublines := lines[1:]
	sort.Strings(sublines)
	return strings.Join(lines, "\n")
}

type MitgliederstatistikTexOutput struct {
	HistogramData []HistogramLine
	Raw           MitgliederstatistikOutput
}

type HistogramLine struct {
	BeginIncluded int
	EndIncluded   int
	Innermap      map[string]int
}

/*

type MitgliederstatistikOutput struct {
	Stichtag Datum
	Statistik map[int]map[string]map[string]int //alter ->geschlecht->status->anzahl
	CollapsedStatus map[int]map[string]int //alter->geschlecht->anzahl
	CollapsedAge map[string]map[string]int //geschlecht->status->anzahl
}

 */

func mitgliederstatistik_to_tex_output(data MitgliederstatistikOutput) MitgliederstatistikTexOutput {
	var lines []HistogramLine
	var min int = 1000000
	var max int = 0
	for k, _ := range data.CollapsedStatus {
		if (k > max) {
			max = k
		}
		if (k < min) {
			min = k
		}
	}

	//group by 5 years (if possible), starting with i % 5 == 1 and ending with i % 5 == 0
	var hgl HistogramLine
	for i := min; i <= max; i++ {
		//if begin
		if (i == min || i%5 == 1) {
			hgl = HistogramLine{}
			hgl.BeginIncluded = i
			hgl.Innermap = make(map[string]int)
		}

		//for all
		var m = hgl.Innermap["Männlich"] + data.CollapsedStatus[i]["Männlich"]
		var w = hgl.Innermap["Weiblich"] + data.CollapsedStatus[i]["Weiblich"]
		var g = m + w
		hgl.Innermap["Männlich"] = m
		hgl.Innermap["Weiblich"] = w
		hgl.Innermap["Gesamt"] = g

		//if last
		if (i == max || i%5 == 0) {
			hgl.EndIncluded = i
			lines = append(lines, hgl)
			//println("Added " + fmt.Sprint(hgl))
		}
	}

	//println("Correcting CollapsedAge")

	for _, statusmap := range data.CollapsedAge {
		var n int = 0
		for _, anzahl := range statusmap {
			n += anzahl
		}
		statusmap["Gesamt"] = n
	}

	var helpermap map[string]int = make(map[string]int)

	for _, statusmap := range data.CollapsedAge {
		for status, anzahl := range statusmap {
			helpermap[status] = helpermap[status] + anzahl
		}
	}
	data.CollapsedAge["Gesamt"] = helpermap

	return MitgliederstatistikTexOutput{
		lines,
		data,
	}
}

func mitgliederstatistik_tex_formatter(data2 MitgliederstatistikOutput) string {
	var data = mitgliederstatistik_to_tex_output(data2)
	//fmt.Print(data.Raw.CollapsedAge)
	var templ = `%Statistik aus dem Mitgliederportal
\newpage
\phantomsection
\addcontentsline{toc}{section}{Mitgliederstatistik}
\markright{Mitgliederstatistik}
{\Huge Mitgliederstatistik (Stichtag {{.Raw.Stichtag.Tag}}.{{.Raw.Stichtag.Monat}}.{{.Raw.Stichtag.Jahr}})}\linebreak

\begin{figure}[h]
\centering
\begin{tikzpicture}
\pgfplotsset{width=16cm, height=10cm}
\begin{axis}[
ybar stacked,
bar width=15pt,
nodes near coords,
enlargelimits=0.15,
legend style={at={(0.5,0.90)},
anchor=north,legend columns=-1},
ylabel={},
symbolic x coords={ {{range $k,$v := .HistogramData}}{{if eq $k 0}}{{else}},{{end}}[{{$v.BeginIncluded}} - {{$v.EndIncluded}}]{{end}} },
xtick=data,
x tick label style={rotate=45,anchor=east},
]
\addplot+[ybar] plot coordinates {
{{range $k,$v := .HistogramData}} ([{{$v.BeginIncluded}} - {{$v.EndIncluded}}],{{$v.Innermap.Männlich}}){{end}}
};
\addplot+[ybar] plot coordinates {
{{range $k,$v := .HistogramData}} ([{{$v.BeginIncluded}} - {{$v.EndIncluded}}],{{$v.Innermap.Weiblich}}){{end}}
};
\legend{\strut Männer, \strut Frauen}
\end{axis}
\end{tikzpicture}
\end{figure}

\begin{tabular}{|l|l|l|l|} \hline
Bezeichnung & Männlich & Weiblich & Gesamt\\ \hline
Gesamt & {{.Raw.CollapsedAge.Männlich.Gesamt}} & {{.Raw.CollapsedAge.Weiblich.Gesamt}} & {{.Raw.CollapsedAge.Gesamt.Gesamt}}\\ \hline
Aktivitas B & {{.Raw.CollapsedAge.Männlich.AktivB}} & {{.Raw.CollapsedAge.Weiblich.AktivB}} & {{.Raw.CollapsedAge.Gesamt.AktivB}}\\ \hline
Aktivitas Ka & {{.Raw.CollapsedAge.Männlich.AktivKA}} & {{.Raw.CollapsedAge.Weiblich.AktivKA}} & {{.Raw.CollapsedAge.Gesamt.AktivKA}}\\ \hline
AHAH & {{.Raw.CollapsedAge.Männlich.AHAH}} & {{.Raw.CollapsedAge.Weiblich.AHAH}} & {{.Raw.CollapsedAge.Gesamt.AHAH}}\\ \hline
& & &\\ \hline
{{range $k,$v := .HistogramData}}  {{$v.BeginIncluded}} - {{$v.EndIncluded}} & {{$v.Innermap.Männlich}} & {{$v.Innermap.Weiblich}} & {{$v.Innermap.Gesamt}} \\ \hline
{{end}}
\end{tabular}
\pagebreak
`

	var templName = "sometemplate"
	t, _ := template.New(templName).Parse(templ)
	b := bytes.NewBufferString("")
	t.Execute(b, data)

	return b.String()
}

func mitgliederstatistik_aggregator(stichtag Datum, filepath string, format string) (string, error) {

	if (format == "json") {
		output, err := mitgliederstatistik_data_collector(stichtag, filepath)
		if (err != nil) {
			return "", err
		} else {
			jsonbytes, err := json.MarshalIndent(output, "", "    ")
			if (err != nil) {
				return "", err
			} else {
				return string(jsonbytes), nil
			}
		}
	} else if (format == "csv") {
		output, err := mitgliederstatistik_data_collector(stichtag, filepath)
		if (err != nil) {
			return "", err
		} else {
			return mitgliederstatistik_csv_formatter(output, ","), nil
		}
	} else if (format == "tex") {
		output, err := mitgliederstatistik_data_collector(stichtag, filepath)
		if (err != nil) {
			return "", err
		} else {
			return mitgliederstatistik_tex_formatter(output), nil
		}
	} else {
		return "", errors.New("Format query parameter war nicht 'json', 'tex', oder 'csv'")
	}
}

func mitgliederstatistik_data_collector(stichtag Datum, filepath string) (MitgliederstatistikOutput, error) {
	sewoberesult, sewoberr := readAndParseSewobe(filepath)
	if (sewoberr != nil) {
		return MitgliederstatistikOutput{}, sewoberr
	}
	var points []StatistikDataPoint
	for _, raw := range sewoberesult {
		if (raw.DATENSATZ["UNTERKATEGORIE"] != "ausgeschlossen" && raw.DATENSATZ["UNTERKATEGORIE"] != "ausgetreten") {
			parseresult, parseerr := statistik_parse(raw, stichtag)
			if (parseerr != nil) {
				return MitgliederstatistikOutput{}, parseerr
			} else {
				points = append(points, parseresult)
			}
		}
	}
	result, folderr := statistik_fold(points)
	if (folderr != nil) {
		return MitgliederstatistikOutput{}, folderr
	} else {
		result.Stichtag = stichtag
		result.FillTransientMaps()
		return result, nil
	}

}

type StatistikDataPoint struct {
	Alter      int
	Geschlecht string //'Männlich' oder 'Weiblich'
	Status     string //'AHAH', 'AktivB' oder 'AktivKA'
}

func statistik_parse(raw Sewoberaw, stichtag Datum) (StatistikDataPoint, error) {
	datum, datumerr := datum_parse(raw.DATENSATZ["GEBURTSDATUM"])
	if (datumerr != nil) {
		jsonraw, _ := json.Marshal(raw)
		return StatistikDataPoint{}, errors.New("raw = " +string(jsonraw) + ", datumerr: "+datumerr.Error())
	}
	alter := calculate_age(datum, stichtag)
	geschlecht := raw.DATENSATZ["GESCHLECHT"]
	if (geschlecht != "Männlich" && geschlecht != "Weiblich") {
		return StatistikDataPoint{}, errors.New("Geschlecht nicht männlich oder weiblich: " + geschlecht)
	}
	status := "N/A"
	kategorie := raw.DATENSATZ["UNTERKATEGORIE"]
	if (raw.DATENSATZ["AHSEIT"] == "0000-00-00") {
		if (strings.Contains(kategorie, "aktiv B") || strings.Contains(kategorie, "vorl. B")) {
			//Berlin?
			status = "AktivB"
		} else if (strings.Contains(kategorie, "aktiv Ka") || strings.Contains(kategorie, "vorl. Ka")) {
			//Karlsruhe?
			status = "AktivKA"
		} else {
			jsonbytes, _ := json.MarshalIndent(raw, "", "    ")
			var jsonstr = string(jsonbytes)
			return StatistikDataPoint{}, errors.New("Unterkategorie enthält kein Wort 'aktiv', aber kein AHSEIT datum gesetzt! Wert war: \n" + jsonstr)
		}
	} else {
		//AH?
		status = "AHAH"
	}
	return StatistikDataPoint{
		alter,
		geschlecht,
		status,
	}, nil
}

func statistik_fold(points []StatistikDataPoint) (MitgliederstatistikOutput, error) {
	var output MitgliederstatistikOutput
	for i := 0; i < len(points); i++ {
		point := points[i]
		if (output.Statistik == nil) {
			output.Statistik = make(map[int]map[string]map[string]int)
		}
		if (output.Statistik[point.Alter] == nil) {
			output.Statistik[point.Alter] = make(map[string]map[string]int)
		}
		old_t1_map := output.Statistik[point.Alter]
		if (old_t1_map[point.Geschlecht] == nil) {
			old_t1_map[point.Geschlecht] = make(map[string]int)
		}
		old_t2_map := old_t1_map[point.Geschlecht]
		old_t3_map := old_t2_map[point.Status]
		old_t3_map = old_t3_map + 1 //add one to counter
		//repack old_t2_map
		old_t2_map[point.Status] = old_t3_map
		//repack old_t1_map
		old_t1_map[point.Geschlecht] = old_t2_map
		//repack output
		output.Statistik[point.Alter] = old_t1_map
	}
	return output, nil
}

type MitgliederstatistikOutput struct {
	Stichtag        Datum
	Statistik       map[int]map[string]map[string]int //alter ->geschlecht->status->anzahl
	CollapsedStatus map[int]map[string]int            //alter->geschlecht->anzahl
	CollapsedAge    map[string]map[string]int         //geschlecht->status->anzahl
}

func (statistik *MitgliederstatistikOutput) FillTransientMaps() {
	statistik.CollapsedAge = make(map[string]map[string]int)
	statistik.CollapsedStatus = make(map[int]map[string]int)
	for alter, geschlechtermap := range statistik.Statistik {
		statistik.CollapsedStatus[alter] = make(map[string]int)
		for geschlecht, statusmap := range geschlechtermap {
			if (statistik.CollapsedAge[geschlecht] == nil) {
				//create map if it not exists
				statistik.CollapsedAge[geschlecht] = make(map[string]int)
			}
			var geschlechterCounter int = 0
			for status, anzahl := range statusmap {
				//für alter->geschlecht->anzahl
				geschlechterCounter += anzahl
				//für geschlecht->status->anzahl
				var old_value int = statistik.CollapsedAge[geschlecht][status]
				//caluclate second transient map
				statistik.CollapsedAge[geschlecht][status] = old_value + anzahl
			}
			//caluclate first transient map
			statistik.CollapsedStatus[alter][geschlecht] = geschlechterCounter
		}
	}
}

type Datum struct {
	Tag   int
	Monat int
	Jahr  int
}

func datum(tag, monat, jahr int) Datum {
	return Datum{
		tag,
		monat,
		jahr,
	}
}

func datum_parse(yyyymmdd string) (Datum, error) {
	t1, e1 := time.Parse(time.RFC3339, yyyymmdd+"T00:00:00+00:00")
	if (e1 != nil) {
		return Datum{}, errors.New("yyyymmdd = " +yyyymmdd + ", error: "+e1.Error())
	} else {
		return Datum{t1.Day(), int(t1.Month()), t1.Year()}, nil
	}
}

func datum_jetzt() Datum {
	return Datum{
		time.Now().Day(),
		int(time.Now().Month()),
		time.Now().Year(),
	}
}

type Sewoberaw struct {
	DATENSATZ map[string]string
}

type AemterInternalLine struct {
	AemterID          string //split string 'AMT' by comma (one internal line per instance)
	Vorname           string //VORNAME-PRIVATPERSON
	Biername          string //BIERNAME
	Nachname          string //NACHNAME-PRIVATPERSON
	Wahlveranstaltung string //NEUWAHL
	Wahlsemester      string //JAHR
}

type InputAemterKonfig struct {
	ExcludeAemter []string
	Gruppen       []InputAemterGruppe
}
type InputAemterGruppe struct {
	Titel        string
	ZeigeNeuwahl bool
	Eintraege    []InputAemterEintrag
}
type InputAemterEintrag struct {
	AemterIDs []string
	Titel     string
}

type AemterKonstellation struct {
	UnzugeordneteAemterIDs []string
	Gruppen                []AemterGruppe
}

type AemterGruppe struct {
	Titel        string
	ZeigeNeuwahl bool
	Bloecke      []AemterBlock
}

type AemterBlock struct {
	Titel  string
	Zeilen []AemterZeile
}

type AemterZeile struct {
	AemterID          string
	Vorname           string
	Biername          string
	Nachname          string
	Wahlveranstaltung string
	Wahlsemester      string
	Vakant            bool
}

func abort(str string) {
	fmt.Println("ERROR")
	fmt.Println(str)
	panic(str)
}

type Datarow struct {
	Vorname  string
	Biername string
	Nachname string
	Tag      int
	Monat    int
	Jahr     int
}

type DatarowAlter struct {
	Vorname  string
	Biername string
	Nachname string
	Tag      int
	Monat    int
	Jahr     int
	Alter    int
}

type By func(p1, p2 *Datarow) bool

func (by By) Sort(rows []Datarow) {
	ps := &rowSorter{
		rows: rows,
		by:   by,
	}
	sort.Sort(ps)
}

type rowSorter struct {
	rows []Datarow
	by   func(p1, p2 *Datarow) bool
}

func (s *rowSorter) Len() int {
	return len(s.rows)
}

func (s *rowSorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

func (s *rowSorter) Less(i, j int) bool {
	return s.by(&s.rows[i], &s.rows[j])
}

func readAndConvert(filepath string) ([]Datarow, error) {
	m, err := readAndParseSewobe(filepath)
	if (err != nil) {
		return make([]Datarow, 0), err
	}
	r := make([]Datarow, 1)
	for _, value := range m {
		tag, _ := strconv.Atoi(value.DATENSATZ["GEBURTSDATUM2"])
		monat, _ := strconv.Atoi(value.DATENSATZ["GEBURTSDATUM3"])
		jahrstr := value.DATENSATZ["GEBURTSDATUM"][0:4]
		jahr, _ := strconv.Atoi(jahrstr)
		r = append(r, Datarow{
			Vorname:  value.DATENSATZ["VORNAME-PRIVATPERSON"],
			Biername: value.DATENSATZ["BIERNAME"],
			Nachname: value.DATENSATZ["NACHNAME-PRIVATPERSON"],
			Tag:      tag,
			Monat:    monat,
			Jahr:     jahr,
		})
	}

	//sort by name (vorname > bier > nachname)
	name := func(p1, p2 *Datarow) bool {
		return p1.Vorname+" ("+p1.Biername+") "+p1.Nachname < p2.Vorname+" ("+p2.Biername+") "+p2.Nachname
	}
	By(name).Sort(r)

	return r, nil
}

func readAndParseSewobe(filepath string) (map[string]Sewoberaw, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return map[string]Sewoberaw{}, errors.New("Error A during readAndParseSewobe!")
	}
	var m map[string]Sewoberaw
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Println(err)
		return map[string]Sewoberaw{}, errors.New("Error B during readAndParseSewobe!")
	}
	return m, nil
}