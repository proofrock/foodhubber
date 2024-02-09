# FoodHubber

<img src="https://eventi.emergency.it/wp-content/uploads/cropped-favicon.png" alt="Emergency" width="200"/>
<img src="https://media.licdn.com/dms/image/C4E0BAQFhNBC2FoSLqw/company-logo_200_200/0/1643361232408/aton_spa_logo?e=2147483647&v=beta&t=Z64YPuG9Az_o9LnDX68tmqzAJ_KHMREjg04uk7UjpFY" alt="Aton" width="200"/>

This project is (or, will be) a web system to manage a Food Hub that distributes food and other goods _pro bono_. These food hubs are put in place in the italian territory by Emergency ONG Onlus.

This project is sponsored by [Emergency ONG Onlus](https://emergency.it) and [Aton Società Benefit](https://www.aton.com).


This project is released as Free and Open Source Software under the [GPLv3 license](https://www.gnu.org/licenses/quick-guide-gplv3.it.html).

**NOTE**: The development is currently happening in the `develop` branch. The first versions will be italian only.

## Installazione (ita)

- [Scaricare](https://github.com/proofrock/foodhubber/releases) il pacchetto per la propria architettura (es. `foodhubber-v0.1.0-win-amd64.zip`);
- scompattarlo in una directory;
- avviare l'eseguibile;
- di default cerca il database nella directory di esecuzione;
  - es. sotto windows dovrebbe essere sufficiente fare doppio click sull'eseguibile.
- collegarsi con un browser a http://localhost:31020

E' attivo un beneficiario con codice `123`.

Per test, potrebbe essere utile cambiare la settimana corrente; può essere fatto con il flag da linea di comando `--force-week`:

```bash
.\foodhubber.exe --force-week=1
```

## Resources

- [Palette](https://kdesign.co/blog/pastel-color-palette-examples/)
