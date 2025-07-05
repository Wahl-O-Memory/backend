# Wahl-O-Memory Backend

Das backend für das [Wahl-O-Memory](https://github.com/Wahl-O-Memory/WahlOMemory_v2) Spiel 

## Technische Grundlagen
Das backend ist in go geschrieben und stellt elections aus dem lokalen Ordner ./elections bereit und provided die svgs die bei den Parteien als logos angegeben sind aus dem ordener ./svgs. Zusätzlich wird für eine einfache Wiederverwendbarkeit eine übersich über bereits vorhandene svgs geboten.

Es existiert ein Dockerfile, um einen Docker-Container zu generieren. Es wird empfohlen die Ordner die die veränderlichen Dateien enthalten als shared Volumes zu exposen.
Befehl zum erstellen des images:
```bash
docker build -t wom_backend .
```
Beispielbefehl zum erstellen eines Conainers aus dem image:
```bash
docker run -d -v $(pwd)/elections:/root/elections -v $(pwd)/svgs:/root/svgs -p 25000:20202 --name wom-backend-test wom_backend
```

**Mögliche configuration-Parameter**
Alle konfigurierbaren Parameter finden sich in der `main.go` unter vden vars
1. Update Intervall für den Elections-folder
2. Port auf dem alle Pfade liegen
3. Relative Ordnerpfade

**Wichtig bei SVG-Problemen:**  
Wenn svgs in der Vorschau nicht richtig angezeigt werden, überprüfen Sie ob die `viewBox` im SVG-Code gesetzt ist:

```svg
<svg height="150" width="300" viewBox="0 0 300 150" ...>
  <!-- SVG-Inhalt -->
</svg>
```

Falls keine `viewBox` vorhanden ist, fügen Sie diese entsprechend der aufgeführten `heigth` und `width` wie im Beispiel oben hinzu.
