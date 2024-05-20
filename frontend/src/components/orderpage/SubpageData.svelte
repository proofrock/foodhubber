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
        ALERT_ERROR,
        ALERT_SUCCESS,
        CALL,
        DIALOG_CONFIRM,
    } from "../../components/utils/Utils.svelte";
    import Divider from "./Divider.svelte";
    import SpaceForFabs from "../app/SpaceForFabs.svelte";

    const dispatch = createEventDispatcher();

    export let order = {};
    export let initData = {};
    export let checkout = {};

    $: recap = [[], []];

    $: toSave = {};
    $: errors = [];
    $: warnings = [];
    $: leftovers = []; // after the order, which articles still are not fully ordered (there's still allowance)

    const residuals = {};

    $: {
        const _recap = [];
        for (let i = 0; i < initData.items.length; i++) {
            const q = order.items[initData.items[i].id.toString()];
            const itm = initData.itemsMap.get(initData.items[i].id);
            if (q > 0) {
                _recap.push(
                    `${q}x ${itm.item}${
                        !!itm.subitem ? "/" + itm.subitem : ""
                    }`,
                );
            }
        }
        const half = Math.ceil(_recap.length / 2);
        recap = [_recap.slice(0, half), _recap.slice(half)];
    }

    onMount(async () => {
        const res = await CALL("getBeneficiary", "GET", null, {
            id: order.beneficiary,
        });
        if (res.isErr) {
            // the check will be done on the server anyway
            return;
        }

        res.payload.allowance.forEach((row) => {
            residuals[row.item] = row.residual;
        });
    });

    $: {
        toSave = {
            checkout: checkout.id,
            operator: checkout.cashier,
            beneficiary: order.beneficiary,
            note: order.note.trim() === "" ? null : order.note,
            rows: Object.entries(order.items)
                .filter((e) => e[1] > 0)
                .map((e) => {
                    return { item: parseInt(e[0]), quantity: e[1] };
                }),
        };

        const _errors = [];
        const _warnings = [];
        const _leftovers = [];

        let isEmpty = true;

        // @ts-ignore
        for (let i = 0; i < toSave.rows.length; i++) {
            // @ts-ignore
            const { item: itemId, quantity } = toSave.rows[i];
            if (quantity === 0) continue; // impossible because of line 72. Anyways...

            // Check the residual allowance
            const item = initData.itemsMap.get(itemId);
            if (residuals.hasOwnProperty(item.item)) {
                residuals[item.item] -= quantity;
            }

            isEmpty = false;
        }

        if (isEmpty) _errors.push("Il ritiro è vuoto.");
        else
            for (let i = 0; i < initData.item_categories.length; i++) {
                const key = initData.item_categories[i];
                if (!!residuals[key]) {
                    if (residuals[key] < 0)
                        _errors.push(
                            `La qt. ritirata per la categoria '${key}' è superiore a quanto previsto.`,
                        );
                    else if (residuals[key] > 0)
                        _leftovers.push([key, residuals[key]]);
                }
            }

        if (_leftovers.length > 0)
            _warnings.push(
                "Ci sono ancora articoli che possono essere ritirati, altrimenti saranno 'persi'.",
            );

        if (_errors.length > 0) errors = _errors;
        if (_warnings.length > 0) warnings = _warnings;
        if (_leftovers.length > 0) leftovers = _leftovers;
    }

    async function save() {
        if (errors.length > 0) {
            ALERT_ERROR(
                "Ci sono degli errori. Non è possibile inviare il ritiro.",
            );
            return;
        }

        if (warnings.length > 0) {
            if (
                !(await DIALOG_CONFIRM(
                    "Ci sono degli avvisi.<br/>Vuoi continuare lo stesso?",
                ))
            )
                return;
        } else {
            if (!(await DIALOG_CONFIRM("Vuoi registrare questo ritiro?")))
                return;
        }

        const res = await CALL("putOrder", "PUT", toSave);
        if (res.isErr)
            await ALERT_ERROR(
                `<p>Salvataggio fallito.</p><p>${res.message}.</p>`,
            );
        else {
            let addendum = "";
            if (res.payload.exceeded_stock)
                addendum =
                    "<p>NOTA: almeno un articolo è andato sottoscorta. Verificare.</p>";

            await ALERT_SUCCESS(
                `<p>Ritiro n° ${res.payload.id} registrato correttamente.</p>${addendum}`,
            );
            dispatch("reset", null);
        }
    }
</script>

<div>&nbsp;</div>

<div class="center">
    <h5>Riepilogo</h5>
</div>
<div class="row">
    <div class="col hide-on-small-and-down m1 l2 xl3" />
    <div class="input-field col s12 m5 l4 xl3">
        {#each recap[0] as rrow}
            <p>{rrow}</p>
        {/each}
    </div>
    <div class="input-field col s12 m5 l4 xl3">
        {#each recap[1] as rrow}
            <p>{rrow}</p>
        {/each}
    </div>
    <div class="col hide-on-small-and-down m1 l2 xl3" />
</div>

<Divider />

<div class="center">
    <h5>Dati finali</h5>
</div>
<div class="row">
    <div class="col hide-on-small-and-down m1 l2 xl3" />
    <div class="input-field col s12 m10 l8 xl6">
        <i class="material-icons prefix">edit</i>
        <input id="notes" type="text" maxlength="64" bind:value={order.note} />
        <label for="notes" class="active">Note</label>
    </div>
    <div class="col hide-on-small-and-down m1 l2 xl3" />
</div>

<Divider />

<div class="center">
    <h5>Verifica</h5>
</div>
{#if errors.length === 0 && warnings.length === 0}
    <div class="row">
        <div class="col s12 center">
            <h6><i>Non ci sono segnalazioni.</i></h6>
        </div>
    </div>
{:else}
    <div class="row">
        <div class="col hide-on-small-and-down m1 l2 xl3" />
        <div class="input-field col s12 m5 l4 xl3">
            <h5>Errori</h5>
            <div class="divider" />
            {#each errors as row}
                <div>&nbsp;</div>
                <div class="border-error">{row}</div>
            {/each}
        </div>
        <div class="input-field col s12 m5 l4 xl3">
            <h5>Avvisi</h5>
            <div class="divider" />
            {#each warnings as row}
                <div>&nbsp;</div>
                <div class="border-warning">
                    {row}
                </div>
            {/each}
        </div>
        <div class="col hide-on-small-and-down m1 l2 xl3" />
    </div>
{/if}

<Divider />

{#if leftovers.length > 0}
    <div class="center">
        <h5>Articoli non ritirati (saranno 'persi')</h5>
    </div>
    <div class="row">
        <div class="col hide-on-small-and-down s1 m2 l3 xl4" />
        <div class="input-field col s10 m8 l6 xl4">
            <table>
                <tr>
                    <th>Categoria</th>
                    <th>Qtà residua</th>
                </tr>
                {#each leftovers as leftover (leftover[0])}
                    <tr>
                        <td>{leftover[0]}</td>
                        <td>{leftover[1]}</td>
                    </tr>
                {/each}
            </table>
        </div>
        <div class="col hide-on-small-and-down s1 m2 l3 xl4" />
    </div>
{/if}

<SpaceForFabs />

<div class="fixed-action-btn sm-fab-left">
    <a
        class="btn-floating btn-large yellow darken-3"
        href="#!"
        on:click={() => {
            order.subpage = 2;
        }}
    >
        <i class="large material-icons">arrow_back</i>
    </a>
</div>

<div class="fixed-action-btn">
    <a class="btn-floating btn-large green" href="#!" on:click={save}>
        <i class="material-icons">send</i>
    </a>
</div>

<style>
    .border-warning {
        border-left: 3px solid gold;
        padding-left: 5px;
    }

    .border-error {
        border-left: 3px solid red;
        padding-left: 5px;
    }
</style>
