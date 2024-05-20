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

    // TODO display item color in the cell?
    import { onDestroy, onMount } from "svelte";
    import {
        ALERT_ERROR,
        CALL,
        DIALOG_CONFIRM,
        DIALOG_PROMPT,
        MAT_ENABLE_SELECT,
        TOAST,
    } from "../components/utils/Utils.svelte";
    import {
        SUB_STOCK_CHANGES,
        UNSUB_STOCK_CHANGES,
    } from "./app/HubChecker.svelte";

    export let initData;
    export let canChangeStock;

    $: data = null;
    $: dataFiltered = null;
    $: filter = "2";
    $: selectedRow = null;

    $: {
        // @ts-ignore
        applyFilter(filter, data);
    }

    function applyFilter() {
        if (data !== null) {
            if (filter === "2") {
                dataFiltered = data;
            } else if (filter === "1") {
                dataFiltered = data.filter((itm) => itm.stock !== "∞");
            } else if (filter === "0") {
                dataFiltered = data.filter(
                    (itm) =>
                        itm.stock !== "∞" && itm.stock <= initData.yellow_limit,
                );
            }
        }
    }

    async function fetchData() {
        if (initData === null) return;

        let __data = await CALL("getStock");
        if (__data.isErr) {
            await ALERT_ERROR(`<p>Init fallito.</p><p>${__data.message}.</p>`);
            return;
        }
        let _data = __data.payload;

        const items_to_display = [];
        initData.items.forEach((item) => {
            const itm = _data.stock.find((i) => item.id === i.item);
            items_to_display.push({
                id: item.id,
                name: item.item + (!!item.subitem ? "/" + item.subitem : ""),
                stock: !!itm ? itm.stock : "∞",
                class: !itm
                    ? ""
                    : itm.stock <= initData.red_limit
                      ? "bd-red"
                      : itm.stock <= initData.yellow_limit
                        ? "bd-yellow"
                        : "",
            });
        });
        data = items_to_display;
    }

    onMount(async () => {
        SUB_STOCK_CHANGES(fetchData);
        await fetchData();
        MAT_ENABLE_SELECT();
    });

    onDestroy(() => {
        UNSUB_STOCK_CHANGES();
    });

    async function setStock(item, stock) {
        if (!canChangeStock) return;

        if (stock !== null) {
            let text = "";
            if (stock.startsWith("+"))
                text = `<p>La giacenza sarà <u>aumentata</u> di ${stock.substring(1)} unità.</p>`;
            else if (stock.startsWith("-"))
                text = `<p>La giacenza sarà <u>diminuita</u> di ${stock.substring(1)} unità.</p>`;
            else
                text = `<p>La giacenza sarà <u>impostata</u> a ${stock} unità.</p>`;

            text += "<p>Vuoi confermare?</p>";
            if (!(await DIALOG_CONFIRM(text))) return;
        }

        const ret = await CALL("setStock", "POST", { item, stock });
        if (ret.isErr)
            await ALERT_ERROR(`<p>Modifica fallita.</p><p>${ret.message}.</p>`);
        else {
            fetchData();
            selectedRow = 0;

            TOAST("Scorta impostata", 1000);
        }
    }

    async function askStock(item) {
        const stock = await DIALOG_PROMPT(`
        <p>Inserisci una quantità, o una differenza (+10, -20...)<p>
        <p>Vuoto: nessuna gestione (scorta infinita)<p>`);
        if (stock === null || stock.trim() === "") return;
        await setStock(item, stock);
    }
</script>

<div>&nbsp;</div>
<div class="row">
    <div class="col hide-on-small-and-down m2 l3 xl4" />
    <div class="input-field col s12 m8 l6 xl4">
        <i class="material-icons prefix">filter_alt</i>
        <select id="filter" bind:value={filter}>
            <option value="2">Mostra tutto</option>
            <option value="1" selected>Solo scorte gestite</option>
            <option value="0">Solo sottoscorta</option>
        </select>
        <label for="filter">Filtra</label>
    </div>
    <div class="col hide-on-small-and-down m2 l3 xl4" />
</div>
{#if dataFiltered !== null}
    <div class="row">
        <div class="col hide-on-small-and-down m1 l2 xl3" />
        <div class="input-field col s12 m10 l8 xl6">
            <table>
                <tr>
                    <th />
                    <th>Scorta Rimanente</th>
                    <th>&nbsp;</th>
                </tr>
                {#each dataFiltered as ing (ing.id)}
                    <tr class={ing.class}>
                        <td>{ing.name}</td>
                        <td>{ing.stock}</td>
                        {#if canChangeStock}
                            <td
                                on:click={() => {
                                    selectedRow =
                                        ing.id === selectedRow ? null : ing.id;
                                }}
                                ><i class="material-icons"
                                    >expand_{ing.id === selectedRow
                                        ? "less"
                                        : "more"}</i
                                ></td
                            >
                        {:else}
                            <td>&nbsp;</td>
                        {/if}
                    </tr>
                    {#if canChangeStock && selectedRow === ing.id}
                        <tr>
                            <td colspan="2">
                                &nbsp;&nbsp;&nbsp;
                                <a
                                    href="#!"
                                    class="btn red"
                                    on:click={() => {
                                        setStock(ing.id, null);
                                    }}
                                    title="Non gestire scorta"
                                >
                                    <i class="material-icons">all_inclusive</i>
                                </a>
                                <a
                                    href="#!"
                                    class="btn blue"
                                    on:click={() => {
                                        setStock(ing.id, "0");
                                    }}
                                    title="Imposta a zero"
                                >
                                    <i class="material-icons">exposure_zero</i>
                                </a>
                                <a
                                    href="#!"
                                    class="btn green"
                                    on:click={() => {
                                        askStock(ing.id);
                                    }}
                                    title="Imposta..."
                                >
                                    <i class="material-icons">edit</i>
                                </a>
                            </td>
                        </tr>
                    {/if}
                {/each}
            </table>
        </div>
        <div class="col hide-on-small-and-down m1 l2 xl3" />
    </div>
{/if}

<div class="fixed-action-btn">
    <a
        class="btn-floating btn-large darken-2"
        href="/api/getStockExcel"
        target="_blank"
        title="Scarica come file Excel"
    >
        <i class="large material-icons"> file_download </i>
    </a>
</div>

<style>
    .bd-yellow {
        background-color: gold;
    }

    .bd-red {
        background-color: red;
        color: white;
    }
</style>
