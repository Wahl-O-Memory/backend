**Wichtig bei SVG-Problemen:**  
Sollten Probleme auftreten, überprüfen Sie ob die `viewBox` im SVG-Code gesetzt ist:

```svg
<svg height="150" width="300" viewBox="0 0 300 150" ...>
  <!-- SVG-Inhalt -->
</svg>
```

Falls keine `viewBox` vorhanden ist, fügen Sie diese entsprechend der aufgeführten `heigth` und `width` wie im Beispiel oben hinzu
