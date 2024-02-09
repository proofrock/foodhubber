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

    import { onMount } from "svelte";
    import { createEventDispatcher } from "svelte";
    import { IS_EL_IN_VIEWPORT } from "../utils/Utils.svelte";

    const dispatch = createEventDispatcher();

    export let itm;
    export let selItem;
    export let order;

    $: qty = 0;
    $: {
        if (order.items) qty = order.items[itm.id.toString()];
    }

    function leftclick(e) {
        const cell = document.getElementById(`cell_${itm.id}`);
        const perc = ((e.clientX - cell.offsetLeft) * 100) / cell.offsetWidth;
        if (perc < 20) {
            if (qty <= 0) return;
            dispatch("decrease", itm.id);
        } else {
            if (max >= 0 && qty >= max) return;
            dispatch("increase", itm.id);
        }
    }

    function rightclick() {
        if (qty <= 0) return;
        dispatch("decrease", itm.id);
    }

    $: {
        if (selItem === itm.id) {
            const cell = document.getElementById(`cell_${itm.id}`);
            if (cell && !IS_EL_IN_VIEWPORT(cell))
                cell.scrollIntoView({ block: "end", behavior: "smooth" });
        }
    }

    let max = -1;

    onMount(async () => {
        max = order.allowance.hasOwnProperty(itm.item)
            ? order.allowance[itm.item]
            : -1;
    });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    id="cell_{itm.id}"
    class="card bd sm-cell-smallgutter"
    class:bd-selected={itm.id === selItem}
    style="background-color: {itm.color};"
    on:click|stopPropagation={leftclick}
    on:contextmenu|preventDefault={rightclick}
>
    <div class="card-content sm-cell-smallcell">
        <span class="card-title">
            <b>{itm.item}</b>
            <span class="right">
                {#if qty === 0}
                    --
                {:else}
                    <b>{qty}</b>
                {/if}
            </span>
        </span>
        <span>
            {#if itm.subitem === null || itm.subitem === ""}
                &nbsp;
            {:else}
                {itm.subitem}
            {/if}
        </span>
        <span class="right">
            {#if max >= 0}
                Rim.: {order.allowance[itm.item]}
            {/if}
        </span>
    </div>
</div>

<style>
    .bd {
        border-left-style: solid;
        border-left-width: 4px;
        border-left-color: silver;
    }

    .bd-selected {
        border-left-color: black;
    }

    .sm-cell-smallgutter {
        margin-top: 4px;
        margin-bottom: 4px;
    }

    .sm-cell-smallcell {
        padding-top: 12px;
        padding-bottom: 12px;
    }
</style>
