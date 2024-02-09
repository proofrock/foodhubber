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
        ALERT_ERROR,
        CALL,
        FORMAT_DATE_TIME,
    } from "../components/utils/Utils.svelte";
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
</script>

<div class="row">&nbsp;</div>
{#if data !== null}
    <div class="row">
        <div class="col hide-on-med-and-down m1 l2 xl3" />
        <div class="input-field col s12 m10 l8 xl6">
            <h6>Client Connessi</h6>
            <div class="divider" />
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
{/if}
