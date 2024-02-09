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

    // Common routines
    import Swal from "sweetalert2";
    // @ts-ignore
    import M from "@materializecss/materialize/dist/js/materialize.min.js";
    import { ERR } from "./I18n.svelte";

    // Everything that's based on Materialize.css

    export const MAT_ENABLE_DROPDOWN = function () {
        setTimeout(() => {
            // To avoid a spurious dropdown showing in Swal2 prompts, see https://github.com/Dogfalo/materialize/issues/6202
            // @ts-ignore
            M.Dropdown.init(
                document.querySelectorAll(
                    ".dropdown-trigger:not(.swal2-select)",
                ),
                { constrainWidth: false },
            );
        }, 0);
    };

    export const MAT_ENABLE_SELECT = function () {
        // effect shown only in the clearance page, but it doesn't initialize the combo box without wrapping it in this. !!
        setTimeout(() => {
            // To avoid a spurious dropdown showing in Swal2 prompts, see https://github.com/Dogfalo/materialize/issues/6202
            // @ts-ignore
            M.FormSelect.init(
                document.querySelectorAll("select:not(.swal2-select)"),
                {},
            );
        }, 0);
    };

    export const MAT_ENABLE_FAB = function () {
        setTimeout(() => {
            // @ts-ignore
            M.FloatingActionButton.init(
                document.querySelectorAll(".fixed-action-btn"),
                { hoverEnabled: false },
            );
        }, 0);
    };

    export const TOAST = function (html, displayLength = 1500) {
        // @ts-ignore
        M.toast({ unsafeHTML: html, displayLength: displayLength });
    };

    export const SCREEN_SIZES = {
        S: 1,
        M: 2,
        L: 3,
        XL: 4,
    };

    export const GET_SCREEN_SIZE = function () {
        let w = window.innerWidth;
        if (w <= 600) return SCREEN_SIZES.S;
        else if (w <= 992) return SCREEN_SIZES.M;
        else if (w <= 1200) return SCREEN_SIZES.L;
        return SCREEN_SIZES.XL;
    };

    // REST stuff

    function mapToUrl(map) {
        let first = true;
        let urlPiece = "";
        for (const [key, value] of Object.entries(map)) {
            if (first) {
                urlPiece += "?";
                first = false;
            } else {
                urlPiece += "&";
            }
            urlPiece += key + "=" + encodeURI(value);
        }
        return urlPiece;
    }

    // @ts-ignore
    const url_prefix = "/api";

    export const CALL = async function (
        srv,
        method = "GET",
        json = null,
        map = null,
    ) {
        let url = url_prefix + "/" + srv;
        if (!!map) url += mapToUrl(map);

        const req = { method: method };
        if (method === "PUT" || method === "POST") {
            req["body"] = !!json ? JSON.stringify(json) : "{}";
            req["headers"] = { "Content-Type": "application/json" };
        }

        try {
            const res = await fetch(url, req);

            const ret = {
                isErr: !res.ok,
                status: res.status,
            };

            if (res.headers.get("Content-Type") == "application/json") {
                if (res.ok) ret.payload = await res.json();
                else {
                    const err = await res.json();
                    let msg = ERR.it[err.code];
                    msg = msg.charAt(0).toUpperCase() + msg.slice(1);
                    if (msg.includes("%s")) msg = msg.replace("%s", err.object);
                    if (!!err.error)
                        console.error("!!ERROR!!" + msg + ": " + err.error);
                    ret.message = msg;
                }
            } else ret.message = await res.text();

            return ret;
        } catch (e) {
            return {
                isErr: true,
                status: 599,
                message: e,
            };
        }
    };

    function cut(s, a, b) {
        return s.substring(a, a + b);
    }

    function cutInt(s, a, b) {
        return parseInt(cut(s, a, b));
    }

    function getNowISO() {
        const dt = new Date();
        const tzoffset = dt.getTimezoneOffset() * 60000;
        return new Date(dt.getTime() - tzoffset)
            .toISOString()
            .replaceAll("-", "")
            .replaceAll(":", "")
            .substring(0, 15);
    }

    export const FORMAT_DATE = function (m = getNowISO()) {
        return `${cutInt(m, 6, 2)}/${cutInt(m, 4, 2)}/${cutInt(m, 0, 4)}`;
    };

    export const FORMAT_TIME = function (m = getNowISO()) {
        return `${cut(m, 9, 2)}:${cut(m, 11, 2)}:${cut(m, 13, 2)}`;
    };

    export const FORMAT_DATE_TIME = function (m = getNowISO()) {
        return FORMAT_DATE(m) + " " + FORMAT_TIME(m);
    };

    export const ENC_HTML = function (text) {
        const textArea = document.createElement("textarea");
        textArea.innerText = text;
        return textArea.innerHTML;
    };

    // Everything that's based on Swal2

    function seemsHTML(text) {
        return text.includes("<") && text.includes(">");
    }

    async function alert(text, icon) {
        if (seemsHTML(text))
            await Swal.fire({
                icon: icon,
                html: text,
            });
        else
            await Swal.fire({
                icon: icon,
                text: text,
            });
    }

    export const ALERT_SUCCESS = async function (text) {
        await alert(text, "success");
    };

    export const ALERT_WARNING = async function (text) {
        await alert(text, "warning");
    };

    export const ALERT_ERROR = async function (text) {
        await alert(text, "error");
    };

    export const ALERT_INFO = async function (text) {
        await alert(text, "info");
    };

    export const DIALOG_CONFIRM = async function (text) {
        const cfg = {
            icon: "question",
            showCancelButton: true,
            confirmButtonText: "Sì",
            cancelButtonText: "No",
            cancelButtonColor: "#FF0000",
        };

        if (seemsHTML(text)) cfg["html"] = text;
        else cfg["text"] = text;

        // @ts-ignore
        return (await Swal.fire(cfg)).isConfirmed;
    };

    export const DIALOG_CHOOSE = async function (text, options) {
        const cfg = {
            icon: "question",
            confirmButtonText: "Ok",
            input: "radio",
            inputOptions: options,
        };

        if (seemsHTML(text)) cfg["html"] = text;
        else cfg["text"] = text;

        // @ts-ignore
        return (await Swal.fire(cfg)).value;
    };

    export const DIALOG_PROMPT = async function (text) {
        const cfg = {
            input: "text",
            confirmationButtonText: "Ok",
        };

        if (seemsHTML(text)) cfg["html"] = text;
        else cfg["text"] = text;

        // @ts-ignore
        return (await Swal.fire(cfg)).value;
    };

    // App-wide stuff

    export const PAGES = {
        ORDER: 1,
        ORDERS_LIST: 2,
        STATS: 3,
        STOCK: 4,
        CONSOLE: 5,
    };

    export const LOGOUT = async function () {
        // @ts-ignore
        const _checkout = window.sessionStorage.getItem("checkout");

        if (!!_checkout) {
            const checkoutId = JSON.parse(_checkout).id;
            await CALL("delSession", "DELETE", null, { id: checkoutId });
        }

        window.sessionStorage.clear(); // leaves the master key, that is in localStorage
        window.location.reload();
    };

    // Utils

    export const IS_EL_IN_VIEWPORT = function (el) {
        const rect = el.getBoundingClientRect();
        return (
            rect.top >= 0 &&
            rect.left >= 0 &&
            rect.bottom <=
                (window.innerHeight || document.documentElement.clientHeight) &&
            rect.right <=
                (window.innerWidth || document.documentElement.clientWidth)
        );
    };

    export const IS_NUMERIC = function (obj) {
        // from https://github.com/jquery/jquery/blob/bf48c21d225c31f0f9b5441d95f73615ca3dcfdb/src/core.js#L206
        return !Array.isArray(obj) && obj - parseFloat(obj) + 1 >= 0;
    };

    export const IS_INTEGER = function (obj) {
        return Number.isInteger(parseFloat(obj));
    };
</script>
