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

    import Item from "./Item.svelte";
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    import {
        DIALOG_CONFIRM,
        GET_SCREEN_SIZE,
        MAT_ENABLE_FAB,
    } from "../../components/utils/Utils.svelte";
    import SpaceForFabs from "../app/SpaceForFabs.svelte";

    const dispatch = createEventDispatcher();

    export let order;
    export let initData;

    $: rows = [];
    $: selItem = 0;
    $: selectedSection = null;
    $: residualAllowance = {};

    const chart = new Map(); // a chart of the grid: for every item [id], its coordinates
    const rowsMap = new Map(); // map of the sections

    function rearrangeCols(colNum) {
        function compose(_list, _dest) {
            for (let i = 0; i < _list.length; i += colNum)
                _dest.push(_list.slice(i, i + colNum));
        }

        const nuRows = [];
        chart.clear();
        rowsMap.clear();
        selectedSection = null;
        if (curColNum >= 3) {
            // No sections
            compose(initData.items, nuRows);

            for (let i = 0; i < nuRows.length; i++)
                for (let j = 0; j < nuRows[i].length; j++)
                    chart.set(nuRows[i][j].id, { r: i, c: j });
        } else {
            // Sections
            for (let i = 0; i < initData.item_types.length; i++) {
                const item_type = initData.item_types[i];
                const filtered_list = initData.items.filter(
                    (li) => li.type === item_type.id,
                );
                if (filtered_list.length > 0) {
                    if (selectedSection === null)
                        selectedSection = item_type.id;
                    const filteredRows = [];
                    compose(filtered_list, filteredRows);
                    const section = {
                        id: item_type.id,
                        title: `Gruppo ${item_type.id}`,
                        color: item_type.color,
                        rows: filteredRows,
                    };
                    nuRows.push(section);
                    rowsMap.set(section.id, section);

                    for (let i = 0; i < filteredRows.length; i++)
                        for (let j = 0; j < filteredRows[i].length; j++)
                            chart.set(filteredRows[i][j].id, { r: i, c: j });
                }
            }
        }
        rows = nuRows;
    }

    let curColNum = -1;

    function resized() {
        const cn = GET_SCREEN_SIZE(); // here the constant returned is the number of columns, conveniently
        if (curColNum !== cn) {
            curColNum = cn;
            rearrangeCols(cn);
        }
    }

    onMount(async () => {
        MAT_ENABLE_FAB();
        window.addEventListener("resize", resized);

        for (let i = 0; i < initData.items.length; i++) {
            const item = initData.items[i];
            if (
                order.allowance.hasOwnProperty(item.item) &&
                order.allowance[item.item] >= 0
            )
                residualAllowance[item.id.toString()] = {
                    limit: order.allowance[item.item],
                    item: item.item,
                };
            else
                residualAllowance[item.id.toString()] = {
                    limit: -1,
                    item: "",
                };
        }

        resized();
    });

    onDestroy(() => {
        window.removeEventListener("resize", resized);
    });

    function recalcResidualAllowance(id, increment) {
        for (const [_id, _obj] of Object.entries(residualAllowance))
            if (
                _obj.item === initData.itemsMap.get(id).item &&
                _id !== id.toString()
            )
                residualAllowance[_id].limit += increment;
        return residualAllowance;
    }

    function increase(data) {
        const id = data.detail;
        order.items[id.toString()]++;
        selItem = id;

        residualAllowance = recalcResidualAllowance(id, -1);
    }

    function decrease(data) {
        const id = data.detail;
        order.items[id.toString()]--;
        selItem = id;

        residualAllowance = recalcResidualAllowance(id, +1);
    }

    function key(evt) {
        if (evt.keyCode >= 37 && evt.keyCode <= 40) {
            // Arrows
            const { r, c } = chart.get(selItem);
            let _rows;
            if (selectedSection !== null)
                // if multi-section (1 or 2 cols) take the subsection
                _rows = rowsMap.get(selectedSection).rows;
            else _rows = rows;
            try {
                if (evt.keyCode === 37) selItem = _rows[r][c - 1].id;
                else if (evt.keyCode === 38) selItem = _rows[r - 1][c].id;
                else if (evt.keyCode === 39) selItem = _rows[r][c + 1].id;
                else if (evt.keyCode === 40) selItem = _rows[r + 1][c].id;
            } catch (e) {}
        } else if (evt.keyCode === 32) {
            // Spacebar
            order.items[selItem.toString()]++;
        } else if (evt.keyCode === 46 || evt.keyCode === 8) {
            // Del or Backspace
            order.items[selItem.toString()] = 0;
        }
    }

    function chSection(selId) {
        if (selId === selectedSection) selectedSection = null;
        else {
            selectedSection = selId;
            const _rows = rowsMap.get(selectedSection).rows;
            selItem = _rows[0][0].id;
        }
    }
</script>

<svelte:window on:keydown|preventDefault|stopPropagation={key} />

<div style="height: 10px;" />

{#if order !== null && initData !== null}
    {#if curColNum >= 3}
        <div>
            {#each rows as row}
                <div class="row sm-div-nomargin">
                    {#each row as item}
                        <div class="col s12 m6 l4 xl3 nomargin">
                            <Item
                                itm={item}
                                bind:selItem
                                bind:limit={residualAllowance[
                                    item.id.toString()
                                ].limit}
                                {order}
                                on:increase={increase}
                                on:decrease={decrease}
                            />
                        </div>
                    {/each}
                </div>
            {/each}
        </div>
    {:else}
        <div>
            {#each rows as section}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <!-- svelte-ignore a11y-no-static-element-interactions -->
                <div
                    class="card bd sm-cell-nomargin"
                    style="background-color: {section.color};"
                    on:click|stopPropagation={() => {
                        chSection(section.id);
                    }}
                    on:contextmenu|preventDefault={() => {}}
                >
                    <div class="card-content">
                        <span class="card-title">{section.title}</span>
                    </div>
                </div>
                {#if selectedSection === section.id}
                    {#each section.rows as row}
                        <div class="row sm-div-nomargin">
                            {#each row as item}
                                <div class="col s12 m6 l4 xl3 nomargin">
                                    <Item
                                        itm={item}
                                        bind:selItem
                                        bind:limit={residualAllowance[
                                            item.id.toString()
                                        ].limit}
                                        {order}
                                        on:increase={increase}
                                        on:decrease={decrease}
                                    />
                                </div>
                            {/each}
                        </div>
                    {/each}
                {/if}
            {/each}
        </div>
    {/if}

    <SpaceForFabs />

    <div class="fixed-action-btn sm-fab-left">
        <a
            class="btn-floating btn-large red"
            href="#!"
            on:click={async () => {
                const txt = "Sei sicuro di voler annullare il ritiro?";
                if (await DIALOG_CONFIRM(txt)) dispatch("reset", null);
            }}
        >
            <i class="large material-icons">delete_sweep</i>
        </a>
    </div>

    <div class="fixed-action-btn">
        <a
            class="btn-floating btn-large green"
            href="#!"
            on:click={() => {
                order.subpage = 3;
            }}
        >
            <i class="material-icons">arrow_forward</i>
        </a>
    </div>
{/if}
