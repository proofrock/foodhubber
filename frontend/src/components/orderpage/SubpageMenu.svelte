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
        GET_SCREEN_SIZE,
        MAT_ENABLE_FAB,
    } from "../../components/utils/Utils.svelte";
    import SpaceForFabs from "../app/SpaceForFabs.svelte";

    const dispatch = createEventDispatcher();

    export let order;
    export let initData;

    $: rows = [];
    let rowsMap = new Map();
    $: selItem = 0;
    $: selectedSection = null;

    const chart = new Map(); // a chart of the grid: for every item [id], its coordinates

    function rearrangeCols(colNum) {
        /* To order by column
        function compose(_list, _dest) {
            let rowNum = Math.ceil(_list.length / colNum);
            let lenOfLastRow = _list.length - colNum * (rowNum - 1);
            let i = 0;
            for (let r = 0; r < rowNum - 1; r++) _dest.push(new Array(colNum));
            _dest.push(new Array(lenOfLastRow));
            for (let c = 0; c < colNum; c++)
                for (let r = 0; r < rowNum; r++)
                    if (r < rowNum - 1 || c < lenOfLastRow) {
                        _dest[r][c] = _list[i];
                        chart.set(_list[i].id, { r: r, c: c });
                        i++;
                    }
        }
        */

        function compose(_list, _dest) {
            for (let i = 0; i < _list.length; i += colNum)
                _dest.push(_list.slice(i, i + colNum));
        }

        chart.clear();
        const nuRows = [];
        const nuRowsMap = new Map();
        selectedSection = null;
        if (curColNum >= 3) {
            // No sections
            compose(initData.items, nuRows);
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
                    const row = {
                        id: item_type.id,
                        title: `Gruppo ${item_type.id}`,
                        color: item_type.color,
                        rows: filteredRows,
                    };
                    nuRows.push(row);
                    nuRowsMap.set(row.id, row);
                }
            }
        }
        rows = nuRows;
        rowsMap = nuRowsMap;
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

        resized();
    });

    onDestroy(() => {
        window.removeEventListener("resize", resized);
    });

    function increase(data) {
        const id = data.detail.toString();
        order.items[id]++;
        selItem = data.detail;
    }

    function decrease(data) {
        const id = data.detail.toString();
        if (order.items[id] > 0) order.items[id]--;
        selItem = data.detail;
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
            on:click={() => {
                dispatch("reset", null);
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
