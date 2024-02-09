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
    import SubpageBeneficiary from "./orderpage/SubpageBeneficiary.svelte";
    import SubpageMenu from "./orderpage/SubpageMenu.svelte";
    import SubpageData from "./orderpage/SubpageData.svelte";

    export let initData;
    export let checkout;
    export let beneficiary;

    $: order = null;

    $: {
        if (order !== null)
            window.sessionStorage.setItem("order", JSON.stringify(order));
    }

    function resetOrder() {
        const _order = {
            subpage: 1,
            beneficiary: "",
            note: "",
            items: {},
            allowance: {},
        };
        for (let i = 0; i < initData.items.length; i++) {
            _order.items[initData.items[i].id.toString()] = 0;
            _order.allowance[initData.items[i].item] = -1;
        }
        order = _order;
    }

    onMount(async () => {
        const _order = window.sessionStorage.getItem("order");
        if (_order === null) resetOrder();
        else {
            order = JSON.parse(_order);
        }
    });

    onDestroy(async () => {
        beneficiary = null;
    });

    $: {
        beneficiary = !!order && !!order.beneficiary ? order.beneficiary : null;
    }
</script>

{#if order !== null && initData !== null}
    {#if order.subpage === 1}
        <SubpageBeneficiary bind:order on:reset={resetOrder} />
    {:else if order.subpage === 2}
        <SubpageMenu bind:order {initData} on:reset={resetOrder} />
    {:else}
        <SubpageData bind:order {initData} {checkout} on:reset={resetOrder} />
    {/if}
{/if}
