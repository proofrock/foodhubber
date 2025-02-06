<script context="module" lang="javascript">
    /*
     * Copyright (C) 2024- Germano Rizzo
     *
     * This file is part of FoodHubber.
     *
     * FoodHubber is free software: you can redistribute it and/or modify
     * it under the terms of the GNU General Public License as published by
     * the Free Software Foundation, either version 3 of the License, or
     * (at your option) any later version.
     *
     * FoodHubber is distributed in the hope that it will be useful,
     * but WITHOUT ANY WARRANTY; without even the implied warranty of
     * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     * GNU General Public License for more details.
     *
     * You should have received a copy of the GNU General Public License
     * along with FoodHubber.  If not, see <http://www.gnu.org/licenses/>.
     */

    export const ERR = {
        it: {
            FHE001: "lettura da %s non riuscita",
            FHE002: "scrittura su %s non riuscita",
            FHE003: "beneficiario non trovato",
            FHE004: "errore residuo sul resultset",
            FHE005: "il parametro '%s' non è numerico",
            FHE006: "il %s non è ben formato",
            FHE007: "in apertura di transazione",
            FHE008: "in chiusura di transazione",
            FHE009: "profilo non trovato",
            FHE010: "errore nel comporre l'excel",
            FHE011: "errore nello scaricare l'excel",
            FHE012: "errore nel comporre il report",
            FHE013: "errore nel preparare uno statement",
            FHE014: "errore nel caricare l'excel",
            FHE015: "errore nel decodificare l'excel",
            FHE100: "il ritiro non ha righe",
            FHE101: "la nota è troppo lunga",
            FHE102: "ritiro non trovato",
            FHE103: "almeno un articolo supera la disponibilità da regole di profilo",
            FHE104: "questa settimana è stato già compiuto un ritiro",
            FHE105: "il numero di ritiri di questo mese ha superato il limite",
            FHE106: "settimana non valida",
            FHE200: "impossibile trovare le intestazioni corrette nell'excel",
        },
        en: {
            FHE001: "cannot read from %s",
            FHE002: "cannot write to %s",
            FHE003: "beneficiary could not be found",
            FHE004: "residual error on resultset",
            FHE005: "parameter '%s' is not a valid number",
            FHE006: "%s is malformed",
            FHE007: "opening transaction",
            FHE008: "committing transaction",
            FHE009: "profile could not be found",
            FHE010: "error in composing excel",
            FHE011: "error in downloading excel",
            FHE012: "error in composing report",
            FHE013: "error in prepararing a statement",
            FHE014: "error in uploading excel",
            FHE015: "error in processing excel",
            FHE100: "pickup has no rows",
            FHE101: "note is too long",
            FHE102: "pickup not found",
            FHE103: "at least one item is not available due to profile rules",
            FHE104: "an order has been already done this week",
            FHE105: "this month, the number of orders exceeded the limit",
            FHE106: "invalid week",
            FHE200: "cannot find correct headers in the excel data",
        },
    };
</script>
