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

    import { createEventDispatcher, onMount } from "svelte";
    import {
        ALERT_INFO,
        FORMAT_TIME,
        LOGOUT,
        MAT_ENABLE_DROPDOWN,
        PAGES,
    } from "../utils/Utils.svelte";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";

    const dispatch = createEventDispatcher();

    let time = "";

    export let initData = null;
    export let checkout = null;
    export let week = null;
    export let beneficiary = null;

    function setTime() {
        time = FORMAT_TIME();
    }

    setTime();

    $: {
        // @ts-ignore
        MAT_ENABLE_DROPDOWN(checkout);
    }

    onMount(() => {
        setInterval(setTime, 1000);
    });

    async function about() {
        await ALERT_INFO(
            `<p><big>FoodHubber ${initData.version}</big></p>
             <hr>
             <p>&copy; 2024- G. Rizzo, col patrocinio di <br/><a href="https://www.emergency.it/" target="_blank">Emergency ONG Onlus</a> 
                e <a href="https://www.aton.com" target="_blank">Aton Società Benefit S.p.A.</a>.</p>
             <p>Questo è software libero:<br/>
             <a href="https://www.gnu.org/licenses/gpl-3.0.en.html" target="_blank">Licenza GPL v3</a>
             , <a href="https://github.com/proofrock/foodhubber" target="_blank">repository dei sorgenti</a></p>`,
        );
    }
</script>

<div class="navbar-fixed">
    <nav>
        <div class="nav-wrapper green darken-1">
            {#if checkout !== null}
                <span class="left hide-on-small-only"
                    >&nbsp;&nbsp;<b>{checkout.cashier}@{checkout.id}</b></span
                >
            {/if}
            <span class="brand-logo center hide-on-small-only"
                ><ConnectionIndicator />FoodHubber{#if !!beneficiary}&nbsp;
                    <small>[Benef.: {beneficiary}]</small>{/if}</span
            >
            <span class="brand-logo left hide-on-med-and-up"
                >&nbsp;<ConnectionIndicator />{#if !!beneficiary}
                    Benef. {beneficiary}{:else}FoodHubber{/if}
                {#if checkout !== null}<sup
                        ><small>&gt;{checkout.cashier}@{checkout.id}</small
                        ></sup
                    >{/if}</span
            >
            {#if checkout !== null}
                <ul id="nav-mobile" class="right">
                    <!-- Dropdown Trigger -->
                    <li class="hide-on-med-and-down">
                        <b
                            >{time}&nbsp;{#if !!week}
                                (Sett. {week}){/if}</b
                        >
                    </li>
                    <li>
                        <a
                            class="dropdown-trigger"
                            href="#!"
                            data-target="pagesMenu"
                            ><i class="material-icons">menu</i></a
                        >
                    </li>
                </ul>
                <!-- Dropdown Structure -->
                <ul id="pagesMenu" class="dropdown-content">
                    <li>
                        <a
                            href="#!"
                            on:click={() => {
                                dispatch("ch_page", PAGES.ORDER);
                            }}>Ordine</a
                        >
                    </li>
                    {#if checkout.can_access_order_list_page}
                        <li>
                            <a
                                href="#!"
                                on:click={() => {
                                    dispatch("ch_page", PAGES.ORDERS_LIST);
                                }}>Lista Ordini</a
                            >
                        </li>
                    {/if}
                    {#if checkout.can_access_stats_page || checkout.can_access_backlogs_page}
                        <li class="divider" />
                    {/if}
                    {#if checkout.can_access_stats_page}
                        <li>
                            <a
                                href="#!"
                                on:click={() => {
                                    dispatch("ch_page", PAGES.STATS);
                                }}>Statistiche</a
                            >
                        </li>
                    {/if}
                    {#if checkout.can_access_stock_page}
                        <li class="divider" />
                        <li>
                            <a
                                href="#!"
                                on:click={() => {
                                    dispatch("ch_page", PAGES.STOCK);
                                }}>Scorte</a
                            >
                        </li>
                    {/if}
                    {#if checkout.can_access_console_page}
                        <li class="divider" />
                        <li>
                            <a
                                href="#!"
                                on:click={() => {
                                    dispatch("ch_page", PAGES.CONSOLE);
                                }}>Sistema</a
                            >
                        </li>
                    {/if}
                    <li class="divider" />
                    <li><a href="#!" on:click={about}>Informazioni</a></li>
                    <li><a href="#!" on:click={LOGOUT}>Esci</a></li>
                </ul>
            {:else}
                <ul class="right">
                    <li>
                        <a href="#!" on:click={about}
                            ><i class="material-icons">info</i></a
                        >
                    </li>
                </ul>
            {/if}
        </div>
    </nav>
</div>
