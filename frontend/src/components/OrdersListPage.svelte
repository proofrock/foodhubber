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

    import { onDestroy, onMount } from "svelte";
    import {
        ALERT_ERROR,
        ALERT_SUCCESS,
        CALL,
        DIALOG_CONFIRM,
        DIALOG_PROMPT,
        ENC_HTML,
        FORMAT_DATE_TIME,
        FORMAT_TIME,
        MAT_ENABLE_DROPDOWN,
    } from "../components/utils/Utils.svelte";
    import {
        SUB_ORDERS_CHANGES,
        UNSUB_ORDERS_CHANGES,
    } from "./app/HubChecker.svelte";
    import SpaceForFabs from "./app/SpaceForFabs.svelte";

    export let initData;
    export let canDelete;

    $: page = 1;
    $: numPages = 0;
    $: data = null;
    $: filter = null;

    async function fetchData() {
        const map = {};
        if (!!filter) map["filter"] = filter.trim();
        const res = await CALL("getOrders", "GET", null, map);
        if (res.isErr) {
            await ALERT_ERROR(`<p>${res.message}.</p>`);
            return;
        }

        numPages = res.payload.numPages;

        const _data = [];
        const indexes = new Map();

        for (let i = 0; i < res.payload.orders.length; i++) {
            const o = res.payload.orders[i];
            const _o = {
                id: o.id,
                tstamp: o.datetime,
                checkout: o.checkout,
                cashier: o.operator,
                beneficiary: o.beneficiary,
                note: o.note,
                rows: [],
            };
            for (let j = 0; j < o.rows.length; j++) {
                _o.rows.push({
                    item: o.rows[j].item,
                    quantity: o.rows[j].qty,
                });
            }
            const idx = _data.push(_o) - 1;
            indexes.set(o.ID, idx);
        }

        data = _data;
    }

    $: {
        // @ts-ignore
        fetchData(page, filter);
    }

    async function reload() {
        page = 1;
        await fetchData();
    }

    onMount(async () => {
        SUB_ORDERS_CHANGES(reload);
        MAT_ENABLE_DROPDOWN();
        await reload();
    });

    onDestroy(() => {
        UNSUB_ORDERS_CHANGES();
    });

    async function filtering() {
        if (filter !== null) filter = null;
        else {
            const _filter = await DIALOG_PROMPT(
                "<p>Inserisci un filtro.</p><p>L'elenco sarà filtrato su id, beneficiario, postazione e operatore.</p>",
            );
            if (!!_filter) filter = _filter;
        }
    }

    let id_open = -1;

    function getInfo() {
        const oi = data.find((ord) => id_open === ord.id);
        let info = "";

        info +=
            "<p><b>Data/ora:</b> " +
            ENC_HTML(FORMAT_DATE_TIME(oi.tstamp)) +
            "</p>";
        info += "<p><b>Operatore:</b> " + ENC_HTML(oi.cashier) + "</p>";
        if (!!oi.note) info += "<p><b>Note:</b> " + ENC_HTML(oi.note) + "</p>";

        const mapp = new Map();

        for (let i = 0; i < oi.rows.length; i++)
            mapp.set(oi.rows[i].item, oi.rows[i].quantity);

        info += "<hr/><table>";

        initData.items.forEach((item) => {
            if (mapp.has(item.id))
                info +=
                    "<tr><td><b>" +
                    ENC_HTML(
                        item.item + (!!item.subitem ? "/" + item.subitem : ""),
                    ) +
                    "</b></td><td>" +
                    ENC_HTML(mapp.get(item.id)) +
                    "</td></tr>";
        });

        info += "</table>";

        return info;
    }

    async function del() {
        if (!canDelete) return;

        if (
            id_open < 0 ||
            !(await DIALOG_CONFIRM("Vuoi davvero cancellare quest'ordine?"))
        )
            return;

        const res = await CALL("delOrder", "DELETE", null, { id: id_open });
        if (res.isErr) await ALERT_ERROR(`<p>${res.message}.</p>`);
        else await ALERT_SUCCESS(`<p>Ordine ${id_open} cancellato</p>`);
        reload();
    }

    function scrollToTop() {
        document
            .getElementById("top")
            .scrollIntoView({ block: "end", behavior: "smooth" });
    }
</script>

<div>&nbsp;</div>
{#if data === null}
    <div class="center"><h5>Caricamento...</h5></div>
{:else if data.length === 0}
    <div class="center"><h5>Nessun ordine da visualizzare</h5></div>
{:else}
    <div id="top" />
    {#if page > 1}
        <div class="center">
            <a
                class="btn blue"
                href="#!"
                on:click={() => {
                    page--;
                }}
            >
                Pagina precedente ({page - 1} su {numPages})
            </a>
        </div>
    {/if}
    <div class="row">
        <div class="col hide-on-med-and-down l1 xl2" />
        <div class="input-field col s12 m12 l10 xl8">
            <table>
                <tr class="row">
                    <th>#</th>
                    <th>Postazione</th>
                    <th class="hide-on-med-and-down">Operatore</th>
                    <th>Beneficiario</th>
                    <th class="hide-on-small-only">Data/ora</th>
                    <th class="hide-on-med-and-up">Ora</th>
                    <th>&nbsp;</th>
                </tr>
                {#each data as order (order.id)}
                    <tr class="row">
                        <td>{order.id}</td>
                        <td>{order.checkout}</td>
                        <td class="hide-on-med-and-down">{order.cashier}</td>
                        <td>
                            <!-- svelte-ignore a11y-click-events-have-key-events -->
                            <!-- svelte-ignore a11y-missing-attribute -->
                            <!-- svelte-ignore a11y-no-static-element-interactions -->
                            <a
                                on:click|stopPropagation={() => {
                                    filter = order.beneficiary;
                                }}>{order.beneficiary}</a
                            ></td
                        >
                        <td class="hide-on-small-only"
                            >{FORMAT_DATE_TIME(order.tstamp)}</td
                        >
                        <td class="hide-on-med-and-up"
                            >{FORMAT_TIME(order.tstamp)}</td
                        >
                        <td
                            on:click={() => {
                                id_open = id_open === order.id ? -1 : order.id;
                            }}
                            ><i class="material-icons"
                                >expand_{id_open === order.id
                                    ? "less"
                                    : "more"}</i
                            ></td
                        >
                    </tr>
                    {#if id_open === order.id}
                        <tr>
                            <td>&nbsp;</td>
                            <td colspan="2">
                                {@html getInfo()}
                            </td>
                            <td>&nbsp;</td>
                            <td>
                                {#if canDelete}
                                    <a
                                        class="btn red"
                                        href="#!"
                                        on:click={del}
                                        title="Cancella ordine"
                                    >
                                        <i class="material-icons"
                                            >delete_forever</i
                                        >
                                    </a>
                                {:else}
                                    &nbsp;
                                {/if}
                            </td>
                        </tr>
                    {/if}
                {/each}
            </table>
        </div>
        <div class="col hide-on-med-and-down l1 xl2" />
    </div>
    <div>&nbsp;</div>
    {#if page < numPages}
        <div class="center">
            <a
                class="btn blue"
                href="#!"
                on:click={() => {
                    page++;
                    scrollToTop();
                }}
            >
                Pagina successiva ({page + 1} su {numPages})
            </a>
        </div>
    {:else}
        <div class="center"><h6><i>Nessun altro ordine</i></h6></div>
    {/if}
{/if}

<SpaceForFabs />

<div class="fixed-action-btn">
    <a
        class="btn-floating btn-large darken-2"
        class:yellow={filter == null}
        class:red={filter != null}
        href="#!"
        on:click={filtering}
    >
        <i class="large material-icons">
            {#if filter == null}search{:else}search_off{/if}
        </i>
    </a>
</div>
