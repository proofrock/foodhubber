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
        CALL,
        MAT_ENABLE_SELECT,
        TOAST,
    } from "../components/utils/Utils.svelte";
    import {
        SUB_ORDERS_CHANGES,
        UNSUB_ORDERS_CHANGES,
    } from "./app/HubChecker.svelte";
    import Divider from "./orderpage/Divider.svelte";
    import { DateInput } from "date-picker-svelte";

    export let initData;

    // XXX Ugly hack: initializing dateFrom/dateTo triggers 2 times the $: block. So, this
    // is used both to fire it one time only, and to avoid firing the toast at first init.
    let countdown = 2;

    $: data = null;
    $: dateFrom = new Date();
    $: dateTo = new Date();

    $: {
        const from = formatDateForService(dateFrom);
        const to = formatDateForService(dateTo);
        if (countdown > 1) {
            countdown--;
        } else if (to < from) {
            ALERT_ERROR("Le date sono in ordine sbagliato.");
            dateFrom = dateTo;
        } else {
            fetchData();
            if (countdown > 0) countdown--;
            else TOAST("Dati aggiornati.");
        }
    }

    function formatDateForService(date) {
        var yyyy = date.getFullYear().toString();
        var mm = (date.getMonth() + 1).toString().padStart(2, "0"); // getMonth() is zero-based
        var dd = date.getDate().toString().padStart(2, "0");
        return `${yyyy}${mm}${dd}`;
    }

    function formatDateForDisplay(date) {
        var yyyy = date.getFullYear().toString();
        var mm = (date.getMonth() + 1).toString().padStart(2, "0"); // getMonth() is zero-based
        var dd = date.getDate().toString().padStart(2, "0");
        return `${dd}/${mm}/${yyyy}`;
    }

    async function fetchData() {
        const from = formatDateForService(dateFrom);
        const to = formatDateForService(dateTo);
        if (to < from) {
            ALERT_ERROR("Le date sono in ordine sbagliato.");
            return;
        }

        const res = await CALL("getStats", "GET", null, { from, to });
        if (res.isErr) {
            ALERT_ERROR(res.message);
            return;
        }

        const _data = res.payload;

        const items_to_display = [];
        initData.items.forEach((item) => {
            const itm = _data.byItem.find((i) => item.id === i.item);
            if (itm)
                items_to_display.push({
                    id: item.id,
                    name:
                        item.item + (!!item.subitem ? "/" + item.subitem : ""),
                    quantity: itm.qty,
                });
        });
        _data.items_to_display = items_to_display;
        _data.from = formatDateForDisplay(dateFrom);
        _data.to = formatDateForDisplay(dateTo);
        data = _data;
    }

    onMount(async () => {
        MAT_ENABLE_SELECT();

        SUB_ORDERS_CHANGES(fetchData);
    });

    onDestroy(() => {
        UNSUB_ORDERS_CHANGES();
    });
</script>

<div>&nbsp;</div>
<div class="row">
    <div class="col hide-on-small-and-down m2 l3 xl4" />
    <div class="input-field inline col s6 m4 l3 xl2">
        <DateInput
            id="from"
            format="dd/MM/yyyy"
            bind:value={dateFrom}
            closeOnSelection={true}
        />
        <label for="from" class="active">Data Inizio</label>
    </div>
    <div class="input-field col s6 m4 l3 xl2">
        <DateInput
            id="to"
            format="dd/MM/yyyy"
            bind:value={dateTo}
            closeOnSelection={true}
        />
        <label for="to" class="active">Data Fine</label>
    </div>
    <div class="col hide-on-small-and-down m2 l3 xl4" />
</div>
<Divider />
{#if !!data}
    <div class="row">
        <div class="col s12">
            <div class="center"><h6>{data.from} - {data.to}</h6></div>
        </div>
    </div>
    <div class="row">
        <div class="col hide-on-small-and-down m1 l2 xl3" />
        <div class="input-field col s12 m5 l4 xl3">
            <div class="center"><h5>Generali</h5></div>
            <Divider />
            <table>
                <tr>
                    <th />
                    <th>Valore</th>
                </tr>
                <tr>
                    <th>Ritiri</th>
                    <td>{data.orders}</td>
                </tr>
                <tr>
                    <th>Articoli</th>
                    <td>{data.items}</td>
                </tr>
            </table>
        </div>
        <div class="input-field col s12 m5 l4 xl3">
            <div class="center"><h5>Per postazione</h5></div>
            <Divider />
            <table>
                <tr>
                    <th />
                    <th>Ritiri</th>
                    <th>Articoli</th>
                </tr>
                {#each data.byCheckout as ck (ck.checkout)}
                    <tr>
                        <td>{ck.checkout}</td>
                        <td>{ck.orders}</td>
                        <td>{ck.items}</td>
                    </tr>
                {/each}
            </table>
        </div>
        <div class="col hide-on-small-and-down m1 l2 xl3" />
    </div>
    <div class="row">
        <div class="col hide-on-small-and-down m3 l4 xl5" />
        <div class="input-field col s12 m6 l4 xl2">
            <div class="center"><h5>Per Articolo</h5></div>
            <Divider />
            <table>
                <tr>
                    <th />
                    <th>Ritirati</th>
                </tr>
                {#each data.items_to_display as itm (itm.id)}
                    <tr>
                        <td>{itm.name}</td>
                        <td>{itm.quantity}</td>
                    </tr>
                {/each}
            </table>
        </div>
        <div class="col hide-on-small-and-down m3 l4 xl5" />
    </div>
{/if}
