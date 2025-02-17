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
        FORMAT_DATE_TIME,
        CALL,
        B2S,
    } from "../utils/Utils.svelte";
    import Divider from "./Divider.svelte";
    import SpaceForFabs from "../app/SpaceForFabs.svelte";

    export let order = {};

    $: candidate = "";
    $: details = null;

    async function submit() {
        const res = await CALL("getBeneficiary", "GET", null, {
            id: candidate,
        });
        if (res.isErr) {
            await ALERT_ERROR(`<p>${res.message}.</p>`);
            return;
        }

        order.beneficiary = candidate;
        details = res.payload;
        details.allowance.forEach((row) => {
            order.allowance[row.item] = row.allowance;
        });
        details.canProceed =
            details.weekIsOk &&
            !details.tooManyOrdersInMonth &&
            !details.tooManyOrdersInWeek;
    }

    function popup() {
        const href = `/api/getBeneficiaryReport?id=${encodeURIComponent(
            order.beneficiary,
        )}`;
        window.open(href, "_blank", "toolbar=0,location=0,menubar=0");
    }
</script>

<div>&nbsp;</div>

<div class="row">
    <div class="col hide-on-small-and-down m3 l4 xl5" />
    <div class="input-field inline col s12 m6 l4 xl2">
        <i class="material-icons prefix">face</i>
        <input id="code" type="text" maxlength="16" bind:value={candidate} />
        <label for="code" class="active">Id beneficiario</label>
        <button class="btn green right" on:click={submit}
            ><i class="material-icons">arrow_forward</i></button
        >
    </div>
    <div class="col hide-on-small-and-down m3 l4 xl5" />
</div>

{#if details}
    {#if details.weekIsOk}
        <Divider />
        <div class="center">
            <h5>Dati del beneficiario</h5>
        </div>
        <div class="row">
            <div class="col hide-on-small-and-down m2 l3 xl4" />
            <div class="input-field col s12 m8 l6 xl4">
                <table>
                    <tr>
                        <td>Profilo</td>
                        <td>{details.profile}</td>
                    </tr>
                    <tr>
                        <td>Ultimo ritiro</td>
                        <td
                            >{#if !!details.lastOrder}N°{details.lastOrder.id} del
                                {FORMAT_DATE_TIME(
                                    details.lastOrder.date,
                                )}{:else}--{/if}</td
                        >
                    </tr>
                    <tr>
                        <td>Già ritirato questa settimana</td>
                        <td>
                            {#if details.tooManyOrdersInWeek}
                                <b class="red-text"
                                    >{B2S(details.tooManyOrdersInWeek)}</b
                                >
                            {:else}
                                {B2S(details.tooManyOrdersInWeek)}
                            {/if}
                        </td>
                    </tr>
                    <tr>
                        <td>Ritiri in questo mese</td>
                        <td>
                            {#if details.tooManyOrdersInMonth}
                                <b class="red-text">{details.ordersInMonth}</b>
                            {:else}
                                {details.ordersInMonth}
                            {/if}
                        </td>
                    </tr>
                    {#if details.canProceed}
                        <tr>
                            <td colspan="2" class="center">
                                <!-- svelte-ignore a11y-click-events-have-key-events -->
                                <!-- svelte-ignore a11y-no-static-element-interactions -->
                                <!-- svelte-ignore a11y-missing-attribute -->
                                <a on:click={popup}>Stampa Scheda</a></td
                            >
                        </tr>
                    {/if}
                </table>
            </div>
            <div class="col hide-on-small-and-down m2 l3 xl4" />
        </div>
    {:else}
        <center
            ><h3>Non è possibile effettuare ritiri questa settimana</h3></center
        >
    {/if}
{/if}

<SpaceForFabs />

<div class="fixed-action-btn">
    <a
        class="btn-floating btn-large green"
        href="#!"
        on:click={() => {
            if (!details) ALERT_ERROR("Beneficiario non valido.");
            else if (!details.canProceed)
                ALERT_ERROR("Non è possibile proseguire.");
            else order.subpage = 2;
        }}
    >
        <i class="material-icons">arrow_forward</i>
    </a>
</div>
