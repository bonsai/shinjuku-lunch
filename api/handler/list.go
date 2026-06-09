package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type listRestaurant struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Area      string  `json:"area"`
	Genre     string  `json:"genre"`
	Price     string  `json:"price"`
	WalkMin   *int    `json:"walk_min,omitempty"`
	Station   *string `json:"station,omitempty"`
	Rating    string  `json:"rating"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Notes     string  `json:"notes"`
}

func loadListData() []listRestaurant {
	raw, err := os.ReadFile("../neon/seed.json")
	if err != nil {
		log.Printf("cannot read seed.json: %v", err)
		return nil
	}
	var seed struct {
		Restaurants []struct {
			Name      string   `json:"name"`
			Area      string   `json:"area"`
			Genre     string   `json:"genre"`
			Notes     *string  `json:"notes,omitempty"`
			Station   *string  `json:"station,omitempty"`
			WalkMin   *int     `json:"walk_min,omitempty"`
			Latitude  *float64 `json:"latitude,omitempty"`
			Longitude *float64 `json:"longitude,omitempty"`
		} `json:"restaurants"`
	}
	json.Unmarshal(raw, &seed)

	var out []listRestaurant
	for i, r := range seed.Restaurants {
		price := ""
		if r.Notes != nil {
			s := *r.Notes
			idx := strings.Index(s, "¥")
			if idx >= 0 {
				end := idx + 1
				for end < len(s) && s[end] >= '0' && s[end] <= '9' {
					end++
				}
				price = s[idx:end]
			}
		}
		notes := ""
		if r.Notes != nil {
			notes = *r.Notes
		}
		out = append(out, listRestaurant{
			ID: i + 1, Name: r.Name, Area: r.Area, Genre: r.Genre,
			Price: price, WalkMin: r.WalkMin, Station: r.Station,
			Latitude: r.Latitude, Longitude: r.Longitude,
			Notes: notes,
		})
	}
	return out
}

func genStars(n int) string {
	return strings.Repeat("★", n) + strings.Repeat("☆", 5-n)
}

func colorForGenre(g string) string {
	m := map[string]string{
		"韓国料理": "#e74c3c",
		"タイ料理": "#2ecc71",
		"和食":    "#3498db",
		"インド料理": "#f39c12",
		"アジア料理": "#9b59b6",
		"カレー":  "#f1c40f",
		"牛丼":    "#8B4513",
		"トルコ料理": "#1abc9c",
		"焼肉":    "#c0392b",
	}
	c, ok := m[g]
	if !ok {
		return "#7f8c8d"
	}
	return c
}

func ListPage(w http.ResponseWriter, r *http.Request) {
	restaurants := loadListData()

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>新宿ランチナビ</title>
<style>
* { margin:0; padding:0; box-sizing:border-box; }
body { font-family:'Segoe UI',system-ui,sans-serif; background:#0f1117; color:#e0e0e0; min-height:100vh; }
.header { background:linear-gradient(135deg,#1a1d2e,#2a1e3c); padding:24px 32px; border-bottom:1px solid #2a2d3a; }
.header h1 { font-size:28px; color:#fff; }
.header p { color:#888; margin-top:4px; font-size:14px; }
.container { max-width:1200px; margin:0 auto; padding:24px 16px; }
.controls { display:flex; gap:12px; flex-wrap:wrap; margin-bottom:24px; }
.controls select,.controls input { background:#1e2030; color:#e0e0e0; border:1px solid #2a2d3a; padding:10px 16px; border-radius:8px; font-size:14px; min-width:140px; }
.controls select:focus,.controls input:focus { outline:none; border-color:#5b6ef5; }
.count { color:#888; font-size:14px; margin-bottom:16px; }
.grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(320px,1fr)); gap:16px; }
.card { background:#1a1d2e; border:1px solid #2a2d3a; border-radius:12px; padding:20px; transition:transform .15s,border-color .15s; cursor:pointer; }
.card:hover { transform:translateY(-2px); border-color:#5b6ef5; }
.card-header { display:flex; justify-content:space-between; align-items:start; margin-bottom:8px; }
.card-name { font-size:18px; font-weight:600; color:#fff; }
.card-genre { display:inline-block; padding:2px 10px; border-radius:99px; font-size:12px; font-weight:600; color:#fff; }
.card-meta { display:flex; gap:16px; font-size:13px; color:#888; margin-bottom:8px; }
.card-meta span { display:flex; align-items:center; gap:4px; }
.card-notes { font-size:13px; color:#666; line-height:1.5; }
.card-actions { margin-top:12px; display:flex; gap:8px; }
.card-actions a { display:inline-block; padding:6px 14px; border-radius:6px; font-size:12px; text-decoration:none; background:#2a2d3a; color:#bbb; transition:background .15s; }
.card-actions a:hover { background:#3a3d4a; color:#fff; }
</style>
</head>
<body>
<div class="header"><h1>🍽 新宿ランチナビ</h1><p>新宿・歌舞伎町・大久保のランチ情報</p></div>
<div class="container">
<div class="controls">
<select id="areaFilter" onchange="applyFilter()"><option value="">すべてのエリア</option>`)
	areas := []string{"歌舞伎町", "大久保", "西新宿", "新宿三丁目", "新宿"}
	for _, a := range areas {
		fmt.Fprintf(&b, `<option value="%s">%s</option>`, a, a)
	}
	b.WriteString(`</select>
<select id="genreFilter" onchange="applyFilter()"><option value="">すべてのジャンル</option>`)
	genres := []string{"韓国料理", "タイ料理", "和食", "インド料理", "アジア料理", "カレー", "牛丼", "トルコ料理", "焼肉"}
	for _, g := range genres {
		fmt.Fprintf(&b, `<option value="%s">%s</option>`, g, g)
	}
	b.WriteString(`</select>
<input type="text" id="searchInput" placeholder="店名で検索..." oninput="applyFilter()">
</div>
`)
	fmt.Fprintf(&b, `<p class="count" id="count">%d件の店舗</p>`, len(restaurants))
	b.WriteString(`<div class="grid" id="grid">
`)

	for _, r := range restaurants {
		col := colorForGenre(r.Genre)
		walkStr := ""
		if r.WalkMin != nil {
			walkStr = fmt.Sprintf("🚶 %d分", *r.WalkMin)
		}
		stationStr := ""
		if r.Station != nil && *r.Station != "" {
			stationStr = fmt.Sprintf("🚉 %s", *r.Station)
		}
		gmapLink := ""
		if r.Latitude != nil && r.Longitude != nil {
			gmapLink = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *r.Latitude, *r.Longitude)
		}
		areaClass := fmt.Sprintf("area-%s", r.Area)
		genreClass := fmt.Sprintf("genre-%s", r.Genre)

		fmt.Fprintf(&b, `<div class="card %s %s" data-area="%s" data-genre="%s" data-name="%s">
<div class="card-header">
<span class="card-name">%s</span>
<span class="card-genre" style="background:%s">%s</span>
</div>
<div class="card-meta">
<span>📍 %s</span>
<span>%s</span>
<span>%s</span>
</div>
<div class="card-notes">%s</div>
<div class="card-actions">`,
			areaClass, genreClass, r.Area, r.Genre, r.Name,
			r.Name, col, r.Genre,
			r.Area, walkStr, stationStr, r.Notes)

		if gmapLink != "" {
			fmt.Fprintf(&b, `<a href="%s" target="_blank">🗺 Google Maps</a>`, gmapLink)
		}
		if r.Price != "" {
			fmt.Fprintf(&b, `<span style="margin-left:auto;font-size:15px;font-weight:600;color:%s">%s</span>`, col, r.Price)
		}
		b.WriteString(`</div></div>`)
	}

	b.WriteString(`</div></div>
<script>
const cards=document.querySelectorAll('.card');
function applyFilter(){
const area=document.getElementById('areaFilter').value;
const genre=document.getElementById('genreFilter').value;
const q=document.getElementById('searchInput').value.toLowerCase();
let n=0;
cards.forEach(c=>{
const ma=!area||c.dataset.area===area;
const mg=!genre||c.dataset.genre===genre;
const mn=!q||c.dataset.name.includes(q);
c.style.display=(ma&&mg&&mn)?'':'none';
if(ma&&mg&&mn)n++;
});
document.getElementById('count').textContent=n+'件の店舗';
}
</script>
</body>
</html>`)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(b.String()))
}
