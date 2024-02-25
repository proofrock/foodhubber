<script>
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

    import {
        ALERT_INFO,
        ALERT_ERROR,
        CALL,
        FORMAT_DATE_TIME,
    } from "../components/utils/Utils.svelte";
    import { ERR } from "../components/utils/I18n.svelte";
    import { onDestroy, onMount } from "svelte";

    export let initData = null;

    $: data = null;

    async function getData() {
        const res = await CALL("getSessions");
        if (res.isErr) await ALERT_ERROR(`<p>${res.message}.</p>`);
        else data = res.payload;
    }

    let interval = null;
    onMount(() => {
        getData();
        interval = setInterval(getData, initData.polling_cycle * 2);
    });

    onDestroy(() => {
        if (interval !== null) clearInterval(interval);
    });

    async function handleSubmitBenefExcel(e) {
        const form = e.currentTarget;
        const url = new URL(form.action);
        // @ts-ignore
        let file = document.getElementById("uplBenefXls").files[0];
        if (!file) {
            await ALERT_ERROR(`<p>Scegliere un file.</p>`);
            return;
        }
        let formData = new FormData();

        formData.append("file", file);

        const fetchOptions = {
            method: form.method,
            body: formData,
        };

        const res = await fetch(url, fetchOptions);
        if (res.status === 200) {
            await ALERT_INFO("<p>Ok, dati caricati.</p>");
        } else {
            const err = await res.json();
            let msg = ERR.it[err.code];
            msg = msg.charAt(0).toUpperCase() + msg.slice(1);
            if (msg.includes("%s")) msg = msg.replace("%s", err.object);
            if (!!err.error)
                console.error("!!ERROR!!" + msg + ": " + err.error);
            await ALERT_ERROR(`<p>${msg}.</p>`);
        }
    }
</script>

<div>&nbsp;</div>
<div class="row">
    <div class="col hide-on-med-and-down s1 m2 l3 xl4" />
    <div class="center col s10 m8 l6 xl4">
        <div>&nbsp;</div>

        <h5 class="center">Caricamento Excel dei beneficiari</h5>
        <div>&nbsp;</div>
        <form
            method="post"
            enctype="multipart/form-data"
            action="/api/setBeneficiariesExcel"
            on:submit|preventDefault={handleSubmitBenefExcel}
        >
            <div class="file-field input-field">
                <div class="btn">
                    <span>File</span>
                    <input type="file" id="uplBenefXls" />
                </div>
                <div class="file-path-wrapper">
                    <input class="file-path" type="text" />
                </div>
            </div>
            <button class="waves-effect waves-light btn green" type="submit"
                >Carica</button
            >
        </form>
    </div>
    <div class="col hide-on-med-and-down s1 m2 l3 xl4" />
</div>
{#if data !== null}
    <div>&nbsp;</div>
    <div class="divider">&nbsp;</div>
    <div class="row">
        <div class="col hide-on-med-and-down m1 l2 xl3" />
        <div class="input-field col s12 m10 l8 xl6">
            <h5 class="center">Client Connessi</h5>
            <div>&nbsp;</div>
            <table>
                {#each data.checkouts as checkout (checkout.id)}
                    <tr>
                        <td>
                            <i class="material-icons">computer</i>
                        </td>
                        <td>{checkout.id} ({checkout.operator})</td>
                        <td>
                            <div
                                class="material-icons"
                                title="Ultima connessione: {FORMAT_DATE_TIME(
                                    checkout.datetime,
                                )}"
                            >
                                {#if checkout.active}done{:else}clear{/if}
                            </div>
                        </td>
                    </tr>
                {/each}
            </table>
        </div>
        <div class="col hide-on-med-and-down m1 l2 xl3" />
    </div>
    <div class="row">&nbsp;</div>
{/if}
